package committransfer

import (
	"fmt"
	"io"
	"os"

	"github.com/pterm/pterm"
)

// PrintHeader prints a stylized boxed banner for command initialization.
func PrintHeader(title string) {
	pterm.DefaultHeader.WithFullWidth().Println(title)
}

// PrintBoxedBanner prints a stylized boxed banner with custom styling.
func PrintBoxedBanner(title string) {
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).Println(title)
}

// BadgeReplaying returns a stylized [REPLAYING] status badge.
func BadgeReplaying() string {
	return pterm.NewStyle(pterm.BgCyan, pterm.FgBlack, pterm.Bold).Sprint(" [REPLAYING] ")
}

// BadgePRCreated returns a stylized [PR CREATED] status badge.
func BadgePRCreated() string {
	return pterm.NewStyle(pterm.BgYellow, pterm.FgBlack, pterm.Bold).Sprint(" [PR CREATED] ")
}

// BadgePRMerged returns a stylized [PR MERGED] status badge.
func BadgePRMerged() string {
	return pterm.NewStyle(pterm.BgGreen, pterm.FgBlack, pterm.Bold).Sprint(" [PR MERGED] ")
}

// BadgeSnapshotSynced returns a stylized [SNAPSHOT SYNCED] status badge.
func BadgeSnapshotSynced() string {
	return pterm.NewStyle(pterm.BgMagenta, pterm.FgWhite, pterm.Bold).Sprint(" [SNAPSHOT SYNCED] ")
}

// FormatCommitCounter returns a colorized [current/total] commit counter.
func FormatCommitCounter(current, total int) string {
	return pterm.NewStyle(pterm.FgLightCyan, pterm.Bold).Sprintf("[%d/%d]", current, total)
}

// LogReplaying writes a commit replay log entry with badge and counter.
func LogReplaying(w io.Writer, prefix, sha, subject string, current, total int) {
	badge := BadgeReplaying()
	counter := FormatCommitCounter(current, total)
	fmt.Fprintf(w, "%s %s %s %s %s\n", prefix, badge, counter, pterm.Cyan(sha), pterm.White(subject))
}

// LogPRCreated writes a PR creation log entry with badge.
func LogPRCreated(w io.Writer, prefix string, prNum int, branch string) {
	badge := BadgePRCreated()
	fmt.Fprintf(w, "%s %s #%d (%s)\n", prefix, badge, prNum, pterm.Yellow(branch))
}

// LogPRMerged writes a PR merged log entry with badge.
func LogPRMerged(w io.Writer, prefix string, prNum int, branch string) {
	badge := BadgePRMerged()
	fmt.Fprintf(w, "%s %s #%d (%s)\n", prefix, badge, prNum, pterm.Green(branch))
}

// LogSnapshotSynced writes a snapshot synced log entry with badge.
func LogSnapshotSynced(w io.Writer, prefix, targetDir string) {
	badge := BadgeSnapshotSynced()
	fmt.Fprintf(w, "%s %s target snapshot synchronized (%s)\n", prefix, badge, targetDir)
}

func countReplayable(commits []SourceCommit) int {
	willReplay := 0
	for _, c := range commits {
		if c.SkipCause == "" {
			willReplay++
		}
	}

	return willReplay
}

func printPlanHeader(w io.Writer, prefix string, willReplay, total, mergeExcluded int, includeMerges bool) {
	fmt.Fprintf(w,
		"%s replaying %d commits onto target (source-considered=%d, merge-excluded=%d)\n",
		prefix, willReplay, total+mergeExcluded, mergeExcluded)
	if mergeExcluded > 0 && !includeMerges {
		fmt.Fprintf(w,
			"%s   note: %d merge commits excluded by --no-include-merges\n",
			prefix, mergeExcluded)
	}
}

