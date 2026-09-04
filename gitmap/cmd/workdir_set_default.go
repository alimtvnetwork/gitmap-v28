// Package cmd — workdir_set_default.go sets the active default work directory.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runWorkDirDefault(target string) error {
	if target != "" {
		return runWorkDirSetDefault(target)
	}

	return runWorkDirShowDefault()
}

func runWorkDirShowDefault() error {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	wd, err := db.GetDefaultWorkDir()
	if err != nil || wd == nil {
		fmt.Println("No default work directory configured.")
		fmt.Println("Run `gitmap workdir set <path>` to configure one.")

		return nil
	}

	fmt.Printf("✓ Default work directory: %s (ID: %d, Label: %s)\n", wd.AbsolutePath, wd.ID, wd.Label)

	return nil
}

func runWorkDirPath() error {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	wd, err := db.GetDefaultWorkDir()
	if err != nil || wd == nil {
		fmt.Println("No default work directory configured.")

		return nil
	}

	fmt.Println(wd.AbsolutePath)

	return nil
}

func runWorkDirSetDefault(target string) error {
	resolvedTarget, err := resolveWorkDirTarget(target)
	if err != nil {
		return err
	}

	absPath, errAbs := filepath.Abs(resolvedTarget)
	if errAbs != nil {
		absPath = resolvedTarget
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	_, _ = db.EnsureWorkDir(absPath, filepath.Base(absPath), true)

	if err := setDefaultWorkDirByPathOrTarget(db, absPath, resolvedTarget); err != nil {
		return err
	}

	fmt.Printf("✓ Default work directory set to: %s\n", absPath)

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
