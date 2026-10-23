package cmdcommit

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/heatmap"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// ErrCommitNotFound indicates that universal SHA resolution could not locate a matching commit.
var ErrCommitNotFound = errors.New("commit not found")

// ResolveUniversalCommitSha resolves any commit reference to a 40-character commit SHA.
func ResolveUniversalCommitSha(repoPath, query string) (string, error) {
	clean := strings.TrimSpace(query)
	if len(clean) == 0 {
		return "", apperror.WrapSimple(ErrCommitNotFound, "empty query")
	}

	if sha, hasMatch := tryStage1GitRevParse(repoPath, clean); hasMatch {
		return sha, nil
	}

	if sha, hasMatch := tryStage2SplitDb(repoPath, clean); hasMatch {
		return sha, nil
	}

	if sha, hasMatch := tryStage3HeatmapLog(clean); hasMatch {
		return sha, nil
	}

	return resolveStage4GitRevList(repoPath, clean)
}

func tryStage1GitRevParse(repoPath, query string) (string, bool) {
	target := query + "^{commit}"
	res := gitutil.RunGitInDir(repoPath, "rev-parse", "--verify", "--quiet", target)
	if res.IsFailure() {
		return "", false
	}

	sha := strings.TrimSpace(string(res.Data))
	if !isFullSha(sha) {
		return "", false
	}

	return strings.ToLower(sha), true
}

func tryStage2SplitDb(repoPath, query string) (string, bool) {
	dbPath := store.ResolveCommitHistoryDbPath(repoPath)
	if !isFileExisting(dbPath) {
		return "", false
	}

	conn, err := openReadOnlyDb(dbPath)
	if err != nil {
		return "", false
	}

	defer conn.Close()

	return querySplitDbSha(conn, query)
}

func querySplitDbSha(conn *sql.DB, query string) (string, bool) {
	q := `SELECT CommitSha FROM CommitHistoryRecord 
		WHERE CommitSha = ? OR ShortSha = ? OR CommitSha LIKE ? || '%' OR CommitSha LIKE '%' || ? 
		LIMIT 1`
	var sha string
	err := conn.QueryRow(q, query, query, query, query).Scan(&sha)
	if err != nil || !isFullSha(sha) {
		return "", false
	}

	return strings.ToLower(sha), true
}

func tryStage3HeatmapLog(query string) (string, bool) {
	dbs, err := heatmap.DiscoverHeatmapDatabases()
	if err != nil || len(dbs) == 0 {
		return "", false
	}

	for _, db := range dbs {
		if sha, hasMatch := scanHeatmapDbForSha(db.Path, query); hasMatch {
			return sha, true
		}
	}

	return "", false
}

func scanHeatmapDbForSha(dbPath, query string) (string, bool) {
	conn, err := openReadOnlyDb(dbPath)
	if err != nil {
		return "", false
	}

	defer conn.Close()

	var sha string
	q := `SELECT CommitSha FROM HeatmapCommitLog WHERE CommitSha LIKE '%' || ? || '%' LIMIT 1`
	if err := conn.QueryRow(q, query).Scan(&sha); err != nil || !isFullSha(sha) {
		return "", false
	}

	return strings.ToLower(sha), true
}

func resolveStage4GitRevList(repoPath, query string) (string, error) {
	res := gitutil.RunGitInDir(repoPath, "rev-list", "--all")
	if res.IsFailure() {
		return "", apperror.WrapSimple(ErrCommitNotFound, "stage 4 git rev-list failed")
	}

	lines := strings.Split(string(res.Data), "\n")
	sha, hasMatch := findShaInRevListLines(lines, query)
	if !hasMatch {
		return "", apperror.WrapSimple(ErrCommitNotFound, "commit not found: "+query)
	}

	return sha, nil
}

func findShaInRevListLines(lines []string, query string) (string, bool) {
	lowQuery := strings.ToLower(query)
	for idx, line := range lines {
		if idx >= 25000 {
			return "", false
		}

		sha := strings.ToLower(strings.TrimSpace(line))
		if isFullSha(sha) && strings.Contains(sha, lowQuery) {
			return sha, true
		}
	}

	return "", false
}

func openReadOnlyDb(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?mode=ro&cache=shared", filepath.ToSlash(dbPath))
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, "openReadOnlyDb")
	}

	return conn, nil
}

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func isFullSha(s string) bool {
	if len(s) != 40 {
		return false
	}

	for _, ch := range s {
		if !isHexRune(ch) {
			return false
		}
	}

	return true
}

func isHexRune(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
