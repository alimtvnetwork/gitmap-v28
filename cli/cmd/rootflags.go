package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
type CloneFlags = cmdclone.CloneFlags

// isLikelyURL is a cheap prefix check used to disambiguate
// "folder name" vs "second URL".
func isLikelyURL(rawURL string) bool {
	return cmdclone.IsLikelyURL(rawURL)
}
