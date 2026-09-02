package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ScanProbeOptions bundles the flags that govern the optional
// background version-probe pass scan kicks off after upserting repos.
// Bundling them keeps parseScanFlags's return list manageable and
// makes the runner-wiring call site read as a single cohesive object.
type ScanProbeOptions struct {
	// Disable suppresses the background probe entirely. Set via --no-probe.
	Disable bool
	// NoWait makes scan return immediately after dispatching jobs;
	// the runner keeps draining in the background until process exit.
	NoWait bool
	// Concurrency overrides the worker count. 0 = use the documented
	// default; negative values disable the runner the same as --no-probe.
	Concurrency int
	// ConcurrencySet records whether the user explicitly passed
	// --probe-workers (or the deprecated --probe-concurrency alias).
	// Used to bypass the auto-trigger ceiling for power users who
	// clearly opted in.
	ConcurrencySet bool
	// Depth is the `--depth N` value forwarded to the shallow-clone
	// fallback inside the background runner. Defaults to
	// constants.ProbeDefaultDepth (1) when no flag was passed.
	Depth int
}

type scanFlagPointers struct {
	cfgFlag           *string
	modeFlag          *string
	outputFlag        *string
	outFileFlag       *string
	outputPathFlag    *string
	manifestFlag      *string
	relRootFlag       *string
	defaultBranchFlag *string
	ghDesktopFlag     *bool
	openFlag          *bool
	quietFlag         *bool
	noVSCodeSyncFlag  *bool
	noAutoTagsFlag    *bool
	reportErrFlag     *bool
	compactFlag       *bool
	fixFlag           *bool
	workersFlag       *int
	concurrencyFlag   *int
	maxDepthFlag      *int
	noProbeFlag       *bool
	noProbeWaitFlag   *bool
	probeConcFlag     *int
	probeWorkersFlag  *int
	probeDepthFlag    *int
}

func registerScanStringFlags(fs *flag.FlagSet, flagPtrs *scanFlagPointers) {
	flagPtrs.cfgFlag = fs.String("config", constants.DefaultConfigPath, constants.FlagDescConfig)
	flagPtrs.modeFlag = fs.String("mode", "", constants.FlagDescMode)
	flagPtrs.outputFlag = fs.String("output", "", constants.FlagDescOutput)
	flagPtrs.outFileFlag = fs.String("out-file", "", constants.FlagDescOutFile)
	flagPtrs.outputPathFlag = fs.String("output-path", "", constants.FlagDescOutputPath)
	flagPtrs.manifestFlag = fs.String(constants.FlagScanManifest, "", constants.FlagDescScanManifest)
	flagPtrs.relRootFlag = fs.String(constants.FlagScanRelativeRoot, "", constants.FlagDescScanRelativeRoot)
	flagPtrs.defaultBranchFlag = fs.String(constants.FlagScanDefaultBranch, "", constants.FlagDescScanDefaultBranch)
}

func registerScanToggles(fs *flag.FlagSet, flagPtrs *scanFlagPointers) {
	flagPtrs.ghDesktopFlag = fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	flagPtrs.openFlag = fs.Bool("open", false, constants.FlagDescOpen)
	flagPtrs.quietFlag = fs.Bool("quiet", false, constants.FlagDescQuiet)
	flagPtrs.noVSCodeSyncFlag = fs.Bool(constants.FlagNoVSCodeSync, false, constants.FlagDescNoVSCodeSync)
	flagPtrs.noAutoTagsFlag = fs.Bool(constants.FlagNoAutoTags, false, constants.FlagDescNoAutoTags)
	flagPtrs.reportErrFlag = fs.Bool(constants.FlagScanReportErrors, false, constants.FlagDescScanReportErrors)
	flagPtrs.compactFlag = fs.Bool(constants.FlagScanCompact, false, constants.FlagDescScanCompact)
	flagPtrs.fixFlag = fs.Bool("fix", false, "Reconcile missing/stale repositories from the gitmap tracking database")
}

func registerScanIntFlags(fs *flag.FlagSet, flagPtrs *scanFlagPointers) {
	flagPtrs.workersFlag = fs.Int(constants.FlagScanWorkers, constants.DefaultScanWorkers, constants.FlagDescScanWorkers)
	flagPtrs.concurrencyFlag = fs.Int(constants.FlagScanWorkersConcurrencyAlias,
		constants.DefaultScanWorkers, constants.FlagDescScanWorkersConcurrencyAlias)
	flagPtrs.maxDepthFlag = fs.Int(constants.FlagScanMaxDepth, constants.DefaultScanMaxDepth, constants.FlagDescScanMaxDepth)
}

