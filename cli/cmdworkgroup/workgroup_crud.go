package cmdworkgroup

import (
	"database/sql"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var openDBFunc = func() (*store.DB, error) {
	return store.OpenDefault()
}

const sqlCreateWorkgroup = `CREATE TABLE IF NOT EXISTS Workgroup (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE COLLATE NOCASE,
	root_directory TEXT NOT NULL COLLATE NOCASE,
	ssh_key_path TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	git_user_name TEXT NOT NULL DEFAULT '',
	git_user_email TEXT NOT NULL DEFAULT '',
	is_default INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

const sqlCreateFolderBinding = `CREATE TABLE IF NOT EXISTS FolderBinding (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	folder_path TEXT NOT NULL UNIQUE COLLATE NOCASE,
	workgroup_id INTEGER NOT NULL DEFAULT 0,
	ssh_key_path TEXT NOT NULL DEFAULT '',
	git_user_name TEXT NOT NULL DEFAULT '',
	git_user_email TEXT NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

// EnsureWorkgroupTables provisions the Workgroup and FolderBinding SQLite schemas.
func EnsureWorkgroupTables(db *store.DB) error {
	if _, err := db.SQL().Exec(sqlCreateWorkgroup); err != nil {
		return apperror.WrapSimple(err, "create Workgroup table")
	}
	if _, err := db.SQL().Exec(sqlCreateFolderBinding); err != nil {
		return apperror.WrapSimple(err, "create FolderBinding table")
	}

	return nil
}

func clearWorkgroupDefaults(db *store.DB) error {
	_, err := db.SQL().Exec("UPDATE Workgroup SET is_default = 0")
	if err != nil {
		return apperror.WrapSimple(err, "clear workgroup defaults")
	}

	return nil
}

func insertWorkgroupRecord(db *store.DB, rec WorkgroupRecord) (int64, error) {
	isDefInt := 0
	if rec.IsDefault {
		isDefInt = 1
	}
	query := `INSERT INTO Workgroup (name, root_directory, ssh_key_path, description, git_user_name, git_user_email, is_default)
VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := db.SQL().Exec(query, rec.Name, rec.RootDirectory, rec.SSHKeyPath, rec.Description, rec.GitUserName, rec.GitUserEmail, isDefInt)
	if err != nil {
		return 0, apperror.WrapSimple(err, "insert workgroup")
	}

	return res.LastInsertId()
}

func scanWorkgroupRow(row *sql.Row) (WorkgroupRecord, error) {
	var r WorkgroupRecord
	var isDefInt int
	err := row.Scan(&r.ID, &r.Name, &r.RootDirectory, &r.SSHKeyPath, &r.Description, &r.GitUserName, &r.GitUserEmail, &isDefInt, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}
	r.IsDefault = (isDefInt == 1)

	return r, nil
}

func scanWorkgroupRows(rows *sql.Rows) ([]WorkgroupRecord, error) {
	records := make([]WorkgroupRecord, 0)
	for rows.Next() {
		var r WorkgroupRecord
		var isDefInt int
		err := rows.Scan(&r.ID, &r.Name, &r.RootDirectory, &r.SSHKeyPath, &r.Description, &r.GitUserName, &r.GitUserEmail, &isDefInt, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scan workgroup row")
		}
		r.IsDefault = (isDefInt == 1)
		records = append(records, r)
	}

	return records, rows.Err()
}

func queryWorkgroups(db *store.DB) ([]WorkgroupRecord, error) {
	query := `SELECT id, name, root_directory, ssh_key_path, description, git_user_name, git_user_email, is_default, created_at, updated_at FROM Workgroup ORDER BY name ASC`
	rows, err := db.SQL().Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query workgroups")
	}
	defer rows.Close()

	return scanWorkgroupRows(rows)
}

// ListWorkgroups returns all registered workgroups.
func ListWorkgroups() ([]WorkgroupRecord, error) {
	db, err := openDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return nil, errEnsure
	}

	return queryWorkgroups(db)
}

// AddWorkgroup persists a new workgroup record.
func AddWorkgroup(rec WorkgroupRecord) (WorkgroupRecord, error) {
	db, err := openDBFunc()
	if err != nil {
		return WorkgroupRecord{}, apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return WorkgroupRecord{}, errEnsure
	}
	if rec.IsDefault {
		_ = clearWorkgroupDefaults(db)
	}
	if abs, errAbs := filepath.Abs(rec.RootDirectory); errAbs == nil {
		rec.RootDirectory = filepath.Clean(abs)
	}
	id, errInsert := insertWorkgroupRecord(db, rec)
	if errInsert != nil {
		return WorkgroupRecord{}, errInsert
	}
	rec.ID = id

	return rec, nil
}

func querySingleWorkgroup(db *store.DB, target string) (WorkgroupRecord, error) {
	query := `SELECT id, name, root_directory, ssh_key_path, description, git_user_name, git_user_email, is_default, created_at, updated_at FROM Workgroup WHERE name = ? COLLATE NOCASE OR root_directory = ? COLLATE NOCASE OR id = ? LIMIT 1`
	row := db.SQL().QueryRow(query, target, target, target)

	return scanWorkgroupRow(row)
}

// GetWorkgroup fetches a workgroup by name, id, or root path.
func GetWorkgroup(target string) (WorkgroupRecord, error) {
	db, err := openDBFunc()
	if err != nil {
		return WorkgroupRecord{}, apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return WorkgroupRecord{}, errEnsure
	}

	return querySingleWorkgroup(db, target)
}

func executeDeleteWorkgroup(db *store.DB, target string) error {
	query := `DELETE FROM Workgroup WHERE name = ? COLLATE NOCASE OR id = ?`
	_, err := db.SQL().Exec(query, target, target)
	if err != nil {
		return apperror.WrapSimple(err, "delete workgroup")
	}

	return nil
}

// DeleteWorkgroup removes a workgroup by name or id.
func DeleteWorkgroup(target string) error {
	db, err := openDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return errEnsure
	}

	return executeDeleteWorkgroup(db, target)
}

func insertFolderBinding(db *store.DB, rec FolderBindingRecord) (int64, error) {
	query := `INSERT OR REPLACE INTO FolderBinding (folder_path, workgroup_id, ssh_key_path, git_user_name, git_user_email) VALUES (?, ?, ?, ?, ?)`
	res, err := db.SQL().Exec(query, rec.FolderPath, rec.WorkgroupID, rec.SSHKeyPath, rec.GitUserName, rec.GitUserEmail)
	if err != nil {
		return 0, apperror.WrapSimple(err, "insert folder binding")
	}

	return res.LastInsertId()
}

// BindFolder persists a folder-to-SSH key binding.
func BindFolder(rec FolderBindingRecord) (FolderBindingRecord, error) {
	db, err := openDBFunc()
	if err != nil {
		return FolderBindingRecord{}, apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return FolderBindingRecord{}, errEnsure
	}
	if abs, errAbs := filepath.Abs(rec.FolderPath); errAbs == nil {
		rec.FolderPath = filepath.Clean(abs)
	}
	id, errInsert := insertFolderBinding(db, rec)
	if errInsert != nil {
		return FolderBindingRecord{}, errInsert
	}
	rec.ID = id

	return rec, nil
}

func executeUnbindFolder(db *store.DB, folderPath string) error {
	query := `DELETE FROM FolderBinding WHERE folder_path = ? COLLATE NOCASE`
	_, err := db.SQL().Exec(query, folderPath)
	if err != nil {
		return apperror.WrapSimple(err, "unbind folder")
	}

	return nil
}

// UnbindFolder removes an explicit folder SSH key binding.
func UnbindFolder(folderPath string) error {
	db, err := openDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return errEnsure
	}
	target := folderPath
	if abs, errAbs := filepath.Abs(folderPath); errAbs == nil {
		target = filepath.Clean(abs)
	}

	return executeUnbindFolder(db, target)
}

func scanFolderBindingRows(rows *sql.Rows) ([]FolderBindingRecord, error) {
	records := make([]FolderBindingRecord, 0)
	for rows.Next() {
		var r FolderBindingRecord
		err := rows.Scan(&r.ID, &r.FolderPath, &r.WorkgroupID, &r.SSHKeyPath, &r.GitUserName, &r.GitUserEmail, &r.CreatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scan folder binding")
		}
		records = append(records, r)
	}

	return records, rows.Err()
}

func queryFolderBindings(db *store.DB) ([]FolderBindingRecord, error) {
	query := `SELECT id, folder_path, workgroup_id, ssh_key_path, git_user_name, git_user_email, created_at FROM FolderBinding ORDER BY folder_path ASC`
	rows, err := db.SQL().Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query folder bindings")
	}
	defer rows.Close()

	return scanFolderBindingRows(rows)
}

// ListFolderBindings lists all explicit folder bindings.
func ListFolderBindings() ([]FolderBindingRecord, error) {
	db, err := openDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return nil, errEnsure
	}

	return queryFolderBindings(db)
}

func findExactFolderBinding(db *store.DB, p string) (FolderBindingRecord, error) {
	query := `SELECT id, folder_path, workgroup_id, ssh_key_path, git_user_name, git_user_email, created_at FROM FolderBinding WHERE folder_path = ? COLLATE NOCASE LIMIT 1`
	row := db.SQL().QueryRow(query, p)
	var r FolderBindingRecord
	err := row.Scan(&r.ID, &r.FolderPath, &r.WorkgroupID, &r.SSHKeyPath, &r.GitUserName, &r.GitUserEmail, &r.CreatedAt)

	return r, err
}

func findWorkgroupForRoot(db *store.DB, p string) (WorkgroupRecord, error) {
	query := `SELECT id, name, root_directory, ssh_key_path, description, git_user_name, git_user_email, is_default, created_at, updated_at FROM Workgroup WHERE root_directory = ? COLLATE NOCASE LIMIT 1`
	row := db.SQL().QueryRow(query, p)

	return scanWorkgroupRow(row)
}

func findDefaultWorkgroup(db *store.DB) (WorkgroupRecord, error) {
	query := `SELECT id, name, root_directory, ssh_key_path, description, git_user_name, git_user_email, is_default, created_at, updated_at FROM Workgroup WHERE is_default = 1 LIMIT 1`
	row := db.SQL().QueryRow(query)

	return scanWorkgroupRow(row)
}

func buildEffectiveFromBinding(b FolderBindingRecord, isInherited bool) EffectiveWorkgroupCredentials {
	return EffectiveWorkgroupCredentials{
		RootDirectory:  b.FolderPath,
		SSHKeyPath:     b.SSHKeyPath,
		GitUserName:    b.GitUserName,
		GitUserEmail:   b.GitUserEmail,
		HasSSHKey:      b.SSHKeyPath != "",
		HasGitIdentity: b.GitUserName != "" || b.GitUserEmail != "",
		IsInherited:    isInherited,
		HasActiveGroup: true,
	}
}

func buildEffectiveFromWorkgroup(w WorkgroupRecord, isInherited bool) EffectiveWorkgroupCredentials {
	return EffectiveWorkgroupCredentials{
		WorkgroupName:  w.Name,
		RootDirectory:  w.RootDirectory,
		SSHKeyPath:     w.SSHKeyPath,
		GitUserName:    w.GitUserName,
		GitUserEmail:   w.GitUserEmail,
		HasSSHKey:      w.SSHKeyPath != "",
		HasGitIdentity: w.GitUserName != "" || w.GitUserEmail != "",
		IsInherited:    isInherited,
		HasActiveGroup: true,
	}
}

func resolveAncestorMatch(db *store.DB, absPath string) (EffectiveWorkgroupCredentials, bool) {
	curr := absPath
	for {
		if b, err := findExactFolderBinding(db, curr); err == nil {
			return buildEffectiveFromBinding(b, curr != absPath), true
		}
		if w, err := findWorkgroupForRoot(db, curr); err == nil {
			return buildEffectiveFromWorkgroup(w, curr != absPath), true
		}
		parent := filepath.Dir(curr)
		if parent == curr || parent == "." || parent == "/" || parent == "" {
			break
		}
		curr = parent
	}

	return EffectiveWorkgroupCredentials{}, false
}

// ResolveWorkgroupForPath resolves the effective workgroup credentials for a directory path.
func ResolveWorkgroupForPath(targetPath string) (EffectiveWorkgroupCredentials, error) {
	db, err := openDBFunc()
	if err != nil {
		return EffectiveWorkgroupCredentials{}, apperror.WrapSimple(err, "open db")
	}
	defer db.Close()

	if errEnsure := EnsureWorkgroupTables(db); errEnsure != nil {
		return EffectiveWorkgroupCredentials{}, errEnsure
	}

	absPath, _ := filepath.Abs(targetPath)
	cleanPath := filepath.Clean(absPath)
	if eff, hasMatch := resolveAncestorMatch(db, cleanPath); hasMatch {
		return eff, nil
	}
	if defWg, errDef := findDefaultWorkgroup(db); errDef == nil {
		return buildEffectiveFromWorkgroup(defWg, true), nil
	}

	return EffectiveWorkgroupCredentials{}, nil
}
