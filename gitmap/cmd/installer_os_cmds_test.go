package cmd

import (
	"testing"
)

func TestInstallerOSCmds(t *testing.T) {
	if err := executeOSInstall([]string{}, "ubuntu"); err == nil {
		t.Fatal("expected error on empty args")
	}
}
