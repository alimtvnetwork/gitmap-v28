package heatmap

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func executeTestSchema(t *testing.T, conn *sql.DB) {
	schema := `CREATE TABLE IF NOT EXISTS HeatmapActionLog (
		HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
		TaskSlug TEXT NOT NULL,
		RepoName TEXT NOT NULL DEFAULT '',
		ActionType TEXT NOT NULL,
		ReasoningState TEXT NOT NULL DEFAULT 'EXPLORATION',
		IsExpectationMet INTEGER NOT NULL DEFAULT 1,
		CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := conn.Exec(schema); err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}
}

func createTestDbWithSchema(t *testing.T, dbPath string) {
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}

	defer conn.Close()

	executeTestSchema(t, conn)
}

func insertTestAction(t *testing.T, dbPath, task, repo, state string, isMet int) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	defer conn.Close()

	query := `INSERT INTO HeatmapActionLog (TaskSlug, RepoName, ActionType, ReasoningState, IsExpectationMet, CreatedAt)
	VALUES (?, ?, 'READ', ?, ?, '2026-10-23 12:00:00')`
	if _, err := conn.Exec(query, task, repo, state, isMet); err != nil {
		t.Fatalf("failed to insert action: %v", err)
	}
}

func TestBuildSlugToAliasIndex(t *testing.T) {
	items := []store.AliasWithRepo{
		{Alias: model.Alias{Alias: "gm", IsPrimary: true}, Slug: "gitmap-v28"},
		{Alias: model.Alias{Alias: "am", IsPrimary: false}, Slug: "antigravity-manager"},
	}

	index := buildSlugToAliasIndex(items)
	if index["gitmap-v28"] != "gm" {
		t.Errorf("expected gm, got %s", index["gitmap-v28"])
	}

	if index["antigravity-manager"] != "am" {
		t.Errorf("expected am, got %s", index["antigravity-manager"])
	}
}

func TestApplyRepoAlias(t *testing.T) {
	index := map[string]string{"gitmap": "gm"}
	db := DiscoveredHeatmapDB{RepoSlug: "gitmap"}
	updated := applyRepoAlias(db, index)
	if updated.RepoAlias != "gm" {
		t.Errorf("expected gm, got %s", updated.RepoAlias)
	}
}

func TestCalculateMetPercent(t *testing.T) {
	if pct := calculateMetPercent(3, 4); pct != 75.0 {
		t.Errorf("expected 75.0, got %f", pct)
	}

	if pct := calculateMetPercent(0, 0); pct != 0.0 {
		t.Errorf("expected 0.0, got %f", pct)
	}
}

func TestParseHeatmapTimestamp(t *testing.T) {
	fallback := time.Now()
	parsed := parseHeatmapTimestamp("2026-10-23 12:00:00", fallback)
	if parsed.Year() != 2026 || parsed.Month() != 10 || parsed.Day() != 23 {
		t.Errorf("unexpected parsed date: %v", parsed)
	}

	if bad := parseHeatmapTimestamp("invalid", fallback); !bad.Equal(fallback) {
		t.Errorf("expected fallback on bad timestamp")
	}
}

func TestQueryTasksFromDb(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	createTestDbWithSchema(t, dbPath)
	insertTestAction(t, dbPath, "task-314", "gitmap", "PASS", 1)
	insertTestAction(t, dbPath, "task-314", "gitmap", "PASS", 1)
	insertTestAction(t, dbPath, "task-313", "gitmap", "FAIL", 0)

	dbEntry := DiscoveredHeatmapDB{Path: dbPath, RepoSlug: "gitmap", RepoAlias: "gm"}
	tasks := queryTasksFromDb(dbEntry, 10)
	assertRecentTasksQueryResult(t, tasks)
}

func assertRecentTasksQueryResult(t *testing.T, tasks []RecentHeatmapTaskRecord) {
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].TaskSlug != "task-314" && tasks[1].TaskSlug != "task-314" {
		t.Errorf("expected task-314 in results")
	}
}

func TestRenderHeatmapOverviewText(t *testing.T) {
	var buf bytes.Buffer
	dbs := []DiscoveredHeatmapDB{
		{Scope: "local_root", RepoSlug: "root", SizeBytes: 1024, ActionCount: 5, TaskCount: 1, LastModified: time.Now()},
	}
	tasks := []RecentHeatmapTaskRecord{
		{TaskSlug: "task-01", RepoSlug: "gitmap", RepoAlias: "gm", TotalActions: 2, MetPercent: 100.0, LatestState: "PASS", LastActiveAt: time.Now()},
	}

	RenderHeatmapOverview(&buf, dbs, tasks, false)
	out := buf.String()
	assertOverviewContains(t, out, "GITMAP HEATMAP OVERVIEW")
	assertOverviewContains(t, out, "Discovered Heatmap Databases:")
	assertOverviewContains(t, out, "Recent Cognitive Heatmap Tasks:")
}

func assertOverviewContains(t *testing.T, out, want string) {
	if !strings.Contains(out, want) {
		t.Errorf("expected output to contain %q", want)
	}
}

func TestRenderOverviewJSON(t *testing.T) {
	var buf bytes.Buffer
	dbs := []DiscoveredHeatmapDB{{Scope: "local_root", RepoSlug: "root", SizeBytes: 1024, ActionCount: 2}}
	tasks := []RecentHeatmapTaskRecord{{TaskSlug: "task-01", RepoSlug: "gitmap", TotalActions: 2, MetPercent: 100.0}}

	if err := renderOverviewJSON(&buf, dbs, tasks, true); err != nil {
		t.Fatalf("unexpected json render error: %v", err)
	}

	assertRcOverviewJSON(t, buf.Bytes())
}

func assertRcOverviewJSON(t *testing.T, data []byte) {
	var payload HeatmapOverviewPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}

	if !payload.IsSuccess || !payload.IsRc {
		t.Errorf("expected isSuccess=true and isRc=true, got success=%v rc=%v", payload.IsSuccess, payload.IsRc)
	}
}

func TestDispatchHeatmapOverview_ContextJSON(t *testing.T) {
	var buf bytes.Buffer
	ctx := WithOverviewJSON(context.Background())
	err := DispatchHeatmapOverview(ctx, &buf, false)
	if err != nil {
		t.Fatalf("unexpected dispatch error: %v", err)
	}

	assertValidJSONOverviewPayload(t, buf.Bytes())
}

func assertValidJSONOverviewPayload(t *testing.T, data []byte) {
	var payload HeatmapOverviewPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("expected valid JSON overview payload, got error: %v", err)
	}

	if !payload.IsSuccess {
		t.Errorf("expected IsSuccess to be true")
	}
}
