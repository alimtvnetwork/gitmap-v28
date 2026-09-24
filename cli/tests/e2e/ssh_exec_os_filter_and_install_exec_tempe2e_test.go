//go:build tempe2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func requireTempE2EEnvSpec152(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("Skipping temporary E2E test: set RUN_TEMP_E2E=1 with -tags=tempe2e to execute on-demand")
	}
}

// TestTempE2E_SSHExecOSFiltering verifies --except-os (unix, win, ubuntu, linux, darwin) and --os filters.
func TestTempE2E_SSHExecOSFiltering(t *testing.T) {
	requireTempE2EEnvSpec152(t)

	sampleConns := []db.SSHConnection{
		{Alias: "win-srv", IPAddress: "10.0.0.1", OS: "windows"},
		{Alias: "linux-worker", IPAddress: "10.0.0.2", OS: "linux"},
		{Alias: "ubuntu-box", IPAddress: "10.0.0.3", OS: "linux"},
		{Alias: "mac-mini", IPAddress: "10.0.0.4", OS: "darwin"},
	}

	// 1. --except-os unix -> must keep only Windows
	exceptUnix := cmdssh.FilterSSHConnectionsByOS(sampleConns, "", "unix")
	if len(exceptUnix) != 1 || exceptUnix[0].Alias != "win-srv" {
		t.Fatalf("expected only win-srv for --except-os unix, got %+v", exceptUnix)
	}

	// 2. --except-os win -> must exclude Windows
	exceptWin := cmdssh.FilterSSHConnectionsByOS(sampleConns, "", "win")
	if len(exceptWin) != 3 {
		t.Fatalf("expected 3 non-windows machines for --except-os win, got %+v", exceptWin)
	}
	for _, c := range exceptWin {
		if c.OS == "windows" {
			t.Fatalf("found windows machine in --except-os win: %+v", c)
		}
	}

	// 3. --except-os ubuntu -> must exclude Linux nodes
	exceptUbuntu := cmdssh.FilterSSHConnectionsByOS(sampleConns, "", "ubuntu")
	if len(exceptUbuntu) != 2 {
		t.Fatalf("expected 2 machines (win + mac) for --except-os ubuntu, got %+v", exceptUbuntu)
	}
	for _, c := range exceptUbuntu {
		if c.OS == "linux" {
			t.Fatalf("found linux machine in --except-os ubuntu: %+v", c)
		}
	}

	// 4. Target OS: --os win -> must keep only Windows
	onlyWin := cmdssh.FilterSSHConnectionsByOS(sampleConns, "win", "")
	if len(onlyWin) != 1 || onlyWin[0].Alias != "win-srv" {
		t.Fatalf("expected only win-srv for --os win, got %+v", onlyWin)
	}

	// 5. Target OS: --os unix -> must keep Linux and Darwin
	onlyUnix := cmdssh.FilterSSHConnectionsByOS(sampleConns, "unix", "")
	if len(onlyUnix) != 3 {
		t.Fatalf("expected 3 unix machines for --os unix, got %+v", onlyUnix)
	}
}

// TestTempE2E_SSHExecMultiCommandFlagParsing verifies argument parsing with --except-os in any position.
func TestTempE2E_SSHExecMultiCommandFlagParsing(t *testing.T) {
	requireTempE2EEnvSpec152(t)

	// Command first, flag second: gitmap ssh exec cmd1,cmd2,cmd3 --except-os unix
	opts1 := cmdssh.ParseSEFlags([]string{"cmd1,cmd2,cmd3", "--except-os", "unix"})
	if opts1.ExceptOS != "unix" {
		t.Fatalf("expected ExceptOS 'unix', got %q", opts1.ExceptOS)
	}
	if len(opts1.Args) != 1 || opts1.Args[0] != "cmd1,cmd2,cmd3" {
		t.Fatalf("expected args ['cmd1,cmd2,cmd3'], got %v", opts1.Args)
	}

	// Flag first, command second: gitmap ssh exec --except-os win cmd1,cmd2,cmd3
	opts2 := cmdssh.ParseSEFlags([]string{"--except-os", "win", "cmd1,cmd2,cmd3"})
	if opts2.ExceptOS != "win" {
		t.Fatalf("expected ExceptOS 'win', got %q", opts2.ExceptOS)
	}
	if len(opts2.Args) != 1 || opts2.Args[0] != "cmd1,cmd2,cmd3" {
		t.Fatalf("expected args ['cmd1,cmd2,cmd3'], got %v", opts2.Args)
	}

	// With ubuntu: gitmap ssh exec cmd1,cmd2,cmd3 --except-os ubuntu
	opts3 := cmdssh.ParseSEFlags([]string{"cmd1,cmd2,cmd3", "--except-os", "ubuntu"})
	if opts3.ExceptOS != "ubuntu" {
		t.Fatalf("expected ExceptOS 'ubuntu', got %q", opts3.ExceptOS)
	}
}

