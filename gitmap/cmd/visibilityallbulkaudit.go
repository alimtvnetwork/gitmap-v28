// Package cmd — visibilityallbulkaudit.go: thin wiring layer that
// persists each `make-all-*` invocation to MakeAllVisibilityRun +
// MakeAllVisibilityResult via the store helpers.
//
// All DB calls are best-effort: a missing/locked audit DB MUST NOT
// abort the user's bulk action. Every error is logged to os.Stderr
// with Code Red context (zero-swallow) but execution continues.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan steps 19-20.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/visibility"
)

// runAudit holds the live audit state for one bulk invocation. db may
// be nil when the audit DB could not be opened — every method on this
// struct must tolerate that and degrade to a no-op.
type runAudit struct {
	db        *store.DB
	runID     int64
	resultIDs map[string]int64 // RepoName → row PK for fast UPDATE
}

// RunID exposes the new audit row's primary key (read-only, zero
// when the audit DB was unreachable) so callers like `--json`
// emitters can include it in their wire output.
func (a *runAudit) RunID() int64 { return a.runID }

// beginRunAudit opens the audit DB, writes the run row + one Pending
// result per matched repo. Returns a non-nil *runAudit even when the
// DB is unreachable so callers can treat audit as fire-and-forget.
func beginRunAudit(ctx ownerContext, target, cmdName, patternsRaw string,
	flags bulkFlags, ownerTotal int, matches []visibility.MatchedRepo,
) *runAudit {
	audit := &runAudit{resultIDs: make(map[string]int64, len(matches))}
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: audit DB open failed: %v (continuing without audit)\n", err)

		return audit
	}
	audit.db = db
	audit.runID = insertRunRow(db, ctx, target, cmdName, patternsRaw, flags, ownerTotal, len(matches))
	persistPendingResults(audit, matches)

	return audit
}

// insertRunRow writes the pre-prompt MakeAllVisibilityRun row.
func insertRunRow(db *store.DB, ctx ownerContext, target, cmdName, patternsRaw string,
	flags bulkFlags, ownerTotal, matchedCount int,
) int64 {
	rec := model.MakeAllVisibilityRunRecord{
		CommandKind:      commandKindFor(cmdName),
		TargetVisibility: target,
		Provider:         ctx.Provider,
		Owner:            ctx.Owner,
		TargetRaw:        ctx.TargetRaw,
		PatternList:      patternsRaw,
		YesFlag:          flags.Yes,
		VerboseFlag:      flags.Verbose,
		OwnerRepoTotal:   ownerTotal,
		MatchedCount:     matchedCount,
		StartedAt:        nowRFC3339(),
	}
	id, err := db.InsertMakeAllVisibilityRun(rec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)

		return 0
	}

	return id
}

// persistPendingResults writes one Pending result row per matched repo
// and stores the assigned IDs on the audit struct keyed by RepoName.
func persistPendingResults(audit *runAudit, matches []visibility.MatchedRepo) {
	if audit.db == nil || audit.runID == 0 {
		return
	}
	rows := make([]model.MakeAllVisibilityResultRecord, 0, len(matches))
	for _, m := range matches {
		rows = append(rows, model.MakeAllVisibilityResultRecord{
			RepoName: m.RepoName, MatchedPattern: m.MatchedPattern,
			Status: constants.ResultStatusPending, StartedAt: nowRFC3339(),
		})
	}
	ids, err := audit.db.InsertMakeAllVisibilityPendingResults(audit.runID, rows)
	if err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)

		return
	}
	for i, m := range matches {
		audit.resultIDs[m.RepoName] = ids[i]
	}
}

// markExcluded flips the difference between `before` and `after` to
// Status='Excluded'. No-op when the user accepted the full set.
func (a *runAudit) markExcluded(before, after []visibility.MatchedRepo) int {
	if a.db == nil {
		return 0
	}
	keep := make(map[string]bool, len(after))
	for _, m := range after {
		keep[m.RepoName] = true
	}
	ids := make([]int64, 0, len(before)-len(after))
	for _, m := range before {
		if keep[m.RepoName] {
			continue
		}
		if id, ok := a.resultIDs[m.RepoName]; ok {
			ids = append(ids, id)
		}
	}
	if err := a.db.MarkMakeAllVisibilityResultsExcluded(ids, nowRFC3339()); err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)
	}

	return len(ids)
}

// updateResult writes the terminal status for one repo. start is
// captured by the caller so the persisted DurationMs reflects the
// actual provider-CLI round-trip.
func (a *runAudit) updateResult(repoName string, st applyStatus, prev, next string, start time.Time) {
	if a.db == nil {
		return
	}
	id, ok := a.resultIDs[repoName]
	if !ok {
		return
	}
	rec := model.MakeAllVisibilityResultRecord{
		ID: id, Status: statusForOutcome(st.outcome),
		PrevVisibility: prev, NewVisibility: next,
		FailureMessage: errString(st.err),
		FinishedAt:     nowRFC3339(),
		DurationMs:     time.Since(start).Milliseconds(),
	}
	if err := a.db.UpdateMakeAllVisibilityResult(rec); err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)
	}
}

// finalize flushes the tally + exit code back to the run row.
func (a *runAudit) finalize(excluded, ok, skipped, failed, exitCode int) {
	if a.db == nil || a.runID == 0 {
		return
	}
	rec := model.MakeAllVisibilityRunRecord{
		ID: a.runID, ExcludedCount: excluded, OkCount: ok,
		SkippedCount: skipped, FailedCount: failed, ExitCode: exitCode,
		FinishedAt: nowRFC3339(),
	}
	if err := a.db.FinalizeMakeAllVisibilityRun(rec); err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: %v\n", err)
	}
}

// commandKindFor maps the dispatcher's command ID to the persisted
// CommandKindEnum value.
func commandKindFor(cmdName string) string {
	switch cmdName {
	case constants.CmdMakeAllPrivate:
		return constants.CommandKindMakeAllPrivate
	case constants.CmdVisibilityUndo:
		return constants.CommandKindVisibilityUndo
	case constants.CmdVisibilityRedo:
		return constants.CommandKindVisibilityRedo
	}

	return constants.CommandKindMakeAllPublic
}

// statusForOutcome maps the apply loop's outcome token to the
// persisted ResultStatusEnum value.
func statusForOutcome(outcome string) string {
	switch outcome {
	case "ok":
		return constants.ResultStatusOk
	case "skip":
		return constants.ResultStatusSkipped
	}

	return constants.ResultStatusFailed
}

// nowRFC3339 returns the current UTC time formatted for SQLite TEXT
// columns (matches the convention used by store/archive_history.go).
func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
