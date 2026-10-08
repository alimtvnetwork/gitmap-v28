package cmdai

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAiTrainFormatting(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "train_test.db")

	db, err := OpenAiAnalysisSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAiAnalysisSplitDBAt failed: %v", err)
	}
	defer db.Close()

	now := time.Now().UTC()
	task := AiAnalysisTask{
		TaskId:           "train-task-100",
		TaskDescription:  "Refactor boolean logic",
		ModelName:        "claude-3-7-sonnet",
		Status:           AiTaskStatusCompleted,
		ReasoningSummary: "Successfully inverted conditions",
		StartedAt:        now.Add(-10 * time.Minute),
		CreatedAt:        now.Add(-10 * time.Minute),
	}
	_ = CreateAiAnalysisTask(db, task)

	lines := []AiAnalysisLine{
		{TaskId: "train-task-100", FilePath: "store/db.go", OperationType: AiOperationRead, StartLine: 1, EndLine: 50, Reasoning: "Analyzed conditions", CreatedAt: now.Add(-8 * time.Minute)},
		{TaskId: "train-task-100", FilePath: "store/db.go", OperationType: AiOperationEdit, StartLine: 25, EndLine: 30, Reasoning: "Inverted flag check", ContentSnippet: "if isEnabled { ... }", CreatedAt: now.Add(-5 * time.Minute)},
	}
	_ = BatchRecordAiAnalysisLines(db, lines)

	chains, extractErr := ExtractDecisionHistory(db, "train-task-100", 1)
	if extractErr != nil {
		t.Fatalf("ExtractDecisionHistory failed: %v", extractErr)
	}
	if len(chains) != 1 {
		t.Fatalf("Expected 1 chain, got %d", len(chains))
	}

	// 1. Verify Markdown
	md := FormatAsMarkdown(chains)
	if !strings.Contains(md, "train-task-100") || !strings.Contains(md, "Inverted flag check") {
		t.Errorf("Markdown output missing expected content: %s", md)
	}

	// 2. Verify JSONL
	jsonlBytes, jsonlErr := FormatAsJsonl(chains)
	if jsonlErr != nil {
		t.Fatalf("FormatAsJsonl failed: %v", jsonlErr)
	}

	var rec LlmConversationRecord
	if err := json.Unmarshal(jsonlBytes[:len(jsonlBytes)-1], &rec); err != nil {
		t.Fatalf("Unmarshal JSONL failed: %v", err)
	}
	if len(rec.Messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(rec.Messages))
	}
	if !strings.Contains(rec.Messages[2].Content, "Inverted flag check") {
		t.Errorf("Assistant content missing reasoning: %s", rec.Messages[2].Content)
	}
}

