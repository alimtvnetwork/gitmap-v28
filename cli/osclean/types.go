package osclean

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// CleanOptions configures system temporary file cleaning.
type CleanOptions struct {
	IsDryRun   bool
	IsVerbose  bool
	IsTempOnly bool
}

// CleanStats records metrics of deleted files and freed bytes.
type CleanStats struct {
	RemovedFilesCount int
	RemovedDirsCount  int
	FreedBytes        int64
	HasErrors         bool
}

// CleanResult encapsulates CleanStats result wrapper.
type CleanResult = result.Result[CleanStats]
