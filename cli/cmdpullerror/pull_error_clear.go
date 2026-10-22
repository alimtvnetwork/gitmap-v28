package cmdpullerror

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func executeClearPullErrors(opts PullErrorOptions) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "open pull split db")
	}
	defer db.Close()

	if err := db.ClearPullErrors(opts.RepoSlug); err != nil {
		return apperror.WrapSimple(err, "clear pull errors")
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
