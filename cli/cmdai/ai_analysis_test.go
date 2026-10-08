package cmdai

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestParseLineRange(t *testing.T) {
	tests := []struct {
		input     string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{"", 0, 0, false},
		{"   ", 0, 0, false},
		{"42", 42, 42, false},
		{"10-50", 10, 50, false},
		{"1-1", 1, 1, false},
		{"0", 0, 0, true},
		{"-5", 0, 0, true},
		{"50-10", 0, 0, true},
		{"abc", 0, 0, true},
		{"10-abc", 0, 0, true},
		{"1-2-3", 0, 0, true},
	}

	for _, tt := range tests {
		start, end, err := ParseLineRange(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseLineRange(%q) err = %v; wantErr = %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && (start != tt.wantStart || end != tt.wantEnd) {
			t.Errorf("ParseLineRange(%q) = (%d, %d); want (%d, %d)", tt.input, start, end, tt.wantStart, tt.wantEnd)
		}
	}
}

func TestAiAnalysisTaskLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_analysis.db")

	db, err := OpenAiAnalysisSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAiAnalysisSplitDBAt failed: %v", err)
	}
	defer db.Close()

	now := time.Now().UTC()
	task := AiAnalysisTask{
		TaskId:          "task-test-01",
		RepoPath:        "d:/work/test",
		TaskDescription: "Verify task lifecycle",
		ModelName:       "gemini-2.5-pro",
		Status:          AiTaskStatusInProgress,
		StartedAt:       now,
		CreatedAt:       now,
	}

	if createErr := CreateAiAnalysisTask(db, task); createErr != nil {
		t.Fatalf("CreateAiAnalysisTask failed: %v", createErr)
	}

	retrieved, getErr := GetAiAnalysisTask(db, "task-test-01")
	if getErr != nil {
		t.Fatalf("GetAiAnalysisTask failed: %v", getErr)
	}
	if retrieved.TaskId != "task-test-01" || retrieved.ModelName != "gemini-2.5-pro" {
		t.Errorf("Retrieved mismatch: %+v", retrieved)
	}

	line1 := AiAnalysisLine{
		TaskId:        "task-test-01",
		FilePath:      "cli/store/test.go",
		OperationType: AiOperationRead,
		StartLine:     1,
		EndLine:       25,
		Reasoning:     "Checking types",
		CreatedAt:     time.Now(),
	}
	if recErr := RecordAiAnalysisLine(db, line1); recErr != nil {
		t.Fatalf("RecordAiAnalysisLine 1 failed: %v", recErr)
	}

	line2 := AiAnalysisLine{
		TaskId:        "task-test-01",
		FilePath:      "cli/store/test.go",
		OperationType: AiOperationEdit,
		StartLine:     30,
		EndLine:       40,
		Reasoning:     "Modifying function",
		CreatedAt:     time.Now(),
	}
	if recErr := RecordAiAnalysisLine(db, line2); recErr != nil {
		t.Fatalf("RecordAiAnalysisLine 2 failed: %v", recErr)
	}

	taskAfter, _ := GetAiAnalysisTask(db, "task-test-01")
	if taskAfter.TotalFilesRead != 1 || taskAfter.TotalFilesModified != 1 {
		t.Errorf("File counters mismatch: read=%d, mod=%d", taskAfter.TotalFilesRead, taskAfter.TotalFilesModified)
	}

	lines, linesErr := GetAiAnalysisLinesForTask(db, "task-test-01")
	if linesErr != nil {
		t.Fatalf("GetAiAnalysisLinesForTask failed: %v", linesErr)
	}
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if finErr := UpdateAiAnalysisTaskStatus(db, "task-test-01", AiTaskStatusCompleted, "Finished successfully"); finErr != nil {
		t.Fatalf("UpdateAiAnalysisTaskStatus failed: %v", finErr)
	}

	taskFinal, _ := GetAiAnalysisTask(db, "task-test-01")
	if taskFinal.Status != AiTaskStatusCompleted || taskFinal.CompletedAt == nil {
		t.Errorf("Expected completed status and timestamp, got status=%s completedAt=%v", taskFinal.Status, taskFinal.CompletedAt)
	}
}

func TestAiAnalysisBatchLinesAndCascade(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "cascade.db")

	db, err := OpenAiAnalysisSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAiAnalysisSplitDBAt failed: %v", err)
	}
	defer db.Close()

	task := AiAnalysisTask{
		TaskId:          "task-cascade-01",
		TaskDescription: "Batch and cascade test",
		ModelName:       "gpt-4o",
		Status:          AiTaskStatusInProgress,
		StartedAt:       time.Now(),
		CreatedAt:       time.Now(),
	}
	_ = CreateAiAnalysisTask(db, task)

	batch := []AiAnalysisLine{
		{TaskId: "task-cascade-01", FilePath: "a.go", OperationType: AiOperationRead, StartLine: 1, EndLine: 10, Reasoning: "inspect a", CreatedAt: time.Now()},
		{TaskId: "task-cascade-01", FilePath: "b.go", OperationType: AiOperationEdit, StartLine: 20, EndLine: 30, Reasoning: "edit b", CreatedAt: time.Now()},
		{TaskId: "task-cascade-01", FilePath: "c.go", OperationType: AiOperationDelete, StartLine: 1, EndLine: 50, Reasoning: "del c", CreatedAt: time.Now()},
	}
	if bErr := BatchRecordAiAnalysisLines(db, batch); bErr != nil {
		t.Fatalf("BatchRecordAiAnalysisLines failed: %v", bErr)
	}

	tCheck, _ := GetAiAnalysisTask(db, "task-cascade-01")
	if tCheck.TotalFilesRead != 1 || tCheck.TotalFilesModified != 1 || tCheck.TotalFilesDeleted != 1 {
		t.Errorf("Batch counters mismatch: read=%d mod=%d del=%d", tCheck.TotalFilesRead, tCheck.TotalFilesModified, tCheck.TotalFilesDeleted)
	}

	// Verify cascade deletion
	delRes := store.ExecWrapper(db, "DELETE FROM AiAnalysisTask WHERE task_id = ?", "task-cascade-01")
	if delRes.IsFailure {
		t.Fatalf("Delete task failed: %v", delRes.Error)
	}

	linesAfter, _ := GetAiAnalysisLinesForTask(db, "task-cascade-01")
	if len(linesAfter) != 0 {
		t.Errorf("Expected 0 lines after cascade delete, got %d", len(linesAfter))
	}
}

func TestDispatchAiAnalysis(t *testing.T) {
	err := DispatchAi([]string{"analysis", "list", "--limit", "1"})
	if err != nil {
		t.Errorf("DispatchAi(analysis list) error: %v", err)
	}
}
