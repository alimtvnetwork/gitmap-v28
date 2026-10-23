package heatmap

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
COALESCE(MAX(ReasoningState), 'EXPLORATION'), COALESCE(MAX(CreatedAt), '') AS LastActive
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
