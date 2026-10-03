package cmdssh

import (
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func TestBuildWindowsCheckCapabilityCmd(t *testing.T) {
	cmd := buildWindowsCheckCapabilityCmd()

	if !strings.Contains(cmd, "OpenSSH.Server") {
		t.Errorf("expected cmd to check OpenSSH.Server, got: %s", cmd)
	}

	if !strings.Contains(cmd, "Get-Service -Name sshd") {
		t.Errorf("expected cmd to check sshd service fallback, got: %s", cmd)
	}
}

func TestBuildWindowsAddCapabilityCmd(t *testing.T) {
	cmd := buildWindowsAddCapabilityCmd()

	if !strings.Contains(cmd, "wuauserv") {
		t.Errorf("expected cmd to handle wuauserv service, got: %s", cmd)
	}

	if !strings.Contains(cmd, "Add-WindowsCapability") {
		t.Errorf("expected cmd to invoke Add-WindowsCapability, got: %s", cmd)
	}
}

func TestParseSSHEnableFlags_Defaults(t *testing.T) {
	opts, err := parseSSHEnableFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.port != 22 {
		t.Errorf("expected default port 22, got %d", opts.port)
	}

	if opts.isForce {
		t.Errorf("expected isForce to be false by default")
	}
}

func TestParseSSHEnableFlags_Custom(t *testing.T) {
	opts, err := parseSSHEnableFlags([]string{"--port", "2222", "--force"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.port != 2222 {
		t.Errorf("expected port 2222, got %d", opts.port)
	}

	if !opts.isForce {
		t.Errorf("expected isForce to be true")
	}
}

func TestHandleCapabilityError_Force(t *testing.T) {
	err := handleCapabilityError(errors.New("test-error"), true)
	if err != nil {
		t.Errorf("expected nil error when isForce is true, got: %v", err)
	}
}

func TestHandleCapabilityError_NoForce(t *testing.T) {
	err := handleCapabilityError(errors.New("test-error"), false)
	if err == nil {
		t.Fatalf("expected error when isForce is false")
	}

	appErr, isApp := err.(*apperror.AppError)
	if !isApp {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}

	if appErr.Type != apperror.ErrorTypeAbort {
		t.Errorf("expected ErrorTypeAbort, got %s", appErr.Type)
	}
}

func TestHandleSSHDServiceError(t *testing.T) {
	appErr := handleSSHDServiceError(errors.New("access denied"), "Start-Service sshd")
	if appErr == nil {
		t.Fatalf("expected non-nil error")
	}

	if appErr.Type != apperror.ErrorTypeAbort {
		t.Errorf("expected ErrorTypeAbort, got %s", appErr.Type)
	}
}
