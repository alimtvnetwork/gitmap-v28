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
		target = "."
	}

	absPath, errAbs := filepath.Abs(target)
	if errAbs != nil {
		return errAbs
	}

	return executeWorkDirAdd(absPath, label)
}

func executeWorkDirAdd(absPath, label string) error {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if existing, errFind := db.GetWorkDirByPath(absPath); errFind == nil && existing != nil {
		fmt.Printf("Work directory already added: %s (ID: %d)\n", existing.AbsolutePath, existing.ID)

		return nil
	}

	return registerNewWorkDir(db, absPath, label)
}

func registerNewWorkDir(db *store.DB, absPath, label string) error {
	all, errList := db.ListWorkDirs()
	isFirst := (errList == nil && len(all) == 0)

	wd, errEnsure := db.EnsureWorkDir(absPath, label, isFirst)
	if errEnsure != nil {
		return errEnsure
	}

	if isFirst {
		fmt.Printf("✓ Work directory registered and marked as default: %s (ID: %d)\n", wd.AbsolutePath, wd.ID)

		return nil
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

	if err := deleteWorkDirByPathOrTarget(db, absPath, target); err != nil {
		return err
	}

	fmt.Printf("✓ Work directory removed: %s\n", target)

	return nil
}

func deleteWorkDirByPathOrTarget(db *store.DB, absPath, target string) error {
	if err := db.DeleteWorkDir(absPath); err == nil {
		return nil
	}

	return db.DeleteWorkDir(target)
}
