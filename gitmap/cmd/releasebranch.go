// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// runReleaseBranch handles the 'release-branch' command.
func runReleaseBranch(args []string) error {
	checkHelp("release-branch", args)
	branch, assets, notes, draft, dryRun, verbose, noCommit, yes := parseReleaseBranchFlags(args)
	_ = verbose

	if len(branch) == 0 {
		return apperror.NewSimple(constants.ErrReleaseBranchUsage, "E9000")
	}

	err := release.ExecuteFromBranch(branch, assets, notes, draft, dryRun, noCommit, yes)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrBareFmt)
	}
	return nil
}

// parseReleaseBranchFlags parses flags for the release-branch command.
func parseReleaseBranchFlags(args []string) (branch, assets, notes string, draft, dryRun, verbose, noCommit, yes bool) {
	fs := flag.NewFlagSet(constants.CmdReleaseBranch, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	notesFlag := fs.String("notes", "", constants.FlagDescNotes)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	noCommitFlag := fs.Bool("no-commit", false, constants.FlagDescNoCommit)
	yesFlag := fs.Bool("yes", false, constants.FlagDescYes)

	// Register -N as shorthand for --notes, -y as shorthand for --yes.
	fs.StringVar(notesFlag, "N", "", constants.FlagDescNotes)
	fs.BoolVar(yesFlag, "y", false, constants.FlagDescYes)

	fs.Parse(args)

	branch = ""
	if fs.NArg() > 0 {
		branch = fs.Arg(0)
	}

	return branch, *assetsFlag, *notesFlag, *draftFlag, *dryRunFlag, *verboseFlag, *noCommitFlag, *yesFlag
}
