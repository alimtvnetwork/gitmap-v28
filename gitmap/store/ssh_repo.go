package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// InsertSSHHost inserts a new SSHHost into the database.
//nolint:revive
func InsertSSHHost(ctx context.Context, host SSHHost, tx *sql.Tx) error {
	query := `
		INSERT INTO ssh_hosts (id, alias, ip, username, created_at)
		VALUES (:id, :alias, :ip, :username, :created_at)
	`
	_, err := tx.ExecContext(ctx, query,
		sql.Named("id", host.ID),
		sql.Named("alias", host.Alias),
		sql.Named("ip", host.IP),
		sql.Named("username", host.Username),
		sql.Named("created_at", host.CreatedAt),
	)
	if err != nil {
		appErr := apperror.Wrap(err, "InsertSSHHost", map[string]any{"id": host.ID})
		appErr.Code = "E_INTERNAL_ERROR"
		return appErr
	}
	return nil
}

// GetHostByAlias retrieves an SSHHost by its alias.
//nolint:revive
func GetHostByAlias(ctx context.Context, alias string, db *sql.DB) (SSHHost, error) {
	query := `SELECT id, alias, ip, username, created_at FROM ssh_hosts WHERE alias = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, alias).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			appErr := apperror.Wrap(apperror.ErrNotFound, "GetHostByAlias", map[string]any{"alias": alias})
			appErr.Code = "E_INTERNAL_ERROR"
			return SSHHost{}, appErr
		}
		appErr := apperror.Wrap(err, "GetHostByAlias", map[string]any{"alias": alias})
		appErr.Code = "E_INTERNAL_ERROR"
		return SSHHost{}, appErr
	}
	return host, nil
}

// DeleteHostByIP deletes an SSHHost by its IP.
//nolint:revive
func DeleteHostByIP(ctx context.Context, ip string, db *sql.DB) error {
	query := `DELETE FROM ssh_hosts WHERE ip = ?`
	res, err := db.ExecContext(ctx, query, ip)
	if err != nil {
		appErr := apperror.Wrap(err, "DeleteHostByIP", map[string]any{"ip": ip})
		appErr.Code = "E_INTERNAL_ERROR"
		return appErr
	}

	// Task 5: Return nil if 0 rows affected (idempotency)
	_, _ = res.RowsAffected() // Idempotent, so we just ignore if 0.

	return nil
}

// ListHosts retrieves all SSH hosts from the database.
//nolint:revive
func ListHosts(ctx context.Context, db *sql.DB) ([]SSHHost, error) {
	query := `SELECT id, alias, ip, username, created_at FROM ssh_hosts ORDER BY created_at DESC`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		appErr := apperror.WrapSimple(err, "ListHosts")
		appErr.Code = "E_INTERNAL_ERROR"
		return nil, appErr
	}
	defer rows.Close()

	var hosts []SSHHost
	for rows.Next() {
		var host SSHHost
		if err := rows.Scan(&host.ID, &host.Alias, &host.IP, &host.Username, &host.CreatedAt); err != nil {
			appErr := apperror.WrapSimple(err, "ListHosts_Scan")
			appErr.Code = "E_INTERNAL_ERROR"
			return nil, appErr
		}
		hosts = append(hosts, host)
	}

	if err := rows.Err(); err != nil {
		appErr := apperror.WrapSimple(err, "ListHosts_Rows")
		appErr.Code = "E_INTERNAL_ERROR"
		return nil, appErr
	}

	if hosts == nil {
		hosts = []SSHHost{}
	}

	return hosts, nil
}
