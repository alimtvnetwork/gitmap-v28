package cmdpipeline

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SQLite hash cache for `pipeline errors all` (program 265, WS7).
// Before the per-repo inspection, the repo's HEAD commit SHA is checked
// against this cache. On a hit the expensive inspection is skipped and
// the repo is marked done (cached). On a miss the repo is inspected and
// the SHA is recorded.

const sqlCreatePeAllCacheTable = `CREATE TABLE IF NOT EXISTS pe_all_commit_cache (
  repo_slug  TEXT PRIMARY KEY,
  commit_sha TEXT NOT NULL,
  checked_at TEXT NOT NULL
);`

func peAllCacheDbPath() string {
	return filepath.Join(store.BinaryDataDir(), "pe_all_commit_cache.db")
}

func openPeAllCache() (*sql.DB, error) {
	dbPath := peAllCacheDbPath()
	if mkdirErr := os.MkdirAll(filepath.Dir(dbPath), 0755); mkdirErr != nil {
		return nil, mkdirErr
	}
	conn, openErr := sql.Open("sqlite", dbPath)
	if openErr != nil {
		return nil, openErr
	}
	if _, execErr := conn.Exec(sqlCreatePeAllCacheTable); execErr != nil {
		_ = conn.Close()

		return nil, execErr
	}

	return conn, nil
}

func getCachedCommitSha(conn *sql.DB, slug string) string {
	var sha string
	row := conn.QueryRow(`SELECT commit_sha FROM pe_all_commit_cache WHERE repo_slug = ?`, slug)
	if scanErr := row.Scan(&sha); scanErr != nil {
		return ""
	}

	return sha
}

func recordCommitSha(conn *sql.DB, slug, sha string) error {
	_, err := conn.Exec(
		`INSERT INTO pe_all_commit_cache (repo_slug, commit_sha, checked_at)
		 VALUES (?, ?, ?)
		 ON CONFLICT(repo_slug) DO UPDATE SET commit_sha = excluded.commit_sha, checked_at = excluded.checked_at`,
		slug, sha, time.Now().UTC().Format(time.RFC3339),
	)
	return err
}

// getRepoHeadSha returns the HEAD commit SHA of the repo at absPath,
// or "" when the path is not a git checkout.
func getRepoHeadSha(absPath string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = absPath
	out, runErr := cmd.Output()
	if runErr != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
