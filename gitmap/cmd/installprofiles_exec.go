package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runInstallProfile(profileName string, opts installOptions) error {
	p, found := FindInstallProfile(profileName)
	if !found {
		fmt.Printf("Unknown installation profile: %s\n", profileName)

		return nil
	}

	installed := loadInstalledLookup()
	if opts.Tree {
		renderProfileTree(p, installed)

		return nil
	}

	return executeProfileWorkflow(p, opts, installed)
}

func executeProfileWorkflow(p InstallProfile, opts installOptions, installed map[string]string) error {
	printProfileStartHeader(p, installed)
	installedCount := executeProfileTools(p, opts)
	printProfileSummary(p, installedCount)

	return nil
}

func printProfileStartHeader(p InstallProfile, installed map[string]string) {
	fmt.Printf("\n=== Installing Profile: %s (%s) ===\n", p.Name, p.Title)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Total Tools: %d\n\n", len(p.Tools))
	renderProfileTreeNodes(p, installed)
	fmt.Println()
}

func executeProfileTools(p InstallProfile, opts installOptions) int {
	installed := loadInstalledLookup()
	successCount := 0
	for idx, tool := range p.Tools {
		stepNum := idx + 1
		if isToolInstalledAlready(tool, installed) {
			announceToolAlreadyInstalled(stepNum, len(p.Tools), tool, installed)
			successCount++
			continue
		}

		installSingleProfileTool(stepNum, len(p.Tools), tool, opts)
		successCount++
	}

	return successCount
}

func isToolInstalledAlready(tool string, installed map[string]string) bool {
	status, _ := resolveToolStatus(tool, installed)

	return status == constants.StatusInstalled
}

func announceToolAlreadyInstalled(step, total int, tool string, installed map[string]string) {
	_, ver := resolveToolStatus(tool, installed)
	fmt.Printf("  ✓ [%d/%d] %s is already installed (%s)\n", step, total, tool, ver)
}

func installSingleProfileTool(step, total int, tool string, opts installOptions) {
	fmt.Printf("\n  → [%d/%d] Installing %s...\n", step, total, tool)
	toolOpts := opts
	toolOpts.Tool = tool
	toolOpts.Yes = true
	executeInstall(toolOpts)
}

func printProfileSummary(p InstallProfile, count int) {
	fmt.Printf("\n" + constants.ColorGreen + "✓" + constants.ColorReset)
	fmt.Printf(" Profile '%s' setup complete! (%d/%d tools processed)\n\n", p.Name, count, len(p.Tools))
}
