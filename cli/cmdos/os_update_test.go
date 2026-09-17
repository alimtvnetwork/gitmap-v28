package cmdos

import (
	"context"
	"os/exec"
	"testing"
)

func setupMockOSCommandRunner() (func() *exec.Cmd, func()) {
	var capturedCmd *exec.Cmd
	orig := defaultOSCommandRunner
	defaultOSCommandRunner = func(cmd *exec.Cmd) error {
		capturedCmd = cmd
		return nil
	}
	getCaptured := func() *exec.Cmd { return capturedCmd }
	cleanup := func() { defaultOSCommandRunner = orig }
	return getCaptured, cleanup
}

func TestExecuteOSUpdate_Mocked(t *testing.T) {
	getCaptured, cleanup := setupMockOSCommandRunner()
	defer cleanup()

	if err := ExecuteOSUpdate(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if getCaptured() == nil {
		t.Fatal("expected non-nil captured cmd")
	}
}

func TestExecuteOSFullUpgrade_Mocked(t *testing.T) {
	getCaptured, cleanup := setupMockOSCommandRunner()
	defer cleanup()

	if err := ExecuteOSFullUpgrade(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if getCaptured() == nil {
		t.Fatal("expected non-nil captured cmd")
	}
}
