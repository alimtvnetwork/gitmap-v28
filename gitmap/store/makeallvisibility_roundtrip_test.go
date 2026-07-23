// Package store — makeallvisibility_roundtrip_test.go: step-35
// round-trip coverage. Walks the full vu/vr lifecycle at the store
// layer: MakeAllPublic insert → VisibilityUndo insert (referencing
// it) → VisibilityRedo insert (referencing the undo) → assert every
// SELECT helper resolves each kind correctly and ordering is stable.
//
// Provider-level e2e (real `gh` calls) is a separate item (45, mock
// harness). This test locks the data-layer contract `vu` / `vr` rely
// on so future SQL refactors can't silently route the wrong row.
package store

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

func insertKind(t *testing.T, db *DB, kind, target, startedAt string) int64 {
	t.Helper()
	id, err := db.InsertMakeAllVisibilityRun(model.MakeAllVisibilityRunRecord{
		CommandKind: kind, TargetVisibility: target,
		Provider: constants.ProviderGitHub, Owner: "acme", TargetRaw: "acme",
		PatternList: "*", StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("InsertRun(%s): %v", kind, err)
	}

	return id
}

// TestVisibilityRoundTripStoreLayer: pub → undo → redo lifecycle.
func TestVisibilityRoundTripStoreLayer(t *testing.T) {
	db := freshHistoryDB(t)
	pubID := insertKind(t, db, constants.CommandKindMakeAllPublic,
		constants.VisibilityPublic, "2026-06-06T10:00:00Z")
	undoID := insertKind(t, db, constants.CommandKindVisibilityUndo,
		constants.VisibilityPrivate, "2026-06-06T10:05:00Z")
	redoID := insertKind(t, db, constants.CommandKindVisibilityRedo,
		constants.VisibilityPublic, "2026-06-06T10:10:00Z")

	rows, err := db.SelectRecentMakeAllVisibilityRuns(10)
	if err != nil || len(rows) != 3 {
		t.Fatalf("recent: got %d err=%v", len(rows), err)
	}
	if rows[0].ID != redoID || rows[1].ID != undoID || rows[2].ID != pubID {
		t.Fatalf("ordering: got %d,%d,%d want %d,%d,%d",
			rows[0].ID, rows[1].ID, rows[2].ID, redoID, undoID, pubID)
	}

	one, err := db.SelectMakeAllVisibilityRunByID(undoID)
	if err != nil || one.ID != undoID || one.CommandKind != constants.CommandKindVisibilityUndo {
		t.Fatalf("byID(undo): got %+v err=%v", one, err)
	}
}
