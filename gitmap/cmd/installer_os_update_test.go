package cmd

import (
	"testing"
)

func TestInstallerOSUpdate(t *testing.T) {
	if err := executeOSUpdate([]string{}, "ubuntu"); err == nil {
		t.Fatal("expected error on empty args")
	}
}
