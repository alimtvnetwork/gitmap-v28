package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

type dbListingOptions struct {
	isAll     bool
	isRefresh bool
	isJson    bool
	isHelp    bool
}

// RunDbListing implements the gitmap db-listing (alias: dbl) command.
func RunDbListing(args []string) *appfault.AppError {
	opts := parseListingArgs(args)
	if opts.isHelp {
		printListingHelp(os.Stdout)
		return nil
	}

	if opts.isRefresh {
		refreshRegistrySizes()
	}

	summaries, err := collectDatabaseSummaries(opts.isAll)
	if err != nil {
		return appfault.NewValidationError(fmt.Sprintf("collect database listing: %v", err))
	}

	if opts.isJson {
		return renderListingJSON(summaries)
	}

	renderListingTable(os.Stdout, summaries)
	return nil
}

func parseListingArgs(args []string) dbListingOptions {
	var opts dbListingOptions
	for _, arg := range args {
		applyListingFlag(arg, &opts)
	}
	return opts
}

func applyListingFlag(arg string, opts *dbListingOptions) {
	switch arg {
	case "--all", "-a":
		opts.isAll = true
	case "--refresh", "-r":
		opts.isRefresh = true
	case "--json", "-j":
		opts.isJson = true
	case "--help", "-h":
		opts.isHelp = true
	}
}

func printListingHelp(out io.Writer) {
	help := `GitMap Database Listing Subsystem

Usage:
  gitmap db-listing [flags]
  gitmap dbl [flags]
  gitmap db listing [flags]

Flags:
  -a, --all        Include discovered split and profile databases
  -r, --refresh    Refresh file size stats and table counts
  -j, --json       Output JSON envelope format
  -h, --help       Show this help text

Examples:
  gitmap db-listing
  gitmap db-listing --all
  gitmap db-listing --refresh
  gitmap db-listing --json
`
	_, _ = fmt.Fprint(out, help)
}

func refreshRegistrySizes() {
	aliases, err := ListAliases()
	if err != nil {
		return
	}

	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return
	}
	conn := openRes.Value
	defer conn.Close()

	for _, a := range aliases {
		refreshOneAliasStat(conn, a)
	}
}

func refreshOneAliasStat(conn *sql.DB, a DatabaseAlias) {
	size := getFileSize(a.DatabasePath)
	if _, err := conn.Exec("UPDATE DatabaseAlias SET FileSizeBytes = ? WHERE DatabaseAliasId = ?", size, a.DatabaseAliasId); err != nil {
		return
	}
}

func collectDatabaseSummaries(includeAll bool) ([]DatabaseAliasSummary, error) {
	aliases, err := ListAliases()
	if err != nil {
		return nil, err
	}

	summaries := buildSummariesFromAliases(aliases)
	if !includeAll {
		return summaries, nil
	}

	return appendDiscoveredDatabases(summaries), nil
}

func buildSummariesFromAliases(aliases []DatabaseAlias) []DatabaseAliasSummary {
	var list []DatabaseAliasSummary
	for _, a := range aliases {
		list = append(list, convertAliasToSummary(a))
	}
	return list
}

func convertAliasToSummary(a DatabaseAlias) DatabaseAliasSummary {
	size := a.FileSizeBytes
	if size == 0 {
		size = getFileSize(a.DatabasePath)
	}
	tableCount := countTablesInDb(a.DatabasePath)

	return DatabaseAliasSummary{
		Alias:          a.AliasName,
		Path:           a.DatabasePath,
		Type:           a.DatabaseType,
		Size:           formatBytes(size),
		FileSizeBytes:  size,
		TableCount:     tableCount,
		SearchCount:    a.SearchHistoryCount,
		IsActive:       a.IsActive,
		Description:    a.Description,
		LastAccessedAt: a.LastAccessedAt,
	}
}

