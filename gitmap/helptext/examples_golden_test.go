// Package helptext_test — golden assertion that every command markdown
// file contains a runnable `## Examples` section (#19). The CI gate
// keeps docs honest: a new command without examples fails the build.
package helptext_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEveryHelpFileHasExamples enforces a runnable-examples section
// on every command markdown file under gitmap/helptext/.
func TestEveryHelpFileHasExamples(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read helptext dir: %v", err)
	}
	var missing []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if isExemptHelpFile(e.Name()) {
			continue
		}
		body, err := os.ReadFile(filepath.Join(".", e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		if !hasExamplesSection(string(body)) {
			missing = append(missing, e.Name())
		}
	}
	if len(missing) > 0 {
		t.Fatalf("help files missing `## Examples` section (%d):\n  - %s",
			len(missing), strings.Join(missing, "\n  - "))
	}
}

// hasExamplesSection returns true when md contains a level-2 heading
// titled Examples followed by a fenced code block (the "runnable" gate).
func hasExamplesSection(md string) bool {
	idx := strings.Index(md, "\n## Examples")
	if idx < 0 {
		return false
	}
	return strings.Contains(md[idx:], "```")
}

// isExemptHelpFile lists markdown files that are not per-command help
// (e.g. shared overview pages) and therefore exempt from the gate.
func isExemptHelpFile(name string) bool {
	switch name {
	case "README.md", "_overview.md":
		return true
	}
	return false
}
