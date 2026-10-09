package types

// AgentActionType represents the category of action performed by an autonomous agent.
type AgentActionType string

const (
	ActionSearch   AgentActionType = "SEARCH"
	ActionRead     AgentActionType = "READ"
	ActionWrite    AgentActionType = "WRITE"
	ActionEdit     AgentActionType = "EDIT"
	ActionExec     AgentActionType = "EXEC"
	ActionLint     AgentActionType = "LINT"
	ActionCheck    AgentActionType = "CHECK"
	ActionClaim    AgentActionType = "CLAIM"
	ActionComplete AgentActionType = "COMPLETE"
	ActionFail     AgentActionType = "FAIL"
	ActionCrash    AgentActionType = "CRASH"
)

// String returns the string representation of AgentActionType.
func (a AgentActionType) String() string {
	return string(a)
}

// IsValid checks if the action type is one of the supported enum values.
func (a AgentActionType) IsValid() bool {
	switch a {
	case ActionSearch, ActionRead, ActionWrite, ActionEdit, ActionExec,
		ActionLint, ActionCheck, ActionClaim, ActionComplete, ActionFail, ActionCrash:
		return true
	default:
		return false
	}
}

// ParentTask represents a root or parent autonomous agent execution task.
type ParentTask struct {
	ParentTaskId      string `json:"parentTaskId"`
	TaskSlug          string `json:"taskSlug"`
	TaskName          string `json:"taskName"`
	RunDirectory      string `json:"runDirectory"`
	RootDbPath        string `json:"rootDbPath"`
	Status            string `json:"status"`
	TotalStepsBudget  int    `json:"totalStepsBudget"`
	CompletedSteps    int    `json:"completedSteps"`
	SpawnedAgentCount int    `json:"spawnedAgentCount"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

// IsFinished reports if the parent task has completed.
func (p ParentTask) IsFinished() bool {
	return p.Status == "COMPLETED" || p.Status == "DONE"
}

// Subtask represents a discrete, atomic work unit delegated to an autonomous worker agent.
type Subtask struct {
	SubtaskId         string `json:"subtaskId"`
	ParentTaskId      string `json:"parentTaskId"`
	TaskCode          string `json:"taskCode"`
	TaskSlug          string `json:"taskSlug"`
	Title             string `json:"title"`
	AssignedAgentRole string `json:"assignedAgentRole"`
	OwnedFilesJson    string `json:"ownedFilesJson"`
	Status            string `json:"status"`
	Evidence          string `json:"evidence"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

// IsPending reports if the subtask is awaiting claim.
func (s Subtask) IsPending() bool {
	return s.Status == "PENDING"
}

// IsDone reports if the subtask is completed.
func (s Subtask) IsDone() bool {
	return s.Status == "DONE" || s.Status == "COMPLETED"
}

// HasEvidence reports if the subtask has verification evidence recorded.
func (s Subtask) HasEvidence() bool {
	return len(s.Evidence) > 0
}

// AgentActionLog captures granular telemetry of file touches, commands, and operations.
type AgentActionLog struct {
	ActionLogId    int64           `json:"actionLogId"`
	SubtaskId      string          `json:"subtaskId"`
	AgentRole      string          `json:"agentRole"`
	ActionType     AgentActionType `json:"actionType"`
	TargetFile     string          `json:"targetFile"`
	StartLine      int             `json:"startLine"`
	EndLine        int             `json:"endLine"`
	QueryOrCommand string          `json:"queryOrCommand"`
	ActionDetails  string          `json:"actionDetails"`
	DurationMs     int64           `json:"durationMs"`
	Status         string          `json:"status"`
	ErrorMessage   string          `json:"errorMessage"`
	CreatedAt      string          `json:"createdAt"`
}

// AgentRegistryEntry registers active worker agents and their split database locations.
type AgentRegistryEntry struct {
	AgentId         string `json:"agentId"`
	ParentTaskId    string `json:"parentTaskId"`
	AgentRole       string `json:"agentRole"`
	AgentSlug       string `json:"agentSlug"`
	SplitDbPath     string `json:"splitDbPath"`
	Status          string `json:"status"`
	LastHeartbeatAt string `json:"lastHeartbeatAt"`
}

// CrashDiagnosis encapsulates autopsy details when an agent or subtask crashes.
type CrashDiagnosis struct {
	AgentRole       string          `json:"agentRole"`
	SubtaskId       string          `json:"subtaskId"`
	LastAction      AgentActionType `json:"lastAction"`
	TargetFile      string          `json:"targetFile"`
	ErrorMessage    string          `json:"errorMessage"`
	LastHeartbeatAt string          `json:"lastHeartbeatAt"`
	DurationMs      int64           `json:"durationMs"`
}

// TaskStatusSummary aggregates rollup metrics and subtask progress for a parent task.
type TaskStatusSummary struct {
	ParentTask      ParentTask `json:"parentTask"`
	Subtasks        []Subtask  `json:"subtasks"`
	PendingCount    int        `json:"pendingCount"`
	InProgressCount int        `json:"inProgressCount"`
	CompletedCount  int        `json:"completedCount"`
	FailedCount     int        `json:"failedCount"`
	IsCompleted     bool       `json:"isCompleted"`
}