func registerScanProbeFlags(fs *flag.FlagSet, flagPtrs *scanFlagPointers) {
	flagPtrs.noProbeFlag = fs.Bool(constants.ScanProbeFlagDisable, false, constants.FlagDescScanProbeDisable)
	flagPtrs.noProbeWaitFlag = fs.Bool(constants.ScanProbeFlagNoWait, false, constants.FlagDescScanProbeNoWait)
	flagPtrs.probeConcFlag = fs.Int(constants.ScanProbeFlagConcurrency,
		constants.ScanProbeDefaultConcurrency, constants.FlagDescScanProbeConcurrency)
	flagPtrs.probeWorkersFlag = fs.Int(constants.ScanProbeFlagProbeWorkers,
		constants.ScanProbeDefaultConcurrency, constants.FlagDescScanProbeProbeWorkers)
	flagPtrs.probeDepthFlag = fs.Int(constants.ScanProbeFlagProbeDepth,
		constants.ProbeDefaultDepth, constants.FlagDescScanProbeProbeDepth)
}

func newScanFlagSet(fs *flag.FlagSet) *scanFlagPointers {
	flagPtrs := &scanFlagPointers{}
	registerScanStringFlags(fs, flagPtrs)
	registerScanToggles(fs, flagPtrs)
	registerScanIntFlags(fs, flagPtrs)
	registerScanProbeFlags(fs, flagPtrs)

	return flagPtrs
}

func resolveScanOutputPath(outputPath, manifest string) string {
	if outputPath == "" && manifest != "" {
		return manifest
	}

	return outputPath
}

// parseScanFlags parses flags for the scan command.
func parseScanFlags(args []string) (dir, configPath, mode, output, outFile, outputPath, relativeRoot, defaultBranch string, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, compact, fix bool, workers, maxDepth int, probeOpts ScanProbeOptions) {
	fs := flag.NewFlagSet(constants.CmdScan, flag.ExitOnError)
	scanFlags := newScanFlagSet(fs)
	fs.Parse(args)

	dir = resolveScanDir(fs)
	probeOpts = resolveScanProbeOptions(fs, scanFlags.noProbeFlag, scanFlags.noProbeWaitFlag,
		scanFlags.probeConcFlag, scanFlags.probeWorkersFlag, scanFlags.probeDepthFlag)
	resolvedWorkers := resolveScanWorkers(fs, scanFlags.workersFlag, scanFlags.concurrencyFlag)
	resolvedOutputPath := resolveScanOutputPath(*scanFlags.outputPathFlag, *scanFlags.manifestFlag)

	return dir, *scanFlags.cfgFlag, *scanFlags.modeFlag, *scanFlags.outputFlag, *scanFlags.outFileFlag, resolvedOutputPath, *scanFlags.relRootFlag, *scanFlags.defaultBranchFlag, *scanFlags.ghDesktopFlag, *scanFlags.openFlag, *scanFlags.quietFlag, *scanFlags.noVSCodeSyncFlag, *scanFlags.noAutoTagsFlag, *scanFlags.reportErrFlag, *scanFlags.compactFlag, *scanFlags.fixFlag, resolvedWorkers, *scanFlags.maxDepthFlag, probeOpts
}

// resolveScanWorkers reconciles --workers (canonical) against the
// deprecated --concurrency alias. Canonical wins when both are set;
// when only --concurrency is set we honor it and emit a one-line
// stderr deprecation notice. Mirrors resolveScanProbeOptions for
// the --probe-workers / --probe-concurrency pair.
func resolveScanWorkers(fs *flag.FlagSet, workers, concurrency *int) int {
	isWorkersSet := wasFlagPassed(fs, constants.FlagScanWorkers)
	isConcSet := wasFlagPassed(fs, constants.FlagScanWorkersConcurrencyAlias)
	if !isWorkersSet && isConcSet {
		fmt.Fprint(os.Stderr, constants.MsgScanWorkersConcurrencyAlias)

		return *concurrency
	}

	return *workers
}

func resolveProbeConcurrency(fs *flag.FlagSet, probeConc, probeWorkers *int) (int, bool) {
	isConcSet := wasFlagPassed(fs, constants.ScanProbeFlagConcurrency)
	isWorkersSet := wasFlagPassed(fs, constants.ScanProbeFlagProbeWorkers)
	if !isWorkersSet && isConcSet {
		fmt.Fprint(os.Stderr, constants.MsgScanProbeConcurrencyAlias)

		return *probeConc, true
	}

	return *probeWorkers, isWorkersSet
}

