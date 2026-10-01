package cmdcache

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SearchCache searches through split database tables.
func SearchCache(opts CacheSearchOptions) *appfault.AppError {
	repoRoot := findRepoRoot()
	rootDB, err := store.OpenRootCacheDB(repoRoot)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	slugs := collectIndexedSlugs(rootDB)
	var allMatches []SearchMatch
	limit := resolveSearchLimit(opts.ResultLimit)

	for _, slug := range slugs {
		matches := searchSlug(slug, repoRoot, opts, limit-len(allMatches))
		allMatches = append(allMatches, matches...)
		if len(allMatches) >= limit {
			break
		}
	}

	renderSearchResults(allMatches, opts)
	return nil
}

func collectIndexedSlugs(rootDB *sql.DB) []string {
	rows, err := rootDB.Query("SELECT DISTINCT FolderSlug FROM Files")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var s string
		if rows.Scan(&s) == nil {
			slugs = append(slugs, s)
		}
	}
	return slugs
}

func searchSlug(slug, repoRoot string, opts CacheSearchOptions, maxRemaining int) []SearchMatch {
	slugDB, err := store.OpenSlugCacheDB(slug, repoRoot)
	if err != nil {
		return nil
	}
	defer slugDB.Close()

	query := "SELECT RelativePath, LineNumber, Content FROM Lines ORDER BY RelativePath, LineNumber"
	rows, queryErr := slugDB.Query(query)
	if queryErr != nil {
		return nil
	}
	defer rows.Close()

	return filterRows(rows, opts, maxRemaining)
}

func filterRows(rows *sql.Rows, opts CacheSearchOptions, maxRemaining int) []SearchMatch {
	var matches []SearchMatch
	for rows.Next() && len(matches) < maxRemaining {
		m, hasMatch := scanAndMatchRow(rows, opts)
		if hasMatch {
			matches = append(matches, m)
		}
	}
	return matches
}

func scanAndMatchRow(rows *sql.Rows, opts CacheSearchOptions) (SearchMatch, bool) {
	var m SearchMatch
	if scanErr := rows.Scan(&m.RelativePath, &m.LineNumber, &m.Content); scanErr != nil {
		return m, false
	}
	if !isFileMatch(m.RelativePath, opts.FileGlobs) || !isContentMatch(m.Content, opts) {
		return m, false
	}
	return m, true
}

func isFileMatch(relPath string, globs []string) bool {
	if len(globs) == 0 {
		return true
	}
	base := filepath.Base(relPath)
	for _, g := range globs {
		clean := strings.TrimSpace(g)
		isMatched, _ := filepath.Match(clean, base)
		if isMatched {
			return true
		}
	}
	return false
}

func isContentMatch(content string, opts CacheSearchOptions) bool {
	for _, pat := range opts.Patterns {
		if hasPatternMatch(content, pat, opts.IsRegex) {
			return true
		}
	}
	return false
}

func hasPatternMatch(content, pat string, isRegex bool) bool {
	if isRegex {
		re, err := regexp.Compile(pat)
		return err == nil && re.MatchString(content)
	}
	return strings.Contains(strings.ToLower(content), strings.ToLower(pat))
}

func resolveSearchLimit(userLimit int) int {
	if userLimit > 0 {
		return userLimit
	}
	return 50
}

func renderSearchResults(matches []SearchMatch, opts CacheSearchOptions) {
	if len(matches) == 0 {
		fmt.Printf("\n%sNo matching lines found in cache.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s  === CACHE SEARCH RESULTS (%d matches) ===%s\n",
		constants.ColorCyan, len(matches), constants.ColorReset)
	for _, m := range matches {
		fmt.Printf("  %s%s:%d%s  %s\n",
			constants.ColorCyan, m.RelativePath, m.LineNumber, constants.ColorReset, strings.TrimSpace(m.Content))
	}
	fmt.Println()
}
