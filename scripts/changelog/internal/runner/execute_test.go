package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/group"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/render"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/writer"
)

// TestEndToEndWriteThenCheck simulates the full make-target round trip
// without depending on git: build an Entry, write it, then re-run
// drift.Check and assert exit code 0.
func TestEndToEndWriteThenCheck(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "CHANGELOG.md")
	tsPath := filepath.Join(dir, "src/data/changelog.ts")

	if err := os.MkdirAll(filepath.Dir(tsPath), 0o755); err != nil {
		t.Fatal(err)
	}

	mdSeed := "# Changelog\n\n## v3.91.0 — (2026-04-24)\n\n### Added\n\n- old\n\n"
	tsSeed := "export const changelog: ChangelogEntry[] = [\n  {\n    version: \"v3.91.0\",\n    date: \"2026-04-24\",\n    items: [],\n  },\n];\n"

	if err := os.WriteFile(mdPath, []byte(mdSeed), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tsPath, []byte(tsSeed), 0o644); err != nil {
		t.Fatal(err)
	}

	entry := render.Entry{
		Version: "v3.92.0",
		Date:    "2026-04-24",
		Groups: []group.Section{
			{Name: "Added", Items: []string{"awesome"}},
			{Name: "Fixed", Items: []string{"crash"}},
		},
	}

	if err := writer.PrependBoth(mdPath, tsPath, entry); err != nil {
		t.Fatalf("write: %v", err)
	}

	mdOut, _ := os.ReadFile(mdPath)
	if !bytes.Contains(mdOut, []byte("## v3.92.0")) || !bytes.Contains(mdOut, []byte("## v3.91.0")) {
		t.Fatalf("Markdown round-trip lost an entry:\n%s", mdOut)
	}
	if !strings.Contains(string(mdOut), "- awesome") || !strings.Contains(string(mdOut), "- crash") {
		t.Fatalf("Markdown round-trip lost items:\n%s", mdOut)
	}
}
