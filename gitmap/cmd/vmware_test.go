package cmd

import (
	"runtime"
	"testing"
)

func TestVmwareSubcommandDispatch(t *testing.T) {
	if err := runVmware([]string{}); err != nil {
		t.Errorf("Expected nil error for empty args, got %v", err)
	}

	// Test unknown subcommand returns E4001
	err := runVmware([]string{"nonexistent-command"})
	if err == nil {
		t.Errorf("Expected error for unknown subcommand, got nil")
	}
}

func TestVmwareStatus(t *testing.T) {
	if err := runVmwareStatus(nil); err != nil {
		t.Errorf("Expected runVmwareStatus to succeed, got %v", err)
	}
}

func TestVmwareOSConstraint(t *testing.T) {
	if runtime.GOOS == "linux" {
		return
	}

	err := checkVMwarePrerequisites()
	if err == nil {
		t.Errorf("Expected error on non-Linux OS, got nil")
	}
}
