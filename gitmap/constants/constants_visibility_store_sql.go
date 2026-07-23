// Package constants — constants_visibility_store_sql.go: INSERT /
// UPDATE statements for the bulk wildcard visibility audit trail.
// Kept in a separate file from the CREATE TABLE schema to honor the
// ≤200-line per-file rule.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan steps 17-18.
package constants

// SQLInsertMakeAllVisibilityRun — pre-prompt INSERT capturing the
// invocation parameters. Counts default to 0 and are flushed by
// SQLUpdateMakeAllVisibilityRunCounts at the end of the run.
const SQLInsertMakeAllVisibilityRun = `INSERT INTO MakeAllVisibilityRun
	(CommandKind, TargetVisibility, Provider, Owner, TargetRaw,
	 PatternList, YesFlag, VerboseFlag, OwnerRepoTotal, MatchedCount,
	 StartedAt)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// SQLInsertMakeAllVisibilityResult — one per matched repo, written
// inside the pre-prompt transaction so a crash leaves an auditable
// 'Pending' trail.
const SQLInsertMakeAllVisibilityResult = `INSERT INTO MakeAllVisibilityResult
	(MakeAllVisibilityRunId, RepoName, MatchedPattern, Status, StartedAt)
	VALUES (?, ?, ?, ?, ?)`

// SQLUpdateMakeAllVisibilityResultExcluded — bulk mark-as-excluded
// after the user trims the matched set via the prompt's exclude grammar.
const SQLUpdateMakeAllVisibilityResultExcluded = `UPDATE MakeAllVisibilityResult
	SET Status = 'Excluded', FinishedAt = ?
	WHERE MakeAllVisibilityResultId = ?`

// SQLUpdateMakeAllVisibilityResult — per-repo terminal status write
// after the apply+verify pipeline finishes for a single repo.
const SQLUpdateMakeAllVisibilityResult = `UPDATE MakeAllVisibilityResult
	SET Status = ?, PrevVisibility = ?, NewVisibility = ?,
	    FailureMessage = ?, FinishedAt = ?, DurationMs = ?
	WHERE MakeAllVisibilityResultId = ?`

// SQLUpdateMakeAllVisibilityRunCounts — final tally flush + exit code.
const SQLUpdateMakeAllVisibilityRunCounts = `UPDATE MakeAllVisibilityRun
	SET ExcludedCount = ?, OkCount = ?, SkippedCount = ?,
	    FailedCount = ?, ExitCode = ?, FinishedAt = ?
	WHERE MakeAllVisibilityRunId = ?`

// Error format strings — Code Red standard (operation + reason).
const (
	ErrMakeAllRunInsertFmt     = "Error: insert MakeAllVisibilityRun failed: %v (operation: SQLInsertMakeAllVisibilityRun, reason: %s)"
	ErrMakeAllResultInsertFmt  = "Error: insert MakeAllVisibilityResult failed: %v (operation: SQLInsertMakeAllVisibilityResult, reason: %s)"
	ErrMakeAllResultUpdateFmt  = "Error: update MakeAllVisibilityResult failed: %v (operation: SQLUpdateMakeAllVisibilityResult, reason: %s)"
	ErrMakeAllRunFinalizeFmt   = "Error: finalize MakeAllVisibilityRun failed: %v (operation: SQLUpdateMakeAllVisibilityRunCounts, reason: %s)"
	ErrMakeAllResultExcludeFmt = "Error: exclude MakeAllVisibilityResult rows failed: %v (operation: SQLUpdateMakeAllVisibilityResultExcluded, reason: %s)"
)

// SQLSelectLatestUndoableRun — picks the most recent run that has at
// least one Ok result with a captured PrevVisibility. Used by
// `gitmap visibility-undo` when no explicit --run is supplied.
const SQLSelectLatestUndoableRun = `SELECT
	MakeAllVisibilityRunId, CommandKind, TargetVisibility, Provider,
	Owner, TargetRaw, OkCount
	FROM MakeAllVisibilityRun
	WHERE OkCount > 0
	ORDER BY MakeAllVisibilityRunId DESC
	LIMIT 1`

// SQLSelectUndoableResultsForRun — Ok results with non-empty Prev/New
// visibility that still need reversing, in deterministic ID order.
const SQLSelectUndoableResultsForRun = `SELECT
	MakeAllVisibilityResultId, RepoName, MatchedPattern,
	PrevVisibility, NewVisibility
	FROM MakeAllVisibilityResult
	WHERE MakeAllVisibilityRunId = ?
	  AND Status = 'Ok'
	  AND PrevVisibility != ''
	  AND PrevVisibility != NewVisibility
	ORDER BY MakeAllVisibilityResultId ASC`

// Error format strings for the select path.
const (
	ErrUndoSelectRunFmt     = "Error: select latest undoable run failed: %v (operation: SQLSelectLatestUndoableRun, reason: %s)"
	ErrUndoSelectResultsFmt = "Error: select undoable results failed: %v (operation: SQLSelectUndoableResultsForRun, reason: %s)"
	ErrUndoNoRunFound       = "Error: no undoable make-all-* run found (operation: visibility-undo, reason: MakeAllVisibilityRun has no row with OkCount>0)"
)

// SQLSelectLatestRunByKind — most recent run with the given CommandKind
// that still has reversible Ok rows. Used by `visibility-redo` (kind=
// 'VisibilityUndo') and by `--run latest` selectors.
const SQLSelectLatestRunByKind = `SELECT
	MakeAllVisibilityRunId, CommandKind, TargetVisibility, Provider,
	Owner, TargetRaw, OkCount
	FROM MakeAllVisibilityRun
	WHERE CommandKind = ? AND OkCount > 0
	ORDER BY MakeAllVisibilityRunId DESC
	LIMIT 1`

// SQLSelectRunByID — exact lookup for the --run <id> selector.
const SQLSelectRunByID = `SELECT
	MakeAllVisibilityRunId, CommandKind, TargetVisibility, Provider,
	Owner, TargetRaw, OkCount
	FROM MakeAllVisibilityRun
	WHERE MakeAllVisibilityRunId = ?`

const (
	ErrUndoBadRunFlagFmt  = "Error: invalid --run value %q: %v (operation: parse-run-flag, reason: %s)"
	ErrUndoRunNotFoundFmt = "Error: run id %d not found in MakeAllVisibilityRun (operation: SQLSelectRunByID, reason: no row)"
	ErrRedoNoRunFound     = "Error: no undoable VisibilityUndo run found (operation: visibility-redo, reason: MakeAllVisibilityRun has no VisibilityUndo row with OkCount>0)"
)

// SQLSelectRecentRuns — newest-first list for `visibility-history`.
const SQLSelectRecentRuns = `SELECT
	MakeAllVisibilityRunId, CommandKind, TargetVisibility, Provider,
	Owner, MatchedCount, OkCount, SkippedCount, FailedCount,
	ExcludedCount, ExitCode, StartedAt, FinishedAt
	FROM MakeAllVisibilityRun
	ORDER BY MakeAllVisibilityRunId DESC
	LIMIT ?`

// SQLSelectRecentRunsBase — column projection + FROM for the
// filtered builder. WHERE / ORDER / LIMIT are appended dynamically
// by `BuildRecentRunsQuery` based on caller-supplied filters.
const SQLSelectRecentRunsBase = `SELECT
	MakeAllVisibilityRunId, CommandKind, TargetVisibility, Provider,
	Owner, MatchedCount, OkCount, SkippedCount, FailedCount,
	ExcludedCount, ExitCode, StartedAt, FinishedAt
	FROM MakeAllVisibilityRun`

// SQL fragments composed by BuildRecentRunsQuery. Centralized so
// the per-clause SQL never sits inline in store code (no magic strings).
const (
	SQLWhereCommandKindEq  = " CommandKind = ?"
	SQLWhereStartedAtGTE   = " StartedAt >= ?"
	SQLOrderRunIDDescLimit = " ORDER BY MakeAllVisibilityRunId DESC LIMIT ?"
	SQLKeywordWHERE        = " WHERE"
	SQLKeywordAND          = " AND"
)

const (
	ErrHistorySelectFmt = "Error: select recent runs failed: %v (operation: SQLSelectRecentRuns, reason: %s)"
	MsgVisHistoryEmpty  = "visibility-history: no make-all-* runs recorded yet\n"
	MsgVisHistoryHeader = "ID    Kind             Owner                 Matched  Ok  Skip Fail Excl Exit  Started\n"
	MsgVisHistoryRowFmt = "%-5d %-16s %-21s %7d %3d %4d %4d %4d %4d  %s\n"
	HistoryDefaultLimit = 20
)

// Dry-run messaging for `vu` / `vr` --dry-run.
const (
	MsgDryRunHeaderFmt = "%s --dry-run: would reverse run #%d (%s/%s) — %d repo(s):\n"
	MsgDryRunRowFmt    = "  %3d/%-3d %-40s : would set visibility -> %s\n"
	MsgDryRunFooterFmt = "%s --dry-run: no mutations performed (re-run without --dry-run to apply)\n"
)
