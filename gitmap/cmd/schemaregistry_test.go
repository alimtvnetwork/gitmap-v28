package cmd

// Versioned schema registry for JSON contract tests. Replaces
// inline `[]string{"name","path","exec"}` literals scattered across
// per-feature contract tests with a single source of truth per
// (schema-name, version) tuple stored in cmd/testdata/schemas/.
//
// Why versioned schemas (not just key lists):
//
//   When a developer legitimately adds a new key (say `enabled`)
//   to startup-list rows, every contract test fails at once. A
//   plain "update the inline literal" workflow encourages copy-
//   paste fixes that miss sibling assertions. Versioning forces a
//   conscious choice: bump to v2 (acknowledging downstream
//   consumers must adapt) or back out the change.
//
// Two escape hatches when drift is detected:
//
//   --accept-schema=NAME@vN     → whitelist the version that the
//                                  test ran against. Use after the
//                                  developer has bumped a vN file
//                                  and wants to confirm "yes, this
//                                  is the new contract".
//
//   --update-schema=NAME        → rewrite the existing latest-
//                                  version file in-place with the
//                                  observed keys. Mirrors the
//                                  GITMAP_UPDATE_GOLDEN pattern.
//
// Env-var equivalents: GITMAP_ACCEPT_SCHEMA, GITMAP_UPDATE_SCHEMA.
// When both env and flag are set, FLAG WINS (documented contract).
//
// Multiple comma-separated values supported:
//   --accept-schema=startup-list@v2,find-next@v2
//   GITMAP_ACCEPT_SCHEMA="latest-branch-no-top@v2,latest-branch-with-top@v2"

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

// schemaDir is the relative path from the cmd package to the
// per-(name,version) JSON files. Declared as a var (not const) so
// the registry's own contract tests can swap in a t.TempDir() and
// exercise --update-schema's file-rewrite path without touching the
// real testdata/schemas/ directory. Production callers MUST NOT
// reassign this — the test that swaps always restores it on cleanup.
//
//nolint:gochecknoglobals // Test-time configurable directory.
var schemaDir = "testdata/schemas"

// schemaFilePattern is the glob used to discover all versions of a
// given schema. `*` is the version segment (e.g. "v1"). Filenames
// are <name>.v<N>.json so a sort.Strings on matches yields version
// order (v1 < v10 lexically — fine until N≥10, see findLatestVersion).
const schemaFilePattern = "%s.v*.json"

// envAcceptSchema / envUpdateSchema are the env-var fallbacks for
// the two test flags. Documented in the package doc comment so
// developers find them via `grep` for either form.
const (
	envAcceptSchema = "GITMAP_ACCEPT_SCHEMA"
	envUpdateSchema = "GITMAP_UPDATE_SCHEMA"
)

// schemaAcceptFlag / schemaUpdateFlag are registered on init.
// Comma-separated lists; each entry is `name@vN` (accept) or `name`
// (update). Flag values OVERRIDE env-var values when both are set.
var (
	//nolint:gochecknoglobals // Test flags must register at init.
	schemaAcceptFlag = flag.String("accept-schema", "",
		"Comma-separated NAME@vN entries whitelisted as the expected schema version when contract tests detect drift. Overrides GITMAP_ACCEPT_SCHEMA.")
	//nolint:gochecknoglobals // Test flags must register at init.
	schemaUpdateFlag = flag.String("update-schema", "",
		"Comma-separated schema names whose latest-version JSON file should be rewritten with the observed keys. Overrides GITMAP_UPDATE_SCHEMA.")
)

// schema is the on-disk shape parsed from cmd/testdata/schemas/
// <name>.v<N>.json files. The leading underscore on _Doc tells
// reviewers it's documentation, not a wire field.
type schema struct {
	Name    string   `json:"name"`
	Version int      `json:"version"`
	Keys    []string `json:"keys"`
	Doc     string   `json:"_doc"`
}

// schemaCache memoizes loaded schemas so a 50-test sweep doesn't
// re-parse the same JSON file 50 times. Tests are single-threaded
// per package by default, but the sync.Mutex protects against
// future `t.Parallel()` adoption.
//
//nolint:gochecknoglobals // Test-only cache.
var (
	schemaCache   = map[string]schema{}
	schemaCacheMu sync.Mutex
)

// loadSchema reads the highest-version <name>.vN.json file from the
// schema directory and returns the parsed schema. Fatals the test
// (rather than returning an error) because a missing schema file
// means the test infrastructure itself is misconfigured — there's
// nothing the test author can do to recover at runtime.
func loadSchema(t *testing.T, name string) schema {
	t.Helper()
	schemaCacheMu.Lock()
	defer schemaCacheMu.Unlock()
	if cached, ok := schemaCache[name]; ok {

		return cached
	}
	path, err := findLatestVersion(name)
	if err != nil {
		t.Fatalf("loadSchema(%q): %v", name, err)
	}
	loaded, err := readSchemaFile(path)
	if err != nil {
		t.Fatalf("loadSchema(%q) reading %s: %v", name, path, err)
	}
	schemaCache[name] = loaded

	return loaded
}

// findLatestVersion globs every <name>.v*.json file and returns
// the path with the highest numeric version. Numeric (not lexical)
// comparison so v10 sorts after v9 — matters once a schema lives
// long enough to accumulate ten revisions.
func findLatestVersion(name string) (string, error) {
	pattern := filepath.Join(resolveSchemaDir(), fmt.Sprintf(schemaFilePattern, name))
	matches, err := filepath.Glob(pattern)
	if err != nil {

		return "", fmt.Errorf("glob %s: %w", pattern, err)
	}
	if len(matches) == 0 {

		return "", fmt.Errorf("no schema files matched %s", pattern)
	}
	sort.Slice(matches, func(i, j int) bool {

		return parseVersionFromPath(matches[i]) < parseVersionFromPath(matches[j])
	})

	return matches[len(matches)-1], nil
}

// parseVersionFromPath extracts the integer N from a path ending
// in `.vN.json`. Returns 0 on parse failure so a malformed filename
// sorts to the bottom and a real version wins.
func parseVersionFromPath(path string) int {
	base := filepath.Base(path)
	// Strip ".json"
	stem := strings.TrimSuffix(base, ".json")
	// Find the last ".v" — split on it.
	idx := strings.LastIndex(stem, ".v")
	if idx < 0 {

		return 0
	}
	var n int
	if _, err := fmt.Sscanf(stem[idx+2:], "%d", &n); err != nil {

		return 0
	}

	return n
}

// readSchemaFile loads + json-decodes one schema file with
// whitespace-tolerant parsing. Returns parse errors verbatim so
// the test failure message names the offending file and the
// json.UnmarshalTypeError field.
func readSchemaFile(path string) (schema, error) {
	raw, err := os.ReadFile(path)
	if err != nil {

		return schema{}, fmt.Errorf("read: %w", err)
	}
	var s schema
	if err := json.Unmarshal(raw, &s); err != nil {

		return schema{}, fmt.Errorf("parse: %w", err)
	}

	return s, nil
}
