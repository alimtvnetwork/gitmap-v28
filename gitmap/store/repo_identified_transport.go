// Package store — Plan 03 step 3 helper: read/write the
// IdentifiedTransport column (migration 008) keyed by either HTTPS
// or SSH URL. cfr/cfrp/clone-now will call these to (a) coerce the
// positional URL to the stored transport before clone, and (b)
// persist whatever transport was actually used after a successful
// clone, so subsequent recloned-after-delete invocations never
// silently downgrade SSH → HTTPS.
//
// Why URL-keyed (not RepoId): the reclone path may run BEFORE any
// scan has produced a RepoId for the destination folder. The URL
// pair (HttpsUrl, SshUrl) is the only stable identifier available
// pre-clone, and `gitutil.ConvertURLToHTTPS` / `ConvertURLToSSH`
// give us both forms from any one input.
package store

import (
	"database/sql"
	"errors"
	"strings"
)

// RepoTransportHTTPS / RepoTransportSSH are the only legal non-empty
// values for Repo.IdentifiedTransport. Empty = "unknown, never
// classified" — backfill from origin URL on the next scan/reclone.
const (
	RepoTransportHTTPS = "https"
	RepoTransportSSH   = "ssh"
)

// LookupRepoIdentifiedTransport returns the stored transport verdict
// for the repo whose HttpsUrl or SshUrl matches `url`. Returns empty
// string + nil error when no row matches OR when the column is empty
// (caller treats both as "unknown"). Any real query error is returned
// per the zero-swallow error policy.
func (db *DB) LookupRepoIdentifiedTransport(url string) (string, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return "", nil
	}
	const q = `SELECT IdentifiedTransport FROM Repo
		WHERE HttpsUrl = ? OR SshUrl = ?
		ORDER BY UpdatedAt DESC LIMIT 1`
	var transport string
	err := db.conn.QueryRow(q, url, url).Scan(&transport)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return transport, nil
}

// SetRepoIdentifiedTransport persists the transport verdict for every
// Repo row whose HttpsUrl or SshUrl matches `url`. No-ops (returns
// nil) when transport is empty — callers shouldn't overwrite a known
// value with "unknown". Returns the number of rows touched so the
// caller can decide whether to fall back to an INSERT (e.g. when the
// row hasn't been scanned yet — out of scope here).
func (db *DB) SetRepoIdentifiedTransport(url, transport string) (int64, error) {
	url = strings.TrimSpace(url)
	transport = strings.TrimSpace(strings.ToLower(transport))
	if url == "" || transport == "" {
		return 0, nil
	}
	if transport != RepoTransportHTTPS && transport != RepoTransportSSH {
		return 0, nil
	}
	const q = `UPDATE Repo SET IdentifiedTransport = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE HttpsUrl = ? OR SshUrl = ?`
	res, err := db.conn.Exec(q, transport, url, url)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// ClassifyURLTransport returns "ssh" for SSH-shorthand or ssh:// URLs,
// "https" for http(s):// URLs, and "" for anything unrecognized.
// Centralised here so the reclone flow and the lookup helper agree
// on classification — no drift with cmd/clonefixrepofoldertransport.go.
func ClassifyURLTransport(url string) string {
	lower := strings.ToLower(strings.TrimSpace(url))
	if strings.HasPrefix(lower, "git@") || strings.HasPrefix(lower, "ssh://") {
		return RepoTransportSSH
	}
	if strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://") {
		return RepoTransportHTTPS
	}
	return ""
}
