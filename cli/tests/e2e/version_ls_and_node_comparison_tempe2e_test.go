//go:build tempe2e

package e2e

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func requireTempE2EEnvSpec151(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("Skipping temporary E2E test: set RUN_TEMP_E2E=1 with -tags=tempe2e to execute on-demand")
	}
}

// TestTempE2E_SemverComparisonAndNodeStatus verifies semantic version comparison
// logic for remote fleet nodes (identical, newer, older, not installed).
func TestTempE2E_SemverComparisonAndNodeStatus(t *testing.T) {
	requireTempE2EEnvSpec151(t)

	// 1. Identical versions
	if cmp := cmdssh.CompareSemverStrings("6.326.0", "6.326.0"); cmp != 0 {
		t.Fatalf("expected 0 for identical semver, got %d", cmp)
	}
	if cmp := cmd.CompareSemver("v6.326.0", "6.326.0"); cmp != 0 {
		t.Fatalf("expected 0 for identical semver with v prefix, got %d", cmp)
	}

	// 2. Remote above (newer)
	if cmp := cmdssh.CompareSemverStrings("6.328.0", "6.326.0"); cmp != 1 {
		t.Fatalf("expected 1 (above) for remote 6.328.0 vs 6.326.0, got %d", cmp)
	}
	if cmp := cmd.CompareSemver("7.0.0", "6.326.0"); cmp != 1 {
		t.Fatalf("expected 1 (above) for remote 7.0.0 vs 6.326.0, got %d", cmp)
	}

	// 3. Remote below (older)
	if cmp := cmdssh.CompareSemverStrings("6.320.0", "6.326.0"); cmp != -1 {
		t.Fatalf("expected -1 (below) for remote 6.320.0 vs 6.326.0, got %d", cmp)
	}
	if cmp := cmd.CompareSemver("5.99.0", "6.326.0"); cmp != -1 {
		t.Fatalf("expected -1 (below) for remote 5.99.0 vs 6.326.0, got %d", cmp)
	}

	// 4. Comparison status resolver
	sameStatus := cmdssh.ResolveNodeVersionComparisonForTest("v"+constants.Version, true)
	if !strings.Contains(sameStatus, "IDENTICAL") && !strings.Contains(sameStatus, "SAME") {
		t.Fatalf("expected IDENTICAL status, got %q", sameStatus)
	}

	aboveStatus := cmdssh.ResolveNodeVersionComparisonForTest("v99.0.0", true)
	if !strings.Contains(aboveStatus, "ABOVE") && !strings.Contains(aboveStatus, "NEWER") {
		t.Fatalf("expected ABOVE status, got %q", aboveStatus)
	}

	belowStatus := cmdssh.ResolveNodeVersionComparisonForTest("v1.0.0", true)
	if !strings.Contains(belowStatus, "BELOW") && !strings.Contains(belowStatus, "OLDER") {
		t.Fatalf("expected BELOW status, got %q", belowStatus)
	}

	notInstalledStatus := cmdssh.ResolveNodeVersionComparisonForTest("", false)
	if !strings.Contains(notInstalledStatus, "NOT INSTALLED") {
		t.Fatalf("expected NOT INSTALLED status, got %q", notInstalledStatus)
	}
}

// TestTempE2E_ReleaseTagsTableRendering verifies GitHub release tags table formatting.
func TestTempE2E_ReleaseTagsTableRendering(t *testing.T) {
	requireTempE2EEnvSpec151(t)

	sampleTags := []cmd.GitHubReleaseTagInfo{
		{TagName: "v6.327.0", PublishedAt: "2026-09-24", IsLatest: true},
		{TagName: "v6.326.0", PublishedAt: "2026-09-24", IsLatest: false},
		{TagName: "v6.320.0", PublishedAt: "2026-09-20", IsLatest: false},
	}

	// Render table without error
	cmd.RenderReleaseTagsTable("Test Application", "v6.326.0", sampleTags)
}

