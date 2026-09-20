package constants

// SQL statements for SplitDatabaseRegistry.
const (
	SQLCreateSplitDatabaseRegistry = `CREATE TABLE IF NOT EXISTS SplitDatabaseRegistry (
    SplitDatabaseRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseType            TEXT NOT NULL,
    DatabaseKey             TEXT NOT NULL,
    DatabasePath            TEXT NOT NULL,
    SizeBytes               INTEGER NOT NULL DEFAULT 0,
    TableCount              INTEGER NOT NULL DEFAULT 0,
    RecordCount             INTEGER NOT NULL DEFAULT 0,
    SchemaVersion           INTEGER NOT NULL DEFAULT 1,
    Status                  TEXT NOT NULL DEFAULT 'active',
    IsActive                INTEGER NOT NULL DEFAULT 1,
    IsAttached              INTEGER NOT NULL DEFAULT 0,
    Description             TEXT NULL,
    Notes                   TEXT NULL,
    Comments                TEXT NULL,
    LastAccessedAt          INTEGER NOT NULL DEFAULT 0,
    LastSyncedAt            INTEGER NOT NULL DEFAULT 0,
    CreatedAt               INTEGER NOT NULL DEFAULT (unixepoch()),
    UpdatedAt               INTEGER NOT NULL DEFAULT (unixepoch()),
    UNIQUE(DatabaseType, DatabaseKey)
);`

	SQLCreateSplitDatabaseRegistryTypeIndex = `CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_Type ON SplitDatabaseRegistry(DatabaseType);`

	SQLCreateSplitDatabaseRegistryStatusIndex = `CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_Status ON SplitDatabaseRegistry(Status);`

	SQLCreateSplitDatabaseRegistryUpdatedAtIndex = `CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_UpdatedAt ON SplitDatabaseRegistry(UpdatedAt);`

	SQLUpsertSplitDatabaseRegistry = `INSERT INTO SplitDatabaseRegistry (
    DatabaseType, DatabaseKey, DatabasePath, SizeBytes, TableCount,
    RecordCount, SchemaVersion, Status, IsActive, IsAttached,
    Description, Notes, Comments, LastAccessedAt, LastSyncedAt,
    CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(DatabaseType, DatabaseKey) DO UPDATE SET
    DatabasePath = excluded.DatabasePath,
    SizeBytes = excluded.SizeBytes,
    TableCount = excluded.TableCount,
    RecordCount = excluded.RecordCount,
    SchemaVersion = excluded.SchemaVersion,
    Status = excluded.Status,
    IsActive = excluded.IsActive,
    IsAttached = excluded.IsAttached,
    Description = CASE WHEN excluded.Description IS NOT NULL AND excluded.Description != '' THEN excluded.Description ELSE SplitDatabaseRegistry.Description END,
    Notes = CASE WHEN excluded.Notes IS NOT NULL AND excluded.Notes != '' THEN excluded.Notes ELSE SplitDatabaseRegistry.Notes END,
    Comments = CASE WHEN excluded.Comments IS NOT NULL AND excluded.Comments != '' THEN excluded.Comments ELSE SplitDatabaseRegistry.Comments END,
    LastAccessedAt = CASE WHEN excluded.LastAccessedAt != 0 THEN excluded.LastAccessedAt ELSE SplitDatabaseRegistry.LastAccessedAt END,
    LastSyncedAt = excluded.LastSyncedAt,
    UpdatedAt = excluded.UpdatedAt;`

	SQLSelectSplitDB = `SELECT SplitDatabaseRegistryId, DatabaseType, DatabaseKey, DatabasePath,
    SizeBytes, TableCount, RecordCount, SchemaVersion, Status,
    IsActive, IsAttached, Description, Notes, Comments,
    LastAccessedAt, LastSyncedAt, CreatedAt, UpdatedAt
FROM SplitDatabaseRegistry
WHERE DatabaseType = ? AND DatabaseKey = ?;`

	SQLSelectListSplitDBByType = `SELECT SplitDatabaseRegistryId, DatabaseType, DatabaseKey, DatabasePath,
    SizeBytes, TableCount, RecordCount, SchemaVersion, Status,
    IsActive, IsAttached, Description, Notes, Comments,
    LastAccessedAt, LastSyncedAt, CreatedAt, UpdatedAt
FROM SplitDatabaseRegistry
WHERE DatabaseType = ?
ORDER BY DatabaseType, DatabaseKey;`

	SQLSelectListAllSplitDB = `SELECT SplitDatabaseRegistryId, DatabaseType, DatabaseKey, DatabasePath,
    SizeBytes, TableCount, RecordCount, SchemaVersion, Status,
    IsActive, IsAttached, Description, Notes, Comments,
    LastAccessedAt, LastSyncedAt, CreatedAt, UpdatedAt
FROM SplitDatabaseRegistry
ORDER BY DatabaseType, DatabaseKey;`

	SQLCreateStartupItem = `CREATE TABLE IF NOT EXISTS StartupItem (
    StartupItemId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name          TEXT UNIQUE NOT NULL,
    TargetType    TEXT NOT NULL,
    TargetPath    TEXT NOT NULL,
    CommandArgs   TEXT NULL,
    IconPath      TEXT NULL,
    RunFrequency  TEXT NOT NULL DEFAULT 'everytime',
    IsActive      INTEGER NOT NULL DEFAULT 1,
    Description   TEXT NULL,
    CreatedAt     INTEGER NOT NULL DEFAULT (unixepoch()),
    UpdatedAt     INTEGER NOT NULL DEFAULT (unixepoch())
);`

	SQLCreateStartupLog = `CREATE TABLE IF NOT EXISTS StartupLog (
    StartupLogId   INTEGER PRIMARY KEY AUTOINCREMENT,
    StartupItemId  INTEGER NOT NULL,
    RunAt          INTEGER NOT NULL DEFAULT (unixepoch()),
    DurationMs     INTEGER NOT NULL DEFAULT 0,
    IsSuccess      INTEGER NOT NULL DEFAULT 1,
    ExitCode       INTEGER NOT NULL DEFAULT 0,
    OutputSummary  TEXT NULL,
    Notes          TEXT NULL,
    Comments       TEXT NULL,
    FOREIGN KEY(StartupItemId) REFERENCES StartupItem(StartupItemId) ON DELETE CASCADE
);`

	// SQL statements for PR Split DB (Spec 129).
	SQLCreatePullRequest = `CREATE TABLE IF NOT EXISTS PullRequest (
    PullRequestId   INTEGER PRIMARY KEY AUTOINCREMENT,
    PrNumber        INTEGER NOT NULL,
    Title           TEXT NOT NULL,
    Description     TEXT NOT NULL DEFAULT '',
    SourceBranch    TEXT NOT NULL,
    TargetBranch    TEXT NOT NULL DEFAULT 'main',
    Status          TEXT NOT NULL DEFAULT 'open',
    MergeCommitSha  TEXT NULL,
    CreatedAt       INTEGER NOT NULL,
    MergedAt        INTEGER NULL,
    ClosedAt        INTEGER NULL,
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    UpdatedAt       INTEGER NOT NULL
);`

	SQLCreatePullRequestIndexes = `CREATE UNIQUE INDEX IF NOT EXISTS IdxPullRequest_PrNumber ON PullRequest (PrNumber);
CREATE INDEX IF NOT EXISTS IdxPullRequest_Status ON PullRequest (Status);
CREATE INDEX IF NOT EXISTS IdxPullRequest_SourceBranch ON PullRequest (SourceBranch);
CREATE INDEX IF NOT EXISTS IdxPullRequest_TargetBranch ON PullRequest (TargetBranch);`

	SQLCreatePrRelease = `CREATE TABLE IF NOT EXISTS PrRelease (
    PrReleaseId     INTEGER PRIMARY KEY AUTOINCREMENT,
    PullRequestId   INTEGER NOT NULL,
    ReleaseTag      TEXT NOT NULL,
    CommitSha       TEXT NOT NULL,
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    CreatedAt       INTEGER NOT NULL,
    FOREIGN KEY (PullRequestId) REFERENCES PullRequest(PullRequestId) ON DELETE CASCADE
);`

	SQLCreatePrReleaseIndexes = `CREATE INDEX IF NOT EXISTS IdxPrRelease_PullRequestId ON PrRelease (PullRequestId);
CREATE INDEX IF NOT EXISTS IdxPrRelease_ReleaseTag ON PrRelease (ReleaseTag);`

	SQLCreatePrBranch = `CREATE TABLE IF NOT EXISTS PrBranch (
    PrBranchId      INTEGER PRIMARY KEY AUTOINCREMENT,
    BranchName      TEXT NOT NULL UNIQUE,
    BranchType      TEXT NOT NULL DEFAULT 'feature',
    IsMerged        INTEGER NOT NULL DEFAULT 0,
    IsDeleted       INTEGER NOT NULL DEFAULT 0,
    CreatedAt       INTEGER NOT NULL,
    MergedAt        INTEGER NULL,
    DeletedAt       INTEGER NULL,
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    UpdatedAt       INTEGER NOT NULL
);`

	SQLCreatePrBranchIndexes = `CREATE UNIQUE INDEX IF NOT EXISTS IdxPrBranch_BranchName ON PrBranch (BranchName);
CREATE INDEX IF NOT EXISTS IdxPrBranch_IsMerged ON PrBranch (IsMerged);
CREATE INDEX IF NOT EXISTS IdxPrBranch_IsDeleted ON PrBranch (IsDeleted);`

	SQLCreatePrViews = `CREATE VIEW IF NOT EXISTS pull_requests AS
SELECT PullRequestId AS id, PrNumber AS pr_number, Title AS title, Description AS description,
       SourceBranch AS source_branch, TargetBranch AS target_branch, Status AS status,
       MergeCommitSha AS merge_commit_sha, CreatedAt AS created_at, MergedAt AS merged_at,
       ClosedAt AS closed_at, Notes AS notes, Comments AS comments, UpdatedAt AS updated_at
FROM PullRequest;

CREATE VIEW IF NOT EXISTS pr_releases AS
SELECT PrReleaseId AS id, PullRequestId AS pr_id, ReleaseTag AS release_tag,
       CommitSha AS commit_sha, Notes AS notes, Comments AS comments, CreatedAt AS created_at
FROM PrRelease;

CREATE VIEW IF NOT EXISTS pr_branches AS
SELECT PrBranchId AS id, BranchName AS branch_name, BranchType AS branch_type,
       IsMerged AS is_merged, IsDeleted AS is_deleted, CreatedAt AS created_at,
       MergedAt AS merged_at, DeletedAt AS deleted_at, Notes AS notes, Comments AS comments,
       UpdatedAt AS updated_at
FROM PrBranch;`

	SQLInsertPullRequest = `INSERT INTO PullRequest (
    PrNumber, Title, Description, SourceBranch, TargetBranch,
    Status, MergeCommitSha, CreatedAt, MergedAt, ClosedAt,
    Notes, Comments, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	SQLSelectPullRequestByNumber = `SELECT PullRequestId, PrNumber, Title, Description,
    SourceBranch, TargetBranch, Status, MergeCommitSha,
    CreatedAt, MergedAt, ClosedAt, Notes, Comments, UpdatedAt
FROM PullRequest
WHERE PrNumber = ?
LIMIT 1;`

	SQLUpdatePullRequestStatus = `UPDATE PullRequest
SET Status = ?,
    MergeCommitSha = CASE WHEN ? != '' THEN ? ELSE MergeCommitSha END,
    MergedAt = CASE WHEN ? = 'merged' AND (MergedAt IS NULL OR MergedAt = 0) THEN ? ELSE MergedAt END,
    ClosedAt = CASE WHEN ? = 'closed' AND (ClosedAt IS NULL OR ClosedAt = 0) THEN ? ELSE ClosedAt END,
    UpdatedAt = ?
WHERE PrNumber = ?;`

	SQLInsertPrRelease = `INSERT INTO PrRelease (
    PullRequestId, ReleaseTag, CommitSha, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?);`

	SQLUpsertPrBranch = `INSERT INTO PrBranch (
    BranchName, BranchType, IsMerged, IsDeleted,
    CreatedAt, MergedAt, DeletedAt, Notes, Comments, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(BranchName) DO UPDATE SET
    BranchType = excluded.BranchType,
    IsMerged = excluded.IsMerged,
    IsDeleted = excluded.IsDeleted,
    MergedAt = CASE WHEN excluded.MergedAt IS NOT NULL AND excluded.MergedAt != 0 THEN excluded.MergedAt ELSE PrBranch.MergedAt END,
    DeletedAt = CASE WHEN excluded.DeletedAt IS NOT NULL AND excluded.DeletedAt != 0 THEN excluded.DeletedAt ELSE PrBranch.DeletedAt END,
    Notes = CASE WHEN excluded.Notes IS NOT NULL AND excluded.Notes != '' THEN excluded.Notes ELSE PrBranch.Notes END,
    Comments = CASE WHEN excluded.Comments IS NOT NULL AND excluded.Comments != '' THEN excluded.Comments ELSE PrBranch.Comments END,
    UpdatedAt = excluded.UpdatedAt;`

	SQLSelectListActivePrBranches = `SELECT PrBranchId, BranchName, BranchType, IsMerged, IsDeleted,
    CreatedAt, MergedAt, DeletedAt, Notes, Comments, UpdatedAt
FROM PrBranch
WHERE IsDeleted = 0 AND IsMerged = 0
ORDER BY CreatedAt DESC;`

	SQLSelectListMergedPrBranches = `SELECT PrBranchId, BranchName, BranchType, IsMerged, IsDeleted,
    CreatedAt, MergedAt, DeletedAt, Notes, Comments, UpdatedAt
FROM PrBranch
WHERE IsMerged = 1 AND IsDeleted = 0
ORDER BY CreatedAt DESC;`

	SQLMarkPrBranchDeleted = `UPDATE PrBranch
SET IsDeleted = 1,
    DeletedAt = ?,
    UpdatedAt = ?
WHERE BranchName = ?;`
)
