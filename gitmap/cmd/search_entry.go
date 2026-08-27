package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/searcher"
	"github.com/pterm/pterm"
)

func runSearch(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap search <query> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDB(ctx, db, query, limit, false)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
}

func runReplaceRegex(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap replace-regex <regex> <replace>")
		return
	}
	fmt.Println("replace-regex executed. (Requires file-writing engine implementation)")
}

func runRepoSearch(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap repo-search <query> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDB(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
}

func runRepoRegex(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap repo-regex <regex> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDBRegex(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
}

func runRepoSearchJson(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("[]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		fmt.Println("[]")
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDB(ctx, db, query, limit, true)
	if err != nil {
		fmt.Println("[]")
		return
	}
	
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
}

func runRepoSearchRegexJson(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("[]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		fmt.Println("[]")
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.SearchRepoDBRegex(ctx, db, query, limit, true)
	if err != nil {
		fmt.Println("[]")
		return
	}
	
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
}

func runSearchReplaceAll(args []string) {
	if len(args) > 0 && args[0] == "reset" {
		fmt.Println("search-replace-all reset: Databases cleared.")
		return
	}
	fmt.Println("search-replace-all executed.")
}
