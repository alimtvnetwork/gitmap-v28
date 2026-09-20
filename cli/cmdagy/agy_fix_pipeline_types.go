package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
)

// SentAgyErrorRecord records metadata for a sent pipeline error batch.
type SentAgyErrorRecord struct {
	Signature string `json:"signature"`
	Repo      string `json:"repo"`
	RunID     uint64 `json:"runId"`
	SHA       string `json:"sha"`
	ErrorHash string `json:"errorHash"`
	SentAt    string `json:"sentAt"`
	SentCount int    `json:"sentCount"`
}

// SentAgyErrorsStore tracks historical sent pipeline error signatures for deduplication.
type SentAgyErrorsStore struct {
	Records map[string]SentAgyErrorRecord `json:"records"`
}

// AgyPromptQueueEntry records prompt queue items for Antigravity sessions.
type AgyPromptQueueEntry struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Prompt    string `json:"prompt"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// AgyPromptQueueFile stores active and queued prompts for Antigravity verification.
type AgyPromptQueueFile struct {
	Active    *AgyPromptQueueEntry  `json:"active,omitempty"`
	Queued    []AgyPromptQueueEntry `json:"queued,omitempty"`
	UpdatedAt string                `json:"updatedAt"`
}

// AgyFixOptions captures all flags for pipeline fix errors agy.
type AgyFixOptions struct {
	Repo          string
	IsForce       bool
	IsDetailed    bool
	IsNoRelease   bool
	CustomPrompt  string
	IsNoClipboard bool
	OutputFile    string
	IsDryRun      bool
	IsAll         bool
	ProjectsCount int
	Limit         int
	IsResetBatch  bool
	IsNoInject    bool
}

// AgyFixDispatchParams packages parameters for executing prompt dispatch.
type AgyFixDispatchParams struct {
	Opts        AgyFixOptions
	StorePath   string
	Sig         string
	ErrHash     string
	Payload     cmdpipeline.PipelineErrorLogsPayload
	ErrorReport string
	HasFailures bool
	TargetDir   string
}

// AgyInjectionModeType indicates whether injection was routed via CLI, IDE, or skipped.
type AgyInjectionModeType string

const (
	AgyInjectionModeNone   AgyInjectionModeType = "none"
	AgyInjectionModeCLI    AgyInjectionModeType = "cli"
	AgyInjectionModeIDE    AgyInjectionModeType = "ide"
	AgyInjectionModeQueued AgyInjectionModeType = "queued"
)

// AgyConvStatusType represents the execution status of an active conversation.
type AgyConvStatusType string

const (
	AgyConvStatusIdle    AgyConvStatusType = "idle"
	AgyConvStatusRunning AgyConvStatusType = "running"
	AgyConvStatusUnknown AgyConvStatusType = "unknown"
)

// AgyConversationExecutionState represents the detected status of a conversation.
type AgyConversationExecutionState struct {
	ConvID         string            `json:"convId"`
	Status         AgyConvStatusType `json:"status"`
	TranscriptPath string            `json:"transcriptPath"`
}

// AgyInjectionResult encapsulates the outcome of an Antigravity injection task.
type AgyInjectionResult struct {
	IsSuccess  bool                 `json:"IsSuccess"`
	Mode       AgyInjectionModeType `json:"Mode"`
	PID        int                  `json:"PID"`
	Message    string               `json:"Message"`
	PromptPath string               `json:"PromptPath"`
	RepoDir    string               `json:"RepoDir"`
}

// AgyAssembledPromptPayload bundles assembled prompt components for pipeline fix dispatch.
type AgyAssembledPromptPayload struct {
	Primary       string `json:"Primary"`
	Followup      string `json:"Followup"`
	PromptContent string `json:"PromptContent"`
	PromptSource  string `json:"PromptSource"`
}
