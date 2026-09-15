package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestHasGHDesktopInstallFlag(t *testing.T) {
	hasLong := hasGHDesktopInstallFlag([]string{"--install"})
	hasShort := hasGHDesktopInstallFlag([]string{"-i"})
	hasNone := hasGHDesktopInstallFlag([]string{"--all"})

	if !hasLong || !hasShort {
		t.Errorf("expected install flag to be detected")
	}
	if hasNone {
		t.Errorf("expected no install flag for --all")
	}
}

func TestIsGHDesktopSubcommand(t *testing.T) {
	isOpt := isGHDesktopSubcommand([]string{"optimize"})
	isClean := isGHDesktopSubcommand([]string{"clean"})
	isGroup := isGHDesktopSubcommand([]string{"group"})
	isDirect := isGHDesktopSubcommand([]string{"/some/path"})

	if !isOpt || !isClean || !isGroup {
		t.Errorf("expected subcommand to be detected")
	}
	if isDirect {
		t.Errorf("direct path should not be recognized as subcommand")
	}
}

func TestResolveGHDesktopTarget_CwdFallback(t *testing.T) {
	cwd := "D:/repos/myrepo"
	targetEmpty := resolveGHDesktopTarget(cwd, []string{})
	targetWithFlag := resolveGHDesktopTarget(cwd, []string{"--install"})

	if targetEmpty != cwd || targetWithFlag != cwd {
		t.Errorf("expected target to fallback to cwd, got %q and %q", targetEmpty, targetWithFlag)
	}
}

func TestRegisterGHDesktop_MissingCLIReturnsValidationError(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("LOCALAPPDATA", tmp)
	t.Setenv("HOME", tmp)
	t.Setenv("PATH", tmp)

	err := registerGHDesktop("nonexistent/path")
	if err == nil {
		t.Fatal("expected error on missing CLI")
	}

	appErr, isAppError := err.(*apperror.AppError)
	if !isAppError {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}

	if appErr.Type != apperror.ErrorTypeValidation {
		t.Errorf("expected validation error type, got %v", appErr.Type)
	}
	if appErr.Message != constants.MsgDesktopNotFound {
		t.Errorf("expected message %q, got %q", constants.MsgDesktopNotFound, appErr.Message)
	}
}
