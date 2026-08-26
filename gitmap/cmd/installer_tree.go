// Package cmd — installer_tree.go renders installer composition tree views.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// InstallerTreeNode represents a node in the installer tree hierarchy.
type InstallerTreeNode struct {
	Title       string
	Description string
	Children    []InstallerTreeNode
}

// printInstallSummaryHeader displays the summary header with a green checkmark.
func printInstallSummaryHeader(slug string) {
	fmt.Printf("\n  %s✓%s Installation Summary: %s%s%s\n",
		constants.ColorGreen,
		constants.ColorReset,
		constants.ColorCyan,
		slug,
		constants.ColorReset,
	)
}

// printInstallerTree renders a recursive tree node hierarchy to stdout.
func printInstallerTree(root InstallerTreeNode, prefix string, isLast bool) {
	if prefix == "" {
		fmt.Printf("  %s%s%s %s%s%s\n",
			constants.ColorWhite, root.Title, constants.ColorReset,
			constants.ColorDim, root.Description, constants.ColorReset)
		printInstallerChildren(root, "  ")
		return
	}
	printInstallerChildren(root, prefix)
}

// printInstallerChildren iterates over children of an installer tree node.
func printInstallerChildren(node InstallerTreeNode, prefix string) {
	childCount := len(node.Children)
	for index, childNode := range node.Children {
		isChildLast := index == childCount-1
		renderInstallerChild(childNode, prefix, isChildLast)
	}
}

// renderInstallerChild formats and prints an individual child node.
func renderInstallerChild(childNode InstallerTreeNode, prefix string, isChildLast bool) {
	connector := constants.TreeBranch
	nextPrefix := prefix + constants.TreePipe
	if isChildLast {
		connector = constants.TreeCorner
		nextPrefix = prefix + constants.TreeSpace
	}
	fmt.Printf("%s%s%s%s %s%s%s %s%s%s\n",
		prefix, constants.ColorCyan, connector, constants.ColorReset,
		constants.ColorWhite, childNode.Title, constants.ColorReset,
		constants.ColorDim, childNode.Description, constants.ColorReset)
	if len(childNode.Children) > 0 {
		printInstallerChildren(childNode, nextPrefix)
	}
}
