package llm

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

const sqlMockHeatmapTable = `CREATE TABLE HeatmapActionLog (
	HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
	TaskSlug TEXT, RepoName TEXT, ActionType TEXT, FilePath TEXT,
	ExpectedOutcome TEXT, ActualOutcome TEXT, IsExpectationMet INTEGER,
	ReasoningState TEXT, ModelName TEXT, ToolName TEXT
);`

func TestResolveHeatmapTrainingSources_ExplicitURL(t *testing.T) {
	testUrl := "https://gitmap.internal/api/heatmap/task-314"
	sources, err := ResolveHeatmapTrainingSources(LlmHeatmapOptions{Source: testUrl})
	if err != nil {
		t.Fatalf("unexpected error resolving URL source: %v", err)
	}

	assertSingleSource(t, sources, testUrl)
}

func TestResolveHeatmapTrainingSources_ExplicitFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample_heatmap.db")
	createEmptyFile(t, filePath)

	sources, err := ResolveHeatmapTrainingSources(LlmHeatmapOptions{Source: filePath})
	if err != nil {
		t.Fatalf("unexpected error resolving file source: %v", err)
	}

	assertSingleSource(t, sources, filepath.ToSlash(filePath))
}

func TestResolveHeatmapTrainingSources_Ambient(t *testing.T) {
	sources, err := ResolveHeatmapTrainingSources(LlmHeatmapOptions{})
	if err != nil {
		t.Fatalf("unexpected error during ambient source resolution: %v", err)
	}

	hasNonNilSources := sources != nil
	if !hasNonNilSources && len(sources) != 0 {
		t.Fatalf("expected valid sources slice, got nil")
	}
}

func TestIngestHeatmapTrainingData_SQLite(t *testing.T) {
	dbPath := createMockHeatmapDb(t)
	examples, err := IngestHeatmapTrainingData(dbPath, 10)
	if err != nil {
		t.Fatalf("unexpected ingestion error: %v", err)
	}

	assertIngestedRecordCount(t, len(examples), 2)
	assertIngestedExampleDetails(t, examples[0], examples[1])
}

func TestIngestHeatmapTrainingData_RemoteURL(t *testing.T) {
	server := createMockHttpServer(t)
	defer server.Close()

	examples, err := IngestHeatmapTrainingData(server.URL, 5)
	if err != nil {
		t.Fatalf("unexpected remote ingestion error: %v", err)
	}

	assertIngestedRecordCount(t, len(examples), 3)
}

func TestPartitionHeatmapExamples(t *testing.T) {
	examples := buildSamplePartitionExamples()
	pos, neg := PartitionHeatmapExamples(examples)

	if len(pos) != 2 {
		t.Errorf("expected 2 positive examples, got %d", len(pos))
	}
	if len(neg) != 1 {
		t.Errorf("expected 1 negative example, got %d", len(neg))
	}
}

func TestFormatHeatmapTrainingDigest(t *testing.T) {
	examples := buildSamplePartitionExamples()
	digest := FormatHeatmapTrainingDigest(examples)

	assertDigestSections(t, digest)
	assertEmptyDigest(t)
}

func TestFormatHeatmapTrainingJsonl(t *testing.T) {
	examples := buildSamplePartitionExamples()
	jsonl, err := FormatHeatmapTrainingJsonl(examples)
	if err != nil {
		t.Fatalf("unexpected jsonl error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(jsonl), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 jsonl lines, got %d", len(lines))
	}
}

func TestRunTrainWithHeatmapFlag(t *testing.T) {
	dbPath := createMockHeatmapDb(t)
	args := []string{"--text-only", "--heatmap", dbPath}
	err := RunTrain(args)
	if err != nil {
		t.Fatalf("expected nil error for RunTrain with --heatmap, got: %v", err)
	}
}

func TestRunLlmTrainAlias(t *testing.T) {
	err := RunLlmTrain([]string{"--text-only"})
	if err != nil {
		t.Fatalf("expected nil error for RunLlmTrain, got: %v", err)
	}
}

