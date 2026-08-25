package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createMockSSHScript(t *testing.T) (string, string) {
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args.txt")

	var scriptPath string
	var scriptContent string
	if os.PathSeparator == '\\' {
		scriptPath = filepath.Join(dir, "ssh.bat")
		scriptContent = fmt.Sprintf("@echo %%* > \"%s\"", argsFile)
	} else {
		scriptPath = filepath.Join(dir, "ssh")
		scriptContent = fmt.Sprintf("#!/bin/sh\necho \"$@\" > \"%s\"", argsFile)
	}

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write mock ssh script: %v", err)
	}

	return dir, argsFile
}

func TestgetLocalPublicKey(t *testing.T) {
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "id_rsa.pub")

	// Missing key
	_, err := getLocalPublicKey(context.Background(), keyPath, false)
	if err == nil {
		t.Errorf("expected error for missing key")
	}

	// Create key
	expectedKey := "ssh-rsa AAAAB3NzaC1yc2E... test@test"
	err = os.WriteFile(keyPath, []byte(expectedKey+"\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write key: %v", err)
	}

	key, err := getLocalPublicKey(context.Background(), keyPath, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if key != expectedKey {
		t.Errorf("unexpected key content: %v", key)
	}
}

func TestappendKeyRemote(t *testing.T) {
	// Test cancelled context
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	target := SSHTarget{Username: "test", IP: "127.0.0.1", Port: 22}
	err := appendKeyRemote(canceledCtx, "ssh-rsa ABC", target)
	if err == nil {
		t.Errorf("expected error with canceled context")
	}

	// Test with mock SSH
	dir, _ := createMockSSHScript(t)
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	err = appendKeyRemote(context.Background(), "ssh-rsa ABC", target)
	if err != nil {
		t.Errorf("expected no error with mock ssh, got %v", err)
	}
}

func TestRunSJAddAuth(t *testing.T) {
	ctx := context.Background()

	// Missing args
	err := runSJAddAuth(nil, []string{}, ctx)
	if err == nil {
		t.Errorf("expected error for missing args")
	}

	// Invalid target format
	err = runSJAddAuth(nil, []string{"@invalid@target@here"}, ctx)
	if err == nil {
		t.Errorf("expected error for invalid target format")
	}

	// Valid target with mock SSH
	dir, _ := createMockSSHScript(t)
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	// Create a dummy local key in a temp location
	tmpDir := t.TempDir()
	keyFile := filepath.Join(tmpDir, "id_rsa.pub")
	_ = os.WriteFile(keyFile, []byte("ssh-rsa DUMMY test\n"), 0644)
	oldHome := os.Getenv("USERPROFILE")
	if oldHome == "" {
		oldHome = os.Getenv("HOME")
	}
	sshDir := filepath.Join(tmpDir, ".ssh")
	_ = os.MkdirAll(sshDir, 0755)
	_ = os.WriteFile(filepath.Join(sshDir, "id_rsa.pub"), []byte("ssh-rsa DUMMY test\n"), 0644)
	os.Setenv("USERPROFILE", tmpDir)
	os.Setenv("HOME", tmpDir)
	defer func() {
		os.Setenv("USERPROFILE", oldHome)
		os.Setenv("HOME", oldHome)
	}()

	err = runSJAddAuth(nil, []string{"test@127.0.0.1"}, ctx)
	if err != nil {
		t.Errorf("expected no error with mock ssh, got %v", err)
	}
}


