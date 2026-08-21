package cluster

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestExecInstall(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() {
		runCmdFunc = origRunCmd
	}()

	runCmdFunc = func(cmd *exec.Cmd) error {
		// Mock package manager detection
		if len(cmd.Args) > 0 && (cmd.Args[0] == "winget" || cmd.Args[0] == "brew") {
			return nil
		}
		if len(cmd.Args) > 0 && (cmd.Args[0] == "choco" || cmd.Args[0] == "apt-get") {
			return errors.New("not found")
		}

		// Mock the actual install command via shell
		cmdStr := strings.Join(cmd.Args, " ")
		if !strings.Contains(cmdStr, "install") {
			return errors.New("unexpected command: " + cmdStr)
		}
		if strings.Contains(cmdStr, "failpkg") {
			cmd.Stderr.Write([]byte("install failed"))
			return &exec.ExitError{}
		}
		return nil
	}

	packages := []string{"goodpkg", "failpkg"}
	results, err := ExecInstall(context.Background(), ClusterNode{}, packages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if !results[0].Succeeded || results[0].PackageName != "goodpkg" {
		t.Errorf("expected goodpkg to succeed, got %v", results[0])
	}
	if results[1].Succeeded || results[1].PackageName != "failpkg" {
		t.Errorf("expected failpkg to fail, got %v", results[1])
	}
	if results[1].Stderr != "install failed" {
		t.Errorf("expected stderr 'install failed', got %q", results[1].Stderr)
	}
}