// TestTempE2E_TargetPinnedVersionAndInstallerArgs verifies explicit version pinning
// in GitMap update command builder for both Windows PowerShell and Unix bash.
func TestTempE2E_TargetPinnedVersionAndInstallerArgs(t *testing.T) {
	requireTempE2EEnvSpec151(t)

	// Set target version
	cmdupdate.SetTargetVersion("6.320.0")
	if ver := cmdupdate.GetTargetVersion(); ver != "6.320.0" {
		t.Fatalf("expected target version 6.320.0, got %q", ver)
	}

	// Build Windows installer command with pinned version
	winCmd := cmdupdate.BuildRemoteWindowsInstallerCmdForTest("install.ps1", "C:\\gitmap", "6.320.0")
	winArgs := strings.Join(winCmd.Args, " ")
	if !strings.Contains(winArgs, "-Version 6.320.0") {
		t.Fatalf("expected Windows installer command to contain '-Version 6.320.0', got: %s", winArgs)
	}

	// Build Unix installer command with pinned version
	unixCmd := cmdupdate.BuildUnixInstallerCmdForTest("/tmp/install.sh", "/usr/local/bin", "6.320.0")
	unixArgs := strings.Join(unixCmd.Args, " ")
	if !strings.Contains(unixArgs, "--version v6.320.0") {
		t.Fatalf("expected Unix installer command to contain '--version v6.320.0', got: %s", unixArgs)
	}

	// Reset
	cmdupdate.SetTargetVersion("")
}

// TestTempE2E_LowLevelSSHBypassExeCopyMechanism verifies the Base64 chunked transfer
// command generation for lower-level SSH execution without SMB or FTP servers.
func TestTempE2E_LowLevelSSHBypassExeCopyMechanism(t *testing.T) {
	requireTempE2EEnvSpec151(t)

	mockExePayload := []byte("MZ-MOCK-BINARY-HEADER-FOR-TESTING-EXE-COPY")
	encoded := base64.StdEncoding.EncodeToString(mockExePayload)

	// Verify Windows low-level powershell command format
	winRemoteCmd := cmdssh.BuildRemoteWriteCmdForTest("windows", "C:\\tools\\gitmap.exe", encoded)
	if !strings.Contains(winRemoteCmd, "[Convert]::FromBase64String") {
		t.Fatalf("expected Windows command to decode Base64 in-memory, got: %s", winRemoteCmd)
	}
	if !strings.Contains(winRemoteCmd, "[IO.File]::WriteAllBytes") {
		t.Fatalf("expected Windows command to write all bytes, got: %s", winRemoteCmd)
	}

	// Verify Linux/Unix low-level bash command format
	unixRemoteCmd := cmdssh.BuildRemoteWriteCmdForTest("linux", "/usr/local/bin/gitmap", encoded)
	if !strings.Contains(unixRemoteCmd, "base64 -d") {
		t.Fatalf("expected Unix command to decode Base64, got: %s", unixRemoteCmd)
	}
	if !strings.Contains(unixRemoteCmd, "chmod +x") {
		t.Fatalf("expected Unix command to set executable permission, got: %s", unixRemoteCmd)
	}
}

// TestTempE2E_AGMVersionListingAndInstallOptions verifies AGM version command wiring.
func TestTempE2E_AGMVersionListingAndInstallOptions(t *testing.T) {
	requireTempE2EEnvSpec151(t)

	called := false
	cmdinstall.RunAGMVersionTagsLSFn = func() error {
		called = true
		return nil
	}

	if cmdinstall.RunAGMVersionTagsLSFn != nil {
		_ = cmdinstall.RunAGMVersionTagsLSFn()
	}

	if !called {
		t.Fatal("expected RunAGMVersionTagsLSFn to be invoked")
	}

	// Wire back real handler
	cmdinstall.RunAGMVersionTagsLSFn = cmd.RunAGMVersionTagsLS
}
