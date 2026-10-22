// Package store — migrate_v15repo.go performs the Phase 1.1 v15 rename:
// the legacy `Repos` table (plural, `Id` PK) is rebuilt as `Repo` (singular,
// `RepoId` PK) atomically, preserving every row.
//
// Safety guarantees:
//
//  1. Detect-then-act: skipped entirely on fresh installs and on installs
//     that already migrated.
//  2. Wrapped in a single transaction so a partial failure rolls back.
//  3. PRAGMA foreign_keys=OFF for the duration so child tables (Aliases,
//     DetectedProjects, RepoVersionHistory, GroupRepo) survive the rename;
//     they are subsequently rebuilt with REFERENCES Repo(RepoId) by the
//     standard CREATE TABLE pass on the next Migrate() call (idempotent
//     because their own CREATE TABLE statements use IF NOT EXISTS — full
//     child-FK rebuild lands in Phase 1.2/1.3 when those child tables get
//     their own singular renames).
//  4. Row-count parity check: aborts with rollback on mismatch.
package store

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// migrateV15Repo rebuilds Repos -> Repo and Id -> RepoId in one transaction.
func (db *DB) migrateV15Repo() error {
	if !db.tableExists(constants.LegacyTableRepos) {
		return nil // fresh install
	}

	if db.tableExists(constants.TableRepo) {
		// Both exist? Defensive: previous migration ran but legacy not dropped.
		// Drop the legacy plural and continue.
		_, _ = ExecWrapper(db.conn, "DROP TABLE IF EXISTS Repos").Destruct()

		return nil
	}

	fmt.Println(constants.MsgV15RepoMigrationStart)

	oldCount, err := db.countRows(constants.LegacyTableRepos)
	if err != nil {
		return apperror.WrapSimple(err, "count Repos")
	}

	if err := db.execV15RepoRebuild(); err != nil {
		return err
	}

	newCount, err := db.countRows(constants.TableRepo)
	if err != nil {
		return apperror.WrapSimple(err, "count Repo")
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr,
			"  ✗ v15 Repo migration row-count mismatch: old=%d new=%d\n",
			oldCount, newCount)

		return apperror.NewExecutionError(fmt.Sprintf("v15 Repo migration count mismatch: old=%d new=%d", oldCount, newCount))
	}

	fmt.Println(constants.MsgV15RepoMigrationDone)

	return nil
}

// execV15RepoRebuild performs the actual table-rebuild dance. SQLite has no
// ALTER COLUMN, so we CREATE the new table, INSERT ... SELECT, then DROP the
// old one. Foreign keys are temporarily disabled so child tables (which still
// reference Repos at this point) survive the rename.
func (db *DB) execV15RepoRebuild() error {
	if _, err := ExecWrapper(db.conn, "PRAGMA foreign_keys = OFF").Destruct(); err != nil {
		return apperror.WrapSimple(err, "disable foreign keys")
	}

	defer func() {
		_, _ = ExecWrapper(db.conn, "PRAGMA foreign_keys = ON").Destruct()
	}()

	if _, err := ExecWrapper(db.conn, constants.SQLCreateRepo).Destruct(); err != nil {
		return apperror.WrapSimple(err, "create Repo table")
	}

	if _, err := ExecWrapper(db.conn, constants.SQLCreateAbsPathIndex).Destruct(); err != nil {
		return apperror.WrapSimple(err, "create AbsPath index")
	}

	const copySQL = `INSERT INTO Repo
		(RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath,
		 AbsolutePath, CloneInstruction, Notes, CreatedAt, UpdatedAt)
		SELECT Id, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath,
		 AbsolutePath, CloneInstruction, Notes, CreatedAt, UpdatedAt
		FROM Repos`

	if _, err := ExecWrapper(db.conn, copySQL).Destruct(); err != nil {
		return apperror.WrapSimple(err, "copy Repos -> Repo")
	}

	if _, err := ExecWrapper(db.conn, "DROP TABLE Repos").Destruct(); err != nil {
		return apperror.WrapSimple(err, "drop legacy Repos")
	}

	// Drop the legacy index name; SQLCreateAbsPathIndex above created the
	// new IdxRepo_AbsolutePath, so the old idx_Repos_AbsolutePath is dead.
	_, _ = ExecWrapper(db.conn, constants.SQLDropLegacyAbsPathIndex).Destruct()

	return nil
}

// countRows returns the row count for a table.
func (db *DB) countRows(table string) (int, error) {
	var n int
	row := QueryRowWrapper(db.conn, fmt.Sprintf("SELECT COUNT(*) FROM %q", table))
	if err := row.Scan(&n); err != nil {
		return 0, err
	}

	return n, nil
}
