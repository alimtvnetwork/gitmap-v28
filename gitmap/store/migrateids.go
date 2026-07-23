package store

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// migrateLegacyIDs detects if the legacy Repos table uses TEXT IDs (UUIDs)
// and rebuilds it with INTEGER PRIMARY KEY AUTOINCREMENT. Runs BEFORE the
// v15 rename so that migrateV15Repo() sees a clean integer-PK Repos table.
// All tables with foreign keys to Repos.Id are dropped and recreated by the
// normal migration.
func (db *DB) migrateLegacyIDs() {
	if !db.hasLegacyTextID(constants.LegacyTableRepos) {
		return
	}

	fmt.Println(constants.MsgLegacyIDMigrationStart)

	db.dropProjectTables()
	db.dropGroupRepos()

	if err := db.rebuildReposTable(); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Legacy ID migration failed: %v\n", err)

		return
	}

	fmt.Println(constants.MsgLegacyIDMigrationDone)
}

// hasLegacyTextID checks if the Id column of a table is TEXT via PRAGMA.
func (db *DB) hasLegacyTextID(table string) bool {
	query := fmt.Sprintf("PRAGMA table_info(%s)", table)

	rows, err := db.conn.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dflt interface{}

		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			continue
		}

		if name == "Id" && colType == "TEXT" {
			return true
		}
	}

	return false
}

// dropProjectTables removes all project detection tables (FK to Repos).
// Includes both v15 singular and legacy plural / pre-Csharp names so any
// installation state can be cleaned before the legacy ID rebuild.
func (db *DB) dropProjectTables() {
	drops := []string{
		constants.SQLDropGoRunnableFile,
		constants.SQLDropGoRunnableFiles,
		constants.SQLDropGoProjectMetadata,
		constants.SQLDropCsharpKeyFile,
		constants.SQLDropCsharpKeyFiles,
		constants.SQLDropCsharpProjectFile,
		constants.SQLDropCsharpProjectFiles,
		constants.SQLDropCsharpProjectMeta,
		constants.SQLDropCsharpProjectMetaLegacy,
		constants.SQLDropDetectedProject,
		constants.SQLDropDetectedProjects,
		constants.SQLDropProjectType,
		constants.SQLDropProjectTypes,
	}

	for _, stmt := range drops {
		if _, err := db.conn.Exec(stmt); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not drop table during migration: %v\n", err)
		}
	}
}

// dropGroupRepos removes the GroupRepos join table (FK to Repos).
func (db *DB) dropGroupRepos() {
	if _, err := db.conn.Exec(constants.SQLDropGroupRepos); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not drop GroupRepos during migration: %v\n", err)
	}
}

// rebuildReposTable recreates the LEGACY Repos table (still plural, still
// `Id` PK) with INTEGER PRIMARY KEY AUTOINCREMENT, preserving all data
// except the old UUID IDs. The subsequent migrateV15Repo() pass then
// renames Repos → Repo and Id → RepoId.
func (db *DB) rebuildReposTable() error {
	const legacyCreate = `CREATE TABLE IF NOT EXISTS Repos (
		Id               INTEGER PRIMARY KEY AUTOINCREMENT,
		Slug             TEXT NOT NULL,
		RepoName         TEXT NOT NULL,
		HttpsUrl         TEXT NOT NULL,
		SshUrl           TEXT NOT NULL,
		Branch           TEXT NOT NULL,
		RelativePath     TEXT NOT NULL,
		AbsolutePath     TEXT NOT NULL,
		CloneInstruction TEXT NOT NULL,
		Notes            TEXT DEFAULT '',
		CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
		UpdatedAt        TEXT DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := db.conn.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("disable foreign keys: %w", err)
	}

	if _, err := db.conn.Exec("ALTER TABLE Repos RENAME TO Repos_legacy"); err != nil {
		return fmt.Errorf("rename Repos to Repos_legacy: %w", err)
	}

	if _, err := db.conn.Exec(legacyCreate); err != nil {
		return fmt.Errorf("create new Repos table: %w", err)
	}

	if _, err := db.conn.Exec(`INSERT INTO Repos (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes)
		SELECT Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes
		FROM Repos_legacy`); err != nil {
		return fmt.Errorf("copy data from Repos_legacy: %w", err)
	}

	if _, err := db.conn.Exec("DROP TABLE IF EXISTS Repos_legacy"); err != nil {
		return fmt.Errorf("drop Repos_legacy: %w", err)
	}

	if _, err := db.conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("re-enable foreign keys: %w", err)
	}

	return nil
}
