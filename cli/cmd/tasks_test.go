// Package cmd — tasks_test.go verifies the tasks command routing and subcommands.
package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestIsLegacyTaskSubcommand(t *testing.T) {
	legacy := []string{"create", "run", "show", "delete"}
	for _, sub := range legacy {
		if !isLegacyTaskSubcommand(sub) {
			t.Errorf("isLegacyTaskSubcommand(%q) = false; want true", sub)
		}
	}

	modern := []string{"list", "history", "undo", "redo", "clear"}
	for _, sub := range modern {
		if isLegacyTaskSubcommand(sub) {
			t.Errorf("isLegacyTaskSubcommand(%q) = true; want false", sub)
		}
	}
}

func TestDispatchTaskSubcommandConstants(t *testing.T) {
	if constants.CmdTasks != "tasks" {
		t.Errorf("CmdTasks = %q; want 'tasks'", constants.CmdTasks)
	}

	if constants.CmdTaskHistory != "history" {
		t.Errorf("CmdTaskHistory = %q; want 'history'", constants.CmdTaskHistory)
	}

	if constants.CmdTaskUndo != "undo" {
		t.Errorf("CmdTaskUndo = %q; want 'undo'", constants.CmdTaskUndo)
	}

	if constants.CmdTaskRedo != "redo" {
		t.Errorf("CmdTaskRedo = %q; want 'redo'", constants.CmdTaskRedo)
	}
}