// TestTempE2E_SSHInstallExecParsingAndCmdGen verifies install-exec parsing and execution generation.
func TestTempE2E_SSHInstallExecParsingAndCmdGen(t *testing.T) {
	requireTempE2EEnvSpec152(t)

	// 1. Parsing args
	args := []string{"./setup.exe", "/SILENT", "/DIR=C:\\App", "--except-os", "unix", "--except", "worker-1,10.0.0.5"}
	opts := cmdssh.ParseInstallExecArgs(args)

	if opts.SetupPath != "./setup.exe" {
		t.Fatalf("expected SetupPath './setup.exe', got %q", opts.SetupPath)
	}
	if len(opts.InstallerArgs) != 2 || opts.InstallerArgs[0] != "/SILENT" {
		t.Fatalf("expected installer args ['/SILENT', '/DIR=C:\\App'], got %v", opts.InstallerArgs)
	}
	if opts.ExceptOS != "unix" {
		t.Fatalf("expected ExceptOS 'unix', got %q", opts.ExceptOS)
	}
	if opts.Except != "worker-1,10.0.0.5" {
		t.Fatalf("expected Except 'worker-1,10.0.0.5', got %q", opts.Except)
	}

	// 2. Command generation for Windows .exe
	winCmd, winShell := cmdssh.BuildRemoteInstallerExecCmd("windows", "C:\\Windows\\Temp\\setup.exe", []string{"/SILENT"}, false)
	if winShell != "ps" {
		t.Fatalf("expected 'ps' shell for Windows installer, got %q", winShell)
	}
	if !strings.Contains(winCmd, "Start-Process") || !strings.Contains(winCmd, "setup.exe") {
		t.Fatalf("expected Start-Process in command, got %q", winCmd)
	}

	// 3. Command generation for Windows .msi with silent flag
	msiCmd, _ := cmdssh.BuildRemoteInstallerExecCmd("windows", "C:\\Windows\\Temp\\app.msi", nil, true)
	if !strings.Contains(msiCmd, "msiexec.exe") || !strings.Contains(msiCmd, "/qn") {
		t.Fatalf("expected msiexec with /qn for silent MSI, got %q", msiCmd)
	}

	// 4. Command generation for Linux script/binary
	unixCmd, unixShell := cmdssh.BuildRemoteInstallerExecCmd("linux", "/tmp/installer.sh", []string{"--install"}, false)
	if unixShell != "bash" {
		t.Fatalf("expected 'bash' shell for Linux installer, got %q", unixShell)
	}
	if !strings.Contains(unixCmd, "chmod +x") || !strings.Contains(unixCmd, "/tmp/installer.sh") {
		t.Fatalf("expected chmod +x in Unix installer cmd, got %q", unixCmd)
	}
}

// TestTempE2E_SSHInstallExecDryRun verifies dry-run execution on mock local file.
func TestTempE2E_SSHInstallExecDryRun(t *testing.T) {
	requireTempE2EEnvSpec152(t)

	tmpDir := t.TempDir()
	dummySetup := filepath.Join(tmpDir, "dummy_installer.exe")
	if err := os.WriteFile(dummySetup, []byte("MOCK-EXE-DATA"), 0755); err != nil {
		t.Fatalf("failed to write dummy setup: %v", err)
	}

	err := cmdssh.RunSSHInstallExecCLI([]string{dummySetup, "/S", "--dry-run", "--except-os", "unix"})
	if err != nil {
		t.Fatalf("expected dry-run to exit cleanly, got error: %v", err)
	}
}
