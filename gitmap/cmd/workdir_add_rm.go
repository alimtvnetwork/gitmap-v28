// Package cmd — workdir_add_rm.go handles adding and removing work directories.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runWorkDirAdd(target, label string) error {
	if target == "" {
		return apperror.New("runWorkDirAdd", "E_INVALID_ARGS", map[string]any{"error": "path required"})
	}

	absPath, errAbs := filepath.Abs(target)
	if errAbs != nil {
		return errAbs
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	wd, errEnsure := db.EnsureWorkDir(absPath, label, false)
	if errEnsure != nil {
		return errEnsure
	}

	fmt.Printf("✓ Work directory registered: %s (ID: %d)\n", wd.AbsolutePath, wd.ID)
	return nil
}

func runWorkDirRm(target string) error {
	if target == "" {
		return apperror.New("runWorkDirRm", "E_INVALID_ARGS", map[string]any{"error": "path or ID required"})
	}

	absPath, _ := filepath.Abs(target)
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if err := db.DeleteWorkDir(absPath); err != nil {
		if errID := db.DeleteWorkDir(target); errID != nil {
			return errID
		}
	}

	fmt.Printf("✓ Work directory removed: %s\n", target)
	return nil
}
