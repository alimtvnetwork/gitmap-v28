package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
)

type sqlContextExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type sqlContextQueryExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func resolveContextRunner(tx any) (sqlContextQueryExecer, error) {
	if runner, isDirect := tx.(sqlContextQueryExecer); isDirect {
		return runner, nil
	}

	if wrap, isTxWrap := tx.(*dbengine.TxWrapper); isTxWrap {
		return wrap.Tx(), nil
	}

	if wrap, isDbWrap := tx.(*dbengine.DBWrapper); isDbWrap {
		return wrap.Conn(), nil
	}

	return nil, apperror.NewValidationError(fmt.Sprintf("unsupported executor type: %T", tx))
}

func resolveContextExecer(tx any) (sqlContextExecer, error) {
	return resolveContextRunner(tx)
}

const sqlInsertSSHHost = `
	INSERT INTO ssh_hosts (id, alias, ip, username, port, encrypted_password, cluster_role, created_at)
	VALUES (:id, :alias, :ip, :username, :port, :encrypted_password, :cluster_role, :created_at)
`

const sqlSelectHostByIP = `SELECT id FROM ssh_hosts WHERE ip = ? LIMIT 1`

const sqlSelectHostByAlias = `SELECT id FROM ssh_hosts WHERE alias = ? LIMIT 1`

const sqlDeleteHostByAliasOrIP = `DELETE FROM ssh_hosts WHERE alias = ? OR ip = ?`

const sqlUpdateHostPasswordByAliasOrIP = `UPDATE ssh_hosts SET encrypted_password = ? WHERE alias = ? OR ip = ?`

const sqlDeleteAllSSHHosts = `DELETE FROM ssh_hosts`

func resolveHostID(id string, ip string) string {
	if id != "" {
		return id
	}

	return fmt.Sprintf("host-%s", ip)
}

func resolveSSHHostCreatedAt(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().UTC()
	}

	return t
}

func resolveHostPort(port int) int {
	if port > 0 {
		return port
	}
	return 22
}

func resolveClusterRole(role string) string {
	if role != "" {
		return role
	}
	return "worker"
}

func buildHostNamedArgs(host SSHHost) []any {
	return []any{
		sql.Named("id", host.ID),
		sql.Named("alias", host.Alias),
		sql.Named("ip", host.IP),
		sql.Named("username", host.Username),
		sql.Named("port", resolveHostPort(host.Port)),
		sql.Named("encrypted_password", host.EncryptedPassword),
		sql.Named("cluster_role", resolveClusterRole(host.ClusterRole)),
		sql.Named("created_at", host.CreatedAt),
	}
}

func wrapHostInsertError(err error, id string) *apperror.AppError {
	appErr := apperror.Wrap(err, "InsertSSHHost", map[string]any{"id": id})
	appErr.Code = "E_INTERNAL_ERROR"

	return appErr
}

func executeSSHHostInsert(ctx context.Context, host SSHHost, execer sqlContextExecer) error {
	_, _ = execer.ExecContext(ctx, SQLCreateSSHHostsTable)
	_, _ = execer.ExecContext(ctx, SQLCreateSSHHistoryTable)
	host.ID = resolveHostID(host.ID, host.IP)
	host.CreatedAt = resolveSSHHostCreatedAt(host.CreatedAt)
	args := buildHostNamedArgs(host)
	_, err := execer.ExecContext(ctx, sqlInsertSSHHost, args...)
	if err != nil {
		return wrapHostInsertError(err, host.ID)
	}

	return nil
}

func ensureSSHTablesTx(tx any) {
	if db, isDB := tx.(*sql.DB); isDB {
		_ = EnsureSSHTables(db)
		return
	}

	if wrap, isDbWrap := tx.(*dbengine.DBWrapper); isDbWrap && wrap != nil {
		_ = EnsureSSHTables(wrap.Conn())
	}
}

// InsertSSHHost inserts a new SSHHost into the database.
func InsertSSHHost(ctx context.Context, host SSHHost, tx any) error {
	ensureSSHTablesTx(tx)
	execer, err := resolveContextExecer(tx)
	if err != nil {
		appErr := apperror.Wrap(err, "InsertSSHHost", map[string]any{"id": host.ID})
		appErr.Code = "E_INTERNAL_ERROR"

		return appErr
	}

	return executeSSHHostInsert(ctx, host, execer)
}

// EnsureSSHHostsTable ensures that the ssh_hosts table exists.
func EnsureSSHHostsTable(db *sql.DB) error {
	return EnsureSSHTables(db)
}

// InsertSSHHostTx inserts a new SSHHost using a dbengine.TxWrapper.
func InsertSSHHostTx(ctx context.Context, host SSHHost, tx *dbengine.TxWrapper) error {
	return InsertSSHHost(ctx, host, tx)
}