// PrintPlan writes a human-readable preview of plan to w (spec §3).
// Returns the number of commits that will actually be replayed (handy
// for the y/N prompt threshold).
func PrintPlan(w io.Writer, plan ReplayPlan, prefix string) int {
	willReplay := countReplayable(plan.Commits)
	printPlanHeader(w, prefix, willReplay, len(plan.Commits), plan.MergeExcluded, plan.IncludeMerges)
	for i, c := range plan.Commits {
		printPlanLine(w, prefix, i+1, len(plan.Commits), c)
	}

	return willReplay
}

// printPlanLine renders one commit's preview row.
func printPlanLine(w io.Writer, prefix string, i, n int, c SourceCommit) {
	padPrefix := "  " + pterm.Magenta(prefix)
	idxStr := FormatCommitCounter(i, n)
	shaStr := pterm.Cyan(c.ShortSHA)
	if c.SkipCause != "" {
		fmt.Fprintf(w, "%s %s %s → -        skipped: %s\n",
			prefix, idxStr, c.ShortSHA, c.SkipCause)

		return
	}

	subject := firstLine(c.Cleaned)
	fmt.Fprintf(w, "%s %s %s  %s\n", padPrefix, idxStr, shaStr, pterm.White(subject))
}

// firstLine returns everything before the first newline.
func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}

	return s
}

// PrintSummary writes the final replay summary (spec §12).
func PrintSummary(w io.Writer, prefix string, res ReplayResult) {
	fmt.Fprintf(w, "%s done: replayed %d, skipped %d (drop=%d, already-replayed=%d, empty=%d)\n",
		prefix, res.Replayed,
		res.SkippedDrop+res.SkippedReplayed+res.SkippedEmpty,
		res.SkippedDrop, res.SkippedReplayed, res.SkippedEmpty)
	if res.Pushed && len(res.NewSHAs) > 0 {
		fmt.Fprintf(w, "%s pushed %d commits\n", prefix, len(res.NewSHAs))
	}
}

func computeReconciliationStatus(considered, accounted int) string {
	if considered != accounted {
		return "discrepancy"
	}

	return "ok"
}

func printReconciliationLine(w io.Writer, prefix string, considered, accounted, excluded int, res ReplayResult, mark string) {
	skipped := res.SkippedDrop + res.SkippedReplayed + res.SkippedEmpty
	fmt.Fprintf(w,
		"%s reconcile: source-considered=%d, replayed=%d, skipped=%d (drop=%d, already-replayed=%d, empty=%d), merge-excluded=%d → accounted=%d [%s]\n",
		prefix, considered, res.Replayed, skipped,
		res.SkippedDrop, res.SkippedReplayed, res.SkippedEmpty,
		excluded, accounted, mark)
}

// PrintReconciliation writes a count-parity line so users can reconcile
// what they saw in `git log` against what landed on the target. When
// source-considered != accounted, a `discrepancy` line is written to
// errW so CI scripts can detect drift.
//
// Issue: .ai-memory/memory/issues/2026-05-09-commit-transfer-count-mismatch.md
func PrintReconciliation(w, errW io.Writer, prefix string, plan ReplayPlan, res ReplayResult) {
	considered := len(plan.Commits) + plan.MergeExcluded
	accounted := res.Replayed + res.SkippedDrop + res.SkippedReplayed +
		res.SkippedEmpty + plan.MergeExcluded
	mark := computeReconciliationStatus(considered, accounted)
	printReconciliationLine(w, prefix, considered, accounted, plan.MergeExcluded, res, mark)
	if considered != accounted && errW != nil {
		fmt.Fprintf(errW,
			"%s reconcile DISCREPANCY: source-considered=%d but accounted=%d (delta=%d)\n",
			prefix, considered, accounted, considered-accounted)
	}
}

// Confirm reads "y" / "yes" from os.Stdin and returns true. Anything
// else (including EOF) returns false.
func Confirm(prefix string) bool {
	fmt.Fprintf(os.Stdout, "%s proceed? [y/N] ", prefix)
	var ans string
	if _, err := fmt.Scanln(&ans); err != nil {
		return false
	}

	switch ans {
	case "y", "Y", "yes", "YES", "Yes":
		return true
	}

	return false
}
