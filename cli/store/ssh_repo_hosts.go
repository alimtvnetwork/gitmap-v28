package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// DeleteAllSSHHosts deletes all SSH hosts from the database and returns rows affected.
func DeleteAllSSHHosts(ctx context.Context, db *sql.DB) (int64, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	res, err := db.ExecContext(ctx, sqlDeleteAllSSHHosts)
	if err != nil {
		return 0, wrapDeleteError(err, "DeleteAllSSHHosts", "all")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, wrapDeleteError(err, "DeleteAllSSHHosts_Rows", "all")
	}

	return affected, nil
}

func safeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func wrapHostByFieldScanError(err error, op string, field string, val string) (SSHHost, error) {
	if errors.Is(err, sql.ErrNoRows) {
		appErr := apperror.Wrap(apperror.ErrNotFound, op, map[string]any{field: val})
		appErr.Code = "E_INTERNAL_ERROR"

		return SSHHost{}, appErr
	}

	appErr := apperror.Wrap(err, op, map[string]any{field: val})
	appErr.Code = "E_INTERNAL_ERROR"

	return SSHHost{}, appErr
}

const sqlSelectHostFields = `SELECT id, alias, ip, username, COALESCE(port, 22), COALESCE(encrypted_password, ''), COALESCE(cluster_role, 'worker'), created_at FROM ssh_hosts`

// GetHostByAlias retrieves an SSHHost by its alias.
func GetHostByAlias(ctx context.Context, alias string, db *sql.DB) (SSHHost, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE alias = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, alias).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.ClusterRole, &host.CreatedAt,
	)
	if err != nil {
		return wrapHostByFieldScanError(err, "GetHostByAlias", "alias", alias)
	}

	return host, nil
}

// GetHostByIP retrieves an SSHHost by its IP.
func GetHostByIP(ctx context.Context, ip string, db *sql.DB) (SSHHost, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE ip = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, ip).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.ClusterRole, &host.CreatedAt,
	)
	if err != nil {
		return wrapHostByFieldScanError(err, "GetHostByIP", "ip", ip)
	}

	return host, nil
}

// GetHostByID retrieves an SSHHost by its ID.
func GetHostByID(ctx context.Context, id string, db *sql.DB) (SSHHost, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE id = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, id).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.ClusterRole, &host.CreatedAt,
	)
	if err != nil {
		return wrapHostByFieldScanError(err, "GetHostByID", "id", id)
	}

	return host, nil
}

// DeleteHostByIP deletes an SSHHost by its IP.
func DeleteHostByIP(ctx context.Context, ip string, db *sql.DB) error {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	query := `DELETE FROM ssh_hosts WHERE ip = ?`
	_, err := db.ExecContext(ctx, query, ip)
	if err != nil {
		appErr := apperror.Wrap(err, "DeleteHostByIP", map[string]any{"ip": ip})
		appErr.Code = "E_INTERNAL_ERROR"

		return appErr
	}

	return nil
}

// UpdateHostPassword updates the encrypted password for a host identified by alias or IP.
func UpdateHostPassword(ctx context.Context, target string, encryptedPass string, db *sql.DB) (int64, error) {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	res, err := db.ExecContext(ctx, sqlUpdateHostPasswordByAliasOrIP, encryptedPass, target, target)
	if err != nil {
		return 0, apperror.WrapSimple(err, "UpdateHostPassword.Exec")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

func scanHostRows(rows *sql.Rows) ([]SSHHost, error) {
	var hosts []SSHHost
	for rows.Next() {
		var host SSHHost
		if err := rows.Scan(&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.ClusterRole, &host.CreatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "ListHosts_Scan")
		}

		hosts = append(hosts, host)
	}

	return hosts, rows.Err()
}

func ensureHostsSlice(hosts []SSHHost) []SSHHost {
	if hosts == nil {
		return []SSHHost{}
	}

	return hosts
}

// ListHostsWrap retrieves all SSH hosts encapsulated in SSHHostSliceWrap.
func ListHostsWrap(ctx context.Context, db *sql.DB) SSHHostSliceWrap {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return result.FailSlice[SSHHost](apperror.WrapSimple(err, "ListHosts"))
	}
	defer rows.Close()
	hosts, err := scanHostRows(rows)
	if err != nil {
		return result.FailSlice[SSHHost](apperror.WrapSimple(err, "ListHosts_Rows"))
	}
	return result.OkSlice(ensureHostsSlice(hosts))
}

// ListHosts retrieves all SSH hosts from the database.
func ListHosts(ctx context.Context, db *sql.DB) ([]SSHHost, error) {
	res := ListHostsWrap(ctx, db)
	if res.HasError() {
		return nil, res.Err
	}
	return res.Items(), nil
}

// ListHostsByRoleWrap retrieves SSH hosts matching the given cluster role encapsulated in SSHHostSliceWrap.
func ListHostsByRoleWrap(ctx context.Context, role string, db *sql.DB) SSHHostSliceWrap {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE cluster_role = ? ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query, role)
	if err != nil {
		return result.FailSlice[SSHHost](apperror.WrapSimple(err, "ListHostsByRole"))
	}
	defer rows.Close()
	hosts, err := scanHostRows(rows)
	if err != nil {
		return result.FailSlice[SSHHost](apperror.WrapSimple(err, "ListHostsByRole_Rows"))
	}
	return result.OkSlice(ensureHostsSlice(hosts))
}

// ListHostsByRole retrieves SSH hosts matching the given cluster role.
func ListHostsByRole(ctx context.Context, role string, db *sql.DB) ([]SSHHost, error) {
	res := ListHostsByRoleWrap(ctx, role, db)
	if res.HasError() {
		return nil, res.Err
	}
	return res.Items(), nil
}
