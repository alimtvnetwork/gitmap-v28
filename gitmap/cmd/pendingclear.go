// Package cmd — `gitmap pending clear` removes orphaned or illegal
// pending tasks so the next clone run is not blocked by a leftover
// entry from an earlier crash.
//
// Modes (see helptext/pending-clear.md for the full contract):
//
//	orphans  — TargetPath is missing on disk
//	illegal  — TargetPath looks like a URL or contains illegal Windows
//	           path characters (`:` after drive letter, `?`, `*`, etc.)
//	all      — every pending task
//	<id>     — a single task by numeric ID
//
// Default mode is `orphans` because it's the safest auto-cleanup.
// Confirmation is required unless --yes is passed; --dry-run previews
// without touching the DB.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// mustParsePendingClearArgs parses and validates arguments or exits on failure.
func mustParsePendingClearArgs(args []string) (string, bool, bool, int64) {
	mode, dryRun, yes, idMatch, err := parsePendingClearArgs(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	return mode, dryRun, yes, idMatch
}

// mustOpenPendingDB opens the store database or exits with a warning.
func mustOpenPendingDB() *store.DB {
	db, dbErr := openDB()
	if dbErr != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingDBOpen, dbErr)
		os.Exit(1)
	}

	return db
}

// mustListPendingTasks fetches pending tasks or exits on database error.
func mustListPendingTasks(db *store.DB) []model.PendingTaskRecord {
	tasks, listErr := db.ListPendingTasks()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskQuery, listErr)
		os.Exit(1)
	}

	return tasks
}

// executePendingClear handles dry-run, confirmation prompt, and candidate deletion.
func executePendingClear(db *store.DB, candidates []pendingClearCandidate, dryRun, yes bool) {
	if dryRun {
		fmt.Printf(constants.MsgPendingClearDryRun, len(candidates))

		return
	}
	if !yes && !confirmPendingClear(len(candidates)) {
		fmt.Print(constants.MsgPendingClearAborted)
		os.Exit(1)
	}
	deleted := deletePendingClearCandidates(db, candidates)
	fmt.Printf(constants.MsgPendingClearDone, deleted, len(candidates))
}

// processPendingClear filters candidates and initiates cleanup workflow.
func processPendingClear(db *store.DB, tasks []model.PendingTaskRecord,
	mode string, idMatch int64, dryRun, yes bool) {
	candidates := selectClearCandidates(tasks, mode, idMatch)
	printPendingClearHeader(mode, len(tasks))
	if len(candidates) == 0 {
		fmt.Printf(constants.MsgPendingClearNoMatches, mode)

		return
	}
	printPendingClearCandidates(candidates)
	executePendingClear(db, candidates, dryRun, yes)
}

// runPendingClear is wired in from runPending when args[0] == "clear".
// args is everything after the "clear" token.
func runPendingClear(args []string) {
	mode, dryRun, yes, idMatch := mustParsePendingClearArgs(args)
	db := mustOpenPendingDB()
	defer db.Close()

	tasks := mustListPendingTasks(db)
	processPendingClear(db, tasks, mode, idMatch, dryRun, yes)
}

// parsePendingClearID validates a string as numeric ID and flags invalid modes.
func parsePendingClearID(a string) (int64, error) {
	id, parseErr := strconv.ParseInt(a, 10, 64)
	isError := parseErr != nil || id <= 0
	if isError && strings.HasPrefix(a, "-") {
		return 0, fmt.Errorf(constants.ErrPendingClearUnknownMode, a)
	}
	if isError {
		return 0, fmt.Errorf(constants.ErrPendingClearBadID, a)
	}

	return id, nil
}

// applyPendingClearIDArg parses and updates mode and idMatch for numeric ID args.
func applyPendingClearIDArg(a string, mode *string, idMatch *int64) error {
	id, err := parsePendingClearID(a)
	if err != nil {
		return err
	}
	*mode = "id"
	*idMatch = id

	return nil
}

// applyPendingClearFlag handles boolean flags for dry-run and confirmation bypass.
func applyPendingClearFlag(a string, dryRun, yes *bool) bool {
	if a == "--dry-run" {
		*dryRun = true
		return true
	}
	if a == "--yes" || a == "-y" {
		*yes = true
		return true
	}

	return false
}

// applyPendingClearArg parses a single command line argument for pending clear.
func applyPendingClearArg(a string, mode *string, dryRun, yes *bool, idMatch *int64) error {
	if applyPendingClearFlag(a, dryRun, yes) {
		return nil
	}
	if a == "orphans" || a == "illegal" || a == "all" {
		*mode = a
		return nil
	}

	return applyPendingClearIDArg(a, mode, idMatch)
}

// parsePendingClearArgs splits args into (mode, dryRun, yes, idMatch, err).
// Mode is one of: orphans (default), illegal, all, or "id" (with idMatch
// holding the parsed numeric ID).
func parsePendingClearArgs(args []string) (string, bool, bool, int64, error) {
	mode := "orphans"
	var dryRun, yes bool
	var idMatch int64
	for _, a := range args {
		if err := applyPendingClearArg(a, &mode, &dryRun, &yes, &idMatch); err != nil {
			return "", false, false, 0, err
		}
	}

	return mode, dryRun, yes, idMatch, nil
}

