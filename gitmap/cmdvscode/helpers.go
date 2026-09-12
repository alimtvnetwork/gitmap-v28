// Package cmdvscode provides VS Code integration, workspace, project manager sync, and duplicate detection commands.
package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/helptext"
)

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
}

func truncateStr(s string, maxLen int) string {
	if len(s) > maxLen && maxLen > 3 {
		return s[:maxLen-3] + "..."
	}

	return s
}
