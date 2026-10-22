// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InsertSSHKey stores a new SSH key record.
func (db *DB) InsertSSHKey(
	name,
	privatePath,
	publicKey,
	fingerprint,
	email string,
) (model.SSHKey, error) {
	_, err := ExecWrapper(db.conn, constants.SQLInsertSSHKey, name, privatePath, publicKey, fingerprint, email).Destruct()
	if err != nil {
		return model.SSHKey{}, apperror.WrapSimple(err, "failed to create SSH key record")
	}

	return db.FindSSHKeyByName(name)
}

// UpdateSSHKey updates an existing SSH key record by name.
func (db *DB) UpdateSSHKey(name, privatePath, publicKey, fingerprint, email string) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpdateSSHKey, privatePath, publicKey, fingerprint, email, name).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to create SSH key record")
	}

	return nil
}

// FindSSHKeyByName retrieves a single SSH key by its name.
func (db *DB) FindSSHKeyByName(name string) (model.SSHKey, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectSSHKeyByName, name)

	return scanOneSSHKey(row)
}

// ListSSHKeysWrap returns all stored SSH keys encapsulated in SSHKeySliceWrap.
func (db *DB) ListSSHKeysWrap() SSHKeySliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllSSHKeys).Destruct()
	if err != nil {
		return result.FailSlice[model.SSHKey](apperror.WrapSimple(err, constants.ErrSSHQuery))
	}
	defer rows.Close()

	keys, err := scanSSHKeyRows(rows)
	if err != nil {
		return result.FailSlice[model.SSHKey](apperror.WrapSimple(err, constants.ErrSSHQuery))
	}

	return result.OkSlice(keys)
}

// ListSSHKeys returns all stored SSH keys ordered by name.
func (db *DB) ListSSHKeys() ([]model.SSHKey, error) {
	res := db.ListSSHKeysWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

// DeleteSSHKey removes an SSH key record by name.
func (db *DB) DeleteSSHKey(name string) error {
	result, err := ExecWrapper(db.conn, constants.SQLDeleteSSHKeyByName, name).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to delete SSH key")
	}

	return checkDeleted(result, name)
}

// SSHKeyExists returns true if an SSH key with the given name exists.
func (db *DB) SSHKeyExists(name string) bool {
	_, err := db.FindSSHKeyByName(name)

	return err == nil
}

// SSHKeyNamesWrap returns a list of all SSH key names encapsulated in result.StringSliceResult.
func (db *DB) SSHKeyNamesWrap() result.StringSliceResult {
	keysRes := db.ListSSHKeysWrap()
	if keysRes.HasError() {
		return result.FailSlice[string](keysRes.AppError())
	}

	keys := keysRes.Items()
	names := make([]string, len(keys))
	for i, k := range keys {
		names[i] = k.Name
	}

	return result.OkSlice(names)
}

// SSHKeyNames returns a list of all SSH key names.
func (db *DB) SSHKeyNames() ([]string, error) {
	res := db.SSHKeyNamesWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
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
