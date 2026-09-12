package cmdinstall

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderProfileTree(p InstallProfile, installed map[string]string) {
	fmt.Printf("\nProfile: %s (%s)\n", p.Name, p.Title)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Total Tools: %d\n", len(p.Tools))
	renderProfileTreeNodes(p, installed)
	fmt.Println()
}

func renderProfileTreeNodes(p InstallProfile, installed map[string]string) {
	toolCount := len(p.Tools)
	for idx, tool := range p.Tools {
		isLast := idx == toolCount-1
		row := formatProfileTreeNode(tool, isLast, installed)
		fmt.Println(row)
	}
}

func formatProfileTreeNode(tool string, isLast bool, installed map[string]string) string {
	connector := resolveTreeConnector(isLast)
	dot, ver := resolveTreeNodeStatus(tool, installed)
	desc := constants.InstallToolDescriptions[tool]

	return fmt.Sprintf("  %s %s %s %s %s", connector, dot, toolStyle.Render(tool), versionStyle.Render(ver), descStyle.Render(desc))
}

func resolveTreeConnector(isLast bool) string {
	if isLast {
		return constants.ColorCyan + constants.TreeCorner + constants.ColorReset
	}

	return constants.ColorCyan + constants.TreeBranch + constants.ColorReset
}

func resolveTreeNodeStatus(tool string, installed map[string]string) (string, string) {
	status, ver := resolveToolStatus(tool, installed)
	if status == constants.StatusInstalled {
		return installedDot, ver
	}

	return missingDot, "—"
}

func renderAllProfilesTree(installed map[string]string) {
	for _, p := range AllInstallProfiles() {
		renderProfileTree(p, installed)
	}
}

func printProfileUsageExamples() {
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap install profile <name> [flags]")
	fmt.Println("  gitmap in <name> [flags]")
	fmt.Println("\nExamples:")
	fmt.Println("  $ gitmap install profile dev")
	fmt.Println("  $ gitmap install profile dev --tree")
	fmt.Println("  $ gitmap in ubuntu")
	fmt.Println("  $ gitmap in dev --tree")
}
