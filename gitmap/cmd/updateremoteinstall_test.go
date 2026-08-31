package cmd

import (
	"runtime"
	"strings"
	"testing"
)

// TestInstallerURLFor verifies the raw.githubusercontent URL is composed
// from the owner constant, the supplied slug, and the platform-correct
// installer filename (install.ps1 on Windows, install.sh elsewhere).
func TestInstallerURLFor(t *testing.T) {
	t.Parallel()

	const slug = "gitmap-v28"
	got := installerURLFor(slug)

	if !strings.HasPrefix(got, "https://raw.githubusercontent.com/") {
		t.Errorf("expected raw.githubusercontent URL, got %q", got)
	}
	if !strings.Contains(got, "/"+slug+"/") {
		t.Errorf("expected slug %q in URL path, got %q", slug, got)
	}

	wantName := "install.sh"
	if runtime.GOOS == "windows" {
		wantName = "install.ps1"
	}
	if !strings.HasSuffix(got, "/"+wantName) {
		t.Errorf("expected suffix /%s, got %q", wantName, got)
	}
}

// TestInstallerURLForVariesByPlatform ensures the installer name segment
// matches GOOS — guards the Windows-vs-POSIX branch.
func TestInstallerURLForVariesByPlatform(t *testing.T) {
	t.Parallel()

	got := installerURLFor("any")
	hasSh := strings.HasSuffix(got, "/install.sh")
	hasPs := strings.HasSuffix(got, "/install.ps1")

	if !hasSh && !hasPs {
		t.Fatalf("URL must end with install.sh or install.ps1, got %q", got)
	}
	if hasSh && hasPs {
		t.Fatalf("URL cannot end with both installers: %q", got)
	}
}

func TestDecodeReleaseTagName(t *testing.T) {
	t.Parallel()

	jsonInput := `{"tag_name": "v6.158.0", "name": "Release v6.158.0"}`
	got := decodeReleaseTagName(strings.NewReader(jsonInput))
	if got != "6.158.0" {
		t.Errorf("expected 6.158.0, got %q", got)
	}
}

func TestDecodeVersionFromMap(t *testing.T) {
	t.Parallel()

	// Upper case Version
	upperJSON := `{"Version": "6.158.0"}`
	if got := decodeVersionFromMap(strings.NewReader(upperJSON)); got != "6.158.0" {
		t.Errorf("expected 6.158.0 for uppercase Version, got %q", got)
	}

	// Lower case version
	lowerJSON := `{"version": "6.158.0"}`
	if got := decodeVersionFromMap(strings.NewReader(lowerJSON)); got != "6.158.0" {
		t.Errorf("expected 6.158.0 for lowercase version, got %q", got)
	}
}
