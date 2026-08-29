// Package cmd — workdir_set_default.go sets the active default work directory.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runWorkDirSetDefault(target string) error {
	resolvedTarget, err := resolveWorkDirTarget(target)
	if err != nil {
		return err
	}

	absPath, _ := filepath.Abs(resolvedTarget)
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if err := setDefaultWorkDirByPathOrTarget(db, absPath, resolvedTarget); err != nil {
		return err
	}

	fmt.Printf("✓ Default work directory set to: %s\n", resolvedTarget)
	return nil
}

func resolveWorkDirTarget(target string) (string, error) {
	if target != "" {
		return target, nil
	}
	return os.Getwd()
}

func setDefaultWorkDirByPathOrTarget(db *store.DB, absPath, target string) error {
	if err := db.SetDefaultWorkDir(absPath); err == nil {
		return nil
	}
	return db.SetDefaultWorkDir(target)
}
