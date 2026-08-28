package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/diff"
)

// runDiff implements `gitmap diff LEFT RIGHT`: a read-only preview
// of what `gitmap merge-*` would change between two folders.
//
// Spec: companion to spec/01-app/97-move-and-merge.md
func runDiff(args []string) error {
	checkHelp("diff", args)
	left, right, walkOpts, printOpts := parseDiffArgs(args)

	leftEP, err := diff.ResolveEndpoint(left)
	if err != nil {
		return apperror.WrapSimple(err, "resolve-endpoint left")
	}
	rightEP, err := diff.ResolveEndpoint(right)
	if err != nil {
		return apperror.WrapSimple(err, "resolve-endpoint right")
	}
	if guardErr := guardDiffPaths(leftEP, rightEP); guardErr != nil {
		return apperror.WrapSimple(guardErr, "guard-paths")
	}

	entries, err := diff.DiffTrees(leftEP.WorkingDir, rightEP.WorkingDir, walkOpts)
	if err != nil {
		return apperror.NewSimple("diff failed: "+constants.LogPrefixDiff, "E9000")
	}
	if reportErr := diff.Report(os.Stdout, entries, printOpts); reportErr != nil {
		return apperror.NewSimple("report failed: "+constants.LogPrefixDiff, "E9000")
	}
	return nil
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
