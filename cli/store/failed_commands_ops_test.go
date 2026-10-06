package store

import (
	"path/filepath"
	"testing"
)

func TestFailedCommandTableAndOperations(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test-errors.db")
	db, err := OpenErrorsSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenErrorsSplitDBAt failed: %v", err)
	}
	defer db.Close()

	rec1 := FailedCommandRecord{
		Command:     "deploy-keyz",
		FullArgs:    "deploy-keyz --all",
		Domain:      "root",
		ErrorCode:   "E1001",
		Message:     "Unknown command: deploy-keyz",
		Suggestions: "deploy-keys, deploy-keys-all",
	}

	id1, rErr := db.RecordFailedCommand(rec1)
	if rErr != nil {
		t.Fatalf("RecordFailedCommand failed: %v", rErr)
	}
	if id1 <= 0 {
		t.Fatalf("expected positive ID, got %d", id1)
	}

	// Record same command again to verify HitCount increment
	id2, rErr2 := db.RecordFailedCommand(rec1)
	if rErr2 != nil {
		t.Fatalf("second RecordFailedCommand failed: %v", rErr2)
	}
	if id2 != id1 {
		t.Fatalf("expected same row ID %d on duplicate command, got %d", id1, id2)
	}

	// Record a second distinct failed command in ssh domain
	rec2 := FailedCommandRecord{
		Command:     "ssh add-keys",
		FullArgs:    "ssh add-keys",
		Domain:      "ssh",
		ErrorCode:   "E1001",
		Message:     "unknown SSH command",
		Suggestions: "gitmap ssh key add, gitmap ssh deploy-keys",
	}
	if _, rErr := db.RecordFailedCommand(rec2); rErr != nil {
		t.Fatalf("RecordFailedCommand rec2 failed: %v", rErr)
	}

	summary, sErr := db.GetFailedCommandSummary(10)
	if sErr != nil {
		t.Fatalf("GetFailedCommandSummary failed: %v", sErr)
	}
	if summary.TotalDistinctCommands != 2 {
		t.Fatalf("expected 2 distinct failed commands, got %d", summary.TotalDistinctCommands)
	}
	if summary.TotalFailedAttempts != 3 {
		t.Fatalf("expected 3 total failed attempts, got %d", summary.TotalFailedAttempts)
	}
	if len(summary.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(summary.Records))
	}
	if summary.Records[0].Command != "deploy-keyz" || summary.Records[0].HitCount != 2 {
		t.Fatalf("expected top record deploy-keyz with HitCount=2, got %+v", summary.Records[0])
	}

	if err := db.ClearFailedCommands(); err != nil {
		t.Fatalf("ClearFailedCommands failed: %v", err)
	}

	distinctAfter, totalAfter, cErr := db.CountFailedCommands()
	if cErr != nil {
		t.Fatalf("CountFailedCommands after clear failed: %v", cErr)
	}
	if distinctAfter != 0 || totalAfter != 0 {
		t.Fatalf("expected 0 counts after clear, got distinct=%d total=%d", distinctAfter, totalAfter)
	}
}
