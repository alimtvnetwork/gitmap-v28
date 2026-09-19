package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runAliasSuggest handles "alias suggest [--apply]".
func runAliasSuggest(args []string) error {
	apply := parseAliasSuggestFlags(args)

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	repos, err := db.ListUnaliasedRepos()
	if err != nil {
		return apperror.WrapSimple(err, "list unaliased repos")
	}

	if len(repos) == 0 {
		fmt.Println(constants.MsgAliasSuggestNone)

		return nil
	}

	created := suggestAliases(db, repos, apply)
	fmt.Printf(constants.MsgAliasSuggestDone, created)
	printHints(aliasSuggestHints())

	return nil
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
		if processSingleAliasSuggestion(db, reader, r, autoApply) {
			created++
		}
	}

	return created
}

func processSingleAliasSuggestion(db *store.DB, reader *bufio.Reader, r store.UnaliasedRepo, autoApply bool) bool {
	suggestion := determineRepoSuggestion(db, r)
	if suggestion == "" {
		return false
	}
	if autoApply || promptAliasSuggestion(reader, r.Slug, suggestion) {
		createSuggestedAlias(db, suggestion, r.ID)

		return true
	}

	return false
}

// determineRepoSuggestion calculates a unique suggestion for a repository.
func determineRepoSuggestion(db *store.DB, r store.UnaliasedRepo) string {
	candidate := GenerateAutoAlias(r.RepoName)
	if candidate == "" {
		candidate = r.RepoName
	}

	return ResolveUniqueAlias(candidate, db.AliasExists)
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
	_, err := db.CreateAliasWithDetails(alias, repoID, true, "auto")
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return
	}

	fmt.Printf(constants.MsgAliasCreated, alias, fmt.Sprintf("%d", repoID))
}
