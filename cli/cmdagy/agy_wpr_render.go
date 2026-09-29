// Package cmdagy — agy_wpr_render.go formats terminal tables and UI status boxes for WPR.
package cmdagy

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RenderWPRAllSummaryTable displays the backup summary of all watched projects.
func RenderWPRAllSummaryTable(summaries []store.WatchPromptsSummary) {
	fmt.Println()
	fmt.Printf("  %s● Antigravity Watch Prompts Running Snapshot (%d project(s))%s\n",
		constants.ColorCyan, len(summaries), constants.ColorReset)
	fmt.Printf("  %-18s %-28s %-10s %-7s %-9s %s\n",
		"PROJECT", "PATH", "STATUS", "WORDS", "PICTURES", "WATCHING")
	fmt.Printf("  %s\n", strings.Repeat("─", 84))

	for _, s := range summaries {
		printSummaryRow(s)
	}

	fmt.Println()
}

func printSummaryRow(s store.WatchPromptsSummary) {
	statusColor := constants.ColorGreen
	if s.PromptStatus != "running" {
		statusColor = constants.ColorYellow
	}

	watchLabel := "YES"
	if !s.IsWatching {
		watchLabel = "NO"
	}

	proj := s.ProjectName
	if len(proj) > 18 {
		proj = proj[:15] + "..."
	}

	path := s.ProjectPath
	if len(path) > 28 {
		path = "..." + path[len(path)-25:]
	}

	fmt.Printf("  %-18s %-28s %s%-10s%s %-7d %-9d %s\n",
		proj, path, statusColor, strings.ToUpper(s.PromptStatus), constants.ColorReset,
		s.WordCount, s.MediaCount, watchLabel)
}

// RenderWPRProjectsTable displays the list of scheduled/monitored projects.
func RenderWPRProjectsTable(projects []WPRProjectStatus, st WPRRuntimeState, isRunning bool) {
	fmt.Println()
	fmt.Printf("  %s╔════ WATCH PROMPTS RUNNING LIST (%d projects) ════╗%s\n\n",
		constants.ColorCyan, len(projects), constants.ColorReset)

	printWatchStateHeader(st, isRunning)
	fmt.Printf("  %-18s %-30s %-10s %-7s %-9s %s\n",
		"PROJECT", "PATH", "STATUS", "WORDS", "PICTURES", "UPDATED")
	fmt.Printf("  %s\n", strings.Repeat("─", 88))

	for _, p := range projects {
		printProjectStatusRow(p)
	}

	fmt.Println()
}

func printWatchStateHeader(st WPRRuntimeState, isRunning bool) {
	if isRunning {
		fmt.Printf("  %s● Watch Loop State:%s ACTIVE (PID: %d, interval: %ds, started: %s)\n\n",
			constants.ColorGreen, constants.ColorReset, st.PID, st.IntervalSec, st.StartedAt)
		return
	}

	fmt.Printf("  %s○ Watch Loop State:%s IDLE (use 'gitmap wpr start' to begin)\n\n",
		constants.ColorYellow, constants.ColorReset)
}

func printProjectStatusRow(p WPRProjectStatus) {
	statusColor := constants.ColorGreen
	if p.PromptStatus != "running" {
		statusColor = constants.ColorYellow
	}

	proj := p.ProjectName
	if len(proj) > 18 {
		proj = proj[:15] + "..."
	}

	path := p.ProjectPath
	if len(path) > 30 {
		path = "..." + path[len(path)-27:]
	}

	up := p.UpdatedAt
	if len(up) > 19 {
		up = up[:19]
	}

	fmt.Printf("  %-18s %-30s %s%-10s%s %-7d %-9d %s\n",
		proj, path, statusColor, strings.ToUpper(p.PromptStatus), constants.ColorReset,
		p.WordCount, p.MediaCount, up)
}

// RenderWPRStatusBox renders comprehensive status including other machines' cached identities.
func RenderWPRStatusBox(st WPRRuntimeState, isRunning bool, machines []WPRMachineCacheEntry) {
	fmt.Printf("\n  %s● Watch Prompts Running Status%s\n", constants.ColorCyan, constants.ColorReset)
	printWatchStateHeader(st, isRunning)

	fmt.Printf("  %s● Other Fleet Machines (Cached from GitMap Home)%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %-12s %-20s %-16s %-10s %s\n", "ALIAS", "MACHINE NAME", "IP ADDRESS", "STATUS", "CACHED AT")
	fmt.Printf("  %s\n", strings.Repeat("─", 74))

	for _, m := range machines {
		up := m.CachedAt
		if len(up) > 19 {
			up = up[:19]
		}
		fmt.Printf("  %-12s %-20s %-16s %-10s %s\n",
			m.Alias, m.MachineName, m.IPAddress, m.Status, up)
	}

	fmt.Println()
}

// RenderWPRLogsTable displays watch logs for a project.
func RenderWPRLogsTable(slug string, logs []store.WatchLogRecord) {
	fmt.Printf("\n  %s● Watch Prompts Logs for [%s] (%d records)%s\n",
		constants.ColorCyan, slug, len(logs), constants.ColorReset)
	fmt.Printf("  %-19s %-16s %-10s %s\n", "TIME", "EVENT", "STATUS", "MESSAGE")
	fmt.Printf("  %s\n", strings.Repeat("─", 70))

	for _, l := range logs {
		timeStr := l.CreatedAt
		if len(timeStr) > 19 {
			timeStr = timeStr[:19]
		}
		fmt.Printf("  %-19s %-16s %-10s %s\n", timeStr, l.Event, l.Status, l.Message)
	}

	fmt.Println()
}

// RenderWPRDeploySummary prints the multi-node deploy results table.
func RenderWPRDeploySummary(summary WPRDeploySummary) {
	fmt.Println()
	fmt.Printf("  %s🚀 GitMap WPR Fleet Deployment Summary (%d/%d nodes successful)%s\n",
		constants.ColorCyan, summary.SuccessCount, summary.TotalNodes, constants.ColorReset)
	fmt.Printf("  %-14s %-18s %-16s %-8s %-10s %s\n",
		"NODE ALIAS", "MACHINE NAME", "IP ADDRESS", "COPIED", "ENQUEUED", "DURATION")
	fmt.Printf("  %s\n", strings.Repeat("─", 78))

	for _, r := range summary.Results {
		copiedMark := "YES"
		if !r.IsCopied {
			copiedMark = "NO"
		}
		enqMark := "YES"
		if !r.IsTaskEnqueued {
			enqMark = "NO"
		}

		fmt.Printf("  %-14s %-18s %-16s %-8s %-10s %dms\n",
			r.NodeAlias, r.MachineName, r.IPAddress, copiedMark, enqMark, r.DurationMs)
	}

	fmt.Println()
}
