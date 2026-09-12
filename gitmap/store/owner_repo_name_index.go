// Package store — owner_repo_name_index.go: per-name index over the
// owner-repo cache, with BaseName/VersionNumber pre-parsed so lookups
// like "highest -vN for base macro-ahk" run as a single SQL query.
//
// Populated by visibilityownerlistcache.go alongside the JSON blob.
package store

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/visibility"
)

// UpsertOwnerRepoNameIndex replaces every row for (provider, owner)
// with one row per name, deriving BaseName/VersionNumber from the
// `-vN` convention via visibility.ParseRepoNameMeta.
func (db *DB) UpsertOwnerRepoNameIndex(
	provider,
	owner string,
	names []string,
	fetchedAt time.Time,
) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "begin tx owner repo index")
	}
	defer tx.Rollback()

	if err := populateOwnerRepoIndexTx(tx, provider, owner, names, fetchedAt); err != nil {
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return apperror.WrapSimple(commitErr, "commit owner repo name index")
	}

	return nil
}

func populateOwnerRepoIndexTx(
	tx *sql.Tx,
	provider,
	owner string,
	names []string,
	fetchedAt time.Time,
) error {
	if _, err := tx.Exec(`DELETE FROM OwnerRepoNameIndex WHERE Provider=? AND Owner=?`, provider, owner); err != nil {
		return apperror.WrapSimple(err, "delete existing owner repo index")
	}

	stmt, err := tx.Prepare(`INSERT INTO OwnerRepoNameIndex
		(Provider, Owner, RepoName, BaseName, VersionNumber, FetchedAt)
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return apperror.WrapSimple(err, "prepare insert owner repo index")
	}
	defer stmt.Close()

	return insertOwnerRepoNames(stmt, provider, owner, names, fetchedAt)
}

func insertOwnerRepoNames(
	stmt *sql.Stmt,
	provider,
	owner string,
	names []string,
	fetchedAt time.Time,
) error {
	ts := fetchedAt.UTC().Format(time.RFC3339Nano)
	for _, n := range names {
		base, ver, ok := visibility.ParseRepoNameMeta(n)
		if !ok {
			base = n
			ver = -1
		}
		if _, err := stmt.Exec(provider, owner, n, base, ver, ts); err != nil {
			return apperror.WrapSimple(err, "insert owner repo index row")
		}
	}

	return nil
}

// LookupHighestVersion returns (RepoName, VersionNumber, true) for
// the highest -vN row whose BaseName equals `base`. Returns
// ("", 0, false) when no versioned sibling exists for that base.
func (db *DB) LookupHighestVersion(provider, owner, base string) (string, int, bool) {
	const q = `SELECT RepoName, VersionNumber FROM OwnerRepoNameIndex
		WHERE Provider=? AND Owner=? AND BaseName=? AND VersionNumber >= 0
		ORDER BY VersionNumber DESC LIMIT 1`
	var name string
	var ver int
	if err := db.conn.QueryRow(q, provider, owner, base).Scan(&name, &ver); err != nil {
		return "", 0, false
	}

	return name, ver, true
}

// EnsureOwnerRepoNameIndex creates the table for callers that need
// it without a full Migrate(). Idempotent.
func (db *DB) EnsureOwnerRepoNameIndex() error {
	_, err := ExecWrapper(db.conn, constants.SQLCreateOwnerRepoNameIndex).Destruct()

	return err
}