// resolveScanProbeOptions reconciles the deprecated --probe-concurrency
// against the unified --probe-workers. The new flag wins when both are
// set; when only the deprecated one is set we honor it and emit a
// one-line stderr deprecation notice. Depth comes through unchanged.
func resolveScanProbeOptions(fs *flag.FlagSet, noProbe, noWait *bool,
	probeConc, probeWorkers, probeDepth *int) ScanProbeOptions {
	conc, isConcSet := resolveProbeConcurrency(fs, probeConc, probeWorkers)

	return ScanProbeOptions{
		Disable:        *noProbe,
		NoWait:         *noWait,
		Concurrency:    conc,
		ConcurrencySet: isConcSet,
		Depth:          *probeDepth,
	}
}

// wasFlagPassed reports whether the named flag was explicitly set on
// the command line (vs left at its default). Go's stdlib flag package
// doesn't surface this directly, so we walk Visit to find out.
func wasFlagPassed(fs *flag.FlagSet, flagName string) bool {
	hasSeen := false
	fs.Visit(func(flagItem *flag.Flag) {
		if flagItem.Name == flagName {
			hasSeen = true
		}
	})

	return hasSeen
}

// resolveScanDir returns the scan directory from positional args or default.
func resolveScanDir(fs *flag.FlagSet) string {
	if fs.NArg() > 0 {
		return fs.Arg(0)
	}

	return constants.DefaultDir
}

// CloneFlags holds all parsed clone-command flags and positional args.
// Exposing the full positional slice (Positional) lets runClone detect
// the multi-URL invocation form documented in spec/01-app/104-clone-multi.md.
type CloneFlags struct {
	Source     string
	FolderName string
	TargetDir  string
	SSHKeyName string
	// DefaultBranch mirrors `gitmap scan --default-branch`: when a
	// manifest row has an unknown / empty Branch (or a non-trustworthy
	// BranchSource like "detached" or "unknown"), the cloner rebuilds
	// the clone instruction as `git clone -b <DefaultBranch> ...`
	// instead of letting the remote's default HEAD decide. Empty keeps
	// the legacy behavior. Same constant powers both flags so the help
	// wording stays byte-identical across surfaces.
	DefaultBranch  string
	Positional     []string
	SafePull       bool
	GHDesktop      bool
	NoReplace      bool
	Verbose        bool
	Audit          bool
	MaxConcurrency int
	// Output selects the per-repo summary format. Empty (default)
	// keeps the legacy terse messages; "terminal" emits the
	// standardized RepoTermBlock right before each clone runs so
	// the shape matches scan/clone-next/clone-from/probe.
	Output string
	// VerifyCmdFaithful enables the dry-run argv-vs-displayed
	// checker. See clonetermverify.go for behavior.
	VerifyCmdFaithful bool
	// VerifyCmdFaithfulExitOnMismatch upgrades the verifier into a
	// hard failure: any divergence sets a sticky bit and the run tail
	// exits with constants.CloneVerifyCmdFaithfulExitCode. Implies
	// VerifyCmdFaithful.
	VerifyCmdFaithfulExitOnMismatch bool
	// PrintCloneArgv dumps the executor's literal argv tokens to
	// stderr. See cloneprintargv.go for behavior.
	PrintCloneArgv bool
	// NoVSCodeSync suppresses the post-clone update of the
	// alefragnani.project-manager projects.json file. Mirrors the
	// flag of the same name on `gitmap scan`. Default false →
	// every successful clone is reflected in the VS Code Project
	// Manager sidebar without an extra command. See
	// spec/01-vscode-project-manager-sync/02-clone-sync.md.
	NoVSCodeSync bool
	// UseSSH forces every direct URL (and the first positional in
	// multi-URL form) to be rewritten into its `git@host:owner/repo.git`
	// SSH-shorthand form before git is invoked. HTTPS and `ssh://` URLs
	// are converted via ConvertURLToSSH; already-SSH-shorthand URLs are
	// normalized (`.git` suffix appended). See `--ssh` in clone.md.
	UseSSH bool
	// UseHTTPS is the symmetric counterpart of UseSSH — forces every
	// URL into `https://host/owner/repo.git` form. Useful in CI/headless
	// environments where the SSH agent isn't unlocked.
	UseHTTPS bool
	// DryRun short-circuits every git clone in this run: the runner
	// prints the exact command + target path but never invokes git.
	// Plumbed through cfr/cfrp as well via parseCloneFixRepoArgs.
	DryRun bool
	// AssumeYes skips the SSH first-connect host-key prompt by asking
	// OpenSSH to accept new host keys. Changed host keys still fail.
	IsAssumeYes bool
	Clean       bool
	MissingOnly bool
	Fix         bool
}

