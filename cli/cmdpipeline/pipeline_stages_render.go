package cmdpipeline

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func renderStagesTerminal(s *pipelinedb.PipelineStageSummary) {
	printStagesHeader(s)
	printStagesTable(s.Jobs)
	printStagesAnalytics(s)
}

func printStagesHeader(s *pipelinedb.PipelineStageSummary) {
	fmt.Printf("\n%s● Pipeline Stage & Job Timings:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Repository:       %s%s%s\n", constants.ColorWhite, s.RepoSlug, constants.ColorReset)
	fmt.Printf("  • Run ID:           #%d (%s)\n", s.RunId, s.WorkflowName)
	fmt.Printf("  • Outcome:          %s / %s\n\n", s.Status, formatConclusionBadge(s.Conclusion))
}

func printStagesTable(jobs []pipelinedb.PipelineJobRecord) {
	fmt.Printf("  %-3s %-38s %-14s %-10s %-20s\n", "#", "STAGE / JOB NAME", "STATUS", "DURATION", "STARTED AT")
	fmt.Printf("  %s\n", strings.Repeat("─", 88))
	for idx, j := range jobs {
		renderStageRow(idx+1, j)
	}
	fmt.Printf("  %s\n", strings.Repeat("─", 88))
}

func renderStageRow(seq int, j pipelinedb.PipelineJobRecord) {
	durStr := formatStageDuration(j.DurationSeconds)
	badge := formatConclusionBadge(j.Conclusion)
	name := truncateStageName(j.JobName, 36)
	startTime := formatIsoTimestampShort(j.StartedAt)
	fmt.Printf("  %-3d %-38s %-14s %-10s %-20s\n", seq, name, badge, durStr, startTime)
}

func printStagesAnalytics(s *pipelinedb.PipelineStageSummary) {
	fmt.Printf("\n%s● Stage Timing Analytics & Full Combined Approx:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Combined Stage Sum:   %s%s%s (%ds across %d stage(s))\n",
		constants.ColorWhite, formatStageDuration(s.StageSumSeconds), constants.ColorReset, s.StageSumSeconds, len(s.Jobs))
	fmt.Printf("  • Wall-Clock Duration:  %s%s%s (%ds actual runtime)\n",
		constants.ColorWhite, formatStageDuration(s.WallClockSeconds), constants.ColorReset, s.WallClockSeconds)
	fmt.Printf("  • Concurrency Speedup:  %s%.2fx%s parallel efficiency\n\n",
		constants.ColorGreen, s.SpeedupMultiplier, constants.ColorReset)
}

func formatStageDuration(secs int) string {
	hasSubMinute := secs < 60
	if hasSubMinute {
		return fmt.Sprintf("%ds", secs)
	}
	mins := secs / 60
	remSecs := secs % 60
	return fmt.Sprintf("%dm %ds", mins, remSecs)
}

func formatConclusionBadge(conclusion string) string {
	if conclusion == "success" {
		return constants.ColorGreen + "✔ success" + constants.ColorReset
	}
	if conclusion == "failure" {
		return constants.ColorRed + "✖ failure" + constants.ColorReset
	}
	if strings.HasPrefix(conclusion, "cancel") {
		return constants.ColorDim + "⊘ canceled" + constants.ColorReset
	}
	if conclusion == "skipped" {
		return constants.ColorDim + "○ skipped" + constants.ColorReset
	}
	return conclusion
}

func truncateStageName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen-3] + "..."
}

func formatIsoTimestampShort(ts string) string {
	if len(ts) < 19 {
		return ts
	}
	return ts[:10] + " " + ts[11:19]
}
