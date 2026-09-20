package cmdautomation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveLocateTarget(t *testing.T) {
	if resolveLocateTarget("") != "vcvarsall.bat" {
		t.Errorf("expected empty target to resolve to vcvarsall.bat")
	}
	if resolveLocateTarget("vcvars") != "vcvarsall.bat" {
		t.Errorf("expected vcvars to resolve to vcvarsall.bat")
	}
	if resolveLocateTarget("msbuild.exe") != "msbuild.exe" {
		t.Errorf("expected msbuild.exe to be preserved")
	}
}

func TestIsVcvarsTarget(t *testing.T) {
	if !isVcvarsTarget("vcvarsall.bat") {
		t.Errorf("expected vcvarsall.bat to match vcvars target")
	}
	if !isVcvarsTarget("vcvars") {
		t.Errorf("expected vcvars to match vcvars target")
	}
	if isVcvarsTarget("other.exe") {
		t.Errorf("expected other.exe to not match vcvars target")
	}
}

func TestScanRootForFile(t *testing.T) {
	tempDir := t.TempDir()
	nested := filepath.Join(tempDir, "sub", "deep")
	_ = os.MkdirAll(nested, 0o755)
	dummy := filepath.Join(nested, "sample_tool.exe")
	_ = os.WriteFile(dummy, []byte("echo tool"), 0o644)

	path, hasFound := scanRootForFile(tempDir, "sample_tool.exe", 3)
	if !hasFound || len(path) == 0 {
		t.Fatalf("expected sample_tool.exe to be found in scanRootForFile")
	}
}

func TestRunLocate_NotFound(t *testing.T) {
	_, err := RunLocate(LocateOptions{Target: "non_existent_tool_xyz123.exe"})
	if err == nil {
		t.Fatalf("expected error for non-existent tool")
	}
}
