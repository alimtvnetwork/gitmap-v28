// Package formatter — terminaltree.go renders folder-tree structures for terminal output.
package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// termPathEntry holds path info for tree rendering.
type termPathEntry struct {
	Path   string
	Branch string
}

// collectTermPaths extracts and sorts paths from records.
func collectTermPaths(records []model.ScanRecord) []termPathEntry {
	entries := make([]termPathEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, termPathEntry{
			Path: r.RelativePath, Branch: r.Branch,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	return entries
}

// termNode represents a folder or repo in the terminal tree.
type termNode struct {
	Name     string
	Children []*termNode
	IsRepo   bool
	Branch   string
}

// buildTermTree constructs a tree from path entries.
func buildTermTree(entries []termPathEntry) *termNode {
	root := &termNode{Name: "."}
	for _, e := range entries {
		insertTermNode(root, e)
	}

	return root
}

// insertTermNode adds a path into the tree.
func insertTermNode(root *termNode, entry termPathEntry) {
	normalized := strings.ReplaceAll(entry.Path, "\\", "/")
	parts := strings.Split(normalized, "/")
	current := root
	for i, part := range parts {
		child := findTermChild(current, part)
		if child == nil {
			child = &termNode{Name: part}
			current.Children = append(current.Children, child)
		}
		if i == len(parts)-1 {
			child.IsRepo = true
			child.Branch = entry.Branch
		}
		current = child
	}
}

// findTermChild looks for a child node by name.
func findTermChild(node *termNode, name string) *termNode {
	for _, c := range node.Children {
		if c.Name == name {
			return c
		}
	}

	return nil
}

// renderTermTree writes the colored tree to the writer.
func renderTermTree(w io.Writer, node *termNode, prefix string) {
	for i, child := range node.Children {
		connector := constants.TreeBranch
		nextPrefix := prefix + constants.TreePipe
		if i == len(node.Children)-1 {
			connector = constants.TreeCorner
			nextPrefix = prefix + constants.TreeSpace
		}
		renderTermNode(w, child, prefix, connector)
		if len(child.Children) > 0 {
			renderTermTree(w, child, nextPrefix)
		}
	}
}

// renderTermNode writes a single colored tree node.
func renderTermNode(w io.Writer, node *termNode, prefix, connector string) {
	if node.IsRepo {
		fmt.Fprintf(w, "%s%s%s 📦 %s%s%s %s(%s)%s\n",
			constants.ColorDim, prefix, connector,
			constants.ColorGreen, node.Name, constants.ColorReset,
			constants.ColorDim, node.Branch, constants.ColorReset)

		return
	}
	fmt.Fprintf(w, "%s%s%s 📁 %s%s%s\n",
		constants.ColorDim, prefix, connector,
		constants.ColorYellow, node.Name, constants.ColorReset)
}
