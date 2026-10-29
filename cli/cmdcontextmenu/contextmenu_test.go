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

func TestRunContextMenu_CleanPowerShell(t *testing.T) {
	mockRunner := func(args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}
	withMockReg("windows", mockRunner, func() {
		if err := RunContextMenu([]string{"clean-powershell", "--dry-run"}); err != nil {
			t.Fatalf("clean-powershell --dry-run failed: %v", err)
		}
		if err := RunContextMenu([]string{"clean-powershell", "--restore", "--dry-run"}); err != nil {
			t.Fatalf("clean-powershell --restore --dry-run failed: %v", err)
		}
		if err := RunContextMenu([]string{"clean-powershell", "--help"}); err != nil {
			t.Fatalf("clean-powershell --help failed: %v", err)
		}
	})
}

func TestRunContextMenu_Customize(t *testing.T) {
	mockRunner := func(args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}
	withMockReg("windows", mockRunner, func() {
		if err := RunContextMenu([]string{"customize"}); err != nil {
			t.Fatalf("customize failed: %v", err)
		}
		if err := RunContextMenu([]string{"customize", "--help"}); err != nil {
			t.Fatalf("customize --help failed: %v", err)
		}
	})
}

func TestHasExplainFlag(t *testing.T) {
	if !hasExplainFlag([]string{"--explain"}) {
		t.Errorf("hasExplainFlag(--explain) = false, want true")
	}

	if hasExplainFlag([]string{"install"}) {
		t.Errorf("hasExplainFlag(install) = true, want false")
	}
}

func TestHasRestoreFlag(t *testing.T) {
	if !hasRestoreFlag([]string{"--restore"}) {
		t.Errorf("hasRestoreFlag(--restore) = false, want true")
	}

	if hasRestoreFlag([]string{"--dry-run"}) {
		t.Errorf("hasRestoreFlag(--dry-run) = true, want false")
	}
}

func TestHasDryRunFlag(t *testing.T) {
	if !hasDryRunFlag([]string{"--dry-run"}) {
		t.Errorf("hasDryRunFlag(--dry-run) = false, want true")
	}

	if hasDryRunFlag([]string{"clean"}) {
		t.Errorf("hasDryRunFlag(clean) = true, want false")
	}
}

func TestHasHelpFlag(t *testing.T) {
	if !hasHelpFlag([]string{"--help"}) {
		t.Errorf("hasHelpFlag(--help) = false, want true")
	}

	if hasHelpFlag([]string{"status"}) {
		t.Errorf("hasHelpFlag(status) = true, want false")
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

func TestHasHideGitFlag(t *testing.T) {
	if !hasHideGitFlag([]string{"--hide-git"}) {
		t.Errorf("hasHideGitFlag(--hide-git) = false, want true")
	}

	if !hasHideGitFlag([]string{"--hide-git-actions"}) {
		t.Errorf("hasHideGitFlag(--hide-git-actions) = false, want true")
	}

	if hasHideGitFlag([]string{"install"}) {
		t.Errorf("hasHideGitFlag(install) = true, want false")
	}
}

func TestRunContextMenu_ToggleGitActions(t *testing.T) {
	mockRunner := func(args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}
	withMockReg("windows", mockRunner, func() {
		if err := RunContextMenu([]string{"hide-git-actions"}); err != nil {
			t.Fatalf("hide-git-actions failed: %v", err)
		}
		if err := RunContextMenu([]string{"show-git-actions"}); err != nil {
			t.Fatalf("show-git-actions failed: %v", err)
		}
		if err := RunContextMenu([]string{"toggle-git-actions"}); err != nil {
			t.Fatalf("toggle-git-actions failed: %v", err)
		}
	})
}
