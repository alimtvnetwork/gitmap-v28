// Package sql provides embedded SQL schemas and parameterized statements for GitMap stores.
package sql

import (
	_ "embed"
)

//go:embed repos.sql
var reposSQL string

//go:embed prompts.sql
var promptsSQL string

//go:embed cluster.sql
var clusterSQL string

// ReposSchema returns the full embedded DDL and statements for repositories, groups, and aliases.
func ReposSchema() string {
	return reposSQL
}

// PromptsSchema returns the full embedded DDL and statements for prompt backups.
func PromptsSchema() string {
	return promptsSQL
}

// ClusterSchema returns the full embedded DDL and statements for cluster fleet orchestration.
func ClusterSchema() string {
	return clusterSQL
}

// HasSchema returns whether a valid schema exists for the requested domain name.
func HasSchema(domain string) bool {
	switch domain {
	case "repos", "prompts", "cluster":
		return true
	default:
		return false
	}
}

// GetSchema resolves the domain SQL schema text if present.
func GetSchema(domain string) (string, bool) {
	switch domain {
	case "repos":
		return ReposSchema(), true
	case "prompts":
		return PromptsSchema(), true
	case "cluster":
		return ClusterSchema(), true
	default:
		return "", false
	}
}

// Repository, Group, Alias, Release, and WorkDirectory query constants.
const (
	SQLUpsertRepo        = "INSERT INTO Repo (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(AbsolutePath) DO UPDATE SET Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl, SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath, CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes, IdentifiedTransport=excluded.IdentifiedTransport, UpdatedAt=CURRENT_TIMESTAMP"
	SQLSelectAllRepos    = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo ORDER BY Slug"
	SQLSelectRepoByID    = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE RepoId = ?"
	SQLSelectRepoBySlug  = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Slug = ?"
	SQLSelectRepoByPath  = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE REPLACE(AbsolutePath, '\\', '/') = ? COLLATE NOCASE"
	SQLSuggestRepoBySlug = "SELECT Slug FROM Repo WHERE Slug LIKE ? OR RepoName LIKE ? ORDER BY Slug ASC LIMIT 10"
	SQLDeleteRepoByPath  = "DELETE FROM Repo WHERE REPLACE(AbsolutePath, '\\', '/') = ? COLLATE NOCASE"
	SQLDeleteRepoBySlug  = "DELETE FROM Repo WHERE Slug = ?"
	SQLDeleteRepoByID    = "DELETE FROM Repo WHERE RepoId = ?"

	SQLInsertGroup       = "INSERT INTO \"Group\" (Name, Description, Color) VALUES (?, ?, ?)"
	SQLSelectAllGroups   = "SELECT GroupId, Name, Description, Color, CreatedAt FROM \"Group\" ORDER BY Name"
	SQLSelectGroupByName = "SELECT GroupId, Name, Description, Color, CreatedAt FROM \"Group\" WHERE Name = ?"
	SQLSelectGroupByID   = "SELECT GroupId, Name, Description, Color, CreatedAt FROM \"Group\" WHERE GroupId = ?"
	SQLDeleteGroup       = "DELETE FROM \"Group\" WHERE Name = ?"
	SQLInsertGroupRepo   = "INSERT OR IGNORE INTO GroupRepo (GroupId, RepoId) VALUES (?, ?)"
	SQLDeleteGroupRepo   = "DELETE FROM GroupRepo WHERE GroupId = ? AND RepoId = ?"
	SQLSelectGroupRepos  = "SELECT r.RepoId, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch, r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes, r.IdentifiedTransport FROM Repo r JOIN GroupRepo gr ON r.RepoId = gr.RepoId WHERE gr.GroupId = ? ORDER BY r.Slug"
	SQLCountGroupRepos   = "SELECT COUNT(*) FROM GroupRepo WHERE GroupId = ?"

	SQLInsertAlias              = "INSERT INTO Alias (Alias, RepoId) VALUES (?, ?)"
	SQLInsertAliasFull          = "INSERT INTO Alias (Alias, RepoId, IsPrimary, Source) VALUES (?, ?, ?, ?)"
	SQLUpdateAlias              = "UPDATE Alias SET RepoId = ? WHERE Alias = ?"
	SQLSelectAllAliases         = "SELECT a.AliasId, a.Alias, a.RepoId, a.IsPrimary, a.Source, a.CreatedAt FROM Alias a ORDER BY a.Alias"
	SQLSelectAliasByName        = "SELECT a.AliasId, a.Alias, a.RepoId, a.IsPrimary, a.Source, a.CreatedAt FROM Alias a WHERE a.Alias = ?"
	SQLSelectAliasesByRepoID    = "SELECT a.AliasId, a.Alias, a.RepoId, COALESCE(a.IsPrimary, 0), COALESCE(a.Source, 'manual'), a.CreatedAt FROM Alias a WHERE a.RepoId = ? ORDER BY a.IsPrimary DESC, a.AliasId ASC"
	SQLDeleteAlias              = "DELETE FROM Alias WHERE Alias = ?"
	SQLSelectAliasWithRepo      = "SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt, r.AbsolutePath, r.Slug FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId WHERE a.Alias = ?"
	SQLSelectAllAliasesWithRepo = "SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt, r.AbsolutePath, r.Slug FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId ORDER BY a.Alias"
	SQLSelectUnaliasedRepos     = "SELECT r.RepoId, r.Slug, r.RepoName FROM Repo r LEFT JOIN Alias a ON r.RepoId = a.RepoId WHERE a.AliasId IS NULL ORDER BY r.Slug"

	SQLUpsertRelease                = "INSERT INTO Release (RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(RepoId, Tag) DO UPDATE SET Version=excluded.Version, Branch=excluded.Branch, SourceBranch=excluded.SourceBranch, CommitSha=excluded.CommitSha, Changelog=excluded.Changelog, Notes=excluded.Notes, IsDraft=excluded.IsDraft, IsPreRelease=excluded.IsPreRelease, IsLatest=excluded.IsLatest, Source=excluded.Source"
	SQLSelectAllReleases            = "SELECT ReleaseId, RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt FROM Release ORDER BY CreatedAt DESC"
	SQLSelectAllReleasesAcrossRepos = "SELECT R.ReleaseId, R.RepoId, P.Slug, P.AbsolutePath, R.Version, R.Tag, R.Branch, R.CommitSha, R.Source, R.IsDraft, R.IsLatest, R.IsPreRelease, R.CreatedAt FROM Release R JOIN Repo P ON P.RepoId = R.RepoId ORDER BY R.CreatedAt DESC, P.Slug ASC"
	SQLSelectReleaseByTag           = "SELECT ReleaseId, RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt FROM Release WHERE Tag = ?"
	SQLClearLatestRelease           = "UPDATE Release SET IsLatest = 0 WHERE IsLatest = 1 AND RepoId = ?"

	SQLUpsertWorkDir     = "INSERT INTO WorkDirectory (AbsolutePath, Label, IsDefault, UpdatedAt) VALUES (?, ?, ?, CURRENT_TIMESTAMP) ON CONFLICT(AbsolutePath) DO UPDATE SET Label = CASE WHEN ? != '' THEN ? ELSE Label END, IsDefault = CASE WHEN ? != 0 THEN ? ELSE IsDefault END, UpdatedAt = CURRENT_TIMESTAMP"
	SQLSelectAllWorkDirs = "SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, CreatedAt, UpdatedAt FROM WorkDirectory ORDER BY IsDefault DESC, UpdatedAt DESC"
	SQLSetDefaultWorkDir = "UPDATE WorkDirectory SET IsDefault = CASE WHEN WorkDirectoryId = ? OR AbsolutePath = ? THEN 1 ELSE 0 END"
	SQLDeleteWorkDir     = "DELETE FROM WorkDirectory WHERE WorkDirectoryId = ? OR AbsolutePath = ?"
)

