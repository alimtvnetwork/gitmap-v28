package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenext"
)

// shouldRunBatch decides whether `gitmap cn` should fan out across
// multiple repos. Three triggers, evaluated in priority order:
//
//  1. Explicit `--csv <path>` flag.
//  2. Explicit `--all` flag.
//  3. Implicit: cwd has no `.git` entry but at least one immediate
//     child directory IS a git repo.
//
// The implicit check (#3) is intentionally cheap: one Stat on the
// cwd's `.git`, then a single short-circuiting ReadDir scan one level
// down (HasGitSubdir bails on the first git child it finds). It does
// NOT walk the full tree — the real walk happens later in
// runCloneNextBatch via clonenext.WalkBatchFromDir, where
// ErrBatchEmpty surfaces if a race emptied the directory between the
// trigger check and the walk.
//
// `cwd` is passed in (rather than re-fetched here) so the caller can
// reuse the value for downstream operations and so this function stays
// trivially testable with a fixture path. An empty `cwd` string skips
// the implicit check — the dispatcher then falls through to the
// single-repo path, which already prints a clear "no remote" error
// when the user really is sitting in a non-repo directory.
func shouldRunBatch(flags CloneNextFlags, cwd string) bool {
	if len(flags.CSVPath) > 0 || flags.All {
		return true
	}
	if len(cwd) == 0 {
		return false
	}
	if clonenext.IsGitRepo(cwd) {
		return false
	}

	return clonenext.HasGitSubdir(cwd)
}

// currentWorkingDir is a thin wrapper over os.Getwd that returns ""
// on error. The dispatcher treats "" as "skip the implicit trigger"
// rather than failing the whole command, because the single-repo path
// produces a clearer error message in the rare case where cwd is
// genuinely unreadable.
func currentWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return cwd
}
