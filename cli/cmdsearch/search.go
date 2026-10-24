package cmdsearch

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/pterm/pterm"

	searchconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/constants"
	searchenums "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunSearch executes indexed keyword and symbol search across repositories.
func RunSearch(args []string) error {
	checkHelp(searchconstants.CmdSearch, args)
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
	flagQuery, positional := parseSearchQueryAndDirFlags(cleanArgs)
	if flagQuery != "" {
		return flagQuery, limit, isAiCaller, true
	}
	if len(positional) == 0 {
		return "", limit, isAiCaller, false
	}

	return positional[0], limit, isAiCaller, true
}

func parseSearchQueryAndDirFlags(args []string) (string, []string) {
	var query string
	var positional []string
	for i := 0; i < len(args); i++ {
		q, p, nextI := matchSearchFlagOrPositional(args, i)
		if q != "" {
			query = q
		}
		if p != "" {
			positional = append(positional, p)
		}
		i = nextI
	}

	return query, positional
}

func matchSearchFlagOrPositional(args []string, i int) (string, string, int) {
	arg := args[i]
	if isSearchQueryArg(arg) && i+1 < len(args) {
		return args[i+1], "", i + 1
	}
	if strings.HasPrefix(arg, searchconstants.FlagArgQueryPrefix) {
		return strings.TrimPrefix(arg, searchconstants.FlagArgQueryPrefix), "", i
	}
	if strings.HasPrefix(arg, searchconstants.FlagArgPatternPrefix) {
		return strings.TrimPrefix(arg, searchconstants.FlagArgPatternPrefix), "", i
	}
	if isSearchDirArg(arg) && i+1 < len(args) {
		return "", "", i + 1
	}
	if !strings.HasPrefix(arg, searchconstants.FlagPrefixDash) {
		return "", arg, i
	}

	return "", "", i
}

func isSearchQueryArg(arg string) bool {
	return slices.Contains([]string{searchconstants.FlagArgQuery, searchconstants.FlagArgQueryShort, searchconstants.FlagArgPattern}, arg)
}

func isSearchDirArg(arg string) bool {
	return slices.Contains([]string{searchconstants.FlagArgRoot, searchconstants.FlagArgDir, searchconstants.FlagArgPath, searchconstants.FlagArgDirShort, searchconstants.FlagArgPathShort}, arg)
}

func parseSearchAiFlag(args []string) (bool, []string) {
	var isAiCaller bool
	var filtered []string
	for _, arg := range args {
		if slices.Contains([]string{searchconstants.FlagArgAI, searchconstants.FlagArgAIShort}, arg) {
			isAiCaller = true
			continue
		}
		filtered = append(filtered, arg)
	}

	return isAiCaller, filtered
}

func printSearchUsage() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println(searchconstants.UsageSearch)
	fmt.Println(searchconstants.UsageSearchClean)
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  " + searchconstants.DescFastIndexedSearch)
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println(searchconstants.ExampleSearch1)
	fmt.Println(searchconstants.ExampleSearch2)
	fmt.Println(searchconstants.ExampleSearch3)
	fmt.Println(searchconstants.UsageSearchClean)
}

func handleSearchCleanCommand() error {
	cleared, err := searcher.CleanSearchHotCache()
	if err != nil {
		return err
	}
	fmt.Printf(searchconstants.MsgSearchCleanSuccess,
		constants.ColorGreen, cleared, constants.ColorReset)

	return nil
}

func executeSearchCommand(query string, limit int, isAiCaller bool) error {
	action, _ := searchenums.ParseSearchAction(query)
	if action.IsHistoryAction() {
		return searcher.RenderAUMSearchHistoryTable(limit)
	}
	if action.IsCleanAction() {
		return handleSearchCleanCommand()
	}
	matchMode := searchenums.SearchMatchModeKeyword.String()
	if cachedRes, dh2d, isHot := searcher.LookupHotCachedSearch(query, matchMode); isHot {
		_, _ = searcher.RecordAUMSearchExecution(query, matchMode, isAiCaller, 0, cachedRes)
		fmt.Printf(searchconstants.FmtHotCacheMatch, dh2d, len(cachedRes))
		renderSearchResults(cachedRes)

		return nil
	}

	return executeSearchDatabaseQuery(query, limit, isAiCaller)
}

func executeSearchDatabaseQuery(query string, limit int, isAiCaller bool) error {
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
	matchMode := searchenums.SearchMatchModeKeyword.String()
	dh2d, _ := searcher.RecordAUMSearchExecution(query, matchMode, isAiCaller, durationMs, res)
	logSearchExecution(query, isAiCaller, durationMs, len(res))
	if searchErr != nil {
		pterm.Error.Println(searchErr)

		return nil
	}
	fmt.Printf(searchconstants.FmtSearchCompleted, dh2d, durationMs, len(res))
	renderSearchResults(res)

	return nil
}

func logSearchExecution(query string, isAiCaller bool, durationMs, count int) {
	category := searchenums.SearchCategoryUser.String()
	if isAiCaller {
		category = searchenums.SearchCategoryAI.String()
	}

	_ = searcher.LogSearchQuery(searcher.SearchLogEntry{
		CategoryCode: category,
		QueryText:    query,
		SearchType:   searchenums.SearchMatchModeKeyword.String(),
		IsAiCaller:   isAiCaller,
		DurationMs:   durationMs,
		ResultCount:  count,
	})
}

func renderSearchResults(res []searcher.SearchResult) {
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf(searchconstants.FmtFoundHeader, r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
}

func runUntrackedSearchFallback(query string, limit int, isAiCaller bool) error {
	cwd, _ := os.Getwd()
	fmt.Printf(searchconstants.FmtSearchingDir, query, cwd)
	matches := findMatchingFilesLocal(cwd, query, limit)
	if len(matches) > 0 {
		printFileMatches(matches)

		return nil
	}
	repoMatches := findMatchingReposInStore(query, limit)
	if len(repoMatches) > 0 {
		printRepoMatches(repoMatches)

		return nil
	}
	printUntrackedSearchHint(query)

	return nil
}

func printFileMatches(matches []string) {
	fmt.Printf(searchconstants.FmtMatchingFilesHeader, constants.ColorGreen, len(matches), constants.ColorReset)
	for _, m := range matches {
		fmt.Printf(searchconstants.FmtBulletItem, constants.ColorCyan, m, constants.ColorReset)
	}
}

func printRepoMatches(matches []string) {
	fmt.Printf(searchconstants.FmtMatchingReposHeader, constants.ColorGreen, len(matches), constants.ColorReset)
	for _, rm := range matches {
		fmt.Printf(searchconstants.FmtBulletItem, constants.ColorCyan, rm, constants.ColorReset)
	}
}

func printUntrackedSearchHint(query string) {
	fmt.Printf(searchconstants.FmtUntrackedNotFound, constants.ColorYellow, query, constants.ColorReset)
	fmt.Println(searchconstants.MsgUntrackedHintTip)
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
		matches = append(matches, fmt.Sprintf(searchconstants.FmtRepoWithLocation, r.RepoName, r.AbsolutePath))
		if len(matches) >= limit {
			break
		}
	}

	return matches
}

func isIgnoredSearchDir(name string) bool {
	return strings.HasPrefix(name, searchconstants.DirPrefixDot) ||
		slices.Contains([]string{searchconstants.DirNodeModules, searchconstants.DirVendor, searchconstants.DirBin}, name)
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
		limit = searchconstants.DefaultLocalSearchLimit
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
