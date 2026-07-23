package dedupe

import (
	"database/sql"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

	_ "modernc.org/sqlite"
)

// openTestDB spins up an in-memory SQLite with just the tables this
// package reads (ShaMap + minimal FK-target chain).
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	for _, ddl := range []string{
		constants.SQLCreateCommitInRunStatus, constants.SQLSeedCommitInRunStatus,
		constants.SQLCreateCommitInInputKind, constants.SQLSeedCommitInInputKind,
		constants.SQLCreateCommitInOutcome, constants.SQLSeedCommitInOutcome,
		constants.SQLCreateCommitInProfile,
		constants.SQLCreateCommitInRun,
		constants.SQLCreateCommitInInputRepo,
		constants.SQLCreateCommitInSourceCommit,
		constants.SQLCreateCommitInRewritten,
		constants.SQLCreateCommitInShaMap,
	} {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatalf("ddl: %v", err)
		}
	}
	return db
}

// TestLookupReturnsMissOnEmptyTable verifies the "proceed" signal: a
// brand-new repo has zero ShaMap rows, so every lookup MUST return
// (Verdict{}, nil) — never an error.
func TestLookupReturnsMissOnEmptyTable(t *testing.T) {
	db := openTestDB(t)
	v, err := Lookup(db, "deadbeef")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if v.IsHit {
		t.Fatalf("expected miss, got hit %+v", v)
	}
}

// TestLookupReturnsHitWithPreviousId proves the hit path returns the
// stored RewrittenCommitId so the caller can FK SkipLog correctly.
func TestLookupReturnsHitWithPreviousId(t *testing.T) {
	db := openTestDB(t)
	rewID := seedRewritten(t, db, "abc999")
	if _, err := db.Exec(`INSERT INTO ShaMap (SourceSha, RewrittenCommitId) VALUES (?, ?)`, "abc999", rewID); err != nil {
		t.Fatalf("seed ShaMap: %v", err)
	}
	v, err := Lookup(db, "abc999")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if !v.IsHit || v.PreviousRewrittenId != rewID {
		t.Fatalf("verdict = %+v, want hit with id %d", v, rewID)
	}
}

// TestLookupRejectsEmptyShaInput is a guardrail: empty input is a
// programming bug, not a miss — surface it loudly.
func TestLookupRejectsEmptyShaInput(t *testing.T) {
	db := openTestDB(t)
	if _, err := Lookup(db, ""); err == nil {
		t.Fatalf("expected error for empty sha")
	}
}

// seedRewritten inserts the FK-target row so ShaMap inserts don't fail
// on REFERENCES RewrittenCommit.
func seedRewritten(t *testing.T, db *sql.DB, sourceSha string) int64 {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO CommitInRun (SourceRepoPath, StartedAt, RunStatusId) VALUES ('p', '2024-01-01', 1)`); err != nil {
		t.Fatalf("seed run: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO InputRepo (CommitInRunId, OrderIndex, OriginalRef, ResolvedPath, InputKindId) VALUES (1, 1, 'x', '/x', 1)`); err != nil {
		t.Fatalf("seed input: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO SourceCommit (InputRepoId, SourceSha, AuthorName, AuthorEmail, AuthorDate, CommitterDate, OriginalMessage, OrderIndex) VALUES (1, ?, 'a', 'a@x', '2024-01-01', '2024-01-01', 'm', 1)`, sourceSha); err != nil {
		t.Fatalf("seed source: %v", err)
	}
	res, err := db.Exec(`INSERT INTO RewrittenCommit (CommitInRunId, SourceCommitId, NewSha, FinalMessage, AppliedAuthorName, AppliedAuthorEmail, AppliedAuthorDate, AppliedCommitterDate, CommitOutcomeId) VALUES (1, 1, 'newsha', 'm', 'a', 'a@x', '2024-01-01', '2024-01-01', 1)`)
	if err != nil {
		t.Fatalf("seed rewritten: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}
