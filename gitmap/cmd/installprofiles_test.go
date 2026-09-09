package cmd

import (
	"testing"
)

func TestInstallProfilesRegistry(t *testing.T) {
	profiles := AllInstallProfiles()
	if len(profiles) < 7 {
		t.Fatalf("expected at least 7 profiles, got %d", len(profiles))
	}
	expectedProfiles := []string{"minimal", "dev", "ubuntu", "ubuntu-dev-ai", "ai", "backend", "fullstack"}
	for _, expected := range expectedProfiles {
		assertProfilePresent(t, profiles, expected)
	}
}

func assertProfilePresent(t *testing.T, profiles []InstallProfile, name string) {
	t.Helper()
	p, found := FindInstallProfile(name)
	if !found {
		t.Errorf("expected profile '%s' to be found", name)

		return
	}
	if len(p.Tools) == 0 {
		t.Errorf("profile '%s' has 0 tools", name)
	}
}

func TestInstallProfileAliases(t *testing.T) {
	aliases := map[string]string{
		"ubuntu-dev": "ubuntu",
		"linux-dev":  "ubuntu",
		"developer":  "dev",
		"min":        "minimal",
		"ml":         "ai",
		"server":     "backend",
		"web":        "fullstack",
	}
	for alias, canonical := range aliases {
		p, found := FindInstallProfile(alias)
		if !found {
			t.Errorf("expected alias '%s' to resolve", alias)
			continue
		}
		if p.Name != canonical {
			t.Errorf("alias '%s' resolved to '%s', expected '%s'", alias, p.Name, canonical)
		}
	}
}

func TestIsInstallProfile(t *testing.T) {
	if !IsInstallProfile("dev") {
		t.Errorf("expected 'dev' to be recognized as profile")
	}
	if !IsInstallProfile("ubuntu") {
		t.Errorf("expected 'ubuntu' to be recognized as profile")
	}
	if !IsInstallProfile("ai") {
		t.Errorf("expected 'ai' to be recognized as profile")
	}
	if IsInstallProfile("non-existent-profile-xyz") {
		t.Errorf("expected unknown string to not be a profile")
	}
}

func TestProfileDotAndBadge(t *testing.T) {
	dotComplete := formatProfileDot(5, 5)
	if dotComplete != installedDot {
		t.Errorf("expected installedDot for complete profile, got %s", dotComplete)
	}
	dotIncomplete := formatProfileDot(3, 5)
	if dotIncomplete != missingDot {
		t.Errorf("expected missingDot for incomplete profile, got %s", dotIncomplete)
	}
	progress := formatProfileProgress(3, 5)
	if progress != "[3/5 tools]" {
		t.Errorf("unexpected progress badge: %s", progress)
	}
}

func TestAgyInstallCommandRegistered(t *testing.T) {
	found := false
	for _, sub := range AgyCmd.Commands() {
		if sub.Name() == "install" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'install' subcommand to be registered in AgyCmd")
	}
}
