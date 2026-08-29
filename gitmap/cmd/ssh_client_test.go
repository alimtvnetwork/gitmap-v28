package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestInteractiveSSHClient(t *testing.T) {
	client := InteractiveSSHClient{
		Stdin:  &bytes.Buffer{},
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}

	if client.Stdin == nil || client.Stdout == nil || client.Stderr == nil {
		t.Errorf("InteractiveSSHClient fields should not be nil")
	}
}

func TestSpawnSSH(t *testing.T) {
	// We can't actually spawn SSH effectively in this unit test without a real server,
	// but we can ensure it compiles and has the correct signature.
	// Running a canceled context is a safe way to test execution framework.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	target := SSHTarget{Username: "test", IP: "127.0.0.1"}
	err := SpawnSSH(ctx, target, []string{"date"})

	if err == nil {
		t.Errorf("expected error due to canceled context")
	}
}

func TestPromptSSHPassword(t *testing.T) {
	ctx := context.Background()
	_, err := PromptSSHPassword(ctx, "Password: ", -1)
	if err == nil {
		t.Errorf("expected error when reading password from invalid fd")
	}

	// Verify error type
	if err != nil {
		if _, ok := err.(interface{ Error() string }); !ok {
			t.Errorf("expected error interface")
		}
	}
}

func TestSSHClient(t *testing.T) {
	target := SSHTarget{Username: "root", IP: "10.0.0.1"}
	expected := []string{"root@10.0.0.1", "ls", "-la"}

	ctx := context.Background()
	dir := t.TempDir()

	var scriptPath string
	var scriptContent string
	argsFile := filepath.Join(dir, "args.txt")

	// Escape path for Windows batch
	escArgsFile := argsFile
	// On Windows we create a .bat file, on Unix a shell script
	if os.PathSeparator == '\\' {
		scriptPath = filepath.Join(dir, "ssh.bat")
		scriptContent = fmt.Sprintf(`@echo %%* > "%s"`, escArgsFile)
	} else {
		scriptPath = filepath.Join(dir, "ssh")
		scriptContent = fmt.Sprintf(`#!/bin/sh
echo "$@" > "%s"`, escArgsFile)
	}

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write mock ssh script: %v", err)
	}

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	args := expected[1:] // first element is target
	if err := SpawnSSH(ctx, target, args); err != nil {
		t.Fatalf("SpawnSSH failed: %v", err)
	}

	out, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("failed to read args file: %v", err)
	}

	outputArgs := string(bytes.TrimSpace(out))
	expectedArgsStr := ""
	for i, arg := range expected {
		if i > 0 {
			expectedArgsStr += " "
		}
		expectedArgsStr += arg
	}

	if outputArgs != expectedArgsStr {
		t.Errorf("expected args %q, got %q", expectedArgsStr, outputArgs)
	}
}
