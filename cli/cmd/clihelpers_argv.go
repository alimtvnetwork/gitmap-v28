// Package cmd — clihelpers_argv.go: generic argv preprocessing predicates.
//
// Kept from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// These are cross-cutting dispatch utilities used by root*.go and many
// command files; argv preprocessing explicitly stays in package cmd.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func isExistingFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isFileFlagWithArg(arg string) bool {
	return arg == "--file" || arg == "--out" || arg == "-o" || arg == "-f"
}

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}

func extractFlagVal(args []string, flagName string) string {
	for i, arg := range args {
		if arg == flagName && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(arg, flagName+"=") {
			return strings.TrimPrefix(arg, flagName+"=")
		}
	}

	return ""
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func isHelpFlag(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}
