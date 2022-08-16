// Code generated by go-mockgen 1.3.4; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package conf

import (
	"context"
	"sync"

	conftypes "github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
)

// MockConfigurationSource is a mock implementation of the
// ConfigurationSource interface (from the package
// github.com/sourcegraph/sourcegraph/internal/conf) used for unit testing.
type MockConfigurationSource struct {
	// ReadFunc is an instance of a mock function object controlling the
	// behavior of the method Read.
	ReadFunc *ConfigurationSourceReadFunc
	// WriteFunc is an instance of a mock function object controlling the
	// behavior of the method Write.
	WriteFunc *ConfigurationSourceWriteFunc
}

// NewMockConfigurationSource creates a new mock of the ConfigurationSource
// interface. All methods return zero values for all results, unless
// overwritten.
func NewMockConfigurationSource() *MockConfigurationSource {
	return &MockConfigurationSource{
		ReadFunc: &ConfigurationSourceReadFunc{
			defaultHook: func(context.Context) (r0 conftypes.RawUnified, r1 error) {
				return
			},
		},
		WriteFunc: &ConfigurationSourceWriteFunc{
			defaultHook: func(context.Context, conftypes.RawUnified) (r0 error) {
				return
			},
		},
	}
}

// NewStrictMockConfigurationSource creates a new mock of the
// ConfigurationSource interface. All methods panic on invocation, unless
// overwritten.
func NewStrictMockConfigurationSource() *MockConfigurationSource {
	return &MockConfigurationSource{
		ReadFunc: &ConfigurationSourceReadFunc{
			defaultHook: func(context.Context) (conftypes.RawUnified, error) {
				panic("unexpected invocation of MockConfigurationSource.Read")
			},
		},
		WriteFunc: &ConfigurationSourceWriteFunc{
			defaultHook: func(context.Context, conftypes.RawUnified) error {
				panic("unexpected invocation of MockConfigurationSource.Write")
			},
		},
	}
}

// NewMockConfigurationSourceFrom creates a new mock of the
// MockConfigurationSource interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockConfigurationSourceFrom(i ConfigurationSource) *MockConfigurationSource {
	return &MockConfigurationSource{
		ReadFunc: &ConfigurationSourceReadFunc{
			defaultHook: i.Read,
		},
		WriteFunc: &ConfigurationSourceWriteFunc{
			defaultHook: i.Write,
		},
	}
}

// ConfigurationSourceReadFunc describes the behavior when the Read method
// of the parent MockConfigurationSource instance is invoked.
type ConfigurationSourceReadFunc struct {
	defaultHook func(context.Context) (conftypes.RawUnified, error)
	hooks       []func(context.Context) (conftypes.RawUnified, error)
	history     []ConfigurationSourceReadFuncCall
	mutex       sync.Mutex
}

// Read delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockConfigurationSource) Read(v0 context.Context) (conftypes.RawUnified, error) {
	r0, r1 := m.ReadFunc.nextHook()(v0)
	m.ReadFunc.appendCall(ConfigurationSourceReadFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Read method of the
// parent MockConfigurationSource instance is invoked and the hook queue is
// empty.
func (f *ConfigurationSourceReadFunc) SetDefaultHook(hook func(context.Context) (conftypes.RawUnified, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Read method of the parent MockConfigurationSource instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ConfigurationSourceReadFunc) PushHook(hook func(context.Context) (conftypes.RawUnified, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ConfigurationSourceReadFunc) SetDefaultReturn(r0 conftypes.RawUnified, r1 error) {
	f.SetDefaultHook(func(context.Context) (conftypes.RawUnified, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ConfigurationSourceReadFunc) PushReturn(r0 conftypes.RawUnified, r1 error) {
	f.PushHook(func(context.Context) (conftypes.RawUnified, error) {
		return r0, r1
	})
}

func (f *ConfigurationSourceReadFunc) nextHook() func(context.Context) (conftypes.RawUnified, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfigurationSourceReadFunc) appendCall(r0 ConfigurationSourceReadFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfigurationSourceReadFuncCall objects
// describing the invocations of this function.
func (f *ConfigurationSourceReadFunc) History() []ConfigurationSourceReadFuncCall {
	f.mutex.Lock()
	history := make([]ConfigurationSourceReadFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfigurationSourceReadFuncCall is an object that describes an invocation
// of method Read on an instance of MockConfigurationSource.
type ConfigurationSourceReadFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 conftypes.RawUnified
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfigurationSourceReadFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfigurationSourceReadFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ConfigurationSourceWriteFunc describes the behavior when the Write method
// of the parent MockConfigurationSource instance is invoked.
type ConfigurationSourceWriteFunc struct {
	defaultHook func(context.Context, conftypes.RawUnified) error
	hooks       []func(context.Context, conftypes.RawUnified) error
	history     []ConfigurationSourceWriteFuncCall
	mutex       sync.Mutex
}

// Write delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockConfigurationSource) Write(v0 context.Context, v1 conftypes.RawUnified) error {
	r0 := m.WriteFunc.nextHook()(v0, v1)
	m.WriteFunc.appendCall(ConfigurationSourceWriteFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Write method of the
// parent MockConfigurationSource instance is invoked and the hook queue is
// empty.
func (f *ConfigurationSourceWriteFunc) SetDefaultHook(hook func(context.Context, conftypes.RawUnified) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Write method of the parent MockConfigurationSource instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ConfigurationSourceWriteFunc) PushHook(hook func(context.Context, conftypes.RawUnified) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ConfigurationSourceWriteFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, conftypes.RawUnified) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ConfigurationSourceWriteFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, conftypes.RawUnified) error {
		return r0
	})
}

func (f *ConfigurationSourceWriteFunc) nextHook() func(context.Context, conftypes.RawUnified) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfigurationSourceWriteFunc) appendCall(r0 ConfigurationSourceWriteFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfigurationSourceWriteFuncCall objects
// describing the invocations of this function.
func (f *ConfigurationSourceWriteFunc) History() []ConfigurationSourceWriteFuncCall {
	f.mutex.Lock()
	history := make([]ConfigurationSourceWriteFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfigurationSourceWriteFuncCall is an object that describes an
// invocation of method Write on an instance of MockConfigurationSource.
type ConfigurationSourceWriteFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 conftypes.RawUnified
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfigurationSourceWriteFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfigurationSourceWriteFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
