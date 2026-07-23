// Package cmd — visibilityhistory.go: `gitmap visibility-history`
// (`vh`) prints the most recent make-all-* / VisibilityUndo /
// VisibilityRedo runs newest-first so users can select a `--run <id>`
// for `vu` / `vr`.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §history.
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runVisibilityHistory is the dispatcher entry point. When `--kind`
// or `--since` is set, the filter is pushed down to SQL (step 39);
// the in-memory `applyHistoryFilters` (step 36) is retained as a
// defense-in-depth second pass that also drops malformed StartedAt
// rows the SQL `>=` comparison would let through as text-compare.
func runVisibilityHistory(args []string) {
	limit := parseHistoryLimit(args)
	filters := parseHistoryFilters(args, time.Now())
	db := openDBOrExit(constants.CmdVisibilityHistory)
	runs, err := loadHistoryRuns(db, limit, filters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(constants.ExitVisAuthFailed)
	}
	runs = applyHistoryFilters(runs, filters, time.Now())
	if len(runs) == 0 {
		fmt.Fprint(os.Stderr, constants.MsgVisHistoryEmpty)
		os.Exit(constants.ExitVisOK)
	}
	printHistory(runs)
	os.Exit(constants.ExitVisOK)
}

// loadHistoryRuns picks the pushdown path when any filter is set,
// otherwise falls back to the unfiltered newest-first SELECT.
func loadHistoryRuns(db *store.DB, limit int, f historyFilters) ([]model.MakeAllVisibilityRunRecord, error) {
	if f.Kind == "" && f.Since == 0 {
		return db.SelectRecentMakeAllVisibilityRuns(limit)
	}
	sinceISO := ""
	if f.Since > 0 {
		sinceISO = time.Now().Add(-f.Since).UTC().Format(time.RFC3339)
	}

	return db.SelectRecentMakeAllVisibilityRunsFiltered(store.RecentRunsFilter{
		Kind: f.Kind, SinceISO: sinceISO, Limit: limit,
	})
}

// parseHistoryLimit accepts `--limit N` (default HistoryDefaultLimit).
// Bad values exit ExitVisBadFlag (zero-swallow).
func parseHistoryLimit(args []string) int {
	for i := 0; i+1 < len(args); i++ {
		if args[i] != "--limit" {
			continue
		}
		n, err := strconv.Atoi(args[i+1])
		if err != nil || n <= 0 {
			fmt.Fprintf(os.Stderr, constants.ErrUndoBadRunFlagFmt, args[i+1], err, "--limit must be positive integer")
			os.Exit(constants.ExitVisBadFlag)
		}

		return n
	}

	return constants.HistoryDefaultLimit
}

// printHistory writes the table to stdout.
func printHistory(runs []model.MakeAllVisibilityRunRecord) {
	fmt.Fprint(os.Stdout, constants.MsgVisHistoryHeader)
	for _, r := range runs {
		owner := truncateCell(r.Provider+"/"+r.Owner, 21)
		fmt.Fprintf(os.Stdout, constants.MsgVisHistoryRowFmt,
			r.ID, truncateCell(r.CommandKind, 16), owner,
			r.MatchedCount, r.OkCount, r.SkippedCount, r.FailedCount,
			r.ExcludedCount, r.ExitCode, r.StartedAt)
	}
}

// truncate clips overflowing cells to keep the table aligned.
func truncateCell(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n-1] + "…"
}
