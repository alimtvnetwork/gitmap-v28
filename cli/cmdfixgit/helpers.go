package cmdfixgit

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

func checkHelp(command string, args []string) {
	if !hasHelpFlag(args) {
		return
	}
	helptext.Print(command)
	cliexit.Exit(0)
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}
	return false
}

func repoRoot() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		cliexit.HandleError(err, 1)
	}
	return filepath.Clean(filepath.FromSlash(strings.TrimSpace(string(out))))
}
