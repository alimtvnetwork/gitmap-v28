package e2e

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/gitlog"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/group"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/render"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/writer"
)

// generate runs the same pipeline as `make changelog` against the
// repository at `repoRoot`, returning the rendered Markdown and
// TypeScript fragments (without the surrounding seed scaffolding) so
// the caller can diff them against fixtures.
func generate(t *testing.T, repoRoot, version, date string) (string, string) {
	t.Helper()

	commits, _, err := gitlog.CommitsSinceLastTag(repoRoot)
	if err != nil {
		t.Fatalf("gitlog: %v", err)
	}

	groups, _ := group.ByPrefix(commits)

	entry := render.Entry{Version: version, Date: date, Groups: groups}

	mdPath, tsPath := repoRoot+"/CHANGELOG.md", repoRoot+"/src/data/changelog.ts"

	err = writer.PrependBoth(mdPath, tsPath, entry)
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	return extractMarkdown(t, mdPath), extractTypeScript(t, tsPath)
}

// extractMarkdown drops the static `# Changelog\n\n` header so the
// fixture only has to track the new entry.
func extractMarkdown(t *testing.T, path string) string {
	t.Helper()

	full := readFile(t, path)

	return strings.TrimPrefix(full, mdSeed)
}

// extractTypeScript drops the array prologue/epilogue so the fixture
// only has to track the new object literal.
func extractTypeScript(t *testing.T, path string) string {
	t.Helper()

	full := readFile(t, path)

	const marker = "export const changelog: ChangelogEntry[] = [\n"
	body := strings.TrimPrefix(full, marker)

	return strings.TrimSuffix(body, "];\n")
}
