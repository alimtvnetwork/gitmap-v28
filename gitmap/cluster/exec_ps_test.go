package cluster

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestExecPS_Windows_Pwsh(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows only test")
	}

	origRunCmd := runCmdFunc
	origLookPath := lookPathFuncVar
	defer func() {
		runCmdFunc = origRunCmd
		lookPathFuncVar = origLookPath
	}()

	lookPathFuncVar = func(file string) (string, error) {
		if file == "pwsh" {
			return "C:\\pwsh.exe", nil
		}
		return "", errors.New("not found")
	}

	runCmdFunc = func(cmd *exec.Cmd) error {
		if !strings.Contains(cmd.Path, "pwsh.exe") {
			t.Errorf("expected pwsh.exe, got %s", cmd.Path)
		}
		cmd.Stdout.Write([]byte("ok"))
		return nil
	}

	out, errOut, code, err := ExecPS(context.Background(), ClusterNode{}, "echo hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Errorf("expected code 0, got %d", code)
	}
	if out != "ok" {
		t.Errorf("expected out 'ok', got %q", out)
	}
	if errOut != "" {
		t.Errorf("expected empty stderr, got %q", errOut)
	}
}

func TestExecPS_Windows_PowershellFallback(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows only test")
	}

	origRunCmd := runCmdFunc
	origLookPath := lookPathFuncVar
	defer func() {
		runCmdFunc = origRunCmd
		lookPathFuncVar = origLookPath
	}()

	lookPathFuncVar = func(file string) (string, error) {
		if file == "pwsh" {
			return "", errors.New("pwsh not found")
		}
		if file == "powershell" {
			return "C:\\powershell.exe", nil
		}
		return "", errors.New("not found")
	}

	runCmdFunc = func(cmd *exec.Cmd) error {
		if !strings.Contains(cmd.Path, "powershell.exe") {
			t.Errorf("expected powershell.exe, got %s", cmd.Path)
		}
		return nil
	}

	_, _, _, err := ExecPS(context.Background(), ClusterNode{}, "echo hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecPS_Unix_Pwsh(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix only test")
	}

	origRunCmd := runCmdFunc
	origLookPath := lookPathFuncVar
	defer func() {
		runCmdFunc = origRunCmd
		lookPathFuncVar = origLookPath
	}()

	lookPathFuncVar = func(file string) (string, error) {
		if file == "pwsh" {
			return "/bin/pwsh", nil
		}
		return "", errors.New("not found")
	}

	runCmdFunc = func(cmd *exec.Cmd) error {
		if !strings.Contains(cmd.Path, "pwsh") {
			t.Errorf("expected pwsh, got %s", cmd.Path)
		}
		return nil
	}

	_, _, _, err := ExecPS(context.Background(), ClusterNode{}, "echo hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecPS_Unix_PwshNotFound(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix only test")
	}

	origRunCmd := runCmdFunc
	origLookPath := lookPathFuncVar
	defer func() {
		runCmdFunc = origRunCmd
		lookPathFuncVar = origLookPath
	}()

	lookPathFuncVar = func(file string) (string, error) {
		return "", errors.New("pwsh not found")
	}

	out, errOut, code, err := ExecPS(context.Background(), ClusterNode{}, "echo hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty out, got %q", out)
	}
	if errOut != "pwsh not found, skipping" {
		t.Errorf("unexpected stderr: %q", errOut)
	}
	if code != 0 {
		t.Errorf("expected code 0, got %d", code)
	}
}
