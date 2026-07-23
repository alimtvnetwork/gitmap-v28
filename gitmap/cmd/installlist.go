// Package cmd — installlist.go renders `gitmap install --list` grouped by
// category with a per-tool installed-status indicator.
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// printInstallListGrouped renders the supported tools grouped by category,
// annotating each row with installed status from the InstalledTool table
// and a PATH probe fallback. Falls back to a flat list if categories are
// undefined.
func printInstallListGrouped() {
	if len(constants.InstallToolCategories) == 0 {
		printInstallListFlat()

		return
	}

	installed := loadInstalledLookup()
	categories := sortedCategoryNames()

	fmt.Print(constants.MsgInstallListHeader)

	for _, cat := range categories {
		printCategoryBlock(cat, constants.InstallToolCategories[cat], installed)
	}

	fmt.Print(constants.MsgInstallListLegend)
}

// printCategoryBlock prints a single category header and its tool rows.
func printCategoryBlock(category string, tools []string, installed map[string]string) {
	underline := strings.Repeat("─", len(category))
	fmt.Printf(constants.MsgInstallListCategory, category, underline)

	for _, tool := range tools {
		printToolRow(tool, installed)
	}
}

// printToolRow prints one tool row with status + version + description.
func printToolRow(tool string, installed map[string]string) {
	desc := constants.InstallToolDescriptions[tool]
	status, version := resolveToolStatus(tool, installed)
	fmt.Printf(constants.MsgInstallListGrouped, status, tool, version, desc)
}

// resolveToolStatus returns (statusGlyph, version) for a tool, preferring
// the recorded DB version, falling back to PATH probe, then unknown.
func resolveToolStatus(tool string, installed map[string]string) (string, string) {
	if version, ok := installed[tool]; ok {
		return constants.StatusInstalled, version
	}
	if isCommandAvailable(tool) {
		return constants.StatusInstalled, "—"
	}

	return constants.StatusNotInstalled, "—"
}

// loadInstalledLookup builds a {tool → version} map from InstalledTool. Any
// DB error degrades to an empty map (PATH probe still runs).
func loadInstalledLookup() map[string]string {
	out := map[string]string{}

	db, err := openDB()
	if err != nil {
		return out
	}
	defer db.Close()

	tools, err := db.ListInstalledTools()
	if err != nil {
		return out
	}

	for _, t := range tools {
		out[t.Tool] = pickDisplayVersion(t)
	}

	return out
}

// pickDisplayVersion prefers the parsed version string, then the raw
// VersionString field, defaulting to a dash.
func pickDisplayVersion(t store.InstalledTool) string {
	if t.VersionString != "" && t.VersionString != "0.0.0" {
		return t.VersionString
	}

	return "—"
}

// sortedCategoryNames returns category keys with Core first, then alpha.
func sortedCategoryNames() []string {
	names := make([]string, 0, len(constants.InstallToolCategories))
	for name := range constants.InstallToolCategories {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		if names[i] == constants.ToolCategoryCore {
			return true
		}
		if names[j] == constants.ToolCategoryCore {
			return false
		}

		return names[i] < names[j]
	})

	return names
}

// printInstallListFlat is the legacy fallback used when categories are empty.
func printInstallListFlat() {
	fmt.Print(constants.MsgInstallListHeader)

	for tool, desc := range constants.InstallToolDescriptions {
		fmt.Printf(constants.MsgInstallListRow, tool, desc)
	}
}
