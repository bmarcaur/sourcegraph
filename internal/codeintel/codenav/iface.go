package codenav

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores/gitserver"
	uploads "github.com/sourcegraph/sourcegraph/internal/codeintel/uploads/shared"
	"github.com/sourcegraph/sourcegraph/lib/codeintel/precise"
)

type UploadService interface {
	GetDumpsWithDefinitionsForMonikers(ctx context.Context, monikers []precise.QualifiedMonikerData) (_ []uploads.Dump, err error)
	GetUploadIDsWithReferences(ctx context.Context, orderedMonikers []precise.QualifiedMonikerData, ignoreIDs []int, repositoryID int, commit string, limit int, offset int) (ids []int, recordsScanned int, totalCount int, err error)
	GetDumpsByIDs(ctx context.Context, ids []int) (_ []uploads.Dump, err error)
	InferClosestUploads(ctx context.Context, repositoryID int, commit, path string, exactPath bool, indexer string) (_ []uploads.Dump, err error)
}

type GitserverClient interface {
	CommitsExist(ctx context.Context, commits []gitserver.RepositoryCommit) ([]bool, error)
}
