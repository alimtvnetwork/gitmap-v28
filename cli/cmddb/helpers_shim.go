package cmddb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PopulateRepoAliases scans unaliased repos and populates auto-generated aliases.
func PopulateRepoAliases(db *store.DB) *apperror.AppError {
	_, err := PopulateRepoAliasesWithCount(db)

	return err
}

// openDb opens the gitmap database from the binary's data directory.
func openDb() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

