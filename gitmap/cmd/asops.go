package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// upsertSingleRepo persists a single repo's ScanRecord and prints a status.
func upsertSingleRepo(rec model.ScanRecord) *apperror.AppError {
	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.MsgDBUpsertFailed, err))
	}
	defer db.Close()

	if err := syncRepoRecord(db, rec); err != nil {
		return err
	}

	fmt.Printf(constants.MsgAsDBSyncedFmt, rec.RepoName, rec.Slug)

	return nil
}

func syncRepoRecord(db *store.DB, rec model.ScanRecord) *apperror.AppError {
	if err := db.Migrate(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.MsgDBUpsertFailed, err))
	}
	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.MsgDBUpsertFailed, err))
	}

	return nil
}

// registerAlias creates or updates the alias mapping using the just-upserted
// repo. With force=false, conflicting aliases (different slug) abort the run.
func registerAlias(name string, rec model.ScanRecord, force bool) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrListDBFailed, err))
	}
	defer db.Close()

	repos, err := db.FindBySlug(rec.Slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrAsResolveFmt, rec.Slug, err)
		fmt.Fprintln(os.Stderr)

		return apperror.NewSimple("fatal error", "E9000")
	}

	return createOrUpdateAliasRow(db, name, repos[0].ID, rec, force)
}

// createOrUpdateAliasRow handles the conflict-detection + write.
func createOrUpdateAliasRow(db *store.DB, name string, repoID int64, rec model.ScanRecord, force bool) *apperror.AppError {
	hasAlias := db.AliasExists(name)
	if !hasAlias {
		return createAsAliasAndReturn(db, name, repoID, rec)
	}

	if err := checkAliasConflict(db, name, rec, force); err != nil {
		return err
	}

	return updateAsAliasAndReturn(db, name, repoID, rec)
}

func updateAsAliasAndReturn(db *store.DB, name string, repoID int64, rec model.ScanRecord) *apperror.AppError {
	if err := db.UpdateAlias(name, repoID); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	fmt.Printf(constants.MsgAsUpdatedFmt, name, rec.RepoName, rec.AbsolutePath)
	fmt.Printf(constants.MsgAsHintNext, name)
	renameVSCodePMByPath(rec.AbsolutePath, name)

	return nil
}

func createAsAliasAndReturn(db *store.DB, name string, repoID int64, rec model.ScanRecord) *apperror.AppError {
	if _, err := db.CreateAlias(name, repoID); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	fmt.Printf(constants.MsgAsRegisteredFmt, rec.RepoName, name, rec.AbsolutePath)
	fmt.Printf(constants.MsgAsHintNext, name)
	renameVSCodePMByPath(rec.AbsolutePath, name)

	return nil
}

func checkAliasConflict(db *store.DB, name string, rec model.ScanRecord, force bool) *apperror.AppError {
	if force {
		return nil
	}
	existing, err := db.ResolveAlias(name)
	if err == nil && existing.Slug != rec.Slug {
		fmt.Fprintf(os.Stderr, constants.ErrAsAliasInUseFmt, name, existing.Slug)
		fmt.Fprintln(os.Stderr)

		return apperror.NewSimple("fatal error", "E9000")
	}

	return nil
}
