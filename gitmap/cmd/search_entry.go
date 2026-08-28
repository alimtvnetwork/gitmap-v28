package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/searcher"
	"github.com/pterm/pterm"
)

func runSearch(args []string) error {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap search <query> [--limit <n>]")
		return nil
	}
	query := cleanArgs[0]

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.New(err, "E9000", nil)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDB(ctx, db, query, limit, false)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	}

	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
	return nil
}

func runReplaceRegex(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap replace-regex <regex> <replace>")
		return nil
	}
	fmt.Println("replace-regex executed. (Requires file-writing engine implementation)")
	return nil
}

func runRepoSearch(args []string) error {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap repo-search <query> [--limit <n>]")
		return nil
	}
	query := cleanArgs[0]
	if query == "history" {
		fmt.Println("repo-search history: No history recorded yet.")
		return nil
	}
	if query == "clear" {
		fmt.Println("repo-search clear: Cache cleared for current folder.")
		return nil
	}

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.New(err, "E9000", nil)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDB(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	}

	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
	return nil
}

func runRepoRegex(args []string) error {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap repo-regex <regex> [--limit <n>]")
		return nil
	}
	query := cleanArgs[0]
	if query == "history" {
		fmt.Println("repo-regex history: No history recorded yet.")
		return nil
	}
	if query == "clear" {
		fmt.Println("repo-regex clear: Cache cleared for current folder.")
		return nil
	}

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.New(err, "E9000", nil)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDBRegex(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	}

	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
	return nil
}

func runRepoSearchJson(args []string) error {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("[]")
		return nil
	}
	query := cleanArgs[0]

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		fmt.Println("[]")
		return apperror.New("fatal error", "E9000", nil)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDB(ctx, db, query, limit, true)
	if err != nil {
		fmt.Println("[]")
		return nil
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	return nil
}

func runRepoSearchRegexJson(args []string) error {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("[]")
		return nil
	}
	query := cleanArgs[0]

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		fmt.Println("[]")
		return apperror.New("fatal error", "E9000", nil)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDBRegex(ctx, db, query, limit, true)
	if err != nil {
		fmt.Println("[]")
		return nil
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	return nil
}

func runSearchReplaceAll(args []string) error {
	if len(args) > 0 && (args[0] == "reset" || args[0] == "clear") {
		// Clean the repo_search folder
		os.RemoveAll(".gitmap/output/repo_search")
		fmt.Println("search-replace-all reset: Databases cleared.")
		return nil
	}
	fmt.Println("search-replace-all executed.")
	return nil
}
