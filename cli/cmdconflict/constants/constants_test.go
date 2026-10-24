package constants

import (
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testFlagSideAndToolConstants(t)
	testFlagControlConstants(t)
	testTokenConstants(t)
	testGitCoreConstants(t)
	testGitMergeAndConfigConstants(t)
	testColorAndGlyphConstants(t)
	testHelpAndGuideConstants(t)
	testMessageConstants(t)
	testErrorConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdConflicts != "conflicts" || CmdConflictsAlias != "conflict" {
		t.Error("conflicts command name or alias mismatch")
	}
	if CmdResolve != "resolve" || CmdResolveAlias != "res" {
		t.Error("resolve command name or alias mismatch")
	}
}

func testFlagSideAndToolConstants(t *testing.T) {
	if FlagTheirs != "theirs" || FlagTheirsLong != "--theirs" || FlagTheirsShort != "-t" {
		t.Error("theirs flag mismatch")
	}
	if FlagOurs != "ours" || FlagOursLong != "--ours" || FlagOursShort != "-o" {
		t.Error("ours flag mismatch")
	}
	if FlagIDE != "ide" || FlagIDELong != "--ide" || FlagIDEShort != "-i" || FlagCodeLong != "--code" {
		t.Error("ide flag mismatch")
	}
}

func testFlagControlConstants(t *testing.T) {
	if FlagAbort != "abort" || FlagAbortLong != "--abort" || FlagAbortShort != "-a" {
		t.Error("abort flag mismatch")
	}
	if FlagRerere != "rerere" || FlagRerereLong != "--rerere" {
		t.Error("rerere flag mismatch")
	}
	if FlagHelp != "help" || FlagHelpLong != "--help" || FlagHelpShort != "-h" {
		t.Error("help flag mismatch")
	}
}

func testTokenConstants(t *testing.T) {
	if TokenTheirs != "theirs" || TokenOurs != "ours" || TokenIDE != "ide" {
		t.Error("theirs/ours/ide token mismatch")
	}
	if TokenCode != "code" || TokenResolve != "resolve" || TokenAbort != "abort" || TokenRerere != "rerere" {
		t.Error("code/resolve/abort/rerere token mismatch")
	}
}

func testGitCoreConstants(t *testing.T) {
	if GitBin != "git" || GitFlagDir != "-C" || GitSubCmdDiff != "diff" {
		t.Error("git binary or diff subcommand mismatch")
	}
	if GitArgNameOnly != "--name-only" || GitDiffFilterUnmerged != "--diff-filter=U" {
		t.Error("git diff args mismatch")
	}
	if GitSubCmdStatus != "status" || GitArgPorcelain != "--porcelain" {
		t.Error("git status args mismatch")
	}
	if GitSubCmdCheckout != "checkout" || GitArgTheirs != "--theirs" || GitArgOurs != "--ours" {
		t.Error("git checkout args mismatch")
	}
}

func testGitMergeAndConfigConstants(t *testing.T) {
	if GitSubCmdMerge != "merge" || GitSubCmdRebase != "rebase" || GitSubCmdCherryPick != "cherry-pick" {
		t.Error("git merge/rebase/cherry-pick mismatch")
	}
	if GitArgAbort != "--abort" || GitConfigRerereEnabled != "rerere.enabled" || GitTrue != "true" {
		t.Error("git abort or rerere config mismatch")
	}
	if GitArgDoubleDash != "--" || GitArgAdd != "add" || GitSubCmdRevParse != "rev-parse" {
		t.Error("git utility flags mismatch")
	}
}

func testColorAndGlyphConstants(t *testing.T) {
	if ColorGreen != "\033[32m" || ColorYellow != "\033[33m" || ColorReset != "\033[0m" {
		t.Error("ANSI color tokens mismatch")
	}
	if GlyphCheck != "✓" || GlyphWarn != "⚠" {
		t.Error("glyph tokens mismatch")
	}
	if TreePrefixBranch != "├──" || TreePrefixLeaf != "└──" {
		t.Error("tree prefix glyphs mismatch")
	}
}

func testHelpAndGuideConstants(t *testing.T) {
	if HelpSectionQuickActions == "" || HelpActionIDE == "" || HelpActionTheirs == "" {
		t.Error("quick action help lines must not be empty")
	}
	if HelpActionOurs == "" || HelpActionAbort == "" || HelpActionCommit == "" {
		t.Error("quick action lines must not be empty")
	}
	if GuideDivider == "" || GuideTitle == "" || GuideStep1 == "" || GuideStep6 == "" {
		t.Error("merge editor guide lines must not be empty")
	}
}

func testMessageConstants(t *testing.T) {
	if MsgNoConflictsDetectedFmt == "" || MsgAcceptedTheirsFmt == "" || MsgAcceptedOursFmt == "" {
		t.Error("message formats must not be empty")
	}
	if MsgNextCommitInstruction == "" || MsgIDELaunchedFmt == "" || MsgMergeAbortedFmt == "" {
		t.Error("message formats must not be empty")
	}
	if MsgRebaseAbortedFmt == "" || MsgCherryPickAbortedFmt == "" || MsgRerereEnabledFmt == "" {
		t.Error("abort/rerere message formats must not be empty")
	}
}

func testErrorConstants(t *testing.T) {
	if ErrNoSupportedIDE == "" || ErrNoAbortInProgressFmt == "" || ErrResolveUsage == "" {
		t.Error("error constants must not be empty")
	}
	if ErrRepoNotFoundFmt == "" || ErrCannotReadRepoStore == "" || ErrCannotListRepos == "" {
		t.Error("error constants must not be empty")
	}
}
