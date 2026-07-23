// Package store — owner_repo_list_cache.go: SQLite-backed TTL cache
// for `gh repo list <owner>` / `glab repo list --group <owner>` JSON
// payloads. The cmd-layer wrapper (visibilityownerlistcache.go) owns
// the TTL policy; this file only persists / fetches the raw JSON blob
// and its FetchedAt timestamp.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §parallel.
package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// LookupOwnerRepoListCache returns (namesJson, fetchedAt, true) when
// a cache row exists for (provider, owner), or ("", zero, false) when
// the row is missing or the timestamp can't be parsed.
func (db *DB) LookupOwnerRepoListCache(provider, owner string) (string, time.Time, bool) {
	const q = `SELECT NamesJson, FetchedAt FROM OwnerRepoListCache
		WHERE Provider = ? AND Owner = ?`
	var rawJSON, rawTime string
	err := db.conn.QueryRow(q, provider, owner).Scan(&rawJSON, &rawTime)
	if err != nil {
		return "", time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339Nano, rawTime)
	if err != nil {
		if t, err = time.Parse(time.RFC3339, rawTime); err != nil {
			return "", time.Time{}, false
		}
	}

	return rawJSON, t, true
}

// UpsertOwnerRepoListCache writes (or replaces) the cache row.
func (db *DB) UpsertOwnerRepoListCache(provider, owner, namesJSON string, fetchedAt time.Time) error {
	const q = `INSERT INTO OwnerRepoListCache (Provider, Owner, NamesJson, FetchedAt)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(Provider, Owner) DO UPDATE SET
			NamesJson = excluded.NamesJson,
			FetchedAt = excluded.FetchedAt`
	_, err := db.conn.Exec(q, provider, owner, namesJSON, fetchedAt.UTC().Format(time.RFC3339Nano))

	return err
}

// PurgeOwnerRepoListCache drops all rows. Exposed for tests + an
// eventual `gitmap cache purge` admin command.
func (db *DB) PurgeOwnerRepoListCache() error {
	_, err := db.conn.Exec(`DELETE FROM OwnerRepoListCache`)

	return err
}

// ensureOwnerRepoListCacheTable runs the CREATE TABLE IF NOT EXISTS
// for callers that need the cache without a full Migrate().
func (db *DB) ensureOwnerRepoListCacheTable() error {
	_, err := db.conn.Exec(constants.SQLCreateOwnerRepoListCache)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return nil
}
