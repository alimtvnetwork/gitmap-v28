package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
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

// parseScanFlags parses flags for the scan command.
func parseScanFlags(args []string) (dir, configPath, mode, output, outFile, outputPath, relativeRoot, defaultBranch string, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors bool, workers, maxDepth int, probeOpts ScanProbeOptions) {
	fs := flag.NewFlagSet(constants.CmdScan, flag.ExitOnError)
	cfgFlag := fs.String("config", constants.DefaultConfigPath, constants.FlagDescConfig)
	modeFlag := fs.String("mode", "", constants.FlagDescMode)
	outputFlag := fs.String("output", "", constants.FlagDescOutput)
	outFileFlag := fs.String("out-file", "", constants.FlagDescOutFile)
	outputPathFlag := fs.String("output-path", "", constants.FlagDescOutputPath)
	// `--manifest` is a unified alias for `--output-path`, mirroring
	// the same flag on `gitmap reclone` so the scan→reclone round-
	// trip uses one vocabulary. When BOTH are passed, `--manifest`
	// wins ONLY if `--output-path` was left at its empty default;
	// otherwise the explicit `--output-path` is preserved (no silent
	// override of a previously-set value).
	manifestFlag := fs.String(constants.FlagScanManifest, "",
		constants.FlagDescScanManifest)
	relRootFlag := fs.String(constants.FlagScanRelativeRoot, "", constants.FlagDescScanRelativeRoot)
	// Empty default → mapper.resolveDefaultBranch falls back to
	// constants.DefaultBranch. We DON'T put "main" here because doing
	// so would make wasFlagPassed-style introspection impossible:
	// the user passing `--default-branch main` would look identical
	// to omitting the flag entirely.
	defaultBranchFlag := fs.String(constants.FlagScanDefaultBranch, "", constants.FlagDescScanDefaultBranch)
	ghDesktopFlag, openFlag, quietFlag := registerScanBoolFlags(fs)
	noVSCodeSyncFlag := fs.Bool(constants.FlagNoVSCodeSync, false, constants.FlagDescNoVSCodeSync)
	noAutoTagsFlag := fs.Bool(constants.FlagNoAutoTags, false, constants.FlagDescNoAutoTags)
	workersFlag := fs.Int(constants.FlagScanWorkers, constants.DefaultScanWorkers, constants.FlagDescScanWorkers)
	concurrencyFlag := fs.Int(constants.FlagScanWorkersConcurrencyAlias,
		constants.DefaultScanWorkers, constants.FlagDescScanWorkersConcurrencyAlias)
	maxDepthFlag := fs.Int(constants.FlagScanMaxDepth, constants.DefaultScanMaxDepth, constants.FlagDescScanMaxDepth)
	reportErrFlag := fs.Bool(constants.FlagScanReportErrors, false, constants.FlagDescScanReportErrors)
	noProbeFlag := fs.Bool(constants.ScanProbeFlagDisable, false, constants.FlagDescScanProbeDisable)
	noProbeWaitFlag := fs.Bool(constants.ScanProbeFlagNoWait, false, constants.FlagDescScanProbeNoWait)
	probeConcFlag := fs.Int(constants.ScanProbeFlagConcurrency,
		constants.ScanProbeDefaultConcurrency, constants.FlagDescScanProbeConcurrency)
	probeWorkersFlag := fs.Int(constants.ScanProbeFlagProbeWorkers,
		constants.ScanProbeDefaultConcurrency, constants.FlagDescScanProbeProbeWorkers)
	probeDepthFlag := fs.Int(constants.ScanProbeFlagProbeDepth,
		constants.ProbeDefaultDepth, constants.FlagDescScanProbeProbeDepth)
	fs.Parse(args)

	dir = resolveScanDir(fs)
	probeOpts = resolveScanProbeOptions(fs, noProbeFlag, noProbeWaitFlag,
		probeConcFlag, probeWorkersFlag, probeDepthFlag)
	resolvedWorkers := resolveScanWorkers(fs, workersFlag, concurrencyFlag)
	resolvedOutputPath := *outputPathFlag
	if resolvedOutputPath == "" && *manifestFlag != "" {
		resolvedOutputPath = *manifestFlag
	}

	return dir, *cfgFlag, *modeFlag, *outputFlag, *outFileFlag, resolvedOutputPath, *relRootFlag, *defaultBranchFlag, *ghDesktopFlag, *openFlag, *quietFlag, *noVSCodeSyncFlag, *noAutoTagsFlag, *reportErrFlag, resolvedWorkers, *maxDepthFlag, probeOpts
}

