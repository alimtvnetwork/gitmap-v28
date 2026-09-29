package cmdfoldertree

import (
	"fmt"
	"path/filepath"
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
	if opts.ShowNumbers {
		return fmt.Sprintf("%d. %s %s %s", root.Sequence, icon, root.Path, gitTag)
	}
	return fmt.Sprintf("%s %s %s", icon, root.Path, gitTag)
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
	if node.GitBranch != "" {
		return fmt.Sprintf("[git: %s]", node.GitBranch)
	}
	return "[git]"
}

func resolveNodeIcon(node *FolderTreeNode) string {
	if node.IsGit {
		return "📦"
	}
	if node.IsDir {
		return "📁"
	}
	return resolveFileIcon(node.Name)
}

func resolveFileIcon(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go", ".ts", ".js", ".py", ".rs", ".cs", ".php", ".c", ".cpp":
		return "💻"
	case ".json", ".yaml", ".yml", ".toml", ".ini", ".env", ".xml":
		return "📜"
	case ".md", ".txt", ".rst", ".doc", ".pdf":
		return "📝"
	case ".png", ".jpg", ".jpeg", ".svg", ".ico", ".webp", ".gif":
		return "🖼️"
	default:
		return "📄"
	}
}