func TestAiExportAndImport(t *testing.T) {
	tempDir := t.TempDir()
	dbPathA := filepath.Join(tempDir, "db_a.db")
	dbPathB := filepath.Join(tempDir, "db_b.db")

	dbA, err := OpenAiAnalysisSplitDBAt(dbPathA)
	if err != nil {
		t.Fatalf("OpenAiAnalysisSplitDBAt A failed: %v", err)
	}
	defer dbA.Close()

	task := AiAnalysisTask{
		TaskId:          "export-task-01",
		TaskDescription: "Export test task",
		ModelName:       "gemini-2.5-pro",
		Status:          AiTaskStatusCompleted,
		StartedAt:       time.Now(),
		CreatedAt:       time.Now(),
	}
	_ = CreateAiAnalysisTask(dbA, task)

	line := AiAnalysisLine{
		TaskId:        "export-task-01",
		FilePath:      "test.go",
		OperationType: AiOperationEdit,
		StartLine:     1,
		EndLine:       10,
		Reasoning:     "Initial write",
		CreatedAt:     time.Now(),
	}
	_ = RecordAiAnalysisLine(dbA, line)

	// Export to JSON
	jsonPath := filepath.Join(tempDir, "export.json")
	optsExport := AiExportOptions{
		Format:   ExportFormatJson,
		Output:   jsonPath,
		RepoRoot: tempDir,
	}
	if err := exportDatabaseOrBundle(dbA, optsExport); err != nil {
		t.Fatalf("exportDatabaseOrBundle failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(jsonPath); err != nil {
		t.Fatalf("Exported JSON does not exist: %v", err)
	}

	// Import into DB B
	dbB, errB := OpenAiAnalysisSplitDBAt(dbPathB)
	if errB != nil {
		t.Fatalf("OpenAiAnalysisSplitDBAt B failed: %v", errB)
	}
	defer dbB.Close()

	optsImport := AiImportOptions{
		InputPath: jsonPath,
		RepoRoot:  tempDir,
	}
	data, _ := os.ReadFile(jsonPath)
	var bundle AiExportBundle
	_ = json.Unmarshal(data, &bundle)
	if err := importBundleIntoDb(dbB, bundle); err != nil {
		t.Fatalf("importBundleIntoDb failed: %v", err)
	}

	importedTask, getErr := GetAiAnalysisTask(dbB, "export-task-01")
	if getErr != nil {
		t.Fatalf("GetAiAnalysisTask on DB B failed: %v", getErr)
	}
	if importedTask.TotalFilesModified != 1 {
		t.Errorf("Expected TotalFilesModified 1, got %d", importedTask.TotalFilesModified)
	}

	// Test ZIP export
	zipPath := filepath.Join(tempDir, "export.zip")
	if err := exportZipBundle(bundle, zipPath); err != nil {
		t.Fatalf("exportZipBundle failed: %v", err)
	}

	r, zipErr := zip.OpenReader(zipPath)
	if zipErr != nil {
		t.Fatalf("zip.OpenReader failed: %v", zipErr)
	}
	defer r.Close()

	if len(r.File) < 2 {
		t.Errorf("Expected at least 2 files in zip, got %d", len(r.File))
	}
	_ = optsImport
}

func exportDatabaseOrBundle(db *sql.DB, opts AiExportOptions) error {
	bundle, err := buildExportBundle(db, opts)
	if err != nil {
		return err
	}
	return exportJsonBundle(bundle, opts.Output)
}

func TestAiClearAndPrune(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "prune.db")

	db, err := OpenAiAnalysisSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAiAnalysisSplitDBAt failed: %v", err)
	}
	defer db.Close()

	// Parse duration tests
	dur1, err1 := ParseRetentionDuration("14d")
	if err1 != nil || dur1 != 14*24*time.Hour {
		t.Errorf("ParseRetentionDuration(14d) = %v, err=%v", dur1, err1)
	}
	dur2, err2 := ParseRetentionDuration("2w")
	if err2 != nil || dur2 != 14*24*time.Hour {
		t.Errorf("ParseRetentionDuration(2w) = %v, err=%v", dur2, err2)
	}
	dur3, err3 := ParseRetentionDuration("24h")
	if err3 != nil || dur3 != 24*time.Hour {
		t.Errorf("ParseRetentionDuration(24h) = %v, err=%v", dur3, err3)
	}

	// Seed tasks
	oldTime := time.Now().UTC().Add(-40 * 24 * time.Hour)
	taskOld := AiAnalysisTask{
		TaskId:          "task-old",
		TaskDescription: "Old task",
		StartedAt:       oldTime,
		CreatedAt:       oldTime,
	}
	_ = CreateAiAnalysisTask(db, taskOld)

	taskNew := AiAnalysisTask{
		TaskId:          "task-new",
		TaskDescription: "New task",
		StartedAt:       time.Now(),
		CreatedAt:       time.Now(),
	}
	_ = CreateAiAnalysisTask(db, taskNew)

	// Prune older than 30d
	pruneCount, pruneErr := executePruneOlderThan(db, 30*24*time.Hour)
	if pruneErr != nil {
		t.Fatalf("executePruneOlderThan failed: %v", pruneErr)
	}
	if pruneCount != 1 {
		t.Errorf("Expected 1 pruned task, got %d", pruneCount)
	}

	// Delete specific task
	delCount, delErr := executeDeleteTask(db, "task-new")
	if delErr != nil {
		t.Fatalf("executeDeleteTask failed: %v", delErr)
	}
	if delCount != 1 {
		t.Errorf("Expected 1 deleted task, got %d", delCount)
	}

	// Test purge all with force
	taskTemp := AiAnalysisTask{TaskId: "task-temp", StartedAt: time.Now(), CreatedAt: time.Now()}
	_ = CreateAiAnalysisTask(db, taskTemp)
	purgeCount, purgeErr := executePurgeAll(db, true)
	if purgeErr != nil {
		t.Fatalf("executePurgeAll failed: %v", purgeErr)
	}
	if purgeCount != 1 {
		t.Errorf("Expected 1 purged task, got %d", purgeCount)
	}
}

func TestDispatchAiTrainAndClear(t *testing.T) {
	errTrain := DispatchAi([]string{"train", "--limit", "1"})
	if errTrain != nil {
		t.Errorf("DispatchAi(train) error: %v", errTrain)
	}

	errClear := DispatchAi([]string{"clear", "--task", "non-existent-task"})
	if errClear != nil {
		t.Errorf("DispatchAi(clear) error: %v", errClear)
	}
}
