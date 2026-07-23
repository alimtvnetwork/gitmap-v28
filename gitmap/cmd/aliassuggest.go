package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runAliasSuggest handles "alias suggest [--apply]".
func runAliasSuggest(args []string) {
	apply := parseAliasSuggestFlags(args)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	repos, err := db.ListUnaliasedRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		exitWith(1)
	}

	if len(repos) == 0 {
		fmt.Println(constants.MsgAliasSuggestNone)

		return
	}

	created := suggestAliases(db, repos, apply)
	fmt.Printf(constants.MsgAliasSuggestDone, created)
	printHints(aliasSuggestHints())
}

// parseAliasSuggestFlags parses flags for alias suggest.
func parseAliasSuggestFlags(args []string) bool {
	fs := flag.NewFlagSet(constants.SubCmdAliasSug, flag.ExitOnError)
	apply := fs.Bool("apply", false, constants.FlagDescAliasApply)
	_ = fs.Parse(args)

	return *apply
}

// suggestAliases proposes aliases for unaliased repos.
func suggestAliases(db *store.DB, repos []store.UnaliasedRepo, autoApply bool) int {
	created := 0
	reader := bufio.NewReader(os.Stdin)

	for _, r := range repos {
		suggestion := r.RepoName
		if db.AliasExists(suggestion) {
			continue
		}

		if autoApply {
			createSuggestedAlias(db, suggestion, r.ID)
			created++

			continue
		}

		if promptAliasSuggestion(reader, r.Slug, suggestion) {
			createSuggestedAlias(db, suggestion, r.ID)
			created++
		}
	}

	return created
}

// promptAliasSuggestion asks the user to accept a suggested alias.
func promptAliasSuggestion(reader *bufio.Reader, slug, suggestion string) bool {
	fmt.Printf(constants.MsgAliasSuggest, slug, suggestion)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// createSuggestedAlias creates an alias and prints confirmation.
func createSuggestedAlias(db *store.DB, alias string, repoID int64) {
	_, err := db.CreateAlias(alias, repoID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return
	}

	fmt.Printf(constants.MsgAliasCreated, alias, fmt.Sprintf("%d", repoID))
}
