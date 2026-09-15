package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runStorageDriveReport(args []string) error {
	path, _ := os.Getwd()
	if len(args) > 0 && !stringsHasPrefix(args[0], "-") {
		path = args[0]
	}

	metrics, err := getDiskSpaceMetrics(path)
	if err != nil {
		return err
	}

	entries := store.CollectAllDatabaseEntries()
	totalDBSize := calculateTotalDBSize(entries)
	printDriveSummary(metrics, totalDBSize, len(entries))

	return nil
}

func calculateTotalDBSize(entries []store.SplitDatabaseEntry) int64 {
	var total int64
	for _, e := range entries {
		total += e.SizeBytes
	}

	return total
}

func formatUintBytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	}

	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1024.0*1024.0))
	}

	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024.0*1024.0*1024.0))
}

func printDriveSummary(info *DiskSpaceInfo, dbBytes int64, dbCount int) {
	fmt.Printf("\n%s Storage Drive Capacity & Diagnostics:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Drive:         %s (format: %s)\n", info.DrivePath, info.Filesystem)
	fmt.Printf("  Total Space:   %s\n", formatUintBytes(info.TotalBytes))
	fmt.Printf("  Used Space:    %s (%.1f%%)\n", formatUintBytes(info.UsedBytes), info.UsedPercent)
	fmt.Printf("  Free Space:    %s\n", formatUintBytes(info.FreeBytes))
	fmt.Println()
	fmt.Printf("  Gitmap DBs:    %d database(s) occupying %s\n", dbCount, cmddb.FormatBytes(dbBytes))
	fmt.Printf("  (run %sgitmap storage ls%s to list all SQLite database files)\n\n", constants.ColorGreen, constants.ColorReset)
}

func runStorageListDatabases() error {
	entries := store.CollectAllDatabaseEntries()
	printDatabaseHeader()

	for _, e := range entries {
		printDatabaseRow(e)
	}

	totalSize := calculateTotalDBSize(entries)
	fmt.Printf("\n  Total: %d databases (%s)\n\n", len(entries), cmddb.FormatBytes(totalSize))

	return nil
}

func printDatabaseHeader() {
	fmt.Printf("\n%s Gitmap SQLite Database Inventory:%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %-12s %-8s %-10s %-8s %-8s %s\n", "NAME", "TYPE", "SIZE", "TABLES", "RECORDS", "PATH")
	fmt.Println("  --------------------------------------------------------------------------------")
}

func printDatabaseRow(e store.SplitDatabaseEntry) {
	sizeStr := cmddb.FormatBytes(e.SizeBytes)
	fmt.Printf("  %-12s %-8s %-10s %-8d %-8d %s\n",
		e.DatabaseKey, e.DatabaseType, sizeStr, e.TableCount, e.RecordCount, e.DatabasePath)
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
