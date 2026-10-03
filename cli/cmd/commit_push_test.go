package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func TestIsCommitPushHelpArg(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "empty args", args: []string{}, want: false},
		{name: "help string", args: []string{"help"}, want: true},
		{name: "double dash help", args: []string{"--help"}, want: true},
		{name: "short help flag", args: []string{"-h"}, want: true},
		{name: "commit message", args: []string{"feat: new feature"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCommitPushHelpArg(tt.args)
			if got != tt.want {
				t.Errorf("isCommitPushHelpArg(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestParseRewriteFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantSha  string
		wantPush bool
	}{
		{name: "empty args", args: []string{}, wantSha: "", wantPush: true},
		{name: "sha only", args: []string{"abc1234"}, wantSha: "abc1234", wantPush: true},
		{name: "sha with --no-push", args: []string{"abc1234", "--no-push"}, wantSha: "abc1234", wantPush: false},
		{name: "local flag", args: []string{"--local", "deadbeef"}, wantSha: "deadbeef", wantPush: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRewriteFlagsMatch(t, tt.args, tt.wantSha, tt.wantPush)
		})
	}
}

func assertRewriteFlagsMatch(t *testing.T, args []string, wantSha string, wantPush bool) {
	gotSha, gotPush := parseRewriteFlags(args)
	if gotSha != wantSha || gotPush != wantPush {
		t.Errorf("flags sha=%q (want %q), push=%v (want %v)", gotSha, wantSha, gotPush, wantPush)
	}
}

func TestIsGitNoiseLine(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"warning: in the working copy of 'foo.txt', LF will be replaced by CRLF", true},
		{"CRLF will be replaced by LF the next time Git touches it", true},
		{"The file will have its original line endings in your working directory", true},
		{"warning: in the working copy of bar.go", true},
		{"To https://github.com/org/repo.git", false},
		{"abc1234..def5678  main -> main", false},
		{"", false},
	}

	for _, tc := range cases {
		got := isGitNoiseLine(tc.line)
		if got != tc.want {
			t.Errorf("isGitNoiseLine(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestHandleNonGitRepoCommitPush_AbortError(t *testing.T) {
	appErr := handleNonGitRepoCommitPush()
	if appErr == nil {
		t.Fatal("expected AppError, got nil")
	}

	if appErr.Code != "E9001" || appErr.Type != apperror.ErrorTypeAbort {
		t.Errorf("unexpected AppError fields: code=%s type=%s", appErr.Code, appErr.Type)
	}

	isAbort := isAbortOrReportedError(appErr)
	if !isAbort {
		t.Error("expected isAbortOrReportedError to be true")
	}
}

func TestNewNonGitRepoAbortError(t *testing.T) {
	appErr := newNonGitRepoAbortError("directory is not a git repo")
	if appErr == nil {
		t.Fatal("expected AppError, got nil")
	}

	if appErr.Op != "commit" || appErr.Code != "E9001" {
		t.Errorf("unexpected AppError op=%s code=%s", appErr.Op, appErr.Code)
	}

	if appErr.Type != apperror.ErrorTypeAbort || appErr.Severity != apperror.SeverityWarn {
		t.Errorf("unexpected AppError type=%s sev=%s", appErr.Type, appErr.Severity)
	}
}

func TestRenderCommitPushSummaryCard(t *testing.T) {
	// Smoke test to ensure rendering doesn't panic
	renderCommitPushSummaryCard("main", "abc1234", "feat: test card", "pushed")
}

func TestHasStagedChangesCP_SafeExecution(t *testing.T) {
	hasChanges, err := hasStagedChangesCP()
	if err != nil {
		t.Logf("hasStagedChangesCP error: %v", err)

		return
	}

	t.Logf("hasStagedChangesCP hasChanges=%v", hasChanges)
}

func TestCountUnpushedCommitsCP_SafeExecution(t *testing.T) {
	count := countUnpushedCommitsCP()
	if count < 0 {
		t.Errorf("countUnpushedCommitsCP() returned negative count: %d", count)
	}
}

func TestRunCommitPushHelpAndValidation(t *testing.T) {
	isExited := false
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = runCommitPush([]string{"--help"})
	if !isExited {
		t.Error("expected exit on --help, but did not exit")
	}

	if err := runCommitPush([]string{}); err == nil {
		t.Error("runCommitPush with empty args expected error, got nil")
	}
}

func TestRunPullCommitPushHelpAndValidation(t *testing.T) {
	isExited := false
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = runPullCommitPush([]string{"--help"})
	if !isExited {
		t.Error("expected exit on --help, but did not exit")
	}

	if err := runPullCommitPush([]string{}); err == nil {
		t.Error("runPullCommitPush with empty args expected error, got nil")
	}
}

func TestRunCommitPushSubcommands(t *testing.T) {
	assertSubcommandHelpAndEmpty(t, runCommitPushBug, "-h")
	assertSubcommandHelpAndEmpty(t, runCommitPushFeature, "--help")
	assertSubcommandHelpAndEmpty(t, runCommitPushRelease, "help")
}

func assertSubcommandHelpAndEmpty(t *testing.T, fn func([]string) error, helpFlag string) {
	isExited := false
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = fn([]string{helpFlag})
	if !isExited {
		t.Errorf("expected exit on %s, but did not exit", helpFlag)
	}

	if err := fn([]string{}); err == nil {
		t.Error("expected error on empty args, got nil")
	}
}

func TestRunRmGitValidation(t *testing.T) {
	assertSubcommandHelpAndEmpty(t, runRmGit, "--help")
	if err := runRmGit([]string{"abc"}); err == nil {
		t.Error("runRmGit short SHA expected error, got nil")
	}
}

func TestRunGitResetValidation(t *testing.T) {
	assertSubcommandHelpAndEmpty(t, runGitReset, "-h")
	if err := runGitReset([]string{"abc"}); err == nil {
		t.Error("runGitReset short SHA expected error, got nil")
	}
}
