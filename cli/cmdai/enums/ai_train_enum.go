package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// LlmTrainingFormatType defines output format for LLM reasoning ingestion.
type LlmTrainingFormatType string

const (
	LlmFormatMarkdown LlmTrainingFormatType = "markdown"
	LlmFormatJsonl    LlmTrainingFormatType = "jsonl"
)

// ExportFormatType defines supported packaging formats for export.
type ExportFormatType string

const (
	ExportFormatJson ExportFormatType = "json"
	ExportFormatZip  ExportFormatType = "zip"
	ExportFormatDb   ExportFormatType = "db"
)

// IsValid reports whether the training format matches a known format.
func (f LlmTrainingFormatType) IsValid() bool {
	switch f {
	case LlmFormatMarkdown, LlmFormatJsonl:
		return true
	default:
		return false
	}
}

// String returns the string representation of the training format.
func (f LlmTrainingFormatType) String() string {
	return string(f)
}

// ParseLlmTrainingFormatType parses a raw string into a validated LlmTrainingFormatType.
func ParseLlmTrainingFormatType(s string) (LlmTrainingFormatType, *apperror.AppError) {
	val := LlmTrainingFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseLlmTrainingFormatType", "E1000", fmt.Sprintf("invalid LLM training format: %q", s)).WithContext("value", s)
}

// AllLlmTrainingFormatTypes returns all valid LlmTrainingFormatType variants.
func AllLlmTrainingFormatTypes() []LlmTrainingFormatType {
	return []LlmTrainingFormatType{LlmFormatMarkdown, LlmFormatJsonl}
}

// IsMarkdown reports whether the format is markdown.
func (f LlmTrainingFormatType) IsMarkdown() bool { return f == LlmFormatMarkdown }

// IsJsonl reports whether the format is jsonl.
func (f LlmTrainingFormatType) IsJsonl() bool { return f == LlmFormatJsonl }

// IsValid reports whether the export format matches a known format.
func (f ExportFormatType) IsValid() bool {
	switch f {
	case ExportFormatJson, ExportFormatZip, ExportFormatDb:
		return true
	default:
		return false
	}
}

// String returns the string representation of the export format.
func (f ExportFormatType) String() string {
	return string(f)
}

// ParseExportFormatType parses a raw string into a validated ExportFormatType.
func ParseExportFormatType(s string) (ExportFormatType, *apperror.AppError) {
	val := ExportFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseExportFormatType", "E1000", fmt.Sprintf("invalid export format: %q", s)).WithContext("value", s)
}

// AllExportFormatTypes returns all valid ExportFormatType variants.
func AllExportFormatTypes() []ExportFormatType {
	return []ExportFormatType{ExportFormatJson, ExportFormatZip, ExportFormatDb}
}

// IsJson reports whether the format is json.
func (f ExportFormatType) IsJson() bool { return f == ExportFormatJson }

// IsZip reports whether the format is zip.
func (f ExportFormatType) IsZip() bool { return f == ExportFormatZip }

// IsDb reports whether the format is db.
func (f ExportFormatType) IsDb() bool { return f == ExportFormatDb }