// PromptBackupBatch, PromptBackupItem, and PromptRestoreLedger query constants.
const (
	SQLInsertPromptBackupBatch       = "INSERT INTO PromptBackupBatch (BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note) VALUES (?, ?, ?, ?, ?, ?, ?)"
	SQLSelectPromptBackupBatchByID   = "SELECT BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note FROM PromptBackupBatch WHERE BatchId = ?"
	SQLSelectListPromptBackupBatches = "SELECT b.BatchId, b.SourcePath, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, b.CreatedAt, b.Note, CASE WHEN r.RestoreId IS NOT NULL THEN 'restored' ELSE 'active' END AS Status FROM PromptBackupBatch b LEFT JOIN (SELECT BatchId, RestoreId FROM PromptRestoreLedger GROUP BY BatchId) r ON b.BatchId = r.BatchId ORDER BY b.CreatedAt DESC"
	SQLCountPromptBackupBatches      = "SELECT COUNT(*) FROM PromptBackupBatch"
	SQLDeletePromptBackupBatchByID   = "DELETE FROM PromptBackupBatch WHERE BatchId = ?"
	SQLDeleteAllPromptBackupBatches  = "DELETE FROM PromptBackupBatch"

	SQLInsertPromptBackupItem           = "INSERT INTO PromptBackupItem (ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	SQLSelectPromptBackupItemsByBatchID = "SELECT ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt FROM PromptBackupItem WHERE BatchId = ? ORDER BY CreatedAt ASC"
	SQLDeletePromptBackupItemsByBatchID = "DELETE FROM PromptBackupItem WHERE BatchId = ?"
	SQLDeleteAllPromptBackupItems       = "DELETE FROM PromptBackupItem"

	SQLInsertPromptRestoreLedger           = "INSERT INTO PromptRestoreLedger (RestoreId, BatchId, RestoredAt, RestoredCount, IsKept, PrunedAt) VALUES (?, ?, ?, ?, ?, ?)"
	SQLSelectExpiredRestoreBatches         = "SELECT BatchId, RestoredAt FROM PromptRestoreLedger WHERE IsKept = 0"
	SQLDeletePromptRestoreLedgersByBatchID = "DELETE FROM PromptRestoreLedger WHERE BatchId = ?"
	SQLDeleteAllPromptRestoreLedgers       = "DELETE FROM PromptRestoreLedger"
)

