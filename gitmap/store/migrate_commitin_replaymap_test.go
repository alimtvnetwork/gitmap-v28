package store

import (
	"sort"
	"strings"
	"testing"
)

// columnInfo is the row shape of `PRAGMA table_info(<name>)`. Mirrors
// SQLite's column_info schema verbatim so the tests can assert exact
// column types, NOT NULL bits, and default-value presence.
type columnInfo struct {
	Name    string
	Type    string
	NotNull int
	Default any
	IsPK    int
}

func tableInfo(t *testing.T, db *DB, table string) []columnInfo {
	t.Helper()
	rows, err := db.conn.Query(`PRAGMA table_info("` + table + `")`)
	if err != nil {
		t.Fatalf("PRAGMA table_info(%s): %v", table, err)
	}
	defer rows.Close()
	var out []columnInfo
	for rows.Next() {
		var (
			ci  columnInfo
			cid int
		)
		if err := rows.Scan(&cid, &ci.Name, &ci.Type, &ci.NotNull, &ci.Default, &ci.IsPK); err != nil {
			t.Fatalf("scan: %v", err)
		}
		out = append(out, ci)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// TestCommitInReplayMapHasAllSpecColumns locks the §9.4 column list
// against PRAGMA. Adding/removing a column from spec MUST be paired
// with a constants edit AND this slice in the same commit.
func TestCommitInReplayMapHasAllSpecColumns(t *testing.T) {
	db := openTempDB(t)
	want := map[string]struct {
		Type    string
		NotNull int
	}{
		"CommitInReplayMapId":   {"INTEGER", 0}, // PK; SQLite reports NotNull=0 for AI PKs
		"CommitInRunId":         {"INTEGER", 1},
		"RewrittenCommitId":     {"INTEGER", 1},
		"SourceTagName":         {"TEXT", 1},
		"SourceTagSha":          {"TEXT", 1},
		"SourceCommitSha":       {"TEXT", 1},
		"DestTagSha":            {"TEXT", 0},
		"DestCommitSha":         {"TEXT", 0},
		"MirroredReleaseBranch": {"TEXT", 0},
		"IsVersionTag":          {"INTEGER", 1},
		"TagReplayOutcomeId":    {"INTEGER", 1},
		"CreatedAt":             {"DATETIME", 1},
	}
	got := tableInfo(t, db, "CommitInReplayMap")
	if len(got) != len(want) {
		t.Fatalf("CommitInReplayMap col count = %d, want %d (got: %+v)", len(got), len(want), got)
	}
	for _, ci := range got {
		w, ok := want[ci.Name]
		if !ok {
			t.Errorf("unexpected column %q", ci.Name)
			continue
		}
		if !strings.EqualFold(ci.Type, w.Type) {
			t.Errorf("col %s: type=%q, want %q", ci.Name, ci.Type, w.Type)
		}
		if ci.NotNull != w.NotNull {
			t.Errorf("col %s: NotNull=%d, want %d", ci.Name, ci.NotNull, w.NotNull)
		}
	}
}

// TestCommitInReplayMapForeignKeysPointAtSpecTargets locks the §9.4
// FK targets via PRAGMA foreign_key_list. Drift here = silent FK loss.
func TestCommitInReplayMapForeignKeysPointAtSpecTargets(t *testing.T) {
	db := openTempDB(t)
	rows, err := db.conn.Query(`PRAGMA foreign_key_list("CommitInReplayMap")`)
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_list: %v", err)
	}
	defer rows.Close()
	got := map[string]string{} // from-col -> to-table
	for rows.Next() {
		var (
			id, seq                                      int
			table, fromCol, toCol, onUpdate, onDel, mtch string
		)
		if err := rows.Scan(&id, &seq, &table, &fromCol, &toCol, &onUpdate, &onDel, &mtch); err != nil {
			t.Fatalf("scan: %v", err)
		}
		got[fromCol] = table
	}
	want := map[string]string{
		"CommitInRunId":      "CommitInRun",
		"RewrittenCommitId":  "RewrittenCommit",
		"TagReplayOutcomeId": "TagReplayOutcome",
	}
	for col, table := range want {
		if got[col] != table {
			t.Errorf("FK %s -> %q, want %q", col, got[col], table)
		}
	}
}

// TestCommitInReplayMapIndexesAreCreated locks §9.4: three indexes
// (SourceTagName, DestCommitSha, MirroredReleaseBranch) MUST exist so
// the §9.5 idempotency lookup and dashboard queries stay O(log N).
func TestCommitInReplayMapIndexesAreCreated(t *testing.T) {
	db := openTempDB(t)
	want := []string{
		"IX_CommitInReplayMap_DestCommitSha",
		"IX_CommitInReplayMap_MirroredReleaseBranch",
		"IX_CommitInReplayMap_SourceTagName",
	}
	rows, err := db.conn.Query(
		`SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='CommitInReplayMap' ORDER BY name`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if strings.HasPrefix(n, "sqlite_autoindex") {
			continue // skip implicit UNIQUE indexes
		}
		got = append(got, n)
	}
	if !equalSorted(got, want) {
		t.Errorf("indexes:\n  got:  %v\n  want: %v", got, want)
	}
}

// TestCommitInReplayMapUniqueOnRunPlusTagName covers the §9.4 UNIQUE
// (CommitInRunId, SourceTagName) constraint. A second insert with the
// same pair MUST fail. This is the "no duplicate tag rows per run"
// guard — separate from the secondary UNIQUE involving RewrittenCommitId.
func TestCommitInReplayMapUniqueOnRunPlusTagName(t *testing.T) {
	db := openTempDB(t)
	seedDeps(t, db)
	insert := `INSERT INTO CommitInReplayMap
		(CommitInRunId, RewrittenCommitId, SourceTagName, SourceTagSha,
		 SourceCommitSha, IsVersionTag, TagReplayOutcomeId)
		VALUES (?, ?, ?, ?, ?, ?, (SELECT TagReplayOutcomeId FROM TagReplayOutcome WHERE Name='Created'))`
	if _, err := db.conn.Exec(insert, 1, 1, "v1.0.0", "tag-1", "src-1", 1); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if _, err := db.conn.Exec(insert, 1, 1, "v1.0.0", "tag-1", "src-1", 1); err == nil {
		t.Fatal("expected UNIQUE (CommitInRunId, SourceTagName) violation")
	}
}

// TestCommitInReplayMapTaggedVersusNonTaggedCommits covers the spec
// promise that NON-TAGGED commits never produce CommitInReplayMap
// rows (§9.2: "no tag, no row") while TAGGED commits do. We seed
// three commits — two with tags, one without — and assert the count
// + the IsVersionTag distribution.
func TestCommitInReplayMapTaggedVersusNonTaggedCommits(t *testing.T) {
	db := openTempDB(t)
	seedDeps(t, db)
	// Commit 1: annotated semver tag → expected row, IsVersionTag=1.
	mustInsertReplay(t, db, 1, 1, "v1.2.3", true)
	// Commit 2: annotated non-version tag → expected row, IsVersionTag=0.
	mustInsertReplay(t, db, 1, 2, "nightly", false)
	// Commit 3: NO tag → caller does not insert. Verify count == 2.
	var n int
	if err := db.conn.QueryRow(
		`SELECT COUNT(*) FROM CommitInReplayMap WHERE CommitInRunId=1`).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 rows for tagged commits, got %d", n)
	}
	// Distribution check.
	verCount := scanInt(t, db,
		`SELECT COUNT(*) FROM CommitInReplayMap WHERE IsVersionTag=1 AND CommitInRunId=1`)
	if verCount != 1 {
		t.Errorf("expected 1 version-tag row, got %d", verCount)
	}
}

// TestCommitInReplayMapForeignKeysAreEnforced verifies that PRAGMA
// foreign_keys=ON (set by openDBAt) actually blocks an insert with a
// dangling RewrittenCommitId.
func TestCommitInReplayMapForeignKeysAreEnforced(t *testing.T) {
	db := openTempDB(t)
	seedDeps(t, db)
	if _, err := db.conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable FK: %v", err)
	}
	_, err := db.conn.Exec(`INSERT INTO CommitInReplayMap
		(CommitInRunId, RewrittenCommitId, SourceTagName, SourceTagSha,
		 SourceCommitSha, IsVersionTag, TagReplayOutcomeId)
		VALUES (1, 99999, 'v9.9.9', 'tag-x', 'src-x', 1,
		 (SELECT TagReplayOutcomeId FROM TagReplayOutcome WHERE Name='Created'))`)
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "foreign key") {
		t.Errorf("expected FK violation, got %v", err)
	}
}
