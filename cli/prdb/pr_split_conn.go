package prdb

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// OpenPrSplitDb opens the split SQLite database for a repository's PR subsystem.
func OpenPrSplitDb(repoSlug, repoRoot string) PrSplitDbWrap {
	dbPath := ResolvePrDbPath(repoSlug, repoRoot)
	slug := SanitizeRepoSlug(repoSlug)

	return openPrSplitDbWithSlug(dbPath, slug)
}

// OpenPrSplitDbAt opens or initializes the PR split DB at an explicit filesystem path.
func OpenPrSplitDbAt(dbPath string) PrSplitDbWrap {
	return openPrSplitDbWithSlug(dbPath, "pr-custom")
}

func openPrSplitDbWithSlug(dbPath, slug string) PrSplitDbWrap {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return result.Fail[*PrSplitDb](apperror.WrapSimple(err, "prdb.mkdir"))
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return result.Fail[*PrSplitDb](apperror.WrapSimple(err, "prdb.open"))
	}

	return initPrSplitConn(conn, slug, dbPath)
}

func initPrSplitConn(conn *sql.DB, repoSlug, dbPath string) PrSplitDbWrap {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return result.Fail[*PrSplitDb](apperror.WrapSimple(err, "prdb.configure"))
	}

	if err := InitPrSchema(conn); err != nil {
		_ = conn.Close()

		return result.Fail[*PrSplitDb](err)
	}

	return result.Ok(&PrSplitDb{conn: conn, RepoSlug: repoSlug, Path: dbPath})
}
