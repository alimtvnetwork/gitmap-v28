package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestNewToolsPackageResolution(t *testing.T) {
	cases := []struct {
		manager string
		tool    string
		want    string
	}{
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
	}

	for _, tc := range cases {
		got := resolvePackageName(tc.manager, tc.tool)
		if got != tc.want {
			t.Errorf("resolvePackageName(%q, %q) = %q; want %q", tc.manager, tc.tool, got, tc.want)
		}
	}
}

func TestResolveToolAlias(t *testing.T) {
	aliasCases := []struct {
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
		{"unknown-tool-xyz", "unknown-tool-xyz"},
	}

	for _, tc := range aliasCases {
		got := resolveToolAlias(tc.input)
		if got != tc.want {
			t.Errorf("resolveToolAlias(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestNewToolsCategoriesAndDescriptions(t *testing.T) {
	tools := []string{
		constants.ToolRust, constants.ToolDotnet, constants.ToolJava, constants.ToolFlutter,
		constants.ToolOllama, constants.ToolLlamaCpp, constants.ToolPythonLibs,
		constants.ToolDocker, constants.ToolKubernetes, constants.ToolJenkins,
		constants.ToolZsh, constants.ToolFlameshot, constants.ToolConemu, constants.ToolVLC,
	}

	for _, tool := range tools {
		desc, exists := constants.InstallToolDescriptions[tool]
		if !exists || desc == "" {
			t.Errorf("missing description for tool %q", tool)
		}
	}

	expectedCategories := []string{
		constants.ToolCategoryLanguages,
		constants.ToolCategoryAI,
		constants.ToolCategoryDevOps,
		constants.ToolCategoryUtilities,
	}

	for _, cat := range expectedCategories {
		toolsInCat, exists := constants.InstallToolCategories[cat]
		if !exists || len(toolsInCat) == 0 {
			t.Errorf("category %q is missing or empty in InstallToolCategories", cat)
		}
	}
}
