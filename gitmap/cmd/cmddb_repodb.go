package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type repoDBRow struct {
	ID        int64
	Slug      string
	DBFile    string
	Size      int64
	FileCount int
	CacheRows int
	Status    string
}

func runDBRepoDB(args []string) error {
	splitDBs := collectSplitDBs()
	repoMap := loadTrackedRepoMap()
	rows := buildRepoDBRows(splitDBs, repoMap)

	printRepoDBHeader(len(rows))
	if len(rows) == 0 {
		printRepoDBEmpty()
		return nil
	}
	printRepoDBTable(rows)
	return nil
}

func loadTrackedRepoMap() map[int64]string {
	repoMap := make(map[int64]string)
	mainDB, err := store.OpenDefault()
	if err != nil {
		return repoMap
	}
	defer mainDB.Close()
	repos, err := mainDB.ListRepos()
	if err != nil {
		return repoMap
	}
	for _, r := range repos {
		repoMap[r.ID] = r.Slug
	}
	return repoMap
}

func buildRepoDBRows(splitDBs []DBFileInfo, repoMap map[int64]string) []repoDBRow {
	var rows []repoDBRow
	for _, s := range splitDBs {
		fc, cc := querySplitDBCounts(s.Path)
		status := "Orphaned"
		if s.RepoID > 0 && repoMap[s.RepoID] != "" {
			status = "Tracked"
		}
		rows = append(rows, repoDBRow{
			ID:        s.RepoID,
			Slug:      s.RepoSlug,
			DBFile:    s.Name,
			Size:      s.Size,
			FileCount: fc,
			CacheRows: cc,
			Status:    status,
		})
	}
	return rows
}

func querySplitDBCounts(path string) (int, int) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return 0, 0
	}
	defer conn.Close()
	fc := querySingleCount(conn, "SELECT COUNT(*) FROM RepoFile")
	cc := querySingleCount(conn, "SELECT COUNT(*) FROM SearchCache")
	return fc, cc
}

func querySingleCount(conn *sql.DB, query string) int {
	var count int
	row := conn.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		return 0
	}
	return count
}

func printRepoDBHeader(count int) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Gitmap Split Repository Databases (repo_search) ──" + constants.ColorReset)
	fmt.Println()
	fmt.Printf("  Discovered %s%d%s split repository database(s):\n\n", constants.ColorWhite, count, constants.ColorReset)
}

func printRepoDBEmpty() {
	fmt.Println("  " + constants.ColorDim + "No split repository databases found in repo_search/." + constants.ColorReset)
	fmt.Println("  " + constants.ColorDim + "Split DBs are generated when indexing repository file contents or searching." + constants.ColorReset)
	fmt.Println()
}

func printRepoDBTable(rows []repoDBRow) {
	fmt.Printf("  %-8s %-26s %-24s %-10s %-8s %-8s %-10s\n",
		"REPO ID", "REPO SLUG", "DB FILE", "SIZE", "FILES", "CACHES", "STATUS")
	fmt.Println("  " + strings.Repeat("─", 98))
	for _, r := range rows {
		idStr := "-"
		if r.ID > 0 {
			idStr = fmt.Sprintf("%d", r.ID)
		}
		statusColor := constants.ColorGreen
		if r.Status == "Orphaned" {
			statusColor = constants.ColorYellow
		}
		fmt.Printf("  %-8s %-26s %-24s %-10s %-8d %-8d %s%-10s%s\n",
			idStr, truncateStr(r.Slug, 25), truncateStr(r.DBFile, 23),
			formatBytes(r.Size), r.FileCount, r.CacheRows,
			statusColor, r.Status, constants.ColorReset)
	}
	fmt.Println()
}
