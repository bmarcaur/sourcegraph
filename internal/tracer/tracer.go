// Package opentracer initializes distributed tracing and log15 behavior. It also updates distributed
// tracing behavior in response to changes in site configuration. When the Init function of this
// package is invoked, opentracing.SetGlobalTracer is called (and subsequently called again after
// every Sourcegraph site configuration change). Importing programs should not invoke
// opentracing.SetGlobalTracer anywhere else.
package tracer

import (
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	"go.uber.org/automaxprocs/maxprocs"
	ddopentracing "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/trace/ot"
	"github.com/sourcegraph/sourcegraph/internal/version"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func init() {
	// Tune GOMAXPROCS for kubernetes. All our binaries import this package,
	// so we tune for all of them.
	//
	// TODO it is surprising that we do this here. We should create a standard
	// import for sourcegraph binaries which would have less surprising
	// behaviour.
	if _, err := maxprocs.Set(); err != nil {
		log15.Error("automaxprocs failed", "error", err)
	}
}

// Options control the behavior of a Tracer
type Options struct {
	Tracer
	externalURL string
	debug       bool
	// these values are not configurable at runtime
	serviceName string
	version     string
	env         string
}

type Tracer string

const (
	none    Tracer = ""
	datadog Tracer = "datadog"
	jaeg    Tracer = "opentracing"
)

// Option follows this idiom https://github.com/tmrts/go-patterns/blob/master/idiom/functional-options.md
type Option func(*Options)

func WithTracer(t Tracer) Option {
	if t != none && t != datadog && t != jaeg {
		log15.Info("impossible tracer ")
		panic("impossible tracer chosen")
	}
	return func(o *Options) {
		o.Tracer = t
	}
}

func WithService(name string) Option {
	return func(o *Options) {
		o.serviceName = name
	}
}

func WithDebug(enabled bool) Option {
	return func(o *Options) {
		o.debug = enabled
	}
}

func WithExternalURL(url string) Option {
	return func(o *Options) {
		o.externalURL = url
	}
}

func WithServiceVersion(v string) Option {
	return func(o *Options) {
		o.version = v
	}
}

// Init should be called from the main function of service
func Init(c conftypes.WatchableSiteConfig, options ...Option) {
	opts := &Options{}
	for _, setter := range options {
		setter(opts)
	}
	if opts.serviceName == "" {
		opts.serviceName = env.MyName
	}
	if opts.version == "" {
		if version.IsDev(version.Version()) {
			opts.env = "dev"
		}
		opts.version = version.Version()
	}

	initTracer(opts, c)
}

// initTracer is a helper that should be called exactly once (from Init).
func initTracer(opts *Options, c conftypes.WatchableSiteConfig) {
	globalTracer := newSwitchableTracer()
	opentracing.SetGlobalTracer(globalTracer)

	// initial tracks if it's our first run of conf.Watch. This is used to
	// prevent logging "changes" when it's the first run.
	initial := true

	// Initially everything is disabled since we haven't read conf yet.
	oldOpts := Options{
		Tracer:      none,
		debug:       false,
		serviceName: opts.serviceName,
		version:     opts.version,
		env:         opts.env,
	}

	// Watch loop
	go c.Watch(func() {
		siteConfig := c.SiteConfig()

		// Set sampling strategy
		samplingStrategy := ot.TraceNone
		shouldLog := false
		tracerType := none
		if tracingConfig := siteConfig.ObservabilityTracing; tracingConfig != nil {
			tracerType = Tracer(tracingConfig.Type)
			switch tracingConfig.Sampling {
			case "all":
				samplingStrategy = ot.TraceAll
			case "selective":
				samplingStrategy = ot.TraceSelective
			}
			shouldLog = tracingConfig.Debug
		}
		if tracePolicy := ot.GetTracePolicy(); tracePolicy != samplingStrategy && !initial {
			log15.Info("opentracing: TracePolicy", "oldValue", tracePolicy, "newValue", samplingStrategy)
		}
		initial = false
		ot.SetTracePolicy(samplingStrategy)

		opts := Options{
			externalURL: siteConfig.ExternalURL,
			Tracer:      tracerType,
			debug:       shouldLog,
			serviceName: opts.serviceName,
			version:     opts.version,
			env:         opts.env,
		}

		if opts == oldOpts {
			// Nothing changed
			return
		}

		oldOpts = opts

		t, closer, err := newTracer(&opts)
		if err != nil {
			log15.Warn("Could not initialize tracer", "error", err.Error())
			return
		}
		globalTracer.set(t, closer, opts.debug)
	})
}

// TODO Use openTelemetry https://github.com/sourcegraph/sourcegraph/issues/27386
func newTracer(opts *Options) (opentracing.Tracer, io.Closer, error) {
	if opts.Tracer == none {
		log15.Info("tracing disabled")
		ddtracer.Stop()
		return opentracing.NoopTracer{}, nil, nil
	}
	if opts.Tracer == datadog {
		log15.Info("datadog: tracing enabled")
		tracer := ddopentracing.New(ddtracer.WithService(opts.serviceName),
			ddtracer.WithDebugMode(opts.debug),
			ddtracer.WithServiceVersion(opts.version), ddtracer.WithEnv(opts.env))
		return tracer, nil, nil
	}
	ddtracer.Stop()

	log15.Info("opentracing: Jaeger enabled")
	cfg, err := jaegercfg.FromEnv()
	cfg.ServiceName = opts.serviceName
	if err != nil {
		return nil, nil, errors.Wrap(err, "jaegercfg.FromEnv failed")
	}
	cfg.Tags = append(cfg.Tags, opentracing.Tag{Key: "service.version", Value: version.Version()})
	if reflect.DeepEqual(cfg.Sampler, &jaegercfg.SamplerConfig{}) {
		// Default sampler configuration for when it is not specified via
		// JAEGER_SAMPLER_* env vars. In most cases, this is sufficient
		// enough to connect Sourcegraph to Jaeger without any env vars.
		cfg.Sampler.Type = jaeger.SamplerTypeConst
		cfg.Sampler.Param = 1
	}
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(log15Logger{}),
		jaegercfg.Metrics(jaegermetrics.NullFactory),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "jaegercfg.NewTracer failed")
	}

	return tracer, closer, nil
}

