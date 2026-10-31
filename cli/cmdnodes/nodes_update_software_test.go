package cmdnodes

import (
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestNormalizeSoftwareApp(t *testing.T) {
	cases := []struct {
		input       string
		expected    string
		isSupported bool
	}{
		{"cursor", "cursor", true},
		{"Cursor", "cursor", true},
		{"cursor-ide", "cursor", true},
		{"vscode", "vscode", true},
		{"code", "vscode", true},
		{"vs-code", "vscode", true},
		{"agm", "agm", true},
		{"ag-manager", "agm", true},
		{"antigravity-manager", "agm", true},
		{"agy", "agy", true},
		{"agy-tools", "agy", true},
		{"antigravity", "agy", true},
		{"antigravity-ide", "agy", true},
		{"ide", "agy", true},
		{"unknown", "", false},
		{"", "", false},
	}

	for _, tc := range cases {
		app, isSupported := normalizeSoftwareApp(tc.input)
		if isSupported != tc.isSupported {
			t.Errorf("normalizeSoftwareApp(%q) isSupported = %v; want %v", tc.input, isSupported, tc.isSupported)
		}
		if app != tc.expected {
			t.Errorf("normalizeSoftwareApp(%q) app = %q; want %q", tc.input, app, tc.expected)
		}
	}
}

func TestResolveSoftwareUpdateCommand_Windows(t *testing.T) {
	cursorCmd := resolveSoftwareUpdateCommand("cursor", true)
	if !strings.Contains(cursorCmd, "Anysphere.Cursor") {
		t.Errorf("expected cursor windows command to contain Anysphere.Cursor, got %q", cursorCmd)
	}

	vscodeCmd := resolveSoftwareUpdateCommand("vscode", true)
	if !strings.Contains(vscodeCmd, "Microsoft.VisualStudioCode") {
		t.Errorf("expected vscode windows command to contain Microsoft.VisualStudioCode, got %q", vscodeCmd)
	}

	agmCmd := resolveSoftwareUpdateCommand("agm", true)
	if !strings.Contains(agmCmd, "install.ps1") {
		t.Errorf("expected agm windows command to contain install.ps1, got %q", agmCmd)
	}

	agyCmd := resolveSoftwareUpdateCommand("agy", true)
	if !strings.Contains(agyCmd, "install.ps1") {
		t.Errorf("expected agy windows command to contain install.ps1, got %q", agyCmd)
	}
}

func TestResolveSoftwareUpdateCommand_Linux(t *testing.T) {
	cursorCmd := resolveSoftwareUpdateCommand("cursor", false)
	if !strings.Contains(cursorCmd, "setup-cursor-ubuntu.py") {
		t.Errorf("expected cursor linux command to contain setup-cursor-ubuntu.py, got %q", cursorCmd)
	}

	vscodeCmd := resolveSoftwareUpdateCommand("vscode", false)
	if !strings.Contains(vscodeCmd, "apt-get") {
		t.Errorf("expected vscode linux command to contain apt-get, got %q", vscodeCmd)
	}

	agmCmd := resolveSoftwareUpdateCommand("agm", false)
	if !strings.Contains(agmCmd, "install.sh") {
		t.Errorf("expected agm linux command to contain install.sh, got %q", agmCmd)
	}

	agyCmd := resolveSoftwareUpdateCommand("agy", false)
	if !strings.Contains(agyCmd, "install.sh") {
		t.Errorf("expected agy linux command to contain install.sh, got %q", agyCmd)
	}
}

func TestParseSoftwareUpdateOptions(t *testing.T) {
	args := []string{"u1", "--dry-run", "--json", "--verbose", "--deploy"}
	opts := parseSoftwareUpdateOptions(args)

	if opts.FilterOpts.Target != "u1" {
		t.Errorf("expected target u1, got %q", opts.FilterOpts.Target)
	}
	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true")
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON true")
	}
	if !opts.IsVerbose {
		t.Errorf("expected IsVerbose true")
	}
	if !opts.IsDeploy {
		t.Errorf("expected IsDeploy true")
	}
}

func TestParseSoftwareUpdateOptions_AllTarget(t *testing.T) {
	args := []string{"all", "-n"}
	opts := parseSoftwareUpdateOptions(args)

	if opts.FilterOpts.Target != "" {
		t.Errorf("expected empty target for all, got %q", opts.FilterOpts.Target)
	}
	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true")
	}
}

func TestBuildDryRunSoftwareResult(t *testing.T) {
	conn := db.SSHConnection{
		Alias:     "u1",
		IPAddress: "192.168.1.101",
		OS:        "linux",
	}
	res := buildDryRunSoftwareResult(conn, "cursor")

	if res.Status != "DRY-RUN" {
		t.Errorf("expected status DRY-RUN, got %q", res.Status)
	}
	if !res.IsSuccess {
		t.Errorf("expected IsSuccess true")
	}
	if !strings.Contains(res.Details, "setup-cursor-ubuntu.py") {
		t.Errorf("expected details to contain command, got %q", res.Details)
	}
}

func TestBuildOfflineSoftwareResult(t *testing.T) {
	conn := db.SSHConnection{
		Alias:     "w1",
		IPAddress: "192.168.1.102",
		OS:        "windows",
	}
	res := buildOfflineSoftwareResult(conn, "vscode")

	if res.Status != "SKIPPED" {
		t.Errorf("expected status SKIPPED, got %q", res.Status)
	}
	if !res.IsOffline {
		t.Errorf("expected IsOffline true")
	}
}

func TestBuildSoftwareUpdateSummary(t *testing.T) {
	results := []SoftwareUpdateNodeResult{
		{NodeAlias: "u1", IsSuccess: true, Status: "SUCCESS"},
		{NodeAlias: "u2", IsOffline: true, Status: "SKIPPED"},
		{NodeAlias: "u3", Status: "FAILED"},
	}

	summary := buildSoftwareUpdateSummary("cursor", results, 500*time.Millisecond)

	if summary.App != "cursor" {
		t.Errorf("expected app cursor, got %q", summary.App)
	}
	if summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Total)
	}
	if summary.Success != 1 {
		t.Errorf("expected success 1, got %d", summary.Success)
	}
	if summary.Skipped != 1 {
		t.Errorf("expected skipped 1, got %d", summary.Skipped)
	}
	if summary.Failed != 1 {
		t.Errorf("expected failed 1, got %d", summary.Failed)
	}
}

func TestRunFleetSoftwareUpdate_Unsupported(t *testing.T) {
	err := RunFleetSoftwareUpdate("nonexistent-tool", []string{"all"})
	if err == nil {
		t.Fatalf("expected error for unsupported tool, got nil")
	}
}

func TestIsNodesUpdateSoftwareRequest(t *testing.T) {
	if !isNodesUpdateSoftwareRequest([]string{"update", "cursor"}) {
		t.Errorf("expected true for update")
	}
	if !isNodesUpdateSoftwareRequest([]string{"up", "vscode"}) {
		t.Errorf("expected true for up")
	}
	if isNodesUpdateSoftwareRequest([]string{}) {
		t.Errorf("expected false for empty args")
	}
	if isNodesUpdateSoftwareRequest([]string{"status"}) {
		t.Errorf("expected false for status")
	}
}
