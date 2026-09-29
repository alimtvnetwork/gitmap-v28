package cmdfoldertree

import "strings"

func exportPaths(root *FolderTreeNode, dirsOnly bool) string {
	var lines []string
	walkCollectPaths(root, dirsOnly, &lines)
	return strings.Join(lines, "\n") + "\n"
}

func walkCollectPaths(node *FolderTreeNode, dirsOnly bool, lines *[]string) {
	if shouldIncludeInPathList(node, dirsOnly) {
		rel := node.RelPath
		if node.IsDir && rel != "." && !strings.HasSuffix(rel, "/") {
			rel += "/"
		}
		*lines = append(*lines, rel)
	}
	for _, child := range node.Children {
		walkCollectPaths(child, dirsOnly, lines)
	}
}

func shouldIncludeInPathList(node *FolderTreeNode, dirsOnly bool) bool {
	if dirsOnly {
		return node.IsDir
	}
	return !node.IsDir
}
