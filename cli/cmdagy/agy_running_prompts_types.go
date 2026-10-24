// Package cmdagy — agy_running_prompts_types.go defines types for prompt copy and multi-format exports.
package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy/enums"
)

// PromptCopyOptions encapsulates parameters for cloning batches or individual prompts.
type PromptCopyOptions struct {
	SourceID        string `json:"sourceId"`
	DestinationID   string `json:"destinationId"`
	DestinationName string `json:"destinationName"`
	Namespace       string `json:"namespace"`
	TargetFile      string `json:"targetFile"`
	IsDryRun        bool   `json:"isDryRun"`
	IsForce         bool   `json:"isForce"`
	IsJSON          bool   `json:"isJSON"`
}

// PromptCopyResult represents the outcome of a prompt copy operation.
type PromptCopyResult struct {
	SourceID        string   `json:"sourceId"`
	DestinationID   string   `json:"destinationId"`
	CopiedCount     int      `json:"copiedCount"`
	DestinationName string   `json:"destinationName"`
	Namespace       string   `json:"namespace"`
	IsSuccess       bool     `json:"isSuccess"`
	Message         string   `json:"message"`
	CopiedItemIDs   []string `json:"copiedItemIds,omitempty"`
}

// PromptExportRequest encapsulates options for generating prompt exports.
type PromptExportRequest struct {
	TargetFile   string                       `json:"targetFile"`
	Format       enums.PromptExportFormatType `json:"format"`
	BatchID      string                       `json:"batchId"`
	MaxWords     int                          `json:"maxWords"`
	IsIdempotent bool                         `json:"isIdempotent"`
	IsFull       bool                         `json:"isFull"`
	IsStdout     bool                         `json:"isStdout"`
}

// PromptExportResult represents the summary of an export execution.
type PromptExportResult struct {
	TargetFile   string                       `json:"targetFile"`
	Format       enums.PromptExportFormatType `json:"format"`
	BatchID      string                       `json:"batchId"`
	TotalPrompts int                          `json:"totalPrompts"`
	TotalBytes   int64                        `json:"totalBytes"`
	IsSuccess    bool                         `json:"isSuccess"`
	SqlContent   string                       `json:"sqlContent,omitempty"`
}
