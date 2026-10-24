package constants

import "testing"

func TestCommandConstants(t *testing.T) {
	if CmdStatus != "status" {
		t.Errorf("expected CmdStatus == 'status', got %q", CmdStatus)
	}

	if CmdStatusAlias != "st" {
		t.Errorf("expected CmdStatusAlias == 'st', got %q", CmdStatusAlias)
	}
}

func TestFlagConstants(t *testing.T) {
	flags := []string{
		FlagGroup,
		FlagGroupShort,
		FlagAll,
		FlagDirty,
		FlagOnlyDirty,
		FlagAhead,
		FlagBehind,
		FlagTable,
		FlagCompact,
		FlagJSON,
		FlagJSONShort,
		FlagHelp,
		FlagHelpShort,
		FlagDirtyShort,
	}

	for _, f := range flags {
		if len(f) == 0 {
			t.Errorf("expected flag constant to be non-empty")
		}
	}
}

func TestTableColumnConstants(t *testing.T) {
	cols := []string{
		TableColRepository,
		TableColBranch,
		TableColStatus,
		TableColSync,
		TableColStash,
		TableColFiles,
	}

	for _, col := range cols {
		if len(col) == 0 {
			t.Errorf("expected column constant to be non-empty")
		}
	}
}

func TestUIStringConstants(t *testing.T) {
	if StatusNotFoundText != "✖ not found" {
		t.Errorf("expected StatusNotFoundText == '✖ not found', got %q", StatusNotFoundText)
	}

	if StatusTableColGap != "   " {
		t.Errorf("expected StatusTableColGap == '   ', got %q", StatusTableColGap)
	}

	if StatusCleanSuccessMsg != "✨ All repositories are !clean" {
		t.Errorf("expected StatusCleanSuccessMsg == '✨ All repositories are !clean', got %q", StatusCleanSuccessMsg)
	}

	if len(StatusHelpTitle) == 0 {
		t.Errorf("expected StatusHelpTitle to be non-empty")
	}
}
