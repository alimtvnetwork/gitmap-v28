package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/diff"
)

// runDiff implements `gitmap diff LEFT RIGHT`: a read-only preview
// of what `gitmap merge-*` would change between two folders.
//
// Spec: companion to spec/01-app/97-move-and-merge.md
func runDiff(args []string) {
	checkHelp("diff", args)
	left, right, walkOpts, printOpts := parseDiffArgs(args)

	leftEP, err := diff.ResolveEndpoint(left)
	if err != nil {
		cliexit.Fail(constants.CmdDiff, "resolve-endpoint", left, err, 1)
	}
	rightEP, err := diff.ResolveEndpoint(right)
	if err != nil {
		cliexit.Fail(constants.CmdDiff, "resolve-endpoint", right, err, 1)
	}
	if guardErr := guardDiffPaths(leftEP, rightEP); guardErr != nil {
		cliexit.Fail(constants.CmdDiff, "guard-paths", left+" vs "+right, guardErr, 1)
	}

	entries, err := diff.DiffTrees(leftEP.WorkingDir, rightEP.WorkingDir, walkOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s diff failed: %v\n", constants.LogPrefixDiff, err)
		os.Exit(1)
	}
	if reportErr := diff.Report(os.Stdout, entries, printOpts); reportErr != nil {
		fmt.Fprintf(os.Stderr, "%s report failed: %v\n", constants.LogPrefixDiff, reportErr)
		os.Exit(1)
	}
}

// parseDiffArgs parses positional + flag arguments.
func parseDiffArgs(args []string) (left, right string, walk diff.WalkOptions, print diff.PrintOptions) {
	fs := flag.NewFlagSet(constants.CmdDiff, flag.ExitOnError)
	jsonOut := fs.Bool(constants.FlagDiffJSON, false, "emit JSON instead of text")
	onlyConf := fs.Bool(constants.FlagDiffOnlyConflicts, false, "show only conflicting files")
	onlyMiss := fs.Bool(constants.FlagDiffOnlyMissing, false, "show only missing-side files")
	includeIdent := fs.Bool(constants.FlagDiffIncludeIdentical, false, "include identical files in output")
	includeVCS := fs.Bool(constants.FlagDiffIncludeVCS, false, "include .git/ in the walk")
	includeNM := fs.Bool(constants.FlagDiffIncludeNodeMods, false, "include node_modules/ in the walk")

	positional := reorderFlagsBeforeArgs(args)
	if err := fs.Parse(positional); err != nil {
		os.Exit(2)
	}
	left, right = extractDiffPositional(fs.Args())
	walk = diff.WalkOptions{IncludeVCS: *includeVCS, IncludeNodeModules: *includeNM}
	print = diff.PrintOptions{
		OnlyConflicts: *onlyConf, OnlyMissing: *onlyMiss,
		IncludeIdentical: *includeIdent, JSON: *jsonOut,
	}

	return left, right, walk, print
}

// extractDiffPositional pulls LEFT and RIGHT from the parsed args.
func extractDiffPositional(rest []string) (string, string) {
	if len(rest) != 2 {
		fmt.Fprintf(os.Stderr, constants.ErrDiffUsageFmt)
		os.Exit(2)
	}

	return rest[0], rest[1]
}

// guardDiffPaths blocks LEFT==RIGHT (nested folders are still useful for
// previewing accidental copies, so we only check exact equality).
func guardDiffPaths(left, right diff.Endpoint) error {
	lAbs, _ := filepath.Abs(left.WorkingDir)
	rAbs, _ := filepath.Abs(right.WorkingDir)
	if lAbs == rAbs {
		return fmt.Errorf(constants.ErrDiffSameFolder, lAbs)
	}

	return nil
}
