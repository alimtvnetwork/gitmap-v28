// Package cmd — clihelpers_argv.go: generic argv preprocessing predicates.
//
// Kept from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// These are cross-cutting dispatch utilities used by root*.go and many
// command files; argv preprocessing explicitly stays in package cmd.
package cmd

import (
	"slices"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isHelpFlag(arg string) bool {
	return slices.Contains([]string{constants.FlagHelpWord, constants.FlagHelp, constants.FlagHelpShort}, arg)
}
