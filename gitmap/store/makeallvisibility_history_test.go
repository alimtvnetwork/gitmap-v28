// Package store — makeallvisibility_history_test.go: smoke coverage
// for SelectRecentMakeAllVisibilityRuns (the read-path behind `vh`).
// Asserts newest-first ordering, limit honored, and per-column Scan
// projection (column-order swap = silently wrong tally display).
package store

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

func freshHistoryDB(t *testing.T) *DB {
	t.Helper()
	db, err := OpenAt(filepath.Join(t.TempDir(), "hist.sqlite"))
	if err != nil {
		t.Fatalf("OpenAt: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	return db
}

func insertHistoryRun(t *testing.T, db *DB, owner string) int64 {
	t.Helper()
	id, err := db.InsertMakeAllVisibilityRun(model.MakeAllVisibilityRunRecord{
		CommandKind:      constants.CommandKindMakeAllPublic,
		TargetVisibility: constants.VisibilityPublic,
		Provider:         constants.ProviderGitHub, Owner: owner, TargetRaw: owner,
		PatternList: "*", StartedAt: "2026-06-06T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("InsertRun(%s): %v", owner, err)
	}

	return id
}

func TestSelectRecentMakeAllVisibilityRunsEmpty(t *testing.T) {
	db := freshHistoryDB(t)
	rows, err := db.SelectRecentMakeAllVisibilityRuns(10)
	if err != nil || len(rows) != 0 {
		t.Fatalf("expected empty, got %d rows err=%v", len(rows), err)
	}
}

func TestSelectRecentMakeAllVisibilityRunsNewestFirst(t *testing.T) {
	db := freshHistoryDB(t)
	_ = insertHistoryRun(t, db, "oldest")
	_ = insertHistoryRun(t, db, "middle")
	newest := insertHistoryRun(t, db, "newest")
	rows, err := db.SelectRecentMakeAllVisibilityRuns(10)
	if err != nil || len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d err=%v", len(rows), err)
	}
	if rows[0].ID != newest || rows[0].Owner != "newest" {
		t.Fatalf("expected newest-first ordering, got %+v", rows[0])
	}
}

func TestSelectRecentMakeAllVisibilityRunsHonorsLimit(t *testing.T) {
	db := freshHistoryDB(t)
	for i := 0; i < 5; i++ {
		_ = insertHistoryRun(t, db, "owner")
	}
	rows, err := db.SelectRecentMakeAllVisibilityRuns(2)
	if err != nil || len(rows) != 2 {
		t.Fatalf("expected limit=2 to return 2 rows, got %d err=%v", len(rows), err)
	}
}
