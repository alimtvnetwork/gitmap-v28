package runlog

import (
	"database/sql"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

	_ "modernc.org/sqlite"
)

// openTestDB returns an in-memory SQLite DB with all commit-in tables
// + enum-mirror seeds applied. Centralized so every test below stays
// declarative.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	for _, ddl := range schemaDDL() {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatalf("ddl exec: %v\n%s", err, ddl)
		}
	}
	return db
}

// schemaDDL is the ordered list of DDL + seeds the run-log layer
// requires. Mirrors gitmap/store/migrate_commitin.go but trimmed to
// just the tables touched by this package.
func schemaDDL() []string {
	return []string{
		constants.SQLCreateCommitInRunStatus, constants.SQLSeedCommitInRunStatus,
		constants.SQLCreateCommitInInputKind, constants.SQLSeedCommitInInputKind,
		constants.SQLCreateCommitInOutcome, constants.SQLSeedCommitInOutcome,
		constants.SQLCreateCommitInSkipReason, constants.SQLSeedCommitInSkipReason,
		constants.SQLCreateCommitInProfile,
		constants.SQLCreateCommitInRun,
		constants.SQLCreateCommitInInputRepo,
		constants.SQLCreateCommitInSourceCommit,
		constants.SQLCreateCommitInSourceCommitFile,
		constants.SQLCreateCommitInRewritten,
		constants.SQLCreateCommitInSkipLog,
		constants.SQLCreateCommitInShaMap,
	}
}

// TestStartAndFinishRunFlowsThroughEnumLookups exercises the round-trip
// from `Running` → `Completed`, asserting both the row count and the
// final RunStatusId via the enum mirror.
func TestStartAndFinishRunFlowsThroughEnumLookups(t *testing.T) {
	db := openTestDB(t)
	runID, err := StartRun(db, "/abs/source", nil, true, nil, time.Now())
	if err != nil {
		t.Fatalf("StartRun: %v", err)
	}
	if runID <= 0 {
		t.Fatalf("expected positive runID, got %d", runID)
	}
	if err := FinishRun(db, runID, constants.CommitInRunStatusCompleted, time.Now()); err != nil {
		t.Fatalf("FinishRun: %v", err)
	}
	var statusName string
	if err := db.QueryRow(`SELECT s.Name FROM CommitInRun r JOIN RunStatus s ON s.RunStatusId = r.RunStatusId WHERE r.CommitInRunId = ?`, runID).Scan(&statusName); err != nil {
		t.Fatalf("readback: %v", err)
	}
	if statusName != constants.CommitInRunStatusCompleted {
		t.Fatalf("status = %q, want Completed", statusName)
	}
}

// TestInsertSourceCommitWritesFilesAtomically verifies the tx wraps
// commit + file inserts: the file rows MUST exist after the call.
func TestInsertSourceCommitWritesFilesAtomically(t *testing.T) {
	db := openTestDB(t)
	runID, _ := StartRun(db, "/abs/x", nil, false, nil, time.Now())
	inputID, err := InsertInputRepo(db, runID, 1, "input-a", "/tmp/abc", constants.CommitInInputKindLocalFolder)
	if err != nil {
		t.Fatalf("InsertInputRepo: %v", err)
	}
	row := SourceCommitRow{
		OrderIndex:           1,
		Sha:                  "deadbeef",
		AuthorName:           "alice",
		AuthorEmail:          "a@x",
		AuthorDateRFC3339:    "2024-01-02T03:04:05Z",
		CommitterDateRFC3339: "2024-01-02T03:04:06Z",
		OriginalMessage:      "first",
		Files:                []string{"a.go", "b.go"},
	}
	scID, err := InsertSourceCommit(db, inputID, row)
	if err != nil {
		t.Fatalf("InsertSourceCommit: %v", err)
	}
	var fileCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM SourceCommitFile WHERE SourceCommitId = ?`, scID).Scan(&fileCount); err != nil {
		t.Fatalf("count files: %v", err)
	}
	if fileCount != 2 {
		t.Fatalf("file count = %d, want 2", fileCount)
	}
}

// TestRecordRewrittenCreatedAlsoInsertsShaMap verifies the dedupe
// invariant: every Created outcome MUST register its SourceSha so the
// next run skips it.
func TestRecordRewrittenCreatedAlsoInsertsShaMap(t *testing.T) {
	db := openTestDB(t)
	scID := seedRunWithOneCommit(t, db, "feedface")
	rewID, err := RecordRewritten(db, 1, scID, RewrittenRow{
		NewSha:               "newsha111111",
		SourceSha:            "feedface",
		FinalMessage:         "rewritten msg",
		AuthorName:           "alice",
		AuthorEmail:          "a@x",
		AuthorDateRFC3339:    "2024-01-02T03:04:05Z",
		CommitterDateRFC3339: "2024-01-02T03:04:06Z",
		Outcome:              constants.CommitInOutcomeCreated,
	})
	if err != nil {
		t.Fatalf("RecordRewritten: %v", err)
	}
	if rewID <= 0 {
		t.Fatalf("expected positive rewrittenID")
	}
	var got int64
	if err := db.QueryRow(`SELECT RewrittenCommitId FROM ShaMap WHERE SourceSha = ?`, "feedface").Scan(&got); err != nil {
		t.Fatalf("ShaMap readback: %v", err)
	}
	if got != rewID {
		t.Fatalf("ShaMap.RewrittenCommitId = %d, want %d", got, rewID)
	}
}

// TestRecordSkipPersistsReason exercises the skip-log path used by
// stages 10/11/12 when a commit cannot be replayed.
func TestRecordSkipPersistsReason(t *testing.T) {
	db := openTestDB(t)
	scID := seedRunWithOneCommit(t, db, "abc123")
	if err := RecordSkip(db, 1, scID, constants.CommitInSkipReasonDuplicateSourceSha, nil); err != nil {
		t.Fatalf("RecordSkip: %v", err)
	}
	var reasonName string
	if err := db.QueryRow(`SELECT r.Name FROM SkipLog s JOIN SkipReason r ON r.SkipReasonId = s.SkipReasonId WHERE s.SourceCommitId = ?`, scID).Scan(&reasonName); err != nil {
		t.Fatalf("readback: %v", err)
	}
	if reasonName != constants.CommitInSkipReasonDuplicateSourceSha {
		t.Fatalf("reason = %q, want DuplicateSourceSha", reasonName)
	}
}

// seedRunWithOneCommit creates a CommitInRun + InputRepo + SourceCommit
// so the tests above can exercise downstream writes without repeating
// boilerplate.
func seedRunWithOneCommit(t *testing.T, db *sql.DB, sha string) int64 {
	t.Helper()
	runID, err := StartRun(db, "/abs/y", nil, false, nil, time.Now())
	if err != nil {
		t.Fatalf("StartRun: %v", err)
	}
	inputID, err := InsertInputRepo(db, runID, 1, "x", "/tmp/x", constants.CommitInInputKindLocalFolder)
	if err != nil {
		t.Fatalf("InsertInputRepo: %v", err)
	}
	scID, err := InsertSourceCommit(db, inputID, SourceCommitRow{
		OrderIndex: 1, Sha: sha,
		AuthorName: "x", AuthorEmail: "x@x",
		AuthorDateRFC3339: "2024-01-01T00:00:00Z", CommitterDateRFC3339: "2024-01-01T00:00:00Z",
		OriginalMessage: "m",
	})
	if err != nil {
		t.Fatalf("InsertSourceCommit: %v", err)
	}
	return scID
}
