package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/searcher"
)

func parseLimit(args []string) (int, []string) {
	limit := 0
	var cleanArgs []string
	for i := 0; i < len(args); i++ {
		if (args[i] == "--limit" || args[i] == "-l") && i+1 < len(args) {
			limit, _ = strconv.Atoi(args[i+1])
			i++
			continue
		}
		if args[i] == "--limit" || args[i] == "-l" {
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

// getRepoDB is used from cmd_db.go

func runFind(args []string) error {
	checkHelp("find", args)
	return executeFindFiles(args, MatchWildcard)
}

func runFindRegex(args []string) error {
	checkHelp("find-regex", args)
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find-regex <regex> [--limit <n>]")
		return nil
	}
	query := cleanArgs[0]

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "error")
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindFileRegex(ctx, db, query, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	}

	for _, r := range res {
		fmt.Println(r.RelativePath)
	}
	return nil
}

func runFindRead(args []string) error {
	checkHelp("find-read", args)
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find-read <query> [--limit <n>]")
		return nil
	}
	query := cleanArgs[0]

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "error")
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, false, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	}

	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Println(r.RelativePath)
		fmt.Println(r.Content)
		fmt.Println()
	}
	return nil
}

func runFindReadJson(args []string) error {
	checkHelp("find-read-json", args)
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
		return nil
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, false, limit, true)
	if err != nil {
		fmt.Println("[]")
		return nil
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	return nil
}

func runFindRegexRead(args []string) error {
	checkHelp("find-regex-read", args)
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println("Usage: gitmap find-regex-read <regex> [--limit <n>]")
		return nil
	}
	query := cleanArgs[0]

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "error")
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, true, limit, true)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	}

	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Println(r.RelativePath)
		fmt.Println(r.Content)
		fmt.Println()
	}
	return nil
}

func runFindRegexReadJson(args []string) error {
	checkHelp("find-regex-read-json", args)
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
		return nil
	}
	defer mainDB.Close()
	defer db.Close()

	res, err := searcher.FindAndRead(ctx, db, query, true, limit, true)
	if err != nil {
		fmt.Println("[]")
		return nil
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	return nil
}

func runFindHelp(args []string) error {
	pterm.DefaultSection.Println("Find Help Options")
	pterm.Println("  gitmap find <query> [--limit <n>]")
	pterm.Println("  gitmap find-regex <regex> [--limit <n>]")
	pterm.Println("  gitmap find-read <query> [--limit <n>]")
	pterm.Println("  gitmap find-read-json <query> [--limit <n>]")
	return nil
}

func runSearchHelp(args []string) error {
	pterm.DefaultSection.Println("Search Help Options")
	pterm.Println("  gitmap search <query>")
	pterm.Println("  gitmap search-replace-all <query>")
	pterm.Println("  gitmap repo-search <query>")
	return nil
}

func runRegexHelp(args []string) error {
	pterm.DefaultSection.Println("Regex Help Options")
	pterm.Println("  gitmap replace-regex <regex>")
	pterm.Println("  gitmap repo-regex <regex>")
	pterm.Println("  gitmap find-regex <regex>")
	return nil
}
