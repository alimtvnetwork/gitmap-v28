package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// InteractiveCloneParams encapsulates arguments for running an interactive clone.
type InteractiveCloneParams struct {
	Cmd      *exec.Cmd
	Record   model.ScanRecord
	URL      string
	Dest     string
	Strategy cloneStrategy
}

func runInteractiveClone(params InteractiveCloneParams) model.CloneResult {
	params.Cmd.Stdin = os.Stdin
	params.Cmd.Stdout = os.Stdout
	params.Cmd.Stderr = os.Stderr
	if err := params.Cmd.Run(); err != nil {
		msg := fmt.Sprintf("git clone failed for %s: url=%q branch=%q dest=%q: %v",
			recordTag(params.Record), params.URL, params.Strategy.branch, params.Dest, err)

		return model.CloneResult{Record: params.Record, IsSuccess: false, Error: msg, Notes: params.Strategy.reason}
	}

	return model.CloneResult{Record: params.Record, IsSuccess: true, Notes: params.Strategy.reason}
}

func isSSHCloneURL(url string) bool {
	lower := strings.ToLower(strings.TrimSpace(url))

	return strings.HasPrefix(lower, constants.PrefixSSH) ||
		strings.HasPrefix(lower, constants.PrefixSSHScheme)
}
