package cmdopen

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func gitTopLevel() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("empty top-level")
	}

	return filepath.Clean(root), nil
}
