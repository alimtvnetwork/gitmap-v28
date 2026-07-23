// Package store — ERD parity test.
//
// Source of truth: every SQLCreate* constant in gitmap/constants/constants_*.go
// MUST be represented as a table block in spec/01-app/gitmap-database-erd.mmd.
//
// Why name-only (not column-level): the test is intentionally a name-set
// equality check. Column tweaks happen frequently and would churn the ERD on
// every migration; whole-table drift is the failure mode that actually hurt us
// (the v3.5.0 ERD was missing 11 tables and nobody noticed for ~6 months —
// see .lovable/memory/03-v3.12.1-session.md).
//
// Spec follow-up: spec/01-app/gitmap-database-erd.mmd.
package store

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

const (
	// erdPathRel is the canonical ERD location relative to the repo root.
	erdPathRel = "spec/01-app/gitmap-database-erd.mmd"
	// constantsDirRel is the directory containing SQLCreate* constants.
	constantsDirRel = "gitmap/constants"
	// erdParityRegenHint is shown when the test fails so the fix is obvious.
	erdParityRegenHint = "Add the missing table block to " + erdPathRel +
		" — even a stub like `TableName { INTEGER TableNameId PK }` is enough."
)

// reCreateTable extracts the table name from `CREATE TABLE IF NOT EXISTS Foo (`
// and `CREATE TABLE IF NOT EXISTS "Foo" (` (the latter for SQL keywords).
var reCreateTable = regexp.MustCompile(`CREATE TABLE IF NOT EXISTS "?([A-Za-z_][A-Za-z0-9_]*)"?`)

// reErdTable matches a top-level table block: `    Foo {` (any indent).
var reErdTable = regexp.MustCompile(`^\s+([A-Z][A-Za-z0-9_]*)\s*\{\s*$`)

// repoRoot walks upward from this test file until it finds the spec/ dir.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed — cannot locate ERD parity test")
	}
	dir := filepath.Dir(thisFile)
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, erdPathRel)); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("could not locate repo root from %s (looking for %s)", thisFile, erdPathRel)
	return ""
}

func collectSQLCreateTables(t *testing.T, root string) map[string]struct{} {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, constantsDirRel))
	if err != nil {
		t.Fatalf("read constants dir: %v", err)
	}
	tables := make(map[string]struct{})
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasPrefix(name, "constants_") || !strings.HasSuffix(name, ".go") {
			continue
		}
		body, err := os.ReadFile(filepath.Join(root, constantsDirRel, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, m := range reCreateTable.FindAllStringSubmatch(string(body), -1) {
			tables[m[1]] = struct{}{}
		}
	}
	return tables
}

func collectErdTables(t *testing.T, root string) map[string]struct{} {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, erdPathRel))
	if err != nil {
		t.Fatalf("read ERD: %v", err)
	}
	tables := make(map[string]struct{})
	for _, line := range strings.Split(string(body), "\n") {
		// Skip mermaid relationship lines; only struct-block headers count.
		if strings.Contains(line, "--") {
			continue
		}
		if m := reErdTable.FindStringSubmatch(line); m != nil {
			tables[m[1]] = struct{}{}
		}
	}
	return tables
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// TestERDMatchesSQLCreate enforces table-name parity between the SQL constants
// and the canonical ERD. See package doc for rationale.
func TestERDMatchesSQLCreate(t *testing.T) {
	root := repoRoot(t)
	sqlTables := collectSQLCreateTables(t, root)
	erdTables := collectErdTables(t, root)

	missingFromErd := make(map[string]struct{})
	for k := range sqlTables {
		if _, ok := erdTables[k]; !ok {
			missingFromErd[k] = struct{}{}
		}
	}
	extraInErd := make(map[string]struct{})
	for k := range erdTables {
		if _, ok := sqlTables[k]; !ok {
			extraInErd[k] = struct{}{}
		}
	}

	if len(missingFromErd) == 0 && len(extraInErd) == 0 {
		return
	}

	if len(missingFromErd) > 0 {
		t.Errorf("ERD parity drift — %d table(s) declared in gitmap/constants/constants_*.go SQLCreate* but missing from %s:\n  %s\n\nFix: %s",
			len(missingFromErd), erdPathRel,
			strings.Join(sortedKeys(missingFromErd), ", "),
			erdParityRegenHint,
		)
	}
	if len(extraInErd) > 0 {
		t.Errorf("ERD parity drift — %d table(s) present in %s but no matching SQLCreate* constant:\n  %s\n\nFix: either add the SQLCreate* constant or remove the orphan block from the ERD.",
			len(extraInErd), erdPathRel,
			strings.Join(sortedKeys(extraInErd), ", "),
		)
	}
}
