package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type SSHConnection struct {
	Alias             string    `json:"alias"`
	IPAddress         string    `json:"ip_address"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"encrypted_password"`
	KeyPath           string    `json:"key_path"`
	OS                string    `json:"os"`
	OSGroup           string    `json:"os_group,omitempty"`
	OSVersion         string    `json:"os_version,omitempty"`
	BuildVersion      string    `json:"build_version,omitempty"`
	FirstRunAt        time.Time `json:"first_run_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

const (
	sqlInsertSSHConnection = `
		INSERT OR REPLACE INTO SSHConnection (
			Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, OSGroup, OSVersion, BuildVersion, FirstRunAt, CreatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	sqlUpdateMatchedSSHConnection = `
		UPDATE SSHConnection SET
			Alias = ?,
			IPAddress = ?,
			Username = ?,
			EncryptedPassword = ?,
			KeyPath = ?,
			OS = ?,
			OSGroup = ?,
			OSVersion = ?,
			BuildVersion = ?
		WHERE Alias = ?
	`
	sqlSelectSSHConnections         = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, COALESCE(OSGroup, ''), COALESCE(OSVersion, ''), COALESCE(BuildVersion, ''), FirstRunAt, CreatedAt FROM SSHConnection`
	sqlSelectSSHConnectionByAlias   = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, COALESCE(OSGroup, ''), COALESCE(OSVersion, ''), COALESCE(BuildVersion, ''), FirstRunAt, CreatedAt FROM SSHConnection WHERE Alias = ? LIMIT 1`
	sqlSelectSSHConnectionByIP      = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, COALESCE(OSGroup, ''), COALESCE(OSVersion, ''), COALESCE(BuildVersion, ''), FirstRunAt, CreatedAt FROM SSHConnection WHERE IPAddress = ? LIMIT 1`
	sqlUpdateSSHConnectionPassword  = `UPDATE SSHConnection SET EncryptedPassword = ? WHERE Alias = ? OR IPAddress = ?`
	sqlUpdateSSHConnectionOS        = `UPDATE SSHConnection SET OS = ? WHERE Alias = ? OR IPAddress = ?`
	sqlUpdateSSHConnectionTelemetry = `UPDATE SSHConnection SET OS = ?, OSVersion = ?, BuildVersion = ? WHERE Alias = ? OR IPAddress = ?`
	sqlDeleteSSHConnection          = `DELETE FROM SSHConnection WHERE Alias = ?`
	sqlDeleteSSHConnectionByTarget  = `DELETE FROM SSHConnection WHERE Alias = ? OR IPAddress = ?`
	sqlDeleteAllSSHConnections      = `DELETE FROM SSHConnection`
)

func InsertOrUpdateSSHConnection(ctx context.Context, db *sql.DB, conn SSHConnection) *apperror.AppError {
	ensureSSHConnectionSchema(ctx, db)
	firstRun, createdAt := resolveConnTimestamps(conn.FirstRunAt, conn.CreatedAt)
	return execUpsertSSHConnection(ctx, db, conn, firstRun, createdAt)
}

func resolveConnTimestamps(firstRun, createdAt time.Time) (time.Time, time.Time) {
	now := time.Now().UTC()
	if firstRun.IsZero() {
		firstRun = now
	}
	if createdAt.IsZero() {
		createdAt = now
	}
	return firstRun, createdAt
}

func execUpsertSSHConnection(
	ctx context.Context,
	db *sql.DB,
	c SSHConnection,
	firstRun, createdAt time.Time,
) *apperror.AppError {
	existing := findExistingMatch(ctx, db, c.Alias, c.IPAddress)
	if existing != nil {
		return updateMatchedConnection(ctx, db, existing, c)
	}
	return insertNewConnection(ctx, db, c, firstRun, createdAt)
}

func findExistingMatch(ctx context.Context, db *sql.DB, alias, ip string) *SSHConnection {
	if found := findByNonEmptyAlias(ctx, db, alias); found != nil {
		return found
	}
	return findByNonEmptyIP(ctx, db, ip)
}

func findByNonEmptyAlias(ctx context.Context, db *sql.DB, alias string) *SSHConnection {
	if alias == "" {
		return nil
	}
	found, err := GetSSHConnectionByAlias(ctx, db, alias)
	if err != nil {
		return nil
	}
	return found
}

func findByNonEmptyIP(ctx context.Context, db *sql.DB, ip string) *SSHConnection {
	if ip == "" {
		return nil
	}
	found, err := GetSSHConnectionByIP(ctx, db, ip)
	if err != nil {
		return nil
	}
	return found
}

func updateMatchedConnection(ctx context.Context, db *sql.DB, existing *SSHConnection, incoming SSHConnection) *apperror.AppError {
	merged := mergeMatchedFields(existing, incoming)
	_, err := db.ExecContext(ctx, sqlUpdateMatchedSSHConnection,
		merged.Alias, merged.IPAddress, merged.Username, merged.EncryptedPassword,
		merged.KeyPath, merged.OS, merged.OSGroup, merged.OSVersion, merged.BuildVersion,
		existing.Alias,
	)
	if err != nil {
		return apperror.WrapSimple(err, "updateMatchedConnection")
	}
	return nil
}

func mergeMatchedFields(existing *SSHConnection, incoming SSHConnection) SSHConnection {
	merged := *existing
	if incoming.Alias != "" {
		merged.Alias = incoming.Alias
	}
	if incoming.IPAddress != "" {
		merged.IPAddress = incoming.IPAddress
	}
	if incoming.Username != "" {
		merged.Username = incoming.Username
	}
	if incoming.EncryptedPassword != "" {
		merged.EncryptedPassword = incoming.EncryptedPassword
	}
	if incoming.KeyPath != "" {
		merged.KeyPath = incoming.KeyPath
	}
	if incoming.OS != "" {
		merged.OS = incoming.OS
	}
	if incoming.OSGroup != "" {
		merged.OSGroup = incoming.OSGroup
	}
	if incoming.OSVersion != "" {
		merged.OSVersion = incoming.OSVersion
	}
	if incoming.BuildVersion != "" {
		merged.BuildVersion = incoming.BuildVersion
	}
	return merged
}

func insertNewConnection(
	ctx context.Context,
	db *sql.DB,
	c SSHConnection,
	firstRun, createdAt time.Time,
) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlInsertSSHConnection,
		c.Alias, c.IPAddress, c.Username, c.EncryptedPassword,
		c.KeyPath, c.OS, c.OSGroup, c.OSVersion, c.BuildVersion,
		firstRun, createdAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "insertNewConnection")
	}
	return nil
}

const sqlCreateSSHConnectionTable = `CREATE TABLE IF NOT EXISTS SSHConnection (
	Alias TEXT PRIMARY KEY,
	IPAddress TEXT NOT NULL,
	Username TEXT NOT NULL,
	EncryptedPassword TEXT NOT NULL,
	KeyPath TEXT,
	OS TEXT DEFAULT 'linux',
	OSGroup TEXT DEFAULT '',
	OSVersion TEXT DEFAULT '',
	BuildVersion TEXT DEFAULT '',
	FirstRunAt TIMESTAMP,
	CreatedAt TIMESTAMP NOT NULL
);`

func ensureSSHConnectionSchema(ctx context.Context, db *sql.DB) {
	_, _ = db.ExecContext(ctx, sqlCreateSSHConnectionTable)
	_, _ = db.ExecContext(ctx, "ALTER TABLE SSHConnection ADD COLUMN OSGroup TEXT DEFAULT ''")
	_, _ = db.ExecContext(ctx, "ALTER TABLE SSHConnection ADD COLUMN OSVersion TEXT DEFAULT ''")
	_, _ = db.ExecContext(ctx, "ALTER TABLE SSHConnection ADD COLUMN BuildVersion TEXT DEFAULT ''")
	_, _ = db.ExecContext(ctx, "ALTER TABLE SSHConnection ADD COLUMN FirstRunAt TIMESTAMP")
}

func GetSSHConnections(ctx context.Context, db *sql.DB) SSHConnectionSliceResult {
	ensureSSHConnectionSchema(ctx, db)

	rows, err := db.QueryContext(ctx, sqlSelectSSHConnections)
	if err != nil {
		return result.FailSlice[SSHConnection](apperror.WrapSimple(err, "GetSSHConnections.Query"))
	}

	defer rows.Close()

	scanRes := scanSSHConnectionRows(rows)
	if scanRes.IsFailure() {
		return scanRes
	}

	return mergeSSHHostsConnections(ctx, db, scanRes.Data)
}

func mergeSSHHostsConnections(ctx context.Context, db *sql.DB, existing []SSHConnection) SSHConnectionSliceResult {
	seen := make(map[string]bool)
	for _, c := range existing {
		seen[c.Alias] = true
		seen[c.IPAddress] = true
	}

	rows, err := db.QueryContext(ctx, `SELECT alias, ip, username, COALESCE(encrypted_password, '') FROM ssh_hosts`)
	if err != nil {
		return result.OkSlice(existing)
	}
	defer rows.Close()

	merged := append([]SSHConnection{}, existing...)
	for rows.Next() {
		merged = scanAndAppendSSHHostRow(rows, seen, merged)
	}

	return result.OkSlice(merged)
}

func applySSHHostCredentials(conn *SSHConnection, encPass, user string) {
	hasMissingPass := conn.EncryptedPassword == "" && encPass != ""
	if hasMissingPass {
		conn.EncryptedPassword = encPass
	}
	hasMissingUser := conn.Username == "" && user != ""
	if hasMissingUser {
		conn.Username = user
	}
}

func updateExistingPassword(merged []SSHConnection, alias, ip, encPass, user string) {
	for i := range merged {
		isMatch := strings.EqualFold(merged[i].Alias, alias) || merged[i].IPAddress == ip
		if isMatch {
			applySSHHostCredentials(&merged[i], encPass, user)
		}
	}
}

func scanAndAppendSSHHostRow(rows *sql.Rows, seen map[string]bool, merged []SSHConnection) []SSHConnection {
	var alias, ip, user, encPass string
	if err := rows.Scan(&alias, &ip, &user, &encPass); err != nil {
		return merged
	}
	hasSeen := seen[alias] || seen[ip]
	if hasSeen {
		updateExistingPassword(merged, alias, ip, encPass, user)
		return merged
	}

	seen[alias] = true
	seen[ip] = true
	initialOS := "linux"
	if strings.EqualFold(user, "Administrator") {
		initialOS = "windows"
	}
	return append(merged, SSHConnection{
		Alias:             alias,
		IPAddress:         ip,
		Username:          user,
		EncryptedPassword: encPass,
		OS:                initialOS,
	})
}

func assignFirstRunOrCreated(c *SSHConnection, firstRun sql.NullTime) {
	if firstRun.Valid {
		c.FirstRunAt = firstRun.Time
		return
	}
	c.FirstRunAt = c.CreatedAt
}

func scanSSHConnectionRow(rows *sql.Rows) (SSHConnection, error) {
	var c SSHConnection
	var firstRun sql.NullTime
	err := rows.Scan(
		&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword,
		&c.KeyPath, &c.OS, &c.OSGroup, &c.OSVersion, &c.BuildVersion,
		&firstRun, &c.CreatedAt,
	)
	if err != nil {
		return c, err
	}
	assignFirstRunOrCreated(&c, firstRun)
	return c, nil
}

func scanSSHConnectionRows(rows *sql.Rows) SSHConnectionSliceResult {
	var conns []SSHConnection
	for rows.Next() {
		c, err := scanSSHConnectionRow(rows)
		if err != nil {
			return result.FailSlice[SSHConnection](apperror.WrapSimple(err, "scanSSHConnectionRows.Scan"))
		}
		conns = append(conns, c)
	}

	if err := rows.Err(); err != nil {
		return result.FailSlice[SSHConnection](apperror.WrapSimple(err, "scanSSHConnectionRows.Rows"))
	}

	return result.OkSlice(conns)
}

func DeleteSSHConnection(ctx context.Context, db *sql.DB, alias string) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlDeleteSSHConnection, alias)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteSSHConnection.Exec")
	}

	return nil
}

// UpdateConnectionOS updates the recorded OS of an SSH connection.
func UpdateConnectionOS(ctx context.Context, db *sql.DB, target, osType string) *apperror.AppError {
	ensureSSHConnectionSchema(ctx, db)
	_, err := db.ExecContext(ctx, sqlUpdateSSHConnectionOS, osType, target, target)
	if err != nil {
		return apperror.WrapSimple(err, "UpdateConnectionOS.Exec")
	}

	return nil
}

// UpdateConnectionTelemetry updates the recorded OS, OSVersion, and build/gitmap version of an SSH connection.
func UpdateConnectionTelemetry(ctx context.Context, db *sql.DB, target, osType, osVersion, buildVersion string) *apperror.AppError {
	ensureSSHConnectionSchema(ctx, db)
	_, err := db.ExecContext(ctx, sqlUpdateSSHConnectionTelemetry, osType, osVersion, buildVersion, target, target)
	if err != nil {
		return apperror.WrapSimple(err, "UpdateConnectionTelemetry.Exec")
	}

	return nil
}

// DeleteSSHConnectionByTarget deletes SSH connections matching alias or IP address.
func DeleteSSHConnectionByTarget(ctx context.Context, db *sql.DB, target string) error {
	_, _ = db.ExecContext(ctx, sqlCreateSSHConnectionTable)
	_, err := db.ExecContext(ctx, sqlDeleteSSHConnectionByTarget, target, target)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteSSHConnectionByTarget.Exec")
	}

	return nil
}

// DeleteAllSSHConnections removes all records from the SSHConnection table.
func DeleteAllSSHConnections(ctx context.Context, db *sql.DB) error {
	_, _ = db.ExecContext(ctx, sqlCreateSSHConnectionTable)
	_, err := db.ExecContext(ctx, sqlDeleteAllSSHConnections)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteAllSSHConnections.Exec")
	}

	return nil
}

// UpdateSSHConnectionPassword updates the encrypted password for a connection identified by alias or IP.
func UpdateSSHConnectionPassword(ctx context.Context, db *sql.DB, target string, encryptedPass string) (int64, error) {
	_, _ = db.ExecContext(ctx, sqlCreateSSHConnectionTable)
	res, err := db.ExecContext(ctx, sqlUpdateSSHConnectionPassword, encryptedPass, target, target)
	if err != nil {
		return 0, apperror.WrapSimple(err, "UpdateSSHConnectionPassword.Exec")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

func scanSSHConnectionRowCtx(row *sql.Row) (*SSHConnection, error) {
	var c SSHConnection
	var firstRun sql.NullTime
	err := row.Scan(
		&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword,
		&c.KeyPath, &c.OS, &c.OSGroup, &c.OSVersion, &c.BuildVersion,
		&firstRun, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	assignFirstRunOrCreated(&c, firstRun)
	return &c, nil
}

// GetSSHConnectionByAlias retrieves an SSHConnection by its Alias.
func GetSSHConnectionByAlias(ctx context.Context, db *sql.DB, alias string) (*SSHConnection, error) {
	ensureSSHConnectionSchema(ctx, db)
	return scanSSHConnectionRowCtx(db.QueryRowContext(ctx, sqlSelectSSHConnectionByAlias, alias))
}

// GetSSHConnectionByIP retrieves an SSHConnection by its IP Address.
func GetSSHConnectionByIP(ctx context.Context, db *sql.DB, ip string) (*SSHConnection, error) {
	ensureSSHConnectionSchema(ctx, db)
	return scanSSHConnectionRowCtx(db.QueryRowContext(ctx, sqlSelectSSHConnectionByIP, ip))
}
