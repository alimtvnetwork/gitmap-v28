package clonepick

// picker_tree.go: list every tracked path in the target repo so the
// --ask picker has something to render. Uses `git ls-tree -r
// --name-only HEAD` against a throwaway partial clone (--no-checkout
// --depth=1 --filter=blob:none) so we never download blob contents
// just to enumerate filenames.
//
// The temp clone is NOT reused for the final sparse-checkout in v1
// (the executor re-clones into the user's --dest). Spec §"--ask
// picker" calls for single-clone reuse as a follow-up optimisation;
// today's behavior is "clone twice, but the first clone is metadata
// only so the bytes-on-the-wire cost is negligible".

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ListRepoPathsKeep is the keep-clone variant: returns the path
// list AND the temp dir that holds the metadata clone. The caller
// owns tmp and must remove it (directly or by promoting it to a
// final destination via PromotePreClonedSrc). On error tmp is
// already cleaned up.
func ListRepoPathsKeep(plan Plan) ([]string, string, error) {
	tmp, err := os.MkdirTemp("", "gitmap-clonepick-ls-")
	if err != nil {
		return nil, "", fmt.Errorf(constants.ErrClonePickPickerLaunch, err)
	}
	if err := metaCloneForListing(plan, tmp); err != nil {
		os.RemoveAll(tmp)

		return nil, "", err
	}
	paths, lsErr := runLsTree(tmp)
	if lsErr != nil {
		os.RemoveAll(tmp)

		return nil, "", lsErr
	}

	return paths, tmp, nil
}

// metaCloneForListing runs the cheapest possible clone (filter=blob,
// no checkout, optional --branch) so ls-tree can walk HEAD without
// pulling any file contents.
func metaCloneForListing(plan Plan, dest string) error {
	args := []string{
		"clone", "--filter=blob:none", "--no-checkout", "--depth=1",
	}
	if len(plan.Branch) > 0 {
		args = append(args, "--branch", plan.Branch)
	}
	args = append(args, plan.RepoUrl, dest)
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(constants.ErrClonePickGitClone, err)
	}

	return nil
}

// runLsTree shells out to `git ls-tree -r --name-only HEAD` inside
// dest and returns one entry per non-empty output line. An empty
// repo (no HEAD) returns ([], nil) so the picker can render the
// "nothing to pick" branch instead of erroring out.
func runLsTree(dest string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", "HEAD")
	cmd.Dir = dest
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(constants.ErrClonePickGitLsTree, err)
	}

	return splitNonEmptyLines(buf.String()), nil
}

// splitNonEmptyLines trims trailing newlines and drops blank entries
// so a stray "" never becomes a phantom row in the picker.
func splitNonEmptyLines(raw string) []string {
	out := make([]string, 0, 32)
	for _, line := range strings.Split(raw, "\n") {
		clean := strings.TrimSpace(line)
		if len(clean) == 0 {
			continue
		}
		out = append(out, clean)
	}

	return out
}

// IsAutoExcluded reports whether path lives under any directory in
// the auto-exclude list (e.g. node_modules/, vendor/). Exported so
// the picker view layer can dim the row without re-implementing the
// match. Comparison is case-sensitive on purpose -- git paths are
// case-sensitive on every platform that matters here.
func IsAutoExcluded(path string) bool {
	for _, prefix := range constants.ClonePickAutoExclude {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}

	return false
}
