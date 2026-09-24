package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	lowQ := strings.ToLower(strings.TrimSpace(query))
	if lowQ == "history" || lowQ == "top" || lowQ == "stats" || lowQ == "dh2d" {
		return searcher.RenderAUMSearchHistoryTable(limit)
	}
	if cachedRes, dh2d, isHot := searcher.LookupHotCachedSearch(query, "keyword"); isHot {
		_, _ = searcher.RecordAUMSearchExecution(query, "keyword", isAiCaller, 0, cachedRes)
		fmt.Printf("⚡ [AUM Hot-Cache %s] Instant match (<0.04ms) — %d result(s)\n", dh2d, len(cachedRes))
		renderSearchResults(cachedRes)
		return nil
	}
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return runUntrackedSearchFallback(query, limit, isAiCaller)
	}
	defer mainDB.Close()
	defer db.Close()

	return runSearchWithMetrics(ctx, db, query, limit, isAiCaller)
}

func runSearchWithMetrics(ctx context.Context, db *sql.DB, query string, limit int, isAiCaller bool) error {
	start := time.Now()
	res, searchErr := searcher.SearchRepoDB(ctx, db, query, limit, false)
	durationMs := int(time.Since(start).Milliseconds())
	dh2d, _ := searcher.RecordAUMSearchExecution(query, "keyword", isAiCaller, durationMs, res)
	logSearchExecution(query, isAiCaller, durationMs, len(res))
	if searchErr != nil {
		pterm.Error.Println(searchErr)
		return nil
	}
	fmt.Printf("🔍 [AUM Search ID: %s] Completed in %d ms (%d matches)\n", dh2d, durationMs, len(res))
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

func runUntrackedSearchFallback(query string, limit int, isAiCaller bool) error {
	cwd, _ := os.Getwd()
	fmt.Printf("Searching for %q in %s...\n", query, cwd)
	matches := findMatchingFilesLocal(cwd, query, limit)
	hasFileMatches := len(matches) > 0
	if hasFileMatches {
		printFileMatches(matches)
		return nil
	}
	repoMatches := findMatchingReposInStore(query, limit)
	hasRepoMatches := len(repoMatches) > 0
	if hasRepoMatches {
		printRepoMatches(repoMatches)
		return nil
	}
	printUntrackedSearchHint(query)
	return nil
}

func printFileMatches(matches []string) {
	fmt.Printf("\n%sFound %d matching file(s):%s\n", constants.ColorGreen, len(matches), constants.ColorReset)
	for _, m := range matches {
		fmt.Printf("  • %s%s%s\n", constants.ColorCyan, m, constants.ColorReset)
	}
}

func printRepoMatches(matches []string) {
	fmt.Printf("\n%sFound %d matching repository(ies):%s\n", constants.ColorGreen, len(matches), constants.ColorReset)
	for _, rm := range matches {
		fmt.Printf("  • %s%s%s\n", constants.ColorCyan, rm, constants.ColorReset)
	}
}

func printUntrackedSearchHint(query string) {
	fmt.Printf("\n%sNo matches found for %q in untracked folder.%s\n", constants.ColorYellow, query, constants.ColorReset)
	fmt.Println("  Tip: run 'gitmap scan' to index repositories or cd into a tracked repository.")
}

func findMatchingReposInStore(query string, limit int) []string {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer mainDB.Close()
	repos, _ := mainDB.ListRepos()
	var matches []string
	cleanQ := strings.ToLower(query)
	for _, r := range repos {
		isRepoMatch := strings.Contains(strings.ToLower(r.RepoName), cleanQ)
		isPathMatch := strings.Contains(strings.ToLower(r.AbsolutePath), cleanQ)
		if !isRepoMatch && !isPathMatch {
			continue
		}
		matches = append(matches, r.RepoName+" ("+r.AbsolutePath+")")
		if len(matches) >= limit {
			break
		}
	}
	return matches
}

func isIgnoredSearchDir(name string) bool {
	return strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "bin"
}

func formatRelOrAbsPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err == nil {
		return rel
	}
	return path
}

func findMatchingFilesLocal(root, query string, limit int) []string {
	if limit <= 0 {
		limit = 20
	}
	var matches []string
	cleanQ := strings.ToLower(query)
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && isIgnoredSearchDir(d.Name()) {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if !strings.Contains(strings.ToLower(d.Name()), cleanQ) {
			return nil
		}
		matches = append(matches, formatRelOrAbsPath(root, path))
		if len(matches) >= limit {
			return filepath.SkipAll
		}
		return nil
	})
	return matches
}
