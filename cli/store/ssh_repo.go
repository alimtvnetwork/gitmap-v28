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

	if wrap, isDbWrap := tx.(*dbengine.DbWrapper); isDbWrap {
		return wrap.Conn(), nil
	}

	return nil, fmt.Errorf("unsupported executor type: %T", tx)
}

func resolveContextExecer(tx any) (sqlContextExecer, error) {
	return resolveContextRunner(tx)
}

const sqlInsertSSHHost = `
	INSERT INTO ssh_hosts (id, alias, ip, username, port, encrypted_password, created_at)
	VALUES (:id, :alias, :ip, :username, :port, :encrypted_password, :created_at)
`

const sqlSelectHostByIP = `SELECT id FROM ssh_hosts WHERE ip = ? LIMIT 1`

const sqlSelectHostByAlias = `SELECT id FROM ssh_hosts WHERE alias = ? LIMIT 1`

const sqlDeleteHostByAliasOrIP = `DELETE FROM ssh_hosts WHERE alias = ? OR ip = ?`

func resolveHostID(id string, ip string) string {
	if id != "" {
		return id
	}

	return fmt.Sprintf("host-%s", ip)
}

func resolveCreatedAt(t time.Time) time.Time {
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

func buildHostNamedArgs(host SSHHost) []any {
	return []any{
		sql.Named("id", host.ID),
		sql.Named("alias", host.Alias),
		sql.Named("ip", host.IP),
		sql.Named("username", host.Username),
		sql.Named("port", resolveHostPort(host.Port)),
		sql.Named("encrypted_password", host.EncryptedPassword),
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
	host.CreatedAt = resolveCreatedAt(host.CreatedAt)
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

	if wrap, isDbWrap := tx.(*dbengine.DbWrapper); isDbWrap && wrap != nil {
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
	createdAt := resolveCreatedAt(host.CreatedAt)
	port := resolveHostPort(host.Port)
	query := `UPDATE ssh_hosts SET alias = ?, username = ?, port = ?, encrypted_password = ?, created_at = ? WHERE id = ?`
	_, err := runner.ExecContext(ctx, query, host.Alias, host.Username, port, host.EncryptedPassword, createdAt, id)
	if err != nil {
		return apperror.Wrap(err, "updateHostForIP", map[string]any{"id": id})
	}

	return nil
}

func updateHostForAlias(ctx context.Context, runner sqlContextQueryExecer, id string, host SSHHost) error {
	createdAt := resolveCreatedAt(host.CreatedAt)
	port := resolveHostPort(host.Port)
	query := `UPDATE ssh_hosts SET ip = ?, username = ?, port = ?, encrypted_password = ?, created_at = ? WHERE id = ?`
	_, err := runner.ExecContext(ctx, query, host.IP, host.Username, port, host.EncryptedPassword, createdAt, id)
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

func resolveHistoryID(id string, hostIP string) string {
	if id != "" {
		return id
	}

	return fmt.Sprintf("hist-%s", hostIP)
}

func recordSSHHistory(ctx context.Context, execer sqlContextExecer, history SSHHistory) error {
	history.ID = resolveHistoryID(history.ID, history.HostIP)
	history.JoinedAt = resolveCreatedAt(history.JoinedAt)
	query := `INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)`
	_, err := execer.ExecContext(ctx, query, history.ID, history.HostIP, history.JoinedAt, history.User)
	if err != nil {
		return apperror.Wrap(err, "recordSSHHistory", map[string]any{"id": history.ID})
	}

	return nil
}

func commitEnrollTx(tx *sql.Tx, hostID string) error {
	if err := tx.Commit(); err != nil {
		return apperror.Wrap(err, "EnrollSSHHost_Commit", map[string]any{"hostId": hostID})
	}

	return nil
}

func runEnrollTransaction(ctx context.Context, tx *sql.Tx, host SSHHost, history SSHHistory) error {
	if err := UpsertSSHHost(ctx, host, tx); err != nil {
		_ = tx.Rollback()

		return apperror.Wrap(err, "EnrollSSHHost_Upsert", map[string]any{"hostId": host.ID})
	}

	if err := recordSSHHistory(ctx, tx, history); err != nil {
		_ = tx.Rollback()

		return apperror.Wrap(err, "EnrollSSHHost_History", map[string]any{"historyId": history.ID})
	}

	return commitEnrollTx(tx, host.ID)
}

// EnrollSSHHost wraps host upsert and history logging into an atomic transaction.
func EnrollSSHHost(ctx context.Context, host SSHHost, history SSHHistory, db *sql.DB) error {
	_ = EnsureSSHTables(db)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap(err, "EnrollSSHHost_BeginTx", map[string]any{"hostId": host.ID})
	}

	return runEnrollTransaction(ctx, tx, host, history)
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

func wrapHostScanError(err error, alias string) (SSHHost, error) {
	return wrapHostByFieldScanError(err, "GetHostByAlias", "alias", alias)
}

const sqlSelectHostFields = `SELECT id, alias, ip, username, COALESCE(port, 22), COALESCE(encrypted_password, ''), created_at FROM ssh_hosts`

// GetHostByAlias retrieves an SSHHost by its alias.
func GetHostByAlias(ctx context.Context, alias string, db *sql.DB) (SSHHost, error) {
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE alias = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, alias).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.CreatedAt,
	)
	if err != nil {
		return wrapHostByFieldScanError(err, "GetHostByAlias", "alias", alias)
	}

	return host, nil
}

// GetHostByIP retrieves an SSHHost by its IP.
func GetHostByIP(ctx context.Context, ip string, db *sql.DB) (SSHHost, error) {
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE ip = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, ip).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.CreatedAt,
	)
	if err != nil {
		return wrapHostByFieldScanError(err, "GetHostByIP", "ip", ip)
	}

	return host, nil
}

// GetHostByID retrieves an SSHHost by its ID.
func GetHostByID(ctx context.Context, id string, db *sql.DB) (SSHHost, error) {
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` WHERE id = ?`

	var host SSHHost
	err := db.QueryRowContext(ctx, query, id).Scan(
		&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.CreatedAt,
	)
	if err != nil {
		return wrapHostByFieldScanError(err, "GetHostByID", "id", id)
	}

	return host, nil
}

// DeleteHostByIP deletes an SSHHost by its IP.
func DeleteHostByIP(ctx context.Context, ip string, db *sql.DB) error {
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

func scanHostRows(rows *sql.Rows) ([]SSHHost, error) {
	var hosts []SSHHost
	for rows.Next() {
		var host SSHHost
		if err := rows.Scan(&host.ID, &host.Alias, &host.IP, &host.Username, &host.Port, &host.EncryptedPassword, &host.CreatedAt); err != nil {
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

// ListHosts retrieves all SSH hosts from the database.
func ListHosts(ctx context.Context, db *sql.DB) ([]SSHHost, error) {
	_ = EnsureSSHTables(db)
	query := sqlSelectHostFields + ` ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListHosts")
	}
	defer rows.Close()
	hosts, err := scanHostRows(rows)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListHosts_Rows")
	}
	return ensureHostsSlice(hosts), nil
}

// LogSSHHistory logs an SSH connection into the database defensively.
func LogSSHHistory(ctx context.Context, h SSHHistory, db *sql.DB) error {
	_ = EnsureSSHTables(db)
	return LogSSHJoin(ctx, h, db)
}
