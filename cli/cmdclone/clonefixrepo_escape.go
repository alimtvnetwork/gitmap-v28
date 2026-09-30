package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func findEnclosingGitRoot(start string) string {
	curr := start
	for hop := 0; hop < 32; hop++ {
		if clonenext.IsGitRepo(curr) {
			return curr
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return ""
}

func findFirstNonRepoAncestor(repoRoot string) string {
	target := filepath.Dir(repoRoot)
	for hop := 0; hop < 32; hop++ {
		if !clonenext.IsGitRepo(target) {
			break
		}
		parent := filepath.Dir(target)
		if parent == target {
			break
		}
		target = parent
	}
	return target
}

// escapeNestedGitRepo walks up from the current working directory to find
// any enclosing git repo, chdir'ing into the first non-repo ancestor before
// the clone step runs.
func escapeNestedGitRepo() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	root := findEnclosingGitRoot(cwd)
	if root == "" {
		return
	}

	target := findFirstNonRepoAncestor(root)
	if target == cwd {
		return
	}

	if chErr := os.Chdir(target); chErr != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCFREscapeChdir, target, chErr)

		return
	}

	fmt.Fprintf(os.Stderr, constants.MsgCFREscapeNested, cwd, target)
}
