// Package cmd — workdir_set_default.go sets the active default work directory.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runWorkDirSetDefault(target string) error {
	if target == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		target = cwd
	}

	absPath, _ := filepath.Abs(target)
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if err := db.SetDefaultWorkDir(absPath); err != nil {
		if errID := db.SetDefaultWorkDir(target); errID != nil {
			return errID
		}
	}

	fmt.Printf("✓ Default work directory set to: %s\n", target)
	return nil
}
