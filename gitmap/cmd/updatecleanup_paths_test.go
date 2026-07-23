package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeriveDeployAppDir(t *testing.T) {
	tests := []struct {
		name     string
		selfPath string
		want     string
	}{
		{name: "path binary outside gitmap-cli dir", selfPath: "E:/bin-run/gitmap.exe", want: "E:/gitmap-cli"},
		{name: "path binary already in gitmap-cli dir", selfPath: "E:/gitmap-cli/gitmap.exe", want: "E:/gitmap-cli"},
		{name: "path binary already in legacy gitmap dir", selfPath: "E:/gitmap/gitmap.exe", want: "E:/gitmap"},
		{name: "empty path", selfPath: "", want: ""},
	}

	for _, tt := range tests {
		got := deriveDeployAppDir(tt.selfPath)
		want := tt.want
		if len(want) > 0 {
			want = filepath.Clean(want)
		}
		if got != want {
			t.Fatalf("deriveDeployAppDir(%q) = %q, want %q", tt.selfPath, got, tt.want)
		}
	}
}

func TestResolveBuildOutputDir(t *testing.T) {
	repoPath := "/repo"
	// Normalize to slash form so the assertion holds on Windows
	// (filepath.Clean uses backslash there, but the test contract
	// is "the right logical path" not "the OS's separator").
	norm := filepath.ToSlash
	if got := norm(resolveBuildOutputDir(repoPath, "./bin")); got != "/repo/bin" {
		t.Fatalf("resolveBuildOutputDir relative = %q, want %q", got, "/repo/bin")
	}
	if got := norm(resolveBuildOutputDir(repoPath, "")); got != "/repo/bin" {
		t.Fatalf("resolveBuildOutputDir default = %q, want %q", got, "/repo/bin")
	}
	if got := norm(resolveBuildOutputDir(repoPath, "/custom/bin")); got != "/custom/bin" {
		t.Fatalf("resolveBuildOutputDir absolute = %q, want %q", got, "/custom/bin")
	}
}

func TestCollectBackupCleanupDirsIncludesPathDerivedDeployAndBuild(t *testing.T) {
	config := updateCleanupConfig{
		DeployPath:  "E:/bin-run",
		BuildOutput: "./bin",
	}

	dirs := collectBackupCleanupDirs("E:/bin-run/gitmap.exe", "/repo", config)
	assertHasCleanupDir(t, dirs, "E:/bin-run")
	assertHasCleanupDir(t, dirs, "E:/gitmap-cli")
	assertHasCleanupDir(t, dirs, "E:/bin-run/gitmap-cli")
	assertHasCleanupDir(t, dirs, "/repo/bin")
}

func TestCollectTempCleanupDirsIncludesTempAndDerivedTargets(t *testing.T) {
	config := updateCleanupConfig{
		DeployPath:  "E:/bin-run",
		BuildOutput: "./bin",
	}

	dirs := collectTempCleanupDirs("E:/bin-run/gitmap.exe", "/repo", config)
	assertHasCleanupDir(t, dirs, os.TempDir())
	assertHasCleanupDir(t, dirs, "E:/gitmap-cli")
	assertHasCleanupDir(t, dirs, "/repo/bin")
}

func assertHasCleanupDir(t *testing.T, dirs []string, want string) {
	t.Helper()
	if hasCleanupDir(dirs, want) {
		return
	}

	t.Fatalf("cleanup dirs %v do not contain %q", dirs, want)
}
