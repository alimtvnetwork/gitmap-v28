package cmdnodes

import (
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestIsNodesCloneCommand(t *testing.T) {
	tests := []struct {
		input    string
		wantKind NodesCloneKind
		wantOk   bool
	}{
		{"clone", CloneKindClone, true},
		{"clone-except-self", CloneKindClone, true},
		{"clone-noself", CloneKindClone, true},
		{"cfr", CloneKindCFR, true},
		{"clone-fix-repo", CloneKindCFR, true},
		{"cfr-except-self", CloneKindCFR, true},
		{"cfrp", CloneKindCFRP, true},
		{"clone-fix-repo-pub", CloneKindCFRP, true},
		{"cfrp-except-self", CloneKindCFRP, true},
		{"status", "", false},
		{"unknown", "", false},
	}

	for _, tt := range tests {
		kind, ok := IsNodesCloneCommand(tt.input)
		if ok != tt.wantOk || kind != tt.wantKind {
			t.Errorf("IsNodesCloneCommand(%q) = (%v, %v); want (%v, %v)",
				tt.input, kind, ok, tt.wantKind, tt.wantOk)
		}
	}
}

func TestParseNodesCloneOptions_ExceptSelf(t *testing.T) {
	raw := []string{"except-self", "ChrisTitusTech/winutil", "--dry-run", "-t", "w1"}
	opts, ok := parseNodesCloneOptions(CloneKindClone, raw)
	if !ok {
		t.Fatalf("expected parse to succeed")
	}
	if opts.TargetFilter != "w1" {
		t.Errorf("TargetFilter = %q; want w1", opts.TargetFilter)
	}
	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun = true")
	}
	if !opts.IsSkipLocal {
		t.Errorf("expected IsSkipLocal = true for except-self")
	}
	if len(opts.PassArgs) == 0 || opts.PassArgs[0] != "ChrisTitusTech/winutil" {
		t.Errorf("PassArgs[0] = %v; want ChrisTitusTech/winutil", opts.PassArgs)
	}
}

func TestParseNodesCloneOptions_TargetDir(t *testing.T) {
	raw := []string{"https://github.com/user/repo", "D:\\custom\\path"}
	opts, ok := parseNodesCloneOptions(CloneKindClone, raw)
	if !ok {
		t.Fatalf("expected parse to succeed")
	}
	if opts.TargetDir != "D:\\custom\\path" {
		t.Errorf("TargetDir = %q; want D:\\custom\\path", opts.TargetDir)
	}
}

func TestBuildRemoteExecString(t *testing.T) {
	opts := NodesCloneOptions{
		Kind:     CloneKindCFR,
		PassArgs: []string{"gitmap.json", "--dry-run"},
		HasFile:  true,
	}
	winCmd := buildRemoteExecString(opts, "gitmap.json", true)
	if !strings.Contains(winCmd, "Set-Location \"D:\\work\"") || !strings.Contains(winCmd, "gitmap cfr") {
		t.Errorf("buildRemoteExecString(win) = %q; want Set-Location and gitmap cfr", winCmd)
	}

	unixCmd := buildRemoteExecString(opts, "gitmap.json", false)
	if !strings.Contains(unixCmd, "cd ~/work") || !strings.Contains(unixCmd, "gitmap cfr") {
		t.Errorf("buildRemoteExecString(unix) = %q; want cd ~/work and gitmap cfr", unixCmd)
	}

	optsCustom := NodesCloneOptions{
		Kind:      CloneKindClone,
		PassArgs:  []string{"https://github.com/user/repo"},
		TargetDir: "D:\\custom\\dir",
	}
	customWinCmd := buildRemoteExecString(optsCustom, "", true)
	if !strings.Contains(customWinCmd, "Set-Location \"D:\\custom\\dir\"") {
		t.Errorf("buildRemoteExecString(customDir) = %q; want custom path", customWinCmd)
	}
}

func TestIsWindowsNode(t *testing.T) {
	if !isWindowsNode(db.SSHConnection{OS: "windows"}) {
		t.Errorf("expected OS windows to be identified as Windows")
	}
	if !isWindowsNode(db.SSHConnection{OS: "win"}) {
		t.Errorf("expected OS win to be identified as Windows")
	}
	if !isWindowsNode(db.SSHConnection{Username: "Administrator", OS: "unknown"}) {
		t.Errorf("expected Administrator username to be identified as Windows")
	}
	if isWindowsNode(db.SSHConnection{OS: "linux", Username: "ubuntu"}) {
		t.Errorf("expected linux node not to be identified as Windows")
	}
}

func TestIsBashMissingError(t *testing.T) {
	err := errors.New("command execution failed: Process exited with status 1 (output: 'bash' is not recognized as an internal or external command)")
	if !isBashMissingError("", err) {
		t.Errorf("expected isBashMissingError to identify Windows bash missing error")
	}
	if !isBashMissingError("bash: command not found", nil) {
		t.Errorf("expected isBashMissingError to identify bash command not found")
	}
	if isBashMissingError("git clone failed", errors.New("fatal: repository not found")) {
		t.Errorf("expected unrelated error not to trigger bash missing")
	}
}

func TestSanitizeErrorAndStdout(t *testing.T) {
	rawErr := "Process exited with status 1 (output: 'bash' is not recognized as an internal or external command,\r\noperable program or batch file.\r\n)"
	cleanErr := sanitizeError(rawErr)
	if strings.Contains(cleanErr, "\n") || strings.Contains(cleanErr, "\r") {
		t.Errorf("expected sanitized error without newlines, got %q", cleanErr)
	}
	if !strings.Contains(cleanErr, "'bash' is not recognized") {
		t.Errorf("expected sanitized error to preserve core message, got %q", cleanErr)
	}

	rawOut := "===================================\nCloned awansoft-v10 successfully."
	cleanOut := sanitizeStdout(rawOut)
	if cleanOut != "Cloned awansoft-v10 successfully." {
		t.Errorf("expected clean stdout, got %q", cleanOut)
	}
}
