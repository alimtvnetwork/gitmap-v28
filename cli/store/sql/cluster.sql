-- Cluster Fleet Schema and Parameterized Queries (ClusterNode, ClusterRun, ClusterExecResult)

-- Schema: ClusterNode Table
CREATE TABLE IF NOT EXISTS ClusterNode (
    NodeId         TEXT PRIMARY KEY,
    Alias          TEXT NOT NULL DEFAULT '',
    DisplayId      INTEGER NOT NULL DEFAULT 0,
    IPAddress      TEXT NOT NULL DEFAULT '',
    NodeRole       TEXT NOT NULL DEFAULT 'client',
    OS             TEXT NOT NULL DEFAULT 'windows',
    JoinedAt       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    LastHeartbeat  TIMESTAMP,
    Status         TEXT NOT NULL DEFAULT 'online',
    PasswordHash   TEXT,
    PackageManager TEXT,
    DefaultPath    TEXT NOT NULL DEFAULT ''
);

-- Schema: ClusterRun Table
CREATE TABLE IF NOT EXISTS ClusterRun (
    ClusterRunId   INTEGER PRIMARY KEY AUTOINCREMENT,
    RunRef         TEXT NOT NULL UNIQUE,
    CommandKind    INTEGER NOT NULL DEFAULT 0,
    RawCommand     TEXT NOT NULL DEFAULT '',
    TargetSelector TEXT NOT NULL DEFAULT '',
    ExceptClause   TEXT,
    StartedAt      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FinishedAt     TIMESTAMP,
    TotalNodes     INTEGER,
    SucceededNodes INTEGER,
    FailedNodes    INTEGER,
    SkippedNodes   INTEGER
);
CREATE INDEX IF NOT EXISTS IdxClusterRun_StartedAt ON ClusterRun(StartedAt);

-- Schema: ClusterExecResult Table
CREATE TABLE IF NOT EXISTS ClusterExecResult (
    ClusterExecResultId INTEGER PRIMARY KEY AUTOINCREMENT,
    ClusterRunId        INTEGER NOT NULL,
    NodeId              TEXT NOT NULL,
    SubCommand          TEXT NOT NULL,
    CommandText         TEXT,
    ResultStatus        INTEGER NOT NULL DEFAULT 0,
    ExitCode            INTEGER,
    Stdout              TEXT,
    Stderr              TEXT,
    StartedAt           TIMESTAMP,
    FinishedAt          TIMESTAMP,
    DurationMs          INTEGER,
    ErrorMessage        TEXT,
    FOREIGN KEY(ClusterRunId) REFERENCES ClusterRun(ClusterRunId) ON DELETE CASCADE,
    FOREIGN KEY(NodeId) REFERENCES ClusterNode(NodeId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS IdxClusterExecResult_ClusterRunId ON ClusterExecResult(ClusterRunId);
CREATE INDEX IF NOT EXISTS IdxClusterExecResult_NodeId ON ClusterExecResult(NodeId);

-- Query: Upsert ClusterNode
INSERT INTO ClusterNode (NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(NodeId) DO UPDATE SET
    Alias = excluded.Alias,
    DisplayId = excluded.DisplayId,
    IPAddress = excluded.IPAddress,
    NodeRole = excluded.NodeRole,
    OS = excluded.OS,
    LastHeartbeat = excluded.LastHeartbeat,
    Status = excluded.Status,
    PasswordHash = excluded.PasswordHash,
    PackageManager = excluded.PackageManager;

-- Query: Select All ClusterNodes
SELECT NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
FROM ClusterNode ORDER BY DisplayId ASC;

-- Query: Select ClusterNode By ID
SELECT NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
FROM ClusterNode WHERE NodeId = ?;

-- Query: Update ClusterNode Password
UPDATE ClusterNode SET PasswordHash = ? WHERE NodeId = ?;

-- Query: Delete ClusterNode
DELETE FROM ClusterNode WHERE NodeId = ?;

-- Query: Insert ClusterRun
INSERT INTO ClusterRun (RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Query: Update ClusterRun
UPDATE ClusterRun
SET FinishedAt = ?, TotalNodes = ?, SucceededNodes = ?, FailedNodes = ?, SkippedNodes = ?
WHERE ClusterRunId = ?;

-- Query: Select ClusterRun By RunRef
SELECT ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
FROM ClusterRun WHERE RunRef = ?;

-- Query: Select All ClusterRuns
SELECT ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
FROM ClusterRun ORDER BY StartedAt DESC;

-- Query: Delete ClusterRuns Before Timestamp
DELETE FROM ClusterRun WHERE StartedAt < ?;

-- Query: Insert ClusterExecResult
INSERT INTO ClusterExecResult (ClusterRunId, NodeId, SubCommand, CommandText, ResultStatus, ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Query: Update ClusterExecResult
UPDATE ClusterExecResult
SET CommandText = ?, ResultStatus = ?, ExitCode = ?, Stdout = ?, Stderr = ?, StartedAt = ?, FinishedAt = ?, DurationMs = ?, ErrorMessage = ?
WHERE ClusterExecResultId = ?;

-- Query: Select ClusterExecResults By Run ID
SELECT ClusterExecResultId, ClusterRunId, NodeId, SubCommand, CommandText, ResultStatus, ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage
FROM ClusterExecResult WHERE ClusterRunId = ? ORDER BY ClusterExecResultId ASC;

-- Query: Count Cluster Runs
SELECT COUNT(*) FROM ClusterRun;

-- Query: Count Cluster Exec Results
SELECT COUNT(*) FROM ClusterExecResult;

-- Query: Count Successful Cluster Commands
SELECT COUNT(*) FROM ClusterExecResult WHERE ResultStatus = 1;

-- Query: Select Most Targeted Cluster Node
SELECT NodeId FROM ClusterExecResult GROUP BY NodeId ORDER BY COUNT(*) DESC LIMIT 1;

-- Query: Select Most Used Cluster SubCommand
SELECT SubCommand FROM ClusterExecResult GROUP BY SubCommand ORDER BY COUNT(*) DESC LIMIT 1;
