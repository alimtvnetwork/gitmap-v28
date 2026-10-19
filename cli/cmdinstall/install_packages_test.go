package cmdinstall

import (
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var newToolList = []string{
	constants.ToolRust, constants.ToolCargo, constants.ToolDotnet, constants.ToolJava, constants.ToolFlutter,
	constants.ToolOllama, constants.ToolLlamaCpp, constants.ToolPythonLibs,
	constants.ToolDocker, constants.ToolKubernetes, constants.ToolJenkins,
	constants.ToolZsh, constants.ToolFlameshot, constants.ToolConemu, constants.ToolVLC,
	constants.ToolNginx, constants.ToolWordPress, constants.ToolLaravel,
	constants.ToolVMware, constants.ToolQBittorrent, constants.ToolUTorrent,
	constants.ToolUbuntuFont, constants.ToolWhatsApp, constants.ToolOneNote,
	constants.ToolLightshot, constants.ToolWindowsTerminal, constants.ToolAria2,
	constants.Tool7Zip, constants.ToolWinRAR, constants.ToolXMind,
	constants.ToolWordWeb, constants.ToolBeyondCompare, constants.ToolVcRedist,
	constants.ToolDirectX, constants.ToolDirectXSdk, constants.ToolStarship,
	constants.ToolOhMyPosh, constants.ToolScoop,
}

var expectedCategoryList = []string{
	constants.ToolCategoryCore,
	constants.ToolCategoryLanguages,
	constants.ToolCategoryAI,
	constants.ToolCategoryDevOps,
	constants.ToolCategoryUtilities,
}

func TestNewToolsCategoriesAndDescriptions(t *testing.T) {
	checkToolDescriptions(t, newToolList)
	checkCategories(t, expectedCategoryList)
}

func verifySingleToolDescription(t *testing.T, tool string) {
	desc, isFound := constants.InstallToolDescriptions[tool]
	hasMissingDesc := (len(desc) == 0)
	isMissing := (!isFound || hasMissingDesc)
	if isMissing {
		t.Errorf("missing description for tool %q", tool)
	}
}

func checkToolDescriptions(t *testing.T, tools []string) {
	for _, tool := range tools {
		verifySingleToolDescription(t, tool)
	}
}

func verifySingleCategory(t *testing.T, cat string) {
	toolsInCat, isFound := constants.InstallToolCategories[cat]
	hasEmptyTools := (len(toolsInCat) == 0)
	isMissing := (!isFound || hasEmptyTools)
	if isMissing {
		t.Errorf("category %q is missing or empty in InstallToolCategories", cat)
	}
}

func checkCategories(t *testing.T, categories []string) {
	for _, cat := range categories {
		verifySingleCategory(t, cat)
	}
}

func verifySingleToolBinary(t *testing.T, tool, want string) {
	got := toolBinaryName(tool)
	if got != want {
		t.Errorf("toolBinaryName(%q) = %q; want %q", tool, got, want)
	}
}

func TestInstallVerifyToolBinaryMapping(t *testing.T) {
	binaryCases := []struct {
		tool string
		want string
	}{
		{constants.ToolNginx, "nginx"},
		{constants.ToolWordPress, "wp"},
		{constants.ToolLaravel, "laravel"},
	}

	for _, tc := range binaryCases {
		verifySingleToolBinary(t, tc.tool, tc.want)
	}
}

func verifySingleVersionFlag(t *testing.T, binary, want string) {
	got := versionFlag(binary)
	if got != want {
		t.Errorf("versionFlag(%q) = %q; want %q", binary, got, want)
	}
}

func TestInstallVerifyVersionFlag(t *testing.T) {
	flagCases := []struct {
		binary string
		want   string
	}{
		{"nginx", "-v"},
		{"wp", "--version"},
		{"laravel", "--version"},
		{"node", "--version"},
	}

	for _, tc := range flagCases {
		verifySingleVersionFlag(t, tc.binary, tc.want)
	}
}

func TestInstallVerifyExpectedExePathNginx(t *testing.T) {
	isNonWindows := (runtime.GOOS != "windows")
	if isNonWindows {
		t.Skip("skipping Windows expectedExePath test on non-windows platform")
	}

	got := expectedExePath(constants.ToolNginx)
	want := `C:\tools\nginx\nginx.exe`
	isMismatch := (got != want)
	if isMismatch {
		t.Errorf("expectedExePath(%q) = %q; want %q", constants.ToolNginx, got, want)
	}
}

func TestIsProtectedWorkDirectory(t *testing.T) {
	protectedPaths := []string{"", ".", "..", "d:/work/something", "/git-work/repo"}
	for _, p := range protectedPaths {
		if !isProtectedWorkDirectory(p) {
			t.Errorf("expected %q to be protected work directory", p)
		}
	}
}

func TestRemovePathTargetsSafetyGuard(t *testing.T) {
	err := removePathTargets([]string{"d:/work/safe"})
	if err == nil {
		t.Fatal("expected safety guard error for protected path, got nil")
	}
}

func TestRenderUninstallDryRun(t *testing.T) {
	plan := UninstallPlan{
		Tool:             "test-tool",
		PackageRemoval:   "test-pkg-rm",
		BinaryWrappers:   []string{"/usr/bin/test-tool"},
		DesktopLaunchers: []string{"test.desktop"},
		SystemIcons:      []string{"test.png"},
		DesktopDock:      "test.dock",
		Configuration:    []string{"~/.config/test"},
	}
	renderUninstallDryRun("test-tool", plan)
}
