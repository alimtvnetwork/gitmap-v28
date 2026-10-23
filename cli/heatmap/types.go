package heatmap

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
