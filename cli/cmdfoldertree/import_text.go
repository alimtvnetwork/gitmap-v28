package cmdfoldertree

import (
	"fmt"
	"path/filepath"
	"strings"
)

func parseTextPathsPayload(data []byte) (*FolderTreeNode, error) {
	lines := strings.Split(string(data), "\n")
	var paths []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		cleaned := stripTreeDecorations(trimmed)
		if cleaned != "" {
			paths = append(paths, cleaned)
		}
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no valid paths found in text input")
	}
	return buildTreeFromPathList(paths), nil
}

func stripTreeDecorations(line string) string {
	res := line
	decorations := []string{"├── ", "└── ", "│   ", "|-- ", "`-- ", "|   ", "📂 ", "📁 ", "📦 ", "💻 ", "📜 ", "📝 ", "🖼️ ", "📄 "}
	for _, d := range decorations {
		res = strings.ReplaceAll(res, d, "")
	}
	res = strings.TrimSpace(res)
	if idx := strings.Index(res, "[git"); idx != -1 {
		res = strings.TrimSpace(res[:idx])
	}
	return res
}

func buildTreeFromPathList(paths []string) *FolderTreeNode {
	root := &FolderTreeNode{
		Name:    ".",
		RelPath: ".",
		IsDir:   true,
	}
	seq := 1
	root.Sequence = seq
	seq++
	for _, p := range paths {
		addPathToTree(root, p, &seq)
	}
	return root
}

func addPathToTree(root *FolderTreeNode, p string, seq *int) {
	clean := filepath.ToSlash(filepath.Clean(p))
	if clean == "." || clean == "" {
		return
	}
	parts := strings.Split(clean, "/")
	isDir := strings.HasSuffix(p, "/") || strings.HasSuffix(p, "\\")
	curr := root
	for i, part := range parts {
		isLast := i == len(parts)-1
		childIsDir := isDir || !isLast
		curr = findOrCreateChild(curr, part, childIsDir, seq)
	}
}

func findOrCreateChild(parent *FolderTreeNode, name string, isDir bool, seq *int) *FolderTreeNode {
	for _, ch := range parent.Children {
		if ch.Name == name {
			return ch
		}
	}
	child := &FolderTreeNode{
		Sequence: *seq,
		Name:     name,
		IsDir:    isDir,
	}
	*seq++
	parent.Children = append(parent.Children, child)
	return child
}