func countTablesInDb(dbPath string) int {
	if _, err := os.Stat(dbPath); err != nil {
		return 0
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return 0
	}
	conn := openRes.Value
	defer conn.Close()

	var count int
	query := `SELECT count(*) FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%'`
	if err := conn.QueryRow(query).Scan(&count); err != nil {
		return 0
	}
	return count
}

func appendDiscoveredDatabases(existing []DatabaseAliasSummary) []DatabaseAliasSummary {
	seenPaths := buildSeenPathMap(existing)
	out := append([]DatabaseAliasSummary(nil), existing...)

	for _, split := range collectSplitDBs() {
		appendSplitDbIfUnseen(&out, seenPaths, split)
	}
	for _, prof := range collectProfileDBs() {
		appendProfileDbIfUnseen(&out, seenPaths, prof)
	}
	return out
}

func buildSeenPathMap(items []DatabaseAliasSummary) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, item := range items {
		m[filepath.Clean(item.Path)] = true
	}
	return m
}

func appendSplitDbIfUnseen(list *[]DatabaseAliasSummary, seen map[string]bool, item DBFileInfo) {
	appendDiscoveredDb(list, seen, item, "split", "split-db")
}

func appendProfileDbIfUnseen(list *[]DatabaseAliasSummary, seen map[string]bool, item DBFileInfo) {
	appendDiscoveredDb(list, seen, item, "profile", "profile")
}

func appendDiscoveredDb(list *[]DatabaseAliasSummary, seen map[string]bool, item DBFileInfo, prefix, dbType string) {
	cleanPath := filepath.Clean(item.Path)
	if seen[cleanPath] {
		return
	}
	seen[cleanPath] = true
	alias := prefix + ":" + strings.TrimSuffix(item.Name, ".db")
	*list = append(*list, DatabaseAliasSummary{
		Alias:         alias,
		Path:          item.Path,
		Type:          dbType,
		Size:          formatBytes(item.Size),
		FileSizeBytes: item.Size,
		TableCount:    countTablesInDb(item.Path),
		SearchCount:   0,
		IsActive:      true,
		Description:   item.Purpose,
	})
}

func renderListingTable(out io.Writer, summaries []DatabaseAliasSummary) {
	if len(summaries) == 0 {
		fmt.Fprintln(out, "No databases registered or discovered.")
		return
	}

	cols := []termout.Column{
		{Title: "ALIAS", MinWidth: 14, Align: termout.AlignLeft},
		{Title: "PATH", MinWidth: 32, Align: termout.AlignLeft},
		{Title: "SIZE", MinWidth: 20, Align: termout.AlignRight},
		{Title: "TABLES", MinWidth: 8, Align: termout.AlignRight},
		{Title: "SEARCHES", MinWidth: 10, Align: termout.AlignRight},
	}

	rows := make([]termout.Row, 0, len(summaries))
	for _, s := range summaries {
		rows = append(rows, buildSummaryRow(s))
	}

	tableStr := termout.RenderTable(termout.TableConfig{
		Columns:    cols,
		Rows:       rows,
		HasBorders: true,
	})
	fmt.Fprint(out, tableStr)
}

func buildSummaryRow(s DatabaseAliasSummary) termout.Row {
	sizeStr := fmt.Sprintf("%s (%d B)", s.Size, s.FileSizeBytes)
	return termout.Row{
		Cells: []string{
			s.Alias,
			s.Path,
			sizeStr,
			strconv.Itoa(s.TableCount),
			strconv.FormatInt(s.SearchCount, 10),
		},
	}
}

func renderListingJSON(summaries []DatabaseAliasSummary) *appfault.AppError {
	var totalBytes int64
	for _, s := range summaries {
		totalBytes += s.FileSizeBytes
	}

	payload := DatabaseListingPayload{
		Databases:  summaries,
		TotalCount: len(summaries),
		TotalBytes: totalBytes,
		TotalSize:  formatBytes(totalBytes),
	}

	err := jsonx.WriteOK(os.Stdout, "db-listing", payload)
	if err != nil {
		return appfault.NewValidationError(fmt.Sprintf("write json envelope: %v", err))
	}
	return nil
}
