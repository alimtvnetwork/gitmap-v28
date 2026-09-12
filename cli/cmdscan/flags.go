package cmdscan

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ScanProbeOptions bundles the flags that govern the optional
// background version-probe pass scan kicks off after upserting repos.
type ScanProbeOptions struct {
	Disable        bool
	NoWait         bool
	Concurrency    int
	ConcurrencySet bool
	Depth          int
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

// ParseScanFlags parses flags for the scan command.
func ParseScanFlags(args []string) (dir, configPath, mode, output, outFile, outputPath, relativeRoot, defaultBranch string, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, compact, fix bool, workers, maxDepth int, probeOpts ScanProbeOptions) {
	fs := flag.NewFlagSet(constants.CmdScan, flag.ExitOnError)
	scanFlags := newScanFlagSet(fs)
	_ = fs.Parse(args)

	dir = resolveScanDir(fs)
	probeOpts = resolveScanProbeOptions(fs, scanFlags.noProbeFlag, scanFlags.noProbeWaitFlag,
		scanFlags.probeConcFlag, scanFlags.probeWorkersFlag, scanFlags.probeDepthFlag)
	resolvedWorkers := resolveScanWorkers(fs, scanFlags.workersFlag, scanFlags.concurrencyFlag)
	resolvedOutputPath := resolveScanOutputPath(*scanFlags.outputPathFlag, *scanFlags.manifestFlag)

	return dir, *scanFlags.cfgFlag, *scanFlags.modeFlag, *scanFlags.outputFlag, *scanFlags.outFileFlag, resolvedOutputPath, *scanFlags.relRootFlag, *scanFlags.defaultBranchFlag, *scanFlags.ghDesktopFlag, *scanFlags.openFlag, *scanFlags.quietFlag, *scanFlags.noVSCodeSyncFlag, *scanFlags.noAutoTagsFlag, *scanFlags.reportErrFlag, *scanFlags.compactFlag, *scanFlags.fixFlag, resolvedWorkers, *scanFlags.maxDepthFlag, probeOpts
}

func parseScanFlags(args []string) (dir, configPath, mode, output, outFile, outputPath, relativeRoot, defaultBranch string, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, compact, fix bool, workers, maxDepth int, probeOpts ScanProbeOptions) {
	return ParseScanFlags(args)
}

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

func wasFlagPassed(fs *flag.FlagSet, flagName string) bool {
	hasSeen := false
	fs.Visit(func(flagItem *flag.Flag) {
		if flagItem.Name == flagName {
			hasSeen = true
		}
	})

	return hasSeen
}

func resolveScanDir(fs *flag.FlagSet) string {
	if fs.NArg() > 0 {
		return fs.Arg(0)
	}

	return constants.DefaultDir
}
