package cmdfoldertree

import (
	"os"
	"path/filepath"
	"strings"
)

func processDirEntry(baseRoot string, parent *FolderTreeNode, currentPath string, entry os.DirEntry, depth int, opts FolderTreeOptions, seq *int) {
	name := entry.Name()
	if isIgnoredEntry(name, opts.IncludeHidden) {
		return
	}
	if !entry.IsDir() && opts.DirsOnly {
		return
	}
	fullPath := filepath.Join(currentPath, name)
	child := createChildNode(baseRoot, fullPath, name, entry.IsDir(), seq)
	parent.Children = append(parent.Children, child)
	if entry.IsDir() {
		walkSubTree(baseRoot, child, fullPath, depth+1, opts, seq)
	}
}

func createChildNode(baseRoot, fullPath, name string, isDir bool, seq *int) *FolderTreeNode {
	rel, _ := filepath.Rel(baseRoot, fullPath)
	isGit, repo, branch := false, "", ""
	if isDir {
		isGit, repo, branch = DetectGitInfo(fullPath)
	}
	node := &FolderTreeNode{
		Sequence:  *seq,
		Name:      name,
		Path:      fullPath,
		RelPath:   filepath.ToSlash(rel),
		IsDir:     isDir,
		IsGit:     isGit,
		GitRepo:   repo,
		GitBranch: branch,
	}
	*seq++
	return node
}

func isIgnoredEntry(name string, includeHidden bool) bool {
	if includeHidden {
		return false
	}
	if name == ".git" || name == ".DS_Store" || name == "node_modules" || name == "vendor" || name == ".turbo" {
		return true
	}
	return strings.HasPrefix(name, ".") && len(name) > 1
}
