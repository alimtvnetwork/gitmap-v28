// Package cmd — clihelpers_argv.go: generic argv preprocessing predicates.
//
// Kept from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// These are cross-cutting dispatch utilities used by root*.go and many
// command files; argv preprocessing explicitly stays in package cmd.
package cmd

func isHelpFlag(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}
