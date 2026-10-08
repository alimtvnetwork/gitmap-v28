package cmdcompletion

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunCompletion(args []string) error {
	checkHelp("completion", args)
	if hasListFlag(args) {
		handleCompletionList(args)

		return nil
	}
	if len(args) > 0 && args[0] == "install" {
		return handleCompletionInstall(args[1:])
	}

	return executeCompletionScript(args)
}

func handleCompletionInstall(args []string) error {
	shell := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		shell = args[0]
	} else {
		shell = completion.DetectShell()
	}

	if err := completion.Install(shell); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to install %s completion: %v\n", shell, err)
		return err
	}

	fmt.Printf("  ✓ GitMap completion and prediction suggestions installed for %s in profile.\n", shell)
	return nil
}

func executeCompletionScript(args []string) error {
	shell := ""
	if len(args) >= 1 {
		shell = args[0]
	} else {
		shell = completion.DetectShell()
	}

	printCompletionScript(shell)

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
		cliexit.HandleError(err, 1)
	}

	fmt.Print(script)
}
