// Package installer — runner_test.go tests shell command execution.
package installer

import (
	"testing"
)

func TestRunner(t *testing.T) {
	if errEmpty := RunInstallerCommand(""); errEmpty == nil {
		t.Fatal("expected error on empty command")
	}

	if errRun := RunInstallerCommand("echo test_runner"); errRun != nil {
		t.Fatalf("RunInstallerCommand failed: %v", errRun)
	}
}
