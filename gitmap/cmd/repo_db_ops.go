package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/repodb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type repoSplitStats struct {
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	RepoSlug    string `json:"repoSlug"`
	RepoID      int64  `json:"repoId"`
	RepoFiles   int    `json:"repoFiles"`
	SearchCache int    `json:"searchCache"`
	FileSeqs    int    `json:"fileSeqs"`
	ScanLogs    int    `json:"scanLogs"`
}

func handleRepoDBStatus(args []string) error {
	db, path, repoID, slug, err := resolveCurrentRepoSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var stats repoSplitStats
	stats.Path = path
	stats.RepoSlug = slug
	stats.RepoID = repoID
	if info, err := os.Stat(path); err == nil {
		stats.Size = info.Size()
	}

	_ = db.QueryRow("SELECT COUNT(*) FROM RepoFile;").Scan(&stats.RepoFiles)
	_ = db.QueryRow("SELECT COUNT(*) FROM SearchCache;").Scan(&stats.SearchCache)
	_ = db.QueryRow("SELECT COUNT(*) FROM FileSequence;").Scan(&stats.FileSeqs)
	_ = db.QueryRow("SELECT COUNT(*) FROM RepoScanLog;").Scan(&stats.ScanLogs)

	if hasArgFlag(args, "--json") {
		return printJSON(stats)
	}

	fmt.Println(constants.ColorCyan + "● Repository Split Database Summary:" + constants.ColorReset)
	fmt.Printf("  • %-20s %s (ID: %d)\n", "Repository:", stats.RepoSlug, stats.RepoID)
	fmt.Printf("  • %-20s %s\n", "Database File:", stats.Path)
	fmt.Printf("  • %-20s %s\n", "File Size:", formatBytes(stats.Size))
	fmt.Printf("  • %-20s %d\n", "Indexed Files:", stats.RepoFiles)
	fmt.Printf("  • %-20s %d\n", "Cached Searches:", stats.SearchCache)
	fmt.Printf("  • %-20s %d\n", "File Sequences:", stats.FileSeqs)
	fmt.Printf("  • %-20s %d\n", "Scan/Sync Logs:", stats.ScanLogs)
	return nil
}

func handleRepoDBLog(args []string) error {
	db, _, _, slug, err := resolveCurrentRepoSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT Action, Status, COALESCE(Details, ''), CreatedAt FROM RepoScanLog ORDER BY Id DESC LIMIT 20;")
	if err != nil {
		fmt.Printf("No scan logs recorded in split database for %s.\n", slug)
		return nil
	}
	defer rows.Close()

	var logs []map[string]string
	for rows.Next() {
		var action, status, details, createdAt string
		if scanErr := rows.Scan(&action, &status, &details, &createdAt); scanErr == nil {
			logs = append(logs, map[string]string{
				"action": action, "status": status, "details": details, "createdAt": createdAt,
			})
		}
	}

	if hasArgFlag(args, "--json") {
		return printJSON(logs)
	}
	if len(logs) == 0 {
		fmt.Printf("No scan logs recorded in split database for %s.\n", slug)
		return nil
	}

	fmt.Printf("\n  %s● Scan/Sync Activity Logs for %s (%d records):%s\n", constants.ColorCyan, slug, len(logs), constants.ColorReset)
	fmt.Printf("    %-20s %-16s %-12s %s\n", "TIMESTAMP", "ACTION", "STATUS", "DETAILS")
	fmt.Printf("    %s\n", strings.Repeat("─", 78))
	for _, l := range logs {
		fmt.Printf("    %-20s %-16s %-12s %s\n", l["createdAt"], l["action"], l["status"], l["details"])
	}
	fmt.Println()
	return nil
}

func handleRepoDBErrorLogs(args []string) error {
	db, _, _, slug, err := resolveCurrentRepoSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT Action, Status, COALESCE(ErrorMessage, ''), CreatedAt FROM RepoScanLog WHERE Status = 'failure' OR ErrorMessage IS NOT NULL ORDER BY Id DESC LIMIT 20;")
	if err != nil {
		fmt.Printf("No error logs found in split database for %s.\n", slug)
		return nil
	}
	defer rows.Close()

	var errs []map[string]string
	for rows.Next() {
		var action, status, errMsg, createdAt string
		if scanErr := rows.Scan(&action, &status, &errMsg, &createdAt); scanErr == nil {
			errs = append(errs, map[string]string{
				"action": action, "status": status, "error": errMsg, "createdAt": createdAt,
			})
		}
	}

	if hasArgFlag(args, "--json") {
		return printJSON(errs)
	}
	if len(errs) == 0 {
		fmt.Printf("No error logs recorded in split database for %s.\n", slug)
		return nil
	}

	fmt.Printf("\n  %s● Recorded Error Logs for %s (%d records):%s\n", constants.ColorRed, slug, len(errs), constants.ColorReset)
	for _, e := range errs {
		fmt.Printf("    [%s] Action: %s | Status: %s\n      Error: %s\n", e["createdAt"], e["action"], e["status"], e["error"])
	}
	fmt.Println()
	return nil
}

func handleRepoDBClear(args []string) error {
	db, _, _, slug, err := resolveCurrentRepoSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !hasArgFlag(args, "-y") && !hasArgFlag(args, "--yes") {
		ok, err := promptConfirm(fmt.Sprintf("Clear search index and cache for %s? [y/N]: ", slug))
		if err != nil || !ok {
			fmt.Println("Clear operation canceled.")
			return nil
		}
	}

	if err := repodb.ClearRepoDB(context.Background(), db); err != nil {
		return err
	}
	fmt.Printf("%s✓ Repository split database cleared for %s.%s\n", constants.ColorGreen, slug, constants.ColorReset)
	return nil
}

func handleRepoDBReset(args []string) error {
	db, _, _, slug, err := resolveCurrentRepoSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !hasArgFlag(args, "-y") && !hasArgFlag(args, "--yes") {
		ok, err := promptConfirm(fmt.Sprintf("Reset repository schema for %s? [y/N]: ", slug))
		if err != nil || !ok {
			fmt.Println("Reset operation canceled.")
			return nil
		}
	}

	if err := repodb.ResetRepoDB(context.Background(), db); err != nil {
		return err
	}
	fmt.Printf("%s✓ Repository split database schema reset for %s.%s\n", constants.ColorGreen, slug, constants.ColorReset)
	return nil
}

func handleRepoDBOptimize(args []string) error {
	db, path, _, slug, err := resolveCurrentRepoSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	reclaimed, err := repodb.OptimizeRepoDB(context.Background(), db, path)
	if err != nil {
		return err
	}
	fmt.Printf("%s✓ Repository split DB optimized for %s.%s Reclaimed: %s (%s)\n",
		constants.ColorGreen, slug, constants.ColorReset, formatBytes(reclaimed), path)
	return nil
}

func resolveCurrentRepoSplitDB() (*sql.DB, string, int64, string, error) {
	cwd, _ := os.Getwd()
	slug := filepath.Base(cwd)
	repoID := int64(1)

	mainDB, err := store.OpenDefault()
	if err == nil {
		defer mainDB.Close()
		if id, findErr := mainDB.GetRepoIDByPath(cwd); findErr == nil && id > 0 {
			repoID = id
		}
	}

	rootDbDir := store.BinaryDataDir()
	dbPath := repodb.ResolveRepoDBPath(rootDbDir, cwd, repoID)
	db, openErr := repodb.OpenRepoDB(context.Background(), rootDbDir, cwd, repoID)
	if openErr != nil {
		return nil, "", 0, "", openErr
	}
	return db, dbPath, repoID, slug, nil
}
