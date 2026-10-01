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
	opts = normalizeSearchOpts(opts)
	repoRoot := findRepoRoot()
	rootDB, err := store.OpenRootCacheDB(repoRoot)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	slugs := collectIndexedSlugs(rootDB)
	var allMatches []SearchMatch

	for _, slug := range slugs {
		maxRem := opts.ResultLimit - len(allMatches)
		matches := searchSlug(slug, repoRoot, opts, maxRem)
		allMatches = append(allMatches, matches...)
		if len(allMatches) >= opts.ResultLimit {
			break
		}
	}
	renderSearchResults(allMatches, opts)
	return nil
}

func normalizeSearchOpts(opts CacheSearchOptions) CacheSearchOptions {
	if opts.LinesToShow <= 0 {
		opts.LinesToShow = 10
	}
	if opts.ResultLimit <= 0 {
		opts.ResultLimit = 20
	}
	return opts
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

	rows, queryErr := slugDB.Query("SELECT RelativePath, LineNumber, Content FROM Lines ORDER BY RelativePath, LineNumber")
	if queryErr != nil {
		return nil
	}
	defer rows.Close()

	regexes := prepareCompiledRegexes(opts.Patterns, opts.IsRegex)
	return scanAndMatchRows(slugDB, rows, opts, regexes, maxRemaining)
}

func prepareCompiledRegexes(patterns []string, isRegex bool) []*regexp.Regexp {
	if !isRegex {
		return nil
	}
	var regexes []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err == nil {
			regexes = append(regexes, re)
		}
	}
	return regexes
}

func scanAndMatchRows(slugDB *sql.DB, rows *sql.Rows, opts CacheSearchOptions, regexes []*regexp.Regexp, maxRemaining int) []SearchMatch {
	var matches []SearchMatch
	for rows.Next() && len(matches) < maxRemaining {
		var m SearchMatch
		if scanErr := rows.Scan(&m.RelativePath, &m.LineNumber, &m.Content); scanErr != nil {
			continue
		}
		if !isRowEligible(m, opts, regexes) {
			continue
		}
		if opts.LinesToShow > 0 {
			m.ContextLines = store.FetchContextLines(slugDB, m.RelativePath, m.LineNumber, opts.LinesToShow)
		}
		matches = append(matches, m)
	}
	return matches
}

func isRowEligible(m SearchMatch, opts CacheSearchOptions, regexes []*regexp.Regexp) bool {
	if !isFileMatch(m.RelativePath, opts.FileGlobs) {
		return false
	}
	return isContentMatch(m.Content, opts.Patterns, regexes, opts.IsRegex)
}

func isFileMatch(relPath string, globs []string) bool {
	if len(globs) == 0 {
		return true
	}
	for _, g := range globs {
		clean := strings.TrimSpace(g)
		if isGlobMatch(relPath, clean) {
			return true
		}
	}
	return false
}

func isGlobMatch(relPath, pattern string) bool {
	base := filepath.Base(relPath)
	isMatched, _ := filepath.Match(pattern, base)
	if isMatched {
		return true
	}
	pathMatched, _ := filepath.Match(pattern, relPath)
	return pathMatched
}

func isContentMatch(content string, patterns []string, regexes []*regexp.Regexp, isRegex bool) bool {
	if isRegex {
		return hasRegexMatch(content, regexes)
	}
	return hasTextMatch(content, patterns)
}

func hasRegexMatch(content string, regexes []*regexp.Regexp) bool {
	for _, re := range regexes {
		if re.MatchString(content) {
			return true
		}
	}
	return false
}

func hasTextMatch(content string, patterns []string) bool {
	lowerContent := strings.ToLower(content)
	for _, pat := range patterns {
		cleanPat := strings.ToLower(strings.TrimSpace(pat))
		if cleanPat != "" && strings.Contains(lowerContent, cleanPat) {
			return true
		}
	}
	return false
}

func renderSearchResults(matches []SearchMatch, opts CacheSearchOptions) {
	if len(matches) == 0 {
		fmt.Printf("\n%sNo matching lines found in cache.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s  === CACHE SEARCH RESULTS (%d matches, context ±%d lines) ===%s\n",
		constants.ColorCyan, len(matches), opts.LinesToShow, constants.ColorReset)
	for _, m := range matches {
		if len(m.ContextLines) > 0 {
			renderMatchWithContext(m)
		} else {
			renderSimpleMatch(m)
		}
	}
	fmt.Println()
}

func renderMatchWithContext(m SearchMatch) {
	fmt.Printf("\n  %s[%s:%d]%s\n", constants.ColorCyan, m.RelativePath, m.LineNumber, constants.ColorReset)
	for _, cl := range m.ContextLines {
		if cl.IsMatch {
			fmt.Printf("  %s> %5d | %s%s\n", constants.ColorGreen, cl.LineNumber, cl.Content, constants.ColorReset)
		} else {
			fmt.Printf("    %5d | %s\n", cl.LineNumber, cl.Content)
		}
	}
}

func renderSimpleMatch(m SearchMatch) {
	fmt.Printf("  %s%s:%d%s  %s\n",
		constants.ColorCyan, m.RelativePath, m.LineNumber, constants.ColorReset, strings.TrimSpace(m.Content))
}
