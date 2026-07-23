package cmd

// clonetermverify_test.go — unit tests for the --verify-cmd-faithful
// machinery. Two scenarios:
//
//   1. Clean match: a CloneTermBlockInput whose buildCloneCommand
//      output equals the executor argv → report.HasMismatch() is
//      false and PrintCmdFaithfulReport writes zero bytes.
//
//   2. Drift: an input whose displayed cmd intentionally diverges
//      from the executor argv (e.g. extra flag, wrong branch) →
//      mismatch list captures every differing position with the
//      right Reason classification.
//
// The intent is not to re-test buildCloneCommand (covered by
// clonetermblock_golden_test.go) but to pin the diff/report
// behavior so a refactor of diffArgvTokens or the report format
// fails loudly.
//
// Test-output hygiene (v4.13.0+): tests that intentionally drive a
// divergent input MUST route the printer through
// PrintCmdFaithfulReportForTest so the emitted "[FAIL]" lines are
// bracketed by an "--- expected mismatch ---" banner. Without that
// wrapper, a green `go test ./...` run prints lines that are byte-
// indistinguishable from a real production verifier failure — which
// has historically trained engineers to ignore "[FAIL]" in test logs.

import (
	"bytes"
	"strings"
	"testing"
)

// TestVerifyCmdFaithful_Match exercises the "no drift" path.
func TestVerifyCmdFaithful_Match(t *testing.T) {
	in := CloneTermBlockInput{
		Index: 1, Name: "repo",
		OriginalURL: "https://x/r.git", TargetURL: "https://x/r.git",
		Dest: "r", CmdBranch: "main",
	}
	executorArgv := []string{"clone", "-b", "main", "https://x/r.git", "r"}
	r := VerifyCmdFaithful(in, executorArgv)
	if r.HasMismatch() {
		t.Fatalf("expected match, got %d mismatches:\n%+v",
			len(r.Mismatches), r.Mismatches)
	}
	var buf bytes.Buffer
	if err := PrintCmdFaithfulReport(&buf, r); err != nil {
		t.Fatalf("print: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty report on match, got:\n%s", buf.String())
	}
}

// TestVerifyCmdFaithful_BranchDrift simulates the canonical regression:
// the executor adds a flag (`--depth=1`) that the displayed cmd forgot.
func TestVerifyCmdFaithful_BranchDrift(t *testing.T) {
	in := CloneTermBlockInput{
		Index: 1, Name: "repo",
		OriginalURL: "https://x/r.git", TargetURL: "https://x/r.git",
		Dest: "r", CmdBranch: "main",
	}
	executorArgv := []string{"clone", "-b", "main", "--depth=1",
		"https://x/r.git", "r"}
	r := VerifyCmdFaithful(in, executorArgv)
	if !r.HasMismatch() {
		t.Fatalf("expected drift, got match. displayed=%q executed=%q",
			r.Displayed, r.Executed)
	}
	// Position-wise diff: an inserted token shifts every subsequent
	// position, so we expect at least one "missing-in-displayed"
	// trailing entry (the last executor token has no displayed
	// counterpart) AND the executed string must contain --depth=1.
	if !strings.Contains(r.Executed, "--depth=1") {
		t.Fatalf("executed should contain --depth=1, got %q", r.Executed)
	}
	var sawMissingInDisplayed bool
	for _, m := range r.Mismatches {
		if m.Reason == "missing-in-displayed" {
			sawMissingInDisplayed = true
		}
	}
	if !sawMissingInDisplayed {
		t.Fatalf("expected ≥1 missing-in-displayed entry, got %+v",
			r.Mismatches)
	}
	var buf bytes.Buffer
	// Use the test-wrapping printer so the captured output is bounded
	// by the "--- expected mismatch ---" banner — keeps go test logs
	// from looking like a real verifier failure.
	if err := PrintCmdFaithfulReportForTest(&buf, r); err != nil {
		t.Fatalf("print: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"--- expected mismatch", "[FAIL]", "verify-cmd-faithful: repo",
		"divergence(s)", "displayed:", "executed:", "--depth=1",
		"--verify-cmd-faithful-exit-on-mismatch",
		"--- end expected mismatch ---",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("report missing %q:\n%s", want, out)
		}
	}
}

// TestVerifyCmdFaithful_WrongBranch ensures position-level diffs are
// classified as "differs" (not as missing-in-X) when both sides have
// a token at the same index.
func TestVerifyCmdFaithful_WrongBranch(t *testing.T) {
	in := CloneTermBlockInput{
		Index: 1, Name: "repo",
		OriginalURL: "https://x/r.git", TargetURL: "https://x/r.git",
		Dest: "r", CmdBranch: "main",
	}
	executorArgv := []string{"clone", "-b", "develop", "https://x/r.git", "r"}
	r := VerifyCmdFaithful(in, executorArgv)
	if !r.HasMismatch() {
		t.Fatal("expected drift on branch-name mismatch")
	}
	if r.Mismatches[0].Reason != "differs" {
		t.Fatalf("expected reason=differs, got %+v", r.Mismatches[0])
	}
}
