package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// GetSetting returns the value for a key, or empty string if not found.
func (db *DB) GetSetting(key string) string {
	var value string
	err := db.conn.QueryRow(constants.SQLSelectSetting, key).Scan(&value)
	if err != nil {
		return ""
	}

	return value
}

// SetSetting upserts a key-value pair in the Settings table.
func (db *DB) SetSetting(key, value string) error {
	_, err := db.conn.Exec(constants.SQLUpsertSetting, key, value)
	if err != nil {
		return fmt.Errorf(constants.ErrDBSettingUpsert, err)
	}

	return nil
}

// DeleteSetting removes a key from the Settings table.
func (db *DB) DeleteSetting(key string) error {
	_, err := db.conn.Exec(constants.SQLDeleteSetting, key)
	if err != nil {
		return fmt.Errorf(constants.ErrDBSettingUpsert, err)
	}

	return nil
}