// resolveScanWorkers reconciles --workers (canonical) against the
// deprecated --concurrency alias. Canonical wins when both are set;
// when only --concurrency is set we honor it and emit a one-line
// stderr deprecation notice. Mirrors resolveScanProbeOptions for
// the --probe-workers / --probe-concurrency pair.
func resolveScanWorkers(fs *flag.FlagSet, workers, concurrency *int) int {
	workersSet := wasFlagPassed(fs, constants.FlagScanWorkers)
	concSet := wasFlagPassed(fs, constants.FlagScanWorkersConcurrencyAlias)
	if !workersSet && concSet {
		fmt.Fprint(os.Stderr, constants.MsgScanWorkersConcurrencyAlias)

		return *concurrency
	}

	return *workers
}

// resolveScanProbeOptions reconciles the deprecated --probe-concurrency
// against the unified --probe-workers. The new flag wins when both are
// set; when only the deprecated one is set we honor it and emit a
// one-line stderr deprecation notice. Depth comes through unchanged.
func resolveScanProbeOptions(fs *flag.FlagSet, noProbe, noWait *bool,
	probeConc, probeWorkers, probeDepth *int) ScanProbeOptions {
	concSet := wasFlagPassed(fs, constants.ScanProbeFlagConcurrency)
	workersSet := wasFlagPassed(fs, constants.ScanProbeFlagProbeWorkers)
	conc := *probeWorkers
	if !workersSet && concSet {
		fmt.Fprint(os.Stderr, constants.MsgScanProbeConcurrencyAlias)
		conc = *probeConc
	}

	return ScanProbeOptions{
		Disable:        *noProbe,
		NoWait:         *noWait,
		Concurrency:    conc,
		ConcurrencySet: workersSet || concSet,
		Depth:          *probeDepth,
	}
}

// wasFlagPassed reports whether the named flag was explicitly set on
// the command line (vs left at its default). Go's stdlib flag package
// doesn't surface this directly, so we walk Visit to find out.
func wasFlagPassed(fs *flag.FlagSet, name string) bool {
	seen := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			seen = true
		}
	})

	return seen
}

// registerScanBoolFlags registers boolean flags for the scan command.
func registerScanBoolFlags(fs *flag.FlagSet) (*bool, *bool, *bool) {
	ghDesktopFlag := fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	openFlag := fs.Bool("open", false, constants.FlagDescOpen)
	quietFlag := fs.Bool("quiet", false, constants.FlagDescQuiet)

	return ghDesktopFlag, openFlag, quietFlag
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
}

