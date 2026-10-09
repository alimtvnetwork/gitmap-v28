package cmdstartup

// Byte-exact JSON contract for `gitmap startup-list --format=json`
// covering the two cases the existing fixtures don't reach:
//
//   1. Multi-entry output — pins the exact bytes between elements
//      (commas, indentation, trailing newline) for a 3-entry list.
//      The structural KeyOrderStable test in the sibling contract
//      file proves keys appear in order, but cannot catch drift in
//      inter-element formatting (e.g., a future change to use
//      `},{` on one line vs. the current `},\n  {`).
//
//   2. Special-character entry — pins the exact escape sequences
//      for embedded quotes, backslashes, control characters, and
//      raw multi-byte UTF-8. Downstream tools that diff JSON output
//      across runs (CI dashboards, jq | sort -u pipelines) need
//      these escapes to be byte-identical or a single-character
//      change in an Exec line shows up as a noisy multi-line diff.
//
// Byte-exact, NOT structural — failure of either test means the
// downstream consumer's diff is no longer stable and the consumer-
// facing changelog must be bumped before regenerating goldens via:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run StartupListJSONBytes
//
// The escape-sequence guarantees in the special-character golden
// come from Go's encoding/json, which `gitmap/stablejson` delegates
// per-value marshaling to. They are stable across Go 1.x releases.
// HTML-escape characters (`<`, `>`, `&`) are deliberately omitted
// so this test doesn't accidentally pin a quirk that would change
// if a future stablejson revision flipped SetEscapeHTML(false).

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/startup"
)

// TestStartupListJSONBytes_MultiEntry pins the exact bytes for a
// representative 3-entry list, including the entry with empty Exec
// (zero-value field rendering must NOT collapse to `null` or be
// omitted by some future omitempty tag).
func TestStartupListJSONBytes_MultiEntry(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: "/bin/b --flag"},
		{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
	}

	assertGoldenBytesDeterministic(t, "startup_list_multi.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStartupListJSON(&buf, entries)

		return buf.Bytes(), err
	})
}

// TestStartupListJSONBytes_SpecialChars pins escape sequences for
// the four character classes that historically cause diff churn:
//
//   - Embedded ASCII double-quote → must escape to \"
//   - Backslash                   → must escape to \\
//   - Control characters (\t, \n) → must escape to \t and \n
//     (NOT to literal whitespace mid-string, which would break
//     line-based downstream parsers)
//   - Sub-\u0020 control byte    → must escape to \u0001
//   - Raw multi-byte UTF-8       → MUST pass through verbatim
//     (no \u-escaping of valid runes — that would make output
//     pointlessly larger and break human-readable diffs)
//
// HTML-escape characters (<, >, &) are intentionally NOT exercised
// here. encoding/json escapes them to \u003c etc. by default, but
// pinning that behavior would couple the contract to an internal
// stdlib detail that's reasonable to flip in the future. Other
// tests cover the schema; this one covers escape stability for the
// classes that downstream pipelines genuinely depend on.
func TestStartupListJSONBytes_SpecialChars(t *testing.T) {
	entries := []startup.Entry{
		{
			// Raw UTF-8 in source — Go source files are UTF-8 and
			// json.Marshal passes valid runes through verbatim.
			Name: "gitmap-中文",
			// Runtime string: /p/"q"\b.desktop → JSON: /p/\"q\"\\b.desktop
			Path: "/p/\"q\"\\b.desktop",
			// Runtime string: a<TAB>b<LF>c<0x01>d → JSON: a\tb\nc\u0001d
			Exec: "a\tb\nc\u0001d",
		},
	}

	assertGoldenBytesDeterministic(t, "startup_list_special.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStartupListJSON(&buf, entries)

		return buf.Bytes(), err
	})
}

// --- local golden helpers (minimal equivalents of cmd's infra) ---

func assertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
	t.Helper()
	first, err := encode()
	if err != nil {
		t.Fatalf("%s: encode run 0: %v", name, err)
	}
	for i := 1; i < 3; i++ {
		got, err := encode()
		if err != nil {
			t.Fatalf("%s: encode run %d: %v", name, i, err)
		}
		if !bytes.Equal(got, first) {
			t.Fatalf("%s: determinism broken — run %d differs from run 0", name, i)
		}
	}
	assertGoldenBytes(t, name, first)
}

func assertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join(startupPackageDir(), "testdata", name)
	if os.Getenv("GITMAP_UPDATE_GOLDEN") == "1" && os.Getenv("GITMAP_ALLOW_GOLDEN_UPDATE") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s", name)
	}
}
