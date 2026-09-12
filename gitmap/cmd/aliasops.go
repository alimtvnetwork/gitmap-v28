package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runAliasSet handles "alias set <alias> <slug>".
func runAliasSet(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple(constants.ErrAliasEmpty, "E9000")
	}

	return executeAliasSet(args[0], args[1])
}

// executeAliasSet resolves the slug and creates or updates the alias.
func executeAliasSet(alias, slug string) error {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrListDBFailed, err))
	}
	defer db.Close()

	return resolveAndPersistAlias(db, alias, slug)
}

func resolveAndPersistAlias(db *store.DB, alias, slug string) error {
	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrAliasRepoMissing, slug)

		return apperror.NewSimple(fmt.Sprintf(constants.ErrAliasRepoMissing, slug), "E9000")
	}

	return persistAliasMapping(db, alias, repos[0].ID, slug)
}

func persistAliasMapping(db *store.DB, alias string, repoID int64, slug string) error {
	hasAlias := db.AliasExists(alias)
	if hasAlias {
		return updateAliasAndReturn(db, alias, repoID, slug)
	}

	return createAliasAndReturn(db, alias, repoID, slug)
}

func updateAliasAndReturn(db *store.DB, alias string, repoID int64, slug string) error {
	if err := db.UpdateAlias(alias, repoID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	fmt.Printf(constants.MsgAliasUpdated, alias, slug)
	printHints(aliasSetHints())

	return nil
}

func createAliasAndReturn(db *store.DB, alias string, repoID int64, slug string) error {
	if _, err := db.CreateAlias(alias, repoID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	fmt.Printf(constants.MsgAliasCreated, alias, slug)
	printHints(aliasSetHints())

	return nil
}

// runAliasRemove handles "alias remove <alias>".
func runAliasRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrAliasEmpty, "E9000")
	}

	return executeAliasRemove(args[0])
}

func executeAliasRemove(alias string) error {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrListDBFailed, err))
	}
	defer db.Close()

	if err := db.DeleteAlias(alias); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	fmt.Printf(constants.MsgAliasRemoved, alias)

	return nil
}

// runAliasList handles "alias list".
func runAliasList() error {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrListDBFailed, err))
	}
	defer db.Close()

	aliases, err := db.ListAliasesWithRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	printAliasList(aliases)
	printHints(aliasListHints())

	return nil
}

// printAliasList renders the alias table to stdout.
func printAliasList(aliases []store.AliasWithRepo) {
	if len(aliases) == 0 {
		fmt.Println("  No aliases defined.")

		return
	}

	fmt.Printf(constants.MsgAliasListHeader, len(aliases))

	for _, a := range aliases {
		fmt.Printf(constants.MsgAliasListRow, a.Alias.Alias, a.Slug)
	}
}

// runAliasShow handles "alias show <alias>".
func runAliasShow(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrAliasEmpty, "E9000")
	}

	return executeAliasShow(args[0])
}

func executeAliasShow(alias string) error {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}
	defer db.Close()

	resolved, err := db.ResolveAlias(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBareFmt, err))
	}

	fmt.Printf(constants.MsgAliasResolved, resolved.Alias.Alias, resolved.AbsolutePath, resolved.Slug)

	return nil
}

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}
