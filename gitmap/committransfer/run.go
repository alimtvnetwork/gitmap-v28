package committransfer

import (
	"fmt"
	"os"
)

// RunRight is the public entry point for `commit-right`. The caller
// has already resolved both endpoints and built the Options struct.
//
// Phase 1 (v3.76.0). Phases 2 + 3 (RunLeft, RunBoth) live in
// runleftboth.go and reuse runOneDirection — RunRight is now a thin
// wrapper kept for callers that pre-date the directional family.
func RunRight(sourceDir, targetDir string, opts Options) error {
	return runOneDirection(sourceDir, targetDir, opts)
}

// maybePush runs `git push` unless --no-push is set, the target is not
// a git repo, or there are no new commits. Returns true on success.
func maybePush(targetDir string, opts Options, newCount int) bool {
	if opts.NoPush || opts.NoCommit || newCount == 0 || opts.DryRun {
		return false
	}
	if _, err := pushHEAD(targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "%s push failed: %v\n", opts.LogPrefix, err)

		return false
	}

	return true
}
