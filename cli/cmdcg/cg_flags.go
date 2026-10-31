// Package cmdcg - cg_flags.go: help-arg detection and flag setup for cg.
// Split from cg.go on 31-Oct-26 so cg.go stays under the 300-line file-size gate.
package cmdcg

import (
	"flag"
	"slices"
)

func isCGHelpArg(arg string) bool {
	return slices.Contains([]string{"help", "-h", "--help"}, arg)
}

func hasCGHelpArg(args []string) bool {
	for _, arg := range args {
		if isCGHelpArg(arg) {
			return true
		}
	}

	return false
}

func initCGFlagSet(fs *flag.FlagSet, opts *cgOptions) {
	fs.BoolVar(&opts.All, "all", false, "Run on all workspaces")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Simulate execution without modifying files")
	fs.StringVar(&opts.Exclude, "except", "", "Exclude repos (comma separated)")
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude repos (comma separated)")
	fs.StringVar(&opts.VersionValue, "version", "", "Starting version value for version.json (default: 1.0.0)")
	fs.StringVar(&opts.VersionValue, "version-value", "", "Starting version value for version.json (default: 1.0.0)")
}
