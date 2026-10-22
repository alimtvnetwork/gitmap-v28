// Package pipelinedb — types.go defines canonical single reusable Result envelopes for pipeline records.
package pipelinedb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// PipelineRunSliceResult is the canonical single reusable result envelope for pipeline run record slices.
	PipelineRunSliceResult = result.WrapSlice[PipelineRunRecord]

	// PipelineErrorSliceResult is the canonical single reusable result envelope for pipeline error record slices.
	PipelineErrorSliceResult = result.WrapSlice[PipelineErrorRecord]

	// PipelineCompactErrorSliceResult is the canonical single reusable result envelope for compact error record slices.
	PipelineCompactErrorSliceResult = result.WrapSlice[PipelineCompactErrorRecord]

	// PipelineRunIdSliceResult is the canonical single reusable result envelope for run ID uint64 slices.
	PipelineRunIdSliceResult = result.WrapSlice[uint64]

	// PipelineRunIdMapResult is the canonical single reusable result envelope for run ID boolean maps.
	PipelineRunIdMapResult = result.WrapMap[uint64, bool]

	// PipelineDbStatsResult is the canonical single reusable result envelope for database statistics.
	PipelineDbStatsResult = result.Wrap[PipelineDbStats]

	// PipelineDatabaseInfoResult is the canonical single reusable result envelope for database telemetry.
	PipelineDatabaseInfoResult = result.Wrap[PipelineDatabaseInfo]
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
