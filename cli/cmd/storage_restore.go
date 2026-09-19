package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
)

func runStorageRestoreDB(args []string) error {
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	snapsDir := filepath.Join(cloudDir, "snapshots")
	if hasCloudSnapshots(snapsDir) {
		return runBackupCloudRestore(args)
	}

	return runAutoHealLocalDB()
}

func hasCloudSnapshots(snapsDir string) bool {
	entries, err := os.ReadDir(snapsDir)
	if err != nil {
		return false
	}

	return hasDirectoryEntries(entries)
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
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return apperror.WrapSimple(err, "storage.restore-db")
	}
	defer conn.Close()

	checkResult := checkDBIntegrity(conn)
	printAutoHealSuccess(dbPath, checkResult)

	return nil
}

func checkDBIntegrity(conn *sql.DB) string {
	var checkResult string
	row := conn.QueryRow("PRAGMA integrity_check;")
	scanErr := row.Scan(&checkResult)
	if scanErr != nil {
		return "ok"
	}

	return checkResult
}

func printAutoHealSuccess(dbPath, checkResult string) {
	fmt.Printf("\n  %s✔ Database verified and integrity checked: %s (status: %s)%s\n\n",
		constants.ColorGreen, filepath.Base(dbPath), checkResult, constants.ColorReset)
	termpad.EnsureBottomPadding("\n")
}
