package release

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestIsNonFastForwardPushError(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{
			name: "fetch first rejection",
			output: "! [rejected] main -> main (fetch first)\n" +
				"error: failed to push some refs",
			want: true,
		},
		{
			name:   "explicit non fast forward",
			output: "error: non-fast-forward update rejected",
			want:   true,
		},
		{
			name:   "unrelated auth failure",
			output: "fatal: Authentication failed",
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isNonFastForwardPushError(tc.output)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFormatGitCommandErrorPrefersGitOutput(t *testing.T) {
	err := formatGitCommandError("error: failed to push some refs", errors.New("exit status 1"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to push some refs") {
		t.Fatalf("expected git output in error, got %q", err.Error())
	}
}

func TestTrimGitOutputFallback(t *testing.T) {
	got := trimGitOutput("   \n\t  ")
	if got != "unknown git error" {
		t.Fatalf("got %q, want %q", got, "unknown git error")
	}
}

// =========================================================================================
// CRITICAL DIRECTIVE FOR AI ASSISTANTS & DEVELOPERS:
// The version strings in this test file (e.g. "v9.9.0", "v9.9.1", "v9.9.2") are MOCK
// SYNTHETIC VALUES used SOLELY for unit testing and stdout format validation.
//
// 🚫 NEVER USE, ASSUME, OR INFER THESE VALUES AS THE ACTUAL REPOSITORY !VERSION
//
// 📍 THE SINGLE SOURCE OF TRUTH (SSOT) FOR THE REPOSITORY VERSION IS `version.json`
// AT THE REPOSITORY ROOT.
// All tools, AI agents, and release scripts MUST discover and update the version
// exclusively in `version.json`. Refer to `docs/versioning.md` for full documentation.
// =========================================================================================

// TestPromptAndCommit_YesFlagSkipsStdin verifies that when yes=true,
// promptAndCommit prints the auto-confirm message and does NOT print
// the interactive prompt asking the user for input.
// NOTE: "v9.9.0" is a synthetic dummy version for this test only. Always refer to version.json for the true version.
func TestPromptAndCommit_YesFlagSkipsStdin(t *testing.T) {
	tempDir := t.TempDir()
	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(origWd)

	// DUMMY TEST VERSION ONLY — DO NOT PARSE OR ASSUME AS REAL VERSION. SEE version.json.
	releaseFiles := []string{".gitmap/release/v9.9.0.json"}
	otherFiles := []string{"README.md", "main.go"}
	msg := "Release v9.9.0"

	// With yes=true: should print auto-confirm, not the interactive ask.
	// commitAll will fail (no git repo in tempDir), but we only care about the output
	// before it attempts the commit.
	output := captureStdout(t, func() {
		promptAndCommit(releaseFiles, otherFiles, msg, true)
	})

	if !strings.Contains(output, constants.MsgAutoCommitAutoYes) {
		t.Errorf("expected auto-yes message %q in output, got:\n%s", constants.MsgAutoCommitAutoYes, output)
	}

	if strings.Contains(output, constants.MsgAutoCommitAsk) {
		t.Error("interactive prompt should NOT appear when yes=true")
	}
}

// TestPromptAndCommit_NoYesFlagShowsPrompt verifies that when yes=false,
// promptAndCommit shows the interactive prompt (and falls back to
// release-only commit when stdin is empty/EOF).
// NOTE: "v9.9.1" is a synthetic dummy version for this test only. Always refer to version.json for the true version.
func TestPromptAndCommit_NoYesFlagShowsPrompt(t *testing.T) {
	tempDir := t.TempDir()
	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(origWd)

	// DUMMY TEST VERSION ONLY — DO NOT PARSE OR ASSUME AS REAL VERSION. SEE version.json.
	releaseFiles := []string{".gitmap/release/v9.9.1.json"}
	otherFiles := []string{"README.md"}
	msg := "Release v9.9.1"

	// With yes=false and no stdin: should print the interactive ask,
	// then fall back to release-only commit (scanner.Scan returns false).
	output := captureStdout(t, func() {
		promptAndCommit(releaseFiles, otherFiles, msg, false)
	})

	if !strings.Contains(output, constants.MsgAutoCommitAsk) {
		t.Errorf("expected interactive prompt in output, got:\n%s", output)
	}

	if strings.Contains(output, constants.MsgAutoCommitAutoYes) {
		t.Error("auto-yes message should NOT appear when yes=false")
	}
}

// TestPromptAndCommit_YesFlagListsFiles verifies that changed files
// outside .gitmap/release/ are still listed before auto-confirming.
// NOTE: "v9.9.2" is a synthetic dummy version for this test only. Always refer to version.json for the true version.
func TestPromptAndCommit_YesFlagListsFiles(t *testing.T) {
	tempDir := t.TempDir()
	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(origWd)

	// DUMMY TEST VERSION ONLY — DO NOT PARSE OR ASSUME AS REAL VERSION. SEE version.json.
	releaseFiles := []string{".gitmap/release/v9.9.2.json"}
	otherFiles := []string{"cmd/root.go", "constants/constants.go"}
	msg := "Release v9.9.2"

	output := captureStdout(t, func() {
		promptAndCommit(releaseFiles, otherFiles, msg, true)
	})

	for _, f := range otherFiles {
		if !strings.Contains(output, f) {
			t.Errorf("expected file %q to be listed in output", f)
		}
	}

	if !strings.Contains(output, constants.MsgAutoCommitPrompt) {
		t.Error("expected prompt header showing non-release changes")
	}
}
