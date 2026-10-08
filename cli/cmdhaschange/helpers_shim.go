package cmdhaschange

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"os/exec"
	"strings"
)

// isWorkingTreeDirty reports whether `git status --porcelain` returns content.
func isWorkingTreeDirty(target string) bool {
	cmd := exec.Command(constants.GitBin, "status", "--porcelain")
	cmd.Dir = target

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

