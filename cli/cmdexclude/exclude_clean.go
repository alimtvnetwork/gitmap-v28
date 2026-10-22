package cmdexclude

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunCleanRepos inspects all repositories in the GitMap store and removes non-owned, cache, or corrupted entries.
func RunCleanRepos(args []string) error {
	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "failed to open database")
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil {
		return apperror.WrapSimple(errList, "failed to list tracked repositories")
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║               GitMap Repository Cleanup & Purge                     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")

	var purged []model.ScanRecord
	for _, r := range repos {
		reason := evaluatePurgeReason(r)
		if reason != "" {
			_, delErr := db.DeleteByPath(r.AbsolutePath)
			if delErr == nil {
				purged = append(purged, r)
				fmt.Printf("  ✗ Purged: %-30s [%s]\n    Path: %s\n", r.Slug, reason, r.AbsolutePath)
			}
		}
	}

	if len(purged) == 0 {
		fmt.Printf("✓ Repository registry is clean (%d repositories tracked, 0 non-owned/cache repos found).\n", len(repos))
	} else {
		fmt.Printf("\n✓ Cleaned %d non-owned / cache / invalid repositories from registry (%d active remaining).\n",
			len(purged), len(repos)-len(purged))
	}

	return nil
}

func evaluatePurgeReason(r model.ScanRecord) string {
	if r.AbsolutePath == "" {
		return "empty path"
	}

	if _, err := os.Stat(r.AbsolutePath); err != nil {
		return "path does not exist"
	}

	if fsutil.IsExcludedScanPath(r.AbsolutePath) {
		return "matches exclusion pattern"
	}

	if !fsutil.IsValidGitRepo(r.AbsolutePath) {
		return "invalid git repository structure"
	}

	lower := strings.ToLower(r.RepoName)
	if (lower == "dst" || lower == "src" || lower == "temp" || lower == "test") && (r.HTTPSUrl == "" && r.SSHUrl == "") {
		return "scratch test directory without remote"
	}

	return ""
}
