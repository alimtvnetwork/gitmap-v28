package cmdssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseOffendingKnownHostsLine(t *testing.T) {
	cases := []struct {
		input    string
		wantPath string
		wantLine int
		wantOk   bool
	}{
		{
			input:    "Offending ECDSA key in C:\\Users\\Administrator/.ssh/known_hosts:15",
			wantPath: "C:\\Users\\Administrator/.ssh/known_hosts",
			wantLine: 15,
			wantOk:   true,
		},
		{
			input:    "Offending key in /home/user/.ssh/known_hosts:42",
			wantPath: "/home/user/.ssh/known_hosts",
			wantLine: 42,
			wantOk:   true,
		},
		{
			input:    "No offending key here",
			wantPath: "",
			wantLine: 0,
			wantOk:   false,
		},
	}

	for _, c := range cases {
		gotPath, gotLine, gotOk := ParseOffendingKnownHostsLine(c.input)
		if gotOk != c.wantOk || gotPath != c.wantPath || gotLine != c.wantLine {
			t.Errorf("ParseOffendingKnownHostsLine(%q) = (%q, %d, %v), want (%q, %d, %v)",
				c.input, gotPath, gotLine, gotOk, c.wantPath, c.wantLine, c.wantOk)
		}
	}
}

func TestIsHostKeyChangedError(t *testing.T) {
	if !isHostKeyChangedError("WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!") {
		t.Error("expected true for identification changed")
	}
	if !isHostKeyChangedError("Host key verification failed.") {
		t.Error("expected true for verification failed")
	}
	if !isHostKeyChangedError("Host key for 192.168.10.155 has changed and you have requested strict checking.") {
		t.Error("expected true for host key changed")
	}
	if isHostKeyChangedError("Connection refused") {
		t.Error("expected false for connection refused")
	}
}

func TestRemoveKnownHostsLine(t *testing.T) {
	dir := t.TempDir()
	khPath := filepath.Join(dir, "known_hosts")
	content := "line1\nline2\nline3\n"
	if err := os.WriteFile(khPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test known_hosts: %v", err)
	}

	if err := RemoveKnownHostsLine(khPath, 2); err != nil {
		t.Fatalf("RemoveKnownHostsLine failed: %v", err)
	}

	data, err := os.ReadFile(khPath)
	if err != nil {
		t.Fatalf("failed to read test known_hosts: %v", err)
	}

	got := strings.TrimSpace(string(data))
	want := "line1\nline3"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildPruneTargets(t *testing.T) {
	targets22 := buildPruneTargets("192.168.1.10", 22)
	if len(targets22) < 1 || targets22[0] != "192.168.1.10" {
		t.Errorf("unexpected targets for port 22: %v", targets22)
	}

	targetsCustom := buildPruneTargets("192.168.1.10", 2222)
	hasCustom := false
	for _, tgt := range targetsCustom {
		if tgt == "[192.168.1.10]:2222" {
			hasCustom = true
		}
	}
	if !hasCustom {
		t.Errorf("expected custom port target in %v", targetsCustom)
	}
}
