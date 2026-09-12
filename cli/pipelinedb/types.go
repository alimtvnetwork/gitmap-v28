// Package pipelinedb — types.go defines canonical single reusable Result envelopes for pipeline records.
package pipelinedb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// PipelineRunSliceResult is the canonical single reusable result envelope for pipeline run record slices.
	PipelineRunSliceResult = result.ResultSlice[PipelineRunRecord]

	// PipelineErrorSliceResult is the canonical single reusable result envelope for pipeline error record slices.
	PipelineErrorSliceResult = result.ResultSlice[PipelineErrorRecord]

	// PipelineCompactErrorSliceResult is the canonical single reusable result envelope for compact error record slices.
	PipelineCompactErrorSliceResult = result.ResultSlice[PipelineCompactErrorRecord]

	// PipelineRunIdSliceResult is the canonical single reusable result envelope for run ID uint64 slices.
	PipelineRunIdSliceResult = result.ResultSlice[uint64]

	// PipelineRunIdMapResult is the canonical single reusable result envelope for run ID boolean maps.
	PipelineRunIdMapResult = result.ResultMap[uint64, bool]

	// PipelineDbStatsResult is the canonical single reusable result envelope for database statistics.
	PipelineDbStatsResult = result.Result[PipelineDbStats]
)
