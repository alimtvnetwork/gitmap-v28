// Package store — makeallvisibility_history_filtered.go: step-39
// SQL-side filter pushdown for `vh`. Avoids loading every
// MakeAllVisibilityRun row into memory when the user passed
// `--kind` and/or `--since` — at thousands of historical runs the
// in-memory filter (step 36 fallback) becomes the bottleneck.
//
// The pure query builder is exported so cmd-side tests can lock the
// WHERE composition without spinning up SQLite.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §history.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// RecentRunsFilter is the pushdown shape — empty fields disable
// their clause (zero value is equivalent to "newest N, no filter").
type RecentRunsFilter struct {
	Kind     string // exact CommandKind match
	SinceISO string // ISO-8601 lower bound (>=)
	Limit    int    // required, must be > 0
}

// BuildRecentRunsQuery composes SQL + positional args for the filter.
// Returns (sql, args). Pure — no DB handle, fully unit-testable.
func BuildRecentRunsQuery(f RecentRunsFilter) (string, []any) {
	sql := constants.SQLSelectRecentRunsBase
	args := make([]any, 0, 3)
	clauses := make([]string, 0, 2)
	if f.Kind != "" {
		clauses = append(clauses, constants.SQLWhereCommandKindEq)
		args = append(args, f.Kind)
	}
	if f.SinceISO != "" {
		clauses = append(clauses, constants.SQLWhereStartedAtGTE)
		args = append(args, f.SinceISO)
	}
	for i, c := range clauses {
		if i == 0 {
			sql += constants.SQLKeywordWHERE + c
			continue
		}
		sql += constants.SQLKeywordAND + c
	}
	sql += constants.SQLOrderRunIDDescLimit
	args = append(args, f.Limit)

	return sql, args
}

// SelectRecentMakeAllVisibilityRunsFiltered runs the built query
// and projects rows through the existing scanner (column order is
// identical to SQLSelectRecentRuns so scanRecentRuns is reused).
func (db *DB) SelectRecentMakeAllVisibilityRunsFiltered(f RecentRunsFilter) ([]model.MakeAllVisibilityRunRecord, error) {
	sql, args := BuildRecentRunsQuery(f)
	rows, err := db.conn.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrHistorySelectFmt, err, err.Error())
	}
	defer rows.Close()

	return scanRecentRuns(rows)
}