// selectClearCandidates filters the full task list down to the rows
// matching the requested mode. Each returned candidate carries a
// human-readable reason that's printed in the preview.
func selectClearCandidates(tasks []model.PendingTaskRecord,
	mode string, idMatch int64) []pendingClearCandidate {
	out := make([]pendingClearCandidate, 0, len(tasks))
	for _, t := range tasks {
		reason, keep := classifyPendingClearTask(t, mode, idMatch)
		if !keep {
			continue
		}
		out = append(out, pendingClearCandidate{task: t, reason: reason})
	}

	return out
}

// classifyIllegalTask checks if a target path matches illegal URL shape or illegal chars.
func classifyIllegalTask(path string) (string, bool) {
	if isURLShapedTarget(path) {
		return constants.MsgPendingClearReasonURL, true
	}
	if hasIllegalPathChar(path) {
		return constants.MsgPendingClearReasonChar, true
	}

	return "", false
}

// classifyIDTask checks if the task ID matches the specified numeric filter.
func classifyIDTask(taskID, targetID int64) (string, bool) {
	if taskID == targetID {
		return constants.MsgPendingClearReasonByID, true
	}

	return "", false
}

// classifyOrphanTask checks if a target path is an orphan directory.
func classifyOrphanTask(path string) (string, bool) {
	if isOrphanTarget(path) {
		return constants.MsgPendingClearReasonOrph, true
	}

	return "", false
}

// classifyPendingClearTask decides whether one task matches the mode
// and returns the reason label for the preview output.
func classifyPendingClearTask(t model.PendingTaskRecord,
	mode string, idMatch int64) (string, bool) {
	if mode == "all" {
		return constants.MsgPendingClearReasonAll, true
	}
	if mode == "id" {
		return classifyIDTask(t.ID, idMatch)
	}
	if mode == "illegal" {
		return classifyIllegalTask(t.TargetPath)
	}

	return classifyOrphanTask(t.TargetPath)
}

// hasSchemePrefix checks if a lowercased path starts with or contains broken URI schemes.
func hasSchemePrefix(lower string) bool {
	for _, scheme := range []string{"http:", "https:", "ssh:", "git:"} {
		if strings.Contains(lower, scheme+`\`) || strings.Contains(lower, scheme+"/") {
			return true
		}
	}

	return false
}

// isURLShapedTarget catches paths shaped like the Windows-corrupted
// targets produced when PowerShell split a comma URL list (issue #11).
// Examples that match: "https:\github.com\...", "git@github.com:...",
// any TargetPath containing "://" anywhere after the first segment.
func isURLShapedTarget(path string) bool {
	if len(path) == 0 {
		return false
	}
	lower := strings.ToLower(path)
	if strings.Contains(lower, "://") || hasSchemePrefix(lower) {
		return true
	}

	return false
}

// hasIllegalPathChar flags Windows-illegal path chars after the drive
// letter (the first `:` at index 1 is legal; any later `:` is not, and
// `?`, `*`, `<`, `>`, `|`, `"` are never legal in a path component).
func hasIllegalPathChar(path string) bool {
	if len(path) == 0 {
		return false
	}
	rest := path
	if len(path) > 2 && path[1] == ':' {
		rest = path[2:]
	}
	if strings.ContainsAny(rest, `:?*<>|"`) {
		return true
	}

	return false
}

// isOrphanTarget returns true when TargetPath is a non-empty path that
// does not exist on disk. Empty paths are NOT treated as orphans —
// some task types legitimately leave TargetPath blank.
func isOrphanTarget(path string) bool {
	if len(path) == 0 {
		return false
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		// Can't even resolve — treat as orphan; user can drop it.
		return true
	}
	_, statErr := os.Stat(abs)

	return os.IsNotExist(statErr)
}

// confirmPendingClear prompts for "yes" on stdin. Anything else cancels.
func confirmPendingClear(count int) bool {
	fmt.Printf(constants.MsgPendingClearConfirm, count)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	return strings.EqualFold(strings.TrimSpace(answer), "yes")
}

// deleteSinglePendingClearCandidate deletes a single pending task and prints status.
func deleteSinglePendingClearCandidate(db *store.DB, c pendingClearCandidate) bool {
	if err := db.DeletePendingTask(c.task.ID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPendingClearDeleteFail, c.task.ID, err)

		return false
	}
	fmt.Printf(constants.MsgPendingClearDeleted, c.task.ID, c.task.TaskTypeName)

	return true
}

// deletePendingClearCandidates iterates the slice and deletes each
// row, printing a per-deletion line and tallying successes. Returns
// the number actually deleted (failures are logged but don't abort).
func deletePendingClearCandidates(db *store.DB,
	candidates []pendingClearCandidate) int {
	deleted := 0
	for _, c := range candidates {
		if deleteSinglePendingClearCandidate(db, c) {
			deleted++
		}
	}

	return deleted
}

// printPendingClearHeader prints the box banner + scan summary.
func printPendingClearHeader(mode string, scanned int) {
	fmt.Print(constants.MsgPendingClearHeader)
	fmt.Printf(constants.MsgPendingClearMode, mode)
	fmt.Printf(constants.MsgPendingClearScanned, scanned)
}

// printPendingClearCandidates prints one bullet per row that will be
// (or would be, in dry-run) deleted.
func printPendingClearCandidates(cands []pendingClearCandidate) {
	for _, c := range cands {
		fmt.Printf(constants.MsgPendingClearCandidate,
			c.task.ID, c.task.TaskTypeName, c.reason, c.task.TargetPath)
	}
}

// pendingClearCandidate pairs a row with its match-reason for output.
type pendingClearCandidate struct {
	task   model.PendingTaskRecord
	reason string
}
