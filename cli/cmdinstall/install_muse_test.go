package cmdinstall

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestMuseResolveInstallCommand(t *testing.T) {
	winBin, winArgs, winErr := ResolveMuseInstallCommand(MusePlatformWindows)
	if winErr != nil {
		t.Fatalf("Windows command resolution failed: %v", winErr)
	}
	if winBin != "powershell" {
		t.Errorf("Expected powershell for Windows, got: %s", winBin)
	}
	winCmdStr := strings.Join(winArgs, " ")
	if !strings.Contains(winCmdStr, "install.ps1") {
		t.Errorf("Expected install.ps1 in Windows args, got: %s", winCmdStr)
	}

	linuxBin, linuxArgs, linuxErr := ResolveMuseInstallCommand(MusePlatformLinux)
	if linuxErr != nil {
		t.Fatalf("Linux command resolution failed: %v", linuxErr)
	}
	if linuxBin != "bash" {
		t.Errorf("Expected bash for Linux, got: %s", linuxBin)
	}
	linuxCmdStr := strings.Join(linuxArgs, " ")
	if !strings.Contains(linuxCmdStr, "install.sh") {
		t.Errorf("Expected install.sh in Linux args, got: %s", linuxCmdStr)
	}

	darwinBin, darwinArgs, darwinErr := ResolveMuseInstallCommand(MusePlatformDarwin)
	if darwinErr != nil {
		t.Fatalf("Darwin command resolution failed: %v", darwinErr)
	}
	if darwinBin != "bash" {
		t.Errorf("Expected bash for Darwin, got: %s", darwinBin)
	}
	darwinCmdStr := strings.Join(darwinArgs, " ")
	if !strings.Contains(darwinCmdStr, "install.sh") {
		t.Errorf("Expected install.sh in Darwin args, got: %s", darwinCmdStr)
	}

	_, _, badErr := ResolveMuseInstallCommand(MusePlatformType("unknown-os"))
	if badErr == nil {
		t.Errorf("Expected error for unsupported platform, got nil")
	}
}

func TestMuseResolvePlatform(t *testing.T) {
	if got := ResolveMusePlatform("win"); got != MusePlatformWindows {
		t.Errorf("Expected windows for 'win', got: %s", got)
	}
	if got := ResolveMusePlatform("windows"); got != MusePlatformWindows {
		t.Errorf("Expected windows for 'windows', got: %s", got)
	}
	if got := ResolveMusePlatform("ubuntu"); got != MusePlatformLinux {
		t.Errorf("Expected linux for 'ubuntu', got: %s", got)
	}
	if got := ResolveMusePlatform("linux"); got != MusePlatformLinux {
		t.Errorf("Expected linux for 'linux', got: %s", got)
	}
	if got := ResolveMusePlatform("macos"); got != MusePlatformDarwin {
		t.Errorf("Expected darwin for 'macos', got: %s", got)
	}
	if got := ResolveMusePlatform("darwin"); got != MusePlatformDarwin {
		t.Errorf("Expected darwin for 'darwin', got: %s", got)
	}
}

func TestMuseDryRun(t *testing.T) {
	optsWin := MuseInstallOptions{
		Platform: MusePlatformWindows,
		DryRun:   true,
	}
	resWin, err := RunMuseInstaller(optsWin)
	if err != nil {
		t.Fatalf("DryRun Windows failed: %v", err)
	}
	if !resWin.Success {
		t.Errorf("Expected Success true for dry-run")
	}
	if !strings.Contains(resWin.CommandExecuted, "install.ps1") {
		t.Errorf("Expected install.ps1 in CommandExecuted, got: %s", resWin.CommandExecuted)
	}

	optsLinux := MuseInstallOptions{
		Platform: MusePlatformLinux,
		DryRun:   true,
	}
	resLinux, err := RunMuseInstaller(optsLinux)
	if err != nil {
		t.Fatalf("DryRun Linux failed: %v", err)
	}
	if !resLinux.Success {
		t.Errorf("Expected Success true for Linux dry-run")
	}
	if !strings.Contains(resLinux.CommandExecuted, "install.sh") {
		t.Errorf("Expected install.sh in Linux CommandExecuted, got: %s", resLinux.CommandExecuted)
	}
}

func TestMuseToolAliases(t *testing.T) {
	if got := resolveToolAlias("muse"); got != constants.ToolMuse {
		t.Errorf("resolveToolAlias('muse') = %s, want %s", got, constants.ToolMuse)
	}
	if got := resolveToolAlias("meta-muse"); got != constants.ToolMetaMuse {
		t.Errorf("resolveToolAlias('meta-muse') = %s, want %s", got, constants.ToolMetaMuse)
	}
	if got := resolveToolAlias("metamuse"); got != constants.ToolMetaMuse {
		t.Errorf("resolveToolAlias('metamuse') = %s, want %s", got, constants.ToolMetaMuse)
	}
}
