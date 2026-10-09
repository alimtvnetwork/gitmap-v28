package cmddesktopsync

import (
	"testing"
)

func TestHasGHDesktopInstallFlag(t *testing.T) {
	hasLong := hasGHDesktopInstallFlag([]string{"--install"})
	hasShort := hasGHDesktopInstallFlag([]string{"-i"})
	hasOther := hasGHDesktopInstallFlag([]string{"--all"})

	if !hasLong || !hasShort {
		t.Errorf("expected install flag to be detected")
	}
	if hasOther {
		t.Errorf("expected no install flag for --all")
	}
}
