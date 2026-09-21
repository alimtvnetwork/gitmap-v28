package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
)

// runSearch executes indexed keyword and symbol search across repositories.
func runSearch(args []string) error {
	checkHelp("search", args)
	query, limit, isAiCaller, hasQuery := parseSearchInvocation(args)
	if !hasQuery {
		printSearchUsage()
		return nil
	}

	return executeSearchCommand(query, limit, isAiCaller)
}

func parseSearchInvocation(args []string) (string, int, bool, bool) {
	isAiCaller, remaining := parseSearchAiFlag(args)
	limit, cleanArgs := parseLimit(remaining)
	if len(cleanArgs) == 0 {
		return "", limit, isAiCaller, false
	}

	return cleanArgs[0], limit, isAiCaller, true
}

func parseSearchAiFlag(args []string) (bool, []string) {
	var isAiCaller bool
	var filtered []string
	for _, arg := range args {
		if arg == "--ai" || arg == "-ai" {
			isAiCaller = true
			continue
		}
		filtered = append(filtered, arg)
	}

	return isAiCaller, filtered
}

func printSearchUsage() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap search <query> [--limit <n>] [--ai]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Fast indexed keyword and symbol search across repositories using SplitDB.")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap search \"Resolve-Version\"")
	fmt.Println("  gitmap search \"AppError\" --limit 10")
	fmt.Println("  gitmap search \"type SearchResult struct\" --ai")
}

func executeSearchCommand(query string, limit int, isAiCaller bool) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "search.getRepoDB")
	}
	defer mainDB.Close()
	defer db.Close()

	return runSearchWithMetrics(ctx, db, query, limit, isAiCaller)
}

func runSearchWithMetrics(ctx context.Context, db *sql.DB, query string, limit int, isAiCaller bool) error {
	start := time.Now()
	res, searchErr := searcher.SearchRepoDB(ctx, db, query, limit, false)
	durationMs := int(time.Since(start).Milliseconds())
	logSearchExecution(query, isAiCaller, durationMs, len(res))
	if searchErr != nil {
		pterm.Error.Println(searchErr)
		return nil
	}

	renderSearchResults(res)
	return nil
}

func logSearchExecution(query string, isAiCaller bool, durationMs, count int) {
	category := "user"
	if isAiCaller {
		category = "ai"
	}

	_ = searcher.LogSearchQuery(searcher.SearchLogEntry{
		CategoryCode: category,
		QueryText:    query,
		SearchType:   "keyword",
		IsAiCaller:   isAiCaller,
		DurationMs:   durationMs,
		ResultCount:  count,
	})
}

func renderSearchResults(res []searcher.SearchResult) {
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf("Found in %s at position %d", r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
}
