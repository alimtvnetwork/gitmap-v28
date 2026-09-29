package cmdfoldertree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ScanFolderTree scans the directory hierarchy and returns the root node.
func ScanFolderTree(rootPath string, opts FolderTreeOptions) (*FolderTreeNode, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, err
	}
	seq := 1
	rootNode := createRootNode(absRoot, info, &seq)
	if rootNode.IsDir {
		walkSubTree(absRoot, rootNode, absRoot, 1, opts, &seq)
	}
	return rootNode, nil
}

func createRootNode(absRoot string, info os.FileInfo, seq *int) *FolderTreeNode {
	isGit, repo, branch := DetectGitInfo(absRoot)
	node := &FolderTreeNode{
		Sequence:  *seq,
		Name:      filepath.Base(absRoot),
		Path:      absRoot,
		RelPath:   ".",
		IsDir:     info.IsDir(),
		IsGit:     isGit,
		GitRepo:   repo,
		GitBranch: branch,
	}
	*seq++
	return node
}

func walkSubTree(baseRoot string, parent *FolderTreeNode, currentPath string, depth int, opts FolderTreeOptions, seq *int) {
	if isDepthExceeded(depth, opts.MaxDepth) {
		return
	}
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return
	}
	sortedEntries := sortDirEntries(entries)
	for _, entry := range sortedEntries {
		processDirEntry(baseRoot, parent, currentPath, entry, depth, opts, seq)
	}
}

func isDepthExceeded(depth, maxDepth int) bool {
	return maxDepth > 0 && depth > maxDepth
}

func sortDirEntries(entries []os.DirEntry) []os.DirEntry {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})
	return entries
}

func processDirEntry(baseRoot string, parent *FolderTreeNode, currentPath string, entry os.DirEntry, depth int, opts FolderTreeOptions, seq *int) {
	name := entry.Name()
	if shouldSkipEntry(name, opts.IncludeHidden) {
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

func shouldSkipEntry(name string, includeHidden bool) bool {
	if includeHidden {
		return false
	}
	if name == ".git" || name == ".DS_Store" || name == "node_modules" || name == "vendor" || name == ".turbo" {
		return true
	}
	return strings.HasPrefix(name, ".") && len(name) > 1
}
