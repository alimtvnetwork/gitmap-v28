package runlog

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

	_ "modernc.org/sqlite"
)

// openTagReplayDB extends openTestDB with the migration-007 tables so
// the tag-replay helpers have everything they reference. Kept local
// to this file to avoid mutating the shared schemaDDL() slice (other
// tests deliberately exercise the v1 surface only).
func openTagReplayDB(t *testing.T) *sql.DB {
	t.Helper()
	db := openTestDB(t)
	for _, ddl := range []string{
		constants.SQLCreateCommitInTagOutcome,
		constants.SQLCreateCommitInReplayMap,
		constants.SQLCreateCommitInReplayMapTagNameIdx,
		constants.SQLCreateCommitInReplayMapDestShaIdx,
		constants.SQLCreateCommitInReplayMapBranchIdx,
		constants.SQLSeedCommitInTagOutcome,
	} {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatalf("tag-replay ddl exec: %v\n%s", err, ddl)
		}
	}
	return db
}

// seedRewrittenRow returns (runID, rewrittenID) ready for use as the
// FK targets on a CommitInReplayMap insert.
func seedRewrittenRow(t *testing.T, db *sql.DB, sourceSha, newSha string) (int64, int64) {
	t.Helper()
	runID, err := StartRun(db, "/abs/source", nil, false, nil, time.Now())
	if err != nil {
		t.Fatalf("StartRun: %v", err)
	}
	inputID, err := InsertInputRepo(db, runID, 1, "input-a", "/tmp/abc",
		constants.CommitInInputKindLocalFolder)
	if err != nil {
		t.Fatalf("InsertInputRepo: %v", err)
	}
	srcID, err := InsertSourceCommit(db, inputID, SourceCommitRow{
		OrderIndex: 1, Sha: sourceSha,
		AuthorName: "A", AuthorEmail: "a@x",
		AuthorDateRFC3339:    time.Now().Format(time.RFC3339),
		CommitterDateRFC3339: time.Now().Format(time.RFC3339),
		OriginalMessage:      "msg",
	})
	if err != nil {
		t.Fatalf("InsertSourceCommit: %v", err)
	}
	rewID, err := RecordRewritten(db, runID, srcID, RewrittenRow{
		NewSha: newSha, SourceSha: sourceSha, FinalMessage: "msg",
		AuthorName: "A", AuthorEmail: "a@x",
		AuthorDateRFC3339:    time.Now().Format(time.RFC3339),
		CommitterDateRFC3339: time.Now().Format(time.RFC3339),
		Outcome:              constants.CommitInOutcomeCreated,
	})
	if err != nil {
		t.Fatalf("RecordRewritten: %v", err)
	}
	return runID, rewID
}

// TestRecordTagReplayCreatedWritesAllColumns is the happy-path R1
// case from spec §9.8: annotated semver tag, branch mirrored,
// outcome Created, both dest SHAs populated.
func TestRecordTagReplayCreatedWritesAllColumns(t *testing.T) {
	db := openTagReplayDB(t)
	runID, rewID := seedRewrittenRow(t, db, "src1", "new1")
	id, err := RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName:         "v1.2.3",
		SourceTagSha:          "tag-sha-1",
		SourceCommitSha:       "src1",
		DestTagSha:            "dest-tag-1",
		DestCommitSha:         "new1",
		MirroredReleaseBranch: "release/v1.2.3",
		IsAnnotated:           true,
		IsVersionTag:          true,
		Outcome:               constants.TagReplayOutcomeCreated,
	})
	if err != nil || id <= 0 {
		t.Fatalf("RecordTagReplay: id=%d err=%v", id, err)
	}
	assertReplayRow(t, db, "v1.2.3", "dest-tag-1", "new1", "release/v1.2.3", 1, "Created")
}

