package cmdinstall

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	catStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3ddc84")).Underline(true).MarginTop(1).MarginBottom(1)
	installedDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("●")
	missingDot   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render("○")
	unknownDot   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c")).Render("?")
	toolStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Width(22)
	versionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Width(14)
	descStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
)

func printInstallListGrouped() {
	if len(constants.InstallToolCategories) == 0 {
		printInstallListFlat()

		return
	}

	installed := loadInstalledLookup()
	printInstallListHeader()
	for _, cat := range sortedCategoryNames() {
		printCategoryBlock(cat, constants.InstallToolCategories[cat], installed)
	}

	printCustomInstallersSection(installed)
	printInstallProfilesSection(installed)
	printInstallListLegend()
}

func printInstallListHeader() {
	title := fmt.Sprintf("Gitmap Supported Tools & Packages (%s):", constants.Version)
	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9")).Render(title))
	recentTags := loadRecentReleaseTags()
	if recentTags != "" {
		tagLine := fmt.Sprintf("Recent Releases: %s", recentTags)
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Italic(true).Render(tagLine))
	}

	fmt.Println()
}

func loadRecentReleaseTags() string {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "tag", "-l", "v*", "--sort=-v:refname")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return formatRecentTags(string(out))
}

func formatRecentTags(raw string) string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var tags []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}

		if len(tags) >= 4 {
			break
		}
	}

	return strings.Join(tags, ", ")
}

func printInstallListLegend() {
	legend := fmt.Sprintf("\nLegend: %s installed   %s not installed   %s unknown\n", installedDot, missingDot, unknownDot)
	fmt.Println(lipgloss.NewStyle().Italic(true).Render(legend))
}

func printCategoryBlock(category string, tools []string, installed map[string]string) {
	fmt.Println(catStyle.Render(category))
	for _, tool := range tools {
		printToolRow(tool, installed)
	}
}

func printToolRow(tool string, installed map[string]string) {
	desc := constants.InstallToolDescriptions[tool]
	status, version := resolveToolStatus(tool, installed)

	var dot string
	switch status {
	case constants.StatusInstalled:
		dot = installedDot
	case constants.StatusNotInstalled:
		dot = missingDot
	default:
		dot = unknownDot
	}

	if version == "" {
		version = "-"
	}

	row := fmt.Sprintf("  %s %s %s %s", dot, toolStyle.Render(tool), versionStyle.Render(version), descStyle.Render(desc))
	fmt.Println(row)
}

func sortedCategoryNames() []string {
	var cats []string
	for cat := range constants.InstallToolCategories {
		if cat != constants.ToolCategoryCore {
			cats = append(cats, cat)
		}
	}

	sort.Strings(cats)

	return append([]string{constants.ToolCategoryCore}, cats...)
}

func printInstallListFlat() {
	var tools []string
	for tool := range constants.InstallToolDescriptions {
		tools = append(tools, tool)
	}

	sort.Strings(tools)

	installed := loadInstalledLookup()
	for _, tool := range tools {
		printToolRow(tool, installed)
	}
}

func loadInstalledLookup() map[string]string {
	out := make(map[string]string)
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return out
	}

	defer splitDB.Close()

	tools, err := splitDB.ListInstalledTools()
	if err != nil {
		return out
	}

	for _, t := range tools {
		out[t.Tool] = t.VersionString
	}

	return out
}

func resolveToolStatus(tool string, installed map[string]string) (string, string) {
	if tool == constants.ToolGitmap {
		return constants.StatusInstalled, constants.Version
	}

	if ver, ok := installed[tool]; ok && ver != "" && ver != "0.0.0" && ver != "found" && ver != "installed" {
		return constants.StatusInstalled, ver
	}

	bin, ver := resolveToolProbeCommand(tool)
	if bin == "" {
		return constants.StatusNotInstalled, "—"
	}

	if ver == "" {
		ver = "installed"
	}

	return constants.StatusInstalled, ver
}

func pickDisplayVersion(t store.InstalledTool) string {
	if t.VersionString == "" || t.VersionString == "0.0.0" {
		return "—"
	}

	return t.VersionString
}

func loadCustomInstallersList() []model.InstallerScript {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return nil
	}

	defer db.Close()
	if errMig := db.MigrateInstallers(); errMig != nil {
		return nil
	}

	list, errList := db.ListInstallers()
	if errList != nil {
		return nil
	}

	return list
}

func printCustomInstallersSection(installed map[string]string) {
	customList := loadCustomInstallersList()
	if len(customList) == 0 {
		return
	}

	fmt.Println(catStyle.Render("Custom Tools"))
	for _, script := range customList {
		printCustomInstallerRow(script, installed)
	}
}

func printCustomInstallerRow(script model.InstallerScript, installed map[string]string) {
	desc := script.Description
	status, ver := resolveCustomToolStatus(script, installed)
	dot := missingDot
	if status == constants.StatusInstalled {
		dot = installedDot
	}

	row := fmt.Sprintf("  %s %s %s %s", dot, toolStyle.Render(script.Slug), versionStyle.Render(ver), descStyle.Render(desc))
	fmt.Println(row)
}

func resolveCustomToolStatus(script model.InstallerScript, installed map[string]string) (string, string) {
	ver := script.Version
	if ver == "" {
		ver = "—"
	}

	if v, ok := installed[script.Slug]; ok && v != "" {
		return constants.StatusInstalled, v
	}

	if isBinaryInPath(script.Slug) {
		return constants.StatusInstalled, ver
	}

	return constants.StatusNotInstalled, ver
}