// parseCloneFlags parses flags for the clone command.
func parseCloneFlags(args []string) CloneFlags {
	fs := flag.NewFlagSet(constants.CmdClone, flag.ExitOnError)
	targetFlag := fs.String("target-dir", constants.DefaultDir, constants.FlagDescTargetDir)
	safePullFlag := fs.Bool("safe-pull", false, constants.FlagDescSafePull)
	ghDesktopFlag := fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	noReplaceFlag := fs.Bool("no-replace", false, constants.FlagDescCloneNoReplace)
	auditFlag := fs.Bool(constants.CloneFlagAudit, false, constants.FlagDescCloneAudit)
	maxConcFlag := fs.Int(constants.CloneFlagMaxConcurrency,
		constants.CloneDefaultMaxConcurrency, constants.FlagDescCloneMaxConcurrency)
	sshKeyFlag := fs.String("ssh-key", "", "SSH key name for clone")
	fs.StringVar(sshKeyFlag, "K", "", "SSH key name (short)")
	// Reuse the scan command's `--default-branch` constant + description
	// verbatim. The two flags share the same role (fallback branch when
	// detection finds nothing); keeping one source of truth means
	// `gitmap scan --help` and `gitmap clone --help` cannot drift.
	defaultBranchFlag := fs.String(constants.FlagScanDefaultBranch, "", constants.FlagDescScanDefaultBranch)
	outputFlag := fs.String(constants.FlagCloneTermOutput, "", constants.FlagDescCloneTermOutput)
	verifyFlag := fs.Bool(constants.FlagCloneVerifyCmdFaithful, false,
		constants.FlagDescCloneVerifyCmdFaithful)
	verifyExitFlag := fs.Bool(constants.FlagCloneVerifyCmdFaithfulExitOnMismatch,
		false, constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch)
	printArgvFlag := fs.Bool(constants.FlagClonePrintArgv, false,
		constants.FlagDescClonePrintArgv)
	noVSCodeSyncFlag := fs.Bool(constants.FlagNoVSCodeSync, false,
		constants.FlagDescNoVSCodeSync)
	debugPathsFlag := fs.Bool(constants.FlagDebugPaths, false,
		constants.FlagDescDebugPaths)
	sshFlag := fs.Bool("ssh", false,
		"Force every clone URL into `git@host:owner/repo.git` SSH-shorthand form before git runs (auto-converts HTTPS / `ssh://` URLs)")
	fs.BoolVar(sshFlag, "sh", false, "Short alias for --ssh")
	httpsFlag := fs.Bool("https", false,
		"Force every clone URL into `https://host/owner/repo.git` form (auto-converts SSH-shorthand / `ssh://` URLs)")
	fs.BoolVar(httpsFlag, "ht", false, "Short alias for --https")
	dryRunFlag := fs.Bool(constants.FlagCloneDryRun, false, constants.FlagDescCloneDryRun)
	fs.BoolVar(dryRunFlag, constants.FlagCloneDryRunShort, false, "Short alias for --dry-run")
	yesFlag := fs.Bool(constants.FlagCloneYes, false, constants.FlagDescCloneYes)
	fs.BoolVar(yesFlag, constants.FlagCloneYesShort, false, constants.FlagDescCloneYes)
	// Reorder so `gitmap clone <url> --ssh` works — Go's flag pkg
	// stops parsing at the first non-flag, which would otherwise
	// silently drop `--ssh` / `--https` / every other bool flag
	// when it follows the URL positional.
	fs.Parse(reorderFlagsBeforeArgs(args))

	applyDebugPathsEnv(*debugPathsFlag)

	return CloneFlags{
		Source:                          resolveCloneSource(fs),
		FolderName:                      resolveCloneFolderName(fs),
		TargetDir:                       *targetFlag,
		SSHKeyName:                      *sshKeyFlag,
		DefaultBranch:                   *defaultBranchFlag,
		Positional:                      fs.Args(),
		SafePull:                        *safePullFlag,
		GHDesktop:                       *ghDesktopFlag,
		NoReplace:                       *noReplaceFlag,
		Verbose:                         *verboseFlag,
		Audit:                           *auditFlag,
		MaxConcurrency:                  *maxConcFlag,
		Output:                          *outputFlag,
		VerifyCmdFaithful:               *verifyFlag,
		VerifyCmdFaithfulExitOnMismatch: *verifyExitFlag,
		PrintCloneArgv:                  *printArgvFlag,
		NoVSCodeSync:                    *noVSCodeSyncFlag,
		UseSSH:                          *sshFlag,
		UseHTTPS:                        *httpsFlag,
		DryRun:                          *dryRunFlag,
		IsAssumeYes:                     *yesFlag,
	}
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
	if fs.NArg() > 1 {
		second := fs.Arg(1)
		if isLikelyURL(second) {
			return ""
		}

		return second
	}

	return ""
}

// isLikelyURL is a cheap prefix check used to disambiguate
// "folder name" vs "second URL" without importing the clone package.
// Mirrors isDirectURL in clone.go — keep both in sync.
func isLikelyURL(s string) bool {
	lower := strings.ToLower(strings.TrimSpace(s))

	return strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "ssh://") ||
		strings.HasPrefix(lower, "git@")
}
