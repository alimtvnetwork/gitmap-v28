// Package cmdinstall — installer_types_test.go tests OS matching, chaining sorting, and validation.
package cmdinstall

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestMatchesTargetOS(t *testing.T) {
	tests := []struct {
		target    InstallerOSType
		currentOS string
		wantMatch bool
	}{
		{InstallerOSAll, "windows", true},
		{InstallerOSAll, "linux", true},
		{InstallerOSWindows, constants.PlatformWindows, true},
		{InstallerOSWindows, constants.PlatformLinux, false},
		{InstallerOSUnix, constants.PlatformLinux, true},
		{InstallerOSUnix, constants.PlatformDarwin, true},
		{InstallerOSUnix, constants.PlatformWindows, false},
		{InstallerOSUbuntu, constants.PlatformLinux, true},
		{InstallerOSCentOS, constants.PlatformLinux, true},
		{"", "windows", true},
	}

	for _, tt := range tests {
		got := MatchesTargetOS(tt.target, tt.currentOS)
		if got != tt.wantMatch {
			t.Errorf("MatchesTargetOS(%q, %q) = %v, want %v", tt.target, tt.currentOS, got, tt.wantMatch)
		}
	}
}

func TestOrderSubInstallers(t *testing.T) {
	chain := []ChainedSubInstaller{
		{SubInstallerId: "step3", SequenceOrder: 3, Command: "echo 3"},
		{SubInstallerId: "step1", SequenceOrder: 1, Command: "echo 1"},
		{SubInstallerId: "step2", SequenceOrder: 2, Command: "echo 2"},
	}

	ordered := OrderSubInstallers(chain)
	if len(ordered) != 3 {
		t.Fatalf("expected 3 items, got %d", len(ordered))
	}
	if ordered[0].SubInstallerId != "step1" || ordered[1].SubInstallerId != "step2" || ordered[2].SubInstallerId != "step3" {
		t.Errorf("chain not ordered properly: %v", ordered)
	}
}

func TestValidateSubInstallerChain_Valid(t *testing.T) {
	chain := []ChainedSubInstaller{
		{SubInstallerId: "step1", Command: "echo 1"},
		{SubInstallerId: "step2", Command: "echo 2"},
	}

	if err := ValidateSubInstallerChain(chain); err != nil {
		t.Fatalf("expected valid chain, got error: %v", err)
	}
}

func TestValidateSubInstallerChain_Errors(t *testing.T) {
	// Empty SubInstallerId
	emptyIdChain := []ChainedSubInstaller{
		{SubInstallerId: "", Command: "echo 1"},
	}
	if err := ValidateSubInstallerChain(emptyIdChain); err == nil {
		t.Error("expected error on empty subInstallerId, got nil")
	}

	// Duplicate SubInstallerId
	dupChain := []ChainedSubInstaller{
		{SubInstallerId: "step1", Command: "echo 1"},
		{SubInstallerId: "step1", Command: "echo 2"},
	}
	if err := ValidateSubInstallerChain(dupChain); err == nil {
		t.Error("expected error on duplicate subInstallerId, got nil")
	}

	// Empty Command
	emptyCmdChain := []ChainedSubInstaller{
		{SubInstallerId: "step1", Command: "   "},
	}
	if err := ValidateSubInstallerChain(emptyCmdChain); err == nil {
		t.Error("expected error on empty command, got nil")
	}
}
