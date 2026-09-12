package cmdinstall

import (
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type pkgTestCase struct {
	manager string
	tool    string
	want    string
}

var chocoTestCases = []pkgTestCase{
	{constants.PkgMgrChocolatey, constants.ToolRust, constants.ChocoPkgRust},
	{constants.PkgMgrChocolatey, constants.ToolDotnet, constants.ChocoPkgDotnet},
	{constants.PkgMgrChocolatey, constants.ToolJava, constants.ChocoPkgJava},
	{constants.PkgMgrChocolatey, constants.ToolFlutter, constants.ChocoPkgFlutter},
	{constants.PkgMgrChocolatey, constants.ToolOllama, constants.ChocoPkgOllama},
	{constants.PkgMgrChocolatey, constants.ToolLlamaCpp, constants.ChocoPkgLlamaCpp},
	{constants.PkgMgrChocolatey, constants.ToolPythonLibs, constants.ChocoPkgPythonLibs},
	{constants.PkgMgrChocolatey, constants.ToolDocker, constants.ChocoPkgDocker},
	{constants.PkgMgrChocolatey, constants.ToolKubernetes, constants.ChocoPkgKubernetes},
	{constants.PkgMgrChocolatey, constants.ToolJenkins, constants.ChocoPkgJenkins},
	{constants.PkgMgrChocolatey, constants.ToolZsh, constants.ChocoPkgZsh},
	{constants.PkgMgrChocolatey, constants.ToolFlameshot, constants.ChocoPkgFlameshot},
	{constants.PkgMgrChocolatey, constants.ToolConemu, constants.ChocoPkgConemu},
	{constants.PkgMgrChocolatey, constants.ToolVLC, constants.ChocoPkgVLC},
	{constants.PkgMgrChocolatey, constants.ToolNginx, constants.ChocoPkgNginx},
	{constants.PkgMgrChocolatey, constants.ToolWordPress, constants.ChocoPkgWordPress},
	{constants.PkgMgrChocolatey, constants.ToolLaravel, constants.ChocoPkgLaravel},
	{constants.PkgMgrChocolatey, constants.ToolVMware, constants.ChocoPkgVMware},
	{constants.PkgMgrChocolatey, constants.ToolQBittorrent, constants.ChocoPkgQBittorrent},
	{constants.PkgMgrChocolatey, constants.ToolUTorrent, constants.ChocoPkgUTorrent},
	{constants.PkgMgrChocolatey, constants.ToolUbuntuFont, constants.ChocoPkgUbuntuFont},
	{constants.PkgMgrChocolatey, constants.ToolWhatsApp, constants.ChocoPkgWhatsApp},
	{constants.PkgMgrChocolatey, constants.ToolOneNote, constants.ChocoPkgOneNote},
	{constants.PkgMgrChocolatey, constants.ToolLightshot, constants.ChocoPkgLightshot},
	{constants.PkgMgrChocolatey, constants.ToolWindowsTerminal, constants.ChocoPkgWindowsTerminal},
	{constants.PkgMgrChocolatey, constants.ToolAria2, constants.ChocoPkgAria2},
	{constants.PkgMgrChocolatey, constants.Tool7Zip, constants.ChocoPkg7Zip},
	{constants.PkgMgrChocolatey, constants.ToolWinRAR, constants.ChocoPkgWinRAR},
	{constants.PkgMgrChocolatey, constants.ToolXMind, constants.ChocoPkgXMind},
	{constants.PkgMgrChocolatey, constants.ToolWordWeb, constants.ChocoPkgWordWeb},
	{constants.PkgMgrChocolatey, constants.ToolBeyondCompare, constants.ChocoPkgBeyondCompare},
	{constants.PkgMgrChocolatey, constants.ToolVcRedist, constants.ChocoPkgVcRedist},
	{constants.PkgMgrChocolatey, constants.ToolDirectX, constants.ChocoPkgDirectX},
	{constants.PkgMgrChocolatey, constants.ToolDirectXSdk, constants.ChocoPkgDirectXSdk},
	{constants.PkgMgrChocolatey, constants.ToolStarship, constants.ChocoPkgStarship},
	{constants.PkgMgrChocolatey, constants.ToolOhMyPosh, constants.ChocoPkgOhMyPosh},
	{constants.PkgMgrChocolatey, constants.ToolScoop, constants.ChocoPkgScoop},
}

func getChocoTestCases() []pkgTestCase {
	return chocoTestCases
}

var wingetTestCases = []pkgTestCase{
	{constants.PkgMgrWinget, constants.ToolRust, constants.WingetPkgRust},
	{constants.PkgMgrWinget, constants.ToolDotnet, constants.WingetPkgDotnet},
	{constants.PkgMgrWinget, constants.ToolJava, constants.WingetPkgJava},
	{constants.PkgMgrWinget, constants.ToolFlutter, constants.WingetPkgFlutter},
	{constants.PkgMgrWinget, constants.ToolOllama, constants.WingetPkgOllama},
	{constants.PkgMgrWinget, constants.ToolLlamaCpp, constants.WingetPkgLlamaCpp},
	{constants.PkgMgrWinget, constants.ToolPythonLibs, constants.WingetPkgPythonLibs},
	{constants.PkgMgrWinget, constants.ToolDocker, constants.WingetPkgDocker},
	{constants.PkgMgrWinget, constants.ToolKubernetes, constants.WingetPkgKubernetes},
	{constants.PkgMgrWinget, constants.ToolJenkins, constants.WingetPkgJenkins},
	{constants.PkgMgrWinget, constants.ToolZsh, constants.WingetPkgZsh},
	{constants.PkgMgrWinget, constants.ToolFlameshot, constants.WingetPkgFlameshot},
	{constants.PkgMgrWinget, constants.ToolConemu, constants.WingetPkgConemu},
	{constants.PkgMgrWinget, constants.ToolVLC, constants.WingetPkgVLC},
	{constants.PkgMgrWinget, constants.ToolNginx, constants.WingetPkgNginx},
	{constants.PkgMgrWinget, constants.ToolWordPress, constants.WingetPkgWordPress},
	{constants.PkgMgrWinget, constants.ToolLaravel, constants.WingetPkgLaravel},
	{constants.PkgMgrWinget, constants.ToolVMware, constants.WingetPkgVMware},
	{constants.PkgMgrWinget, constants.ToolQBittorrent, constants.WingetPkgQBittorrent},
	{constants.PkgMgrWinget, constants.ToolUTorrent, constants.WingetPkgUTorrent},
	{constants.PkgMgrWinget, constants.ToolUbuntuFont, constants.WingetPkgUbuntuFont},
	{constants.PkgMgrWinget, constants.ToolWhatsApp, constants.WingetPkgWhatsApp},
	{constants.PkgMgrWinget, constants.ToolOneNote, constants.WingetPkgOneNote},
	{constants.PkgMgrWinget, constants.ToolLightshot, constants.WingetPkgLightshot},
	{constants.PkgMgrWinget, constants.ToolWindowsTerminal, constants.WingetPkgWindowsTerminal},
	{constants.PkgMgrWinget, constants.ToolAria2, constants.WingetPkgAria2},
	{constants.PkgMgrWinget, constants.Tool7Zip, constants.WingetPkg7Zip},
	{constants.PkgMgrWinget, constants.ToolWinRAR, constants.WingetPkgWinRAR},
	{constants.PkgMgrWinget, constants.ToolXMind, constants.WingetPkgXMind},
	{constants.PkgMgrWinget, constants.ToolWordWeb, constants.WingetPkgWordWeb},
	{constants.PkgMgrWinget, constants.ToolBeyondCompare, constants.WingetPkgBeyondCompare},
	{constants.PkgMgrWinget, constants.ToolVcRedist, constants.WingetPkgVcRedist},
	{constants.PkgMgrWinget, constants.ToolDirectX, constants.WingetPkgDirectX},
	{constants.PkgMgrWinget, constants.ToolDirectXSdk, constants.WingetPkgDirectXSdk},
	{constants.PkgMgrWinget, constants.ToolStarship, constants.WingetPkgStarship},
	{constants.PkgMgrWinget, constants.ToolOhMyPosh, constants.WingetPkgOhMyPosh},
	{constants.PkgMgrWinget, constants.ToolScoop, constants.WingetPkgScoop},
}

func getWingetTestCases() []pkgTestCase {
	return wingetTestCases
}

var aptTestCases = []pkgTestCase{
	{constants.PkgMgrApt, constants.ToolRust, constants.AptPkgRust},
	{constants.PkgMgrApt, constants.ToolDotnet, constants.AptPkgDotnet},
	{constants.PkgMgrApt, constants.ToolJava, constants.AptPkgJava},
	{constants.PkgMgrApt, constants.ToolFlutter, constants.AptPkgFlutter},
	{constants.PkgMgrApt, constants.ToolOllama, constants.AptPkgOllama},
	{constants.PkgMgrApt, constants.ToolLlamaCpp, constants.AptPkgLlamaCpp},
	{constants.PkgMgrApt, constants.ToolPythonLibs, constants.AptPkgPythonLibs},
	{constants.PkgMgrApt, constants.ToolDocker, constants.AptPkgDocker},
	{constants.PkgMgrApt, constants.ToolKubernetes, constants.AptPkgKubernetes},
	{constants.PkgMgrApt, constants.ToolJenkins, constants.AptPkgJenkins},
	{constants.PkgMgrApt, constants.ToolZsh, constants.AptPkgZsh},
	{constants.PkgMgrApt, constants.ToolFlameshot, constants.AptPkgFlameshot},
	{constants.PkgMgrApt, constants.ToolConemu, constants.AptPkgConemu},
	{constants.PkgMgrApt, constants.ToolVLC, constants.AptPkgVLC},
	{constants.PkgMgrApt, constants.ToolNginx, constants.AptPkgNginx},
	{constants.PkgMgrApt, constants.ToolWordPress, constants.AptPkgWordPress},
	{constants.PkgMgrApt, constants.ToolLaravel, constants.AptPkgLaravel},
	{constants.PkgMgrApt, constants.ToolVMware, constants.AptPkgVMware},
	{constants.PkgMgrApt, constants.ToolQBittorrent, constants.AptPkgQBittorrent},
	{constants.PkgMgrApt, constants.ToolUTorrent, constants.AptPkgUTorrent},
	{constants.PkgMgrApt, constants.ToolOBS, constants.AptPkgOBS},
	{constants.PkgMgrApt, constants.ToolDbeaver, constants.AptPkgDbeaver},
	{constants.PkgMgrApt, constants.ToolUbuntuFont, constants.AptPkgUbuntuFont},
	{constants.PkgMgrApt, constants.ToolWhatsApp, constants.AptPkgWhatsApp},
	{constants.PkgMgrApt, constants.ToolOneNote, constants.AptPkgOneNote},
	{constants.PkgMgrApt, constants.ToolLightshot, constants.AptPkgLightshot},
	{constants.PkgMgrApt, constants.ToolWindowsTerminal, constants.AptPkgWindowsTerminal},
	{constants.PkgMgrApt, constants.ToolAria2, constants.AptPkgAria2},
	{constants.PkgMgrApt, constants.Tool7Zip, constants.AptPkg7Zip},
	{constants.PkgMgrApt, constants.ToolWinRAR, constants.AptPkgWinRAR},
	{constants.PkgMgrApt, constants.ToolXMind, constants.AptPkgXMind},
	{constants.PkgMgrApt, constants.ToolWordWeb, constants.AptPkgWordWeb},
	{constants.PkgMgrApt, constants.ToolBeyondCompare, constants.AptPkgBeyondCompare},
	{constants.PkgMgrApt, constants.ToolVcRedist, constants.AptPkgVcRedist},
	{constants.PkgMgrApt, constants.ToolDirectX, constants.AptPkgDirectX},
	{constants.PkgMgrApt, constants.ToolDirectXSdk, constants.AptPkgDirectXSdk},
	{constants.PkgMgrApt, constants.ToolStarship, constants.AptPkgStarship},
	{constants.PkgMgrApt, constants.ToolOhMyPosh, constants.AptPkgOhMyPosh},
	{constants.PkgMgrApt, constants.ToolScoop, constants.AptPkgScoop},
}

func getAptTestCases() []pkgTestCase {
	return aptTestCases
}

var brewTestCases = []pkgTestCase{
	{constants.PkgMgrBrew, constants.ToolRust, constants.BrewPkgRust},
	{constants.PkgMgrBrew, constants.ToolDotnet, constants.BrewPkgDotnet},
	{constants.PkgMgrBrew, constants.ToolJava, constants.BrewPkgJava},
	{constants.PkgMgrBrew, constants.ToolFlutter, constants.BrewPkgFlutter},
	{constants.PkgMgrBrew, constants.ToolOllama, constants.BrewPkgOllama},
	{constants.PkgMgrBrew, constants.ToolLlamaCpp, constants.BrewPkgLlamaCpp},
	{constants.PkgMgrBrew, constants.ToolPythonLibs, constants.BrewPkgPythonLibs},
	{constants.PkgMgrBrew, constants.ToolDocker, constants.BrewPkgDocker},
	{constants.PkgMgrBrew, constants.ToolKubernetes, constants.BrewPkgKubernetes},
	{constants.PkgMgrBrew, constants.ToolJenkins, constants.BrewPkgJenkins},
	{constants.PkgMgrBrew, constants.ToolZsh, constants.BrewPkgZsh},
	{constants.PkgMgrBrew, constants.ToolFlameshot, constants.BrewPkgFlameshot},
	{constants.PkgMgrBrew, constants.ToolConemu, constants.BrewPkgConemu},
	{constants.PkgMgrBrew, constants.ToolVLC, constants.BrewPkgVLC},
	{constants.PkgMgrBrew, constants.ToolNginx, constants.BrewPkgNginx},
	{constants.PkgMgrBrew, constants.ToolWordPress, constants.BrewPkgWordPress},
	{constants.PkgMgrBrew, constants.ToolLaravel, constants.BrewPkgLaravel},
	{constants.PkgMgrBrew, constants.ToolVMware, constants.BrewPkgVMware},
	{constants.PkgMgrBrew, constants.ToolQBittorrent, constants.BrewPkgQBittorrent},
	{constants.PkgMgrBrew, constants.ToolUTorrent, constants.BrewPkgUTorrent},
	{constants.PkgMgrBrew, constants.ToolUbuntuFont, constants.BrewPkgUbuntuFont},
	{constants.PkgMgrBrew, constants.ToolWhatsApp, constants.BrewPkgWhatsApp},
	{constants.PkgMgrBrew, constants.ToolOneNote, constants.BrewPkgOneNote},
	{constants.PkgMgrBrew, constants.ToolLightshot, constants.BrewPkgLightshot},
	{constants.PkgMgrBrew, constants.ToolWindowsTerminal, constants.BrewPkgWindowsTerminal},
	{constants.PkgMgrBrew, constants.ToolAria2, constants.BrewPkgAria2},
	{constants.PkgMgrBrew, constants.Tool7Zip, constants.BrewPkg7Zip},
	{constants.PkgMgrBrew, constants.ToolWinRAR, constants.BrewPkgWinRAR},
	{constants.PkgMgrBrew, constants.ToolXMind, constants.BrewPkgXMind},
	{constants.PkgMgrBrew, constants.ToolWordWeb, constants.BrewPkgWordWeb},
	{constants.PkgMgrBrew, constants.ToolBeyondCompare, constants.BrewPkgBeyondCompare},
	{constants.PkgMgrBrew, constants.ToolVcRedist, constants.BrewPkgVcRedist},
	{constants.PkgMgrBrew, constants.ToolDirectX, constants.BrewPkgDirectX},
	{constants.PkgMgrBrew, constants.ToolDirectXSdk, constants.BrewPkgDirectXSdk},
	{constants.PkgMgrBrew, constants.ToolStarship, constants.BrewPkgStarship},
	{constants.PkgMgrBrew, constants.ToolOhMyPosh, constants.BrewPkgOhMyPosh},
	{constants.PkgMgrBrew, constants.ToolScoop, constants.BrewPkgScoop},
}

func getBrewTestCases() []pkgTestCase {
	return brewTestCases
}

var snapTestCases = []pkgTestCase{
	{constants.PkgMgrSnap, constants.ToolOBS, constants.SnapPkgOBS},
	{constants.PkgMgrSnap, constants.ToolDbeaver, constants.SnapPkgDbeaver},
	{constants.PkgMgrSnap, constants.ToolUbuntuFont, constants.SnapPkgUbuntuFont},
	{constants.PkgMgrSnap, constants.ToolWhatsApp, constants.SnapPkgWhatsApp},
	{constants.PkgMgrSnap, constants.ToolOneNote, constants.SnapPkgOneNote},
	{constants.PkgMgrSnap, constants.ToolLightshot, constants.SnapPkgLightshot},
	{constants.PkgMgrSnap, constants.ToolWindowsTerminal, constants.SnapPkgWindowsTerminal},
	{constants.PkgMgrSnap, constants.ToolAria2, constants.SnapPkgAria2},
	{constants.PkgMgrSnap, constants.Tool7Zip, constants.SnapPkg7Zip},
	{constants.PkgMgrSnap, constants.ToolWinRAR, constants.SnapPkgWinRAR},
	{constants.PkgMgrSnap, constants.ToolXMind, constants.SnapPkgXMind},
	{constants.PkgMgrSnap, constants.ToolWordWeb, constants.SnapPkgWordWeb},
	{constants.PkgMgrSnap, constants.ToolBeyondCompare, constants.SnapPkgBeyondCompare},
	{constants.PkgMgrSnap, constants.ToolVcRedist, constants.SnapPkgVcRedist},
	{constants.PkgMgrSnap, constants.ToolDirectX, constants.SnapPkgDirectX},
	{constants.PkgMgrSnap, constants.ToolDirectXSdk, constants.SnapPkgDirectXSdk},
	{constants.PkgMgrSnap, constants.ToolStarship, constants.SnapPkgStarship},
	{constants.PkgMgrSnap, constants.ToolOhMyPosh, constants.SnapPkgOhMyPosh},
	{constants.PkgMgrSnap, constants.ToolScoop, constants.SnapPkgScoop},
}

func getSnapTestCases() []pkgTestCase {
	return snapTestCases
}

func checkPkgCases(t *testing.T, cases []pkgTestCase) {
	for _, tc := range cases {
		got := resolvePackageName(tc.manager, tc.tool)
		isMismatch := (got != tc.want)
		if isMismatch {
			t.Errorf("resolvePackageName(%q, %q) = %q; want %q", tc.manager, tc.tool, got, tc.want)
		}
	}
}

func TestNewToolsPackageResolution(t *testing.T) {
	checkPkgCases(t, getChocoTestCases())
	checkPkgCases(t, getWingetTestCases())
	checkPkgCases(t, getAptTestCases())
	checkPkgCases(t, getBrewTestCases())
	checkPkgCases(t, getSnapTestCases())
}

func TestSpecificPackageMappingsParity(t *testing.T) {
	if constants.AptPkgOBS != "obs-studio" {
		t.Errorf("expected AptPkgOBS 'obs-studio', got %q", constants.AptPkgOBS)
	}

	if constants.SnapPkgOBS != "obs-studio" {
		t.Errorf("expected SnapPkgOBS 'obs-studio', got %q", constants.SnapPkgOBS)
	}

	if constants.AptPkgDbeaver != "dbeaver-ce" {
		t.Errorf("expected AptPkgDbeaver 'dbeaver-ce', got %q", constants.AptPkgDbeaver)
	}

	if constants.SnapPkgDbeaver != "dbeaver-ce" {
		t.Errorf("expected SnapPkgDbeaver 'dbeaver-ce', got %q", constants.SnapPkgDbeaver)
	}
}

var aliasTestCases = []struct {
	input string
	want  string
}{
	{"k8s", constants.ToolKubernetes},
	{"kubectl", constants.ToolKubernetes},
	{"dotnet-sdk", constants.ToolDotnet},
	{"jdk", constants.ToolJava},
	{"openjdk", constants.ToolJava},
	{"cargo", constants.ToolRust},
	{"rustup", constants.ToolRust},
	{"llamacpp", constants.ToolLlamaCpp},
	{"ngx", constants.ToolNginx},
	{"engine-x", constants.ToolNginx},
	{"wp", constants.ToolWordPress},
	{"wp-cli", constants.ToolWordPress},
	{"wpcli", constants.ToolWordPress},
	{"artisan", constants.ToolLaravel},
	{"laravel-installer", constants.ToolLaravel},
	{"open-vm-tools", constants.ToolVMware},
	{"vmtools", constants.ToolVMware},
	{"vmware-tools", constants.ToolVMware},
	{"vm", constants.ToolVMware},
	{"qtorrent", constants.ToolQBittorrent},
	{"qbittorrent", constants.ToolQBittorrent},
	{"qbit", constants.ToolQBittorrent},
	{"utorrent", constants.ToolUTorrent},
	{"u-torrent", constants.ToolUTorrent},
	{"uttorrent", constants.ToolUTorrent},
	{"code", constants.ToolVSCode},
	{"vs-code", constants.ToolVSCode},
	{"nodejs", constants.ToolNodeJS},
	{"golang", constants.ToolGo},
	{"github-cli", constants.ToolGHCLI},
	{"c++", constants.ToolCPP},
	{"mingw", constants.ToolCPP},
	{"gcc", constants.ToolCPP},
	{"postgres", constants.ToolPostgreSQL},
	{"psql", constants.ToolPostgreSQL},
	{"mongo", constants.ToolMongoDB},
	{"pwsh", constants.ToolPowerShell},
	{"choco", constants.ToolChocolatey},
	{"notepad++", constants.ToolNpp},
	{"notepadpp", constants.ToolNpp},
	{"obs-studio", constants.ToolOBS},
	{"wt", constants.ToolWindowsTerminal},
	{"dart", constants.ToolFlutter},
	{"c#", constants.ToolDotnet},
	{"csharp", constants.ToolDotnet},
	{"local-llm", constants.ToolOllama},
	{"llm", constants.ToolOllama},
	{"llama", constants.ToolLlamaCpp},
	{"docker-desktop", constants.ToolDocker},
	{"starship", constants.ToolStarship},
	{"ohmyposh", constants.ToolOhMyPosh},
	{"omp", constants.ToolOhMyPosh},
	{"scoop", constants.ToolScoop},
	{"wa", constants.ToolWhatsApp},
	{"onenote", constants.ToolOneNote},
	{"unknown-tool-xyz", "unknown-tool-xyz"},
}

func TestResolveToolAlias(t *testing.T) {
	for _, tc := range aliasTestCases {
		got := resolveToolAlias(tc.input)
		isMismatch := (got != tc.want)
		if isMismatch {
			t.Errorf("resolveToolAlias(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

var newToolList = []string{
	constants.ToolRust, constants.ToolDotnet, constants.ToolJava, constants.ToolFlutter,
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

func checkToolDescriptions(t *testing.T, tools []string) {
	for _, tool := range tools {
		desc, isFound := constants.InstallToolDescriptions[tool]
		hasMissingDesc := (len(desc) == 0)
		isMissing := (!isFound || hasMissingDesc)
		if isMissing {
			t.Errorf("missing description for tool %q", tool)
		}
	}
}

func checkCategories(t *testing.T, categories []string) {
	for _, cat := range categories {
		toolsInCat, isFound := constants.InstallToolCategories[cat]
		hasEmptyTools := (len(toolsInCat) == 0)
		isMissing := (!isFound || hasEmptyTools)
		if isMissing {
			t.Errorf("category %q is missing or empty in InstallToolCategories", cat)
		}
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
		got := toolBinaryName(tc.tool)
		isMismatch := (got != tc.want)
		if isMismatch {
			t.Errorf("toolBinaryName(%q) = %q; want %q", tc.tool, got, tc.want)
		}
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
		got := versionFlag(tc.binary)
		isMismatch := (got != tc.want)
		if isMismatch {
			t.Errorf("versionFlag(%q) = %q; want %q", tc.binary, got, tc.want)
		}
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
