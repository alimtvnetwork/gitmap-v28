package clonenow

// Executor tests focus on the two behaviors that don't require a
// live `git` binary in PATH:
//
//   - PickURL mode + fallback rules (Row.PickURL is pure).
//   - The "non-empty dest -> skip" idempotency rule (executeRow's
//     skip branch is reachable without invoking git).
//
// We deliberately do NOT spawn `git clone` here -- that would make
// the test require network access and a writable tmp under git's
// notion of safe.directory. The skip path is the one that matters
// most for re-run safety and is fully covered.

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestPickURL_Modes(t *testing.T) {
	full := Row{HTTPSUrl: "https://x/a.git", SSHUrl: "git@x:a.git"}
	if got := full.PickURL(constants.CloneNowModeHTTPS); got != full.HTTPSUrl {
		t.Errorf("https mode: %q", got)
	}
	if got := full.PickURL(constants.CloneNowModeSSH); got != full.SSHUrl {
		t.Errorf("ssh mode: %q", got)
	}

	// Fallback: ssh requested but only https present.
	httpsOnly := Row{HTTPSUrl: "https://x/a.git"}
	if got := httpsOnly.PickURL(constants.CloneNowModeSSH); got != httpsOnly.HTTPSUrl {
		t.Errorf("ssh->https fallback: %q", got)
	}
	// Fallback: https requested but only ssh present.
	sshOnly := Row{SSHUrl: "git@x:a.git"}
	if got := sshOnly.PickURL(constants.CloneNowModeHTTPS); got != sshOnly.SSHUrl {
		t.Errorf("https->ssh fallback: %q", got)
	}
	// No URLs at all -> empty (executor reports as failed).
	if got := (Row{}).PickURL(constants.CloneNowModeHTTPS); got != "" {
		t.Errorf("empty pick: %q", got)
	}
}

// TestExecuteRow_NonRepoDestFails pins the new safer behavior:
// a populated non-git directory is now reported as failed (with
// MsgCloneNowNotARepo) instead of the old silent skip. This is
// the "we never silently hide drift" contract from the docstring.
func TestExecuteRow_NonRepoDestFails(t *testing.T) {
	tmp := t.TempDir()
	dest := filepath.Join(tmp, "existing")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dest, "marker"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	row := Row{HTTPSUrl: "https://example.com/x.git", RelativePath: "existing"}
	plan := Plan{Mode: constants.CloneNowModeHTTPS, OnExists: constants.CloneNowOnExistsSkip}
	res := executeRow(row, plan, tmp)
	if res.Status != constants.CloneNowStatusFailed {
		t.Errorf("status = %q, want failed", res.Status)
	}
	if res.Detail != constants.MsgCloneNowNotARepo {
		t.Errorf("detail = %q", res.Detail)
	}
}

func TestExecuteRow_NoURLIsFailure(t *testing.T) {
	plan := Plan{Mode: constants.CloneNowModeHTTPS, OnExists: constants.CloneNowOnExistsSkip}
	res := executeRow(Row{RelativePath: "z"}, plan, t.TempDir())
	if res.Status != constants.CloneNowStatusFailed {
		t.Errorf("status = %q, want failed", res.Status)
	}
	if res.Detail != constants.MsgCloneNowNoURL {
		t.Errorf("detail = %q", res.Detail)
	}
}

func TestRender_DryRunBytes(t *testing.T) {
	plan := Plan{
		Source: "/tmp/scan.json", Format: "json", Mode: "ssh",
		Rows: []Row{
			{RepoName: "a", SSHUrl: "git@x:a.git", RelativePath: "src/a", Branch: "main"},
		},
	}
	var buf bytes.Buffer
	if err := Render(&buf, plan); err != nil {
		t.Fatalf("Render: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"git@x:a.git", "src/a", "main", "mode=ssh"} {
		if !bytesContains(got, want) {
			t.Errorf("render missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderSummary_Tally(t *testing.T) {
	results := []Result{
		{Status: constants.CloneNowStatusOK, URL: "u1", Dest: "d1"},
		{Status: constants.CloneNowStatusSkipped, URL: "u2", Dest: "d2", Detail: "dest exists"},
		{Status: constants.CloneNowStatusFailed, URL: "u3", Dest: "d3", Detail: "boom"},
	}
	var buf bytes.Buffer
	if err := RenderSummary(&buf, results); err != nil {
		t.Fatalf("RenderSummary: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"1 ok", "1 skipped", "1 failed", "boom"} {
		if !bytesContains(got, want) {
			t.Errorf("summary missing %q in:\n%s", want, got)
		}
	}
}

// bytesContains is a tiny helper so the assertions read like
// English. strings.Contains would be fine; this avoids importing
// "strings" just for tests.
func bytesContains(haystack, needle string) bool {
	return len(needle) == 0 || indexOf(haystack, needle) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}

	return -1
}