// TestRecordTagReplayDryRunWritesNullDestColumns covers R6: dry-run
// inserts the row with all dest-side fields NULL.
func TestRecordTagReplayDryRunWritesNullDestColumns(t *testing.T) {
	db := openTagReplayDB(t)
	runID, rewID := seedRewrittenRow(t, db, "src2", "new2")
	if _, err := RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName: "v1.0.0", SourceTagSha: "tag-sha-2",
		SourceCommitSha: "src2", IsAnnotated: true, IsVersionTag: true,
		Outcome: constants.TagReplayOutcomeCreatedDryRun,
	}); err != nil {
		t.Fatalf("RecordTagReplay: %v", err)
	}
	var dt, dc, mb sql.NullString
	if err := db.QueryRow(`SELECT DestTagSha, DestCommitSha, MirroredReleaseBranch
		FROM CommitInReplayMap WHERE SourceTagName = ?`, "v1.0.0").
		Scan(&dt, &dc, &mb); err != nil {
		t.Fatalf("readback: %v", err)
	}
	if dt.Valid || dc.Valid || mb.Valid {
		t.Fatalf("expected NULLs on dry-run, got %+v %+v %+v", dt, dc, mb)
	}
}

// TestLookupTagReplayHitsOnPriorCreated covers R8: a re-run finds the
// previous Created row and gets its dest SHAs back.
func TestLookupTagReplayHitsOnPriorCreated(t *testing.T) {
	db := openTagReplayDB(t)
	runID, rewID := seedRewrittenRow(t, db, "src3", "new3")
	_, _ = RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName: "v2.0.0", SourceTagSha: "tag-sha-3",
		SourceCommitSha: "src3", DestTagSha: "dest-3", DestCommitSha: "new3",
		IsAnnotated: true, IsVersionTag: true, Outcome: constants.TagReplayOutcomeCreated,
	})
	got, err := LookupTagReplay(db, "v2.0.0", "tag-sha-3")
	if err != nil {
		t.Fatalf("LookupTagReplay: %v", err)
	}
	if got.DestTagSha != "dest-3" || got.DestCommitSha != "new3" {
		t.Fatalf("lookup got %+v", got)
	}
}

// TestLookupTagReplayMissesOnFailedOutcome locks §9.5: rows with
// Failed / Skipped / CreatedDryRun outcomes do NOT count as hits.
func TestLookupTagReplayMissesOnFailedOutcome(t *testing.T) {
	db := openTagReplayDB(t)
	runID, rewID := seedRewrittenRow(t, db, "src4", "new4")
	_, _ = RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName: "v3.0.0", SourceTagSha: "tag-sha-4",
		SourceCommitSha: "src4", IsAnnotated: true, IsVersionTag: true,
		Outcome: constants.TagReplayOutcomeFailed,
	})
	if _, err := LookupTagReplay(db, "v3.0.0", "tag-sha-4"); !errors.Is(err, ErrTagReplayMiss) {
		t.Fatalf("expected ErrTagReplayMiss, got %v", err)
	}
}

// TestLookupTagReplayMissesOnUnknownTag is the cold-cache case.
func TestLookupTagReplayMissesOnUnknownTag(t *testing.T) {
	db := openTagReplayDB(t)
	if _, err := LookupTagReplay(db, "vNope", "tag-sha-x"); !errors.Is(err, ErrTagReplayMiss) {
		t.Fatalf("expected ErrTagReplayMiss, got %v", err)
	}
}

// TestRecordTagReplayUniqueConstraintBlocksDuplicateInRun locks the
// (CommitInRunId, SourceTagName) UNIQUE per spec §9.4.
func TestRecordTagReplayUniqueConstraintBlocksDuplicateInRun(t *testing.T) {
	db := openTagReplayDB(t)
	runID, rewID := seedRewrittenRow(t, db, "src5", "new5")
	if _, err := RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName: "v4.0.0", SourceTagSha: "t1", SourceCommitSha: "src5",
		IsAnnotated: true, IsVersionTag: true, Outcome: constants.TagReplayOutcomeCreated,
	}); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if _, err := RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName: "v4.0.0", SourceTagSha: "t1", SourceCommitSha: "src5",
		IsAnnotated: true, IsVersionTag: true, Outcome: constants.TagReplayOutcomeCreated,
	}); err == nil {
		t.Fatal("expected UNIQUE violation, got nil")
	}
}

