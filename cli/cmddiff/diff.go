package cmddiff

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/diff"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// runDiff implements `gitmap diff LEFT RIGHT`: a read-only preview
// of what `gitmap merge-*` would change between two folders,
// or single-target diff against HEAD in a git repository.
//
// Spec: companion to 02-spec/01-app/97-move-and-merge.md
func runDiff(args []string) error {
	checkHelp(constants.CmdDiff, args)
	left, right, walkOpts, printOpts, err := parseDiffArgs(args)
	if err != nil {
		return err
	}

	if left == "" && right == "" {
		return nil
	}

	leftEP, rightEP, err := resolveAndGuardEndpoints(left, right)
	if err != nil {
		return err
	}

	return executeDiffAndReport(leftEP, rightEP, walkOpts, printOpts)
}

func resolveEndpointSide(target, wrapToken string) (diff.Endpoint, error) {
	ep, err := diff.ResolveEndpoint(target)
	if err != nil {
		return diff.Endpoint{}, apperror.WrapSimple(err, wrapToken)
	}

	return ep, nil
}

func resolveAndGuardEndpoints(left, right string) (diff.Endpoint, diff.Endpoint, error) {
	leftEP, err := resolveEndpointSide(left, constants.WrapResolveEndpointLeft)
	if err != nil {
		return diff.Endpoint{}, diff.Endpoint{}, err
	}

	rightEP, err := resolveEndpointSide(right, constants.WrapResolveEndpointRight)
	if err != nil {
		return diff.Endpoint{}, diff.Endpoint{}, err
	}

	if guardErr := guardDiffPaths(leftEP, rightEP); guardErr != nil {
		return diff.Endpoint{}, diff.Endpoint{}, apperror.WrapSimple(guardErr, constants.WrapGuardPaths)
	}

	return leftEP, rightEP, nil
}

func executeDiffAndReport(leftEP, rightEP diff.Endpoint, walkOpts diff.WalkOptions, printOpts diff.PrintOptions) error {
	entries, err := diff.DiffTrees(leftEP.WorkingDir, rightEP.WorkingDir, walkOpts)
	if err != nil {
		return apperror.WrapSimple(err, constants.WrapDiffFailed)
	}

	if reportErr := diff.Report(os.Stdout, entries, printOpts); reportErr != nil {
		return apperror.WrapSimple(reportErr, constants.WrapReportFailed)
	}

	return nil
}

type diffFlagBindings struct {
	jsonOut      *bool
	onlyConf     *bool
	onlyMiss     *bool
	includeIdent *bool
	includeVCS   *bool
	includeNM    *bool
}

func initDiffFlags(fs *flag.FlagSet) diffFlagBindings {
	return diffFlagBindings{
		jsonOut:      fs.Bool(constants.FlagDiffJSON, false, constants.FlagDescJSON),
		onlyConf:     fs.Bool(constants.FlagDiffOnlyConflicts, false, constants.FlagDescOnlyConflicts),
		onlyMiss:     fs.Bool(constants.FlagDiffOnlyMissing, false, constants.FlagDescOnlyMissing),
		includeIdent: fs.Bool(constants.FlagDiffIncludeIdentical, false, constants.FlagDescIncludeIdentical),
		includeVCS:   fs.Bool(constants.FlagDiffIncludeVCS, false, constants.FlagDescIncludeVCS),
		includeNM:    fs.Bool(constants.FlagDiffIncludeNodeMods, false, constants.FlagDescIncludeNodeMods),
	}
}

// parseDiffArgs parses positional + flag arguments.
func parseDiffArgs(args []string) (string, string, diff.WalkOptions, diff.PrintOptions, error) {
	fs := flag.NewFlagSet(constants.CmdDiff, flag.ContinueOnError)
	flags := initDiffFlags(fs)
	positional := flagutil.ReorderFlagsBeforeArgs(args)
	if err := fs.Parse(positional); err != nil {
		return "", "", diff.WalkOptions{}, diff.PrintOptions{}, apperror.WrapSimple(err, "parse flags")
	}

	left, right, err := extractDiffPositional(fs.Args())
	if err != nil {
		return "", "", diff.WalkOptions{}, diff.PrintOptions{}, err
	}

	walk := diff.WalkOptions{IncludeVCS: *flags.includeVCS, IncludeNodeModules: *flags.includeNM}
	print := diff.PrintOptions{
		OnlyConflicts: *flags.onlyConf, OnlyMissing: *flags.onlyMiss,
		IncludeIdentical: *flags.includeIdent, JSON: *flags.jsonOut,
	}

	return left, right, walk, print, nil
}

// extractDiffPositional pulls LEFT and RIGHT from the parsed args.
func extractDiffPositional(rest []string) (string, string, error) {
	if len(rest) >= constants.ExpectedDiffPositionalArgs {
		return rest[constants.IndexDiffLeft], rest[constants.IndexDiffRight], nil
	}

	target := "."
	if len(rest) == 1 {
		target = rest[0]
	}

	return "", "", runSinglePathDiff(target)
}

// guardDiffPaths blocks LEFT==RIGHT (nested folders are still useful for
// previewing accidental copies, so we only check exact equality).
func guardDiffPaths(left, right diff.Endpoint) error {
	lAbs, _ := filepath.Abs(left.WorkingDir)
	rAbs, _ := filepath.Abs(right.WorkingDir)
	if lAbs == rAbs {
		return apperror.NewExecutionError(fmt.Sprintf(constants.ErrDiffSameFolder, lAbs))
	}

	return nil
}
