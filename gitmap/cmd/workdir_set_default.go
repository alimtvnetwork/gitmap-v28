// Package cmd — workdir_set_default.go sets the active default work directory.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runWorkDirSetDefault(target string) error {
	if target == "" {
		return apperror.New("runWorkDirSetDefault", "E_INVALID_ARGS", map[string]any{"error": "target path or ID required"})
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
