package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
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

// parseCloneNextFlags parses flags for the clone-next command.
func parseCloneNextFlags(args []string) CloneNextFlags {
	fs := flag.NewFlagSet(constants.CmdCloneNext, flag.ExitOnError)
	delFlag := fs.Bool("delete", false, constants.FlagDescCloneNextDelete)
	kpFlag := fs.Bool("keep", false, constants.FlagDescCloneNextKeep)
	noDesk := fs.Bool("no-desktop", false, constants.FlagDescCloneNextNoDesktop)
	createRem := fs.Bool("create-remote", false, constants.FlagDescCloneNextCreateRemote)
	sshKey := fs.String("ssh-key", "", "SSH key name for clone")
	fs.StringVar(sshKey, "K", "", "SSH key name (short)")
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	csvPath := fs.String("csv", "", constants.FlagDescCloneNextCSV)
	allFlag := fs.Bool("all", false, constants.FlagDescCloneNextAll)
	// Force-flatten: long --force and short -f both bind to the same
	// bool var so either form is canonical (mirrors the --ssh-key/-K
	// pairing convention used elsewhere in this flagset).
	forceFlag := fs.Bool("force", false, constants.FlagDescCloneNextForce)
	fs.BoolVar(forceFlag, "f", false, constants.FlagDescCloneNextForce)
	// Batch worker-pool size. Reuses the same flag name as `gitmap clone`
	// so users learn one name (--max-concurrency / -j is reserved for a
	// later short alias if needed). Default 1 = sequential.
	maxConcFlag := fs.Int(constants.CloneFlagMaxConcurrency,
		constants.CloneDefaultMaxConcurrency, constants.FlagDescCloneMaxConcurrency)
	noProgressFlag := fs.Bool(constants.FlagCloneNextNoProgress, false,
		constants.FlagDescCloneNextNoProgress)
	reportErrFlag := fs.Bool(constants.FlagScanReportErrors, false, constants.FlagDescScanReportErrors)
	dryRunFlag := fs.Bool(constants.FlagCloneNextDryRun, false, constants.FlagDescCloneNextDryRun)
	outputFlag := fs.String(constants.FlagCloneNextOutput, "", constants.FlagDescCloneNextOutput)
	verifyFlag := fs.Bool(constants.FlagCloneVerifyCmdFaithful, false,
		constants.FlagDescCloneVerifyCmdFaithful)
	verifyExitFlag := fs.Bool(constants.FlagCloneVerifyCmdFaithfulExitOnMismatch,
		false, constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch)
	printArgvFlag := fs.Bool(constants.FlagClonePrintArgv, false,
		constants.FlagDescClonePrintArgv)
	noVSCodeSyncFlag := fs.Bool(constants.FlagNoVSCodeSync, false,
		constants.FlagDescNoVSCodeSync)
	// Reorder so flags placed AFTER the positional version (e.g.
	// `gitmap cn v+1 -f`) are still recognized. Go's stdlib flag
	// parser stops at the first non-flag arg, so without this the
	// `-f` in the screenshot above was silently dropped.
	fs.Parse(reorderFlagsBeforeArgs(args))

	resolvedConc, ok := cloneconcurrency.Resolve(*maxConcFlag)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrCloneMaxConcurrencyInvalid, *maxConcFlag)
		os.Exit(1)
	}

	out := CloneNextFlags{
		Delete:                          *delFlag,
		Keep:                            *kpFlag,
		NoDesktop:                       *noDesk,
		CreateRemote:                    *createRem,
		SSHKeyName:                      *sshKey,
		Verbose:                         *verboseFlag,
		CSVPath:                         *csvPath,
		All:                             *allFlag,
		Force:                           *forceFlag,
		MaxConcurrency:                  resolvedConc,
		NoProgress:                      *noProgressFlag,
		ReportErrors:                    *reportErrFlag,
		DryRun:                          *dryRunFlag,
		Output:                          *outputFlag,
		VerifyCmdFaithful:               *verifyFlag,
		VerifyCmdFaithfulExitOnMismatch: *verifyExitFlag,
		PrintCloneArgv:                  *printArgvFlag,
		NoVSCodeSync:                    *noVSCodeSyncFlag,
	}
	if fs.NArg() > 0 {
		out.VersionArg = fs.Arg(0)
	}

	return out
}
