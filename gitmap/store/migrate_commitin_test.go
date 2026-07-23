package store

import (
	"path/filepath"
	"sort"
	"testing"
)

// openTempDB opens a brand-new SQLite DB under t.TempDir() and runs
// the full Migrate() pipeline, returning the *DB for table inspection.
// Centralized here so future commit-in tests can reuse it without
// duplicating the boilerplate.
func openTempDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "gitmap.sqlite")
	db, err := openDBAt(dbPath)
	if err != nil {
		t.Fatalf("openDBAt: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return db
}

// TestCommitInMigrationCreatesAllTables locks the spec §4 table list
// against what Migrate() actually creates. Adding/removing a table
// from the spec MUST update this list in the same commit.
func TestCommitInMigrationCreatesAllTables(t *testing.T) {
	db := openTempDB(t)
	wantTables := []string{
		"CommitInRun", "RunStatus",
		"InputRepo", "InputKind",
		"SourceCommit", "SourceCommitFile",
		"RewrittenCommit", "CommitOutcome",
		"SkipLog", "SkipReason",
		"ShaMap",
		"Profile", "ProfileExclusion", "ExclusionKind",
		"ProfileMessageRule", "MessageRuleKind",
		"FunctionIntelLanguage", "ConflictMode",
		// Migration 007 — tag replay map (spec §09).
		"CommitInReplayMap", "TagReplayOutcome",
	}
	for _, name := range wantTables {
		if !db.tableExists(name) {
			t.Errorf("commit-in: table %q is missing after Migrate()", name)
		}
	}
}

// TestCommitInMigrationSeedsEnumMirrors verifies that every enum-mirror
// table is seeded with the spec's exact member set (and nothing extra).
// The Go-side parity test in gitmap/cmd/commitin/enums_test.go locks
// the typed enums; this test locks the SQL seeds. Both must agree.
func TestCommitInMigrationSeedsEnumMirrors(t *testing.T) {
	db := openTempDB(t)
	cases := []struct {
		table string
		want  []string
	}{
		{"RunStatus", []string{"Completed", "Failed", "PartiallyFailed", "Pending", "Running"}},
		{"InputKind", []string{"GitUrl", "LocalFolder", "VersionedSibling"}},
		{"CommitOutcome", []string{"Created", "Failed", "Skipped"}},
		{"SkipReason", []string{"DryRun", "DuplicateSourceSha", "EmptyAfterMessageRules", "ExcludedAllFiles"}},
		{"ExclusionKind", []string{"PathFile", "PathFolder"}},
		{"MessageRuleKind", []string{"Contains", "EndsWith", "StartsWith"}},
		{"FunctionIntelLanguage", []string{"CSharp", "Go", "Java", "JavaScript", "Php", "Python", "Rust", "TypeScript"}},
		{"ConflictMode", []string{"ForceMerge", "Prompt"}},
		{"TagReplayOutcome", []string{"AlreadyExists", "Created", "CreatedDryRun", "Failed", "Skipped"}},
	}
	for _, tc := range cases {
		got := selectNames(t, db, tc.table)
		if !equalSorted(got, tc.want) {
			t.Errorf("commit-in seed mismatch on %s:\n  got : %v\n  want: %v",
				tc.table, got, tc.want)
		}
	}
}

// TestCommitInMigrationIsIdempotent asserts that re-running Migrate()
// on an already-migrated DB is a no-op and does not duplicate seed
// rows (INSERT OR IGNORE contract).
func TestCommitInMigrationIsIdempotent(t *testing.T) {
	db := openTempDB(t)
	before := selectNames(t, db, "RunStatus")
	if err := db.migrateCommitIn(); err != nil {
		t.Fatalf("second migrateCommitIn: %v", err)
	}
	after := selectNames(t, db, "RunStatus")
	if !equalSorted(before, after) {
		t.Errorf("commit-in migration is not idempotent:\n  before: %v\n  after : %v",
			before, after)
	}
}

// ---- helpers ------------------------------------------------------

func selectNames(t *testing.T, db *DB, table string) []string {
	t.Helper()
	rows, err := db.conn.Query("SELECT Name FROM " + table)
	if err != nil {
		t.Fatalf("SELECT Name FROM %s: %v", table, err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan %s: %v", table, err)
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func equalSorted(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
