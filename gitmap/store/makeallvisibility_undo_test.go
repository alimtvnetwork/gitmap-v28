// Package store — makeallvisibility_undo_test.go: round-trip tests
// for the 4 SELECT helpers in makeallvisibility_undo.go. Each test
// seeds a fresh on-disk SQLite DB via the canonical Insert/Update
// path, then asserts that the SELECT projects every column into the
// right struct field (guards against a future SQL column-order swap
// silently corrupting undo decisions).
package store

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// freshUndoDB returns a migrated DB at a per-test temp path.
func freshUndoDB(t *testing.T) *DB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "undo.sqlite")
	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAt: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	return db
}

// seedRun inserts a Run row + one Ok result with pub→pri transition
// and finalizes counts so it qualifies as undoable (OkCount>0).
func seedRun(t *testing.T, db *DB, kind, owner string) int64 {
	t.Helper()
	runID, err := db.InsertMakeAllVisibilityRun(model.MakeAllVisibilityRunRecord{
		CommandKind: kind, TargetVisibility: constants.VisibilityPrivate,
		Provider: constants.ProviderGitHub, Owner: owner, TargetRaw: owner,
		PatternList: "*", MatchedCount: 1, StartedAt: "2026-06-06T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("InsertRun: %v", err)
	}
	ids, err := db.InsertMakeAllVisibilityPendingResults(runID,
		[]model.MakeAllVisibilityResultRecord{{RepoName: "repo1", MatchedPattern: "*", StartedAt: "x"}})
	if err != nil {
		t.Fatalf("InsertPending: %v", err)
	}
	if err := db.UpdateMakeAllVisibilityResult(model.MakeAllVisibilityResultRecord{
		ID: ids[0], Status: constants.ResultStatusOk,
		PrevVisibility: constants.VisibilityPublic, NewVisibility: constants.VisibilityPrivate,
		FinishedAt: "y", DurationMs: 1,
	}); err != nil {
		t.Fatalf("UpdateResult: %v", err)
	}
	if err := db.FinalizeMakeAllVisibilityRun(model.MakeAllVisibilityRunRecord{
		ID: runID, OkCount: 1, ExitCode: 0, FinishedAt: "z",
	}); err != nil {
		t.Fatalf("Finalize: %v", err)
	}

	return runID
}

func TestSelectLatestUndoableMakeAllVisibilityRunReturnsNewest(t *testing.T) {
	db := freshUndoDB(t)
	_ = seedRun(t, db, constants.CommandKindMakeAllPrivate, "alpha")
	newest := seedRun(t, db, constants.CommandKindMakeAllPublic, "beta")
	got, err := db.SelectLatestUndoableMakeAllVisibilityRun()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got.ID != newest || got.Owner != "beta" || got.OkCount != 1 {
		t.Fatalf("expected newest run (id=%d owner=beta okCount=1), got %+v", newest, got)
	}
}

func TestSelectLatestUndoableEmptyReturnsZero(t *testing.T) {
	db := freshUndoDB(t)
	got, err := db.SelectLatestUndoableMakeAllVisibilityRun()
	if err != nil || got.ID != 0 {
		t.Fatalf("expected (zero, nil), got (%+v, %v)", got, err)
	}
}

func TestSelectMakeAllVisibilityRunByIDExactMatch(t *testing.T) {
	db := freshUndoDB(t)
	want := seedRun(t, db, constants.CommandKindMakeAllPrivate, "gamma")
	got, err := db.SelectMakeAllVisibilityRunByID(want)
	if err != nil || got.ID != want || got.Owner != "gamma" {
		t.Fatalf("by-id mismatch: %+v err=%v", got, err)
	}
}

func TestSelectMakeAllVisibilityRunByIDUnknownReturnsZero(t *testing.T) {
	db := freshUndoDB(t)
	got, err := db.SelectMakeAllVisibilityRunByID(9999)
	if err != nil || got.ID != 0 {
		t.Fatalf("expected (zero, nil) for unknown id, got (%+v, %v)", got, err)
	}
}

func TestSelectLatestMakeAllVisibilityRunByKindFilters(t *testing.T) {
	db := freshUndoDB(t)
	pub := seedRun(t, db, constants.CommandKindMakeAllPublic, "p1")
	_ = seedRun(t, db, constants.CommandKindMakeAllPrivate, "p2")
	got, err := db.SelectLatestMakeAllVisibilityRunByKind(constants.CommandKindMakeAllPublic)
	if err != nil || got.ID != pub {
		t.Fatalf("kind filter wrong: id=%d (want %d) err=%v", got.ID, pub, err)
	}
}

func TestSelectUndoableResultsForRunReturnsOkRows(t *testing.T) {
	db := freshUndoDB(t)
	runID := seedRun(t, db, constants.CommandKindMakeAllPrivate, "delta")
	rows, err := db.SelectUndoableResultsForRun(runID)
	if err != nil || len(rows) != 1 {
		t.Fatalf("expected 1 undoable row, got %d err=%v", len(rows), err)
	}
	r := rows[0]
	if r.RepoName != "repo1" || r.PrevVisibility != constants.VisibilityPublic ||
		r.NewVisibility != constants.VisibilityPrivate {
		t.Fatalf("column projection wrong: %+v", r)
	}
}
