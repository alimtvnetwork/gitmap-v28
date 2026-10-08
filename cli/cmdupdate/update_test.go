package cmdupdate

import (
	"encoding/json"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v6.512.0", "6.512.0"},
		{"V6.512.0", "6.512.0"},
		{"6.512.0", "6.512.0"},
		{"  v1.2.3  ", "1.2.3"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := NormalizeVersion(tt.input)
		if actual != tt.expected {
			t.Errorf("NormalizeVersion(%q) = %q, expected %q", tt.input, actual, tt.expected)
		}
	}
}

func TestFormatVersionTag(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"6.512.0", "v6.512.0"},
		{"v6.512.0", "v6.512.0"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := FormatVersionTag(tt.input)
		if actual != tt.expected {
			t.Errorf("FormatVersionTag(%q) = %q, expected %q", tt.input, actual, tt.expected)
		}
	}
}

func TestIsAlreadyUpdated(t *testing.T) {
	// Same version without force -> skip
	if !IsAlreadyUpdated("v6.512.0", "v6.512.0", false) {
		t.Errorf("expected true when versions match and force is false")
	}

	if !IsAlreadyUpdated("6.512.0", "v6.512.0", false) {
		t.Errorf("expected true when normalized versions match")
	}

	// Force flag overrides skip
	if IsAlreadyUpdated("v6.512.0", "v6.512.0", true) {
		t.Errorf("expected false when force is true")
	}

	// Different version -> do not skip
	if IsAlreadyUpdated("v6.512.0", "v6.513.0", false) {
		t.Errorf("expected false when versions differ")
	}

	// Empty target -> do not skip
	if IsAlreadyUpdated("v6.512.0", "", false) {
		t.Errorf("expected false when target is empty")
	}
}

func TestParseUpdateOptions(t *testing.T) {
	args := []string{"--version", "v6.515.0", "-f", "-j", "-q", "-n"}
	opts := ParseUpdateOptions(args)

	if opts.TargetVersion != "v6.515.0" {
		t.Errorf("expected TargetVersion = v6.515.0, got %s", opts.TargetVersion)
	}

	if !opts.IsForce {
		t.Errorf("expected IsForce = true")
	}

	if !opts.IsJSON {
		t.Errorf("expected IsJSON = true")
	}

	if !opts.IsQuiet {
		t.Errorf("expected IsQuiet = true")
	}

	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun = true")
	}

	// Positional version test
	posArgs := []string{"6.520.0"}
	posOpts := ParseUpdateOptions(posArgs)
	if posOpts.TargetVersion != "6.520.0" {
		t.Errorf("expected positional target 6.520.0, got %s", posOpts.TargetVersion)
	}
}

func TestMatchPlatformAsset(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		platform string
		arch     string
		expected bool
	}{
		{"gitmap-v28.0.0-windows-amd64.zip", "v28.0.0", "windows", "amd64", true},
		{"gitmap-v28.0.0-windows-arm64.zip", "v28.0.0", "windows", "arm64", true},
		{"gitmap-v28.0.0-windows-amd64.zip", "v28.0.0", "windows", "arm64", false},
		{"gitmap-v28.0.0-linux-amd64.tar.gz", "v28.0.0", "linux", "amd64", true},
		{"gitmap-v28.0.0-darwin-arm64.tar.gz", "v28.0.0", "darwin", "arm64", true},
		{"gitmap-v28.0.0-linux-amd64.tar.gz", "v28.0.0", "windows", "amd64", false},
	}

	for _, tt := range tests {
		actual := MatchPlatformAsset(tt.name, tt.tag, tt.platform, tt.arch)
		if actual != tt.expected {
			t.Errorf("MatchPlatformAsset(%q, %q, %q, %q) = %v, expected %v", tt.name, tt.tag, tt.platform, tt.arch, actual, tt.expected)
		}
	}
}

func TestAlreadyUpdatedJSONSerialization(t *testing.T) {
	res := AlreadyUpdatedResult{
		Status:         "already_updated",
		CurrentVersion: "v6.512.0",
		TargetVersion:  "v6.512.0",
		Updated:        false,
		Message:        "GitMap is already on the target version.",
	}

	data, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if m["status"] != "already_updated" {
		t.Errorf("expected status already_updated, got %v", m["status"])
	}
	if m["updated"] != false {
		t.Errorf("expected updated false, got %v", m["updated"])
	}
}
