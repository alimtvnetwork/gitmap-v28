// Package pipelinedb — types.go defines canonical single reusable Result envelopes for pipeline records.
package pipelinedb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// Canonical single reusable monadic containers for pipeline records.
	PipelineRunSliceWrap          = result.WrapSlice[PipelineRunRecord]
	PipelineErrorSliceWrap        = result.WrapSlice[PipelineErrorRecord]
	PipelineCompactErrorSliceWrap = result.WrapSlice[PipelineCompactErrorRecord]
	PipelineRunIdSliceWrap        = result.WrapSlice[uint64]
	PipelineRunIdMapWrap          = result.WrapMap[uint64, bool]
	PipelineDbStatsWrap           = result.Wrap[PipelineDbStats]
	PipelineDatabaseInfoWrap      = result.Wrap[PipelineDatabaseInfo]

	// Backward-compatibility aliases for legacy callers.
	PipelineRunSliceResult          = PipelineRunSliceWrap
	PipelineErrorSliceResult        = PipelineErrorSliceWrap
	PipelineCompactErrorSliceResult = PipelineCompactErrorSliceWrap
	PipelineRunIdSliceResult        = PipelineRunIdSliceWrap
	PipelineRunIdMapResult          = PipelineRunIdMapWrap
	PipelineDbStatsResult           = PipelineDbStatsWrap
	PipelineDatabaseInfoResult      = PipelineDatabaseInfoWrap
)

// PipelineDatabaseInfo encapsulates diagnostic telemetry and metadata for the pipeline SQLite database.
type PipelineDatabaseInfo struct {
	Path         string `json:"path"`
	RelativePath string `json:"relativePath"`
	SizeBytes    int64  `json:"sizeBytes"`
	HumanSize    string `json:"humanSize"`
	IsExisting   bool   `json:"isExisting"`
	TotalRuns    int    `json:"totalRuns"`
	FailedRuns   int    `json:"failedRuns"`
	ErrorCount   int    `json:"errorCount"`
	JobCount     int    `json:"jobCount"`
	SegmentCount int    `json:"segmentCount"`
	LastUpdated  string `json:"lastUpdated"`
}
