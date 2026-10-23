package heatmap

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func generateHeatmapSuggestions(dbs []DiscoveredHeatmapDB, tasks []RecentHeatmapTaskRecord, isRc bool) []string {
	var suggestions []string
	if len(tasks) > 0 {
		topTask := tasks[0]
		suggestions = append(suggestions, buildInspectSuggestion(topTask, isRc))
		suggestions = append(suggestions, buildTreeSuggestion(topTask, isRc))
	}

	suggestions = append(suggestions, buildSyncSuggestion(tasks, isRc))
	suggestions = append(suggestions, buildQuerySuggestion(isRc))

	return suggestions
}

func buildInspectSuggestion(task RecentHeatmapTaskRecord, isRc bool) string {
	if isRc {
		return fmt.Sprintf("gitmap rc heatmap status --task %s --repo %s", task.TaskSlug, task.RepoSlug)
	}

	return fmt.Sprintf("gitmap heatmap status --task %s", task.TaskSlug)
}

func buildTreeSuggestion(task RecentHeatmapTaskRecord, isRc bool) string {
	if isRc {
		return fmt.Sprintf("gitmap rc heatmap tree --task %s --repo %s", task.TaskSlug, task.RepoSlug)
	}

	return fmt.Sprintf("gitmap heatmap tree --task %s", task.TaskSlug)
}

func buildSyncSuggestion(tasks []RecentHeatmapTaskRecord, isRc bool) string {
	repo := "gitmap"
	task := "sync-task"
	if len(tasks) > 0 {
		repo = tasks[0].RepoSlug
		task = tasks[0].TaskSlug
	}

	return fmt.Sprintf("gitmap rc heatmap sync --repo %s --task %s", repo, task)
}

func buildQuerySuggestion(isRc bool) string {
	if isRc {
		return `gitmap rc heatmap search "fix" --json`
	}

	return `gitmap heatmap search "fix" --json`
}

func formatHeatmapRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}

	diff := time.Since(t)
	if diff < time.Minute {
		return fmt.Sprintf("%ds ago", int(diff.Seconds()))
	}

	if diff < time.Hour {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	}

	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}

	return fmt.Sprintf("%dd ago", int(diff.Hours()/24.0))
}

func renderOverviewHeader(out io.Writer, isRc bool) {
	title := "GITMAP HEATMAP OVERVIEW & TELEMETRY DISCOVERY"
	if isRc {
		title = "GITMAP RC HEATMAP OVERVIEW & MULTI-AI DISCOVERY"
	}

	fmt.Fprintf(out, "%s\n", termout.RoleSGRBold(termout.RoleAccent)+title+termout.Reset())
	currentAI := store.DetectCurrentAIName()
	currentRepo := store.ResolveAmbientRepoSlug("")
	fmt.Fprintf(out, "Ambient Context: AI: %s | Repo: %s\n\n", currentAI, currentRepo)
}

func renderDiscoveredDatabasesTable(out io.Writer, dbs []DiscoveredHeatmapDB) {
	fmt.Fprintln(out, "Discovered Heatmap Databases:")
	fmt.Fprintf(out, "  %-16s %-8s %-18s %-9s %-7s %-10s %s\n", "REPO", "ALIAS", "SCOPE", "ACTIONS", "TASKS", "SIZE", "LAST UPDATED")
	fmt.Fprintf(out, "  %-16s %-8s %-18s %-9s %-7s %-10s %s\n", "----------------", "--------", "------------------", "---------", "-------", "----------", "-------------------")
	for _, db := range dbs {
		printDiscoveredDbRow(out, db)
	}

	fmt.Fprintln(out)
}

func printDiscoveredDbRow(out io.Writer, db DiscoveredHeatmapDB) {
	alias := "-"
	if db.RepoAlias != "" {
		alias = db.RepoAlias
	}

	timeStr := db.LastModified.Format("2006-01-02 15:04:05")
	sizeStr := formatHeatmapFileSize(db.SizeBytes)
	fmt.Fprintf(out, "  %-16s %-8s %-18s %-9d %-7d %-10s %s\n", db.RepoSlug, alias, db.Scope, db.ActionCount, db.TaskCount, sizeStr, timeStr)
}

