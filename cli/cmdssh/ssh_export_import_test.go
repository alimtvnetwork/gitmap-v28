package cmdssh

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseExportAllOptions(t *testing.T) {
	cases := []struct {
		args       []string
		wantTarget string
		wantDryRun bool
		wantForce  bool
	}{
		{[]string{}, "all", false, false},
		{[]string{"nodes"}, "all", false, false},
		{[]string{"all"}, "all", false, false},
		{[]string{"devbox", "--dry-run"}, "devbox", true, false},
		{[]string{"--target", "worker-1", "-f"}, "worker-1", false, true},
		{[]string{"nodes", "--force"}, "all", false, true},
	}

	for _, tc := range cases {
		opts := parseExportAllOptions(tc.args)
		if opts.Target != tc.wantTarget {
			t.Errorf("parseExportAllOptions(%v).Target = %q; want %q", tc.args, opts.Target, tc.wantTarget)
		}
		if opts.IsDryRun != tc.wantDryRun {
			t.Errorf("parseExportAllOptions(%v).IsDryRun = %v; want %v", tc.args, opts.IsDryRun, tc.wantDryRun)
		}
		if opts.IsForce != tc.wantForce {
			t.Errorf("parseExportAllOptions(%v).IsForce = %v; want %v", tc.args, opts.IsForce, tc.wantForce)
		}
	}
}

func TestParseImportAllOptions(t *testing.T) {
	cases := []struct {
		args       []string
		wantTarget string
		wantLocal  bool
		wantDryRun bool
	}{
		{[]string{"node", "u2"}, "u2", false, false},
		{[]string{"node", "192.168.1.5"}, "192.168.1.5", false, false},
		{[]string{"main"}, "main", false, false},
		{[]string{"--local-bundle"}, "", true, false},
		{[]string{"node", "u2", "--dry-run"}, "u2", false, true},
		{[]string{"-t", "node1"}, "node1", false, false},
	}

	for _, tc := range cases {
		opts := parseImportAllOptions(tc.args)
		if opts.Target != tc.wantTarget {
			t.Errorf("parseImportAllOptions(%v).Target = %q; want %q", tc.args, opts.Target, tc.wantTarget)
		}
		if opts.IsLocalBundle != tc.wantLocal {
			t.Errorf("parseImportAllOptions(%v).IsLocalBundle = %v; want %v", tc.args, opts.IsLocalBundle, tc.wantLocal)
		}
		if opts.IsDryRun != tc.wantDryRun {
			t.Errorf("parseImportAllOptions(%v).IsDryRun = %v; want %v", tc.args, opts.IsDryRun, tc.wantDryRun)
		}
	}
}

func TestGitmapExportBundleSerialization(t *testing.T) {
	bundle := GitmapExportBundle{
		Version:    "v6.273.0",
		ExportedAt: "2026-09-20T01:00:00Z",
		SourceNode: "test-node",
		Config:     map[string]any{"defaultMode": "ssh"},
		KnownHosts: "192.168.1.5 ssh-ed25519 AAAAC3NzaC1yc2EAAAADAQAB",
	}

	data, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded GitmapExportBundle
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.SourceNode != "test-node" || decoded.Version != "v6.273.0" {
		t.Errorf("bundle mismatch: %+v", decoded)
	}
}

func TestBuildRemoteWriteCmd(t *testing.T) {
	unixCmd := buildRemoteExportWriteCmd(".gitmap", "test.json", "dGVzdA==", false)
	if !strings.Contains(unixCmd, "mkdir -p") || !strings.Contains(unixCmd, "base64 -d") {
		t.Errorf("unexpected unixCmd: %s", unixCmd)
	}

	winCmd := buildRemoteExportWriteCmd(".gitmap", "test.json", "dGVzdA==", true)
	if !strings.Contains(winCmd, "powershell") || !strings.Contains(winCmd, "FromBase64String") {
		t.Errorf("unexpected winCmd: %s", winCmd)
	}
}

func TestDispatchExportImportSSH(t *testing.T) {
	exp := dispatchExportImportSSH("export-all", []string{"--help"})
	if !exp.IsMatched() {
		t.Errorf("expected export-all to match")
	}

	imp := dispatchExportImportSSH("import-all", []string{"--help"})
	if !imp.IsMatched() {
		t.Errorf("expected import-all to match")
	}

	unknown := dispatchExportImportSSH("unknown", []string{})
	if unknown.IsMatched() {
		t.Errorf("expected unknown to not match")
	}
}

func TestFormatExportMacroNames(t *testing.T) {
	short := formatExportMacroNames([]string{"a", "b"})
	if short != "a, b" {
		t.Errorf("expected 'a, b', got %q", short)
	}

	long := formatExportMacroNames([]string{"a", "b", "c", "d", "e"})
	if long != "a, b, c, +2 more" {
		t.Errorf("expected 'a, b, c, +2 more', got %q", long)
	}
}

func TestFilterNonEmptyLogs(t *testing.T) {
	input := []string{"log1", "", "log2", "", "log3"}
	filtered := filterNonEmptyLogs(input)
	if len(filtered) != 3 || filtered[0] != "log1" || filtered[1] != "log2" || filtered[2] != "log3" {
		t.Errorf("expected 3 items, got %v", filtered)
	}
}
