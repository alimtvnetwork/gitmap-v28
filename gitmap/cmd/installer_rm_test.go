package cmd

import (
	"testing"
)

func TestInstallerRmCmd(t *testing.T) {
	if err := executeInstallerRm([]string{}); err == nil {
		t.Fatal("expected error on empty args")
	}
}
