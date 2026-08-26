// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/config"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// runRelease handles the 'release' command.
func runRelease(args []string) {
	if tryCrossDirRelease(args) {
		return
	}
	checkHelp("release", args)
	v, assets, commit, branch, bump, notes, targets, zg, zi, bundle, draft, dryRun, verbose, comp, cs, bin, listTgts, noComm, yes := parseReleaseFlags(args)
	_ = verbose
	if listTgts {
		printListTargets(targets)
		return
	}
	if handleOutsideRepoRelease(args, v, bump, commit, branch, yes) {
		return
	}
	performInsideRepoRelease(v, assets, commit, branch, bump, notes, targets, zg, zi, bundle, draft, dryRun, verbose, comp, cs, bin, noComm, yes)
}

func handleOutsideRepoRelease(args []string, version, bump, commit, branch string, yes bool) bool {
	if release.IsInsideGitRepo() {
		return false
	}
	if tryRunReleaseInRecentClone(args) || (shouldAutoBumpMinor(version, bump, commit, branch) && tryRunReleaseScanDir(yes)) {
		return true
	}
	runReleaseSelf(args)
	return true
}

func performInsideRepoRelease(version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes bool) {
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

func loadReleaseConfig() model.Config {
	cfg, cfgErr := config.LoadFromFile(constants.DefaultConfigPath)
	if cfgErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not load config: %v\n", cfgErr)
	}
	return cfg
}

func buildReleaseOptions(version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes bool, cfg model.Config) release.Options {
	return release.Options{
		Version: version, Assets: assets, Commit: commit, Branch: branch,
		Bump: bump, Notes: notes, Targets: targets, ConfigTargets: cfg.Release.Targets,
		ZipGroups: zipGroups, ZipItems: zipItems, BundleName: bundleName,
		IsDraft: draft, DryRun: dryRun, Verbose: verbose,
		Compress: compress || cfg.Release.IsCompress, Checksums: checksums || cfg.Release.HasChecksums,
		Bin: bin, NoCommit: noCommit, Yes: yes,
	}
}

// executeRelease builds options and runs the release workflow.
func executeRelease(version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes bool) {
	cfg := loadReleaseConfig()
	opts := buildReleaseOptions(version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes, cfg)
	if err := release.Execute(opts); err != nil {
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

type releaseFlagHolders struct {
	assets, commit, branch, bump, notes, targets, bundle                         *string
	draft, dryRun, verbose, compress, checksums, bin, listTargets, noCommit, yes *bool
	zgGroups                                                                     zipGroupFlag
	zgItems                                                                      zipItemFlag
}

func initStringFlags(fs *flag.FlagSet, h *releaseFlagHolders) {
	h.assets = fs.String("assets", "", constants.FlagDescAssets)
	h.commit = fs.String("commit", "", constants.FlagDescCommit)
	h.branch = fs.String("branch", "", constants.FlagDescRelBranch)
	h.bump = fs.String(constants.FlagBump, "", constants.FlagDescBump)
	h.notes = fs.String("notes", "", constants.FlagDescNotes)
	h.targets = fs.String("targets", "", constants.FlagDescTargets)
	h.bundle = fs.String("bundle", "", constants.FlagDescZGBundle)
	fs.StringVar(h.notes, "N", "", constants.FlagDescNotes)
}

func initBoolFlags(fs *flag.FlagSet, h *releaseFlagHolders) {
	h.draft = fs.Bool("draft", false, constants.FlagDescDraft)
	h.dryRun = fs.Bool("dry-run", false, constants.FlagDescDryRun)
	h.verbose = fs.Bool("verbose", false, constants.FlagDescVerbose)
	h.compress = fs.Bool("compress", false, constants.FlagDescCompress)
	h.checksums = fs.Bool("checksums", false, constants.FlagDescChecksums)
	h.bin = fs.Bool("bin", false, constants.FlagDescBin)
	h.listTargets = fs.Bool("list-targets", false, constants.FlagDescListTargets)
	h.noCommit = fs.Bool("no-commit", false, constants.FlagDescNoCommit)
	h.yes = fs.Bool("yes", false, constants.FlagDescYes)
	fs.BoolVar(h.bin, "b", false, constants.FlagDescBin)
	fs.BoolVar(h.yes, "y", false, constants.FlagDescYes)
}

func newReleaseFlagSet() (*flag.FlagSet, *releaseFlagHolders) {
	fs := flag.NewFlagSet(constants.CmdRelease, flag.ExitOnError)
	h := &releaseFlagHolders{}
	initStringFlags(fs, h)
	initBoolFlags(fs, h)
	fs.Var(&h.zgGroups, "zip-group", constants.FlagDescZGZipGroup)
	fs.Var(&h.zgItems, "Z", constants.FlagDescZGZipItem)
	return fs, h
}

// parseReleaseFlags parses flags for the release command.
func parseReleaseFlags(args []string) (version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, listTargets, noCommit, yes bool) {
	fs, h := newReleaseFlagSet()
	fs.Parse(reorderFlagsBeforeArgs(args))
	if fs.NArg() > 0 {
		version = normalizeVersion(fs.Arg(0))
	}
	if forceYesOverride {
		*h.yes = true
	}
	return version, *h.assets, *h.commit, *h.branch, *h.bump, *h.notes, *h.targets, []string(h.zgGroups), []string(h.zgItems), *h.bundle, *h.draft, *h.dryRun, *h.verbose, *h.compress, *h.checksums, *h.bin, *h.listTargets, *h.noCommit, *h.yes
}

const (
	versionPrefix    = "v"
	versionTrimChars = "vV"
)

// normalizeVersion strips all leading 'v' or 'V' characters and prepends exactly one 'v'.
func normalizeVersion(v string) string {
	if len(v) == 0 {
		return v
	}
	return versionPrefix + strings.TrimLeft(v, versionTrimChars)
}

// forceYesOverride, when true, forces parseReleaseFlags to return
// yes=true regardless of CLI args. Set by `gitmap pr` / `gitmap prc`
// before delegating to runRelease so the post-release auto-commit
// prompt cannot stall the pipeline.
var forceYesOverride bool

// printListTargets resolves and prints the target matrix, then returns.
func printListTargets(flagTargets string) {
	cfg := loadReleaseConfig()
	targets, err := release.ResolveTargets(flagTargets, cfg.Release.Targets)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}
	printTargetDetails(flagTargets, cfg.Release.Targets, targets)
}

func printTargetDetails(flagTargets string, configTargets []model.ReleaseTarget, targets []release.BuildTarget) {
	source := resolveTargetSource(flagTargets, configTargets)
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
