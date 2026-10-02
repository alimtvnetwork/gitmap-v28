package cmdcache

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunCacheCLI routes gitmap cache commands.
func RunCacheCLI(args []string) *appfault.AppError {
	hasArgs := len(args) > 0
	if !hasArgs {
		return PrintCacheHelp()
	}
	subCmd := strings.ToLower(args[0])
	subArgs := args[1:]
	return dispatchSubCmd(subCmd, subArgs)
}

func dispatchSubCmd(subCmd string, subArgs []string) *appfault.AppError {
	switch subCmd {
	case "help", "-h", "--help":
		return PrintCacheHelp()
	case "create", "add":
		return CreateCache(subArgs)
	case "ls", "list":
		return runCacheList()
	case "remove", "rm":
		return runCacheRemove(subArgs)
	case "recache", "reconcile", "sync":
		return ReconcileCache(findRepoRoot())
	default:
		return routeSearchCmd(subCmd, subArgs)
	}
}

func routeSearchCmd(subCmd string, subArgs []string) *appfault.AppError {
	switch subCmd {
	case "search":
		return runSearchSingle(subArgs)
	case "search-multi":
		return runSearchMulti(subArgs, false)
	case "search-multi-grep":
		return runSearchMulti(subArgs, true)
	default:
		return appfault.NewSimple("unknown cache subcommand: "+subCmd, "E1030")
	}
}

func runCacheList() *appfault.AppError {
	repoRoot := findRepoRoot()
	summaries, _ := store.ListCachedRepos(repoRoot)
	hasSummaries := len(summaries) > 0
	if hasSummaries {
		renderRepoSummaries(summaries)
	}
	return listAndRenderCacheFiles(repoRoot)
}

func listAndRenderCacheFiles(repoRoot string) *appfault.AppError {
	rootDb, err := store.OpenRootCacheDB(repoRoot)
	hasErr := err != nil
	if hasErr {
		return err
	}
	defer rootDb.Close()

	files, queryErr := store.ListCacheFiles(rootDb)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return queryErr
	}
	renderCachedFiles(files)
	return nil
}

func renderRepoSummaries(summaries []store.CachedRepoSummary) {
	fmt.Printf("\n%s  === CACHED REPOSITORIES (%d total) ===%s\n",
		constants.ColorCyan, len(summaries), constants.ColorReset)
	for _, s := range summaries {
		fmt.Printf("  • %-20s | %4d files | %8d B | last: %s\n",
			s.RepoSlug, s.FileCount, s.TotalBytes, s.LastIndexed)
	}
}

func renderCachedFiles(files []store.CacheFileRecord) {
	fmt.Printf("\n%s  === CURRENT REPO CACHED FILES (%d total) ===%s\n",
		constants.ColorCyan, len(files), constants.ColorReset)
	for _, f := range files {
		keepTag := ""
		if f.IsKeep {
			keepTag = " [keep]"
		}
		fmt.Printf("  • %-45s | %6d B | slug: %-10s%s\n",
			f.RelativePath, f.FileSize, f.FolderSlug, keepTag)
	}
	fmt.Println()
}

func runCacheRemove(args []string) *appfault.AppError {
	hasArgs := len(args) > 0
	if !hasArgs {
		return appfault.NewSimple("Usage: gitmap cache remove <path>", "E1031")
	}
	repoRoot := findRepoRoot()
	targets := extractTargetParts(strings.Join(args, ","))
	queueId, tasksDb := enqueueCacheTask("remove", strings.Join(targets, ","))
	defer closeTasksDB(tasksDb)

	return executeCacheRemove(targets, repoRoot, queueId, tasksDb)
}

func executeCacheRemove(targets []string, repoRoot, queueId string, tasksDb *store.TasksSplitDB) *appfault.AppError {
	rootDb, err := store.OpenRootCacheDB(repoRoot)
	hasErr := err != nil
	if hasErr {
		logAndFailCacheTask(tasksDb, queueId, err, "cache_remove_open_db")
		return err
	}
	defer rootDb.Close()

	for _, target := range targets {
		removeSingleTarget(rootDb, target, repoRoot)
	}
	completeCacheTask(tasksDb, queueId)
	fmt.Printf("%s✓ Removed %d target(s) from cache database.%s\n", constants.ColorGreen, len(targets), constants.ColorReset)
	return nil
}

func removeSingleTarget(rootDb *sql.DB, target, repoRoot string) {
	cleanRel := filepath.ToSlash(target)
	slug := resolveFileSlug(cleanRel)
	_ = store.DeleteCacheFile(rootDb, cleanRel)
	removeLinesFromSlug(slug, cleanRel, repoRoot)
}

func runSearchSingle(args []string) *appfault.AppError {
	hasArgs := len(args) > 0
	if !hasArgs {
		return appfault.NewSimple("Usage: gitmap cache search <query> [glob] [--lines N] [--limit N]", "E1032")
	}
	opts := parseSearchFlags(args)
	return SearchCache(opts)
}

func runSearchMulti(args []string, isRegex bool) *appfault.AppError {
	hasArgs := len(args) > 0
	if !hasArgs {
		return appfault.NewSimple("Usage: gitmap cache search-multi <queries...> [-file-pattern glob]", "E1033")
	}
	opts := parseSearchFlags(args)
	opts.IsRegex = isRegex
	return SearchCache(opts)
}

// PrintCacheHelp outputs usage documentation for gitmap cache.
func PrintCacheHelp() *appfault.AppError {
	printHelpOverview()
	printHelpCommands()
	printHelpExamples()
	return nil
}

func printHelpOverview() {
	fmt.Printf(`
%s  GitMap Repository Split-DB Cache Engine (cache)%s
  High-speed indexed repository cache utilizing SQLite Split-DB architecture:
  • Root metadata & file index: .gitmap/cache/repos/<slug>/sql.db
  • Top-level folder databases: .gitmap/cache/repos/<slug>/<folder-slug>.db
  • Exclusion gates: files > 200KB, binaries, large JSON, .git/, node_modules/
  • Zero SHA hashes during scan (relies on filesystem mtime)
`, constants.ColorCyan, constants.ColorReset)
}

func printHelpCommands() {
	fmt.Printf(`  COMMANDS:
    • gitmap cache create [paths...]             Index target path(s) into Split-DB
    • gitmap cache add <paths...>                Add specific files or paths to cache
    • gitmap cache ls / list                     List all cached repositories and files
    • gitmap cache remove / rm <paths...>        Remove files or paths from cache
    • gitmap cache search "<text>" [glob]        Fast text search with context lines
    • gitmap cache search-multi "t1", "t2"       Multi-term search with file pattern
    • gitmap cache search-multi-grep "<regex>"   Regular expression search
    • gitmap cache recache / reconcile / sync    Update cache based on filesystem mtime
    • gitmap cache help                          Show this help documentation
`)
}

func printHelpExamples() {
	fmt.Printf(`  FLAGS & EXAMPLES:
    • gitmap cache create .
    • gitmap cache create src,pkg/api --keep
    • gitmap cache create "a.json", "b.json"
    • gitmap cache search "TODO" "*.go" --lines 10 --limit 20
    • gitmap cache search "ErrNotFound" -file-pattern (fp) "a*.md", "b*.md"
    • gitmap cache search-multi "func", "return" -fp "*.go" --lines 5
    • gitmap cache search-multi-grep "AppError.*Simple" -fp "*.go"
    • gitmap cache reconcile
`)
}
