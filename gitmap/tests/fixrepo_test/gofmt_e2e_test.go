package fixrepo_test

// End-to-end test: build the gitmap binary, run `gitmap fix-repo --all`
// against a fixture git repo whose tracked Go files contain
// column-aligned map literals straddling the version-token width
// boundary, then assert `gofmt -l .` reports zero output. This is the
// regression test for the v4.8.0 / v4.9.0 post-rewrite gofmt step
// (see .lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md).
//
// The test is skipped when go/gofmt/git aren't on PATH so it doesn't
// false-fail in restricted CI environments. On standard ubuntu-latest
// runners (and the local dev box) all three are present.

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestFixRepoGofmtCleanAfterRewrite is the headline assertion: after
// fix-repo bumps `repo-v9` -> `repo-v12` inside an aligned map literal,
// `gofmt -l .` MUST be silent. Pre-v4.8.0 this would fail because the
// byte-level rewriter widened one key by 2 chars without re-padding
// the surrounding rows.
func TestFixRepoGofmtCleanAfterRewrite(t *testing.T) {
	requireToolsOrSkip(t, "go", "gofmt", "git")

	bin := buildGitmapBinary(t)
	repo := setupFixtureRepo(t, "repo-v9", 12)

	stdoutPath := filepath.Join(t.TempDir(), "stdout")
	stderrPath := filepath.Join(t.TempDir(), "stderr")

	stdoutF, err := os.Create(stdoutPath)
	if err != nil {
		t.Fatalf("create stdout capture file: %v", err)
	}
	stderrF, err := os.Create(stderrPath)
	if err != nil {
		t.Fatalf("create stderr capture file: %v", err)
	}

	cmd := exec.Command(bin, "fix-repo", "--all", "--verbose")
	cmd.Dir = repo
	cmd.Stdout = stdoutF
	cmd.Stderr = stderrF
	err = cmd.Run()
	stdoutF.Close()
	stderrF.Close()

	if err != nil {
		// Read whatever we captured for the failure log.
		out, _ := os.ReadFile(stdoutPath)
		errOut, _ := os.ReadFile(stderrPath)
		t.Fatalf("fix-repo failed: %v\nstdout=%s\nstderr=%s", err, out, errOut)
	}

	stdoutBytes, _ := os.ReadFile(stdoutPath)
	stderrBytes, _ := os.ReadFile(stderrPath)
	out := append(stdoutBytes, stderrBytes...)

	t.Logf("fix-repo output:\n%s", out)

	if !strings.Contains(string(out), "gofmt:") {
		t.Errorf("fix-repo output missing 'gofmt:' summary line; v4.8.0 step did not run.\n%s", out)
	}

	dirty := runGofmtList(t, repo)
	if dirty != "" {
		t.Fatalf("gofmt -l . reported dirty files after fix-repo:\n%s\n--- repo tree dump ---\n%s",
			dirty, dumpGoFiles(t, repo))
	}
}

// requireToolsOrSkip skips the test when any required external tool is
// missing. Keeps CI green on environments without a Go toolchain while
// still hard-failing on standard runners that DO have it.
func requireToolsOrSkip(t *testing.T, tools ...string) {
	t.Helper()
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("required tool %q not on PATH; skipping e2e", tool)
		}
	}
}

// buildGitmapBinary compiles the gitmap binary into the test's temp
// dir and returns its absolute path. Each test gets a fresh build so
// stale binaries can't leak between runs. On Windows the produced
// executable MUST carry the `.exe` suffix or `exec.Command` cannot
// find it via PATH-style lookup — Go's `go build -o foo` writes the
// exact name you ask for and won't append the suffix for you.
func buildGitmapBinary(t *testing.T) string {
	t.Helper()
	repoRoot := findRepoRoot(t)
	binName := "gitmap"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	bin := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(repoRoot, "gitmap")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go build gitmap: %v", err)
	}

	return bin
}

// findRepoRoot walks up from CWD until it finds a `gitmap/` directory
// (the project's canonical layout). Required because t.TempDir() does
// not give the test its source-relative location.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "gitmap", "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatalf("could not locate repo root from %s", dir)
	return ""
}

// setupFixtureRepo, runGofmtList, dumpGoFiles live in
// fixture_helpers_test.go (same package). They build a realistic
// bare-remote fixture whose URL ends in `<base>-vN.git` so fix-repo's
// identity resolver picks the canonical version. Earlier stub copies
// of these helpers used to live here and silently shadowed the real
// ones — see .lovable/memory/issues/2026-05-02-fixrepo-helper-dup.md.
