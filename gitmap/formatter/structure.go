// Package formatter — structure.go renders a folder-tree Markdown file.
package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// WriteStructure writes a Markdown folder tree of discovered repos.
func WriteStructure(w io.Writer, records []model.ScanRecord) error {
	fmt.Fprintln(w, constants.StructureTitle)
	fmt.Fprintln(w)
	fmt.Fprintln(w, constants.StructureDescription)
	fmt.Fprintln(w)
	paths := collectPaths(records)
	tree := buildTree(paths)

	return renderTree(w, tree, "")
}

// pathEntry pairs a relative path with its record.
type pathEntry struct {
	Path   string
	Branch string
	URL    string
}

// collectPaths extracts and sorts paths from records.
func collectPaths(records []model.ScanRecord) []pathEntry {
	entries := make([]pathEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, pathEntry{
			Path: r.RelativePath, Branch: r.Branch, URL: r.HTTPSUrl,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	return entries
}

// treeNode represents a folder or repo in the tree.
type treeNode struct {
	Name     string
	Children []*treeNode
	IsRepo   bool
	Branch   string
	URL      string
}

// buildTree constructs a tree from sorted path entries.
func buildTree(entries []pathEntry) *treeNode {
	root := &treeNode{Name: "."}
	for _, e := range entries {
		insertPath(root, e)
	}

	return root
}

// insertPath adds a single path into the tree.
func insertPath(root *treeNode, entry pathEntry) {
	normalized := strings.ReplaceAll(entry.Path, "\\", "/")
	parts := strings.Split(normalized, "/")
	current := root
	for i, part := range parts {
		child := findChild(current, part)
		if child == nil {
			child = &treeNode{Name: part}
			current.Children = append(current.Children, child)
		}
		if i == len(parts)-1 {
			child.IsRepo = true
			child.Branch = entry.Branch
			child.URL = entry.URL
		}
		current = child
	}
}

// findChild looks for an existing child node by name.
func findChild(node *treeNode, name string) *treeNode {
	for _, c := range node.Children {
		if c.Name == name {
			return c
		}
	}

	return nil
}

// renderTree writes the tree as indented Markdown lines.
func renderTree(w io.Writer, node *treeNode, prefix string) error {
	for i, child := range node.Children {
		connector := constants.TreeBranch
		nextPrefix := prefix + constants.TreePipe
		if i == len(node.Children)-1 {
			connector = constants.TreeCorner
			nextPrefix = prefix + constants.TreeSpace
		}
		err := renderNode(w, child, prefix, connector)
		if err != nil {
			return err
		}
		if len(child.Children) > 0 {
			err = renderTree(w, child, nextPrefix)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// renderNode writes a single tree node line.
func renderNode(w io.Writer, node *treeNode, prefix, connector string) error {
	label := node.Name
	if node.IsRepo {
		label = fmt.Sprintf(constants.StructureRepoFmt,
			node.Name, node.Branch, node.URL)
	}
	_, err := fmt.Fprintf(w, "%s%s %s\n", prefix, connector, label)

	return err
}
