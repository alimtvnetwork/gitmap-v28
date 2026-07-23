// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/config"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// runRelease handles the 'release' command.
//
// Form 1 — `gitmap r vX.Y.Z`         : in-place release of the current repo.
// Form 2 — `gitmap r <repo> vX.Y.Z`  : cross-dir — chdir into <repo>, fetch
//   - pull --rebase, release, chdir back.
//
// See `releaserebase.go` for the cross-dir machinery.
func runRelease(args []string) {
	if tryCrossDirRelease(args) {
		return
	}
	checkHelp("release", args)

	version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, bin, listTargets, noCommit, yes := parseReleaseFlags(args)
	_ = verbose

	if listTargets {
		printListTargets(targets)

		return
	}

	// Auto-fallback when not inside a Git repo.
	if !release.IsInsideGitRepo() {
		if tryRunReleaseInRecentClone(args) {
			return
		}
		if shouldAutoBumpMinor(version, bump, commit, branch) && tryRunReleaseScanDir(yes) {
			return
		}
		runReleaseSelf(args)

		return
	}

	requireOnline()
	bump = applyBareReleaseAutoBump(version, bump, commit, branch, yes)
	validateReleaseFlags(version, bump, commit, branch)
	executeRelease(version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes)
}

// applyBareReleaseAutoBump injects bump=minor when no explicit version/bump
// was provided, after confirming with the user (skipped with -y).
func applyBareReleaseAutoBump(version, bump, commit, branch string, yes bool) string {
	if !shouldAutoBumpMinor(version, bump, commit, branch) {
		return bump
	}

	current, next, ok := peekNextMinorVersion()
	if !ok {
		return bump
	}
	if !confirmAutoBump(current, next, yes) {
		fmt.Fprint(os.Stderr, constants.MsgReleaseAutoBumpAborted)
		os.Exit(1)
	}

	return constants.BumpMinor
}

// executeRelease builds options and runs the release workflow.
func executeRelease(version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes bool) {
	cfg, cfgErr := config.LoadFromFile(constants.DefaultConfigPath)
	if cfgErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not load config: %v\n", cfgErr)
	}

	opts := release.Options{
		Version: version, Assets: assets,
		Commit: commit, Branch: branch,
		Bump: bump, Notes: notes, Targets: targets,
		ConfigTargets: cfg.Release.Targets,
		ZipGroups:     zipGroups,
		ZipItems:      zipItems,
		BundleName:    bundleName,
		IsDraft:       draft, DryRun: dryRun,
		Verbose:   verbose,
		Compress:  compress || cfg.Release.Compress,
		Checksums: checksums || cfg.Release.Checksums,
		Bin:       bin,
		NoCommit:  noCommit,
		Yes:       yes,
	}
	err := release.Execute(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	persistReleaseToDB()
}

// validateReleaseFlags checks for mutually exclusive flags.
func validateReleaseFlags(version, bump, commit, branch string) {
	if len(bump) > 0 && len(version) > 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseBumpConflict)
		os.Exit(1)
	}
	if len(commit) > 0 && len(branch) > 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseCommitBranch)
		os.Exit(1)
	}
}

// zipGroupFlag collects multiple --zip-group values.
type zipGroupFlag []string

func (z *zipGroupFlag) String() string { return fmt.Sprintf("%v", *z) }
func (z *zipGroupFlag) Set(val string) error {
	*z = append(*z, val)

	return nil
}

// zipItemFlag collects multiple -Z values.
type zipItemFlag []string

func (z *zipItemFlag) String() string { return fmt.Sprintf("%v", *z) }
func (z *zipItemFlag) Set(val string) error {
	*z = append(*z, val)

	return nil
}

// parseReleaseFlags parses flags for the release command.
func parseReleaseFlags(args []string) (version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, listTargets, noCommit, yes bool) {
	fs := flag.NewFlagSet(constants.CmdRelease, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	commitFlag := fs.String("commit", "", constants.FlagDescCommit)
	branchFlag := fs.String("branch", "", constants.FlagDescRelBranch)
	bumpFlag := fs.String(constants.FlagBump, "", constants.FlagDescBump)
	notesFlag := fs.String("notes", "", constants.FlagDescNotes)
	targetsFlag := fs.String("targets", "", constants.FlagDescTargets)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	compressFlag := fs.Bool("compress", false, constants.FlagDescCompress)
	checksumsFlag := fs.Bool("checksums", false, constants.FlagDescChecksums)
	binFlag := fs.Bool("bin", false, constants.FlagDescBin)
	listTargetsFlag := fs.Bool("list-targets", false, constants.FlagDescListTargets)
	bundleFlag := fs.String("bundle", "", constants.FlagDescZGBundle)
	noCommitFlag := fs.Bool("no-commit", false, constants.FlagDescNoCommit)
	yesFlag := fs.Bool("yes", false, constants.FlagDescYes)

	fs.BoolVar(binFlag, "b", false, constants.FlagDescBin)
	fs.BoolVar(yesFlag, "y", false, constants.FlagDescYes)

	var zgGroups zipGroupFlag
	var zgItems zipItemFlag
	fs.Var(&zgGroups, "zip-group", constants.FlagDescZGZipGroup)
	fs.Var(&zgItems, "Z", constants.FlagDescZGZipItem)
	fs.StringVar(notesFlag, "N", "", constants.FlagDescNotes)

	fs.Parse(reorderFlagsBeforeArgs(args))
	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}
	// Process-level override set by pull-release / prc / similar
	// wrappers so the auto-commit prompt is never asked even when
	// the user did not type -y. Belt-and-suspenders with the
	// ensureYesForward arg injection.
	if forceYesOverride {
		*yesFlag = true
	}
	return version, *assetsFlag, *commitFlag, *branchFlag, *bumpFlag, *notesFlag, *targetsFlag, []string(zgGroups), []string(zgItems), *bundleFlag, *draftFlag, *dryRunFlag, *verboseFlag, *compressFlag, *checksumsFlag, *binFlag, *listTargetsFlag, *noCommitFlag, *yesFlag
}

// forceYesOverride, when true, forces parseReleaseFlags to return
// yes=true regardless of CLI args. Set by `gitmap pr` / `gitmap prc`
// before delegating to runRelease so the post-release auto-commit
// prompt cannot stall the pipeline.
var forceYesOverride bool

// printListTargets resolves and prints the target matrix, then returns.
func printListTargets(flagTargets string) {
	cfg, cfgErr := config.LoadFromFile(constants.DefaultConfigPath)
	if cfgErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not load config: %v\n", cfgErr)
	}

	targets, err := release.ResolveTargets(flagTargets, cfg.Release.Targets)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	source := resolveTargetSource(flagTargets, cfg.Release.Targets)

	fmt.Printf(constants.MsgListTargetsHeader, len(targets))
	fmt.Printf(constants.MsgListTargetsSource, source)

	for _, t := range targets {
		fmt.Printf(constants.MsgListTargetsRow, t.GOOS, t.GOARCH)
	}
}

// resolveTargetSource returns a human-readable label for the active target source.
func resolveTargetSource(flagTargets string, configTargets []model.ReleaseTarget) string {
	if len(flagTargets) > 0 {
		return "--targets flag"
	}

	if len(configTargets) > 0 {
		return "config.json (release.targets)"
	}

	return "built-in defaults"
}
