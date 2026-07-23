package fixrepo_test

// Fixture builder for the fix-repo gofmt e2e test. Creates a throwaway
// git repo whose origin URL ends in `<base>-vN` and whose tracked Go
// files contain column-aligned map literals + const blocks crossing
// the version-token width boundary. fix-repo is expected to rewrite
// every `<base>-vN` token to `<base>-v<current>` and then gofmt-clean
// the result.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupFixtureRepo materializes the fixture under a fresh temp dir and
// returns the work-tree path. The bare remote URL ends in
// `<base>-v<current>.git` so fix-repo's identity resolver picks
// <current> as the canonical version and rewrites prior `-vN` tokens.
func setupFixtureRepo(t *testing.T, base string, current int) string {
	t.Helper()
	root := t.TempDir()
	bare := filepath.Join(root, fmt.Sprintf("%s-v%d.git", baseName(base), current))
	work := filepath.Join(root, "work")
	mustGit(t, root, "init", "--bare", bare)
	mustGit(t, root, "clone", bare, work)
	mustGit(t, work, "config", "user.email", "fix@repo.test")
	mustGit(t, work, "config", "user.name", "FixRepo Test")
	// fix-repo's identity resolver only accepts HTTPS / SSH remote
	// URL shapes (see parseRemoteURL in cmd/fixrepo_identity.go) — a
	// raw local-filesystem path triggers E_NO_REMOTE. Override the
	// origin URL with a fake HTTPS one whose final path segment is
	// `<base>-v<current>.git` so the resolver picks <current> as the
	// canonical version, while git itself still pushes/fetches from
	// the local bare repo we cloned from (we never actually run a
	// network op — fix-repo only reads the URL string).
	fakeURL := fmt.Sprintf("https://example.com/fixture/%s-v%d.git",
		baseName(base), current)
	mustGit(t, work, "remote", "set-url", "origin", fakeURL)
	writeFixtureFiles(t, work, baseName(base))
	mustGit(t, work, "add", "-A")
	mustGit(t, work, "commit", "-m", "fixture")

	return work
}

// baseName strips a trailing `-vN` from `base` so the same string can
// be used both as the repo name (`repo-v9`) and as the rewrite base
// (`repo`). Mirrors fix-repo's own Split-RepoVersion logic.
func baseName(s string) string {
	idx := strings.LastIndex(s, "-v")
	if idx < 0 {
		return s
	}
	rest := s[idx+2:]
	for _, c := range rest {
		if c < '0' || c > '9' {
			return s
		}
	}

	return s[:idx]
}

// writeFixtureFiles writes the two Go source files that exercise the
// gofmt-alignment regression: a map literal whose keys span v8/v9/v10
// widths, and a const block whose names embed the same tokens.
func writeFixtureFiles(t *testing.T, work, base string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(work, "go.mod"),
		[]byte("module fixture\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(work, "aligned_map.go"),
		[]byte(alignedMapSource(base)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(work, "aligned_const.go"),
		[]byte(alignedConstSource(base)), 0o644); err != nil {
		t.Fatal(err)
	}
}

// alignedMapSource returns a Go file whose map literal is gofmt-clean
// at v8/v9/v10 widths. The padding column is set by the longest key
// (`<base>-v9 `), so a v9->v12 rewrite shifts the column by 2 chars
// and is what trips gofmt -l . pre-fix.
func alignedMapSource(base string) string {
	return fmt.Sprintf(`package fixture

// AlignedKeys exercises gofmt's column padding: keys end in -vN tokens
// that fix-repo will rewrite, and the value column must stay aligned.
var AlignedKeys = map[string]string{
	"%[1]s-v8":  "eight",
	"%[1]s-v9":  "nine",
	"%[1]s-v10": "ten",
}
`, base)
}

// alignedConstSource exercises the same regression inside a const
// block (the most common shape inside gitmap/constants/).
func alignedConstSource(base string) string {
	return fmt.Sprintf(`package fixture

const (
	Repo8  = "%[1]s-v8"
	Repo9  = "%[1]s-v9"
	Repo10 = "%[1]s-v10"
)
`, base)
}

// mustGit runs git with t.Fatal-on-failure semantics. Stdout is
// suppressed so the test log stays readable; stderr propagates so a
// real failure shows the underlying git error.
func mustGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v in %s: %v", args, dir, err)
	}
}

// runGofmtList executes `gofmt -l .` in dir and returns the trimmed
// output. Empty string means every Go file is gofmt-clean.
func runGofmtList(t *testing.T, dir string) string {
	t.Helper()
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("gofmt -l .: %v", err)
	}

	return strings.TrimSpace(string(out))
}

// dumpGoFiles concatenates every .go file under dir with a header so
// failures surface the exact byte-level layout that gofmt rejected.
func dumpGoFiles(t *testing.T, dir string) string {
	t.Helper()
	var b strings.Builder
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".go") {
			return nil
		}
		data, _ := os.ReadFile(p)
		fmt.Fprintf(&b, "=== %s ===\n%s\n", p, data)

		return nil
	})

	return b.String()
}