func renderRecentTasksTable(out io.Writer, tasks []RecentHeatmapTaskRecord) {
	if len(tasks) == 0 {
		fmt.Fprintln(out, "Recent Cognitive Heatmap Tasks:")
		fmt.Fprintln(out, "  (no recent tasks recorded)")
		fmt.Fprintln(out)
		return
	}

	fmt.Fprintln(out, "Recent Cognitive Heatmap Tasks:")
	fmt.Fprintf(out, "  %-36s %-14s %-8s %-9s %-8s %-12s %s\n", "TASK SLUG", "REPO", "ALIAS", "ACTIONS", "MET %", "LAST STATE", "UPDATED")
	fmt.Fprintf(out, "  %-36s %-14s %-8s %-9s %-8s %-12s %s\n", "------------------------------------", "--------------", "--------", "---------", "--------", "------------", "-----------")
	for _, t := range tasks {
		printRecentTaskRow(out, t)
	}

	fmt.Fprintln(out)
}

func printRecentTaskRow(out io.Writer, t RecentHeatmapTaskRecord) {
	alias := "-"
	if t.RepoAlias != "" {
		alias = t.RepoAlias
	}

	ago := formatHeatmapRelativeTime(t.LastActiveAt)
	pctStr := fmt.Sprintf("%.1f%%", t.MetPercent)
	stateStyled := formatReasoningState(t.LatestState)
	fmt.Fprintf(out, "  %-36s %-14s %-8s %-9d %-8s %-12s %s\n", t.TaskSlug, t.RepoSlug, alias, t.TotalActions, pctStr, stateStyled, ago)
}

func renderDynamicSuggestions(out io.Writer, suggestions []string) {
	if len(suggestions) == 0 {
		return
	}

	fmt.Fprintln(out, "Dynamic Suggestions for AI & Developers:")
	for _, s := range suggestions {
		fmt.Fprintf(out, "  %s %s\n", termout.RoleSGR(termout.RoleAccent)+"▶"+termout.Reset(), s)
	}

	fmt.Fprintln(out)
}

// RenderHeatmapOverview outputs formatted terminal overview across databases and tasks.
func RenderHeatmapOverview(out io.Writer, dbs []DiscoveredHeatmapDB, tasks []RecentHeatmapTaskRecord, isRc bool) {
	renderOverviewHeader(out, isRc)
	renderDiscoveredDatabasesTable(out, dbs)
	renderRecentTasksTable(out, tasks)
	suggestions := generateHeatmapSuggestions(dbs, tasks, isRc)
	renderDynamicSuggestions(out, suggestions)
}

func sumTotalActions(dbs []DiscoveredHeatmapDB) int64 {
	var total int64
	for _, db := range dbs {
		total += db.ActionCount
	}

	return total
}

func buildOverviewPayload(dbs []DiscoveredHeatmapDB, tasks []RecentHeatmapTaskRecord, isRc bool) HeatmapOverviewPayload {
	return HeatmapOverviewPayload{
		IsSuccess:      true,
		AmbientAI:      store.DetectCurrentAIName(),
		AmbientRepo:    store.ResolveAmbientRepoSlug(""),
		IsRc:           isRc,
		TotalDatabases: len(dbs),
		TotalActions:   sumTotalActions(dbs),
		Databases:      dbs,
		RecentTasks:    tasks,
		Suggestions:    generateHeatmapSuggestions(dbs, tasks, isRc),
	}
}

func renderOverviewJSON(out io.Writer, dbs []DiscoveredHeatmapDB, tasks []RecentHeatmapTaskRecord, isRc bool) error {
	payload := buildOverviewPayload(dbs, tasks, isRc)
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.overview.json")
	}

	_, err = fmt.Fprintln(out, string(data))

	return err
}
