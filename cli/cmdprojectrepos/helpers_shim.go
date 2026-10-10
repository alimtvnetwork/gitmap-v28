package cmdprojectrepos

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"os"
	"strings"
)

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}

func printHints(hints []hintEntry) {
	fmt.Fprint(os.Stderr, constants.MsgHintHeader)
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, constants.MsgHintRowFmt, h.command, h.description)
	}
}

// projectReposHints returns hints shown after go-repos, node-repos, etc.
func projectReposHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupAdd, constants.HintGroupAddDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
		{constants.HintPullGroup, constants.HintPullGroupDesc},
	}
}

type hintEntry struct {
	command     string
	description string
}
