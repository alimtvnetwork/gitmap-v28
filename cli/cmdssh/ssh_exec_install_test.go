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
	if hasWinGitmap == false || hasWinVersion == false || winShell != "cmd" {
		t.Errorf("getGitmapCheckCmd(windows) = (%q, %q); expected gitmap and cmd", winCmd, winShell)
	}

	unixCmd, unixShell := getGitmapCheckCmd("linux")
	hasPath := strings.Contains(unixCmd, "export PATH=")
	hasGitmap := strings.Contains(unixCmd, "gitmap")
	if hasPath == false || hasGitmap == false || unixShell != "sh" {
		t.Errorf("getGitmapCheckCmd(linux) = (%q, %q); expected export PATH, gitmap, and sh", unixCmd, unixShell)
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

func TestParseInstallExecArgs_DefaultOSAndForceAll(t *testing.T) {
	// 1. .exe automatically targets Windows
	optsExe := ParseInstallExecArgs([]string{"./setup.exe", "/SILENT"})
	if optsExe.TargetOS != "win" {
		t.Errorf("expected TargetOS='win' for .exe, got %q", optsExe.TargetOS)
	}
	if optsExe.IsForceAll {
		t.Errorf("expected IsForceAll=false by default")
	}

	// 2. .sh automatically targets Unix
	optsSh := ParseInstallExecArgs([]string{"./bootstrap.sh"})
	if optsSh.TargetOS != "unix" {
		t.Errorf("expected TargetOS='unix' for .sh, got %q", optsSh.TargetOS)
	}

	// 3. --force-all bypasses default OS targeting
	optsForce := ParseInstallExecArgs([]string{"./setup.exe", "--force-all"})
	if !optsForce.IsForceAll {
		t.Errorf("expected IsForceAll=true on --force-all")
	}
	if optsForce.TargetOS != "" {
		t.Errorf("expected TargetOS='' when --force-all is set, got %q", optsForce.TargetOS)
	}

	// 4. Explicit --os wins
	optsExplicit := ParseInstallExecArgs([]string{"./setup.exe", "--os", "linux"})
	if optsExplicit.TargetOS != "linux" {
		t.Errorf("expected TargetOS='linux' when explicitly set, got %q", optsExplicit.TargetOS)
	}

	// 5. Default IsSilent is true, overridden by --no-silent
	optsSilent := ParseInstallExecArgs([]string{"./setup.exe"})
	if !optsSilent.IsSilent {
		t.Errorf("expected IsSilent=true by default")
	}
	optsNoSilent := ParseInstallExecArgs([]string{"./setup.exe", "--no-silent"})
	if optsNoSilent.IsSilent {
		t.Errorf("expected IsSilent=false when --no-silent is passed")
	}
}

func TestBuildRemoteInstallerExecCmd_PayloadDetection(t *testing.T) {
	// 1. NSIS detection
	nsisData := []byte("...NullsoftInst...")
	cmdNSIS, shellNSIS := BuildRemoteInstallerExecCmdWithPayload("windows", "C:\\Temp\\setup.exe", nil, true, nsisData)
	if shellNSIS != "ps" || !strings.Contains(cmdNSIS, "'-ArgumentList', '/S'") && !strings.Contains(cmdNSIS, "-ArgumentList '/S'") {
		t.Errorf("expected NSIS to use /S argument, got cmd: %s", cmdNSIS)
	}

	// 2. Inno Setup detection
	innoData := []byte("...Inno Setup...")
	cmdInno, _ := BuildRemoteInstallerExecCmdWithPayload("windows", "C:\\Temp\\setup.exe", nil, true, innoData)
	if !strings.Contains(cmdInno, "/VERYSILENT") {
		t.Errorf("expected Inno Setup to use /VERYSILENT, got: %s", cmdInno)
	}

	// 3. MSI
	cmdMSI, _ := BuildRemoteInstallerExecCmdWithPayload("windows", "C:\\Temp\\setup.msi", nil, true, nil)
	if !strings.Contains(cmdMSI, "msiexec.exe") {
		t.Errorf("expected MSI to invoke msiexec.exe, got: %s", cmdMSI)
	}
}
