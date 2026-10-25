package heatmap

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// HeatmapLogPayload represents the JSON payload ingested by gitmap heatmap log --json.
type HeatmapLogPayload struct {
	Task             string `json:"task,omitempty"`
	TaskSlug         string `json:"taskSlug,omitempty"`
	Repo             string `json:"repo,omitempty"`
	RepoName         string `json:"repoName,omitempty"`
	Action           string `json:"action,omitempty"`
	ActionType       string `json:"actionType,omitempty"`
	File             string `json:"file,omitempty"`
	FilePath         string `json:"filePath,omitempty"`
	Lines            string `json:"lines,omitempty"`
	Code             string `json:"code,omitempty"`
	CodeSnippet      string `json:"codeSnippet,omitempty"`
	Cmd              string `json:"cmd,omitempty"`
	CommandText      string `json:"commandText,omitempty"`
	State            string `json:"state,omitempty"`
	ReasoningState   string `json:"reasoningState,omitempty"`
	Reason           string `json:"reason,omitempty"`
	ReasoningNotes   string `json:"reasoningNotes,omitempty"`
	Expected         string `json:"expected,omitempty"`
	ExpectedOutcome  string `json:"expectedOutcome,omitempty"`
	Actual           string `json:"actual,omitempty"`
	ActualOutcome    string `json:"actualOutcome,omitempty"`
	IsMet            *bool  `json:"met,omitempty"`
	IsExpectationMet *bool  `json:"isExpectationMet,omitempty"`
	Model            string `json:"model,omitempty"`
	ModelName        string `json:"modelName,omitempty"`
	Tool             string `json:"tool,omitempty"`
	ToolName         string `json:"toolName,omitempty"`
	AI               string `json:"ai,omitempty"`
	AiName           string `json:"aiName,omitempty"`
	Parent           string `json:"parent,omitempty"`
	ParentLogId      *int64 `json:"parentLogId,omitempty"`
}

// HeatmapLogJSONInput models single-shot JSON string input for gitmap heatmap log.
type HeatmapLogJSONInput struct {
	Slug        string `json:"slug,omitempty"`
	TaskUuid    string `json:"taskUuid,omitempty"`
	SubtaskName string `json:"subtaskName,omitempty"`
	Status      string `json:"status,omitempty"`
	FilesRead   string `json:"filesRead,omitempty"`
	FilesEdited string `json:"filesEdited,omitempty"`
	ModelName   string `json:"modelName,omitempty"`
	ToolName    string `json:"toolName,omitempty"`
	Roadmap     string `json:"roadmap,omitempty"`
}

// HeatmapResponseEnvelope represents standard JSON response emitted to stdout.
type HeatmapResponseEnvelope struct {
	IsSuccess  bool   `json:"isSuccess"`
	StatusCode int    `json:"statusCode"`
	TaskUuid   string `json:"taskUuid,omitempty"`
	Slug       string `json:"slug,omitempty"`
	ModelName  string `json:"modelName,omitempty"`
	ToolName   string `json:"toolName,omitempty"`
	RecordedAt string `json:"recordedAt"`
	Message    string `json:"message,omitempty"`
}

// ResolveTask returns the effective task identifier.
func (p HeatmapLogPayload) ResolveTask() string {
	if p.Task != "" {
		return p.Task
	}

	return p.TaskSlug
}

// ResolveRepo returns the effective repository name.
func (p HeatmapLogPayload) ResolveRepo() string {
	if p.Repo != "" {
		return p.Repo
	}

	return p.RepoName
}

// ResolveAction returns the effective action type.
func (p HeatmapLogPayload) ResolveAction() string {
	if p.Action != "" {
		return p.Action
	}

	return p.ActionType
}

// ResolveFile returns the effective file path.
func (p HeatmapLogPayload) ResolveFile() string {
	if p.File != "" {
		return p.File
	}

	return p.FilePath
}

// ResolveCode returns the effective code snippet.
func (p HeatmapLogPayload) ResolveCode() string {
	if p.Code != "" {
		return p.Code
	}

	return p.CodeSnippet
}

// ResolveCmd returns the effective command text.
func (p HeatmapLogPayload) ResolveCmd() string {
	if p.Cmd != "" {
		return p.Cmd
	}

	return p.CommandText
}

// ResolveState returns the effective reasoning state.
func (p HeatmapLogPayload) ResolveState() string {
	if p.State != "" {
		return p.State
	}

	return p.ReasoningState
}

// ResolveReason returns the effective reasoning notes.
func (p HeatmapLogPayload) ResolveReason() string {
	if p.Reason != "" {
		return p.Reason
	}

	return p.ReasoningNotes
}

// ResolveExpected returns the expected outcome text.
func (p HeatmapLogPayload) ResolveExpected() string {
	if p.Expected != "" {
		return p.Expected
	}

	return p.ExpectedOutcome
}

// ResolveActual returns the actual outcome text.
func (p HeatmapLogPayload) ResolveActual() string {
	if p.Actual != "" {
		return p.Actual
	}

	return p.ActualOutcome
}

func (p HeatmapLogPayload) resolveModelName() string {
	if p.Model != "" {
		return p.Model
	}

	return p.ModelName
}

