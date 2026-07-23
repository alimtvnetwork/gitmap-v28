package cmd

// CLI entry point for `gitmap clone-pick <repo-url> <paths>` (spec
// 100, v3.153.0+). Sparse-checkout a subset of a git repo into the
// current working directory (or --dest), and auto-save the selection
// to the CloneInteractiveSelection table.
//
// Exit codes:
//
//   0   -- dry-run rendered OR clone succeeded
//   1   -- runtime failure (git, fs, db)
//   2   -- bad CLI usage (missing args, invalid flag value)
//   130 -- user canceled the picker
//
// Flag binding lives in clonepick_flags.go, the picker glue in
// clonepick_picker.go, and the side-effecting executor in
// clonepick_execute.go -- keeping this file focused on the
// dispatcher itself (and under the strict 200-line cap).

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runClonePick is the dispatcher entry registered in rootcore.go.
func runClonePick(args []string) {
	checkHelp("clone-pick", args)

	parsed := parseClonePickFlags(args)
	setCmdFaithfulVerify(parsed.VerifyCmdFaithful)
	setCmdFaithfulExitOnMismatch(parsed.VerifyCmdFaithfulExitOnMismatch)
	setCmdPrintArgv(parsed.PrintCloneArgv)

	plan, replayId, err := buildClonePickPlan(parsed)
	if err != nil {
		cliexit.Fail(constants.CmdClonePick, "parse-args", parsed.RawURL, err, 2)
	}
	plan = maybeRunClonePickPicker(plan, parsed.Flags.Ask)

	if plan.DryRun {
		runClonePickDryRun(plan, parsed)

		return
	}

	if parsed.Output == constants.OutputTerminal {
		printClonePickTermBlock(plan)
	}
	runClonePickExecute(plan, parsed.NoVSCodeSync, replayId)
}

// runClonePickDryRun handles the --dry-run branch. Split out so the
// dispatcher stays under the function-length cap.
func runClonePickDryRun(plan clonepick.Plan, parsed clonePickParsed) {
	if parsed.Output == constants.OutputTerminal {
		printClonePickTermBlock(plan)
		maybeExitOnCmdFaithfulMismatch()

		return
	}
	if err := clonepick.Render(os.Stdout, plan); err != nil {
		cliexit.Fail(constants.CmdClonePick, "render-dry-run", parsed.RawURL, err, 1)
	}
	maybeExitOnCmdFaithfulMismatch()
}

// buildClonePickPlan picks between the parse path (positional args)
// and the replay path (--replay <id|name> hits the DB). Returns the
// Plan plus the replayed SelectionId (0 for fresh runs) so the
// executor can bump CreatedAt without re-deriving the id.
func buildClonePickPlan(parsed clonePickParsed) (clonepick.Plan, int64, error) {
	if len(parsed.Flags.Replay) == 0 {
		plan, err := clonepick.ParseArgs(parsed.RawURL, parsed.RawPaths, parsed.Flags)

		return plan, 0, err
	}
	loader, err := openDB()
	if err != nil {
		return clonepick.Plan{}, 0, err
	}
	plan, replayId, loadErr := clonepick.LoadFromDB(loader, parsed.Flags.Replay)
	if loadErr != nil {
		return clonepick.Plan{}, 0, loadErr
	}
	applyReplayOverrides(&plan, parsed.Flags)

	return plan, replayId, nil
}

// applyReplayOverrides lets runtime-only flags from THIS invocation
// override the persisted Plan -- spec rule: replay reproduces the
// selection, not the verbosity / dest choice.
func applyReplayOverrides(plan *clonepick.Plan, flags clonepick.Flags) {
	plan.DryRun = flags.DryRun
	plan.Quiet = flags.Quiet
	plan.Force = flags.Force
	if len(flags.Dest) > 0 && flags.Dest != "." {
		plan.DestDir = flags.Dest
	}
}
