// Package store — types.go centralizes domain models, supporting structs, and Result container aliases for store operations.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ScanRecordSlice is a slice of ScanRecord.
	ScanRecordSlice = []model.ScanRecord

	// ScanRecordWrap is the canonical single reusable result container for ScanRecord.
	ScanRecordWrap = result.Wrap[model.ScanRecord]

	// ScanRecordSliceWrap is the canonical single reusable collection container for ScanRecord.
	ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]

	// ScanRecordMapWrap is the canonical single reusable keyed lookup container for ScanRecord.
	ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]

	// RepoPathWrap is the canonical single reusable result container for repository paths.
	RepoPathWrap = result.Wrap[string]

	// HeatmapActionLogRecordSliceWrap is the canonical collection container for HeatmapActionLogRecord.
	HeatmapActionLogRecordSliceWrap = result.WrapSlice[HeatmapActionLogRecord]

	// HeatmapExecutionTreeNodeSliceWrap is the canonical collection container for HeatmapExecutionTreeNode.
	HeatmapExecutionTreeNodeSliceWrap = result.WrapSlice[HeatmapExecutionTreeNode]

	// HeatmapStatusSummaryWrap is the canonical result container for HeatmapStatusSummary.
	HeatmapStatusSummaryWrap = result.Wrap[HeatmapStatusSummary]

	// HeatmapTaskSummaryWrap is the canonical result container for HeatmapTaskSummary.
	HeatmapTaskSummaryWrap = result.Wrap[HeatmapTaskSummary]

	// FileTouchCountSliceWrap is the canonical collection container for FileTouchCount.
	FileTouchCountSliceWrap = result.WrapSlice[FileTouchCount]

	// DuplicateRepoGroupSliceWrap is the canonical collection container for DuplicateRepoGroup.
	DuplicateRepoGroupSliceWrap = result.WrapSlice[DuplicateRepoGroup]

	// RowsAffectedWrap is the canonical result container for affected row counts.
	RowsAffectedWrap = result.Wrap[int64]

	// RepoIDWrap is the canonical result container for repository IDs.
	RepoIDWrap = result.Wrap[int64]

	// RepoCountWrap is the canonical result container for integer repo counts.
	RepoCountWrap = result.Wrap[int]
)

// HeatmapActionLogRecord models an action log entry for heatmap rendering.
type HeatmapActionLogRecord struct {
	ID        int64  `json:"id"`
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	Path      string `json:"path"`
}

// HeatmapExecutionTreeNode models a node in the execution tree.
type HeatmapExecutionTreeNode struct {
	NodeID   string `json:"nodeId"`
	ParentID string `json:"parentId"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

// HeatmapStatusSummary models status aggregation metrics.
type HeatmapStatusSummary struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// HeatmapTaskSummary models task summary metrics.
type HeatmapTaskSummary struct {
	TaskID    string `json:"taskId"`
	State     string `json:"state"`
	DurationMs int64 `json:"durationMs"`
}

// FileTouchCount models file touch frequency.
type FileTouchCount struct {
	FilePath   string `json:"filePath"`
	TouchCount int    `json:"touchCount"`
}

// DuplicateRepoGroup models a group of duplicate repositories.
type DuplicateRepoGroup struct {
	CanonicalPath string   `json:"canonicalPath"`
	DuplicatePaths []string `json:"duplicatePaths"`
}
