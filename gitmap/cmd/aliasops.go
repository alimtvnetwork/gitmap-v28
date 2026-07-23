package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runAliasSet handles "alias set <alias> <slug>".
func runAliasSet(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
		os.Exit(1)
	}

	alias := args[0]
	slug := args[1]

	executeAliasSet(alias, slug)
}

// executeAliasSet resolves the slug and creates or updates the alias.
func executeAliasSet(alias, slug string) {
	if code := executeAliasSetCode(alias, slug); code != 0 {
		os.Exit(code)
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

	if db.AliasExists(alias) {
		if err := db.UpdateAlias(alias, repoID); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
			return 1
		}
		fmt.Printf(constants.MsgAliasUpdated, alias, slug)
		printHints(aliasSetHints())
		return 0
	}

	if _, err := db.CreateAlias(alias, repoID); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		return 1
	}
	fmt.Printf(constants.MsgAliasCreated, alias, slug)
	printHints(aliasSetHints())
	return 0
}

// runAliasRemove handles "alias remove <alias>".
func runAliasRemove(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
		os.Exit(1)
	}

	alias := args[0]
	if code := runAliasRemoveCode(alias); code != 0 {
		os.Exit(code)
	}
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
func runAliasList() {
	if code := runAliasListCode(); code != 0 {
		os.Exit(code)
	}
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
func runAliasShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
		os.Exit(1)
	}

	alias := args[0]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	resolved, err := db.ResolveAlias(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		exitWith(1)
	}

	fmt.Printf(constants.MsgAliasResolved, resolved.Alias, resolved.AbsolutePath, resolved.Slug)
}

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}
