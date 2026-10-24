package constants

import (
	"strings"
	"testing"
)

func TestCommandConstants(t *testing.T) {
	if CmdPlans != "plans" {
		t.Errorf("expected CmdPlans to be 'plans', got %q", CmdPlans)
	}
	if CmdPlansAlias != "plan" {
		t.Errorf("expected CmdPlansAlias to be 'plan', got %q", CmdPlansAlias)
	}
}

func TestFlagConstants(t *testing.T) {
	flags := []string{
		FlagLimit, FlagAll, FlagPending, FlagCompleted,
		FlagJSON, FlagYAML, FlagYML, FlagHelp,
		FlagStatus, FlagFormat, FlagDir,
	}

	for _, f := range flags {
		if !strings.HasPrefix(f, "--") {
			t.Errorf("flag %q should have '--' prefix", f)
		}
	}

	shortFlags := []string{
		FlagLimitShort, FlagAllShort, FlagPendingShort,
		FlagCompletedShort, FlagJSONShort, FlagHelpShort,
	}

	for _, sf := range shortFlags {
		if !strings.HasPrefix(sf, "-") || strings.HasPrefix(sf, "--") {
			t.Errorf("short flag %q should start with single '-'", sf)
		}
	}
}

func TestDirectoryConstants(t *testing.T) {
	if DirAiMemory != ".ai-memory" {
		t.Errorf("expected .ai-memory, got %q", DirAiMemory)
	}
	if DirPlans != "plans" {
		t.Errorf("expected plans, got %q", DirPlans)
	}
	if DirPending != "pending" {
		t.Errorf("expected pending, got %q", DirPending)
	}
	if DirCompleted != "completed" {
		t.Errorf("expected completed, got %q", DirCompleted)
	}
	if DirSubtasks != "subtasks" {
		t.Errorf("expected subtasks, got %q", DirSubtasks)
	}
	if ExtMarkdown != ".md" {
		t.Errorf("expected .md, got %q", ExtMarkdown)
	}
}

func TestTableColumnConstants(t *testing.T) {
	cols := []string{ColID, ColTitle, ColStatus, ColFilePath}
	for _, c := range cols {
		if len(c) == 0 {
			t.Errorf("column header should not be empty")
		}
	}
}

func TestUIAndColorConstants(t *testing.T) {
	colors := []string{ColorHeader, ColorBorder, ColorStripe}
	for _, c := range colors {
		if !strings.HasPrefix(c, "#") {
			t.Errorf("hex color %q should start with #", c)
		}
	}

	if len(UsagePlans) == 0 {
		t.Errorf("UsagePlans should not be empty")
	}
	if !strings.Contains(UsagePlans, CmdPlans) {
		t.Errorf("UsagePlans should mention %q", CmdPlans)
	}
}
