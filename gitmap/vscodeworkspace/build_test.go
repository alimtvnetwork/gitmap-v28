package vscodeworkspace

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildSortsAndDedupes asserts that Build drops duplicate paths
// (regardless of input order) and emits the remaining folders sorted
// by Name so two runs over the same repo set produce byte-identical
// `.code-workspace` files.
func TestBuildSortsAndDedupes(t *testing.T) {
	in := []Folder{
		{Name: "zeta", Path: "/repos/zeta"},
		{Name: "alpha", Path: "/repos/alpha"},
		{Name: "alpha-dup", Path: "/repos/alpha"},
		{Name: "mid", Path: "/repos/mid"},
	}

	got := Build(in)
	if len(got.Folders) != 3 {
		t.Fatalf("len = %d, want 3 (after dedupe)", len(got.Folders))
	}
	wantOrder := []string{"alpha", "mid", "zeta"}
	for i, f := range got.Folders {
		if f.Name != wantOrder[i] {
			t.Errorf("Folders[%d].Name = %q, want %q", i, f.Name, wantOrder[i])
		}
	}
	if got.Settings == nil {
		t.Errorf("Settings = nil, want non-nil empty map")
	}
}

// TestEncodeMatchesVSCodeShape locks the byte shape of the emitted
// JSON: tab-indented, newline-terminated, "settings": {} present.
func TestEncodeMatchesVSCodeShape(t *testing.T) {
	ws := Build([]Folder{{Name: "demo", Path: "/r/demo"}})

	got, err := Encode(ws)
	if err != nil {
		t.Fatalf("Encode err = %v", err)
	}
	if !strings.HasSuffix(string(got), "\n") {
		t.Errorf("output missing trailing newline: %q", got)
	}
	if !strings.Contains(string(got), "\t\"folders\"") {
		t.Errorf("expected tab-indented folders key, got: %s", got)
	}
	if !strings.Contains(string(got), "\"settings\": {}") {
		t.Errorf("expected empty settings object, got: %s", got)
	}

	// Round-trip: must parse back into the same struct.
	var parsed Workspace
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatalf("re-parse err = %v", err)
	}
	if len(parsed.Folders) != 1 || parsed.Folders[0].Name != "demo" {
		t.Errorf("round-trip mismatch: %+v", parsed)
	}
}

// TestRelativizeProducesForwardSlashes verifies paths land relative
// to baseDir AND use forward slashes (VS Code's preferred form on
// every OS, including Windows).
func TestRelativizeProducesForwardSlashes(t *testing.T) {
	base := string(filepath.Separator) + "work"
	in := []Folder{
		{Name: "a", Path: filepath.Join(base, "alpha")},
		{Name: "b", Path: filepath.Join(base, "nested", "beta")},
	}

	got, err := Relativize(in, base)
	if err != nil {
		t.Fatalf("Relativize err = %v", err)
	}
	if got[0].Path != "alpha" {
		t.Errorf("got[0].Path = %q, want %q", got[0].Path, "alpha")
	}
	if got[1].Path != "nested/beta" {
		t.Errorf("got[1].Path = %q, want %q", got[1].Path, "nested/beta")
	}
}