// ResolveModel returns the effective model name or AI provenance.
func (p HeatmapLogPayload) ResolveModel() string {
	if name := p.resolveModelName(); name != "" {
		return name
	}

	if p.AI != "" {
		return p.AI
	}

	return p.AiName
}

// ResolveTool returns the effective tool environment name.
func (p HeatmapLogPayload) ResolveTool() string {
	if p.Tool != "" {
		return p.Tool
	}

	return p.ToolName
}

// HasMetExpectation reports whether expectation conditions were met.
func (p HeatmapLogPayload) HasMetExpectation() bool {
	if p.IsMet != nil {
		return *p.IsMet
	}

	if p.IsExpectationMet != nil {
		return *p.IsExpectationMet
	}

	return true
}

// TaskItem represents an architectural plan or telemetry parent task.
type TaskItem struct {
	TaskSequence int    `json:"taskSequence"`
	TaskSlug     string `json:"taskSlug"`
	Title        string `json:"title"`
	TotalActions int    `json:"totalActions"`
	IsActive     bool   `json:"isActive"`
	IsCompleted  bool   `json:"isCompleted"`
}

// TaskInventory holds categorized tasks for status reporting.
type TaskInventory struct {
	Completed []TaskItem `json:"completed"`
	Pending   []TaskItem `json:"pending"`
}

// HeatmapAiConfusion represents an unresolved AI confusion or telemetry alert.
type HeatmapAiConfusion struct {
	ConfusionId   int    `json:"confusionId"`
	TaskSlug      string `json:"taskSlug"`
	FilePath      string `json:"filePath"`
	StartLine     int    `json:"startLine"`
	EndLine       int    `json:"endLine"`
	ReasoningNote string `json:"reasoningNote"`
	IsResolved    bool   `json:"isResolved"`
}

// ConfusionAlertItem is an alias for HeatmapAiConfusion.
type ConfusionAlertItem = HeatmapAiConfusion

// SuggestedCommandsPayload contains actionable one-liner CLI commands.
type SuggestedCommandsPayload struct {
	TreeViewCommand        string `json:"treeViewCommand"`
	SearchErrorCommand     string `json:"searchErrorCommand"`
	AiConfusionSyncCommand string `json:"aiConfusionSyncCommand"`
}

// HeatmapTelemetryPayload contains action metrics and reasoning state breakdown.
type HeatmapTelemetryPayload struct {
	TotalActions         int            `json:"totalActions"`
	ExpectationsMetCount int            `json:"expectationsMetCount"`
	ExpectationsMetPct   float64        `json:"expectationsMetPct"`
	ReasoningStates      map[string]int `json:"reasoningStates"`
}

// HeatmapStatusJSONPayload is emitted when gitmap heatmap status --json is requested.
type HeatmapStatusJSONPayload struct {
	IsSuccess         bool                     `json:"isSuccess"`
	WorkspaceSlug     string                   `json:"workspaceSlug"`
	ActiveTaskSlug    string                   `json:"activeTaskSlug"`
	Telemetry         HeatmapTelemetryPayload  `json:"telemetry"`
	CompletedTasks    []TaskItem               `json:"completedTasks"`
	PendingTasks      []TaskItem               `json:"pendingTasks"`
	ConfusionAlerts   []HeatmapAiConfusion     `json:"confusionAlerts"`
	SuggestedCommands SuggestedCommandsPayload `json:"suggestedCommands"`
}

func sanitizeTaskSlices(inv TaskInventory) ([]TaskItem, []TaskItem) {
	completed := inv.Completed
	if completed == nil {
		completed = []TaskItem{}
	}

	pending := inv.Pending
	if pending == nil {
		pending = []TaskItem{}
	}

	return completed, pending
}

func sanitizeAlerts(alerts []HeatmapAiConfusion) []HeatmapAiConfusion {
	if alerts == nil {
		return []HeatmapAiConfusion{}
	}

	return alerts
}

func sanitizeReasoningStates(states map[string]int) map[string]int {
	if states == nil {
		return make(map[string]int)
	}

	return states
}

func buildTelemetryPayload(status store.HeatmapStatusSummary) HeatmapTelemetryPayload {
	return HeatmapTelemetryPayload{
		TotalActions:         status.TotalActions,
		ExpectationsMetCount: status.ExpectationsMetCount,
		ExpectationsMetPct:   status.ExpectationsMetPct,
		ReasoningStates:      sanitizeReasoningStates(status.ReasoningStateCounts),
	}
}

// BuildHeatmapStatusPayload constructs a validated HeatmapStatusJSONPayload.
func BuildHeatmapStatusPayload(status store.HeatmapStatusSummary, wsSlug, activeSlug string, inv TaskInventory, alerts []HeatmapAiConfusion, cmds SuggestedCommandsPayload) HeatmapStatusJSONPayload {
	completed, pending := sanitizeTaskSlices(inv)

	return HeatmapStatusJSONPayload{
		IsSuccess:         true,
		WorkspaceSlug:     wsSlug,
		ActiveTaskSlug:    activeSlug,
		Telemetry:         buildTelemetryPayload(status),
		CompletedTasks:    completed,
		PendingTasks:      pending,
		ConfusionAlerts:   sanitizeAlerts(alerts),
		SuggestedCommands: cmds,
	}
}
