// Package cmdvscode provides VS Code integration, workspace, project manager sync, and duplicate detection commands.
package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

func isHelpToken(a string) bool {
	return a == "--help" || a == "-h" || a == "help"
}

func printVSCodeHelpAndExit(command string) {
	if command == "vscode" || command == "code" {
		RenderVSCodeHelp()
		cliexit.Exit(0)
	}

	helptext.Print(command)
	cliexit.Exit(0)
}

func checkHelp(command string, args []string) {
	for _, a := range args {
		if isHelpToken(a) {
			printVSCodeHelpAndExit(command)
		}
	}
}

func truncateStr(s string, maxLen int) string {
	if len(s) > maxLen && maxLen > 3 {
		return s[:maxLen-3] + "..."
	}

	return s
}
