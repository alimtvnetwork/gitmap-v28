package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/completion"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runCompletion handles the "completion" subcommand.
func runCompletion(args []string) {
	checkHelp("completion", args)

	if hasListFlag(args) {
		handleCompletionList(args)

		return
	}

	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrCompUsage)
		os.Exit(1)
	}

	printCompletionScript(args[0])
}

// hasListFlag checks if any --list-* flag is present.
func hasListFlag(args []string) bool {
	for _, a := range args {
		if a == constants.CompListRepos || a == constants.CompListGroups ||
			a == constants.CompListCommands || a == constants.CompListAliases ||
			a == constants.CompListZipGroups || a == constants.CompListSSHKeys {
			return true
		}
	}

	return false
}

// handleCompletionList routes to the appropriate list printer.
func handleCompletionList(args []string) {
	for _, a := range args {
		switch a {
		case constants.CompListRepos:
			printCompletionRepos()

			return
		case constants.CompListGroups:
			printCompletionGroups()

			return
		case constants.CompListCommands:
			printCompletionCommands()

			return
		case constants.CompListAliases:
			printCompletionAliases()

			return
		case constants.CompListZipGroups:
			printCompletionZipGroups()

			return
		case constants.CompListSSHKeys:
			printCompletionSSHKeys()

			return
		case constants.CompListHelpGroups:
			printCompletionHelpGroups()

			return
		}
	}
}

// printCompletionRepos prints all repo slugs, one per line.
func printCompletionRepos() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	repos, err := db.ListRepos()
	if err != nil {
		return
	}

	for _, r := range repos {
		fmt.Println(r.Slug)
	}
}

// printCompletionGroups prints all group names, one per line.
func printCompletionGroups() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	groups, err := db.ListGroups()
	if err != nil {
		return
	}

	for _, g := range groups {
		fmt.Println(g.Name)
	}
}

// printCompletionCommands prints all command names, one per line.
func printCompletionCommands() {
	for _, cmd := range completion.AllCommands() {
		fmt.Println(cmd)
	}
}

// printCompletionAliases prints all alias names, one per line.
func printCompletionAliases() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	aliases, err := db.ListAliases()
	if err != nil {
		return
	}

	for _, a := range aliases {
		fmt.Println(a.Alias)
	}
}

// printCompletionZipGroups prints all zip group names, one per line.
func printCompletionZipGroups() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	groups, err := db.ListZipGroups()
	if err != nil {
		return
	}

	for _, g := range groups {
		fmt.Println(g.Name)
	}
}

// printCompletionSSHKeys prints all SSH key names, one per line.
func printCompletionSSHKeys() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	names, err := db.SSHKeyNames()
	if err != nil {
		return
	}

	for _, n := range names {
		fmt.Println(n)
	}
}

// printCompletionHelpGroups prints all help group keywords, one per line.
func printCompletionHelpGroups() {
	for _, g := range constants.HelpGroupKeys {
		fmt.Println(g)
	}
}

// printCompletionScript outputs the shell completion script.
func printCompletionScript(shell string) {
	script, err := completion.Generate(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCompUnknownShell, shell)
		os.Exit(1)
	}

	fmt.Print(script)
}
