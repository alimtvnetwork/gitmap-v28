package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// upsertSingleRepo persists a single repo's ScanRecord and prints a status.
func upsertSingleRepo(rec model.ScanRecord) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		os.Exit(1)
	}

	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgAsDBSyncedFmt, rec.RepoName, rec.Slug)
}

// registerAlias creates or updates the alias mapping using the just-upserted
// repo. With force=false, conflicting aliases (different slug) abort the run.
func registerAlias(name string, rec model.ScanRecord, force bool) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	repos, err := db.FindBySlug(rec.Slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrAsResolveFmt, rec.Slug, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	repoID := repos[0].ID
	createOrUpdateAliasRow(db, name, repoID, rec, force)
}

// createOrUpdateAliasRow handles the conflict-detection + write.
func createOrUpdateAliasRow(db *store.DB, name string, repoID int64, rec model.ScanRecord, force bool) {
	if !db.AliasExists(name) {
		if _, err := db.CreateAlias(name, repoID); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
			os.Exit(1)
		}
		fmt.Printf(constants.MsgAsRegisteredFmt, rec.RepoName, name, rec.AbsolutePath)
		fmt.Printf(constants.MsgAsHintNext, name)
		renameVSCodePMByPath(rec.AbsolutePath, name)

		return
	}

	if !force {
		existing, err := db.ResolveAlias(name)
		if err == nil && existing.Slug != rec.Slug {
			fmt.Fprintf(os.Stderr, constants.ErrAsAliasInUseFmt, name, existing.Slug)
			fmt.Fprintln(os.Stderr)
			os.Exit(1)
		}
	}

	if err := db.UpdateAlias(name, repoID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgAsUpdatedFmt, name, rec.RepoName, rec.AbsolutePath)
	fmt.Printf(constants.MsgAsHintNext, name)
	renameVSCodePMByPath(rec.AbsolutePath, name)
}
