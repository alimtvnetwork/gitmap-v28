package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/group"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/render"
)

func TestPrependBothInsertsAtTopWithoutLosingExisting(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "CHANGELOG.md")
	tsPath := filepath.Join(dir, "changelog.ts")

	mdSeed := "# Changelog\n\n## v3.91.0 — (2026-04-24)\n\n### Added\n\n- existing\n\n"
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
			{Name: "Added", Items: []string{"brand new"}},
		},
	}

	if err := PrependBoth(mdPath, tsPath, entry); err != nil {
		t.Fatalf("PrependBoth: %v", err)
	}

	mdOut, _ := os.ReadFile(mdPath)
	tsOut, _ := os.ReadFile(tsPath)

	if !strings.Contains(string(mdOut), "## v3.92.0") || !strings.Contains(string(mdOut), "## v3.91.0") {
		t.Fatalf("Markdown lost an entry:\n%s", mdOut)
	}
	if strings.Index(string(mdOut), "v3.92.0") > strings.Index(string(mdOut), "v3.91.0") {
		t.Fatalf("v3.92.0 must appear before v3.91.0:\n%s", mdOut)
	}

	if !strings.Contains(string(tsOut), `version: "v3.92.0"`) || !strings.Contains(string(tsOut), `version: "v3.91.0"`) {
		t.Fatalf("TypeScript lost an entry:\n%s", tsOut)
	}
	if strings.Index(string(tsOut), "v3.92.0") > strings.Index(string(tsOut), "v3.91.0") {
		t.Fatalf("v3.92.0 must appear before v3.91.0:\n%s", tsOut)
	}
}