func assertSingleSource(t *testing.T, sources []string, expected string) {
	if len(sources) != 1 {
		t.Fatalf("expected exactly 1 source, got %d", len(sources))
	}
	if sources[0] != expected {
		t.Errorf("expected source %q, got %q", expected, sources[0])
	}
}

func createEmptyFile(t *testing.T, path string) {
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
}

func createMockHeatmapDb(t *testing.T) string {
	dbPath := filepath.Join(t.TempDir(), "test_heatmap.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	defer conn.Close()

	initMockTable(t, conn)
	insertMockActionRows(t, conn)

	return dbPath
}

func initMockTable(t *testing.T, conn *sql.DB) {
	if _, err := conn.Exec(sqlMockHeatmapTable); err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
}

func insertMockActionRows(t *testing.T, conn *sql.DB) {
	query := `INSERT INTO HeatmapActionLog 
	(TaskSlug, RepoName, ActionType, FilePath, ExpectedOutcome, ActualOutcome, IsExpectationMet, ReasoningState, ModelName, ToolName)
	VALUES 
	('task-1', 'gitmap', 'REPLACE', 'cli/err.go', 'zero errors', 'compile fail', 0, 'REFACTOR', 'claude-3-5', 'replace'),
	('task-1', 'gitmap', 'MODIFY', 'cli/main.go', 'clean pass', 'clean pass', 1, 'VERIFICATION', 'claude-3-5', 'write_file');`
	if _, err := conn.Exec(query); err != nil {
		t.Fatalf("failed to insert action logs: %v", err)
	}
}

func createMockHttpServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := buildSamplePartitionExamples()
		bytes, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bytes)
	}))
}

func buildSamplePartitionExamples() []HeatmapTrainingExample {
	return []HeatmapTrainingExample{
		makeSampleExample("task-1", "cli/a.go", "pass", "pass", "EXECUTION", true),
		makeSampleExample("task-1", "cli/b.go", "pass", "pass", "EXECUTION", true),
		makeSampleExample("task-2", "cli/c.go", "success", "divergence", "RECOVERY", false),
	}
}

func makeSampleExample(task, file, exp, act, state string, isMet bool) HeatmapTrainingExample {
	return HeatmapTrainingExample{
		TaskSlug:         task,
		RepoSlug:         "gitmap",
		ActionType:       "MODIFY",
		FilePath:         file,
		ExpectedOutcome:  exp,
		ActualOutcome:    act,
		IsExpectationMet: isMet,
		ReasoningState:   state,
		ModelName:        "claude-3-5",
		ToolName:         "edit_file",
	}
}

func assertIngestedRecordCount(t *testing.T, actual, expected int) {
	if actual != expected {
		t.Fatalf("expected %d records, got %d", expected, actual)
	}
}

func assertIngestedExampleDetails(t *testing.T, first, second HeatmapTrainingExample) {
	if !first.IsExpectationMet {
		t.Errorf("expected first record (latest) to be positive, got negative")
	}
	if second.IsExpectationMet {
		t.Errorf("expected second record to be negative, got positive")
	}
}

func assertDigestSections(t *testing.T, digest string) {
	if !strings.Contains(digest, "HEATMAP COGNITIVE DECISION TRACES") {
		t.Errorf("expected header in digest")
	}
	if !strings.Contains(digest, "Positive Reinforcement Traces") {
		t.Errorf("expected positive reinforcement section in digest")
	}
	if !strings.Contains(digest, "Negative Anti-Pattern Traces") {
		t.Errorf("expected negative section in digest")
	}
}

func assertEmptyDigest(t *testing.T) {
	emptyDigest := FormatHeatmapTrainingDigest(nil)
	if !strings.Contains(emptyDigest, "No heatmap training actions recorded") {
		t.Errorf("expected empty message, got %q", emptyDigest)
	}
}
