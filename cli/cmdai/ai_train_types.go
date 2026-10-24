package cmdai

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai/enums"
)

// LlmTrainingFormatType is the canonical LLM training format type, defined in enums.
type LlmTrainingFormatType = enums.LlmTrainingFormatType

const (
	LlmFormatMarkdown = enums.LlmFormatMarkdown
	LlmFormatJsonl    = enums.LlmFormatJsonl
)

// LlmMessageItem represents a single role turn in conversational training schemas.
type LlmMessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LlmConversationRecord represents a JSONL fine-tuning record.
type LlmConversationRecord struct {
	Messages []LlmMessageItem `json:"messages"`
}

// TaskDecisionChain encapsulates a task session with its chronological line decisions.
type TaskDecisionChain struct {
	Task  AiAnalysisTask   `json:"task"`
	Lines []AiAnalysisLine `json:"lines"`
}

// LlmTrainOptions configures the reasoning extraction and formatting.
type LlmTrainOptions struct {
	TaskId   string
	Format   LlmTrainingFormatType
	Limit    int
	Output   string
	RepoRoot string
}

// ExportFormatType is the canonical export format type, defined in enums.
type ExportFormatType = enums.ExportFormatType

const (
	ExportFormatJson = enums.ExportFormatJson
	ExportFormatZip  = enums.ExportFormatZip
	ExportFormatDb   = enums.ExportFormatDb
)

// AiExportBundle represents a serialized collection of tasks and lines.
type AiExportBundle struct {
	ExportedAt time.Time        `json:"exportedAt"`
	RepoSlug   string           `json:"repoSlug"`
	Tasks      []AiAnalysisTask `json:"tasks"`
	Lines      []AiAnalysisLine `json:"lines"`
}

// AiExportOptions configures export operations.
type AiExportOptions struct {
	TaskId   string
	Format   ExportFormatType
	Output   string
	RepoRoot string
}

// AiImportOptions configures import operations.
type AiImportOptions struct {
	InputPath string
	RepoRoot  string
}

// AiClearOptions configures pruning and truncation.
type AiClearOptions struct {
	TaskId    string
	OlderThan time.Duration
	All       bool
	Force     bool
	RepoRoot  string
}