type log15Logger struct{}

func (l log15Logger) Error(msg string) { log15.Error(msg) }

func (l log15Logger) Infof(msg string, args ...interface{}) {
	log15.Info(fmt.Sprintf(msg, args...))
}

// move to OpenTelemetry https://github.com/sourcegraph/sourcegraph/issues/27386
// switchableTracer implements opentracing.Tracer. The underlying opentracer used is switchable (set via
// the `set` method).
type switchableTracer struct {
	mu           sync.RWMutex
	opentracer   opentracing.Tracer
	tracerCloser io.Closer
	log          bool
}

// move to OpenTelemetry https://github.com/sourcegraph/sourcegraph/issues/27386
func newSwitchableTracer() *switchableTracer {
	return &switchableTracer{opentracer: opentracing.NoopTracer{}}
}

func (t *switchableTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.log {
		log15.Info("opentracing: StartSpan", "operationName", operationName, "opentracer", fmt.Sprintf("%T", t.opentracer))
	}
	return t.opentracer.StartSpan(operationName, opts...)
}

func (t *switchableTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.log {
		log15.Info("opentracing: Inject", "opentracer", fmt.Sprintf("%T", t.opentracer))
	}
	return t.opentracer.Inject(sm, format, carrier)
}

func (t *switchableTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.log {
		log15.Info("opentracing: Extract", "tracer", fmt.Sprintf("%T", t.opentracer))
	}
	return t.opentracer.Extract(format, carrier)
}

func (t *switchableTracer) set(tracer opentracing.Tracer, tracerCloser io.Closer, log bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if tc := t.tracerCloser; tc != nil {
		// Close the old tracerCloser outside the critical zone
		go tc.Close()
	}

	t.tracerCloser = tracerCloser
	t.opentracer = tracer
	t.log = log
}
