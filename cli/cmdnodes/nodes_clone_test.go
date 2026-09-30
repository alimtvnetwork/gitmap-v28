package cmdnodes

import (
	"strings"
	"testing"
)

func TestIsNodesCloneCommand(t *testing.T) {
	tests := []struct {
		input    string
		wantKind NodesCloneKind
		wantOk   bool
	}{
		{"clone", CloneKindClone, true},
		{"cfr", CloneKindCFR, true},
		{"clone-fix-repo", CloneKindCFR, true},
		{"cfrp", CloneKindCFRP, true},
		{"clone-fix-repo-pub", CloneKindCFRP, true},
		{"cfr-pub", CloneKindCFRP, true},
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

func TestParseNodesCloneOptions(t *testing.T) {
	raw := []string{"ChrisTitusTech/winutil", "--dry-run", "-t", "w1", "--skip-local"}
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
		t.Errorf("expected IsSkipLocal = true")
	}
	if len(opts.PassArgs) == 0 || opts.PassArgs[0] != "ChrisTitusTech/winutil" {
		t.Errorf("PassArgs[0] = %v; want ChrisTitusTech/winutil", opts.PassArgs)
	}
}

func TestResolveRemoteWorkDir(t *testing.T) {
	if got := ResolveRemoteWorkDir("windows"); got != "D:/work" {
		t.Errorf("ResolveRemoteWorkDir(windows) = %q; want D:/work", got)
	}
	if got := ResolveRemoteWorkDir("win"); got != "D:/work" {
		t.Errorf("ResolveRemoteWorkDir(win) = %q; want D:/work", got)
	}
	if got := ResolveRemoteWorkDir("linux"); got != "~/work" {
		t.Errorf("ResolveRemoteWorkDir(linux) = %q; want ~/work", got)
	}
	if got := ResolveRemoteWorkDir("darwin"); got != "~/work" {
		t.Errorf("ResolveRemoteWorkDir(darwin) = %q; want ~/work", got)
	}
}

func TestResolveRemoteDestPath(t *testing.T) {
	winDest := ResolveRemoteDestPath("windows", "gitmap.json")
	if !strings.Contains(winDest, "gitmap.json") || !strings.Contains(winDest, "work") {
		t.Errorf("ResolveRemoteDestPath(windows) = %q; want path with work and gitmap.json", winDest)
	}
	unixDest := ResolveRemoteDestPath("linux", "gitmap.json")
	if unixDest != "~/work/gitmap.json" {
		t.Errorf("ResolveRemoteDestPath(linux) = %q; want ~/work/gitmap.json", unixDest)
	}
}

func TestBuildRemoteExecString(t *testing.T) {
	opts := NodesCloneOptions{
		Kind:     CloneKindCFR,
		PassArgs: []string{"gitmap.json", "--dry-run"},
		HasFile:  true,
	}
	winCmd := buildRemoteExecString(opts, "gitmap.json", true)
	if !strings.Contains(winCmd, "Set-Location D:\\work") || !strings.Contains(winCmd, "gitmap cfr") {
		t.Errorf("buildRemoteExecString(win) = %q; want Set-Location and gitmap cfr", winCmd)
	}

	unixCmd := buildRemoteExecString(opts, "gitmap.json", false)
	if !strings.Contains(unixCmd, "cd ~/work") || !strings.Contains(unixCmd, "gitmap cfr") {
		t.Errorf("buildRemoteExecString(unix) = %q; want cd ~/work and gitmap cfr", unixCmd)
	}

	optsNoFile := NodesCloneOptions{
		Kind:     CloneKindClone,
		PassArgs: []string{"user/repo"},
		HasFile:  false,
	}
	noFileCmd := buildRemoteExecString(optsNoFile, "", true)
	if noFileCmd != "gitmap clone user/repo" {
		t.Errorf("buildRemoteExecString(noFile) = %q; want gitmap clone user/repo", noFileCmd)
	}
}

func TestPrintNodesCloneHelp(t *testing.T) {
	if err := PrintNodesCloneHelp(CloneKindClone); err != nil {
		t.Errorf("PrintNodesCloneHelp(clone) error = %v", err)
	}
	if err := PrintNodesCloneHelp(CloneKindCFR); err != nil {
		t.Errorf("PrintNodesCloneHelp(cfr) error = %v", err)
	}
	if err := PrintNodesCloneHelp(CloneKindCFRP); err != nil {
		t.Errorf("PrintNodesCloneHelp(cfrp) error = %v", err)
	}
}
