package render

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/group"
)

func sampleEntry() Entry {
	return Entry{
		Version: "v3.92.0",
		Date:    "2026-04-24",
		Groups: []group.Section{
			{Name: "Added", Items: []string{"new flag --foo"}},
			{Name: "Fixed", Items: []string{"crash on empty input"}},
		},
	}
}

func TestMarkdownIncludesHeaderAndSections(t *testing.T) {
	out := Markdown(sampleEntry())

	wantParts := []string{
		"## v3.92.0 — (2026-04-24)",
		"### Added",
		"- new flag --foo",
		"### Fixed",
		"- crash on empty input",
	}
	for _, p := range wantParts {
		if !strings.Contains(out, p) {
			t.Fatalf("Markdown missing %q\n---\n%s", p, out)
		}
	}
}

func TestTypeScriptIsValidObjectFragment(t *testing.T) {
	out := TypeScript(sampleEntry())

	wantParts := []string{
		`version: "v3.92.0"`,
		`date: "2026-04-24"`,
		`"Added: new flag --foo"`,
		`"Fixed: crash on empty input"`,
	}
	for _, p := range wantParts {
		if !strings.Contains(out, p) {
			t.Fatalf("TypeScript missing %q\n---\n%s", p, out)
		}
	}

	if !strings.HasPrefix(out, "  {\n") || !strings.HasSuffix(out, "  },\n") {
		t.Fatalf("TypeScript fragment must be a comma-terminated object literal:\n%s", out)
	}
}
