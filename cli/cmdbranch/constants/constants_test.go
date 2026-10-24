package constants

import (
	"strings"
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testFlagConstants(t)
	testFlagArgConstants(t)
	testHelpConstants(t)
	testErrorAndMessageConstants(t)
	testSentinelConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdBranch != "branch" || CmdBranchAlias != "b" {
		t.Error("branch command and alias mismatch")
	}
	if CmdSubDefault != "default" || CmdSubDefaultAlias != "def" {
		t.Error("subcommand default and alias mismatch")
	}
	if CmdSubLatest != "latest-branch" || CmdSubLatestAlias != "lb" {
		t.Error("subcommand latest-branch and alias mismatch")
	}
	if CmdSubList != "list" || CmdSubCurrent != "current" || CmdSubHelp != "help" {
		t.Error("subcommands list/current/help mismatch")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagJSON != "json" || FlagHelp != "help" {
		t.Error("flag json/help mismatch")
	}
	if FlagSort != "sort" || FlagFilter != "filter" || FlagFormat != "format" {
		t.Error("flag sort/filter/format mismatch")
	}
	if FlagDescJSON == "" || FlagDescHelp == "" {
		t.Error("flag descriptions must not be empty")
	}
}

func testFlagArgConstants(t *testing.T) {
	if FlagArgJSON != "--json" || FlagArgJSONShort != "-j" {
		t.Error("json flag arguments mismatch")
	}
	if FlagArgHelp != "--help" || FlagArgHelpShort != "-h" {
		t.Error("help flag arguments mismatch")
	}
	if FlagHelpEntryJSON != "-j, --json" || FlagHelpEntryHelp != "-h, --help" {
		t.Error("help flag entry formatting mismatch")
	}
}

func testHelpConstants(t *testing.T) {
	if !strings.Contains(HelpTitleBranch, "gitmap branch") {
		t.Error("help title mismatch")
	}
	if !strings.HasPrefix(HelpUsageDefault, "gitmap branch") {
		t.Error("default usage line mismatch")
	}
	if !strings.HasPrefix(HelpUsageDefAlias, "gitmap b") {
		t.Error("def alias usage line mismatch")
	}
	if !strings.HasPrefix(HelpUsageLatestBranch, "gitmap latest-branch") {
		t.Error("latest-branch usage line mismatch")
	}
	if HelpSectionOperations == "" || HelpTipResolveDefault == "" || HelpTipOfflineFallback == "" {
		t.Error("help section and tips must not be empty")
	}
}

func testErrorAndMessageConstants(t *testing.T) {
	if ErrCodeE9000 != "E9000" || ErrFatalMessage != "fatal error" {
		t.Error("error code or fatal message mismatch")
	}
	if !strings.Contains(ErrBranchUnknownSubcommand, "Unknown branch subcommand") {
		t.Error("unknown subcommand error mismatch")
	}
	if !strings.Contains(ErrBranchMissingSubcommand, "Missing subcommand") {
		t.Error("missing subcommand error mismatch")
	}
	if !strings.Contains(ErrBranchNotRepo, "Not inside a git work tree") {
		t.Error("not inside repo error mismatch")
	}
	if !strings.Contains(MsgBranchDefaultSwitching, "Switching to default branch") {
		t.Error("switching default branch message mismatch")
	}
}

func testSentinelConstants(t *testing.T) {
	if SentinelCurrentDir != "." {
		t.Error("sentinel current directory mismatch")
	}
	if SentinelDefaultBranchFallback != "main" {
		t.Error("sentinel default branch fallback mismatch")
	}
}
