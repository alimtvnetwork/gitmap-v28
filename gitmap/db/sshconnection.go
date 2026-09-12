package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type SSHConnection struct {
	Alias             string    `json:"alias"`
	IPAddress         string    `json:"ip_address"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"encrypted_password"`
	KeyPath           string    `json:"key_path"`
	OS                string    `json:"os"`
	CreatedAt         time.Time `json:"created_at"`
}

const (
	sqlUpsertSSHConnection = `
		INSERT INTO SSHConnection (
			Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, CreatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(Alias) DO UPDATE SET
			IPAddress = excluded.IPAddress,
			Username = excluded.Username,
			EncryptedPassword = excluded.EncryptedPassword,
			KeyPath = excluded.KeyPath,
			OS = excluded.OS
	`
	sqlSelectSSHConnections = `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, CreatedAt FROM SSHConnection`
	sqlDeleteSSHConnection  = `DELETE FROM SSHConnection WHERE Alias = ?`
)

func InsertOrUpdateSSHConnection(ctx context.Context, db *sql.DB, conn SSHConnection) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlUpsertSSHConnection,
		conn.Alias,
		conn.IPAddress,
		conn.Username,
		conn.EncryptedPassword,
		conn.KeyPath,
		conn.OS,
		conn.CreatedAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "InsertOrUpdateSSHConnection.Exec")
	}

	return nil
}

func GetSSHConnections(ctx context.Context, db *sql.DB) ([]SSHConnection, *apperror.AppError) {
	rows, err := db.QueryContext(ctx, sqlSelectSSHConnections)
	if err != nil {
		return nil, apperror.WrapSimple(err, "GetSSHConnections.Query")
	}

	defer rows.Close()

	return scanSSHConnectionRows(rows)
}

func scanSSHConnectionRows(rows *sql.Rows) ([]SSHConnection, *apperror.AppError) {
	var conns []SSHConnection
	for rows.Next() {
		var c SSHConnection
		if err := rows.Scan(&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword, &c.KeyPath, &c.OS, &c.CreatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "scanSSHConnectionRows.Scan")
		}

		conns = append(conns, c)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "scanSSHConnectionRows.Rows")
	}

	return conns, nil
}

func DeleteSSHConnection(ctx context.Context, db *sql.DB, alias string) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlDeleteSSHConnection, alias)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteSSHConnection.Exec")
	}

	return nil
}
