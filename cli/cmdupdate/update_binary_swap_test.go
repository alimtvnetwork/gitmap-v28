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

	if !MatchPlatformAsset("gitmap-linux-amd64.tar.gz", "v6.590.0", "linux", "amd64") {
		t.Error("expected gitmap-linux-amd64.tar.gz to match linux/amd64")
	}

	if !MatchPlatformAsset("gitmap-darwin-arm64.tar.gz", "v6.590.0", "darwin", "arm64") {
		t.Error("expected gitmap-darwin-arm64.tar.gz to match darwin/arm64")
	}

	if MatchPlatformAsset("cli-updater-v6.593.1-windows-amd64.exe", "v6.593.1", "windows", "amd64") {
		t.Error("cli-updater must NEVER match windows/amd64 platform asset")
	}

	if MatchPlatformAsset("cli-updater-v6.593.1-linux-amd64", "v6.593.1", "linux", "amd64") {
		t.Error("cli-updater must NEVER match linux/amd64 platform asset")
	}
}

func TestFindExecutableAsset_PrefersGitmapZipOverCliUpdaterExe(t *testing.T) {
	assets := []GitHubReleaseAsset{
		{Name: "cli-updater-v6.593.1-windows-amd64.exe", Size: 6200000},
		{Name: "gitmap-v6.593.1-windows-amd64.zip", Size: 18000000},
	}
	asset, found := FindExecutableAsset(assets, "v6.593.1", "windows", "amd64")
	if !found {
		t.Fatal("expected to find executable asset")
	}

	if asset.Name != "gitmap-v6.593.1-windows-amd64.zip" {
		t.Fatalf("expected gitmap-v6.593.1-windows-amd64.zip, got %s", asset.Name)
	}
}

func TestIsDirectBinaryAsset(t *testing.T) {
	if !isDirectBinaryAsset("gitmap.exe", "windows") {
		t.Error("expected gitmap.exe to be detected as direct binary")
	}

	if isDirectBinaryAsset("cli-updater-v6.593.1-windows-amd64.exe", "windows") {
		t.Error("cli-updater must NEVER be detected as direct binary")
	}

	if isDirectBinaryAsset("gitmap.zip", "windows") {
		t.Error("expected gitmap.zip not to be direct binary")
	}

	if isDirectBinaryAsset("gitmap.tar.gz", "linux") {
		t.Error("expected gitmap.tar.gz not to be direct binary")
	}

	if !isDirectBinaryAsset("gitmap-linux-amd64", "linux") {
		t.Error("expected gitmap-linux-amd64 to be detected as direct binary")
	}
}

func TestExtractDirectExecutable(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "dummy-binary.exe")
	if err := os.WriteFile(srcFile, []byte("echo binary"), 0755); err != nil {
		t.Fatalf("failed to create dummy binary: %v", err)
	}

	outDir := filepath.Join(tmpDir, "out")
	_ = os.MkdirAll(outDir, 0755)

	dst, err := extractDirectExecutable(srcFile, outDir)
	if err != nil {
		t.Fatalf("unexpected copy error: %v", err)
	}

	if _, statErr := os.Stat(dst); statErr != nil {
		t.Fatalf("destination binary does not exist: %v", statErr)
	}
}
