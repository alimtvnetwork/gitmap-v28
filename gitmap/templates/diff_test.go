package templates

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTemp drops content into a temp file under t.TempDir() and returns
// the absolute path. Centralized so tests stay one-liner clean.
func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write %q: %v", p, err)
	}

	return p
}

// TestDiffMissingFile pins the "file does not exist" branch: every
// template line should appear as `+` under the banner. This is the
// "would-create" UX path.
func TestDiffMissingFile(t *testing.T) {
	body := []byte("*.log\nnode_modules/\n")
	res, err := Diff(filepath.Join(t.TempDir(), "absent.gitignore"), "ignore/node", body)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.Status != DiffMissingFile {
		t.Fatalf("status = %v, want DiffMissingFile", res.Status)
	}
	if len(res.Hunks) != 3 || res.Hunks[1] != "+*.log" || res.Hunks[2] != "+node_modules/" {
		t.Errorf("unexpected hunks: %v", res.Hunks)
	}
}

// TestDiffMissingBlock pins the "file exists but no gitmap block"
// branch — `add` would insert the block, so every template line is `+`.
func TestDiffMissingBlock(t *testing.T) {
	p := writeTemp(t, ".gitignore", "# user-managed line\n*.swp\n")
	body := []byte("*.log\n")
	res, err := Diff(p, "ignore/node", body)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.Status != DiffMissingBlock {
		t.Fatalf("status = %v, want DiffMissingBlock", res.Status)
	}
	if len(res.Hunks) != 2 || res.Hunks[1] != "+*.log" {
		t.Errorf("unexpected hunks: %v", res.Hunks)
	}
}

// TestDiffNoChange pins the idempotency contract: when the on-disk
// block body byte-equals the template body, Status == DiffNoChange and
// Hunks is empty. Drives the script-friendly exit-code-0 path.
func TestDiffNoChange(t *testing.T) {
	prior := "# >>> gitmap:ignore/node >>>\n*.log\n# <<< gitmap:ignore/node <<<\n"
	p := writeTemp(t, ".gitignore", prior)
	res, err := Diff(p, "ignore/node", []byte("*.log\n"))
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.Status != DiffNoChange {
		t.Fatalf("status = %v, want DiffNoChange", res.Status)
	}
	if len(res.Hunks) != 0 {
		t.Errorf("expected empty hunks, got %v", res.Hunks)
	}
}

// TestDiffBlockChanged pins the "block exists but body differs" branch.
// Output is a flat removal-then-addition (we don't claim LCS alignment;
// the marker block is small and honesty beats false sophistication).
func TestDiffBlockChanged(t *testing.T) {
	prior := "# >>> gitmap:ignore/node >>>\n*.log\n# <<< gitmap:ignore/node <<<\n"
	p := writeTemp(t, ".gitignore", prior)
	res, err := Diff(p, "ignore/node", []byte("*.log\nnode_modules/\n"))
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.Status != DiffBlockChanged {
		t.Fatalf("status = %v, want DiffBlockChanged", res.Status)
	}
	want := []string{
		"@@ gitmap:ignore/node @@",
		"-*.log",
		"+*.log",
		"+node_modules/",
	}
	if len(res.Hunks) != len(want) {
		t.Fatalf("hunks = %v, want %v", res.Hunks, want)
	}
	for i := range want {
		if res.Hunks[i] != want[i] {
			t.Errorf("hunk[%d] = %q, want %q", i, res.Hunks[i], want[i])
		}
	}
}

// TestDiffPreservesBlankLines guards splitDiffLines against the easy
// mistake of collapsing intra-body blank lines. Blank lines in
// templates are visual separators users actively rely on.
func TestDiffPreservesBlankLines(t *testing.T) {
	body := []byte("a\n\nb\n")
	res, err := Diff(filepath.Join(t.TempDir(), "x.gitignore"), "ignore/test", body)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(res.Hunks) != 4 || res.Hunks[2] != "+" {
		t.Errorf("blank line lost in hunks: %v", res.Hunks)
	}
}
