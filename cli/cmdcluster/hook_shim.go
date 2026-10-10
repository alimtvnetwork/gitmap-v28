package cmdcluster

import "github.com/alimtvnetwork/gitmap-v28/cli/cliexit"

// CheckHelpFn is wired by cmd/di_hooks.go.
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

// HasHelpFlagFn is wired by cmd/di_hooks.go.
var HasHelpFlagFn func(args []string) bool

func hasHelpFlag(args []string) bool {
	if HasHelpFlagFn != nil {
		return HasHelpFlagFn(args)
	}
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}
	return false
}

// CheckHelpOrEmptyFn is wired by cmd/di_hooks.go.
var CheckHelpOrEmptyFn func(command string, args []string)

func CheckHelpOrEmpty(command string, args []string) {
	if CheckHelpOrEmptyFn != nil {
		CheckHelpOrEmptyFn(command, args)
		return
	}
	if len(args) == 0 {
		cliexit.Exit(0)
		return
	}
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			cliexit.Exit(0)
			return
		}
	}
}
