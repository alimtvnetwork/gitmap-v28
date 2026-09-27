# 167: Data Contracts & Database Schema

## 1. Running Prompts Backup SQLite Split-DB Schema

Default Location: `data/backup-prompts/sql.db` (resolved relative to binary/data folder or explicitly via `-file` / `-f`).

### 1.1 `PromptBackupBatch` Table
```sql
CREATE TABLE IF NOT EXISTS PromptBackupBatch (
    BatchId TEXT PRIMARY KEY,
    SourcePath TEXT NOT NULL DEFAULT '',
    TotalPrompts INTEGER NOT NULL DEFAULT 0,
    RunningCount INTEGER NOT NULL DEFAULT 0,
    EnqueuedCount INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Note TEXT NOT NULL DEFAULT ''
);
```

### 1.2 `PromptBackupItem` Table
```sql
CREATE TABLE IF NOT EXISTS PromptBackupItem (
    ItemId TEXT PRIMARY KEY,
    BatchId TEXT NOT NULL,
    ProjectName TEXT NOT NULL,
    ProjectPath TEXT NOT NULL,
    ProjectId TEXT NOT NULL,
    ConversationId TEXT NOT NULL,
    SequenceId TEXT NOT NULL DEFAULT '',
    PromptText TEXT NOT NULL,
    PromptStatus TEXT NOT NULL, -- 'running' or 'enqueued'
    WordCount INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_prompt_backup_item_batch ON PromptBackupItem(BatchId);
CREATE INDEX IF NOT EXISTS idx_prompt_backup_item_project ON PromptBackupItem(ProjectId);
```

### 1.3 `PromptRestoreLedger` Sub-Table
Tracks restored batches and items to support 1-day retention auto-pruning.
```sql
CREATE TABLE IF NOT EXISTS PromptRestoreLedger (
    RestoreId TEXT PRIMARY KEY,
    BatchId TEXT NOT NULL,
    RestoredAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    RestoredCount INTEGER NOT NULL DEFAULT 0,
    IsKept INTEGER NOT NULL DEFAULT 0, -- 1 when --keep/-k flag is passed
    PrunedAt TEXT NOT NULL DEFAULT '',
    FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_prompt_restore_batch ON PromptRestoreLedger(BatchId);
```

## 2. Go Domain Models & Structs

```go
type PromptStatusType string

const (
    PromptStatusRunning  PromptStatusType = "running"
    PromptStatusEnqueued PromptStatusType = "enqueued"
)

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
    CreatedAt      string           `json:"createdAt"`
}

type PromptBackupSummary struct {
    BatchId        string                `json:"batchId"`
    TotalPrompts   int                   `json:"totalPrompts"`
    RunningCount   int                   `json:"runningCount"`
    EnqueuedCount  int                   `json:"enqueuedCount"`
    DatabasePath   string                `json:"databasePath"`
    DatabaseSize   int64                 `json:"databaseSizeBytes"`
    CreatedAt      string                `json:"createdAt"`
    Items          []RunningPromptRecord `json:"items,omitempty"`
}

type RestoreOptions struct {
    IsKeep         bool   `json:"isKeep"`
    IsJSON         bool   `json:"isJSON"`
    TargetFile     string `json:"targetFile"`
}

type ShutdownUntilGreenConfig struct {
    ProjectTargets []string `json:"projectTargets"`
    IntervalSec    int      `json:"intervalSeconds"`
    IsRunningOnly  bool     `json:"isRunningOnly"`
}
```
