package cmdpullerror

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func executeClearPullErrors(opts PullErrorOptions) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return fmt.Errorf("open pull split db: %w", err)
	}
	defer db.Close()

	if err := db.ClearPullErrors(opts.RepoSlug); err != nil {
		return fmt.Errorf("clear pull errors: %w", err)
	}

	renderClearPullErrorsSuccess(opts.RepoSlug)

	return nil
}

func renderClearPullErrorsSuccess(target string) {
	if target == "all" || target == "" {
		fmt.Println("  ✓ Cleared all recorded pull errors from gitmap-pull.db.")

		return
	}

	fmt.Printf("  ✓ Cleared recorded pull errors for '%s' from gitmap-pull.db.\n", target)
}
