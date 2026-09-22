package cmdos

import (
	"testing"
)

func TestThemeModeConstants(t *testing.T) {
	if ThemeModeDark != "dark" {
		t.Errorf("expected dark, got %s", ThemeModeDark)
	}

	if ThemeModeLight != "light" {
		t.Errorf("expected light, got %s", ThemeModeLight)
	}
}

func TestRunOSThemeRouting(t *testing.T) {
	if err := runOSThemeCommand([]string{"help"}); err != nil {
		t.Errorf("expected nil error for help, got: %v", err)
	}

	err := runOSThemeCommand([]string{"invalid-theme-xyz"})
	if err == nil {
		t.Errorf("expected error for invalid theme mode, got nil")
	}
}
