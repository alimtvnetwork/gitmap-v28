package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAGYEnhancements_E2E(t *testing.T) {
	if os.Getenv("HEAVY_E2E") == "" {
		t.Skip("Skipping heavy e2e test, set HEAVY_E2E=1 to run")
	}

	tempDir := t.TempDir()

	// 1. Test create-repo --common
	testRepoDir := filepath.Join(tempDir, "existing-folder")
	os.MkdirAll(testRepoDir, 0755)

	cmdCreate := exec.Command("gitmap", "create-repo", testRepoDir, "--common")
	cmdCreate.Dir = tempDir
	if err := cmdCreate.Run(); err != nil {
		t.Fatalf("Failed to run create-repo --common: %v", err)
	}

	// 2. Test agy conv ls
	outJson := filepath.Join(tempDir, "out.json")
	cmdLs := exec.Command("gitmap", "agy", "conv", "ls", "5", "-f", outJson)
	cmdLs.Dir = testRepoDir
	if err := cmdLs.Run(); err != nil {
		t.Fatalf("Failed to run agy conv ls: %v", err)
	}

	if _, err := os.Stat(outJson); os.IsNotExist(err) {
		t.Fatalf("Expected %s to exist", outJson)
	}

	// 3. Test agy ipt
	cmdIpt := exec.Command("gitmap", "agy", "ipt", "Hello world", "-p", "test")
	cmdIpt.Dir = testRepoDir
	if err := cmdIpt.Run(); err != nil {
		t.Fatalf("Failed to run agy ipt: %v", err)
	}
}
