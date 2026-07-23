package cmd

// CLI entry point for `gitmap clone-from <file>`. Parses flags,
// loads the plan, renders dry-run by default, executes only when
// --execute is passed, then prints the summary.
//
// Exit codes (mirrors gitmap conventions):
//
//   0 — dry-run completed OR every row was ok/skipped on execute
//   1 — file open / parse error, OR any row failed on execute
//   2 — bad CLI usage (missing <file> argument)
//
// The split between exit-1 (something failed during the operation)
// and exit-2 (caller error) lets shell scripts distinguish "you
// invoked me wrong" from "I tried but git rejected one of the
// URLs" — the second is recoverable by editing the input file,
// the first needs a different command.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// cloneFromFlags holds parsed CLI inputs. Grouped in a struct
// (rather than passed individually) so future additions don't
// churn every helper signature.
type cloneFromFlags struct {
	file     string
	execute  bool
	quiet    bool
	noReport bool
	// emitSchema, when non-empty, short-circuits the command: the
	// requested JSON Schema (kind = "report" | "input") is written
	// to stdout and the process exits 0. No <file> argument is
	// required in this mode — useful for CI tooling that wants to
	// validate exported manifests without running a clone first.
	emitSchema        string
	output            string
	checkout          string
	verifyCmdFaithful bool
	// verifyCmdFaithfulExitOnMismatch upgrades the verifier into a
	// hard failure: any divergence sets a sticky bit and the run tail
	// exits with constants.CloneVerifyCmdFaithfulExitCode. Implies
	// verifyCmdFaithful.
	verifyCmdFaithfulExitOnMismatch bool
	// printCloneArgv dumps the executor argv to stderr.
	printCloneArgv bool
	// noVSCodeSync suppresses the post-clone update of the
	// alefragnani.project-manager projects.json file. Mirrors
	// `gitmap scan --no-vscode-sync`. Default false. See
	// spec/01-vscode-project-manager-sync/02-clone-sync.md.
	noVSCodeSync bool
	// maxConcurrency is the resolved worker-pool size. The parser
	// runs cloneconcurrency.Resolve so the value seen here is
	// always >=1 (0=auto becomes NumCPU at parse time). Increasing
	// N preserves the on-disk hierarchy because each worker still
	// uses the row's Dest / DeriveDest verbatim — only progress-line
	// timing changes.
	maxConcurrency int
}

// runCloneFrom is the dispatcher entry. checkHelp handles `--help`
// per the project's help-system convention before any flag parsing
// so unparseable flags don't suppress the help text.
func runCloneFrom(args []string) {
	checkHelp("clone-from", args)
	cfg := parseCloneFromFlags(args)
	if cfg.emitSchema != "" {
		runCloneFromEmitSchema(cfg.emitSchema)

		return
	}
	setCmdFaithfulVerify(cfg.verifyCmdFaithful)
	setCmdFaithfulExitOnMismatch(cfg.verifyCmdFaithfulExitOnMismatch)
	setCmdPrintArgv(cfg.printCloneArgv)
	plan, err := clonefrom.ParseFile(cfg.file)
	if err != nil {
		cliexit.Fail(constants.CmdCloneFrom, "parse-manifest", cfg.file, err, 1)
	}
	applyCheckoutDefault(&plan, cfg.checkout)
	if !cfg.execute {
		runCloneFromDry(plan, cfg)
		maybeExitOnCmdFaithfulMismatch()

		return
	}
	runCloneFromExecute(plan, cfg)
}

// runCloneFromEmitSchema lives in clonefrom_emitschema.go to keep
// this file under the 200-line cap.

// applyCheckoutDefault and validateCheckoutFlag live in
// clonefrom_checkout.go to keep this file under the 200-line cap.

// parseCloneFromFlags lives in clonefrom_flags.go to keep this file
// under the 200-line cap (mem://style/code-constraints, item 3).

// validateCheckoutFlag lives in clonefrom_checkout.go.

// runCloneFromDry renders the dry-run preview and exits with the
// dry-run conventional code (0 = "I would do these things"). No
// side effects — by design, a dry-run never touches the network
// or the filesystem outside of READING the input file.
func runCloneFromDry(plan clonefrom.Plan, cfg cloneFromFlags) {
	render := pickCloneFromRenderer(cfg.output)
	if err := render(os.Stdout, plan); err != nil {
		cliexit.Fail(constants.CmdCloneFrom, "render-dry-run", cfg.file, err, 1)
	}
}

