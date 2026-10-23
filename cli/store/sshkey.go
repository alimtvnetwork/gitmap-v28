// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var sshKeyAlterCols = []string{
	"ALTER TABLE SshKey ADD COLUMN EncryptedPrivateKey TEXT DEFAULT ''",
	"ALTER TABLE SshKey ADD COLUMN KeyType TEXT DEFAULT 'ed25519'",
	"ALTER TABLE SshKey ADD COLUMN Provider TEXT DEFAULT 'github'",
	"ALTER TABLE SshKey ADD COLUMN IsDefault INTEGER DEFAULT 0",
}

func (db *DB) ensureSSHKeyColumns() {
	if db == nil || db.conn == nil {
		return
	}

	for _, col := range sshKeyAlterCols {
		_ = execAlterColumn(db.conn, col)
	}
}

// EnsureSSHKeyColumns ensures the SshKey table has all encrypted vault columns.
func (db *DB) EnsureSSHKeyColumns() error {
	db.ensureSSHKeyColumns()

	return nil
}

// InsertSSHKey stores a new SSH key record, optionally with encrypted private key.
func (db *DB) InsertSSHKey(
	name,
	privatePath,
	publicKey,
	fingerprint,
	email string,
	encrypted ...string,
) (model.SSHKey, error) {
	db.ensureSSHKeyColumns()

	hasEnc := len(encrypted) > 0
	if hasEnc {
		k := model.SSHKey{
			Name:                name,
			PrivatePath:         privatePath,
			PublicKey:           publicKey,
			Fingerprint:         fingerprint,
			Email:               email,
			EncryptedPrivateKey: encrypted[0],
		}

		return db.InsertSSHKeyWithVault(k)
	}

	_, err := ExecWrapper(db.conn, constants.SQLInsertSSHKey, name, privatePath, publicKey, fingerprint, email).Destruct()
	if err != nil {
		return model.SSHKey{}, apperror.WrapSimple(err, "failed to create SSH key record")
	}

	return db.FindSSHKeyByName(name)
}

// InsertSSHKeyWithVault inserts an SSH key with all encrypted vault fields.
func (db *DB) InsertSSHKeyWithVault(k model.SSHKey) (model.SSHKey, error) {
	db.ensureSSHKeyColumns()
	defInt, keyType, prov := resolveSSHKeyDefaults(k)

	_, err := ExecWrapper(
		db.conn, constants.SQLInsertSSHKeyWithVault,
		k.Name, k.PrivatePath, k.PublicKey, k.Fingerprint, k.Email,
		k.EncryptedPrivateKey, keyType, prov, defInt,
	).Destruct()
	if err != nil {
		return model.SSHKey{}, apperror.WrapSimple(err, "failed to insert SSH key with vault")
	}

	return db.FindSSHKeyByName(k.Name)
}

func resolveSSHKeyDefaults(k model.SSHKey) (int, string, string) {
	defInt := 0
	if k.IsDefault {
		defInt = 1
	}

	keyType := k.KeyType
	if keyType == "" {
		keyType = "ed25519"
	}

	prov := k.Provider
	if prov == "" {
		prov = "github"
	}

	return defInt, keyType, prov
}

// UpdateSSHKey updates an existing SSH key record by name, optionally with encrypted key.
func (db *DB) UpdateSSHKey(
	name,
	privatePath,
	publicKey,
	fingerprint,
	email string,
	encrypted ...string,
) error {
	db.ensureSSHKeyColumns()

	hasEnc := len(encrypted) > 0
	if hasEnc {
		return db.UpdateSSHKeyWithVault(name, privatePath, publicKey, fingerprint, email, encrypted[0])
	}

	_, err := ExecWrapper(db.conn, constants.SQLUpdateSSHKey, privatePath, publicKey, fingerprint, email, name).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to update SSH key record")
	}

	return nil
}

// UpdateSSHKeyWithVault updates an SSH key record including EncryptedPrivateKey.
func (db *DB) UpdateSSHKeyWithVault(name, privatePath, publicKey, fingerprint, email, encPrivKey string) error {
	db.ensureSSHKeyColumns()

	query := `UPDATE SshKey SET PrivatePath = ?, PublicKey = ?, Fingerprint = ?, Email = ?, EncryptedPrivateKey = ? WHERE Name = ?`
	_, err := ExecWrapper(db.conn, query, privatePath, publicKey, fingerprint, email, encPrivKey, name).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to update SSH key record with vault")
	}

	return nil
}

// FindSSHKeyByName retrieves a single SSH key by its name.
func (db *DB) FindSSHKeyByName(name string) (model.SSHKey, error) {
	db.ensureSSHKeyColumns()
	row := QueryRowWrapper(db.conn, constants.SQLSelectSSHKeyByName, name)

	return scanOneSSHKey(row)
}

// GetDefaultSSHKey retrieves the default SSH key from the database.
func (db *DB) GetDefaultSSHKey() (model.SSHKey, error) {
	db.ensureSSHKeyColumns()
	row := QueryRowWrapper(db.conn, constants.SQLSelectDefaultSSHKey)

	return scanOneSSHKey(row)
}

// ListSSHKeysWrap returns all stored SSH keys encapsulated in SSHKeySliceWrap.
func (db *DB) ListSSHKeysWrap() SSHKeySliceWrap {
	db.ensureSSHKeyColumns()
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
	var encPriv, keyType, prov sql.NullString
	var defInt sql.NullInt64

	err := row.Scan(
		&k.ID, &k.Name, &k.PrivatePath, &k.PublicKey, &k.Fingerprint, &k.Email,
		&encPriv, &keyType, &prov, &defInt, &k.CreatedAt,
	)
	if err != nil {
		return model.SSHKey{}, err
	}

	assignScannedVaultFields(&k, encPriv, keyType, prov, defInt)

	return k, nil
}

func assignScannedVaultFields(k *model.SSHKey, encPriv, keyType, prov sql.NullString, defInt sql.NullInt64) {
	k.EncryptedPrivateKey = encPriv.String
	k.KeyType = keyType.String
	k.Provider = prov.String
	k.IsDefault = (defInt.Int64 == 1)
}

// scanSSHKeyRows scans multiple SSH key rows.
func scanSSHKeyRows(rows *sql.Rows) ([]model.SSHKey, error) { //nolint:unparam // error kept for interface consistency
	var keys []model.SSHKey

	for rows.Next() {
		var k model.SSHKey
		var encPriv, keyType, prov sql.NullString
		var defInt sql.NullInt64

		err := rows.Scan(
			&k.ID, &k.Name, &k.PrivatePath, &k.PublicKey, &k.Fingerprint, &k.Email,
			&encPriv, &keyType, &prov, &defInt, &k.CreatedAt,
		)
		if err != nil {
			continue
		}

		assignScannedVaultFields(&k, encPriv, keyType, prov, defInt)
		keys = append(keys, k)
	}

	return keys, nil
}
