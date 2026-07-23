package cmd

// clonenextterm.go — `--output terminal` adapter for `gitmap cn`.
// Bridges the dispatcher (clonenext.go) and the shared
// render.RenderRepoTermBlock helper. Kept in its own file so the
// dispatcher stays focused on flow control and the rendering details
// (branch detection, fallback URLs, command shape) live alongside the
// other render-adjacent helpers.

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
)

// maybePrintCloneNextTermBlock emits the standardized RepoTermBlock
// for `gitmap cn` ONLY when `--output terminal` was passed. Any other
// value (including the empty default) is a no-op so existing CI logs
// and screenshots stay byte-identical.
//
// originalURL is the cn source repo's discovered remote (HTTPS/SSH).
// targetURL is the rewritten next-version URL that will actually be
// cloned. The standardized block surfaces both so the user can spot
// any URL-rewrite surprises before the clone starts.
func maybePrintCloneNextTermBlock(flags CloneNextFlags, name, branch, originalURL, targetURL, dest string) {
	if flags.Output != constants.OutputTerminal {
		return
	}
	block := render.RepoTermBlock{
		Index:        1,
		Name:         name,
		Branch:       branch,
		BranchSource: branchSourceLabel(branch),
		OriginalURL:  originalURL,
		TargetURL:    targetURL,
		CloneCommand: fmt.Sprintf("%s %s %s %s", constants.GitBin, constants.GitClone, targetURL, dest),
	}
	_ = render.RenderRepoTermBlock(os.Stdout, block)
}

// currentBranch shells out to `git -C <dir> rev-parse --abbrev-ref HEAD`
// and returns the branch name. Returns "" on any error so the caller
// (and downstream RenderRepoTermBlock) renders the "(unknown)"
// placeholder rather than failing the whole command.
func currentBranch(dir string) string {
	cmd := exec.Command(constants.GitBin, "-C", dir, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// branchSourceLabel returns "HEAD" for any non-empty branch (cn always
// reads from HEAD via rev-parse) and "" for empty so the renderer
// drops the parenthesized source segment entirely.
func branchSourceLabel(branch string) string {
	if len(strings.TrimSpace(branch)) == 0 {
		return ""
	}

	return "HEAD"
}
