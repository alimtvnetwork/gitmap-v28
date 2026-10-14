package cmdreconcile

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

var promptChoiceMap = map[string]string{
	"1":         "stash",
	"stash":     "stash",
	"2":         "wip",
	"wip":       "wip",
	"3":         "discard",
	"discard":   "discard",
	"s":         "skip",
	"skip":      "skip",
	"a":         "all-stash",
	"all":       "all-stash",
	"all-stash": "all-stash",
}

func resolveDirtyFiles(item *cmdremediation.RemediationItem) []string {
	if len(item.RepoPath) == 0 {
		return item.Files
	}

	diag := gitutil.InspectDirtyState(item.RepoPath)
	if len(diag.AllFiles) > 0 {
		item.Files = diag.AllFiles
		return diag.AllFiles
	}

	return item.Files
}

func resolvePromptChoice(choice string) (string, bool) {
	isQuit := choice == "q" || choice == "quit" || choice == "exit"
	if isQuit {
		return "", true
	}

	action, hasAction := promptChoiceMap[choice]
	if hasAction {
		return action, false
	}

	return "stash", false
}
