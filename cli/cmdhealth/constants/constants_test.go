package constants

import (
	"strings"
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testFlagConstants(t)
	testThresholdConstants(t)
	testGhCommandConstants(t)
	testNoteConstants(t)
	testUiAndWorstConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdHealth != "health" || CmdHealthAlias != "hlth" {
		t.Error("health command mismatch")
	}
	if EnvGitmapWorkdir != "GITMAP_WORKDIR" {
		t.Error("GITMAP_WORKDIR env var mismatch")
	}
	if DefaultWorkspaceDir != "workspace" || DefaultReposDir != "repos" {
		t.Error("default workspace directory mismatch")
	}
	if FallbackWorkdir != "." || WrapHealthFlags != "health.flags" {
		t.Error("fallback workdir or wrap token mismatch")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagWorkdir != "workdir" || FlagHelp != "help" {
		t.Error("flag name mismatch")
	}
	if FlagArgWorkdir != "--workdir" || FlagArgHelp != "--help" || FlagArgHelpShort != "-h" {
		t.Error("flag argument mismatch")
	}
	if FlagDescWorkdir != "repo workdir" {
		t.Error("flag description mismatch")
	}
}

func testThresholdConstants(t *testing.T) {
	if MaxDirtyReposDisplay != 5 || MaxRedPipelinesDisplay != 5 {
		t.Error("display limits threshold mismatch")
	}
	if MaxPipelineProbedRepos != 50 || StaleDbDaysCutoff != 30 || HoursPerDay != 24 {
		t.Error("probed repos or cutoff threshold mismatch")
	}
}

func testGhCommandConstants(t *testing.T) {
	if CmdGh != "gh" || ArgGhRun != "run" || ArgGhList != "list" {
		t.Error("gh command tokens mismatch")
	}
	if FlagGhRepo != "-R" || FlagGhLimit != "--limit" || ValGhLimitOne != "1" {
		t.Error("gh repo/limit flags mismatch")
	}
	if FlagGhJSON != "--json" || ValGhConclusion != "conclusion" || JsonConclusionFailure != "\"failure\"" {
		t.Error("gh json conclusion mismatch")
	}
	if PathSeparator != "/" {
		t.Error("path separator mismatch")
	}
}

func testNoteConstants(t *testing.T) {
	if NoteGhNotInstalled != "gh not installed; pipeline check skipped" {
		t.Error("NoteGhNotInstalled mismatch")
	}
	if NoteWorkdirUnreadable != "workdir unreadable; pipeline check skipped" {
		t.Error("NoteWorkdirUnreadable mismatch")
	}
	if NoteNodeRegistryUnreachable != "node registry unreachable" {
		t.Error("NoteNodeRegistryUnreachable mismatch")
	}
	if !strings.Contains(NoteFleetSkewHint, "gitmap nodes version") {
		t.Error("NoteFleetSkewHint mismatch")
	}
}

func testUiAndWorstConstants(t *testing.T) {
	if HeaderWorkstationHealth != "workstation health:" || ItemSeparator != ", " {
		t.Error("ui header or item separator mismatch")
	}
	if !strings.Contains(WorstTriageFmt, "pe triage -> %s (%d failing)") {
		t.Error("WorstTriageFmt mismatch")
	}
	if WorstDbDoctorRepair != "db doctor --repair (stale databases)" {
		t.Error("WorstDbDoctorRepair mismatch")
	}
	if !strings.Contains(WorstDirtyWorkingTreesFmt, "dirty working trees (%d)") {
		t.Error("WorstDirtyWorkingTreesFmt mismatch")
	}
	if WorstNothingUrgent != "nothing urgent — all checks green" {
		t.Error("WorstNothingUrgent mismatch")
	}
	if !strings.Contains(HelpTextHealth, "gitmap health") {
		t.Error("HelpTextHealth mismatch")
	}
}
