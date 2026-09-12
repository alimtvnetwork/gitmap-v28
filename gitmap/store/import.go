package store

import (
	"context"
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ImportAll restores a DatabaseExport into the database using upsert/insert-or-ignore semantics atomically.
func (db *DB) ImportAll(data model.DatabaseExport) error {
	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	appErr = wrap.WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError {
		return db.executeImportAllTx(tx.Tx(), data)
	})
	if appErr != nil {
		return appErr
	}

	return nil
}

func (db *DB) executeImportAllTx(tx *sql.Tx, data model.DatabaseExport) *apperror.AppError {
	if err := db.importPrimaryEntitiesTx(tx, data); err != nil {
		return err
	}

	return db.importSecondaryEntitiesTx(tx, data)
}

func (db *DB) importPrimaryEntitiesTx(tx *sql.Tx, data model.DatabaseExport) *apperror.AppError {
	if err := db.importReposTx(tx, data.Repos); err != nil {
		return apperror.WrapSimple(err, "repos")
	}

	if err := db.importGroupsTx(tx, data.Groups); err != nil {
		return apperror.WrapSimple(err, "groups")
	}

	if err := db.importReleasesTx(tx, data.Releases); err != nil {
		return apperror.WrapSimple(err, "releases")
	}

	return nil
}

func (db *DB) importSecondaryEntitiesTx(tx *sql.Tx, data model.DatabaseExport) *apperror.AppError {
	if err := db.importHistoryTx(tx, data.History); err != nil {
		return apperror.WrapSimple(err, "history")
	}

	if err := db.importBookmarksTx(tx, data.Bookmarks); err != nil {
		return apperror.WrapSimple(err, "bookmarks")
	}

	return nil
}

// importReposTx upserts all repos by ID within a transaction or runner.
func (db *DB) importReposTx(runner sqlExecutor, repos []model.ScanRecord) error {
	for _, r := range repos {
		_, err := ExecWrapper(runner, constants.SQLUpsertRepo,
			r.Slug, r.RepoName, r.HTTPSUrl, r.SSHUrl,
			r.Branch, r.RelativePath, r.AbsolutePath,
			r.CloneInstruction, r.Notes, r.Transport).Destruct()
		if err != nil {
			return err
		}
	}

	return nil
}

// importGroupsTx creates groups and links repos by slug within a transaction.
func (db *DB) importGroupsTx(tx *sql.Tx, groups []model.GroupExport) error {
	for _, ge := range groups {
		if err := db.importOneGroupTx(tx, ge); err != nil {
			return err
		}
	}

	return nil
}

// importOneGroupTx creates a group and links its member repos within a transaction.
func (db *DB) importOneGroupTx(tx *sql.Tx, ge model.GroupExport) error {
	_, err := ExecWrapper(tx, constants.SQLImportInsertGroup,
		ge.Name, ge.Description, ge.Color).Destruct()
	if err != nil {
		return err
	}
	group, err := findGroupByNameRunner(tx, ge.Name)
	if err != nil {
		return err
	}

	return db.linkGroupReposRunner(tx, group.ID, ge.RepoSlugs)
}

// linkGroupReposRunner links repos to a group by resolving slugs.
func (db *DB) linkGroupReposRunner(runner sqlQueryerExecutor, groupID int64, slugs []string) error {
	for _, slug := range slugs {
		repos, err := findBySlugRunner(runner, slug)
		if err != nil || len(repos) == 0 {
			continue
		}

		_, err = ExecWrapper(runner, constants.SQLInsertGroupRepo, groupID, repos[0].ID).Destruct()
		if err != nil {
			return err
		}
	}

	return nil
}

// importReleasesTx upserts all release records within a transaction.
func (db *DB) importReleasesTx(runner sqlExecutor, releases []model.ReleaseRecord) error {
	for _, r := range releases {
		if err := upsertReleaseTx(runner, r); err != nil {
			return err
		}
	}

	return nil
}

// importHistoryTx inserts history records, ignoring duplicates.
func (db *DB) importHistoryTx(runner sqlExecutor, records []model.CommandHistoryRecord) error {
	for _, r := range records {
		_, err := ExecWrapper(runner,
			"INSERT OR IGNORE INTO CommandHistory (Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			r.Command, r.Alias, r.Args, r.Flags,
			r.StartedAt, r.FinishedAt, r.DurationMs, r.ExitCode, r.Summary, r.RepoCount).Destruct()
		if err != nil {
			return err
		}
	}

	return nil
}

// importBookmarksTx inserts bookmarks, ignoring duplicates by name.
func (db *DB) importBookmarksTx(runner sqlExecutor, bookmarks []model.BookmarkRecord) error {
	for _, b := range bookmarks {
		_, err := ExecWrapper(runner, constants.SQLImportInsertBookmark,
			b.Name, b.Command, b.Args, b.Flags).Destruct()
		if err != nil {
			return err
		}
	}

	return nil
}