// ClusterNode, ClusterRun, and ClusterExecResult query constants.
const (
	SQLUpsertClusterNode         = "INSERT INTO ClusterNode (NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(NodeId) DO UPDATE SET Alias = excluded.Alias, DisplayId = excluded.DisplayId, IPAddress = excluded.IPAddress, NodeRole = excluded.NodeRole, OS = excluded.OS, LastHeartbeat = excluded.LastHeartbeat, Status = excluded.Status, PasswordHash = excluded.PasswordHash, PackageManager = excluded.PackageManager"
	SQLSelectListClusterNodes    = "SELECT NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager FROM ClusterNode ORDER BY DisplayId ASC"
	SQLSelectClusterNodeByID     = "SELECT NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager FROM ClusterNode WHERE NodeId = ?"
	SQLUpdateClusterNodePassword = "UPDATE ClusterNode SET PasswordHash = ? WHERE NodeId = ?"
	SQLDeleteClusterNode         = "DELETE FROM ClusterNode WHERE NodeId = ?"

	SQLInsertClusterRun        = "INSERT INTO ClusterRun (RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	SQLUpdateClusterRun        = "UPDATE ClusterRun SET FinishedAt = ?, TotalNodes = ?, SucceededNodes = ?, FailedNodes = ?, SkippedNodes = ? WHERE ClusterRunId = ?"
	SQLSelectClusterRunByRef   = "SELECT ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes FROM ClusterRun WHERE RunRef = ?"
	SQLSelectListClusterRuns   = "SELECT ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause, StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes FROM ClusterRun ORDER BY StartedAt DESC"
	SQLDeleteClusterRunsBefore = "DELETE FROM ClusterRun WHERE StartedAt < ?"

	SQLInsertClusterExecResult         = "INSERT INTO ClusterExecResult (ClusterRunId, NodeId, SubCommand, CommandText, ResultStatus, ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	SQLUpdateClusterExecResult         = "UPDATE ClusterExecResult SET CommandText = ?, ResultStatus = ?, ExitCode = ?, Stdout = ?, Stderr = ?, StartedAt = ?, FinishedAt = ?, DurationMs = ?, ErrorMessage = ? WHERE ClusterExecResultId = ?"
	SQLSelectClusterExecResultsByRunID = "SELECT ClusterExecResultId, ClusterRunId, NodeId, SubCommand, CommandText, ResultStatus, ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage FROM ClusterExecResult WHERE ClusterRunId = ? ORDER BY ClusterExecResultId ASC"

	SQLCountClusterRuns              = "SELECT COUNT(*) FROM ClusterRun"
	SQLCountClusterExecResults       = "SELECT COUNT(*) FROM ClusterExecResult"
	SQLCountClusterSuccessCommands   = "SELECT COUNT(*) FROM ClusterExecResult WHERE ResultStatus = 1"
	SQLSelectMostTargetedClusterNode = "SELECT NodeId FROM ClusterExecResult GROUP BY NodeId ORDER BY COUNT(*) DESC LIMIT 1"
	SQLSelectMostUsedClusterSubCmd   = "SELECT SubCommand FROM ClusterExecResult GROUP BY SubCommand ORDER BY COUNT(*) DESC LIMIT 1"
)
