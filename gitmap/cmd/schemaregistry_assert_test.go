package cmd

// Schema-registry assertion helpers + drift-handling. Split from
// schemaregistry_test.go to keep both files under the 200-line
// budget. Public surface for callers:
//
//   assertSchemaKeysArray(t, raw, schemaName)
//       — for `--format=json` array outputs (every object must
//         match the schema key list)
//   assertSchemaKeysObject(t, raw, schemaName)
//       — for `--format=json` single-object outputs
//   assertSchemaKeysSlice(t, schemaName) []string
//       — raw key list for callers (e.g. CSV header check or
//         JSONL per-line key-order test) that drive their own
//         comparison loop.
//
// On drift (extra/missing/reordered keys), each assert routes
// through handleSchemaDrift, which:
//
//   1. If --update-schema=NAME (or env) is set, rewrites
//      <name>.v<latest>.json with the observed keys and PASSES
//      the test. Mirrors GITMAP_UPDATE_GOLDEN.
//   2. Else if --accept-schema=NAME@vN matches the loaded
//      version, PASSES the test (developer has consciously
//      bumped and is acknowledging).
//   3. Else FAILS with a message that prints both the expected
//      and observed key lists AND the exact CLI commands needed
//      to either accept or update.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// assertSchemaKeysArray loads schema `name` and delegates to the
// existing assertEveryObjectKeysExact helper. On drift, hands off
// to handleSchemaDrift which decides whether to accept, update, or
// fail. Returns nothing — failures are reported via t.Errorf so a
// snapshot test can keep running and report multiple drifts at once.
func assertSchemaKeysArray(t *testing.T, raw []byte, name string) {
	t.Helper()
	expected := loadSchema(t, name)
	observed := readEveryObjectKeys(t, raw)
	if len(observed) == 0 {
		t.Fatalf("schema %q: expected at least one object", name)
	}
	for i, got := range observed {
		if !equalStringSlices(got, expected.Keys) {
			handleSchemaDrift(t, expected, got, fmt.Sprintf("array[%d]", i))
			return
		}
	}
}

// assertSchemaKeysFirstObject loads schema `name` and asserts the
// FIRST top-level object's keys match. Handles both shapes the
// existing readFirstObjectKeys helper supports:
//
//   - {...}        — single-object output (e.g. latest-branch)
//   - [{...}, ...] — array output (checks only the first object;
//     use assertSchemaKeysArray to check every row)
//
// Routes through the same accept/update flow as the array variant.
func assertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	expected := loadSchema(t, name)
	got := readFirstObjectKeys(t, raw)
	if !equalStringSlices(got, expected.Keys) {
		handleSchemaDrift(t, expected, got, "first object")
	}
}

// assertSchemaKeysSlice returns the schema's key slice for callers
// that drive their own comparison (CSV header check, JSONL per-line
// key-order assertion). Routing through this function instead of
// inlining the file load means a CSV-only consumer still triggers
// the schema cache hit and still benefits from --update-schema.
func assertSchemaKeysSlice(t *testing.T, name string) []string {
	t.Helper()
	return loadSchema(t, name).Keys
}

// handleSchemaDrift is the single decision point for "the test
// observed keys that don't match the stored schema". One of three
// outcomes — fully driven by accept/update flags + env vars.
func handleSchemaDrift(t *testing.T, expected schema, observed []string, where string) {
	t.Helper()
	if shouldUpdateSchema(expected.Name) {
		if err := writeSchemaFile(expected, observed); err != nil {
			t.Fatalf("--update-schema=%q write failed: %v", expected.Name, err)
		}
		t.Logf("--update-schema=%q: rewrote v%d with observed keys %v",
			expected.Name, expected.Version, observed)
		return
	}
	if isSchemaAccepted(expected.Name, expected.Version) {
		t.Logf("--accept-schema=%s@v%d: drift acknowledged (observed %v)",
			expected.Name, expected.Version, observed)
		return
	}
	t.Errorf("schema drift in %s for %q (loaded v%d)\n  expected: %v\n  observed: %v\n%s",
		where, expected.Name, expected.Version, expected.Keys, observed,
		schemaDriftHowToFix(expected.Name, expected.Version))
}

// schemaDriftHowToFix prints the two ways out of a drift failure
// so the developer doesn't have to grep for the flag name. The
// suggested next-version (N+1) makes the bump-and-acknowledge
// workflow a one-line copy.
func schemaDriftHowToFix(name string, currentVersion int) string {
	return fmt.Sprintf("To accept this drift, either:\n"+
		"  1. Auto-update v%d in place: GITMAP_UPDATE_SCHEMA=%s go test ./cmd/...\n"+
		"  2. Bump to v%d (create %s/%s.v%d.json) and run with: -accept-schema=%s@v%d",
		currentVersion, name,
		currentVersion+1, schemaDir, name, currentVersion+1, name, currentVersion+1)
}

// shouldUpdateSchema returns true when the current test run wants
// to rewrite this schema's stored expectation. Flag overrides env.
// Empty/missing values are treated identically (no opt-in).
func shouldUpdateSchema(name string) bool {
	return listContains(*schemaUpdateFlag, name) ||
		listContains(os.Getenv(envUpdateSchema), name)
}

// isSchemaAccepted returns true when the developer has whitelisted
// this exact (name, version) tuple. Version match is strict — a
// `--accept-schema=startup-list@v3` does NOT acknowledge v2 drift,
// because the whole point is to confirm "I know I'm running against
// v3". Flag overrides env.
func isSchemaAccepted(name string, version int) bool {
	want := fmt.Sprintf("%s@v%d", name, version)
	return listContains(*schemaAcceptFlag, want) ||
		listContains(os.Getenv(envAcceptSchema), want)
}

// listContains splits a comma-separated string and returns true
// when `want` matches any trimmed entry. Centralized so the
// accept/update parsing rules stay identical.
func listContains(commaList, want string) bool {
	if commaList == "" {
		return false
	}
	for _, raw := range strings.Split(commaList, ",") {
		if strings.TrimSpace(raw) == want {
			return true
		}
	}
	return false
}

// writeSchemaFile rewrites the latest-version file for `s` with
// the observed `keys`, preserving the `_doc` field so the in-file
// reviewer guidance survives an --update-schema run. Sorted keys
// are NOT used — the wire order matters and must be preserved.
func writeSchemaFile(s schema, observedKeys []string) error {
	path := filepath.Join(resolveSchemaDir(), fmt.Sprintf("%s.v%d.json", s.Name, s.Version))
	updated := schema{
		Name:    s.Name,
		Version: s.Version,
		Keys:    observedKeys,
		Doc:     s.Doc,
	}
	body, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	// Invalidate cache so a subsequent loadSchema in the same run
	// sees the freshly written file.
	schemaCacheMu.Lock()
	delete(schemaCache, s.Name)
	schemaCacheMu.Unlock()
	return nil
}
