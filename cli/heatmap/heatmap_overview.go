package heatmap

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	_ "modernc.org/sqlite"
)

type heatmapOverviewContextKey struct{}

// RecentHeatmapTaskRecord contains execution metrics for a recent cognitive task.
type RecentHeatmapTaskRecord struct {
	TaskSlug     string    `json:"taskSlug"`
	RepoSlug     string    `json:"repoSlug"`
	RepoAlias    string    `json:"repoAlias,omitempty"`
	TotalActions int64     `json:"totalActions"`
	MetPercent   float64   `json:"metPercent"`
	LatestState  string    `json:"latestState"`
	LastActiveAt time.Time `json:"lastActiveAt"`
}

// HeatmapOverviewPayload models the JSON envelope when --json is requested.
type HeatmapOverviewPayload struct {
	IsSuccess      bool                      `json:"isSuccess"`
	AmbientAI      string                    `json:"ambientAi"`
	AmbientRepo    string                    `json:"ambientRepo"`
	IsRc           bool                      `json:"isRc"`
	TotalDatabases int                       `json:"totalDatabases"`
	TotalActions   int64                     `json:"totalActions"`
	Databases      []DiscoveredHeatmapDB     `json:"databases"`
	RecentTasks    []RecentHeatmapTaskRecord `json:"recentTasks"`
	Suggestions    []string                  `json:"suggestions"`
}

// WithOverviewJSON returns a context configured to output JSON overview format.
func WithOverviewJSON(ctx context.Context) context.Context {
	return context.WithValue(ctx, heatmapOverviewContextKey{}, true)
}

func isOverviewJSONRequested(ctx context.Context) bool {
	if val, hasVal := ctx.Value(heatmapOverviewContextKey{}).(bool); hasVal && val {
		return true
	}

	for _, a := range os.Args {
		lower := strings.ToLower(strings.TrimSpace(a))
		if lower == "--json" || lower == "-j" || lower == "--format=json" {
			return true
		}
	}

	return false
}

func queryTasksFromDb(db DiscoveredHeatmapDB, perDbLimit int) []RecentHeatmapTaskRecord {
	dsn := fmt.Sprintf("file:%s?mode=ro&cache=shared", filepath.ToSlash(db.Path))
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil
	}

	defer conn.Close()

	return executeRecentTasksQuery(conn, db, perDbLimit)
}

func executeRecentTasksQuery(conn *sql.DB, db DiscoveredHeatmapDB, limit int) []RecentHeatmapTaskRecord {
	query := `SELECT TaskSlug, COALESCE(RepoName, ''), COUNT(*),
COALESCE(SUM(CASE WHEN IsExpectationMet = 1 THEN 1 ELSE 0 END), 0),
COALESCE(MAX(ReasoningState), 'EXPLORATION'), COALESCE(MAX(CreatedAt), '')
FROM HeatmapActionLog GROUP BY TaskSlug ORDER BY LastActive DESC LIMIT ?`
	rows, err := conn.Query(query, limit)
	if err != nil {
		return nil
	}

	defer rows.Close()

	return scanRecentTaskRows(rows, db)
}

func scanRecentTaskRows(rows *sql.Rows, db DiscoveredHeatmapDB) []RecentHeatmapTaskRecord {
	var records []RecentHeatmapTaskRecord
	for rows.Next() {
		var task, repo, state, createdStr string
		var total, met int64
		if err := rows.Scan(&task, &repo, &total, &met, &state, &createdStr); err != nil {
			continue
		}

		records = append(records, buildRecentTaskRecord(task, repo, state, createdStr, total, met, db))
	}

	return records
}

func calculateMetPercent(met, total int64) float64 {
	if total <= 0 {
		return 0.0
	}

	return (float64(met) / float64(total)) * 100.0
}

func buildRecentTaskRecord(task, repo, state, createdStr string, total, met int64, db DiscoveredHeatmapDB) RecentHeatmapTaskRecord {
	effectiveRepo := repo
	if effectiveRepo == "" {
		effectiveRepo = db.RepoSlug
	}

	return RecentHeatmapTaskRecord{
		TaskSlug:     task,
		RepoSlug:     effectiveRepo,
		RepoAlias:    db.RepoAlias,
		TotalActions: total,
		MetPercent:   calculateMetPercent(met, total),
		LatestState:  state,
		LastActiveAt: parseHeatmapTimestamp(createdStr, db.LastModified),
	}
}

func parseHeatmapTimestamp(raw string, fallback time.Time) time.Time {
	if raw == "" {
		return fallback
	}

	formats := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02T15:04:05"}
	for _, fmtStr := range formats {
		if t, err := time.Parse(fmtStr, raw); err == nil {
			return t
		}
	}

	return fallback
}

func resolveTaskLimit(limit int) int {
	if limit <= 0 {
		return 5
	}

	return limit
}

func sortTasksByLastActive(tasks []RecentHeatmapTaskRecord) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].LastActiveAt.After(tasks[j].LastActiveAt)
	})
}

func trimRecentTasks(tasks []RecentHeatmapTaskRecord, limit int) []RecentHeatmapTaskRecord {
	hasSeen := make(map[string]bool)
	var result []RecentHeatmapTaskRecord
	for _, t := range tasks {
		if hasSeen[t.TaskSlug] {
			continue
		}

		hasSeen[t.TaskSlug] = true
		result = append(result, t)
		if len(result) >= limit {
			return result
		}
	}

	return result
}

// QueryRecentHeatmapTasks aggregates recent task metrics across discovered databases.
func QueryRecentHeatmapTasks(limit int) ([]RecentHeatmapTaskRecord, error) {
	effectiveLimit := resolveTaskLimit(limit)
	dbs, err := DiscoverHeatmapDatabases()
	if err != nil {
		return nil, err
	}

	enriched := EnrichHeatmapDatabasesWithAliases(dbs)
	var allTasks []RecentHeatmapTaskRecord
	for _, db := range enriched {
		allTasks = append(allTasks, queryTasksFromDb(db, effectiveLimit)...)
	}

	sortTasksByLastActive(allTasks)

	return trimRecentTasks(allTasks, effectiveLimit), nil
}

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
		fmt.Fprintln(out, "Recent Cognitive Heatmap Tasks:\n  (no recent tasks recorded)\n")
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

// DispatchHeatmapOverview coordinates database discovery, task aggregation, and rendering.
func DispatchHeatmapOverview(ctx context.Context, out io.Writer, isRc bool) error {
	dbs, err := DiscoverHeatmapDatabases()
	if err != nil {
		return err
	}

	enriched := EnrichHeatmapDatabasesWithAliases(dbs)
	tasks, _ := QueryRecentHeatmapTasks(5)
	if isOverviewJSONRequested(ctx) {
		return renderOverviewJSON(out, enriched, tasks, isRc)
	}

	RenderHeatmapOverview(out, enriched, tasks, isRc)

	return nil
}
