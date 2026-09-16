package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runCompletion(args []string) error {
	checkHelp("completion", args)
	if hasListFlag(args) {
		handleCompletionList(args)

		return nil
	}

	return executeCompletionScript(args)
}

func executeCompletionScript(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrCompUsage)
		cliexit.HandleError(nil, 1)
	}

	printCompletionScript(args[0])

	return nil
}

func hasListFlag(args []string) bool {
	for _, a := range args {
		if isCompletionListToken(a) {
			return true
		}
	}

	return false
}

func isCompletionListToken(a string) bool {
	return a == constants.CompListRepos || a == constants.CompListGroups ||
		a == constants.CompListCommands || a == constants.CompListAliases ||
		a == constants.CompListZipGroups || a == constants.CompListSSHKeys ||
		a == constants.CompListHelpGroups
}

func handleCompletionList(args []string) {
	printers := completionPrinters()
	for _, a := range args {
		if fn, ok := printers[a]; ok {
			fn()

			return
		}
	}
}

func completionPrinters() map[string]func() {
	return map[string]func(){
		constants.CompListRepos:      printCompletionRepos,
		constants.CompListGroups:     printCompletionGroups,
		constants.CompListCommands:   printCompletionCommands,
		constants.CompListAliases:    printCompletionAliases,
		constants.CompListZipGroups:  printCompletionZipGroups,
		constants.CompListSSHKeys:    printCompletionSSHKeys,
		constants.CompListHelpGroups: printCompletionHelpGroups,
	}
}

func printCompletionScript(shell string) {
	script, err := completion.Generate(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCompUnknownShell, shell)
		cliexit.HandleError(nil, 1)
	}

	fmt.Print(script)
}
