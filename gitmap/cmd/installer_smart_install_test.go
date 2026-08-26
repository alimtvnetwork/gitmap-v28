package cmd

import (
	"testing"
)

func TestInstallerSmartInstall(t *testing.T) {
	if err := executeSmartInstall([]string{}); err == nil {
		t.Fatal("expected error on empty args")
	}

	osTarget := detectCurrentHostOSTarget()
	if osTarget == "" {
		t.Fatal("expected non-empty host OS target")
	}
}
