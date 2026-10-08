package cmddir

import (
	"io/fs"
	"os"
)

func handleFindDirSkip(name string) error {
	if name == ".git" || name == ".tmp" || name == "node_modules" || name == ".gitmap" {
		return fs.SkipDir
	}

	return nil
}

func handleFindWalkErr(err error) error {
	if err == nil || os.IsNotExist(err) {
		return nil
	}

	return err
}
