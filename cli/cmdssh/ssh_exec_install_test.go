package cmdssh

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestGetGitmapCheckCmd(t *testing.T) {
	winCmd, winShell := getGitmapCheckCmd("windows")
	hasWinGitmap := strings.Contains(winCmd, "gitmap")
	hasWinVersion := strings.Contains(winCmd, "version")
	if hasWinGitmap == false || hasWinVersion == false || winShell != "" {
		t.Errorf("getGitmapCheckCmd(windows) = (%q, %q); expected gitmap and version check", winCmd, winShell)
	}

	unixCmd, unixShell := getGitmapCheckCmd("linux")
	hasPath := strings.Contains(unixCmd, "export PATH=")
	hasVersion := strings.Contains(unixCmd, "version")
	hasGitmap := strings.Contains(unixCmd, "gitmap")
	if hasPath == false || hasVersion == false || hasGitmap == false || unixShell != "bash" {
		t.Errorf("getGitmapCheckCmd(linux) = (%q, %q); expected export PATH, gitmap, version, and bash", unixCmd, unixShell)
	}
}

func TestGetGitmapInstallCmd(t *testing.T) {
	winCmd, winShell := getGitmapInstallCmd("windows")
	hasWinUrl := strings.Contains(winCmd, constants.SelfInstallRemotePwsh)
	if hasWinUrl == false || winShell != "ps" {
		t.Errorf("getGitmapInstallCmd(windows) = (%q, %q); expected Pwsh URL and ps", winCmd, winShell)
	}

	unixCmd, unixShell := getGitmapInstallCmd("linux")
	hasUnixUrl := strings.Contains(unixCmd, constants.SelfInstallRemoteBash)
	if hasUnixUrl == false || unixShell != "bash" {
		t.Errorf("getGitmapInstallCmd(linux) = (%q, %q); expected Bash URL and bash", unixCmd, unixShell)
	}
}

func TestWrapUnixPath(t *testing.T) {
	raw := "gitmap status"
	wrapped := wrapUnixPath(raw)
	hasPrefix := strings.HasPrefix(wrapped, "export PATH=")
	if hasPrefix == false {
		t.Errorf("wrapUnixPath(%q) = %q; expected export PATH prefix", raw, wrapped)
	}

	idempotent := wrapUnixPath(wrapped)
	if idempotent != wrapped {
		t.Errorf("wrapUnixPath idempotent failed: got %q, want %q", idempotent, wrapped)
	}
}
