// Package cmdinstaller — helpers.go provides shared utilities for installer subcommands.
package cmdinstaller

import (
	"github.com/spf13/cobra"
)

// Subcommands returns all subcommands registered under installer.
func Subcommands() []*cobra.Command {
	if installerCmd == nil {
		return nil
	}

	return installerCmd.Commands()
}

func resolveProfileTreeInternal(slug string) (any, bool) {
	if ResolveProfileTreeFn != nil {
		return ResolveProfileTreeFn(slug)
	}

	return nil, false
}

func printProfileTreeInternal(profile any) {
	if PrintProfileTreeFn != nil {
		PrintProfileTreeFn(profile)
	}
}

func printProfileInstallSummaryInternal(slug string) {
	if PrintProfileInstallSummaryFn != nil {
		PrintProfileInstallSummaryFn(slug)
	}
}

var separateFlagAndPositionalArgs = SeparateFlagAndPositionalArgs
