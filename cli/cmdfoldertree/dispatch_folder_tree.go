package cmdfoldertree

import (
	"os"
	"strings"
)

// DispatchFolderTree routes folder-tree, tree-search, and tree-learn commands.
func DispatchFolderTree(command string) (bool, error) {
	if isFolderTreeCommand(command) {
		args := os.Args[2:]
		return true, RunFolderTree(args)
	}

	if isTreeSearchCommand(command) {
		args := os.Args[2:]
		return true, RunTreeSearchDispatch(command, args)
	}

	if isTreeLearnCommand(command) {
		args := os.Args[2:]
		return true, RunTreeLearn(args)
	}

	return false, nil
}

func isFolderTreeCommand(cmd string) bool {
	switch strings.ToLower(cmd) {
	case "folder-tree", "foldertree", "ft":
		return true
	default:
		return false
	}
}

func isTreeSearchCommand(cmd string) bool {
	switch strings.ToLower(cmd) {
	case "tree-search", "treesearch", "ts",
		"tree-search-startswith", "tree-search-starts-with", "tree-search-prefix", "tss",
		"tree-search-contains", "tree-search-contain", "tree-search-substr", "tsc",
		"tree-search-endswith", "tree-search-ends-with", "tree-search-suffix", "tse",
		"tree-search-grep", "tree-search-regex", "tsg":
		return true
	default:
		return false
	}
}

func isTreeLearnCommand(cmd string) bool {
	switch strings.ToLower(cmd) {
	case "tree-learn", "treelearn", "tl":
		return true
	default:
		return false
	}
}
