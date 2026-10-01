package cmdcache

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunCacheCLI routes gitmap cache commands.
func RunCacheCLI(args []string) error {
	if len(args) == 0 {
		return PrintCacheHelp()
	}
	subCmd := strings.ToLower(args[0])
	subArgs := args[1:]

	switch subCmd {
	case "help", "-h", "--help":
		return PrintCacheHelp()
	case "create":
		return CreateCache(subArgs)
	case "add":
		return CreateCache(subArgs)
	case "ls", "list":
		return runCacheList()
	case "remove", "rm":
		return runCacheRemove(subArgs)
	case "search":
		return runSearchSingle(subArgs)
	case "search-multi":
		return runSearchMulti(subArgs, false)
	case "search-multi-grep":
		return runSearchMulti(subArgs, true)
	case "recache", "reconcile", "sync":
		return CreateCache([]string{"."})
	default:
		return apperror.NewSimple("unknown cache subcommand: "+subCmd, "E1030")
	}
}

func runCacheList() error {
	repoRoot := findRepoRoot()
	rootDB, err := OpenRootCacheDB(repoRoot)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	files, queryErr := ListCacheFiles(rootDB)
	if queryErr != nil {
		return queryErr
	}

	fmt.Printf("\n%s  === CACHED FILES IN SPLIT-DB (%d total) ===%s\n",
		constants.ColorCyan, len(files), constants.ColorReset)
	for _, f := range files {
		fmt.Printf("  • %-45s | %6d B | slug: %-10s\n", f.RelativePath, f.FileSize, f.FolderSlug)
	}
	fmt.Println()
	return nil
}

func runCacheRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap cache remove <path>", "E1031")
	}
	repoRoot := findRepoRoot()
	rootDB, err := OpenRootCacheDB(repoRoot)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	for _, target := range args {
		if _, delErr := rootDB.Exec("DELETE FROM Files WHERE RelativePath = ?", target); delErr != nil {
			return apperror.WrapSimple(delErr, "cache_remove")
		}
	}
	fmt.Printf("%s✓ Removed %d file(s) from cache database.%s\n", constants.ColorGreen, len(args), constants.ColorReset)
	return nil
}

func runSearchSingle(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap cache search <query> [glob] [--lines N] [--limit N]", "E1032")
	}
	opts := parseSearchFlags(args)
	return SearchCache(opts)
}

func runSearchMulti(args []string, isRegex bool) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap cache search-multi <queries...> [-file-pattern glob]", "E1033")
	}
	opts := parseSearchFlags(args)
	opts.IsRegex = isRegex
	return SearchCache(opts)
}

func parseSearchFlags(args []string) CacheSearchOptions {
	opts := CacheSearchOptions{ResultLimit: 50}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--lines" && i+1 < len(args) {
			opts.LinesToShow, _ = strconv.Atoi(args[i+1])
			i++
		} else if a == "--limit" && i+1 < len(args) {
			opts.ResultLimit, _ = strconv.Atoi(args[i+1])
			i++
		} else if (a == "-file-pattern" || a == "-fp" || a == "--file-pattern") && i+1 < len(args) {
			opts.FileGlobs = append(opts.FileGlobs, strings.Split(args[i+1], ",")...)
			i++
		} else if strings.HasPrefix(a, "-") {
			continue
		} else if len(opts.Patterns) == 0 {
			opts.Patterns = append(opts.Patterns, splitQuotedTerms(a)...)
		} else if len(opts.FileGlobs) == 0 && (strings.Contains(a, "*") || strings.Contains(a, ".")) {
			opts.FileGlobs = append(opts.FileGlobs, a)
		} else {
			opts.Patterns = append(opts.Patterns, a)
		}
	}
	return opts
}

func splitQuotedTerms(raw string) []string {
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.Trim(strings.TrimSpace(p), `"'`)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// PrintCacheHelp outputs usage documentation for gitmap cache.
func PrintCacheHelp() error {
	fmt.Printf(`
  GitMap Repository Split-DB Cache Engine (cache)
    • gitmap cache create [paths...]             - Index target paths (files <= 200KB) into Split-DB
    • gitmap cache ls                            - List all indexed files in cache
    • gitmap cache add <paths...>                - Add specific files to cache
    • gitmap cache remove (rm) <paths...>        - Remove files from cache
    • gitmap cache search "query" [pattern]      - Full-text search across cached files
    • gitmap cache search-multi "q1, q2"         - Multi-query search across cached files
    • gitmap cache search-multi-grep "r1, r2"    - Multi-regex search across cached files
    • gitmap cache recache / reconcile / sync    - Re-index current repository into cache
`)
	return nil
}
