package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

func executeMove(db *store.DB, rec model.ScanRecord, destPath string, opts moveOpts) {
	if err := fsutil.SafeRename(rec.AbsolutePath, destPath); err != nil {
		fmt.Printf("mv: physical move failed: %v\n", err)
		return
	}

	newRepoName := filepath.Base(destPath)
	if err := updateRepoInDB(db, rec.ID, destPath, newRepoName); err != nil {
		fmt.Printf("mv: database update error: %v\n", err)
	}

	syncExternalMove(rec.AbsolutePath, destPath, newRepoName, opts)
	printMoveSuccess(rec.Slug, rec.AbsolutePath, destPath)
}

func syncExternalMove(oldPath, newPath, newName string, opts moveOpts) {
	if !opts.noVSCode {
		_ = vscodepm.UpdateRootPath(oldPath, newPath, newName)
	}
	if !opts.noDesktop {
		_ = desktop.UpdateRepoPath(oldPath, newPath)
	}
}

func printMoveSuccess(slug, oldPath, newPath string) {
	fmt.Println()
	fmt.Printf("  ✔ Moved %q successfully:\n", slug)
	fmt.Printf("     From: %s\n", oldPath)
	fmt.Printf("     To:   %s\n", newPath)
	fmt.Println()
}
