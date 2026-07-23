package cmd

// Shared helpers for JSON schema contract tests. Used by:
//
//   - startuplistjson_contract_test.go
//   - latestbranchjson_contract_test.go
//   - findnextjson_contract_test.go
//
// The point of these tests is to make output drift LOUD: any rename,
// reorder, addition, or removal of a JSON field that downstream
// consumers depend on must fail CI before it ships. Two complementary
// strictness levels are offered:
//
//  1. assertGoldenBytes — byte-exact comparison against a committed
//     fixture under testdata/. Catches whitespace, indentation, key
//     order, AND field changes. Used for canonical fixtures (empty
//     list, single hand-authored entry) where the exact output is a
//     deliberate publication.
//
//  2. assertObjectKeyOrder / assertObjectKeyExactSet — structural
//     checks that pin the schema (which keys exist, in what
//     declaration order) without locking down the values. Used for
//     fixtures with variable / time-dependent data so the test
//     stays useful when an underlying field's display format
//     changes for a non-breaking reason.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/goldenguard"
)

// goldenDir is the conventional Go location for committed test
// fixtures — the `testdata/` directory is ignored by `go build` and
// `go vet`, so files there are test-only.
const goldenDir = "testdata"

// assertGoldenBytes byte-compares `got` against the committed file
// `testdata/<name>`. On mismatch it prints a unified-style diff so
// the failure is actionable. To regenerate intentional fixture
// changes, set BOTH GITMAP_UPDATE_GOLDEN=1 AND the dedicated
// GITMAP_ALLOW_GOLDEN_UPDATE=1 (see goldenguard) — the dual gate
// keeps a single stray env var from rewriting fixtures in CI.
func assertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := resolveGoldenPath(name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		mustWriteGolden(t, path, got)

		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v "+
			"(run with GITMAP_UPDATE_GOLDEN=1 and "+
			"GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s\n--- want\n%s\n--- got\n%s",
			name, string(want), string(got))
	}
}

// mustWriteGolden writes a regenerated fixture and fails the test
// loudly so a CI run can never silently pass on a regenerate.
func mustWriteGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir testdata: %v", err)
	}
	if err := os.WriteFile(path, got, 0o644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
	t.Fatalf("regenerated golden %s — re-run without "+
		"GITMAP_UPDATE_GOLDEN to confirm", path)
}

// assertObjectKeyOrder used to live here; removed once no test
// referenced it. Snapshot suites pin shape via assertGoldenBytes
// instead. Re-introduce alongside its first caller if structural-
// only ordering checks are needed again.

// readFirstObjectKeys streams tokens from `raw` and returns the keys
// of the first top-level object encountered. Handles both shapes:
//
//   - {...}             → keys of the object itself.
//   - [{...}, ...]      → keys of the first object in the array.
//
// Anything else (number, string, null, bool, empty array) returns
// an empty slice so the caller's comparison fails with a clear
// "got: []" message rather than a panic.
func readFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := skipUntilFirstObjectStart(dec); err != nil {

		return nil
	}

	return collectObjectKeys(t, dec)
}

// skipUntilFirstObjectStart advances the decoder past wrapper tokens
// (top-level `[`) until it reaches the `{` that opens the first
// object. Returns io.EOF (or similar) if no object is found.
func skipUntilFirstObjectStart(dec *json.Decoder) error {
	for {
		tok, err := dec.Token()
		if err != nil {

			return err
		}
		if delim, ok := tok.(json.Delim); ok && delim == '{' {

			return nil
		}
	}
}

// collectObjectKeys, equalStringSlices, and skipOneValue have been
// consolidated into jsonsnapshot_helpers_test.go (canonical home
// for shared JSON-test helpers). trimTrailingNewline used to live
// here for byte-comparison normalization but was removed once no
// test required it — snapshot helpers compare against goldens with
// a stable trailing terminator.
