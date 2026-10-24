package constants

import (
	"strings"
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testFlagConstants(t)
	testFlagArgConstants(t)
	testGitTokensConstants(t)
	testPromptConstants(t)
	testMessageConstants(t)
	testErrorConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdPrune != "prune" || CmdPruneAlias != "prn" {
		t.Error("prune command and alias mismatch")
	}
	if HelpPrune == "" || HelpUsage == "" {
		t.Error("help strings must not be empty")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagDryRun != "dry-run" || FlagConfirm != "confirm" {
		t.Error("dry-run and confirm flags mismatch")
	}
	if FlagRemote != "remote" || FlagHelp != "help" {
		t.Error("remote and help flags mismatch")
	}
	if FlagDescDryRun == "" || FlagDescConfirm == "" {
		t.Error("flag descriptions must not be empty")
	}
	if FlagDescRemote == "" || FlagDescHelp == "" {
		t.Error("flag descriptions must not be empty")
	}
}

func testFlagArgConstants(t *testing.T) {
	if FlagArgDryRun != "--dry-run" || FlagArgConfirm != "--confirm" {
		t.Error("flag argument strings mismatch")
	}
	if FlagArgRemote != "--remote" || FlagArgHelp != "--help" {
		t.Error("flag argument strings mismatch")
	}
	if FlagArgHelpShort != "-h" || FlagHelpEntryHelp != "-h, --help" {
		t.Error("short help flag arguments mismatch")
	}
}

func testGitTokensConstants(t *testing.T) {
	if GitBin != "git" || GitBranch != "branch" {
		t.Error("git binary and branch tokens mismatch")
	}
	if GitBranchListFlag != "--list" || GitBranchDeleteFlag != "-D" {
		t.Error("git branch flag tokens mismatch")
	}
	if GitPush != "push" || GitOrigin != "origin" || GitPushDeleteFlag != "--delete" {
		t.Error("git push tokens mismatch")
	}
	if ReleaseBranchPrefix != "release/" || GitWildcard != "*" {
		t.Error("release branch prefix or wildcard mismatch")
	}
}

func testPromptConstants(t *testing.T) {
	if PromptAnswerY != "y" || PromptAnswerYUpper != "Y" {
		t.Error("prompt answer Y mismatch")
	}
	if PromptAnswerN != "n" || PromptAnswerNUpper != "N" {
		t.Error("prompt answer N mismatch")
	}
	if Newline != "\n" || EmptyString != "" {
		t.Error("delimiters mismatch")
	}
}

func testMessageConstants(t *testing.T) {
	if !strings.Contains(MsgPruneStaleHeader, "Stale release branches") {
		t.Error("stale header message mismatch")
	}
	if !strings.Contains(MsgPruneDryRunHint, "--confirm") {
		t.Error("dry run hint message mismatch")
	}
	if !strings.Contains(MsgPruneSummary, "Summary") {
		t.Error("summary message mismatch")
	}
	if !strings.Contains(MsgPruneNone, "No stale release branches") {
		t.Error("none found message mismatch")
	}
	if !strings.Contains(MsgPrunePrompt, "(y/N)") {
		t.Error("prompt message mismatch")
	}
}

func testErrorConstants(t *testing.T) {
	if ErrCodePruneFailed != "EPRUNE01" {
		t.Error("error code mismatch")
	}
	if !strings.Contains(ErrPruneListBranches, "failed to list branches") {
		t.Error("list branches error format mismatch")
	}
	if !strings.Contains(ErrPruneDeleteBranch, "Failed to delete") {
		t.Error("delete branch error format mismatch")
	}
}
