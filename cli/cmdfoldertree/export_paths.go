package cmdfoldertree

import "strings"

func exportPaths(root *FolderTreeNode, dirsOnly bool) string {
	var lines []string
	walkCollectPaths(root, dirsOnly, &lines)
	return strings.Join(lines, "\n") + "\n"
}

func walkCollectPaths(node *FolderTreeNode, dirsOnly bool, lines *[]string) {
	if isIncludedInPathList(node, dirsOnly) {
		*lines = append(*lines, formatExportPath(node))
	}
	for _, child := range node.Children {
		walkCollectPaths(child, dirsOnly, lines)
	}
}

func formatExportPath(node *FolderTreeNode) string {
	rel := node.RelPath
	if node.IsDir && rel != "." && !strings.HasSuffix(rel, "/") {
		return rel + "/"
	}
	return rel
}

func isIncludedInPathList(node *FolderTreeNode, dirsOnly bool) bool {
	if dirsOnly {
		return node.IsDir
	}
	return !node.IsDir
}
