package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SSHKnownHost represents a known host entry stored in the SQLite database.
type SSHKnownHost struct {
	ID          string    `json:"id"`
	Host        string    `json:"host"`
	KeyType     string    `json:"key_type"`
	PublicKey   string    `json:"public_key"`
	Fingerprint string    `json:"fingerprint"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const sqlUpsertSSHKnownHost = `
	INSERT INTO ssh_known_hosts (id, host, key_type, public_key, fingerprint, comment, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		host = excluded.host,
		key_type = excluded.key_type,
		public_key = excluded.public_key,
		fingerprint = excluded.fingerprint,
		comment = excluded.comment,
		updated_at = excluded.updated_at
`

const sqlSelectListKnownHosts = `
	SELECT id, host, key_type, public_key, fingerprint, COALESCE(comment, ''), created_at, updated_at
	FROM ssh_known_hosts ORDER BY host ASC
`

const sqlDeleteKnownHost = `
	DELETE FROM ssh_known_hosts
	WHERE id = ? OR host = ? OR host LIKE ? OR fingerprint = ?
`

const sqlGetKnownHost = `
	SELECT id, host, key_type, public_key, fingerprint, COALESCE(comment, ''), created_at, updated_at
	FROM ssh_known_hosts
	WHERE id = ? OR host = ? OR host LIKE ?
	LIMIT 1
`

// UpsertSSHKnownHost inserts or updates a known host record in SQLite.
func UpsertSSHKnownHost(ctx context.Context, kh SSHKnownHost, db *sql.DB) error {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	now := time.Now().UTC()
	if kh.CreatedAt.IsZero() {
		kh.CreatedAt = now
	}
	kh.UpdatedAt = now
	_, err := db.ExecContext(ctx, sqlUpsertSSHKnownHost,
		kh.ID, kh.Host, kh.KeyType, kh.PublicKey, kh.Fingerprint, kh.Comment, kh.CreatedAt, kh.UpdatedAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "UpsertSSHKnownHost.Exec")
	}
	return nil
}

// ListSSHKnownHosts retrieves all known hosts from SQLite.
func ListSSHKnownHosts(ctx context.Context, db *sql.DB) ([]SSHKnownHost, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	rows, err := db.QueryContext(ctx, sqlSelectListKnownHosts)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListSSHKnownHosts.Query")
	}
	defer rows.Close()
	return scanKnownHostRows(rows)
}

func scanKnownHostRows(rows *sql.Rows) ([]SSHKnownHost, error) {
	var result []SSHKnownHost
	for rows.Next() {
		var h SSHKnownHost
		if err := rows.Scan(&h.ID, &h.Host, &h.KeyType, &h.PublicKey, &h.Fingerprint, &h.Comment, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "scanKnownHostRows")
		}
		result = append(result, h)
	}
	return result, rows.Err()
}

// DeleteSSHKnownHostByTarget removes a known host by id, host, or fingerprint.
func DeleteSSHKnownHostByTarget(ctx context.Context, target string, db *sql.DB) (int64, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	pattern := "%" + target + "%"
	res, err := db.ExecContext(ctx, sqlDeleteKnownHost, target, target, pattern, target)
	if err != nil {
		return 0, apperror.WrapSimple(err, "DeleteSSHKnownHostByTarget.Exec")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

// GetSSHKnownHost finds a known host by target or host prefix.
func GetSSHKnownHost(ctx context.Context, target string, db *sql.DB) (SSHKnownHost, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	pattern := "%" + target + "%"
	var h SSHKnownHost
	err := db.QueryRowContext(ctx, sqlGetKnownHost, target, target, pattern).Scan(
		&h.ID, &h.Host, &h.KeyType, &h.PublicKey, &h.Fingerprint, &h.Comment, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		return SSHKnownHost{}, apperror.WrapSimple(err, "GetSSHKnownHost.QueryRow")
	}
	return h, nil
}
