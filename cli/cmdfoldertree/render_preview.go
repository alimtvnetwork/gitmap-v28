package cmdfoldertree

import (
	"fmt"
	"strings"
)

// RenderPreview returns the numbered one-line gap format matching the screenshot.
func RenderPreview(root *FolderTreeNode, opts FolderTreeOptions) string {
	nodes := collectNodesFlat(root, opts.DirsOnly)
	var sb strings.Builder
	for i, node := range nodes {
		gitTag := formatGitTag(node)
		nameLine := formatPreviewName(i+1, node.Name, gitTag, opts.ShowNumbers)
		sb.WriteString(nameLine + "\n")
		sb.WriteString(node.Path + "\n\n")
	}
	return sb.String()
}

func formatPreviewName(seq int, name, gitTag string, showNumbers bool) string {
	prefix := ""
	if showNumbers {
		prefix = fmt.Sprintf("%d. ", seq)
	}
	if gitTag != "" {
		return fmt.Sprintf("%s%s %s", prefix, name, gitTag)
	}
	return fmt.Sprintf("%s%s", prefix, name)
}

func collectNodesFlat(root *FolderTreeNode, dirsOnly bool) []*FolderTreeNode {
	var list []*FolderTreeNode
	if !dirsOnly || root.IsDir {
		list = append(list, root)
	}
	walkCollect(root.Children, dirsOnly, &list)
	return list
}

func walkCollect(children []*FolderTreeNode, dirsOnly bool, list *[]*FolderTreeNode) {
	for _, child := range children {
		if !dirsOnly || child.IsDir {
			*list = append(*list, child)
		}
		if len(child.Children) > 0 {
			walkCollect(child.Children, dirsOnly, list)
		}
	}
}
