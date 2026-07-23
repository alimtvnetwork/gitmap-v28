// Package cmd — integration tests for `stale`, `orphans`, `dedupe`,
// `size` using two real git repos in a temp directory. Skips when
// the host `git` binary is missing (e.g. minimal CI containers).
// v6.71.0.
package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"testing"
)

// makeRepo initialises a real git repo at dir with a single commit
// whose contents differ when uniqueBody is true.
func makeRepo(t *testing.T, dir string, uniqueBody bool) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	body := "shared\n"
	if uniqueBody {
		body = "unique-" + filepath.Base(dir) + "\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q", "-b", "main")
	run("add", ".")
	run("commit", "-q", "-m", "init")
}

// requireGit skips the test when git is unavailable on PATH.
func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
}

// TestHygieneIntegrationScansAndProbes covers stale/dedupe/size shared
// helpers against two repos with identical contents (so dedupe groups
// them) inside a temp root.
func TestHygieneIntegrationScansAndProbes(t *testing.T) {
	requireGit(t)
	root := t.TempDir()
	a := filepath.Join(root, "a")
	b := filepath.Join(root, "b")
	makeRepo(t, a, false)
	makeRepo(t, b, false)
	if err := os.MkdirAll(filepath.Join(root, "not-a-repo"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	repos := scanForReposParallel(root)
	if len(repos) != 2 {
		t.Fatalf("scanForReposParallel got %d repos, want 2: %v", len(repos), repos)
	}
	for _, r := range repos {
		if _, ok := lastCommitTime(r); !ok {
			t.Fatalf("lastCommitTime(%s) failed", r)
		}
		if sz := dirSize(filepath.Join(r, ".git")); sz <= 0 {
			t.Fatalf("dirSize(%s) = %d, want > 0", r, sz)
		}
	}
	groups := map[string][]string{}
	for _, r := range repos {
		sha, ok := headTreeSHA(r)
		if !ok {
			t.Fatalf("headTreeSHA(%s) failed", r)
		}
		groups[sha] = append(groups[sha], r)
	}
	dupes := filterDuplicateGroups(groups)
	if len(dupes) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(dupes))
	}
}

// TestHygieneIntegrationOrphanProbe verifies originURL on a repo with
// a configured remote, and gitURLToHTTPS conversion edge cases.
func TestHygieneIntegrationOrphanProbe(t *testing.T) {
	requireGit(t)
	root := t.TempDir()
	a := filepath.Join(root, "a")
	makeRepo(t, a, true)
	cmd := exec.Command("git", "-C", a, "remote", "add", "origin", "git@github.com:owner/repo.git")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("remote add: %v\n%s", err, out)
	}
	u, ok := originURL(a)
	if !ok || u != "git@github.com:owner/repo.git" {
		t.Fatalf("originURL = %q ok=%v", u, ok)
	}
	if got := gitURLToHTTPS(u); got != "https://github.com/owner/repo" {
		t.Fatalf("gitURLToHTTPS = %q", got)
	}
}

// TestParseHygieneFormat exercises the format flag parser.
func TestParseHygieneFormat(t *testing.T) {
	cases := map[string]hygieneFormat{
		"":      hygieneFormatTable,
		"table": hygieneFormatTable,
		"json":  hygieneFormatJSON,
		"csv":   hygieneFormatCSV,
	}
	for in, want := range cases {
		got, err := parseHygieneFormat(in)
		if err != nil || got != want {
			t.Fatalf("parseHygieneFormat(%q) = %q,%v want %q", in, got, err, want)
		}
	}
	if _, err := parseHygieneFormat("xml"); err == nil {
		t.Fatalf("expected error for invalid format")
	}
}

// TestEmitJSONAndCSV captures stdout and decodes the emitted payload
// to confirm the schema used by stale/dedupe/size/orphans is parseable.
func TestEmitJSONAndCSV(t *testing.T) {
	type row struct {
		Path string `json:"path"`
	}
	withStdout(t, func() { emitJSON([]row{{Path: "a"}, {Path: "b"}}) }, func(buf []byte) {
		var got []row
		if err := json.Unmarshal(buf, &got); err != nil {
			t.Fatalf("json: %v\n%s", err, buf)
		}
		if len(got) != 2 || got[0].Path != "a" {
			t.Fatalf("unexpected json: %+v", got)
		}
	})
	withStdout(t, func() { emitCSV([]string{"path"}, [][]string{{"a"}, {"b"}}) }, func(buf []byte) {
		r := csv.NewReader(bytes.NewReader(buf))
		recs, err := r.ReadAll()
		if err != nil {
			t.Fatalf("csv: %v", err)
		}
		if len(recs) != 3 || recs[0][0] != "path" || recs[2][0] != "b" {
			t.Fatalf("unexpected csv: %v", recs)
		}
	})
}

// withStdout redirects os.Stdout for the duration of fn and feeds the
// captured bytes to inspect. Restores stdout on every exit path.
func withStdout(t *testing.T, fn func(), inspect func([]byte)) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()
	done := make(chan []byte, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		done <- buf.Bytes()
	}()
	fn()
	_ = w.Close()
	inspect(<-done)
}
