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

// runPendingClear is wired in from runPending when args[0] == "clear".
// args is everything after the "clear" token.
func runPendingClear(args []string) {
	mode, dryRun, yes, idMatch, err := parsePendingClearArgs(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	db, dbErr := openDB()
	if dbErr != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingDBOpen, dbErr)
		os.Exit(1)
	}
	defer db.Close()

	tasks, listErr := db.ListPendingTasks()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskQuery, listErr)
		os.Exit(1)
	}

	candidates := selectClearCandidates(tasks, mode, idMatch)
	printPendingClearHeader(mode, len(tasks))
	if len(candidates) == 0 {
		fmt.Printf(constants.MsgPendingClearNoMatches, mode)

		return
	}
	printPendingClearCandidates(candidates)

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

// parsePendingClearArgs splits args into (mode, dryRun, yes, idMatch, err).
// Mode is one of: orphans (default), illegal, all, or "id" (with idMatch
// holding the parsed numeric ID).
func parsePendingClearArgs(args []string) (
	mode string, dryRun, yes bool, idMatch int64, err error) {
	mode = "orphans"
	for _, a := range args {
		switch a {
		case "--dry-run":
			dryRun = true
		case "--yes", "-y":
			yes = true
		case "orphans", "illegal", "all":
			mode = a
		default:
			id, parseErr := strconv.ParseInt(a, 10, 64)
			if parseErr != nil || id <= 0 {
				if strings.HasPrefix(a, "-") {
					err = fmt.Errorf(
						constants.ErrPendingClearUnknownMode, a)

					return
				}
				err = fmt.Errorf(constants.ErrPendingClearBadID, a)

				return
			}
			mode = "id"
			idMatch = id
		}
	}

	return
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

// classifyPendingClearTask decides whether one task matches the mode
// and returns the reason label for the preview output.
func classifyPendingClearTask(t model.PendingTaskRecord,
	mode string, idMatch int64) (string, bool) {
	switch mode {
	case "all":
		return constants.MsgPendingClearReasonAll, true
	case "id":
		if t.ID == idMatch {
			return constants.MsgPendingClearReasonByID, true
		}

		return "", false
	case "illegal":
		if isURLShapedTarget(t.TargetPath) {
			return constants.MsgPendingClearReasonURL, true
		}
		if hasIllegalPathChar(t.TargetPath) {
			return constants.MsgPendingClearReasonChar, true
		}

		return "", false
	default: // orphans
		if isOrphanTarget(t.TargetPath) {
			return constants.MsgPendingClearReasonOrph, true
		}

		return "", false
	}
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
	if strings.Contains(lower, "://") {
		return true
	}
	// "https:\github.com\..." — colon-immediately-after-scheme even
	// though Windows broke the slashes.
	for _, scheme := range []string{"http:", "https:", "ssh:", "git:"} {
		if strings.Contains(lower, scheme+`\`) ||
			strings.Contains(lower, scheme+"/") {
			return true
		}
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

// deletePendingClearCandidates iterates the slice and deletes each
// row, printing a per-deletion line and tallying successes. Returns
// the number actually deleted (failures are logged but don't abort).
func deletePendingClearCandidates(db *store.DB,
	candidates []pendingClearCandidate) int {
	deleted := 0
	for _, c := range candidates {
		if err := db.DeletePendingTask(c.task.ID); err != nil {
			fmt.Fprintf(os.Stderr,
				constants.ErrPendingClearDeleteFail, c.task.ID, err)

			continue
		}
		deleted++
		fmt.Printf(constants.MsgPendingClearDeleted,
			c.task.ID, c.task.TaskTypeName)
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
