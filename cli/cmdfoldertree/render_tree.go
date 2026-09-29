package cmdfoldertree

import (
	"fmt"
	"strings"
)

// RenderTree returns the formatted emoji tree string representation.
func RenderTree(root *FolderTreeNode, opts FolderTreeOptions) string {
	var sb strings.Builder
	sb.WriteString(formatRootHeader(root, opts))
	sb.WriteString("\n")
	renderChildrenTree(&sb, root.Children, "", opts)
	return sb.String()
}

func formatRootHeader(root *FolderTreeNode, opts FolderTreeOptions) string {
	icon := "📂"
	gitTag := formatGitTag(root)
	tagPart := ""
	if gitTag != "" {
		tagPart = " " + gitTag
	}
	if opts.ShowNumbers {
		return fmt.Sprintf("%d. %s %s%s", root.Sequence, icon, root.Path, tagPart)
	}
	return fmt.Sprintf("%s %s%s", icon, root.Path, tagPart)
}

func renderChildrenTree(sb *strings.Builder, children []*FolderTreeNode, indent string, opts FolderTreeOptions) {
	for i, child := range children {
		isLast := i == len(children)-1
		branch := "├── "
		nextIndent := indent + "│   "
		if isLast {
			branch = "└── "
			nextIndent = indent + "    "
		}
		sb.WriteString(indent + branch + formatTreeNode(child, opts) + "\n")
		if len(child.Children) > 0 {
			renderChildrenTree(sb, child.Children, nextIndent, opts)
		}
	}
}

func formatTreeNode(node *FolderTreeNode, opts FolderTreeOptions) string {
	icon := resolveNodeIcon(node)
	gitTag := formatGitTag(node)
	numPrefix := ""
	if opts.ShowNumbers {
		numPrefix = fmt.Sprintf("%d. ", node.Sequence)
	}
	if gitTag != "" {
		return fmt.Sprintf("%s%s %s %s", numPrefix, icon, node.Name, gitTag)
	}
	return fmt.Sprintf("%s%s %s", numPrefix, icon, node.Name)
}

func formatGitTag(node *FolderTreeNode) string {
	if !node.IsGit {
		return ""
	}
	if node.GitRepo != "" && node.GitBranch != "" && node.GitRepo != node.Name {
		return fmt.Sprintf("[git: %s (%s)]", node.GitRepo, node.GitBranch)
	}
	if node.GitBranch != "" {
		return fmt.Sprintf("[git: %s]", node.GitBranch)
	}
	if node.GitRepo != "" {
		return fmt.Sprintf("[git: %s]", node.GitRepo)
	}
	return "[git]"
}
