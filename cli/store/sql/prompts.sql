-- Prompt Backup Batch, Item, and Restore Ledger Schema and Parameterized Queries

-- Schema: PromptBackupBatch Table
CREATE TABLE IF NOT EXISTS PromptBackupBatch (
    BatchId       TEXT PRIMARY KEY,
    SourcePath    TEXT NOT NULL DEFAULT '',
    TotalPrompts  INTEGER NOT NULL DEFAULT 0,
    RunningCount  INTEGER NOT NULL DEFAULT 0,
    EnqueuedCount INTEGER NOT NULL DEFAULT 0,
    CreatedAt     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Note          TEXT NOT NULL DEFAULT ''
);

-- Schema: PromptBackupItem Table
CREATE TABLE IF NOT EXISTS PromptBackupItem (
    ItemId         TEXT PRIMARY KEY,
    BatchId        TEXT NOT NULL,
    ProjectName    TEXT NOT NULL,
    ProjectPath    TEXT NOT NULL,
    ProjectId      TEXT NOT NULL,
    ConversationId TEXT NOT NULL,
    SequenceId     TEXT NOT NULL DEFAULT '',
    PromptText     TEXT NOT NULL,
    PromptStatus   TEXT NOT NULL,
    WordCount      INTEGER NOT NULL DEFAULT 0,
    MediaPathsJson TEXT NOT NULL DEFAULT '[]',
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS IdxPromptBackupItem_BatchId ON PromptBackupItem(BatchId);
CREATE INDEX IF NOT EXISTS IdxPromptBackupItem_ProjectId ON PromptBackupItem(ProjectId);

-- Schema: PromptRestoreLedger Table
CREATE TABLE IF NOT EXISTS PromptRestoreLedger (
    RestoreId     TEXT PRIMARY KEY,
    BatchId       TEXT NOT NULL,
    RestoredAt    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    RestoredCount INTEGER NOT NULL DEFAULT 0,
    IsKept        INTEGER NOT NULL DEFAULT 0,
    PrunedAt      TEXT NOT NULL DEFAULT '',
    FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS IdxPromptRestoreLedger_BatchId ON PromptRestoreLedger(BatchId);

-- Query: Insert PromptBackupBatch
INSERT INTO PromptBackupBatch (BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- Query: Select PromptBackupBatch By ID
SELECT BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note
FROM PromptBackupBatch WHERE BatchId = ?;

-- Query: List PromptBackupBatches With Restore Status
SELECT b.BatchId, b.SourcePath, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, b.CreatedAt, b.Note,
    CASE WHEN r.RestoreId IS NOT NULL THEN 'restored' ELSE 'active' END AS Status
FROM PromptBackupBatch b
LEFT JOIN (SELECT BatchId, RestoreId FROM PromptRestoreLedger GROUP BY BatchId) r ON b.BatchId = r.BatchId
ORDER BY b.CreatedAt DESC;

-- Query: Count PromptBackupBatches
SELECT COUNT(*) FROM PromptBackupBatch;

-- Query: Delete PromptBackupBatch By ID
DELETE FROM PromptBackupBatch WHERE BatchId = ?;

-- Query: Delete All PromptBackupBatches
DELETE FROM PromptBackupBatch;

-- Query: Insert PromptBackupItem
INSERT INTO PromptBackupItem (ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Query: Select PromptBackupItems By Batch ID
SELECT ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt
FROM PromptBackupItem WHERE BatchId = ? ORDER BY CreatedAt ASC;

-- Query: Delete PromptBackupItems By Batch ID
DELETE FROM PromptBackupItem WHERE BatchId = ?;

-- Query: Delete All PromptBackupItems
DELETE FROM PromptBackupItem;

-- Query: Insert PromptRestoreLedger
INSERT INTO PromptRestoreLedger (RestoreId, BatchId, RestoredAt, RestoredCount, IsKept, PrunedAt)
VALUES (?, ?, ?, ?, ?, ?);

-- Query: Select Expired Restore Batches
SELECT BatchId, RestoredAt FROM PromptRestoreLedger WHERE IsKept = 0;

-- Query: Delete PromptRestoreLedgers By Batch ID
DELETE FROM PromptRestoreLedger WHERE BatchId = ?;

-- Query: Delete All PromptRestoreLedgers
DELETE FROM PromptRestoreLedger;
