package cmdcontextmenu

import (
	"testing"
)

func TestRunContextMenu_Help(t *testing.T) {
	if err := RunContextMenu([]string{"--help"}); err != nil {
		t.Fatalf("RunContextMenu(--help) failed: %v", err)
	}

	if err := RunContextMenu([]string{"-h"}); err != nil {
		t.Fatalf("RunContextMenu(-h) failed: %v", err)
	}
}

func TestRunContextMenu_Empty(t *testing.T) {
	if err := RunContextMenu([]string{}); err != nil {
		t.Fatalf("RunContextMenu() failed: %v", err)
	}
}

func TestRunContextMenu_Status(t *testing.T) {
	if err := RunContextMenu([]string{"status"}); err != nil {
		t.Fatalf("RunContextMenu(status) failed: %v", err)
	}
}

func TestRunContextMenu_Unknown(t *testing.T) {
	if err := RunContextMenu([]string{"unknown-subcommand"}); err != nil {
		t.Fatalf("RunContextMenu(unknown) failed: %v", err)
	}
}

func TestRunContextMenu_Wizard(t *testing.T) {
	if err := RunContextMenu([]string{"wizard"}); err != nil {
		t.Fatalf("RunContextMenu(wizard) failed: %v", err)
	}
}

func TestHasExplainFlag(t *testing.T) {
	if !hasExplainFlag([]string{"--explain"}) {
		t.Errorf("hasExplainFlag(--explain) = false, want true")
	}

	if hasExplainFlag([]string{"install"}) {
		t.Errorf("hasExplainFlag(install) = true, want false")
	}
}

func TestIsContextMenuHelp(t *testing.T) {
	if !isContextMenuHelp([]string{"help"}) {
		t.Errorf("isContextMenuHelp(help) = false, want true")
	}

	if isContextMenuHelp([]string{"install"}) {
		t.Errorf("isContextMenuHelp(install) = true, want false")
	}
}
