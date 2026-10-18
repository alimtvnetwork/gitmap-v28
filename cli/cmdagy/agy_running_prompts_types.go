// Package cmdagy — agy_running_prompts_types.go defines types for prompt copy and multi-format exports.
package cmdagy

import (
	"strings"
)

// PromptExportFormatType defines supported export serialization formats.
type PromptExportFormatType string

const (
	// PromptExportFormatTypeJSON indicates export to JSON format.
	PromptExportFormatTypeJSON PromptExportFormatType = "json"

	// PromptExportFormatTypeSQL indicates export to SQL format.
	PromptExportFormatTypeSQL PromptExportFormatType = "sql"

	// PromptExportFormatTypeDB indicates export to SQLite .db format.
	PromptExportFormatTypeDB PromptExportFormatType = "db"
)

const (
	// PromptExportFormatJSON alias for PromptExportFormatTypeJSON.
	PromptExportFormatJSON = PromptExportFormatTypeJSON

	// PromptExportFormatSQL alias for PromptExportFormatTypeSQL.
	PromptExportFormatSQL = PromptExportFormatTypeSQL

	// PromptExportFormatDB alias for PromptExportFormatTypeDB.
	PromptExportFormatDB = PromptExportFormatTypeDB
)

// ParsePromptExportFormatType parses raw format string into PromptExportFormatType.
func ParsePromptExportFormatType(raw string) PromptExportFormatType {
	lower := strings.ToLower(strings.TrimSpace(raw))

	switch lower {
	case "json", ".json":
		return PromptExportFormatTypeJSON
	case "sql", ".sql":
		return PromptExportFormatTypeSQL
	case "db", ".db", "sqlite", ".sqlite":
		return PromptExportFormatTypeDB
	default:
		return PromptExportFormatType(lower)
	}
}

// IsValid reports whether format is recognized.
func (f PromptExportFormatType) IsValid() bool {
	return f == PromptExportFormatTypeJSON || f == PromptExportFormatTypeSQL || f == PromptExportFormatTypeDB
}

// IsJSON reports whether format is JSON.
func (f PromptExportFormatType) IsJSON() bool {
	return f == PromptExportFormatTypeJSON
}

// IsSQL reports whether format is SQL.
func (f PromptExportFormatType) IsSQL() bool {
	return f == PromptExportFormatTypeSQL
}

// IsDB reports whether format is SQLite DB.
func (f PromptExportFormatType) IsDB() bool {
	return f == PromptExportFormatTypeDB
}

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
	TargetFile   string                 `json:"targetFile"`
	Format       PromptExportFormatType `json:"format"`
	BatchID      string                 `json:"batchId"`
	MaxWords     int                    `json:"maxWords"`
	IsIdempotent bool                   `json:"isIdempotent"`
	IsFull       bool                   `json:"isFull"`
	IsStdout     bool                   `json:"isStdout"`
}

// PromptExportResult represents the summary of an export execution.
type PromptExportResult struct {
	TargetFile   string                 `json:"targetFile"`
	Format       PromptExportFormatType `json:"format"`
	BatchID      string                 `json:"batchId"`
	TotalPrompts int                    `json:"totalPrompts"`
	TotalBytes   int64                  `json:"totalBytes"`
	IsSuccess    bool                   `json:"isSuccess"`
	SqlContent   string                 `json:"sqlContent,omitempty"`
}
