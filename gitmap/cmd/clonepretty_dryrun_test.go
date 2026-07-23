// Package cmd — unit tests for runCloneCommandPretty's --dry-run
// short circuit (v6.49.0). Covers: (1) no `git` subprocess is spawned,
// (2) the header includes the exact argv the user would see,
// (3) the yellow dry-run noop sentinel is emitted, (4) the flag
// resets cleanly so subsequent tests are not affected.
package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestRunCloneCommandPrettyDryRun verifies that toggling
// SetCloneDryRun(true) makes runCloneCommandPretty print the header
// + dry-run noop sentinel and return nil without executing git.
func TestRunCloneCommandPrettyDryRun(t *testing.T) {
	SetCloneDryRun(true)
	SetCloneSpinnerOff(true)
	t.Cleanup(func() {
		SetCloneDryRun(false)
		SetCloneSpinnerOff(false)
	})

	const (
		url  = "https://github.com/acme/repo-v9.git"
		dest = "/tmp/gitmap-dryrun-fixture-does-not-exist"
	)
	out := captureClonePrettyStdout(t, func() {
		if err := runCloneCommandPretty(url, dest); err != nil {
			t.Fatalf("dry-run returned error: %v", err)
		}
	})

	if !strings.Contains(out, url) {
		t.Errorf("header missing url %q in output:\n%s", url, out)
	}
	if !strings.Contains(out, dest) {
		t.Errorf("header missing dest %q in output:\n%s", dest, out)
	}
	if !strings.Contains(out, constants.GitClone) {
		t.Errorf("header missing %q (git subcommand) in output:\n%s",
			constants.GitClone, out)
	}
	if !strings.Contains(out, constants.MsgCloneDryRunNoop) {
		t.Errorf("dry-run sentinel missing in output:\n%s", out)
	}
	if _, err := os.Stat(dest); err == nil {
		t.Errorf("dry-run should not have created %q", dest)
	}
}

// TestRunCloneCommandPrettyDryRunFlagDefaultsOff guards against a
// regression where the package-level flag would leak across tests.
func TestRunCloneCommandPrettyDryRunFlagDefaultsOff(t *testing.T) {
	if IsCloneDryRun() {
		t.Fatalf("cloneDryRunFlag should default to false at package init")
	}
}

func TestWithSSHAcceptNewAddsOption(t *testing.T) {
	got := withSSHAcceptNew("ssh -i ~/.ssh/id_ed25519")
	if !strings.Contains(got, constants.SSHStrictHostKeyAcceptNew) {
		t.Fatalf("missing accept-new option in %q", got)
	}
}

func TestWithSSHAcceptNewKeepsExplicitStrictHostKey(t *testing.T) {
	existing := "ssh -o StrictHostKeyChecking=no"
	if got := withSSHAcceptNew(existing); got != existing {
		t.Fatalf("withSSHAcceptNew() = %q, want %q", got, existing)
	}
}

// captureClonePrettyStdout redirects os.Stdout for the duration of fn and
// returns whatever was written. Failures during pipe setup fail the
// test outright — without stdout we cannot assert anything useful.
func captureClonePrettyStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string, 1)
	go func() {
		buf, readErr := io.ReadAll(r)
		if readErr != nil {
			done <- ""
			return
		}
		done <- string(buf)
	}()

	fn()
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("close pipe writer: %v", closeErr)
	}
	return <-done
}
