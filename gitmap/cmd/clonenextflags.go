package cmd

import (
	"flag"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CloneNextFlags bundles every parsed flag from the clone-next command so
// the dispatcher in runCloneNext can branch on batch vs single mode without
// a 9-arg return list.
type CloneNextFlags struct {
	VersionArg   string
	Delete       bool
	Keep         bool
	NoDesktop    bool
	CreateRemote bool
	SSHKeyName   string
	Verbose      bool
	CSVPath      string
	All          bool
	// Force forces a flat clone even when the user's cwd IS the target
	// folder. Triggers a chdir-to-parent before the existence check (to
	// release Windows file locks) and DISABLES the versioned-folder
	// fallback so the user gets either a flat layout or a clear error.
	// See spec/01-app/87-clone-next-flatten.md.
	Force bool
	// MaxConcurrency is the worker-pool size for batch mode (--all / --csv).
	// 1 (the default) preserves the historical sequential behavior so
	// stdout ordering of per-repo lines is deterministic. Values >1 fan
	// repos out across a bounded pool that mirrors the main cloner's
	// pattern (see gitmap/cloner/concurrent.go). Ignored in single-repo
	// mode where there is only one unit of work.
	MaxConcurrency int
	// NoProgress suppresses the live per-repo progress line printed
	// by the batch collector as workers finish. The final summary
	// (ok/failed/skipped totals) always prints regardless. Default
	// false so users get progress feedback out-of-the-box.
	NoProgress bool
	// ReportErrors enables a JSON failure report at command exit
	// when any per-repo clone fails. Off by default; mirrors the
	// `gitmap scan --errors-report` flag for consistent UX.
	ReportErrors bool
	// DryRun, when true, prints the would-be `git clone` commands
	// (single-repo + batch) and skips ALL side effects — no actual
	// clone, no folder removal, no DB write, no GH Desktop / VS Code
	// launch, no shell handoff. See FlagCloneNextDryRun.
	DryRun bool
	// Output selects the per-repo summary format. Empty keeps the
	// legacy terse stage messages; "terminal" additionally emits
	// the standardized RepoTermBlock right before the clone, so the
	// shape matches scan/clone-from/probe.
	Output string
	// VerifyCmdFaithful enables the dry-run argv-vs-displayed checker.
	VerifyCmdFaithful bool
	// VerifyCmdFaithfulExitOnMismatch upgrades the verifier into a
	// hard failure: any divergence sets a sticky bit and the run tail
	// exits with constants.CloneVerifyCmdFaithfulExitCode. Implies
	// VerifyCmdFaithful.
	VerifyCmdFaithfulExitOnMismatch bool
	// PrintCloneArgv dumps the executor argv to stderr.
	PrintCloneArgv bool
	// NoVSCodeSync suppresses the post-clone update of the
	// alefragnani.project-manager projects.json file. Mirrors
	// `gitmap scan --no-vscode-sync`. Default false. See
	// spec/01-vscode-project-manager-sync/02-clone-sync.md.
	NoVSCodeSync bool
}

type cloneNextPointers struct {
	delFlag          *bool
	kpFlag           *bool
	noDesk           *bool
	createRem        *bool
	sshKey           *string
	verboseFlag      *bool
	csvPath          *string
	allFlag          *bool
	forceFlag        *bool
	maxConcFlag      *int
	noProgressFlag   *bool
	reportErrFlag    *bool
	dryRunFlag       *bool
	outputFlag       *string
	verifyFlag       *bool
	verifyExitFlag   *bool
	printArgvFlag    *bool
	noVSCodeSyncFlag *bool
}

// parseCloneNextFlags parses flags for the clone-next command.
func parseCloneNextFlags(args []string) CloneNextFlags {
	fs := flag.NewFlagSet(constants.CmdCloneNext, flag.ExitOnError)
	var p cloneNextPointers
	bindCloneNextFlags(fs, &p)
	fs.Parse(reorderFlagsBeforeArgs(args))
	conc := resolveCloneNextConcurrency(*p.maxConcFlag)
	out := buildCloneNextFlags(&p, conc)
	if fs.NArg() > 0 {
		out.VersionArg = fs.Arg(0)
	}
	return out
}

func resolveCloneNextConcurrency(maxConc int) int {
	resolvedConc, ok := cloneconcurrency.Resolve(maxConc)
	if !ok {
		return 1
	}
	return resolvedConc
}

func bindCloneNextFlags(fs *flag.FlagSet, p *cloneNextPointers) {
	bindCloneNextBasicFlags(fs, p)
	bindCloneNextBatchFlags(fs, p)
	bindCloneNextVerifyFlags(fs, p)
}

func bindCloneNextBasicFlags(fs *flag.FlagSet, p *cloneNextPointers) {
	p.delFlag = fs.Bool("delete", false, constants.FlagDescCloneNextDelete)
	p.kpFlag = fs.Bool("keep", false, constants.FlagDescCloneNextKeep)
	p.noDesk = fs.Bool("no-desktop", false, constants.FlagDescCloneNextNoDesktop)
	p.createRem = fs.Bool("create-remote", false, constants.FlagDescCloneNextCreateRemote)
	p.sshKey = fs.String("ssh-key", "", "SSH key name for clone")
	fs.StringVar(p.sshKey, "K", "", "SSH key name (short)")
	p.verboseFlag = fs.Bool("verbose", false, constants.FlagDescVerbose)
}

func bindCloneNextBatchFlags(fs *flag.FlagSet, p *cloneNextPointers) {
	p.csvPath = fs.String("csv", "", constants.FlagDescCloneNextCSV)
	p.allFlag = fs.Bool("all", false, constants.FlagDescCloneNextAll)
	p.forceFlag = fs.Bool("force", false, constants.FlagDescCloneNextForce)
	fs.BoolVar(p.forceFlag, "f", false, constants.FlagDescCloneNextForce)
	p.maxConcFlag = fs.Int(constants.CloneFlagMaxConcurrency, constants.CloneDefaultMaxConcurrency, constants.FlagDescCloneMaxConcurrency)
	p.noProgressFlag = fs.Bool(constants.FlagCloneNextNoProgress, false, constants.FlagDescCloneNextNoProgress)
	p.reportErrFlag = fs.Bool(constants.FlagScanReportErrors, false, constants.FlagDescScanReportErrors)
}

func bindCloneNextVerifyFlags(fs *flag.FlagSet, p *cloneNextPointers) {
	p.dryRunFlag = fs.Bool(constants.FlagCloneNextDryRun, false, constants.FlagDescCloneNextDryRun)
	p.outputFlag = fs.String(constants.FlagCloneNextOutput, "", constants.FlagDescCloneNextOutput)
	p.verifyFlag = fs.Bool(constants.FlagCloneVerifyCmdFaithful, false, constants.FlagDescCloneVerifyCmdFaithful)
	p.verifyExitFlag = fs.Bool(constants.FlagCloneVerifyCmdFaithfulExitOnMismatch, false, constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch)
	p.printArgvFlag = fs.Bool(constants.FlagClonePrintArgv, false, constants.FlagDescClonePrintArgv)
	p.noVSCodeSyncFlag = fs.Bool(constants.FlagNoVSCodeSync, false, constants.FlagDescNoVSCodeSync)
}

func buildCloneNextFlags(p *cloneNextPointers, maxConcurrency int) CloneNextFlags {
	out := CloneNextFlags{
		Delete:         *p.delFlag,
		Keep:           *p.kpFlag,
		NoDesktop:      *p.noDesk,
		CreateRemote:   *p.createRem,
		SSHKeyName:     *p.sshKey,
		Verbose:        *p.verboseFlag,
		CSVPath:        *p.csvPath,
		All:            *p.allFlag,
		Force:          *p.forceFlag,
		MaxConcurrency: maxConcurrency,
	}
	populateCloneNextOptionalFlags(&out, p)
	return out
}

func populateCloneNextOptionalFlags(out *CloneNextFlags, p *cloneNextPointers) {
	out.NoProgress = *p.noProgressFlag
	out.ReportErrors = *p.reportErrFlag
	out.DryRun = *p.dryRunFlag
	out.Output = *p.outputFlag
	out.VerifyCmdFaithful = *p.verifyFlag
	out.VerifyCmdFaithfulExitOnMismatch = *p.verifyExitFlag
	out.PrintCloneArgv = *p.printArgvFlag
	out.NoVSCodeSync = *p.noVSCodeSyncFlag
}
