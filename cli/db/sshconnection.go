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
	OSVersion         string    `json:"os_version,omitempty"`
	FirstRunAt        time.Time `json:"first_run_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

const (
	sqlUpsertSSHConnection = `
		INSERT INTO SSHConnection (
			Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, OSVersion, FirstRunAt, CreatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(Alias) DO UPDATE SET
			IPAddress = excluded.IPAddress,
			Username = excluded.Username,
			EncryptedPassword = excluded.EncryptedPassword,
			KeyPath = excluded.KeyPath,
			OS = excluded.OS,
			OSVersion = CASE WHEN excluded.OSVersion != '' THEN excluded.OSVersion ELSE SSHConnection.OSVersion END,
			FirstRunAt = COALESCE(SSHConnection.FirstRunAt, excluded.FirstRunAt)
	`
	sqlSelectSSHConnections        = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, COALESCE(OSVersion, ''), COALESCE(FirstRunAt, CreatedAt), CreatedAt FROM SSHConnection`
	sqlSelectSSHConnectionByAlias  = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, COALESCE(OSVersion, ''), COALESCE(FirstRunAt, CreatedAt), CreatedAt FROM SSHConnection WHERE Alias = ? LIMIT 1`
	sqlSelectSSHConnectionByIP     = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, COALESCE(OSVersion, ''), COALESCE(FirstRunAt, CreatedAt), CreatedAt FROM SSHConnection WHERE IPAddress = ? LIMIT 1`
	sqlUpdateSSHConnectionPassword = `UPDATE SSHConnection SET EncryptedPassword = ? WHERE Alias = ? OR IPAddress = ?`
	sqlDeleteSSHConnection         = `DELETE FROM SSHConnection WHERE Alias = ?`
	sqlDeleteSSHConnectionByTarget = `DELETE FROM SSHConnection WHERE Alias = ? OR IPAddress = ?`
	sqlDeleteAllSSHConnections     = `DELETE FROM SSHConnection`
)

func InsertOrUpdateSSHConnection(ctx context.Context, db *sql.DB, conn SSHConnection) *apperror.AppError {
	ensureSSHConnectionSchema(ctx, db)
	now := time.Now().UTC()
	firstRun := conn.FirstRunAt
	if firstRun.IsZero() {
		firstRun = now
	}
	createdAt := conn.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	_, err := db.ExecContext(ctx, sqlUpsertSSHConnection,
		conn.Alias,
		conn.IPAddress,
		conn.Username,
		conn.EncryptedPassword,
		conn.KeyPath,
		conn.OS,
		conn.OSVersion,
		firstRun,
		createdAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "InsertOrUpdateSSHConnection.Exec")
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
	OSVersion TEXT DEFAULT '',
	FirstRunAt TIMESTAMP,
	CreatedAt TIMESTAMP NOT NULL
);`

func ensureSSHConnectionSchema(ctx context.Context, db *sql.DB) {
	_, _ = db.ExecContext(ctx, sqlCreateSSHConnectionTable)
	_, _ = db.ExecContext(ctx, "ALTER TABLE SSHConnection ADD COLUMN OSVersion TEXT DEFAULT ''")
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
	return append(merged, SSHConnection{
		Alias:             alias,
		IPAddress:         ip,
		Username:          user,
		EncryptedPassword: encPass,
		OS:                "linux",
	})
}

func scanSSHConnectionRows(rows *sql.Rows) SSHConnectionSliceResult {
	var conns []SSHConnection
	for rows.Next() {
		var c SSHConnection
		if err := rows.Scan(&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword, &c.KeyPath, &c.OS, &c.OSVersion, &c.FirstRunAt, &c.CreatedAt); err != nil {
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

// GetSSHConnectionByAlias retrieves an SSHConnection by its Alias.
func GetSSHConnectionByAlias(ctx context.Context, db *sql.DB, alias string) (*SSHConnection, error) {
	ensureSSHConnectionSchema(ctx, db)
	var c SSHConnection
	err := db.QueryRowContext(ctx, sqlSelectSSHConnectionByAlias, alias).Scan(
		&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword, &c.KeyPath, &c.OS, &c.OSVersion, &c.FirstRunAt, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetSSHConnectionByIP retrieves an SSHConnection by its IP Address.
func GetSSHConnectionByIP(ctx context.Context, db *sql.DB, ip string) (*SSHConnection, error) {
	ensureSSHConnectionSchema(ctx, db)
	var c SSHConnection
	err := db.QueryRowContext(ctx, sqlSelectSSHConnectionByIP, ip).Scan(
		&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword, &c.KeyPath, &c.OS, &c.OSVersion, &c.FirstRunAt, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
