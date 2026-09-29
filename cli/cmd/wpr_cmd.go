// Package cmd — wpr_cmd.go exposes top-level watch-prompts-running CLI command delegation.
package cmd

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"

// RunWPR delegates execution directly to cmdagy.RunWPRCLI.
func RunWPR(args []string) error {
	return cmdagy.RunWPRCLI(args)
}
