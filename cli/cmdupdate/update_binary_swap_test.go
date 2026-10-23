package cmdupdate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMatchPlatformAsset_DirectBinary(t *testing.T) {
	if !MatchPlatformAsset("gitmap-windows-amd64.exe", "v6.590.0", "windows", "amd64") {
		t.Error("expected gitmap-windows-amd64.exe to match windows/amd64")
	}

	if !MatchPlatformAsset("gitmap-linux-amd64", "v6.590.0", "linux", "amd64") {
		t.Error("expected gitmap-linux-amd64 to match linux/amd64")
	}

	if !MatchPlatformAsset("gitmap-darwin-arm64", "v6.590.0", "darwin", "arm64") {
		t.Error("expected gitmap-darwin-arm64 to match darwin/arm64")
	}
}

func TestIsDirectBinaryAsset(t *testing.T) {
	if !isDirectBinaryAsset("gitmap.exe") {
		t.Error("expected gitmap.exe to be detected as direct binary")
	}

	if isDirectBinaryAsset("gitmap.zip") {
		t.Error("expected gitmap.zip not to be direct binary")
	}

	if isDirectBinaryAsset("gitmap.tar.gz") {
		t.Error("expected gitmap.tar.gz not to be direct binary")
	}

	if !isDirectBinaryAsset("gitmap-linux-amd64") {
		t.Error("expected gitmap-linux-amd64 to be detected as direct binary")
	}
}

func TestCopyDirectBinary(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "dummy-binary.exe")
	if err := os.WriteFile(srcFile, []byte("echo binary"), 0755); err != nil {
		t.Fatalf("failed to create dummy binary: %v", err)
	}

	outDir := filepath.Join(tmpDir, "out")
	_ = os.MkdirAll(outDir, 0755)

	dst, err := copyDirectBinary(srcFile, outDir)
	if err != nil {
		t.Fatalf("unexpected copy error: %v", err)
	}

	if _, statErr := os.Stat(dst); statErr != nil {
		t.Fatalf("destination binary does not exist: %v", statErr)
	}
}
