package cmd

import (
	"fmt"
	"os/exec"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/charmbracelet/lipgloss"
)

var (
	catStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3ddc84")).Underline(true).MarginTop(1).MarginBottom(1)
	installedDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("●")
	missingDot   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render("○")
	unknownDot   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c")).Render("?")
	toolStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Width(22)
	versionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Width(12)
	descStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
)

func printInstallListGrouped() {
	if len(constants.InstallToolCategories) == 0 {
		printInstallListFlat()
		return
	}

	installed := loadInstalledLookup()
	categories := sortedCategoryNames()

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9")).Render("Gitmap Supported Tools & Packages:"))

	for _, cat := range categories {
		printCategoryBlock(cat, constants.InstallToolCategories[cat], installed)
	}

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
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	return cats
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
	db, err := store.OpenDefault()
	if err != nil {
		return out
	}
	defer db.Close()

	tools, err := db.ListInstalledTools()
	if err != nil {
		return out
	}

	for _, t := range tools {
		out[t.Tool] = t.VersionString
	}
	return out
}

func resolveToolStatus(tool string, installed map[string]string) (string, string) {
	if ver, ok := installed[tool]; ok {
		return constants.StatusInstalled, ver
	}
	// Fallback to path check for generic tools if not tracked in DB
	if isBinaryInPath(tool) {
		return constants.StatusInstalled, "found"
	}
	return constants.StatusNotInstalled, "—"
}
func isBinaryInPath(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}
func pickDisplayVersion(t store.InstalledTool) string {
	if t.VersionString == "" || t.VersionString == "0.0.0" {
		return "—"
	}
	return t.VersionString
}
