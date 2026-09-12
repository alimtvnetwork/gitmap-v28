package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	profileNameStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Width(12)
	profileBadgeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Width(13)
)

// printInstallProfilesSection renders the Installation Profiles table.
func printInstallProfilesSection(installed map[string]string) {
	profiles := AllInstallProfiles()
	if len(profiles) == 0 {
		return
	}

	fmt.Println(catStyle.Render("Installation Profiles"))
	for _, p := range profiles {
		printSingleProfileRow(p, installed)
	}

	printProfileUsageExamples()
}

// printInstallProfilesOnly prints the profiles when `gitmap install profile` is invoked.
func printInstallProfilesOnly() {
	installed := loadInstalledLookup()
	printInstallProfilesSection(installed)
}

func printSingleProfileRow(p InstallProfile, installed map[string]string) {
	total := len(p.Tools)
	count := countProfileInstalledTools(p, installed)
	dot := formatProfileDot(count, total)
	badge := formatProfileProgress(count, total)
	desc := p.Description
	row := fmt.Sprintf("  %s %s %s %s", dot, profileNameStyle.Render(p.Name), profileBadgeStyle.Render(badge), descStyle.Render(desc))
	fmt.Println(row)
}

func countProfileInstalledTools(p InstallProfile, installed map[string]string) int {
	count := 0
	for _, tool := range p.Tools {
		status, _ := resolveToolStatus(tool, installed)
		if status == constants.StatusInstalled {
			count++
		}
	}

	return count
}

func formatProfileDot(count, total int) string {
	if count == total {
		return installedDot
	}

	return missingDot
}

func formatProfileProgress(count, total int) string {
	return fmt.Sprintf("[%d/%d tools]", count, total)
}
