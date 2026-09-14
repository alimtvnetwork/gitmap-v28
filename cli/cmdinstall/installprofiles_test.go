package cmdinstall

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var allExpectedProfiles = []string{
	"minimal", "base", "git-compact", "advance", "cpp-dx", "small-dev",
	"dev", "dev-advance", "terminal", "web-dev", "devops",
	"ubuntu", "ai", "ai-tools", "antigravity-suite", "backend", "fullstack",
}

func TestInstallProfilesRegistry(t *testing.T) {
	profiles := AllInstallProfiles()
	if len(profiles) < 17 {
		t.Fatalf("expected at least 17 profiles, got %d", len(profiles))
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
	"ubuntu-dev":          "ubuntu",
	"linux-dev":           "ubuntu",
	"developer":           "dev",
	"dev-stack":           "dev",
	"min":                 "minimal",
	"basic":               "minimal",
	"ml":                  "ai",
	"ai-dev":              "ai",
	"llm":                 "ai",
	"server":              "backend",
	"back":                "backend",
	"web":                 "fullstack",
	"full":                "fullstack",
	"workstation":         "base",
	"daily":               "base",
	"gitcompact":          "git-compact",
	"profile-git":         "git-compact",
	"profile-git-compact": "git-compact",
	"advanced":            "advance",
	"cppdx":               "cpp-dx",
	"directx":             "cpp-dx",
	"smalldev":            "small-dev",
	"slim-dev":            "small-dev",
	"simple-dev":          "small-dev",
	"simpledev":           "small-dev",
	"devadvance":          "dev-advance",
	"dev-plus":            "dev-advance",
	"term":                "terminal",
	"cli":                 "terminal",
	"terminal-profile":    "terminal",
	"profile-terminal":    "terminal",
	"terminalprofile":     "terminal",
	"terminal-essentials": "terminal",
	"webdev":              "web-dev",
	"frontend":            "web-dev",
	"infra":               "devops",
	"cloud":               "devops",
	"all-ai":              "ai-tools",
	"aitools":             "ai-tools",
	"ag-suite":            "antigravity-suite",
	"antigravitysuite":    "antigravity-suite",
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
	"dev", "ubuntu", "ai", "ai-tools", "antigravity-suite", "base", "git-compact",
	"advance", "cpp-dx", "small-dev", "dev-advance",
	"terminal", "terminal-profile", "profile-terminal",
	"web-dev", "devops",
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

func TestTerminalProfileComposition(t *testing.T) {
	comp, isFound := resolveProfileTree("terminal-profile")
	if !isFound {
		t.Fatalf("expected terminal-profile tree to be found")
	}
	if comp.Name != "terminal" {
		t.Errorf("expected name 'terminal', got '%s'", comp.Name)
	}
	if len(comp.Tools) == 0 {
		t.Errorf("expected terminal composition to contain tools")
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
	for _, sub := range cmdagy.AgyCmd.Commands() {
		if sub.Name() == "install" {
			isFound = true
			break
		}
	}

	if !isFound {
		t.Errorf("expected 'install' subcommand to be registered in AgyCmd")
	}
}

func TestGitIsNotProfile(t *testing.T) {
	if IsInstallProfile("git") {
		t.Errorf("expected 'git' NOT to be recognized as an install profile alias")
	}
}

func TestResolveProfileInstalledBadge(t *testing.T) {
	prof := InstallProfile{
		Name:  "test-prof",
		Tools: []string{"git"},
	}
	installedMap := map[string]string{"git": "2.40.0"}
	badge := resolveProfileInstalledBadge(prof, installedMap)
	if badge != " [✔ Already Installed]" {
		t.Errorf("expected badge ' [✔ Already Installed]', got %q", badge)
	}

	emptyMap := map[string]string{}
	badgeEmpty := resolveProfileInstalledBadge(prof, emptyMap)
	if badgeEmpty != "" {
		t.Errorf("expected empty badge, got %q", badgeEmpty)
	}
}

func TestLinuxSmallDevPartition(t *testing.T) {
	tools := resolveLinuxSmallDevTools()
	for _, tool := range tools {
		if isWindowsOnlyTool(tool) {
			t.Errorf("expected no Windows tool in Linux small-dev, found %s", tool)
		}
	}
	assertLinuxSmallDevTools(t, tools)
}

func assertLinuxSmallDevTools(t *testing.T, tools []string) {
	t.Helper()
	assertContainsTool(t, tools, constants.ToolGitHubDesktop)
	assertContainsTool(t, tools, constants.ToolGitCompact)
	assertContainsTool(t, tools, constants.ToolGo)
	assertContainsTool(t, tools, constants.ToolRust)
	assertContainsTool(t, tools, constants.ToolPHP)
	assertContainsTool(t, tools, constants.ToolPython)
}

func TestLinuxDevPartition(t *testing.T) {
	tools := resolveLinuxDevTools()
	for _, tool := range tools {
		if isWindowsOnlyTool(tool) {
			t.Errorf("expected no Windows tool in Linux dev, found %s", tool)
		}
	}
	assertLinuxDevTools(t, tools)
}

func assertLinuxDevTools(t *testing.T, tools []string) {
	t.Helper()
	assertContainsTool(t, tools, constants.ToolNodeJS)
	assertContainsTool(t, tools, constants.ToolPnpm)
	assertContainsTool(t, tools, constants.ToolYarn)
	assertContainsTool(t, tools, constants.ToolAntigravity)
	assertContainsTool(t, tools, constants.ToolAgManager)
}

func TestLinuxTerminalProfileModernUtils(t *testing.T) {
	tools := resolveUnixTerminalTools()
	assertContainsTool(t, tools, constants.ToolJq)
	assertContainsTool(t, tools, constants.ToolYq)
	assertContainsTool(t, tools, constants.ToolZellij)
}

func TestSpecialProfileTrees(t *testing.T) {
	aiComp, hasAi := resolveProfileTree("ai-tools")
	if !hasAi || aiComp.Name != "ai-tools" {
		t.Errorf("expected ai-tools profile tree to resolve, got %v", hasAi)
	}

	agComp, hasAg := resolveProfileTree("antigravity-suite")
	if !hasAg || agComp.Name != "antigravity-suite" {
		t.Errorf("expected antigravity-suite profile tree to resolve, got %v", hasAg)
	}
}

func TestUbuntuTreeCompositions(t *testing.T) {
	sdev, hasSdev := resolveProfileTree("small-dev")
	if !hasSdev || len(sdev.Tools) != 6 {
		t.Errorf("expected small-dev tree to have 6 tools, got %d", len(sdev.Tools))
	}

	dev, hasDev := resolveProfileTree("dev")
	if !hasDev || len(dev.Tools) != 5 {
		t.Errorf("expected dev tree to have 5 tools, got %d", len(dev.Tools))
	}
}

func isWindowsOnlyTool(tool string) bool {
	return tool == constants.ToolConemu || tool == constants.ToolWinRAR || tool == constants.ToolNpp
}

func assertContainsTool(t *testing.T, tools []string, target string) {
	t.Helper()
	for _, tool := range tools {
		if tool == target {
			return
		}
	}
	t.Errorf("expected tools to contain %s", target)
}
