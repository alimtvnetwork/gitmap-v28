package cmd

// clonepick_execute.go: side-effecting half of `gitmap clone-pick`
// (open DB, run the sparse-checkout, persist the selection, sync to
// VS Code Project Manager). Split out of clonepick.go so the
// dispatcher file stays under the 200-line cap.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runClonePickExecute opens the DB (best-effort), runs the
// sparse-checkout, and translates the Result to an exit code.
// replayId is non-zero when the Plan came from --replay; on success
// CreatedAt is bumped so most-recently-replayed sorts to the top.
func runClonePickExecute(plan clonepick.Plan, noVSCodeSync bool, replayId int64) {
	progress := io.Writer(os.Stderr)
	if plan.Quiet {
		progress = io.Discard
	}

	db, dbErr := openDB()
	if dbErr != nil {
		// DB open failure is non-fatal -- clone still proceeds, just
		// without persistence. Per the zero-swallow policy we surface
		// the error to stderr so it isn't silently dropped.
		fmt.Fprintln(os.Stderr, dbErr)
	}

	result := clonepick.Execute(plan, db, progress)
	announceClonePickPersistence(plan, result, replayId, db)

	if result.Status == clonepick.StatusFailed {
		maybeExitOnCmdFaithfulMismatch()
		os.Exit(1)
	}

	// VS Code Project Manager sync. Result.Detail carries the
	// resolved destination path on success; fall back to the plan's
	// DestDir when empty.
	syncClonePickResultToVSCodePM(plan, result, noVSCodeSync)

	if plan.DestDir != "." && plan.DestDir != "" {
		WriteShellHandoff(result.Detail)
	}
	maybeExitOnCmdFaithfulMismatch()
}

// announceClonePickPersistence prints the saved/replayed line and,
// for replays, bumps the CreatedAt column. Split out so the main
// executor stays under the function-length cap.
func announceClonePickPersistence(plan clonepick.Plan, result clonepick.Result, replayId int64, db *store.DB) {
	name := plan.Name
	if len(name) == 0 {
		name = "(unnamed)"
	}
	switch {
	case replayId > 0 && result.Status == clonepick.StatusOK:
		fmt.Fprintf(os.Stderr, constants.MsgClonePickReplayed,
			replayId, plan.RepoCanonicalId, name)
		if !plan.DryRun && db != nil {
			if err := clonepick.TouchAfterReplay(db, replayId, plan.DryRun); err != nil {
				cliexit.Reportf(constants.CmdClonePick, "touch-replay",
					strconv.FormatInt(replayId, 10), err)
			}
		}
	case result.SelectionId > 0:
		fmt.Fprintf(os.Stderr, constants.MsgClonePickSaved,
			result.SelectionId, plan.RepoCanonicalId, name)
	}
}

// syncClonePickResultToVSCodePM registers a successful sparse-checkout
// dest in projects.json. The pick name (when set) wins over the folder
// basename so users see their alias in the Project Manager sidebar.
func syncClonePickResultToVSCodePM(plan clonepick.Plan, result clonepick.Result, skip bool) {
	dest := result.Detail
	if dest == "" {
		dest = plan.DestDir
	}
	abs, err := filepath.Abs(dest)
	if err != nil {
		abs = dest
	}
	name := plan.Name
	if name == "" {
		name = filepath.Base(abs)
	}
	syncSingleClonedRepoToVSCodePM(abs, name, skip)
}
