// Package store — owner_repo_name_index.go: per-name index over the
// owner-repo cache, with BaseName/VersionNumber pre-parsed so lookups
// like "highest -vN for base macro-ahk" run as a single SQL query.
//
// Populated by visibilityownerlistcache.go alongside the JSON blob.
package store

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/visibility"
)

// UpsertOwnerRepoNameIndex replaces every row for (provider, owner)
// with one row per name, deriving BaseName/VersionNumber from the
// `-vN` convention via visibility.ParseRepoNameMeta.
func (db *DB) UpsertOwnerRepoNameIndex(provider, owner string, names []string, fetchedAt time.Time) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM OwnerRepoNameIndex WHERE Provider=? AND Owner=?`, provider, owner); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO OwnerRepoNameIndex
		(Provider, Owner, RepoName, BaseName, VersionNumber, FetchedAt)
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ts := fetchedAt.UTC().Format(time.RFC3339Nano)
	for _, n := range names {
		base, ver, ok := visibility.ParseRepoNameMeta(n)
		if !ok {
			base = n
			ver = -1
		}
		if _, err := stmt.Exec(provider, owner, n, base, ver, ts); err != nil {
			return err
		}
	}

	return tx.Commit()
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
	_, err := db.conn.Exec(constants.SQLCreateOwnerRepoNameIndex)

	return err
}
