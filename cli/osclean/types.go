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

// CleanWrap encapsulates CleanStats result wrapper.
type CleanWrap = result.Wrap[CleanStats]

// CleanResult is a backward-compatible alias for CleanWrap.
type CleanResult = CleanWrap

// DevCleanSummaryWrap encapsulates DevCleanSummary result wrapper.
type DevCleanSummaryWrap = result.Wrap[DevCleanSummary]

// TerminalCleanSummaryWrap encapsulates TerminalCleanSummary result wrapper.
type TerminalCleanSummaryWrap = result.Wrap[TerminalCleanSummary]
