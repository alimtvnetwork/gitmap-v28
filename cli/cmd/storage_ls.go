package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func runStorageListDatabases() error {
	dataDir := store.BinaryDataDir()
	entries := collectStorageEntries()
	printStorageInventoryHeader(dataDir)

	cfg := buildStorageTableConfig(dataDir, entries)
	termtable.PrintTable(cfg)
	printStorageInventoryFooter(entries)

	return nil
}

func printStorageInventoryFooter(entries []store.SplitDatabaseEntry) {
	totalSize := calculateTotalDBSize(entries)
	msg := fmt.Sprintf("\n  Total: %d databases (%s)\n", len(entries), cmddb.FormatBytes(totalSize))
	fmt.Print(msg)
	termpad.EnsureBottomPadding(msg)
}

func collectStorageEntries() []store.SplitDatabaseEntry {
	entries := store.CollectAllDatabaseEntries()
	seen := buildDatabaseSeenMap(entries)
	extra := collectRepoScopedEntries(seen)

	return append(entries, extra...)
}

func buildDatabaseSeenMap(entries []store.SplitDatabaseEntry) map[string]bool {
	seen := make(map[string]bool, len(entries))
	for _, e := range entries {
		seen[filepath.Clean(e.DatabasePath)] = true
	}

	return seen
}

func collectRepoScopedEntries(seen map[string]bool) []store.SplitDatabaseEntry {
	var list []store.SplitDatabaseEntry
	appendDiscoveredEntries(&list, collectUserPipelineEntries(seen))
	appendDiscoveredEntries(&list, collectRepoDbEntries(seen))

	return list
}

func appendDiscoveredEntries(dest *[]store.SplitDatabaseEntry, src []store.SplitDatabaseEntry) {
	if len(src) == 0 {
		return
	}

	*dest = append(*dest, src...)
}

func collectUserPipelineEntries(seen map[string]bool) []store.SplitDatabaseEntry {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	dir := filepath.Join(home, ".gitmap", "pipeline")

	return scanDbFilesInDir(dir, "pipeline", "User home pipeline split database", seen)
}

func collectRepoDbEntries(seen map[string]bool) []store.SplitDatabaseEntry {
	var results []store.SplitDatabaseEntry
	dirs := candidateRepoDbDirs()
	for _, dir := range dirs {
		found := scanDbFilesInDir(dir, "repodb", "Repository-scoped split database", seen)
		appendDiscoveredEntries(&results, found)
	}

	return results
}

func candidateRepoDbDirs() []string {
	dataDir := store.BinaryDataDir()
	baseDir := filepath.Dir(dataDir)
	dirs := []string{
		filepath.Join(".", "repodb"),
		filepath.Join(".", ".gitmap", "repodb"),
		filepath.Join(dataDir, "repodb"),
		filepath.Join(baseDir, "repodb"),
	}

	return appendUserHomeRepoDbDir(dirs)
}

func appendUserHomeRepoDbDir(dirs []string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return dirs
	}

	return append(dirs, filepath.Join(home, ".gitmap", "repodb"))
}

func scanDbFilesInDir(dir, dbType, desc string, seen map[string]bool) []store.SplitDatabaseEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var matched []store.SplitDatabaseEntry
	for _, e := range entries {
		appendIfValidEntry(&matched, inspectDbDirEntry(dir, e, dbType, desc, seen))
	}

	return matched
}

func appendIfValidEntry(list *[]store.SplitDatabaseEntry, item *store.SplitDatabaseEntry) {
	if item == nil {
		return
	}

	*list = append(*list, *item)
}

func inspectDbDirEntry(
	dir string,
	e os.DirEntry,
	dbType, desc string,
	seen map[string]bool,
) *store.SplitDatabaseEntry {
	isDbFile := !e.IsDir() && strings.HasSuffix(e.Name(), ".db")
	if !isDbFile {
		return nil
	}

	fullPath := filepath.Clean(filepath.Join(dir, e.Name()))
	if seen[fullPath] {
		return nil
	}
	seen[fullPath] = true

	slug := strings.TrimSuffix(e.Name(), ".db")
	entry := store.InspectSplitDBFile(dbType, slug, fullPath, desc)

	return &entry
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
