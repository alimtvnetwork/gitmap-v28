package cmdcursor

import (
	"strings"
	"testing"
)

func TestCursorMigrate_Defaults(t *testing.T) {
	opts, err := parseCursorMigrateArgs([]string{})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if opts.targetNode != "u1" {
		t.Errorf("expected default node 'u1', got: %s", opts.targetNode)
	}

	if opts.isSkipSettings {
		t.Errorf("expected isSkipSettings to default to false")
	}

	if opts.isSettingsOnly {
		t.Errorf("expected isSettingsOnly to default to false")
	}

	if opts.isDryRun {
		t.Errorf("expected isDryRun to default to false")
	}

	if opts.isNoBackup {
		t.Errorf("expected isNoBackup to default to false")
	}

	if opts.isForce {
		t.Errorf("expected isForce to default to false")
	}
}

func TestCursorMigrate_FlagParsing(t *testing.T) {
	args := []string{
		"-n", "u2",
		"--exclude", "myproject,otherrepo",
		"-o", "/tmp/bundle.tar.gz",
		"--skip-settings",
		"-d",
		"--no-backup",
		"-f",
	}

	opts, err := parseCursorMigrateArgs(args)
	if err != nil {
		t.Fatalf("unexpected error parsing flags: %v", err)
	}

	if opts.targetNode != "u2" {
		t.Errorf("expected targetNode 'u2', got: %s", opts.targetNode)
	}

	if opts.exclude != "myproject,otherrepo" {
		t.Errorf("expected exclude 'myproject,otherrepo', got: %s", opts.exclude)
	}

	if opts.outputPath != "/tmp/bundle.tar.gz" {
		t.Errorf("expected outputPath '/tmp/bundle.tar.gz', got: %s", opts.outputPath)
	}

	if !opts.isSkipSettings {
		t.Errorf("expected isSkipSettings to be true")
	}

	if !opts.isDryRun {
		t.Errorf("expected isDryRun to be true")
	}

	if !opts.isNoBackup {
		t.Errorf("expected isNoBackup to be true")
	}

	if !opts.isForce {
		t.Errorf("expected isForce to be true")
	}
}

func TestCursorMigrate_MutuallyExclusiveError(t *testing.T) {
	args := []string{"--skip-settings", "--settings-only"}

	_, err := parseCursorMigrateArgs(args)
	if err == nil {
		t.Fatalf("expected mutual exclusivity error, got nil")
	}

	if !strings.Contains(err.Error(), "cannot specify both --skip-settings and --settings-only") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCursorMigrate_BuildMigrationPyArgs(t *testing.T) {
	opts := cursorMigrateOptions{
		targetNode:     "u3",
		exclude:        "skip-repo",
		outputPath:     "/tmp/custom.tar.gz",
		isSettingsOnly: true,
		isDryRun:       true,
		isNoBackup:     true,
		isForce:        true,
	}

	pyArgs := buildMigrationPyArgs("migrate.py", opts)
	joined := strings.Join(pyArgs, " ")

	expectedTokens := []string{
		"migrate.py",
		"--sync",
		"--node u3",
		"--exclude skip-repo",
		"--output /tmp/custom.tar.gz",
		"--settings-only",
		"--dry-run",
		"--no-backup",
		"--force",
	}

	for _, token := range expectedTokens {
		if !strings.Contains(joined, token) {
			t.Errorf("expected joined args to contain '%s', got: %s", token, joined)
		}
	}
}

func TestCursorMigrate_HelpRequested(t *testing.T) {
	if !isMigrateHelpRequested([]string{"help"}) {
		t.Errorf("expected help to be detected for 'help'")
	}

	if !isMigrateHelpRequested([]string{"--help"}) {
		t.Errorf("expected help to be detected for '--help'")
	}

	if !isMigrateHelpRequested([]string{"-h"}) {
		t.Errorf("expected help to be detected for '-h'")
	}

	if isMigrateHelpRequested([]string{"--node", "u1"}) {
		t.Errorf("expected help not to be detected for regular args")
	}
}
