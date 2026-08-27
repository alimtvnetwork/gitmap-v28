package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/repodb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/searcher"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/pterm/pterm"
)

func parseLimit(args []string) (int, []string) {
	limit := 0
	var cleanArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--limit" || args[i] == "-l" {
			if i+1 < len(args) {
				limit, _ = strconv.Atoi(args[i+1])
				i++
			}
			continue
		}
		if strings.HasPrefix(args[i], "--limit=") {
			limit, _ = strconv.Atoi(strings.TrimPrefix(args[i], "--limit="))
			continue
		}
		cleanArgs = append(cleanArgs, args[i])
	}
	return limit, cleanArgs
}

func getFindDB(ctx context.Context) (*store.DB, *sql.DB, error) {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return nil, nil, err
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		mainDB.Close()
		return nil, nil, err
	}
	
	repos, err := mainDB.FindByPath(cwd)
	if err != nil || len(repos) == 0 {
		mainDB.Close()
		return nil, nil, fmt.Errorf("current directory is not a tracked gitmap repository. run 'gitmap scan' first")
	}
	
	repoDB, err := repodb.OpenRepoDB(ctx, constants.DefaultOutputDir, repos[0].AbsolutePath, repos[0].ID)
	if err != nil {
		mainDB.Close()
		return nil, nil, err
	}
	
	return mainDB, repoDB, nil
}

func runFind(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find <query> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getFindDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindFile(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		fmt.Println(r.RelativePath)
	}
}

func runFindRegex(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find-regex <regex> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getFindDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindFileRegex(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		fmt.Println(r.RelativePath)
	}
}

func runFindRead(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find-read <query> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getFindDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, false, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Println(r.RelativePath)
		fmt.Println(r.Content)
		fmt.Println()
	}
}

func runFindReadJson(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("[]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getFindDB(ctx)
	if err != nil {
		fmt.Println("[]")
		return
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, false, limit, true)
	if err != nil {
		fmt.Println("[]")
		return
	}
	
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
}

func runFindRegexRead(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find-regex-read <regex> [--limit <n>]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getFindDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, true, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Println(r.RelativePath)
		fmt.Println(r.Content)
		fmt.Println()
	}
}

func runFindRegexReadJson(args []string) {
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("[]")
		return
	}
	query := cleanArgs[0]
	
	ctx := context.Background()
	mainDB, db, err := getFindDB(ctx)
	if err != nil {
		fmt.Println("[]")
		return
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, true, limit, true)
	if err != nil {
		fmt.Println("[]")
		return
	}
	
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
}

func runFindHelp(args []string) {
	pterm.DefaultSection.Println("Find Help Options")
	pterm.Println("  gitmap find <query> [--limit <n>]")
	pterm.Println("  gitmap find-regex <regex> [--limit <n>]")
	pterm.Println("  gitmap find-read <query> [--limit <n>]")
	pterm.Println("  gitmap find-read-json <query> [--limit <n>]")
}

func runSearchHelp(args []string) {
	pterm.DefaultSection.Println("Search Help Options")
	pterm.Println("  gitmap search <query>")
	pterm.Println("  gitmap search-replace-all <query>")
	pterm.Println("  gitmap repo-search <query>")
}

func runRegexHelp(args []string) {
	pterm.DefaultSection.Println("Regex Help Options")
	pterm.Println("  gitmap replace-regex <regex>")
	pterm.Println("  gitmap repo-regex <regex>")
	pterm.Println("  gitmap find-regex <regex>")
}
