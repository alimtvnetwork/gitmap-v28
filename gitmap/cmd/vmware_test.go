package cmd

import (
	"errors"
	"runtime"
	"strings"
	"testing"
)

func TestVmwareSubcommandDispatch(t *testing.T) {
	if err := runVmware([]string{}); err != nil {
		t.Errorf("Expected nil error for empty args, got %v", err)
	}

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

func TestVmwareInstallDryRun(t *testing.T) {
	if err := runVmware([]string{"install", "--dry-run"}); err != nil {
		t.Errorf("Expected nil error for install --dry-run, got %v", err)
	}

	if err := runVmware([]string{"in", "-n"}); err != nil {
		t.Errorf("Expected nil error for in -n, got %v", err)
	}
}

func TestVmwareSharedEnableDryRun(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	if err := runVmware([]string{"shared", "enable", "--dry-run"}); err != nil {
		t.Errorf("Expected nil error for shared enable --dry-run, got %v", err)
	}
}

func TestFormatVmwareMountErrorDiagnostics(t *testing.T) {
	fakeErr := errors.New("exit status 149")
	out := []byte("mount failed: Error -107 cannot open connection!")
	formatted := formatVmwareMountError(out, fakeErr)

	if !strings.Contains(formatted, "Diagnostic & Remediation") {
		t.Errorf("Expected diagnostic section in formatted error, got %s", formatted)
	}

	if !strings.Contains(formatted, "Virtual Machine Settings") && !strings.Contains(formatted, "Shared Folders") {
		t.Errorf("Expected Shared Folders instructions in formatted error, got %s", formatted)
	}
}
