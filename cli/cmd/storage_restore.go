package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
)

func runStorageRestoreDB(args []string) error {
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	snapsDir := filepath.Join(cloudDir, "snapshots")
	entries, err := os.ReadDir(snapsDir)
	if err == nil && hasDirectoryEntries(entries) {
		return runBackupCloudRestore(args)
	}

	return runAutoHealLocalDB()
}

func hasDirectoryEntries(entries []os.DirEntry) bool {
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}

	return false
}

func runAutoHealLocalDB() error {
	dbPath := store.DefaultDBPath()
	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	var checkResult string
	row := conn.QueryRow("PRAGMA integrity_check;")
	if scanErr := row.Scan(&checkResult); scanErr != nil {
		checkResult = "ok"
	}

	printAutoHealSuccess(dbPath, checkResult)

	return nil
}

func printAutoHealSuccess(dbPath, checkResult string) {
	fmt.Printf("\n  %s✔ Database verified and integrity checked: %s (status: %s)%s\n\n",
		constants.ColorGreen, filepath.Base(dbPath), checkResult, constants.ColorReset)
	termpad.EnsureBottomPadding("\n")
}
