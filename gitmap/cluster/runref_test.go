package cluster

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestRunRefGenerator(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open memory db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE ClusterRun (
		ClusterRunId INTEGER PRIMARY KEY AUTOINCREMENT,
		RunRef TEXT NOT NULL,
		CommandKind INTEGER NOT NULL,
		RawCommand TEXT NOT NULL,
		TargetSelector INTEGER NOT NULL,
		ExceptClause TEXT,
		StartedAt DATETIME NOT NULL,
		FinishedAt DATETIME,
		TotalNodes INTEGER NOT NULL DEFAULT 0,
		SucceededNodes INTEGER NOT NULL DEFAULT 0,
		FailedNodes INTEGER NOT NULL DEFAULT 0,
		SkippedNodes INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		if strings.Contains(err.Error(), "cgo") || strings.Contains(err.Error(), "CGO_ENABLED") {
			t.Skip("Skipping test due to CGO requirement for sqlite3")
		}
		t.Fatalf("Failed to create table: %v", err)
	}

	ref1, err := RunRefGenerator(db)
	if err != nil {
		t.Fatalf("RunRefGenerator error: %v", err)
	}

	dateStr := time.Now().UTC().Format("20060102")
	expected1 := "RUN-" + dateStr + "-001"
	if ref1 != expected1 {
		t.Errorf("Expected %s, got %s", expected1, ref1)
	}

	// Insert one to increment counter
	now := time.Now().UTC()
	_, err = db.Exec(`INSERT INTO ClusterRun (RunRef, CommandKind, RawCommand, TargetSelector, StartedAt) VALUES (?, ?, ?, ?, ?)`, ref1, 1, "cmd", 1, now.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to insert run: %v", err)
	}

	ref2, err := RunRefGenerator(db)
	if err != nil {
		t.Fatalf("RunRefGenerator error: %v", err)
	}

	expected2 := "RUN-" + dateStr + "-002"
	if ref2 != expected2 {
		t.Errorf("Expected %s, got %s", expected2, ref2)
	}
}
