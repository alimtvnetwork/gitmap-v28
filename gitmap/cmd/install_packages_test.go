package cmd

import (
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type pkgTestCase struct {
	manager string
	tool    string
	want    string
}

func getChocoTestCases() []pkgTestCase {
	return []pkgTestCase{
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
	}
}

func getWingetTestCases() []pkgTestCase {
	return []pkgTestCase{
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
	}
}

func getAptTestCases() []pkgTestCase {
	return []pkgTestCase{
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
	}
}

func getBrewTestCases() []pkgTestCase {
	return []pkgTestCase{
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
	}
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
