package cmdinstaller

import (
	"testing"
)

func TestInstallerPin_CLI_Helpers(t *testing.T) {
	// Test pinning
	err := PinVersion("test-tool", "v2.5.0")
	if err != nil {
		t.Fatalf("unexpected error pinning test-tool: %v", err)
	}

	ver, err := GetPinnedVersion("test-tool")
	if err != nil {
		t.Fatalf("unexpected error getting pinned version: %v", err)
	}
	if ver != "v2.5.0" {
		t.Fatalf("expected v2.5.0, got: %s", ver)
	}

	pins, errList := ListPinnedVersions()
	if errList != nil {
		t.Fatalf("unexpected error listing pins: %v", errList)
	}
	if pins["test-tool"] != "v2.5.0" {
		t.Fatalf("expected test-tool in pins map, got: %+v", pins)
	}

	// Test unpinning
	err = UnpinVersion("test-tool")
	if err != nil {
		t.Fatalf("unexpected error unpinning test-tool: %v", err)
	}

	ver, err = GetPinnedVersion("test-tool")
	if err != nil {
		t.Fatalf("unexpected error getting unpinned version: %v", err)
	}
	if ver != "" {
		t.Fatalf("expected empty version after unpin, got: %s", ver)
	}
}

func TestResolvePinnedVersion_Fallback(t *testing.T) {
	pins := map[string]string{
		"docker": "24.0.5",
		"node":   "v20.10.0",
	}

	ver, isPinned := resolvePinnedVersion("docker", "Docker CLI", pins)
	if !isPinned || ver != "24.0.5" {
		t.Fatalf("expected 24.0.5 pinned, got %s, %v", ver, isPinned)
	}

	ver, isPinned = resolvePinnedVersion("unknown", "unregistered", pins)
	if isPinned || ver != "-" {
		t.Fatalf("expected - not pinned, got %s, %v", ver, isPinned)
	}
}
