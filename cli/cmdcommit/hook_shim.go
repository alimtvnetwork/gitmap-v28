package cmdcommit

import "github.com/alimtvnetwork/gitmap-v28/cli/cliexit"

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
		return
	}
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			cliexit.Exit(0)
			return
		}
	}
}