type cloneFlagPointers struct {
	targetFlag        *string
	safePullFlag      *bool
	ghDesktopFlag     *bool
	verboseFlag       *bool
	noReplaceFlag     *bool
	cleanFlag         *bool
	missingOnlyFlag   *bool
	auditFlag         *bool
	maxConcFlag       *int
	sshKeyFlag        *string
	defaultBranchFlag *string
	outputFlag        *string
	verifyFlag        *bool
	verifyExitFlag    *bool
	printArgvFlag     *bool
	noVSCodeSyncFlag  *bool
	debugPathsFlag    *bool
	sshFlag           *bool
	httpsFlag         *bool
	dryRunFlag        *bool
	yesFlag           *bool
	fixFlag           *bool
}

func registerCloneStringFlags(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) {
	flagPtrs.targetFlag = fs.String("target-dir", constants.DefaultDir, constants.FlagDescTargetDir)
	flagPtrs.sshKeyFlag = fs.String("ssh-key", "", "SSH key name for clone")
	fs.StringVar(flagPtrs.sshKeyFlag, "K", "", "SSH key name (short)")
	flagPtrs.defaultBranchFlag = fs.String(constants.FlagScanDefaultBranch, "", constants.FlagDescScanDefaultBranch)
	flagPtrs.outputFlag = fs.String(constants.FlagCloneTermOutput, "", constants.FlagDescCloneTermOutput)
	flagPtrs.maxConcFlag = fs.Int(constants.CloneFlagMaxConcurrency,
		constants.CloneDefaultMaxConcurrency, constants.FlagDescCloneMaxConcurrency)
}

func registerCloneToggles(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) {
	flagPtrs.safePullFlag = fs.Bool("safe-pull", false, constants.FlagDescSafePull)
	flagPtrs.ghDesktopFlag = fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	flagPtrs.verboseFlag = fs.Bool("verbose", false, constants.FlagDescVerbose)
	flagPtrs.noReplaceFlag = fs.Bool("no-replace", false, constants.FlagDescCloneNoReplace)
	flagPtrs.cleanFlag = fs.Bool("clean", false, "Forcefully delete the local folder and re-clone")
	flagPtrs.missingOnlyFlag = fs.Bool("missing-only", false, "Skip existing directories entirely")
	flagPtrs.auditFlag = fs.Bool(constants.CloneFlagAudit, false, constants.FlagDescCloneAudit)
	flagPtrs.noVSCodeSyncFlag = fs.Bool(constants.FlagNoVSCodeSync, false, constants.FlagDescNoVSCodeSync)
	flagPtrs.debugPathsFlag = fs.Bool(constants.FlagDebugPaths, false, constants.FlagDescDebugPaths)
	flagPtrs.fixFlag = fs.Bool("fix", false, "Remove repeated projects across tools")
	fs.BoolVar(flagPtrs.fixFlag, "repeat-fix", false, "Alias for --fix")
}

func registerCloneExecutionFlags(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) {
	flagPtrs.verifyFlag = fs.Bool(constants.FlagCloneVerifyCmdFaithful, false,
		constants.FlagDescCloneVerifyCmdFaithful)
	flagPtrs.verifyExitFlag = fs.Bool(constants.FlagCloneVerifyCmdFaithfulExitOnMismatch,
		false, constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch)
	flagPtrs.printArgvFlag = fs.Bool(constants.FlagClonePrintArgv, false,
		constants.FlagDescClonePrintArgv)
	flagPtrs.sshFlag = fs.Bool("ssh", false,
		"Force every clone URL into `git@host:owner/repo.git` SSH-shorthand form before git runs (auto-converts HTTPS / `ssh://` URLs)")
	fs.BoolVar(flagPtrs.sshFlag, "sh", false, "Short alias for --ssh")
	flagPtrs.httpsFlag = fs.Bool("https", false,
		"Force every clone URL into `https://host/owner/repo.git` form (auto-converts SSH-shorthand / `ssh://` URLs)")
	fs.BoolVar(flagPtrs.httpsFlag, "ht", false, "Short alias for --https")
	flagPtrs.dryRunFlag = fs.Bool(constants.FlagCloneDryRun, false, constants.FlagDescCloneDryRun)
	fs.BoolVar(flagPtrs.dryRunFlag, constants.FlagCloneDryRunShort, false, "Short alias for --dry-run")
	flagPtrs.yesFlag = fs.Bool(constants.FlagCloneYes, false, constants.FlagDescCloneYes)
	fs.BoolVar(flagPtrs.yesFlag, constants.FlagCloneYesShort, false, constants.FlagDescCloneYes)
}

