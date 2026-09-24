package cmdssh

import (
	"strings"
	"testing"
)

func TestSplitRemoteDestPath(t *testing.T) {
	// Windows tests
	dirWin, fileWin := splitRemoteDestPath(`C:\Windows\Temp\setup.exe`, true)
	if dirWin != `C:\Windows\Temp` || fileWin != "setup.exe" {
		t.Fatalf("expected C:\\Windows\\Temp and setup.exe, got %q and %q", dirWin, fileWin)
	}

	dirSlash, fileSlash := splitRemoteDestPath("C:/Windows/Temp/setup.exe", true)
	if dirSlash != `C:\Windows\Temp` || fileSlash != "setup.exe" {
		t.Fatalf("expected normalized C:\\Windows\\Temp, got %q and %q", dirSlash, fileSlash)
	}

	dirBare, fileBare := splitRemoteDestPath("setup.exe", true)
	if dirBare != `C:\Windows\Temp` || fileBare != "setup.exe" {
		t.Fatalf("expected fallback C:\\Windows\\Temp, got %q and %q", dirBare, fileBare)
	}

	// Unix tests
	dirUnix, fileUnix := splitRemoteDestPath("/tmp/setup.sh", false)
	if dirUnix != "/tmp" || fileUnix != "setup.sh" {
		t.Fatalf("expected /tmp and setup.sh, got %q and %q", dirUnix, fileUnix)
	}

	dirUnixSub, fileUnixSub := splitRemoteDestPath("/opt/deploy/bin/app", false)
	if dirUnixSub != "/opt/deploy/bin" || fileUnixSub != "app" {
		t.Fatalf("expected /opt/deploy/bin and app, got %q and %q", dirUnixSub, fileUnixSub)
	}

	dirUnixBare, fileUnixBare := splitRemoteDestPath("app", false)
	if dirUnixBare != "/tmp" || fileUnixBare != "app" {
		t.Fatalf("expected /tmp and app, got %q and %q", dirUnixBare, fileUnixBare)
	}
}

func TestBuildTarExtractCmd(t *testing.T) {
	winCmd := buildTarExtractCmd(`C:\Windows\Temp`, true)
	if !strings.Contains(winCmd, "cmd.exe") || !strings.Contains(winCmd, "tar.exe") {
		t.Fatalf("expected cmd.exe and tar.exe in winCmd, got %q", winCmd)
	}

	unixCmd := buildTarExtractCmd("/tmp", false)
	if !strings.Contains(unixCmd, "sh -c") || !strings.Contains(unixCmd, "tar -xf") {
		t.Fatalf("expected sh -c and tar -xf in unixCmd, got %q", unixCmd)
	}
}

func TestBuildDirectStreamCmd(t *testing.T) {
	winCmd := buildDirectStreamCmd(`C:\Windows\Temp\setup.exe`, true)
	if !strings.Contains(winCmd, "OpenStandardInput") {
		t.Fatalf("expected OpenStandardInput in winCmd, got %q", winCmd)
	}

	unixCmd := buildDirectStreamCmd("/tmp/setup.sh", false)
	if !strings.Contains(unixCmd, "cat >") || !strings.Contains(unixCmd, "chmod 755") {
		t.Fatalf("expected cat > and chmod 755 in unixCmd, got %q", unixCmd)
	}
}
