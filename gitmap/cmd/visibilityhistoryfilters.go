// Package cmd — visibilityhistoryfilters.go: step-36 post-fetch
// filters for `vh`. Pure helpers (no DB, no I/O) so the filter
// policy is unit-testable in isolation. SQL-side filtering is a
// follow-up if `vh` ever paginates beyond MaxFilterBacklog rows.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §history.
package cmd

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// historyFilters holds the parsed `--kind <K>` / `--since <dur>` pair.
// Zero value = no filtering (both branches are no-ops).
type historyFilters struct {
	Kind  string        // exact CommandKind match; "" disables
	Since time.Duration // include only StartedAt >= now-Since; 0 disables
}

// parseHistoryFilters extracts `--kind <K>` and `--since <dur>` from
// args. Unknown tokens are ignored (limit is parsed elsewhere).
func parseHistoryFilters(args []string, now time.Time) historyFilters {
	var f historyFilters
	for i := 0; i+1 < len(args); i++ {
		switch args[i] {
		case "--kind":
			f.Kind = strings.TrimSpace(args[i+1])
		case "--since":
			if d, err := time.ParseDuration(args[i+1]); err == nil && d > 0 {
				f.Since = d
			}
		}
	}
	_ = now // accepted for symmetry with applyHistoryFilters

	return f
}

// applyHistoryFilters returns the subset of `runs` matching f. Stable,
// preserves input order (caller pre-sorts newest-first).
func applyHistoryFilters(runs []model.MakeAllVisibilityRunRecord, f historyFilters, now time.Time) []model.MakeAllVisibilityRunRecord {
	if f.Kind == "" && f.Since == 0 {
		return runs
	}
	cutoff := now.Add(-f.Since)
	out := make([]model.MakeAllVisibilityRunRecord, 0, len(runs))
	for _, r := range runs {
		if f.Kind != "" && r.CommandKind != f.Kind {
			continue
		}
		if f.Since > 0 {
			ts, err := time.Parse(time.RFC3339, r.StartedAt)
			if err != nil || ts.Before(cutoff) {
				continue
			}
		}
		out = append(out, r)
	}

	return out
}