func scanHostID(row *sql.Row) (string, bool, error) {
	var id string
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	return id, true, nil
}

func findHostByField(ctx context.Context, runner sqlContextQueryExecer, query string, val string) (string, bool, error) {
	if val == "" {
		return "", false, nil
	}

	return scanHostID(runner.QueryRowContext(ctx, query, val))
}

func updateHostForIP(ctx context.Context, runner sqlContextQueryExecer, id string, host SSHHost) error {
	role := resolveClusterRole(host.ClusterRole)
	createdAt := resolveSSHHostCreatedAt(host.CreatedAt)
	port := resolveHostPort(host.Port)
	query := `UPDATE ssh_hosts SET alias = ?, username = ?, port = ?, encrypted_password = CASE WHEN ? != '' THEN ? ELSE encrypted_password END, cluster_role = ?, created_at = ? WHERE id = ?`
	_, err := runner.ExecContext(ctx, query, host.Alias, host.Username, port, host.EncryptedPassword, host.EncryptedPassword, role, createdAt, id)
	if err != nil {
		return apperror.Wrap(err, "updateHostForIP", map[string]any{"id": id})
	}

	return nil
}

func updateHostForAlias(ctx context.Context, runner sqlContextQueryExecer, id string, host SSHHost) error {
	role := resolveClusterRole(host.ClusterRole)
	createdAt := resolveSSHHostCreatedAt(host.CreatedAt)
	port := resolveHostPort(host.Port)
	query := `UPDATE ssh_hosts SET ip = ?, username = ?, port = ?, encrypted_password = CASE WHEN ? != '' THEN ? ELSE encrypted_password END, cluster_role = ?, created_at = ? WHERE id = ?`
	_, err := runner.ExecContext(ctx, query, host.IP, host.Username, port, host.EncryptedPassword, host.EncryptedPassword, role, createdAt, id)
	if err != nil {
		return apperror.Wrap(err, "updateHostForAlias", map[string]any{"id": id})
	}

	return nil
}

func upsertHostByAliasOrInsert(ctx context.Context, host SSHHost, runner sqlContextQueryExecer) error {
	idByAlias, isFoundByAlias, err := findHostByField(ctx, runner, sqlSelectHostByAlias, host.Alias)
	if err != nil {
		return apperror.Wrap(err, "UpsertSSHHost_FindAlias", map[string]any{"alias": host.Alias})
	}

	if isFoundByAlias {
		return updateHostForAlias(ctx, runner, idByAlias, host)
	}

	return executeSSHHostInsert(ctx, host, runner)
}

func executeSSHHostUpsert(ctx context.Context, host SSHHost, runner sqlContextQueryExecer) error {
	_, _ = runner.ExecContext(ctx, SQLCreateSSHHostsTable)
	_, _ = runner.ExecContext(ctx, SQLCreateSSHHistoryTable)

	idByIP, isFoundByIP, err := findHostByField(ctx, runner, sqlSelectHostByIP, host.IP)
	if err != nil {
		return apperror.Wrap(err, "UpsertSSHHost_FindIP", map[string]any{"ip": host.IP})
	}

	if isFoundByIP {
		return updateHostForIP(ctx, runner, idByIP, host)
	}

	return upsertHostByAliasOrInsert(ctx, host, runner)
}

// UpsertSSHHost inserts or updates an SSHHost record idempotently.
func UpsertSSHHost(ctx context.Context, host SSHHost, tx any) error {
	ensureSSHTablesTx(tx)
	runner, err := resolveContextRunner(tx)
	if err != nil {
		appErr := apperror.Wrap(err, "UpsertSSHHost", map[string]any{"id": host.ID})
		appErr.Code = "E_INTERNAL_ERROR"

		return appErr
	}

	return executeSSHHostUpsert(ctx, host, runner)
}

func wrapDeleteError(err error, op string, target string) *apperror.AppError {
	appErr := apperror.Wrap(err, op, map[string]any{"target": target})
	appErr.Code = "E_INTERNAL_ERROR"

	return appErr
}

// DeleteHostByAliasOrIP deletes SSH hosts matching the given alias or IP and returns rows affected.
func DeleteHostByAliasOrIP(ctx context.Context, target string, db *sql.DB) (int64, error) {
	_ = EnsureSSHTables(db)
	res, err := db.ExecContext(ctx, sqlDeleteHostByAliasOrIP, target, target)
	if err != nil {
		return 0, wrapDeleteError(err, "DeleteHostByAliasOrIP", target)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, wrapDeleteError(err, "DeleteHostByAliasOrIP_Rows", target)
	}

	return affected, nil
}
