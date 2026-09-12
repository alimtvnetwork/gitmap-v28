package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunOSUsage(t *testing.T) {
	if err := runOS(nil); err != nil {
		t.Fatalf("runOS(nil) error: %v", err)
	}

	if err := runOS([]string{}); err != nil {
		t.Fatalf("runOS([]string{}) error: %v", err)
	}
}

func TestRunOSStatus(t *testing.T) {
	if err := runOSStatus(nil); err != nil {
		t.Fatalf("runOSStatus error: %v", err)
	}
}

func TestRunOSUnknownSubcmd(t *testing.T) {
	err := runOS([]string{"unknownxyz"})
	if err == nil {
		t.Fatalf("expected error for unknown subcommand, got nil")
	}
}

func TestFixLinkWithTarget(t *testing.T) {
	dir := t.TempDir()
	linkPath := filepath.Join(dir, "my_link")
	targetPath := filepath.Join(dir, "target.txt")
	_ = os.WriteFile(targetPath, []byte("hello"), 0644)

	danglingTarget := filepath.Join(dir, "missing.txt")
	if err := os.Symlink(danglingTarget, linkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	args := []string{linkPath, "--target", targetPath}
	if err := runOSFixLink(args); err != nil {
		t.Fatalf("runOSFixLink failed: %v", err)
	}

	assertLinkTarget(t, linkPath, targetPath)
}

func TestFixLinkDryRun(t *testing.T) {
	dir := t.TempDir()
	linkPath := filepath.Join(dir, "test_link")
	dangling := filepath.Join(dir, "missing.txt")
	if err := os.Symlink(dangling, linkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	args := []string{linkPath, "--target", dir, "--dry-run"}
	if err := runOSFixLink(args); err != nil {
		t.Fatalf("runOSFixLink dry run error: %v", err)
	}

	assertLinkTarget(t, linkPath, dangling)
}

func TestFixLinkDirectory(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "valid_target.txt")
	_ = os.WriteFile(targetPath, []byte("data"), 0644)

	validLink := filepath.Join(dir, "valid_link")
	if err := os.Symlink(targetPath, validLink); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	args := []string{dir, "--recursive"}
	if err := runOSFixLink(args); err != nil {
		t.Fatalf("directory runOSFixLink error: %v", err)
	}
}

func TestFixLinkJSON(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "file.txt")
	_ = os.WriteFile(targetPath, []byte("data"), 0644)

	linkPath := filepath.Join(dir, "json_link")
	if err := os.Symlink(targetPath, linkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	args := []string{linkPath, "--json"}
	if err := runOSFixLink(args); err != nil {
		t.Fatalf("runOSFixLink --json error: %v", err)
	}
}

func assertLinkTarget(t *testing.T, linkPath, expectedTarget string) {
	t.Helper()
	actual, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("readlink error: %v", err)
	}

	if actual != expectedTarget {
		t.Fatalf("link target = %q, want %q", actual, expectedTarget)
	}
}
