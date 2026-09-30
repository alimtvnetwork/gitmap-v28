// Package model — pendingtask.go defines pending and completed task records.
package model

// PendingTaskRecord represents a task awaiting execution.
type PendingTaskRecord struct {
	ID               int64  `json:"id"`
	TaskTypeId       int64  `json:"taskTypeId"`
	TaskTypeName     string `json:"taskTypeName"`
	TargetPath       string `json:"targetPath"`
	WorkingDirectory string `json:"workingDirectory,omitempty"`
	SourceCommand    string `json:"sourceCommand"`
	CommandArgs      string `json:"commandArgs,omitempty"`
	FailureReason    string `json:"failureReason,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
}

// CompletedTaskRecord represents a successfully executed task.
type CompletedTaskRecord struct {
	ID               int64  `json:"id"`
	OriginalTaskId   int64  `json:"originalTaskId"`
	TaskTypeId       int64  `json:"taskTypeId"`
	TaskTypeName     string `json:"taskTypeName"`
	TargetPath       string `json:"targetPath"`
	WorkingDirectory string `json:"workingDirectory,omitempty"`
	SourceCommand    string `json:"sourceCommand"`
	CommandArgs      string `json:"commandArgs,omitempty"`
	CompletedAt      string `json:"completedAt,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
}

// TaskHistoryRecord represents an audit history record across macro, ssh, and installer subsystems.
type TaskHistoryRecord struct {
	TaskHistoryId  int64  `json:"taskHistoryId"`
	TaskId         string `json:"taskId"`
	Section        string `json:"section"`
	Action         string `json:"action"`
	Target         string `json:"target"`
	ForwardPayload string `json:"forwardPayload,omitempty"`
	InversePayload string `json:"inversePayload,omitempty"`
	Status         string `json:"status"`
	RestoredAt     string `json:"restoredAt,omitempty"`
	ExecutedAt     string `json:"executedAt,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
}
