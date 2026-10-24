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
	testErrorAndWrapConstants(t)
	testPositionalAndExitConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdDiff != "diff" || CmdDiffAlias != "df" {
		t.Error("diff command mismatch")
	}
	if CmdCompare != "compare" || CmdCompareAlias != "comp" {
		t.Error("compare command mismatch")
	}
	if CmdDiffProfiles != "diff-profiles" || CmdDiffProfilesAlias != "dp" {
		t.Error("diff-profiles command mismatch")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagDiffJSON != "json" || FlagDiffOnlyConflicts != "only-conflicts" {
		t.Error("json/conflicts flag mismatch")
	}
	if FlagDiffOnlyMissing != "only-missing" || FlagDiffIncludeIdentical != "include-identical" {
		t.Error("missing/identical flag mismatch")
	}
	if FlagDiffIncludeVCS != "include-vcs" || FlagDiffIncludeNodeMods != "include-node-modules" {
		t.Error("vcs/node_modules flag mismatch")
	}
	if FlagDiffStat != "stat" || FlagDiffCached != "cached" || FlagDiffStaged != "staged" {
		t.Error("stat/cached/staged flag mismatch")
	}
	if FlagDiffNameOnly != "name-only" || FlagDiffHelp != "help" {
		t.Error("name-only/help flag mismatch")
	}
}

func testFlagArgConstants(t *testing.T) {
	if FlagArgStat != "--stat" || FlagArgCached != "--cached" || FlagArgStaged != "--staged" {
		t.Error("flag argument prefix mismatch")
	}
	if FlagArgNameOnly != "--name-only" || FlagArgJSON != "--json" || FlagArgJSONShort != "-j" {
		t.Error("name-only/json flag argument mismatch")
	}
	if FlagArgHelp != "--help" || FlagArgHelpShort != "-h" {
		t.Error("help flag argument mismatch")
	}
	if FlagHelpEntryCachedStaged != "--cached, --staged" || FlagHelpEntryJSON != "-j, --json" {
		t.Error("help entry formatting mismatch")
	}
}

func testHelpConstants(t *testing.T) {
	if !strings.Contains(HelpTitleDiff, "gitmap diff") {
		t.Error("help title mismatch")
	}
	if !strings.HasPrefix(HelpUsageDiffTarget, "gitmap diff") {
		t.Error("diff target usage line mismatch")
	}
	if !strings.HasPrefix(HelpUsageDiffProfiles, "gitmap diff-profiles") {
		t.Error("diff profiles usage line mismatch")
	}
	if HelpSectionRepoDiff == "" || HelpSectionProfiles == "" {
		t.Error("help sections must not be empty")
	}
	if HelpTipUncommitted == "" || HelpTipDrift == "" {
		t.Error("help tips must not be empty")
	}
}

func testErrorAndWrapConstants(t *testing.T) {
	if WrapResolveEndpointLeft != "resolve-endpoint left" {
		t.Error("WrapResolveEndpointLeft mismatch")
	}
	if WrapResolveEndpointRight != "resolve-endpoint right" {
		t.Error("WrapResolveEndpointRight mismatch")
	}
	if WrapGuardPaths != "guard-paths" {
		t.Error("WrapGuardPaths mismatch")
	}
	if WrapDiffFailed != "diff failed: [diff]" || WrapReportFailed != "report failed: [diff]" {
		t.Error("wrap diff/report failed mismatch")
	}
	if LogPrefixDiff != "[diff]" {
		t.Error("LogPrefixDiff mismatch")
	}
	if !strings.Contains(ErrDiffUsageFmt, "Usage:") || !strings.Contains(ErrDiffSameFolder, "same folder") {
		t.Error("error formatting message mismatch")
	}
}

func testPositionalAndExitConstants(t *testing.T) {
	if ExpectedDiffPositionalArgs != 2 {
		t.Errorf("expected 2 positional args, got %d", ExpectedDiffPositionalArgs)
	}
	if ArgsSubcommandOffset != 2 {
		t.Errorf("expected subcommand offset 2, got %d", ArgsSubcommandOffset)
	}
	if IndexDiffLeft != 0 || IndexDiffRight != 1 {
		t.Error("diff left/right indices mismatch")
	}
	if ExitUsageError != 2 {
		t.Errorf("expected ExitUsageError to be 2, got %d", ExitUsageError)
	}
}
