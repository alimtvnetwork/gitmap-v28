package e2e

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/workspace"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestPromptModeOnClobberAbortsRun seeds the destination repo with a
// commit that touches `a.txt`, then runs commit-in with --conflict
// Prompt against an input whose first commit also writes a *different*
// `a.txt`. DetectClobbers must fire, finalize.Resolve(Prompt) must
// return Abort, and the orchestrator must propagate
// CommitInExitConflictAborted (=8) without producing any new commits.
func TestPromptModeOnClobberAbortsRun(t *testing.T) {
	src := NewRepo(t, "src")
	src.Commit("a.txt", "DEST-version\n", "dest-seed",
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	startCount := len(src.LogFirstParent(t))

	input := NewRepo(t, "input")
	input.Commit("a.txt", "INPUT-version\n", "would-clobber",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	raw := NewRawArgs(src.Path, input.Path)
	raw.ConflictMode = constants.CommitInConflictModePrompt
	res := Run(t, raw)

	if res.ExitCode != constants.CommitInExitConflictAborted {
		t.Fatalf("exit=%d, want CommitInExitConflictAborted (%d)\nstderr=%s",
			res.ExitCode, constants.CommitInExitConflictAborted, res.Stderr)
	}
	if got := len(src.LogFirstParent(t)); got != startCount {
		t.Fatalf("dst gained %d commits during aborted run", got-startCount)
	}
}

// TestForceMergeOnClobberOverwritesAndContinues mirrors the previous
// test but with --conflict ForceMerge. The orchestrator must take the
// source side, log the clobber, and exit cleanly with the destination
// gaining one new commit. The post-run blob at `a.txt` must be the
// INPUT version, proving the merge resolution actually applied.
func TestForceMergeOnClobberOverwritesAndContinues(t *testing.T) {
	src := NewRepo(t, "src")
	src.Commit("a.txt", "DEST-version\n", "dest-seed",
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

	input := NewRepo(t, "input")
	input.Commit("a.txt", "INPUT-version\n", "force-merge-me",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	raw := NewRawArgs(src.Path, input.Path)
	raw.ConflictMode = constants.CommitInConflictModeForceMerge
	res := Run(t, raw)

	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s\nstdout=%s",
			res.ExitCode, res.Stderr, res.Stdout)
	}
	src.AssertCommitCount(t, 2)
	src.AssertHasSubject(t, "force-merge-me")
	// Audit log: stdout must mention the force-merge clobber.
	if !strings.Contains(res.Stdout, "force-merge clobbering") {
		t.Errorf("stdout missing force-merge clobber audit line\nstdout=%s", res.Stdout)
	}
	// Final blob at a.txt must be the INPUT side.
	got := readBlobAtHead(t, src, "a.txt")
	if got != "INPUT-version\n" {
		t.Fatalf("a.txt @ HEAD = %q, want INPUT-version", got)
	}
}

// TestSecondConcurrentRunFailsWithLockBusy plants a lock file owned by
// PID 1 (always alive in any Linux container) and then attempts a
// commit-in run. AcquireLock must observe the live PID, refuse to
// proceed, and the orchestrator must exit with CommitInExitLockBusy.
//
// Skipped on non-Linux because PID 1 may not be a long-running process
// under macOS/Windows test runners.
func TestSecondConcurrentRunFailsWithLockBusy(t *testing.T) {
	if !isLinux() {
		t.Skip("PID-1-always-alive trick is Linux-specific")
	}
	src := NewRepo(t, "src")
	input := NewRepo(t, "input")
	input.Commit("a.txt", "1\n", "seed", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	// Force the workspace into existence so we know where the lock lives.
	paths, err := workspace.EnsureWorkspace(src.Path)
	if err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}
	if err := os.WriteFile(paths.LockFile, []byte("1"), 0o644); err != nil {
		t.Fatalf("plant lock: %v", err)
	}
	defer os.Remove(paths.LockFile)

	res := Run(t, NewRawArgs(src.Path, input.Path))
	if res.ExitCode != constants.CommitInExitLockBusy {
		t.Fatalf("exit=%d, want CommitInExitLockBusy (%d)\nstderr=%s",
			res.ExitCode, constants.CommitInExitLockBusy, res.Stderr)
	}
}
