package cmdcluster

import (
	"slices"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// HasHelpFlagFn is wired by cmd/di_hooks.go.
var HasHelpFlagFn func(args []string) bool

func hasHelpFlag(args []string) bool {
	if HasHelpFlagFn != nil {
		return HasHelpFlagFn(args)
	}
	for _, arg := range args {
		if slices.Contains(clusterHelpFlagTokens, arg) {
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
		if slices.Contains(clusterHelpFlagTokens, arg) {
			cliexit.Exit(0)
			return
		}
	}
}
