// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// InsertSSHKey stores a new SSH key record.
func (db *DB) InsertSSHKey(name, privatePath, publicKey, fingerprint, email string) (model.SSHKey, error) {
	_, err := db.conn.Exec(constants.SQLInsertSSHKey, name, privatePath, publicKey, fingerprint, email)
	if err != nil {
		return model.SSHKey{}, fmt.Errorf(constants.ErrSSHCreate, err)
	}

	return db.FindSSHKeyByName(name)
}

// UpdateSSHKey updates an existing SSH key record by name.
func (db *DB) UpdateSSHKey(name, privatePath, publicKey, fingerprint, email string) error {
	_, err := db.conn.Exec(constants.SQLUpdateSSHKey, privatePath, publicKey, fingerprint, email, name)
	if err != nil {
		return fmt.Errorf(constants.ErrSSHCreate, err)
	}

	return nil
}

// FindSSHKeyByName retrieves a single SSH key by its name.
func (db *DB) FindSSHKeyByName(name string) (model.SSHKey, error) {
	row := db.conn.QueryRow(constants.SQLSelectSSHKeyByName, name)

	return scanOneSSHKey(row)
}

// ListSSHKeys returns all stored SSH keys ordered by name.
func (db *DB) ListSSHKeys() ([]model.SSHKey, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllSSHKeys)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrSSHQuery, err)
	}
	defer rows.Close()

	return scanSSHKeyRows(rows)
}

// DeleteSSHKey removes an SSH key record by name.
func (db *DB) DeleteSSHKey(name string) error {
	result, err := db.conn.Exec(constants.SQLDeleteSSHKeyByName, name)
	if err != nil {
		return fmt.Errorf(constants.ErrSSHDelete, err)
	}

	return checkDeleted(result, name)
}

// SSHKeyExists returns true if an SSH key with the given name exists.
func (db *DB) SSHKeyExists(name string) bool {
	_, err := db.FindSSHKeyByName(name)

	return err == nil
}

// SSHKeyNames returns a list of all SSH key names.
func (db *DB) SSHKeyNames() ([]string, error) {
	keys, err := db.ListSSHKeys()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(keys))
	for i, k := range keys {
		names[i] = k.Name
	}

	return names, nil
}

// scanOneSSHKey scans a single SSH key row.
func scanOneSSHKey(row *sql.Row) (model.SSHKey, error) {
	var k model.SSHKey

	err := row.Scan(&k.ID, &k.Name, &k.PrivatePath, &k.PublicKey, &k.Fingerprint, &k.Email, &k.CreatedAt)
	if err != nil {
		return model.SSHKey{}, err
	}

	return k, nil
}

// scanSSHKeyRows scans multiple SSH key rows.
func scanSSHKeyRows(rows *sql.Rows) ([]model.SSHKey, error) { //nolint:unparam // error kept for interface consistency
	var keys []model.SSHKey

	for rows.Next() {
		var k model.SSHKey

		err := rows.Scan(&k.ID, &k.Name, &k.PrivatePath, &k.PublicKey, &k.Fingerprint, &k.Email, &k.CreatedAt)
		if err != nil {
			continue
		}

		keys = append(keys, k)
	}

	return keys, nil
}
