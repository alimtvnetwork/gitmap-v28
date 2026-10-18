package cmdspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsExecutableFile(t *testing.T) {
	dir := t.TempDir()
	exePath := filepath.Join(dir, "sender.exe")
	if writeErr := os.WriteFile(exePath, []byte("echo"), 0755); writeErr != nil {
		t.Fatal(writeErr)
	}

	isExe := isExecutableFile(exePath)
	if !isExe {
		t.Fatalf("expected %s to be executable", exePath)
	}

	isDirExe := isExecutableFile(dir)
	if isDirExe {
		t.Fatalf("directory should not be executable")
	}
}

func TestCheckWindowsExecutable(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"test.exe", true},
		{"test.cmd", true},
		{"test.bat", true},
		{"test", true},
		{"test.txt", false},
	}
	for _, tc := range cases {
		got := checkWindowsExecutable(tc.path)
		if got != tc.want {
			t.Errorf("checkWindowsExecutable(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestResolveGoEmailReaderBinCandidates(t *testing.T) {
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(binDir, "GoEmailSender.exe")
	if err := os.WriteFile(binPath, []byte("echo"), 0755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("GOPATH", dir)
	resolved := resolveGoEmailReaderBin()
	if len(resolved) == 0 {
		t.Fatalf("expected binary resolution from GOPATH, got empty")
	}
}

func TestResolveCandidatePath(t *testing.T) {
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatal(err)
	}
	binFile := filepath.Join(binDir, "goemailsender.exe")
	if err := os.WriteFile(binFile, []byte("echo"), 0755); err != nil {
		t.Fatal(err)
	}

	found := probeCandidatePath("goemailsender", "", dir)
	if len(found) == 0 {
		t.Fatalf("probeCandidatePath failed to find %s", binFile)
	}
}

func TestWindowsExecutableExt(t *testing.T) {
	hasExe := HasWindowsExecutableExt("bin.exe")
	if !hasExe {
		t.Errorf("expected bin.exe to have valid windows extension")
	}
	hasTxt := HasWindowsExecutableExt("doc.txt")
	if hasTxt {
		t.Errorf("doc.txt should not have windows executable extension")
	}
}
