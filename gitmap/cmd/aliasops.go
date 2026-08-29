package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runAliasSet handles "alias set <alias> <slug>".
func runAliasSet(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple(constants.ErrAliasEmpty, "E9000")
	}

	alias := args[0]
	slug := args[1]

	executeAliasSet(alias, slug)
	return nil
}

// executeAliasSet resolves the slug and creates or updates the alias.
func executeAliasSet(alias, slug string) {
	if code := executeAliasSetCode(alias, slug); code != 0 {
		cliexit.HandleError(nil, code)
	}
}

// executeAliasSetCode performs the work and returns an exit code so
// deferred db.Close runs before any process exit.
func executeAliasSetCode(alias, slug string) int {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		return 1
	}
	defer db.Close()

	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrAliasRepoMissing, slug)
		return 1
	}

	repoID := repos[0].ID

	if db.AliasExists(alias) == true {
		return updateAliasAndReturn(db, alias, repoID, slug)
	}

	return createAliasAndReturnCode(db, alias, repoID, slug)
}

func updateAliasAndReturn(db *store.DB, alias string, repoID int64, slug string) int {
	if err := db.UpdateAlias(alias, repoID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		return 1
	}
	fmt.Printf(constants.MsgAliasUpdated, alias, slug)
	printHints(aliasSetHints())
	return 0
}

func createAliasAndReturnCode(db *store.DB, alias string, repoID int64, slug string) int {
	if _, err := db.CreateAlias(alias, repoID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		return 1
	}
	fmt.Printf(constants.MsgAliasCreated, alias, slug)
	printHints(aliasSetHints())
	return 0
}

// runAliasRemove handles "alias remove <alias>".
func runAliasRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrAliasEmpty, "E9000")
	}

	alias := args[0]
	if code := runAliasRemoveCode(alias); code != 0 {
		cliexit.HandleError(nil, code)
	}
	return nil
}

// runAliasRemoveCode returns an exit code so deferred db.Close runs
// before any process exit.
func runAliasRemoveCode(alias string) int {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		return 1
	}
	defer db.Close()

	if err := db.DeleteAlias(alias); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		return 1
	}
	fmt.Printf(constants.MsgAliasRemoved, alias)
	return 0
}

// runAliasList handles "alias list".
func runAliasList() error {
	if code := runAliasListCode(); code != 0 {
		cliexit.HandleError(nil, code)
	}
	return nil
}

// runAliasListCode returns an exit code so deferred db.Close runs
// before any process exit.
func runAliasListCode() int {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		return 1
	}
	defer db.Close()

	aliases, err := db.ListAliasesWithRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		return 1
	}
	printAliasList(aliases)
	printHints(aliasListHints())
	return 0
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

	alias := args[0]

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}
	defer db.Close()

	resolved, err := db.ResolveAlias(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		exitWith(1)
	}

	fmt.Printf(constants.MsgAliasResolved, resolved.Alias, resolved.AbsolutePath, resolved.Slug)
	return nil
}

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}
