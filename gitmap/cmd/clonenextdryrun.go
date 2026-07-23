package cmd

// `--dry-run` glue for `gitmap cn` (v3.132.0+).
//
// Split out of clonenext.go and clonenextbatch.go to keep both files
// under the 200-line per-file budget. The single-repo path calls
// printCloneNextDryRun + os.Exit(0) right before the actual
// runGitClone invocation; the batch path replaces its
// processBatchRepos pass with previewDryRunBatch which emits one
// preview line per repo and then exits.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// printCloneNextDryRun previews the single-repo clone and exits 0.
// Mirrors the exact `git clone <url> <dest>` invocation runGitClone
// would build, so users can copy/paste the printed line into a shell.
// Header + footer bracket the preview so it stands out in scrollback.
func printCloneNextDryRun(url, dest string) {
	fmt.Print(constants.MsgCloneNextDryRunHeader)
	fmt.Printf(constants.MsgCloneNextDryRunCmd,
		constants.GitBin, constants.GitClone, url, dest)
	fmt.Printf(constants.MsgCloneNextDryRunFooter, 1)
	maybeExitOnCmdFaithfulMismatch()
	os.Exit(0)
}

// previewDryRunBatch handles `--dry-run` for `cn --all` and
// `cn --csv`. For each repo we resolve the same target version the
// real batch path would (via processOneBatchRepo's helpers) so the
// printed clone command matches what an actual run produces — then
// we exit 0 without calling any of the side-effecting helpers.
//
// Failures during target resolution print as `[skip]` lines so the
// preview is honest about which repos would error out at run time.
func previewDryRunBatch(csvPath string, walkAll bool) {
	repos, err := loadBatchRepos(csvPath, walkAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextBatchLoad, err)
		os.Exit(1)
	}
	fmt.Print(constants.MsgCloneNextDryRunHeader)
	count := emitDryRunRows(repos)
	fmt.Printf(constants.MsgCloneNextDryRunFooter, count)
	maybeExitOnCmdFaithfulMismatch()
	os.Exit(0)
}

// emitDryRunRows iterates repos and prints one preview line per
// resolvable target. Returns the number of clone commands actually
// printed (skipped rows don't count toward the footer total).
func emitDryRunRows(repos []string) int {
	count := 0
	for _, repoPath := range repos {
		url, dest, ok := resolveDryRunTarget(repoPath)
		if !ok {
			fmt.Printf("  ⊘ %s — cannot resolve next version, would skip\n",
				filepath.Base(repoPath))
			continue
		}
		fmt.Printf(constants.MsgCloneNextDryRunCmd,
			constants.GitBin, constants.GitClone, url, dest)
		count++
	}

	return count
}

// resolveDryRunTarget computes the (url, dest) pair the real batch
// path would clone for a single repo. Returns ok=false on any
// resolution error so the caller can emit a skip line instead. No
// network or filesystem mutations — pure planning.
func resolveDryRunTarget(repoPath string) (url, dest string, ok bool) {
	parsed, _, err := readRepoVersion(repoPath)
	if err != nil {
		return "", "", false
	}
	target, err := clonenext.ResolveTarget(parsed, "v++")
	if err != nil {
		return "", "", false
	}
	state, err := clonenext.ReadLocalRepoState(repoPath)
	if err != nil || len(state.OriginURL) == 0 {
		return "", "", false
	}
	dest = filepath.Join(filepath.Dir(repoPath),
		fmt.Sprintf("%s-v%d", parsed.BaseName, target))

	return state.OriginURL, dest, true
}
