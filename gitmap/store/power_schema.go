package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	sqlCreatePowerSetting = `CREATE TABLE IF NOT EXISTS PowerSetting (
	PowerSettingId INTEGER PRIMARY KEY AUTOINCREMENT,
	ProfileName TEXT NOT NULL UNIQUE,
	Platform TEXT NOT NULL,
	DisplayTimeoutMinutes INTEGER NOT NULL,
	SleepTimeoutMinutes INTEGER NOT NULL,
	DiskTimeoutMinutes INTEGER NOT NULL DEFAULT 0,
	IsNeverSleep INTEGER NOT NULL DEFAULT 0,
	IsLockDisabled INTEGER NOT NULL DEFAULT 0,
	IsActive INTEGER NOT NULL DEFAULT 0,
	Notes TEXT NULL,
	Comments TEXT NULL,
	UpdatedAt TEXT NOT NULL,
	CreatedAt TEXT NOT NULL
);`

	sqlCreatePowerSettingHistory = `CREATE TABLE IF NOT EXISTS PowerSettingHistory (
	PowerSettingHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
	Action TEXT NOT NULL,
	Platform TEXT NOT NULL,
	DisplayTimeoutMinutes INTEGER NOT NULL,
	SleepTimeoutMinutes INTEGER NOT NULL,
	IsNeverSleep INTEGER NOT NULL DEFAULT 0,
	Notes TEXT NULL,
	Comments TEXT NULL,
	CreatedAt TEXT NOT NULL
);`
)

// EnsurePowerTables ensures that power configuration and history tables exist.
func EnsurePowerTables(conn *sql.DB) error {
	if conn == nil {
		return apperror.NewSimple("database connection is nil", "E_NIL_CONN")
	}

	if _, err := conn.Exec(sqlCreatePowerSetting); err != nil {
		return apperror.WrapSimple(err, "store.createPowerSettingTable")
	}

	if _, err := conn.Exec(sqlCreatePowerSettingHistory); err != nil {
		return apperror.WrapSimple(err, "store.createPowerSettingHistoryTable")
	}

	return nil
}