func newCloneFlagSet(fs *flag.FlagSet) *cloneFlagPointers {
	flagPtrs := &cloneFlagPointers{}
	registerCloneStringFlags(fs, flagPtrs)
	registerCloneToggles(fs, flagPtrs)
	registerCloneExecutionFlags(fs, flagPtrs)

	return flagPtrs
}

func populateCloneToggles(cloneOpts *CloneFlags, flagPtrs *cloneFlagPointers) {
	cloneOpts.SafePull = *flagPtrs.safePullFlag
	cloneOpts.GHDesktop = *flagPtrs.ghDesktopFlag
	cloneOpts.NoReplace = *flagPtrs.noReplaceFlag
	cloneOpts.Verbose = *flagPtrs.verboseFlag
	cloneOpts.Audit = *flagPtrs.auditFlag
	cloneOpts.MaxConcurrency = *flagPtrs.maxConcFlag
	cloneOpts.NoVSCodeSync = *flagPtrs.noVSCodeSyncFlag
	cloneOpts.Clean = *flagPtrs.cleanFlag
	cloneOpts.MissingOnly = *flagPtrs.missingOnlyFlag
	cloneOpts.Fix = *flagPtrs.fixFlag
}

func populateCloneExecutionFlags(cloneOpts *CloneFlags, flagPtrs *cloneFlagPointers) {
	cloneOpts.VerifyCmdFaithful = *flagPtrs.verifyFlag
	cloneOpts.VerifyCmdFaithfulExitOnMismatch = *flagPtrs.verifyExitFlag
	cloneOpts.PrintCloneArgv = *flagPtrs.printArgvFlag
	cloneOpts.UseSSH = *flagPtrs.sshFlag
	cloneOpts.UseHTTPS = *flagPtrs.httpsFlag
	cloneOpts.DryRun = *flagPtrs.dryRunFlag
	cloneOpts.IsAssumeYes = *flagPtrs.yesFlag
}

func buildCloneFlags(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) CloneFlags {
	cloneOpts := CloneFlags{
		Source:        resolveCloneSource(fs),
		FolderName:    resolveCloneFolderName(fs),
		TargetDir:     *flagPtrs.targetFlag,
		SSHKeyName:    *flagPtrs.sshKeyFlag,
		DefaultBranch: *flagPtrs.defaultBranchFlag,
		Positional:    fs.Args(),
		Output:        *flagPtrs.outputFlag,
	}
	populateCloneToggles(&cloneOpts, flagPtrs)
	populateCloneExecutionFlags(&cloneOpts, flagPtrs)

	return cloneOpts
}

// parseCloneFlags parses flags for the clone command.
func parseCloneFlags(args []string) CloneFlags {
	fs := flag.NewFlagSet(constants.CmdClone, flag.ExitOnError)
	flagPtrs := newCloneFlagSet(fs)
	fs.Parse(reorderFlagsBeforeArgs(args))
	applyDebugPathsEnv(*flagPtrs.debugPathsFlag)

	return buildCloneFlags(fs, flagPtrs)
}

// resolveCloneSource returns the clone source from positional args.
func resolveCloneSource(fs *flag.FlagSet) string {
	if fs.NArg() > 0 {
		return fs.Arg(0)
	}

	return ""
}

// resolveCloneFolderName returns the optional folder name (second positional arg).
// When the second positional looks like a URL, it's NOT a folder name —
// callers must treat the full positional list as a multi-URL batch instead.
func resolveCloneFolderName(fs *flag.FlagSet) string {
	if fs.NArg() <= 1 {
		return ""
	}
	secondArg := fs.Arg(1)
	if isLikelyURL(secondArg) {
		return ""
	}

	return secondArg
}

// isLikelyURL is a cheap prefix check used to disambiguate
// "folder name" vs "second URL" without importing the clone package.
// Mirrors isDirectURL in clone.go — keep both in sync.
func isLikelyURL(rawURL string) bool {
	lowerURL := strings.ToLower(strings.TrimSpace(rawURL))

	return strings.HasPrefix(lowerURL, "https://") ||
		strings.HasPrefix(lowerURL, "http://") ||
		strings.HasPrefix(lowerURL, "ssh://") ||
		strings.HasPrefix(lowerURL, "git@")
}
