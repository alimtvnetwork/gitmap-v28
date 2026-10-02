package cmdpull

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// OptimizeRedundantRepos checks SQLite for duplicate repositories (both remote URL and path-level) and cleans them before pulling.
func OptimizeRedundantRepos(db *store.DB, isQuiet bool) (*store.DeduplicationSummary, error) {
	if db == nil {
		return nil, nil
	}

	summary, err := db.DeduplicateRepos(false)
	if err != nil {
		return nil, err
	}

	if !isQuiet && summary != nil && summary.RowsPurged > 0 {
		fmt.Fprintf(os.Stderr, "  %s✓ SQLite: Optimized %d redundant repository record(s) across %d group(s) before pull.%s\n\n",
			constants.ColorGreen, summary.RowsPurged, summary.GroupsFound, constants.ColorReset)
	}

	return summary, nil
}
