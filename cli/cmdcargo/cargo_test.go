package cmdcargo

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestExtractInstallFlag(t *testing.T) {
	cases := []struct {
		args        []string
		wantHas     bool
		wantCleaned []string
	}{
		{
			args:        []string{"build", "--release"},
			wantHas:     false,
			wantCleaned: []string{"build", "--release"},
		},
		{
			args:        []string{"--install", "fmt", "--", "--check"},
			wantHas:     true,
			wantCleaned: []string{"fmt", "--", "--check"},
		},
		{
			args:        []string{"-i", "check"},
			wantHas:     true,
			wantCleaned: []string{"check"},
		},
	}

	for _, tc := range cases {
		cleaned, has := extractInstallFlag(tc.args)
		if has != tc.wantHas {
			t.Errorf("extractInstallFlag(%v) has = %v, want %v", tc.args, has, tc.wantHas)
		}
		if len(cleaned) != len(tc.wantCleaned) {
			t.Errorf("extractInstallFlag(%v) cleaned len = %d, want %d", tc.args, len(cleaned), len(tc.wantCleaned))
		}
	}
}

func TestFormatCargoInstallSuggestions(t *testing.T) {
	for _, targetOS := range []string{constants.PlatformWindows, constants.PlatformDarwin, constants.PlatformLinux} {
		out := FormatCargoInstallSuggestions(targetOS)
		if !strings.Contains(out, "gitmap install cargo") {
			t.Errorf("expected %q to contain 'gitmap install cargo' for OS %s", out, targetOS)
		}
		if !strings.Contains(out, "rustup.rs") {
			t.Errorf("expected %q to contain rustup link for OS %s", out, targetOS)
		}
	}
}

func TestNewMissingCargoError(t *testing.T) {
	err := NewMissingCargoError()
	if err.Code != "E7100" {
		t.Errorf("expected error code E7100, got %s", err.Code)
	}
}

func TestSubcommandMatchers(t *testing.T) {
	if !isStatusSubcommand("status") || !isStatusSubcommand("st") || !isStatusSubcommand("info") {
		t.Errorf("expected status subcommands to match")
	}
	if !isInstallSubcommand("install") || !isInstallSubcommand("in") {
		t.Errorf("expected install subcommands to match")
	}
}
