package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfoldertree"
)

func dispatchFolderTree(command string) (bool, error) {
	if !isFolderTreeCommand(command) {
		return false, nil
	}
	args := os.Args[2:]
	return true, cmdfoldertree.RunFolderTree(args)
}

func isFolderTreeCommand(cmd string) bool {
	switch cmd {
	case "folder-tree", "foldertree", "ft":
		return true
	default:
		return false
	}
}
