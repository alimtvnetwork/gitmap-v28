package cmd

import (
	"reflect"
	"testing"
)

func TestParseCommitFlags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantClean  []string
		wantPush   bool
		wantDryRun bool
	}{
		{
			name:       "simple message",
			args:       []string{"feat: add new feature"},
			wantClean:  []string{"feat: add new feature"},
			wantPush:   false,
			wantDryRun: false,
		},
		{
			name:       "message with push long flag",
			args:       []string{"fix: resolve bug", "--push"},
			wantClean:  []string{"fix: resolve bug"},
			wantPush:   true,
			wantDryRun: false,
		},
		{
			name:       "message with push short flag",
			args:       []string{"fix: resolve bug", "-p"},
			wantClean:  []string{"fix: resolve bug"},
			wantPush:   true,
			wantDryRun: false,
		},
		{
			name:       "message with dry-run long flag",
			args:       []string{"refactor: clean up", "--dry-run"},
			wantClean:  []string{"refactor: clean up"},
			wantPush:   false,
			wantDryRun: true,
		},
		{
			name:       "message with dry-run short flag",
			args:       []string{"refactor: clean up", "-n"},
			wantClean:  []string{"refactor: clean up"},
			wantPush:   false,
			wantDryRun: true,
		},
		{
			name:       "both push and dry-run flags",
			args:       []string{"--dry-run", "docs: update guide", "-p"},
			wantClean:  []string{"docs: update guide"},
			wantPush:   true,
			wantDryRun: true,
		},
		{
			name:       "empty args",
			args:       []string{},
			wantClean:  nil,
			wantPush:   false,
			wantDryRun: false,
		},
		{
			name:       "flag with -m explicitly",
			args:       []string{"-m", "chore: maintenance", "-p"},
			wantClean:  []string{"-m", "chore: maintenance"},
			wantPush:   true,
			wantDryRun: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clean, push, dryRun := parseCommitFlags(tt.args)
			if !reflect.DeepEqual(clean, tt.wantClean) {
				t.Errorf("parseCommitFlags() clean = %v, want %v", clean, tt.wantClean)
			}
			if push != tt.wantPush {
				t.Errorf("parseCommitFlags() push = %v, want %v", push, tt.wantPush)
			}
			if dryRun != tt.wantDryRun {
				t.Errorf("parseCommitFlags() dryRun = %v, want %v", dryRun, tt.wantDryRun)
			}
		})
	}
}

func TestRunCommitHelp(t *testing.T) {
	// Calling with --help should return nil without executing git
	err := runCommit([]string{"--help"})
	if err != nil {
		t.Errorf("runCommit(--help) unexpected error: %v", err)
	}

	errShort := runCommit([]string{"-h"})
	if errShort != nil {
		t.Errorf("runCommit(-h) unexpected error: %v", errShort)
	}
}
