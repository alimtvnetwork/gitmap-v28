package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func TestParseCommitFlags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantClean  []string
		wantPush   bool
		wantDryRun bool
	}{
		{name: "simple", args: []string{"feat: add"}, wantClean: []string{"feat: add"}},
		{name: "push", args: []string{"fix: bug", "-p"}, wantClean: []string{"fix: bug"}, wantPush: true},
		{name: "dry-run", args: []string{"clean", "-n"}, wantClean: []string{"clean"}, wantDryRun: true},
		{name: "both", args: []string{"-n", "docs", "-p"}, wantClean: []string{"docs"}, wantPush: true, wantDryRun: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCommitFlagsMatch(t, tt.args, tt.wantClean, tt.wantPush, tt.wantDryRun)
		})
	}
}

func assertCommitFlagsMatch(t *testing.T, args, wantClean []string, wantPush, wantDryRun bool) {
	clean, push, dryRun := parseCommitFlags(args)
	if !reflect.DeepEqual(clean, wantClean) {
		t.Errorf("parseCommitFlags() clean = %v, want %v", clean, wantClean)
	}

	if push != wantPush || dryRun != wantDryRun {
		t.Errorf("flags push=%v (want %v), dryRun=%v (want %v)", push, wantPush, dryRun, wantDryRun)
	}
}

func TestRunCommitHelp(t *testing.T) {
	err := runCommit([]string{"--help"})
	if err != nil {
		t.Errorf("runCommit(--help) unexpected error: %v", err)
	}

	errShort := runCommit([]string{"-h"})
	if errShort != nil {
		t.Errorf("runCommit(-h) unexpected error: %v", errShort)
	}
}

func TestIsAllCommitRequested(t *testing.T) {
	cases := []struct {
		args []string
		want bool
	}{
		{[]string{"all"}, true},
		{[]string{"--all"}, true},
		{[]string{"-a"}, true},
		{[]string{"feat: msg"}, false},
		{[]string{}, false},
	}

	for _, tc := range cases {
		assertIsAllCommitRequested(t, tc.args, tc.want)
	}
}

func assertIsAllCommitRequested(t *testing.T, args []string, want bool) {
	got := isAllCommitRequested(args)
	if got != want {
		t.Errorf("isAllCommitRequested(%v) = %v, want %v", args, got, want)
	}
}

func TestStripAllFlags(t *testing.T) {
	cases := []struct {
		args []string
		want []string
	}{
		{[]string{"all", "feat: msg"}, []string{"feat: msg"}},
		{[]string{"--all", "-m", "fix"}, []string{"-m", "fix"}},
		{[]string{"all"}, []string{"chore: commit pending changes"}},
		{[]string{}, []string{"chore: commit pending changes"}},
	}

	for _, tc := range cases {
		assertStripAllFlags(t, tc.args, tc.want)
	}
}

func assertStripAllFlags(t *testing.T, args, want []string) {
	got := stripAllFlags(args)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("stripAllFlags(%v) = %v, want %v", args, got, want)
	}
}

func TestHandleNonGitRepoCommit_AbortError(t *testing.T) {
	tempDir := t.TempDir()
	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origWd) }()

	err := handleNonGitRepoCommit([]string{"msg"}, false, false)
	assertNonGitAbortError(t, err)
}

func assertNonGitAbortError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr, isApp := err.(*apperror.AppError)
	if !isApp || appErr == nil {
		t.Fatalf("expected AppError, got %T", err)
	}

	if appErr.Code != "E9001" || appErr.Type != apperror.ErrorTypeAbort {
		t.Errorf("unexpected AppError fields: code=%s type=%s", appErr.Code, appErr.Type)
	}
}

func TestHandleNonGitRepoCommit_DryRunChildRepo(t *testing.T) {
	tempDir := t.TempDir()
	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origWd) }()

	childGit := filepath.Join(tempDir, "child-repo", ".git")
	_ = os.MkdirAll(childGit, 0755)

	err := handleNonGitRepoCommit([]string{"all"}, false, true)
	if err != nil {
		t.Errorf("expected nil for dry run on child repos, got %v", err)
	}
}
