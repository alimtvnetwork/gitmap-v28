package cmdhaschange

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// isWorkingTreeDirty reports whether `git status --porcelain` returns content.
func isWorkingTreeDirty(target string) (bool, error) {
	cmd := exec.Command(constants.GitBin, "status", "--porcelain")
	cmd.Dir = target

	out, err := cmd.Output()
	if err != nil {
		return false, apperror.WrapSimple(err, "git status porcelain in "+target)
	}

	isDirty := len(strings.TrimSpace(string(out))) > 0

	return isDirty, nil
}
