package cmdfoldertree

import (
	"os"
)

func DispatchFolderTree(command string) (bool, error) {
	if !isFolderTreeCommand(command) {
		return false, nil
	}
	args := os.Args[2:]
	return true, RunFolderTree(args)
}

func isFolderTreeCommand(cmd string) bool {
	switch cmd {
	case "folder-tree", "foldertree", "ft":
		return true
	default:
		return false
	}
}
