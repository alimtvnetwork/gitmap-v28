package cmd

// Sort + render helpers for the regoldens diff summary. Split from
// regoldens_diff.go so each file stays well under the 200-line cap
// and each function under the 15-line cap.

import (
	"fmt"
	"os"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// sortGoldenDiffEntries orders entries deterministically: by status
// (deletions last, additions first) then by path. Stable output
// helps reviewers diff two regoldens runs.
func sortGoldenDiffEntries(entries []goldenDiffEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].status != entries[j].status {
			return goldenDiffStatusRank(entries[i].status) < goldenDiffStatusRank(entries[j].status)
		}
		return entries[i].path < entries[j].path
	})
}

// goldenDiffStatusRank assigns a sort weight per status letter so
// the summary reads added → modified → renamed → deleted.
func goldenDiffStatusRank(status string) int {
	switch status {
	case "A":
		return 0
	case "M":
		return 1
	case "R":
		return 2
	case "D":
		return 3
	}
	return 4
}

// printGoldenDiffEntries emits one line per file plus an aggregate
// totals line. Output shape depends on mode: "short" omits +/-
// counts and rename details; "full" includes both.
func printGoldenDiffEntries(entries []goldenDiffEntry, mode string) {
	totals := goldenDiffTotals{count: len(entries)}
	for _, e := range entries {
		printOneDiffEntry(e, mode)
		totals.accumulate(e)
	}
	fmt.Fprintf(os.Stdout, constants.MsgRegoldensDiffTotals,
		totals.count, totals.added, totals.modified, totals.renamed,
		totals.deleted, totals.linesAdded, totals.linesDeleted)
}

// printOneDiffEntry renders a single entry per the requested mode.
// The full mode appends a rename-source line for R entries.
func printOneDiffEntry(e goldenDiffEntry, mode string) {
	if mode == constants.RegoldensDiffModeShort {
		fmt.Fprintf(os.Stdout, constants.MsgRegoldensDiffLineShort, e.status, e.path)
		return
	}
	fmt.Fprintf(os.Stdout, constants.MsgRegoldensDiffLineFull,
		e.status, e.path, e.added, e.deleted)
	if e.status == "R" && e.renamedFrom != "" {
		fmt.Fprintf(os.Stdout, constants.MsgRegoldensDiffRenameFull, e.renamedFrom)
	}
}

// goldenDiffTotals aggregates per-status counts and total +/- lines.
type goldenDiffTotals struct {
	count        int
	added        int
	modified     int
	renamed      int
	deleted      int
	linesAdded   int
	linesDeleted int
}

// accumulate folds one entry into the running totals. Status counts
// are mutually exclusive per file; line counts always accrue.
func (t *goldenDiffTotals) accumulate(e goldenDiffEntry) {
	switch e.status {
	case "A":
		t.added++
	case "D":
		t.deleted++
	case "R":
		t.renamed++
	default:
		t.modified++
	}
	t.linesAdded += e.added
	t.linesDeleted += e.deleted
}
