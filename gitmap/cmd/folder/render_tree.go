package folder

import (
	"fmt"
	"sort"
	"strings"
)

// TreeNode represents an individual directory or file in a hierarchical folder tree.
type TreeNode struct {
	Name     string
	IsDir    bool
	Meta     *FileMeta
	Children map[string]*TreeNode
}

// NewTreeNode initializes an empty tree node.
func NewTreeNode(name string, isDir bool) *TreeNode {
	return &TreeNode{
		Name:     name,
		IsDir:    isDir,
		Children: make(map[string]*TreeNode),
	}
}

// BuildTree constructs a hierarchical TreeNode graph from a flat list of FileMeta items.
func BuildTree(rootName string, files []*FileMeta) *TreeNode {
	root := NewTreeNode(rootName, true)
	for _, f := range files {
		insertNode(root, f)
	}

	return root
}

func insertNode(root *TreeNode, meta *FileMeta) {
	cleanPath := strings.TrimPrefix(meta.Path, "./")
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	parts := strings.Split(cleanPath, "/")

	curr := root
	for i := 0; i < len(parts)-1; i++ {
		p := parts[i]
		if _, exists := curr.Children[p]; !exists {
			curr.Children[p] = NewTreeNode(p, true)
		}

		curr = curr.Children[p]
	}

	leafName := parts[len(parts)-1]
	leaf := NewTreeNode(leafName, false)
	leaf.Meta = meta
	curr.Children[leafName] = leaf
}

// RenderTree converts a TreeNode hierarchy into a formatted ASCII/Unicode directory tree.
func RenderTree(root *TreeNode, isDetailed bool) string {
	var sb strings.Builder
	sb.WriteString(root.Name + "/\n")
	renderTreeLevel(root, "", &sb, isDetailed)

	return sb.String()
}

func renderTreeLevel(node *TreeNode, prefix string, sb *strings.Builder, isDetailed bool) {
	children := sortTreeChildren(node.Children)
	count := len(children)

	for i, child := range children {
		isLast := (i == count-1)
		connector := "├── "
		subPrefix := "│   "
		if isLast {
			connector = "└── "
			subPrefix = "    "
		}

		if child.IsDir {
			sb.WriteString(fmt.Sprintf("%s%s%s/\n", prefix, connector, child.Name))
			renderTreeLevel(child, prefix+subPrefix, sb, isDetailed)
		} else {
			metaInfo := formatMetaSuffix(child.Meta, isDetailed)
			sb.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, connector, child.Name, metaInfo))
		}
	}
}

func formatMetaSuffix(meta *FileMeta, isDetailed bool) string {
	if !isDetailed || meta == nil {
		return ""
	}

	var parts []string
	if meta.Sequence > 0 {
		parts = append(parts, fmt.Sprintf("seq: %02d", meta.Sequence))
	}

	if !meta.IsBinary && meta.LinesOfCode > 0 {
		parts = append(parts, fmt.Sprintf("%d lines", meta.LinesOfCode))
	}

	if meta.SizeFormatted != "" {
		parts = append(parts, meta.SizeFormatted)
	}

	if len(parts) == 0 {
		return ""
	}

	return " (" + strings.Join(parts, ", ") + ")"
}

func sortTreeChildren(childMap map[string]*TreeNode) []*TreeNode {
	var dirs []*TreeNode
	var files []*TreeNode

	for _, node := range childMap {
		if node.IsDir {
			dirs = append(dirs, node)
		} else {
			files = append(files, node)
		}
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name < dirs[j].Name })
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	return append(dirs, files...)
}
