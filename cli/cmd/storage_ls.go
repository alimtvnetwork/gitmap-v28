package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func runStorageListDatabases() error {
	dataDir := store.BinaryDataDir()
	entries := store.CollectAllDatabaseEntries()
	printStorageInventoryHeader(dataDir)

	cfg := buildStorageTableConfig(dataDir, entries)
	termtable.PrintTable(cfg)

	totalSize := calculateTotalDBSize(entries)
	fmt.Printf("\n  Total: %d databases (%s)\n\n", len(entries), cmddb.FormatBytes(totalSize))

	return nil
}

func printStorageInventoryHeader(dataDir string) {
	fmt.Printf("\n%s Gitmap SQLite Database Inventory:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %sRoot Folder:%s %s\n\n", constants.ColorDim, constants.ColorReset, dataDir)
}

func buildStorageTableConfig(dataDir string, entries []store.SplitDatabaseEntry) termtable.TableConfig {
	var rows []termtable.Row
	for _, e := range entries {
		rows = append(rows, buildStorageTableRow(dataDir, e))
	}

	return termtable.TableConfig{
		Columns: storageTableColumns(),
		Rows:    rows,
	}
}

func storageTableColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "NAME", Align: termtable.AlignLeft},
		{Title: "TYPE", Align: termtable.AlignLeft},
		{Title: "SIZE", Align: termtable.AlignRight},
		{Title: "TABLES", Align: termtable.AlignRight},
		{Title: "RECORDS", Align: termtable.AlignRight},
		{Title: "PATH", Align: termtable.AlignLeft},
	}
}

func buildStorageTableRow(dataDir string, e store.SplitDatabaseEntry) termtable.Row {
	relPath := resolveRelativeDBPath(dataDir, e.DatabasePath)

	return termtable.Row{
		Cells: []string{
			e.DatabaseKey,
			e.DatabaseType,
			cmddb.FormatBytes(e.SizeBytes),
			fmt.Sprintf("%d", e.TableCount),
			fmt.Sprintf("%d", e.RecordCount),
			relPath,
		},
	}
}

func resolveRelativeDBPath(dataDir, fullPath string) string {
	rel, err := filepath.Rel(dataDir, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fullPath
	}

	return rel
}
