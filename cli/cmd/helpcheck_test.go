package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func setupTestExit(isExited *bool, exitCode *int) func() {
	prevExit := cliexit.SetExitFunc(func(code int) {
		*exitCode = code
		*isExited = true
	})

	return func() { cliexit.SetExitFunc(prevExit) }
}

func assertHelpRequested(t *testing.T, args []string) {
	t.Helper()
	if !IsHelpRequestedOrEmpty(args) {
		t.Errorf("expected help requested for %v", args)
	}
}

func assertHelpIgnored(t *testing.T, args []string) {
	t.Helper()
	if IsHelpRequestedOrEmpty(args) {
		t.Errorf("expected help ignored for %v", args)
	}
}

func TestIsHelpRequestedOrEmpty_EmptyAndNil(t *testing.T) {
	assertHelpRequested(t, []string{})
	assertHelpRequested(t, nil)
}

func TestIsHelpRequestedOrEmpty_HelpFlags(t *testing.T) {
	flags := [][]string{
		{"-h"},
		{"--help"},
		{"help"},
		{"subcmd", "-h"},
		{"subcmd", "--help"},
	}
	for _, args := range flags {
		assertHelpRequested(t, args)
	}
}

func TestIsHelpRequestedOrEmpty_NonHelpArgs(t *testing.T) {
	cases := [][]string{
		{"status"},
		{"clone", "git@github.com:foo/bar"},
		{"--verbose"},
	}
	for _, args := range cases {
		assertHelpIgnored(t, args)
	}
}

func TestCheckHelpOrEmpty_ExitsOnHelp(t *testing.T) {
	var exitCode int
	var isExited bool
	defer setupTestExit(&isExited, &exitCode)()

	CheckHelpOrEmpty("status", []string{"--help"})
	if !isExited {
		t.Errorf("expected exit on help, but did not exit")
	}
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestCheckHelpOrEmpty_ExitsOnEmpty(t *testing.T) {
	var exitCode int
	var isExited bool
	defer setupTestExit(&isExited, &exitCode)()

	CheckHelpOrEmpty("status", []string{})
	if !isExited {
		t.Errorf("expected exit on empty args, but did not exit")
	}
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestCheckHelpOrEmpty_ContinuesOnArgs(t *testing.T) {
	var exitCode int
	var isExited bool
	defer setupTestExit(&isExited, &exitCode)()

	CheckHelpOrEmpty("status", []string{"some-arg"})
	if isExited {
		t.Errorf("expected no exit for non-help args, but exited")
	}
}
