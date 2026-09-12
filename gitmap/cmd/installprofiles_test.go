package cmd

import (
	"testing"
)

var allExpectedProfiles = []string{
	"minimal", "base", "git-compact", "advance", "cpp-dx", "small-dev",
	"dev", "dev-advance", "terminal", "web-dev", "devops",
	"ubuntu", "ai", "backend", "fullstack",
}

func TestInstallProfilesRegistry(t *testing.T) {
	profiles := AllInstallProfiles()
	if len(profiles) < 15 {
		t.Fatalf("expected at least 15 profiles, got %d", len(profiles))
	}

	for _, expected := range allExpectedProfiles {
		assertProfilePresent(t, profiles, expected)
	}
}

func assertProfilePresent(t *testing.T, profiles []InstallProfile, name string) {
	t.Helper()
	p, isFound := FindInstallProfile(name)
	if !isFound {
		t.Errorf("expected profile '%s' to be found", name)

		return
	}

	if len(p.Tools) == 0 {
		t.Errorf("profile '%s' has 0 tools", name)
	}
}

var profileAliasesMap = map[string]string{
	"ubuntu-dev":  "ubuntu",
	"linux-dev":   "ubuntu",
	"developer":   "dev",
	"dev-stack":   "dev",
	"min":         "minimal",
	"basic":       "minimal",
	"ml":          "ai",
	"ai-dev":      "ai",
	"llm":         "ai",
	"server":      "backend",
	"back":        "backend",
	"web":         "fullstack",
	"full":        "fullstack",
	"workstation": "base",
	"daily":       "base",
	"git":         "git-compact",
	"gitcompact":  "git-compact",
	"advanced":    "advance",
	"cppdx":       "cpp-dx",
	"directx":     "cpp-dx",
	"smalldev":    "small-dev",
	"slim-dev":    "small-dev",
	"devadvance":  "dev-advance",
	"dev-plus":    "dev-advance",
	"term":        "terminal",
	"cli":         "terminal",
	"webdev":      "web-dev",
	"frontend":    "web-dev",
	"infra":       "devops",
	"cloud":       "devops",
}

func TestInstallProfileAliases(t *testing.T) {
	for alias, canonical := range profileAliasesMap {
		assertAliasResolves(t, alias, canonical)
	}
}

func assertAliasResolves(t *testing.T, alias, canonical string) {
	t.Helper()
	p, isFound := FindInstallProfile(alias)
	if !isFound {
		t.Errorf("expected alias '%s' to resolve", alias)

		return
	}

	if p.Name != canonical {
		t.Errorf("alias '%s' resolved to '%s', expected '%s'", alias, p.Name, canonical)
	}
}

var expectedKnownProfiles = []string{
	"dev", "ubuntu", "ai", "base", "git-compact",
	"advance", "cpp-dx", "small-dev", "dev-advance",
	"terminal", "web-dev", "devops",
}

func TestIsInstallProfile(t *testing.T) {
	for _, name := range expectedKnownProfiles {
		if !IsInstallProfile(name) {
			t.Errorf("expected %q to be recognized as profile", name)
		}
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
	isFound := false
	for _, sub := range AgyCmd.Commands() {
		if sub.Name() == "install" {
			isFound = true
			break
		}
	}

	if !isFound {
		t.Errorf("expected 'install' subcommand to be registered in AgyCmd")
	}
}
