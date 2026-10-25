package cmdinstall

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestRunInstallWizard(t *testing.T) {
	err := RunInstallWizard()
	if err != nil {
		t.Fatalf("RunInstallWizard() failed: %v", err)
	}
}

func TestIsContextMenuActive(t *testing.T) {
	_ = IsContextMenuActive()
}

func TestContextMenuAliases(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"context-menu", constants.ToolCtx},
		{"contextmenu", constants.ToolCtx},
		{"ctx", constants.ToolCtx},
	}

	for _, tt := range tests {
		got := ResolveToolAlias(tt.input)
		if got != tt.want {
			t.Errorf("ResolveToolAlias(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestInstallWizardCommand(t *testing.T) {
	err := RunInstall([]string{"--wizard"})
	if err != nil {
		t.Fatalf("RunInstall(--wizard) failed: %v", err)
	}
}
