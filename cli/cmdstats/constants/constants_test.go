package constants

import (
	"strings"
	"testing"
)

func TestCommandNames(t *testing.T) {
	if CmdStats != "stats" {
		t.Errorf("CmdStats = %q, want stats", CmdStats)
	}

	if CmdStatsAlias != "ss" {
		t.Errorf("CmdStatsAlias = %q, want ss", CmdStatsAlias)
	}
}

func TestFlagsAndDescriptions(t *testing.T) {
	if FlagCommand != "command" || FlagJSON != "json" || FlagHelp != "help" {
		t.Errorf("unexpected flag names: %s, %s, %s", FlagCommand, FlagJSON, FlagHelp)
	}

	if FlagArgCommand != "--command" || FlagArgJSON != "--json" || FlagArgHelp != "--help" {
		t.Errorf("unexpected flag args: %s, %s, %s", FlagArgCommand, FlagArgJSON, FlagArgHelp)
	}

	if FlagDescStatsCommand == "" || FlagDescJSON == "" || FlagDescHelp == "" {
		t.Errorf("unexpected empty flag descriptions")
	}
}

func TestWireContractKeys(t *testing.T) {
	topKeys := []string{
		WireKeyTotalCommands,
		WireKeyUniqueCommands,
		WireKeyTotalSuccess,
		WireKeyTotalFail,
		WireKeyOverallFailRate,
		WireKeyAvgDurationMs,
		WireKeyCommands,
	}

	for _, k := range topKeys {
		if k == "" {
			t.Errorf("unexpected empty top-level wire key")
		}
	}

	cmdKeys := []string{
		WireKeyCmdCommand,
		WireKeyCmdTotalRuns,
		WireKeyCmdSuccessCount,
		WireKeyCmdFailCount,
		WireKeyCmdFailRate,
		WireKeyCmdAvgDurationMs,
		WireKeyCmdMinDurationMs,
		WireKeyCmdMaxDurationMs,
		WireKeyCmdLastUsed,
	}

	for _, k := range cmdKeys {
		if k == "" {
			t.Errorf("unexpected empty per-command wire key")
		}
	}
}

func TestAllSSHForwardCandidates(t *testing.T) {
	candidates := AllSSHForwardCandidates()
	if len(candidates) != 11 {
		t.Fatalf("expected 11 SSH forward candidates, got %d", len(candidates))
	}

	seen := make(map[string]bool)
	for _, c := range candidates {
		if c == "" {
			t.Errorf("unexpected empty candidate string")
		}

		if seen[c] {
			t.Errorf("duplicate candidate string: %s", c)
		}
		seen[c] = true
	}

	expectedPrefix := "fix-auth"
	if candidates[0] != expectedPrefix {
		t.Errorf("first candidate = %s, want %s", candidates[0], expectedPrefix)
	}
}

func TestErrorCodesAndPatterns(t *testing.T) {
	if ErrCodeE1154 != "E1154" || ErrCodeE1155 != "E1155" || ErrCodeE1156 != "E1156" {
		t.Errorf("unexpected error codes")
	}

	if !strings.Contains(ErrPatternScanError, "Scan error") {
		t.Errorf("unexpected ErrPatternScanError: %s", ErrPatternScanError)
	}

	if !strings.Contains(ErrPatternDriverValue, "driver.Value") {
		t.Errorf("unexpected ErrPatternDriverValue: %s", ErrPatternDriverValue)
	}

	if OpStatsOpenDB == "" || OpStatsLegacyData == "" || OpStatsQuery == "" || DomainStats == "" {
		t.Errorf("unexpected empty op/domain identifiers")
	}

	if SchemaStatsFilename != "stats.schema.json" {
		t.Errorf("unexpected schema filename: %s", SchemaStatsFilename)
	}
}

func TestExitCodes(t *testing.T) {
	if ExitSuccess != 0 {
		t.Errorf("ExitSuccess = %d, want 0", ExitSuccess)
	}

	if ExitGeneralError != 1 {
		t.Errorf("ExitGeneralError = %d, want 1", ExitGeneralError)
	}
}
