// Package constants — constants_visibility_store.go: SQLite schema for
// the bulk wildcard visibility audit trail (spec/01-app/116, plan
// steps 15-16). Two tables:
//
//	MakeAllVisibilityRun     — one row per `make-all-*` invocation.
//	MakeAllVisibilityResult  — one row per (Run, repo) pair.
//
// Result rows are written in three phases by the apply loop:
//
//  1. Pre-prompt: insert one row per matched repo with Status='Pending'
//     so a crash mid-prompt still leaves an auditable trail.
//  2. Post-exclusion: flip excluded rows to Status='Excluded'.
//  3. Per-apply: UPDATE the row with the final Status (Ok/Skipped/Failed),
//     PrevVisibility, NewVisibility, FailureMessage.
//
// Strict PascalCase column names + INTEGER PRIMARY KEY AUTOINCREMENT per
// project-wide DB convention (mem://tech/database-architecture).
package constants

// Bulk-visibility table names.
const (
	TableMakeAllVisibilityRun    = "MakeAllVisibilityRun"
	TableMakeAllVisibilityResult = "MakeAllVisibilityResult"
)

// CommandKindEnum — discriminates which of the four bulk command IDs
// kicked off the run. Stored as TEXT in MakeAllVisibilityRun.CommandKind.
const (
	CommandKindMakeAllPublic  = "MakeAllPublic"
	CommandKindMakeAllPrivate = "MakeAllPrivate"
	CommandKindVisibilityUndo = "VisibilityUndo"
	CommandKindVisibilityRedo = "VisibilityRedo"
)

// ResultStatusEnum — terminal status of a single repo within a run.
// 'Pending' is the pre-apply placeholder; everything else is terminal.
const (
	ResultStatusPending  = "Pending"
	ResultStatusOk       = "Ok"
	ResultStatusSkipped  = "Skipped"
	ResultStatusFailed   = "Failed"
	ResultStatusExcluded = "Excluded"
)

// SQLCreateMakeAllVisibilityRun — one row per invocation.
// PatternList is the raw comma-separated input (audit-faithful).
// MatchedCount / ExcludedCount / OkCount / SkippedCount / FailedCount
// are populated incrementally so a query mid-run shows live progress.
const SQLCreateMakeAllVisibilityRun = `CREATE TABLE IF NOT EXISTS MakeAllVisibilityRun (
	MakeAllVisibilityRunId INTEGER PRIMARY KEY AUTOINCREMENT,
	CommandKind            TEXT NOT NULL,
	TargetVisibility       TEXT NOT NULL,
	Provider               TEXT NOT NULL,
	Owner                  TEXT NOT NULL,
	TargetRaw              TEXT NOT NULL,
	PatternList            TEXT NOT NULL,
	YesFlag                INTEGER NOT NULL DEFAULT 0,
	VerboseFlag            INTEGER NOT NULL DEFAULT 0,
	OwnerRepoTotal         INTEGER NOT NULL DEFAULT 0,
	MatchedCount           INTEGER NOT NULL DEFAULT 0,
	ExcludedCount          INTEGER NOT NULL DEFAULT 0,
	OkCount                INTEGER NOT NULL DEFAULT 0,
	SkippedCount           INTEGER NOT NULL DEFAULT 0,
	FailedCount            INTEGER NOT NULL DEFAULT 0,
	ExitCode               INTEGER NOT NULL DEFAULT 0,
	StartedAt              TEXT NOT NULL,
	FinishedAt             TEXT DEFAULT '',
	CreatedAt              TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQLCreateMakeAllVisibilityResult — one row per (Run, repo).
const SQLCreateMakeAllVisibilityResult = `CREATE TABLE IF NOT EXISTS MakeAllVisibilityResult (
	MakeAllVisibilityResultId INTEGER PRIMARY KEY AUTOINCREMENT,
	MakeAllVisibilityRunId    INTEGER NOT NULL,
	RepoName                  TEXT NOT NULL,
	MatchedPattern            TEXT NOT NULL,
	Status                    TEXT NOT NULL DEFAULT 'Pending',
	PrevVisibility            TEXT DEFAULT '',
	NewVisibility             TEXT DEFAULT '',
	FailureMessage            TEXT DEFAULT '',
	StartedAt                 TEXT DEFAULT '',
	FinishedAt                TEXT DEFAULT '',
	DurationMs                INTEGER NOT NULL DEFAULT 0,
	CreatedAt                 TEXT DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (MakeAllVisibilityRunId) REFERENCES MakeAllVisibilityRun(MakeAllVisibilityRunId)
)`

// Helper index on the result table for fast per-run queries
// (history view, retry-failed selection).
const SQLCreateMakeAllVisibilityResultRunIndex = `CREATE INDEX IF NOT EXISTS IdxMakeAllVisibilityResultRun
	ON MakeAllVisibilityResult(MakeAllVisibilityRunId)`