// pickCloneFromRenderer dispatches between the legacy and the
// standardized terminal renderer based on --output. Anything other
// than "terminal" (including the empty default) keeps the legacy
// output so existing scripts and doc snippets stay byte-identical.
func pickCloneFromRenderer(output string) func(io.Writer, clonefrom.Plan) error {
	if output == constants.OutputTerminal {
		return clonefrom.RenderTerminal
	}

	return clonefrom.Render
}

// runCloneFromExecute is the side-effecting branch. Picks the
// progress writer based on --quiet, executes the plan, writes the
// CSV report (unless --no-report), prints the summary, then
// translates the result tally to an exit code.
func runCloneFromExecute(plan clonefrom.Plan, cfg cloneFromFlags) {
	progress := io.Writer(os.Stderr)
	if cfg.quiet {
		progress = io.Discard
	}
	// `--output terminal`: stream one standardized RepoTermBlock per
	// row via ExecuteWithHooks's BeforeRow callback — printed
	// IMMEDIATELY before that row's `git clone` shells out. This
	// interleaves per-repo previews with clone progress instead of
	// dumping every block upfront, matching the URL-driven `clone
	// <urls...>` behavior. Dry-run still uses RenderTerminal upfront
	// (no execution to interleave with). A nil hook keeps the legacy
	// path byte-identical for callers that didn't opt in.
	var hook clonefrom.BeforeRowHook
	if cfg.output == constants.OutputTerminal {
		hook = printCloneFromTermBlockRow
	}
	// Dispatch sequential vs parallel on the resolved worker count.
	// 0=auto becomes NumCPU at parse time, so any value reaching
	// here is >=1. The concurrent runner short-circuits to
	// ExecuteWithHooks for workers <=1.
	var results []clonefrom.Result
	if cfg.maxConcurrency > 1 {
		fmt.Fprintf(os.Stderr, constants.MsgCloneConcurrencyEnabledFmt, cfg.maxConcurrency)
		results = clonefrom.ExecuteWithHooksConcurrent(plan, "", progress, hook, cfg.maxConcurrency)
	} else {
		results = clonefrom.ExecuteWithHooks(plan, "", progress, hook)
	}
	csvPath, jsonPath := writeCloneFromReports(results, cfg)
	syncCloneFromResultsToVSCodePM(results, cfg.noVSCodeSync)
	if cfg.output == constants.OutputTerminal {
		if err := clonefrom.RenderSummaryTerminal(os.Stdout, results, csvPath, jsonPath); err != nil {
			cliexit.Reportf(constants.CmdCloneFrom, "render-summary", csvPath, err)
		}
	} else if err := clonefrom.RenderSummary(os.Stdout, results, csvPath); err != nil {
		cliexit.Reportf(constants.CmdCloneFrom, "render-summary", csvPath, err)
	}
	maybeExitOnCmdFaithfulMismatch()
	os.Exit(cloneFromExitCode(results))
}

// writeCloneFromReports lives in clonefrom_reports.go to keep this
// file under the project's 200-line cap.

// cloneFromExitCode returns 1 if any row failed, else 0. Skipped
// rows are NOT failures — re-running an idempotent plan with all
// destinations already in place is a successful no-op.
func cloneFromExitCode(results []clonefrom.Result) int {
	for _, r := range results {
		if r.Status == constants.CloneFromStatusFailed {

			return 1
		}
	}

	return 0
}

// syncCloneFromResultsToVSCodePM filters status=ok rows and pushes
// them into projects.json in a single batched Sync. Mirrors the
// pattern used by executeClone for manifest-style inputs.
func syncCloneFromResultsToVSCodePM(results []clonefrom.Result, skip bool) {
	pairs := make([]vscodepm.Pair, 0, len(results))
	for _, r := range results {
		if r.Status != constants.CloneFromStatusOK {
			continue
		}
		abs, err := filepath.Abs(r.Dest)
		if err != nil {
			abs = r.Dest
		}
		name := filepath.Base(abs)
		pairs = append(pairs, buildClonePMPair(abs, name))
	}
	syncClonedReposToVSCodePM(pairs, skip)
}
