package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func runInteractiveClone(cmd *exec.Cmd, rec model.ScanRecord, url, dest string,
	strategy cloneStrategy) model.CloneResult {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("git clone failed for %s: url=%q branch=%q dest=%q: %v",
			recordTag(rec), url, strategy.branch, dest, err)
		return model.CloneResult{Record: rec, Success: false, Error: msg, Notes: strategy.reason}
	}

	return model.CloneResult{Record: rec, Success: true, Notes: strategy.reason}
}

func isSSHCloneURL(url string) bool {
	lower := strings.ToLower(strings.TrimSpace(url))
	return strings.HasPrefix(lower, constants.PrefixSSH) ||
		strings.HasPrefix(lower, constants.PrefixSSHScheme)
}
