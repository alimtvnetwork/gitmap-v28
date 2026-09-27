// Package store — backup_prompts_split_types.go defines domain models for running prompts backup & restore.
package store

// PromptStatusType defines running or enqueued prompt status.
type PromptStatusType string

const (
	// PromptStatusRunning indicates an active executing prompt.
	PromptStatusRunning PromptStatusType = "running"
	// PromptStatusEnqueued indicates a pending queued prompt.
	PromptStatusEnqueued PromptStatusType = "enqueued"
)

// RunningPromptRecord models a single prompt from active execution or queue.
type RunningPromptRecord struct {
	ProjectName    string           `json:"projectName"`
	ProjectPath    string           `json:"projectPath"`
	ProjectId      string           `json:"projectId"`
	ConversationId string           `json:"conversationId"`
	SequenceId     string           `json:"sequenceId"`
	Prompt         string           `json:"prompt"`
	Snippet        string           `json:"snippet"`
	Status         PromptStatusType `json:"status"`
	WordCount      int              `json:"wordCount"`
	Node           string           `json:"node,omitempty"`
	Host           string           `json:"host,omitempty"`
	CreatedAt      string           `json:"createdAt"`
}

// PromptBackupSummary holds metadata and items for a prompt backup batch.
type PromptBackupSummary struct {
	BatchId       string                `json:"batchId"`
	TotalPrompts  int                   `json:"totalPrompts"`
	RunningCount  int                   `json:"runningCount"`
	EnqueuedCount int                   `json:"enqueuedCount"`
	DatabasePath  string                `json:"databasePath"`
	DatabaseSize  int64                 `json:"databaseSizeBytes"`
	CreatedAt     string                `json:"createdAt"`
	Node          string                `json:"node,omitempty"`
	Host          string                `json:"host,omitempty"`
	Items         []RunningPromptRecord `json:"items,omitempty"`
}

// PromptBackupBatchRecord models a database row in PromptBackupBatch.
type PromptBackupBatchRecord struct {
	BatchID       string `json:"batchId"`
	SourcePath    string `json:"sourcePath"`
	TotalPrompts  int    `json:"totalPrompts"`
	RunningCount  int    `json:"runningCount"`
	EnqueuedCount int    `json:"enqueuedCount"`
	CreatedAt     string `json:"createdAt"`
	Note          string `json:"note"`
	Status        string `json:"status"`
	Node          string `json:"node,omitempty"`
	Host          string `json:"host,omitempty"`
}

// PromptRestoreLedgerRecord models a database row in PromptRestoreLedger.
type PromptRestoreLedgerRecord struct {
	RestoreID     string `json:"restoreId"`
	BatchID       string `json:"batchId"`
	RestoredAt    string `json:"restoredAt"`
	RestoredCount int    `json:"restoredCount"`
	IsKept        bool   `json:"isKept"`
	PrunedAt      string `json:"prunedAt"`
}

// RestoreOptions encapsulates configuration for running-prompts restore.
type RestoreOptions struct {
	IsKeep     bool   `json:"isKeep"`
	IsJSON     bool   `json:"isJSON"`
	IsSSH      bool   `json:"isSSH"`
	TargetFile string `json:"targetFile"`
}
