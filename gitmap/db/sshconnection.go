package db

import (
	"context"
	"database/sql"
	"time"
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

func InsertOrUpdateSSHConnection(ctx context.Context, db *sql.DB, conn SSHConnection) error {
	query := `
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
	_, err := db.ExecContext(ctx, query,
		conn.Alias,
		conn.IPAddress,
		conn.Username,
		conn.EncryptedPassword,
		conn.KeyPath,
		conn.OS,
		conn.CreatedAt,
	)
	return err
}

func GetSSHConnections(ctx context.Context, db *sql.DB) ([]SSHConnection, error) {
	query := `SELECT Alias, IPAddress, Username, EncryptedPassword, KeyPath, OS, CreatedAt FROM SSHConnection`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conns []SSHConnection
	for rows.Next() {
		var c SSHConnection
		if err := rows.Scan(&c.Alias, &c.IPAddress, &c.Username, &c.EncryptedPassword, &c.KeyPath, &c.OS, &c.CreatedAt); err != nil {
			return nil, err
		}
		conns = append(conns, c)
	}
	return conns, rows.Err()
}

func DeleteSSHConnection(ctx context.Context, db *sql.DB, alias string) error {
	query := `DELETE FROM SSHConnection WHERE Alias = ?`
	_, err := db.ExecContext(ctx, query, alias)
	return err
}
