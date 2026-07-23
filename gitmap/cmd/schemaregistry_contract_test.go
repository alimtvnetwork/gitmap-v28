package cmd

// Contract for the schema registry. Drift-handling is exercised
// via a swap of schemaDir to a t.TempDir() so writes stay out of
// the real testdata/schemas/ tree.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestSchemaRegistry_ParseVersionFromPath(t *testing.T) {
	cases := []struct {
		path string
		want int
	}{
		{"x/foo.v1.json", 1},
		{"x/foo.v9.json", 9},
		{"x/foo.v10.json", 10},
		{"x/long-dashes.v42.json", 42},
		{"x/no-version.json", 0},
		{"x/foo.vXYZ.json", 0},
	}
	for _, tc := range cases {
		if got := parseVersionFromPath(tc.path); got != tc.want {
			t.Errorf("path=%q want %d got %d", tc.path, tc.want, got)
		}
	}
}

func TestSchemaRegistry_FindLatestPicksHighestVersion(t *testing.T) {
	dir := makeTempSchemaDir(t)
	writeSchemaFor(t, dir, "demo", 1, []string{"a"})
	writeSchemaFor(t, dir, "demo", 9, []string{"a", "b"})
	writeSchemaFor(t, dir, "demo", 10, []string{"a", "b", "c"})
	writeSchemaFor(t, dir, "other", 99, []string{"x"})
	got, err := findLatestVersion("demo")
	if err != nil {
		t.Fatalf("findLatestVersion: %v", err)
	}
	if !strings.HasSuffix(got, "demo.v10.json") {
		t.Fatalf("want demo.v10.json, got %s", got)
	}
}

func TestSchemaRegistry_FindLatestErrorsOnMissing(t *testing.T) {
	makeTempSchemaDir(t)
	_, err := findLatestVersion("nonexistent-schema")
	if err == nil || !strings.Contains(err.Error(), "no schema files matched") {
		t.Fatalf("want 'no schema files matched' error, got %v", err)
	}
}

func TestSchemaRegistry_ListContains(t *testing.T) {
	cases := []struct {
		list, want string
		expect     bool
	}{
		{"", "anything", false},
		{"foo", "foo", true},
		{"foo", "bar", false},
		{"foo,bar,baz", "bar", true},
		{"foo, bar , baz", "bar", true},
		{"foo@v1,bar@v2", "bar@v2", true},
		{"foo@v1,bar@v2", "bar@v3", false},
	}
	for _, tc := range cases {
		if got := listContains(tc.list, tc.want); got != tc.expect {
			t.Errorf("listContains(%q,%q): want %v got %v",
				tc.list, tc.want, tc.expect, got)
		}
	}
}

// TestSchemaRegistry_WriteSchemaPreservesDoc proves --update-schema
// rewrites keys but keeps _doc intact, so reviewer guidance survives.
func TestSchemaRegistry_WriteSchemaPreservesDoc(t *testing.T) {
	dir := makeTempSchemaDir(t)
	writeSchemaFor(t, dir, "demo", 1, []string{"a", "b"})
	loaded := loadSchema(t, "demo")
	if loaded.Doc == "" {
		t.Fatalf("setup: doc must be non-empty")
	}
	if err := writeSchemaFile(loaded, []string{"a", "b", "c"}); err != nil {
		t.Fatalf("write: %v", err)
	}
	reloaded := loadSchema(t, "demo")
	if !equalStringSlices(reloaded.Keys, []string{"a", "b", "c"}) {
		t.Fatalf("keys not updated: %v", reloaded.Keys)
	}
	if reloaded.Doc != loaded.Doc {
		t.Fatalf("doc lost\n  before: %q\n  after: %q", loaded.Doc, reloaded.Doc)
	}
}

// TestSchemaRegistry_AcceptIsVersionStrict pins that NAME@v3 in the
// accept list does NOT acknowledge v2 drift — the whole point of
// the version is to confirm "I know which version I'm running against".
func TestSchemaRegistry_AcceptIsVersionStrict(t *testing.T) {
	t.Setenv(envAcceptSchema, "demo@v3")
	if !isSchemaAccepted("demo", 3) {
		t.Fatalf("v3 should be accepted")
	}
	if isSchemaAccepted("demo", 2) {
		t.Fatalf("v2 must NOT be accepted via @v3 entry")
	}
	if isSchemaAccepted("other", 3) {
		t.Fatalf("name mismatch must not match")
	}
}

// TestSchemaRegistry_FlagAndEnvBothHonored verifies env and flag
// values are additive (both honored when both set). When the same
// name appears in both, either source matches — the documented
// "flag wins" rule applies to conflicts on the SAME entry, not to
// disjoint entries.
func TestSchemaRegistry_FlagAndEnvBothHonored(t *testing.T) {
	t.Setenv(envUpdateSchema, "from-env")
	previous := *schemaUpdateFlag
	t.Cleanup(func() { *schemaUpdateFlag = previous })
	*schemaUpdateFlag = "from-flag"
	if !shouldUpdateSchema("from-flag") {
		t.Fatalf("flag value must be honored")
	}
	if !shouldUpdateSchema("from-env") {
		t.Fatalf("env value must also be honored")
	}
}

// TestSchemaRegistry_ProductionSchemasParse asserts the four real
// schema files load cleanly. Catches a malformed JSON file before
// it blows up an unrelated contract test downstream.
func TestSchemaRegistry_ProductionSchemasParse(t *testing.T) {
	for _, name := range []string{
		"startup-list", "find-next",
		"latest-branch-no-top", "latest-branch-with-top",
	} {
		t.Run(name, func(t *testing.T) {
			s := loadSchema(t, name)
			if s.Name != name {
				t.Fatalf("name %q != filename %q", s.Name, name)
			}
			if s.Version < 1 || len(s.Keys) == 0 {
				t.Fatalf("invalid schema: version=%d keys=%v", s.Version, s.Keys)
			}
		})
	}
}

// makeTempSchemaDir swaps schemaDir for a t.TempDir(), restores
// the original on cleanup, and clears the schema cache so freshly
// written fixtures aren't shadowed by previously-loaded entries.
func makeTempSchemaDir(t *testing.T) string {
	t.Helper()
	previous := schemaDir
	dir := t.TempDir()
	schemaDir = dir
	clearSchemaCache()
	t.Cleanup(func() {
		schemaDir = previous
		clearSchemaCache()
	})

	return dir
}

// clearSchemaCache empties the memoized schema map so a test that
// rewrites a file gets a fresh read on the next loadSchema call.
func clearSchemaCache() {
	schemaCacheMu.Lock()
	defer schemaCacheMu.Unlock()
	for k := range schemaCache {
		delete(schemaCache, k)
	}
}

// writeSchemaFor writes a synthetic schema file with a non-empty
// _doc field so the preserve-doc test has something to verify.
func writeSchemaFor(t *testing.T, dir, name string, version int, keys []string) {
	t.Helper()
	body := map[string]any{
		"name":    name,
		"version": version,
		"keys":    keys,
		"_doc":    "synthetic test schema for " + name,
	}
	raw, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(dir, name+".v"+strconv.Itoa(version)+".json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
