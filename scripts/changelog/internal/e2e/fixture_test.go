package e2e

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFixtureMarkdownAndTypeScript builds a deterministic git history
// (same commits, same author timestamps, same hash tiebreakers), runs
// the full changelog pipeline, and asserts both regenerated outputs
// match the committed golden fixtures byte-for-byte.
//
// If this test ever fails, either:
//
//   - You changed the renderer / sort contract on purpose — regenerate
//     fixtures with `UPDATE_FIXTURES=1 go test ./internal/e2e/...`,
//     review the diff, and commit the new files; OR
//   - Something regressed silently — DO NOT update the fixtures, fix
//     the regression. The whole point of this test is to catch that.
func TestFixtureMarkdownAndTypeScript(t *testing.T) {
	repoRoot := t.TempDir()

	r := newGitRepo(t, repoRoot)
	seedFiles(t, repoRoot)
	seedHistory(r)
	r.finalize()

	gotMD, gotTS := generate(t, repoRoot, "v1.0.0", "2026-04-24")

	assertFixture(t, "changelog.md.golden", gotMD)
	assertFixture(t, "changelog.ts.golden", gotTS)
}

// seedHistory creates a known git history covering every relevant case:
// multiple sections, items needing lexicographic sort, an un-prefixed
// commit that must be skipped, and two commits sharing a timestamp so
// the hash tiebreak is exercised. Pinned in chronological order.
func seedHistory(r *gitRepo) {
	r.commit("chore: bootstrap", 1700000000)
	r.tag("v0.9.0")
	r.commit("feat: zeta capability", 1700000100)
	r.commit("feat: alpha capability", 1700000200)
	r.commit("Changes", 1700000300) // un-prefixed → skipped
	r.commit("fix: zulu defect", 1700000400)
	r.commit("fix: alpha defect", 1700000400) // same ts as above → hash tiebreak
	r.commit("docs: tidy README", 1700000500)
	r.commit("perf: faster scan", 1700000600)
}

func assertFixture(t *testing.T, name, got string) {
	t.Helper()

	path := filepath.Join("testdata", name)

	if os.Getenv("UPDATE_FIXTURES") == "1" {
		writeFile(t, path, got)
		t.Logf("updated fixture %s", path)

		return
	}

	want := readFile(t, path)
	if got == want {
		return
	}

	t.Fatalf("fixture %s drift:\n--- want ---\n%s\n--- got ---\n%s", path, want, got)
}
