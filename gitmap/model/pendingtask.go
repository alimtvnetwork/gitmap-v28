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
