package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestAntigravityDesktopDownloadUrl(t *testing.T) {
	winUrl := getAntigravityDesktopDownloadUrl("windows")
	if !strings.Contains(winUrl, "Antigravity-x64.exe") {
		t.Fatalf("expected windows URL to point to Antigravity-x64.exe, got %s", winUrl)
	}

	if !strings.Contains(winUrl, "antigravity.google") {
		t.Fatalf("expected windows URL from antigravity.google, got %s", winUrl)
	}

	linuxUrl := getAntigravityDesktopDownloadUrl("linux")
	if !strings.Contains(linuxUrl, "Antigravity.tar.gz") {
		t.Fatalf("expected linux URL to point to Antigravity.tar.gz, got %s", linuxUrl)
	}

	if !strings.Contains(linuxUrl, "antigravity.google") {
		t.Fatalf("expected linux URL from antigravity.google, got %s", linuxUrl)
	}
}

func TestAntigravityToolRoutingAndAliases(t *testing.T) {
	tests := []struct {
		alias    string
		expected string
	}{
		{"antigravity", constants.ToolAntigravity},
		{"antigravity-ide", constants.ToolAntigravity},
		{"antigravity-desktop", constants.ToolAntigravity},
		{"ag", constants.ToolAntigravity},
		{"agy", constants.ToolAgy},
		{"antigravity-cli", constants.ToolAgy},
	}

	for _, tc := range tests {
		actual := resolveToolAlias(tc.alias)
		if actual != tc.expected {
			t.Errorf("expected alias %q to resolve to %q, got %q", tc.alias, tc.expected, actual)
		}
	}
}

func TestAntigravityBinaryAndProbeMapping(t *testing.T) {
	if bin := toolBinaryName(constants.ToolAntigravity); bin != "antigravity" {
		t.Errorf("expected ToolAntigravity binary name 'antigravity', got %q", bin)
	}

	if bin := toolBinaryName(constants.ToolAgy); bin != "agy" {
		t.Errorf("expected ToolAgy binary name 'agy', got %q", bin)
	}

	cfgAntigravity, hasAntigravity := toolProbeMap[constants.ToolAntigravity]
	if !hasAntigravity || len(cfgAntigravity.bins) == 0 {
		t.Fatalf("expected ToolAntigravity in toolProbeMap")
	}

	cfgAgy, hasAgy := toolProbeMap[constants.ToolAgy]
	if !hasAgy || len(cfgAgy.bins) == 0 {
		t.Fatalf("expected ToolAgy in toolProbeMap")
	}
}