// TestRecordTagReplayRejectsLightweightVersionTag locks the strict-
// semver gate at the persistence boundary: a caller asserting
// IsVersionTag=true on a lightweight tag (IsAnnotated=false) MUST be
// rejected with ErrLightweightVersionTag and NOT insert a row.
func TestRecordTagReplayRejectsLightweightVersionTag(t *testing.T) {
	db := openTagReplayDB(t)
	runID, rewID := seedRewrittenRow(t, db, "src6", "new6")
	_, err := RecordTagReplay(db, runID, rewID, TagReplayFacts{
		SourceTagName: "v5.0.0", SourceTagSha: "t-lw", SourceCommitSha: "src6",
		IsAnnotated: false, IsVersionTag: true,
		Outcome: constants.TagReplayOutcomeCreated,
	})
	if !errors.Is(err, ErrLightweightVersionTag) {
		t.Fatalf("expected ErrLightweightVersionTag, got %v", err)
	}
	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM CommitInReplayMap WHERE SourceTagName=?`, "v5.0.0",
	).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 0 {
		t.Fatalf("rejected insert leaked a row: count=%d", n)
	}
}

// TestClassifyVersionTagStrictMatrix locks the canonical strict-semver
// classifier: annotated AND name-matches → true; everything else false.
func TestClassifyVersionTagStrictMatrix(t *testing.T) {
	cases := []struct {
		name        string
		isAnnotated bool
		want        bool
	}{
		{"v1.2.3", true, true},
		{"v1.2.3", false, false}, // lightweight rejected even with valid name
		{"1.2.3", true, true},
		{"1.2.3", false, false},
		{"nightly", true, false}, // annotated but not semver
		{"nightly", false, false},
		{"v1.2", true, false}, // annotated but not full semver
		{"", true, false},
	}
	for _, tc := range cases {
		got := ClassifyVersionTag(tc.name, tc.isAnnotated)
		if got != tc.want {
			t.Errorf("ClassifyVersionTag(%q, annotated=%v) = %v, want %v",
				tc.name, tc.isAnnotated, got, tc.want)
		}
	}
}

// TestIsAnnotatedSemverVersionTagMatrix locks the version-tag detector
// against the spec §9.4 examples plus a few border cases.
func TestIsAnnotatedSemverVersionTagMatrix(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"v1.2.3", true},
		{"1.2.3", true},
		{"v0.0.0", true},
		{"v10.20.30", true},
		{"v1.0.0-rc.1", true},
		{"v2.0.0+build.7", true},
		{"v1.0.0-alpha+sha.abc", true},
		{"v1.2", false},
		{"v1.2.3.4", false},
		{"nightly", false},
		{"release-1.0", false},
		{"", false},
		{"v01.2.3", false}, // leading zero in MAJOR is invalid SemVer
	}
	for _, tc := range cases {
		if got := IsAnnotatedSemverVersionTag(tc.name); got != tc.want {
			t.Errorf("IsAnnotatedSemverVersionTag(%q) = %v, want %v",
				tc.name, got, tc.want)
		}
	}
}

// assertReplayRow is a compact readback helper kept under 15 lines.
func assertReplayRow(t *testing.T, db *sql.DB, name, dt, dc, mb string, isVer int, outcome string) {
	t.Helper()
	var (
		gotDT, gotDC, gotMB sql.NullString
		gotIsVer            int
		gotOutcome          string
	)
	q := `SELECT m.DestTagSha, m.DestCommitSha, m.MirroredReleaseBranch,
		m.IsVersionTag, o.Name
		FROM CommitInReplayMap m
		JOIN TagReplayOutcome  o ON o.TagReplayOutcomeId = m.TagReplayOutcomeId
		WHERE m.SourceTagName = ?`
	if err := db.QueryRow(q, name).Scan(&gotDT, &gotDC, &gotMB, &gotIsVer, &gotOutcome); err != nil {
		t.Fatalf("readback %s: %v", name, err)
	}
	if gotDT.String != dt || gotDC.String != dc || gotMB.String != mb ||
		gotIsVer != isVer || gotOutcome != outcome {
		t.Fatalf("row mismatch for %s:\n got DT=%q DC=%q MB=%q isVer=%d outcome=%s",
			name, gotDT.String, gotDC.String, gotMB.String, gotIsVer, gotOutcome)
	}
}
