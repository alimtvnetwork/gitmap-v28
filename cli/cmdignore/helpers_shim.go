package cmdignore

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/templates"
)

func insideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func gitTopLevel() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func outcomeVerb(o templates.MergeOutcomeType) string {
	switch o {
	case templates.MergeCreated:
		return "created"
	case templates.MergeInserted:
		return "inserted"
	case templates.MergeUpdated:
		return "updated"
	default:
		return "processed"
	}
}
