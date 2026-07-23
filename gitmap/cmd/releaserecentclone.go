// Package cmd — releaserecentclone.go: the auto-cd-into-most-recent-
// clone fallback for `gitmap release`.
//
// When the user runs `gitmap r vX.Y.Z` from a parent directory that
// is itself NOT a git repo (typical right after `gitmap clone` /
// `cn` / `cfrp`), we look up the most recently cloned repo in the
// SQLite DB and chdir into it before delegating to the normal
// release pipeline. After the release returns we chdir back so the
// shell's working directory is unchanged.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// tryRunReleaseInRecentClone attempts the most-recent-clone fallback.
// Returns true when it owned the command (success OR explicit failure
// after chdir), false when the caller should continue to the next
// fallback (no recent clone recorded, or the recorded path is gone).
func tryRunReleaseInRecentClone(args []string) bool {
	target, ok := lookupRecentClone()
	if !ok {
		return false
	}

	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve cwd: %v\n", err)

		return false
	}

	if err := os.Chdir(target.AbsolutePath); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not switch to %s: %v\n", target.AbsolutePath, err)

		return false
	}

	fmt.Printf(constants.MsgReleaseAutoCdRecent, target.AbsolutePath, target.ClonedAt)

	defer restoreDir(originalDir)

	if !release.IsInsideGitRepo() {
		fmt.Fprintf(os.Stderr, "  Warning: %s is no longer a git repo\n", target.AbsolutePath)

		return false
	}

	runRelease(args)

	return true
}

// lookupRecentClone reads the most recent clone from the DB. Soft-fails
// (returns ok=false) on every error path so the caller can fall through.
func lookupRecentClone() (target struct {
	AbsolutePath string
	ClonedAt     string
}, ok bool) {
	db, err := openDB()
	if err != nil {
		return target, false
	}
	defer db.Close()

	rec, found, err := db.MostRecentClone()
	if err != nil || !found {
		return target, false
	}

	if _, statErr := os.Stat(rec.AbsolutePath); statErr != nil {
		return target, false
	}

	target.AbsolutePath = rec.AbsolutePath
	target.ClonedAt = rec.ClonedAt

	return target, true
}

// restoreDir chdir's back to originalDir, warning on failure.
func restoreDir(originalDir string) {
	if err := os.Chdir(originalDir); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not return to %s: %v\n", originalDir, err)

		return
	}
	fmt.Printf(constants.MsgReleaseAutoCdReturn, originalDir)
}
