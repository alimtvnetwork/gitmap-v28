package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isObjectOfType(dbConn *sql.DB, name, objType string) bool {
	var count int
	query := "SELECT 1 FROM sqlite_master WHERE type=? AND name=?"
	err := dbConn.QueryRow(query, objType, name).Scan(&count)
	hasMatch := (err == nil && count == 1)

	return hasMatch
}

func dropViewIfPresent(dbConn *sql.DB, viewName string) error {
	hasView := isObjectOfType(dbConn, viewName, "view")
	if !hasView {
		return nil
	}

	_, errExec := dbConn.Exec("DROP VIEW IF EXISTS " + viewName)
	if errExec != nil {
		return apperror.WrapSimple(errExec, "dropViewIfPresent")
	}

	return nil
}

func copyLegacyWorkDirs(dbConn *sql.DB) error {
	const sqlBackfill = `INSERT OR IGNORE INTO WorkDirectory (
		WorkDirectoryId, AbsolutePath, Label, IsDefault, HasActiveBinding, CreatedAt, UpdatedAt
	) SELECT id, absolute_path, COALESCE(label, ''), is_default, 1, created_at, updated_at FROM work_directories;`

	_, err := dbConn.Exec(sqlBackfill)
	if err != nil {
		return apperror.WrapSimple(err, "copyLegacyWorkDirs")
	}

	return nil
}

func dropLegacyWorkDirsTable(dbConn *sql.DB) error {
	_, err := dbConn.Exec("DROP TABLE IF EXISTS work_directories")
	if err != nil {
		return apperror.WrapSimple(err, "dropLegacyWorkDirsTable")
	}

	return nil
}

func migrateLegacyWorkDirsData(dbConn *sql.DB) error {
	hasTable := isObjectOfType(dbConn, "work_directories", "table")
	if !hasTable {
		return nil
	}

	if errCopy := copyLegacyWorkDirs(dbConn); errCopy != nil {
		return errCopy
	}

	return dropLegacyWorkDirsTable(dbConn)
}

func ensureWorkDirsCompatView(dbConn *sql.DB) error {
	hasView := isObjectOfType(dbConn, constants.ViewWorkDirectories, "view")
	if hasView {
		return nil
	}

	_, err := dbConn.Exec(constants.SQLCreateWorkDirectoriesCompatView)
	if err != nil {
		return apperror.WrapSimple(err, "ensureWorkDirsCompatView")
	}

	return nil
}

func bumpSchemaVersionTo34(dbConn *sql.DB) error {
	hasSetting := isObjectOfType(dbConn, constants.TableSetting, "table")
	if !hasSetting {
		return nil
	}

	const sqlBump = `INSERT INTO Setting (Key, Value) VALUES ('schema_version', '34')
		ON CONFLICT(Key) DO UPDATE SET Value='34'`
	if _, err := dbConn.Exec(sqlBump); err != nil {
		return apperror.WrapSimple(err, "bumpSchemaVersionTo34")
	}

	return nil
}

var workgroupDDLStatements = []string{
	constants.SQLCreateWorkDirectoryTable,
	constants.SQLCreateWorkDirectoryAbsPathIndex,
	constants.SQLCreateWorkDirectoryIsDefaultIndex,
	constants.SQLCreateFolderWorkDirTable,
	constants.SQLCreateFolderWorkDirFolderPathIndex,
	constants.SQLCreateFolderWorkDirWorkDirIdIndex,
	constants.SQLCreateFolderWorkDirPriorityIndex,
	constants.SQLCreateFolderGroupTable,
	constants.SQLCreateFolderGroupGroupNameIndex,
	constants.SQLCreateFolderGroupFolderPathIndex,
	constants.SQLCreateSshKeyAssociationTable,
	constants.SQLCreateSshKeyAssocTargetKeyIndex,
	constants.SQLCreateSshKeyAssocSshKeyIdIndex,
	constants.SQLCreateSshKeyAssocFolderPathIndex,
	constants.SQLCreateSshKeyAssocGroupNameIndex,
	constants.SQLCreateSshKeyAssocWorkDirIdIndex,
	constants.SQLCreateSshKeyAssocRepoIdIndex,
	constants.SQLCreateSshKeyAssocPriorityIndex,
	constants.SQLCreateWorkDirectoriesCompatView,
}

func getWorkgroupDDLStatements() []string {
	return workgroupDDLStatements
}

func executeStatements(dbConn *sql.DB, stmts []string) error {
	for _, stmt := range stmts {
		if _, err := dbConn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "executeStatements")
		}
	}

	return nil
}

func runWorkgroupDDLAndMigration(dbConn *sql.DB) error {
	if errDrop := dropViewIfPresent(dbConn, constants.TableWorkDirectory); errDrop != nil {
		return errDrop
	}

	if errDDL := executeStatements(dbConn, getWorkgroupDDLStatements()); errDDL != nil {
		return errDDL
	}

	return migrateLegacyWorkDirsData(dbConn)
}

// Migration_AddFolderWorkDirAndSshKeyAutomation executes the schema version 34 migration.
func Migration_AddFolderWorkDirAndSshKeyAutomation(dbConn *sql.DB) error {
	if err := runWorkgroupDDLAndMigration(dbConn); err != nil {
		return err
	}

	if errCompat := ensureWorkDirsCompatView(dbConn); errCompat != nil {
		return errCompat
	}

	return bumpSchemaVersionTo34(dbConn)
}

// MigrateWorkgroup applies the v34 workgroup migration to the active database.
func (db *DB) MigrateWorkgroup() error {
	return Migration_AddFolderWorkDirAndSshKeyAutomation(db.conn)
}

// EnsureWorkgroupSchema guarantees that all v34 workgroup tables, indexes, and views exist.
func (db *DB) EnsureWorkgroupSchema() error {
	return Migration_AddFolderWorkDirAndSshKeyAutomation(db.conn)
}
