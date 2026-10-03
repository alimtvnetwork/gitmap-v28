package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func TestIsCommitPushHelpArg(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "empty args",
			args: []string{},
			want: false,
		},
		{
			name: "help string",
			args: []string{"help"},
			want: true,
		},
		{
			name: "double dash help",
			args: []string{"--help"},
			want: true,
		},
		{
			name: "short help flag",
			args: []string{"-h"},
			want: true,
		},
		{
			name: "commit message",
			args: []string{"feat: new feature"},
			want: false,
		},
		{
			name: "other flag",
			args: []string{"--dry-run"},
			want: false,
		},
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
		{
			name:     "empty args",
			args:     []string{},
			wantSha:  "",
			wantPush: true,
		},
		{
			name:     "sha only",
			args:     []string{"abc1234"},
			wantSha:  "abc1234",
			wantPush: true,
		},
		{
			name:     "sha with --no-push flag",
			args:     []string{"abc1234", "--no-push"},
			wantSha:  "abc1234",
			wantPush: false,
		},
		{
			name:     "--local flag with sha",
			args:     []string{"--local", "deadbeef"},
			wantSha:  "deadbeef",
			wantPush: false,
		},
		{
			name:     "-n flag with sha",
			args:     []string{"-n", "feedbeef"},
			wantSha:  "feedbeef",
			wantPush: false,
		},
		{
			name:     "arbitrary flags ignored",
			args:     []string{"--verbose", "cafebabe"},
			wantSha:  "cafebabe",
			wantPush: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSha, gotPush := parseRewriteFlags(tt.args)
			if gotSha != tt.wantSha {
				t.Errorf("parseRewriteFlags() gotSha = %q, want %q", gotSha, tt.wantSha)
			}

			if gotPush != tt.wantPush {
				t.Errorf("parseRewriteFlags() gotPush = %v, want %v", gotPush, tt.wantPush)
			}
		})
	}
}

func TestHasStagedChangesCP_SafeExecution(t *testing.T) {
	hasChanges, err := hasStagedChangesCP()
	if err != nil {
		t.Logf("hasStagedChangesCP returned error in test env: %v", err)

		return
	}

	t.Logf("hasStagedChangesCP returned hasChanges=%v", hasChanges)
}

func TestCountUnpushedCommitsCP_SafeExecution(t *testing.T) {
	count := countUnpushedCommitsCP()
	if count < 0 {
		t.Errorf("countUnpushedCommitsCP() returned negative count: %d", count)
	}

	t.Logf("countUnpushedCommitsCP returned count=%d", count)
}

func TestRunCommitPushHelpAndValidation(t *testing.T) {
	var isExited bool
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = runCommitPush([]string{"--help"})
	if !isExited {
		t.Errorf("expected exit on --help, but did not exit")
	}

	if err := runCommitPush([]string{}); err == nil {
		t.Error("runCommitPush with empty args expected error, got nil")
	}
}

func TestRunPullCommitPushHelpAndValidation(t *testing.T) {
	var isExited bool
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = runPullCommitPush([]string{"--help"})
	if !isExited {
		t.Errorf("expected exit on --help, but did not exit")
	}

	if err := runPullCommitPush([]string{}); err == nil {
		t.Error("runPullCommitPush with empty args expected error, got nil")
	}
}

func TestRunCommitPushSubcommandsHelpAndValidation(t *testing.T) {
	var isExited bool
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = runCommitPushBug([]string{"-h"})
	if !isExited {
		t.Errorf("expected exit on -h, but did not exit")
	}

	if err := runCommitPushBug([]string{}); err == nil {
		t.Error("runCommitPushBug empty args expected error, got nil")
	}

	isExited = false
	_ = runCommitPushFeature([]string{"--help"})
	if !isExited {
		t.Errorf("expected exit on --help, but did not exit")
	}

	if err := runCommitPushFeature([]string{}); err == nil {
		t.Error("runCommitPushFeature empty args expected error, got nil")
	}

	isExited = false
	_ = runCommitPushRelease([]string{"help"})
	if !isExited {
		t.Errorf("expected exit on help, but did not exit")
	}

	if err := runCommitPushRelease([]string{}); err == nil {
		t.Error("runCommitPushRelease empty args expected error, got nil")
	}
}

func TestRunRmGitAndResetValidation(t *testing.T) {
	var isExited bool
	prevExit := cliexit.SetExitFunc(func(code int) {
		isExited = true
	})
	defer cliexit.SetExitFunc(prevExit)

	_ = runRmGit([]string{"--help"})
	if !isExited {
		t.Errorf("expected exit on --help, but did not exit")
	}

	if err := runRmGit([]string{}); err == nil {
		t.Error("runRmGit empty args expected error, got nil")
	}

	if err := runRmGit([]string{"abc"}); err == nil {
		t.Error("runRmGit short SHA expected error, got nil")
	}

	isExited = false
	_ = runGitReset([]string{"-h"})
	if !isExited {
		t.Errorf("expected exit on -h, but did not exit")
	}

	if err := runGitReset([]string{}); err == nil {
		t.Error("runGitReset empty args expected error, got nil")
	}

	if err := runGitReset([]string{"abc"}); err == nil {
		t.Error("runGitReset short SHA expected error, got nil")
	}
}
