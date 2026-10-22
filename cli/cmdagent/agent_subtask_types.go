package cmdagent

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

// SubtaskAddOptions defines options for adding subtasks to a parent task database.
type SubtaskAddOptions struct {
	ParentId   string
	JsonInput  string
	Title      string
	TaskCode   string
	Slug       string
	OwnedFiles string
	IsJson     bool
}

// SubtaskClaimOptions defines options for atomically claiming the next pending subtask.
type SubtaskClaimOptions struct {
	AgentRole string
	TaskId    string
}

// SubtaskStartOptions defines options for transitioning a claimed subtask to in-progress.
type SubtaskStartOptions struct {
	SubtaskId string
	TaskCode  string
	Slug      string
	AgentRole string
	TaskId    string
	IsJson    bool
}

// SubtaskCompleteOptions defines options for completing a subtask with verification evidence.
type SubtaskCompleteOptions struct {
	SubtaskId string
	TaskCode  string
	Slug      string
	AgentRole string
	Evidence  string
	TaskId    string
	IsJson    bool
}

// SubtaskFailOptions defines options for failing a subtask with diagnostic reasoning.
type SubtaskFailOptions struct {
	SubtaskId string
	TaskCode  string
	Slug      string
	AgentRole string
	Reason    string
	TaskId    string
	IsJson    bool
}

// SubtaskLsOptions defines options for listing subtasks.
type SubtaskLsOptions struct {
	TaskId string
	IsJson bool
}

// SubtaskStatusSummary contains aggregate counts for subtasks in a task database.
type SubtaskStatusSummary struct {
	TotalCount      int `json:"totalCount"`
	PendingCount    int `json:"pendingCount"`
	InProgressCount int `json:"inProgressCount"`
	CompletedCount  int `json:"completedCount"`
	FailedCount     int `json:"failedCount"`
}

// SubtaskLsEnvelope is the standard JSON response envelope for `gitmap agent subtask ls --json`.
type SubtaskLsEnvelope struct {
	Status       string               `json:"status"`
	DatabasePath string               `json:"databasePath"`
	ParentTaskId string               `json:"parentTaskId"`
	Summary      SubtaskStatusSummary `json:"summary"`
	Subtasks     []types.Subtask      `json:"subtasks"`
}

// SubtaskEnvelope is a type alias for SubtaskLsEnvelope for backward compatibility.
type SubtaskEnvelope = SubtaskLsEnvelope

// SubtaskAddEnvelope is the standard JSON response envelope for `gitmap agent subtask add`.
type SubtaskAddEnvelope struct {
	Status       string          `json:"status"`
	DatabasePath string          `json:"databasePath"`
	ParentTaskId string          `json:"parentTaskId"`
	AddedCount   int             `json:"addedCount"`
	Subtasks     []types.Subtask `json:"subtasks"`
}

// SubtaskLifecycleEnvelope is the standard compact JSON response envelope for lifecycle commands.
type SubtaskLifecycleEnvelope struct {
	Status       string `json:"status"`
	SubtaskId    string `json:"subtaskId"`
	AgentRole    string `json:"agentRole"`
	DatabasePath string `json:"databasePath"`
	Details      string `json:"details,omitempty"`
}
