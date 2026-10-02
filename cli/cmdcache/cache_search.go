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
	rootDb, err := store.OpenRootCacheDB(repoRoot)
	hasErr := err != nil
	if hasErr {
		return err
	}
	defer rootDb.Close()

	allMatches := searchAllSlugs(rootDb, repoRoot, opts)
	renderSearchResults(allMatches, opts)
	return nil
}

func searchAllSlugs(rootDb *sql.DB, repoRoot string, opts CacheSearchOptions) []SearchMatch {
	slugs := collectIndexedSlugs(rootDb)
	var allMatches []SearchMatch
	for _, slug := range slugs {
		maxRem := opts.ResultLimit - len(allMatches)
		matches := searchSlug(slug, repoRoot, opts, maxRem)
		allMatches = append(allMatches, matches...)
		hasReachedLimit := len(allMatches) >= opts.ResultLimit
		if hasReachedLimit {
			break
		}
	}
	return allMatches
}

func collectIndexedSlugs(rootDb *sql.DB) []string {
	rows, err := rootDb.Query("SELECT DISTINCT FolderSlug FROM Files")
	hasErr := err != nil
	if hasErr {
		return nil
	}
	defer rows.Close()
	return scanSlugs(rows)
}

func scanSlugs(rows *sql.Rows) []string {
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
	slugDb, err := store.OpenSlugCacheDB(slug, repoRoot)
	hasErr := err != nil
	if hasErr {
		return nil
	}
	defer slugDb.Close()
	return executeSlugSearch(slugDb, opts, maxRemaining)
}

func executeSlugSearch(slugDb *sql.DB, opts CacheSearchOptions, maxRemaining int) []SearchMatch {
	rows, queryErr := slugDb.Query("SELECT RelativePath, LineNumber, Content FROM Lines ORDER BY RelativePath, LineNumber")
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil
	}
	defer rows.Close()

	regexes := prepareCompiledRegexes(opts.Patterns, opts.IsRegex)
	return scanAndMatchRows(slugDb, rows, opts, regexes, maxRemaining)
}

func prepareCompiledRegexes(patterns []string, isRegex bool) []*regexp.Regexp {
	if !isRegex {
		return nil
	}
	var regexes []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		hasErr := err != nil
		if !hasErr {
			regexes = append(regexes, re)
		}
	}
	return regexes
}

func scanAndMatchRows(slugDb *sql.DB, rows *sql.Rows, opts CacheSearchOptions, regexes []*regexp.Regexp, maxRemaining int) []SearchMatch {
	var matches []SearchMatch
	for rows.Next() && len(matches) < maxRemaining {
		m, hasMatch := scanMatchRow(slugDb, rows, opts, regexes)
		if hasMatch {
			matches = append(matches, m)
		}
	}
	return matches
}

func scanMatchRow(slugDb *sql.DB, rows *sql.Rows, opts CacheSearchOptions, regexes []*regexp.Regexp) (SearchMatch, bool) {
	var m SearchMatch
	scanErr := rows.Scan(&m.RelativePath, &m.LineNumber, &m.Content)
	hasScanErr := scanErr != nil
	if hasScanErr {
		return m, false
	}
	isEligible := isRowEligible(m, opts, regexes)
	if !isEligible {
		return m, false
	}
	m.ContextLines = attachContextLines(slugDb, m.RelativePath, m.LineNumber, opts.LinesToShow)
	return m, true
}

func attachContextLines(slugDb *sql.DB, relPath string, lineNum, linesToShow int) []store.CacheContextLine {
	if linesToShow <= 0 {
		return nil
	}
	return store.FetchContextLines(slugDb, relPath, lineNum, linesToShow)
}

func isRowEligible(m SearchMatch, opts CacheSearchOptions, regexes []*regexp.Regexp) bool {
	isMatchedFile := isFileMatch(m.RelativePath, opts.FileGlobs)
	if !isMatchedFile {
		return false
	}
	return isContentMatch(m.Content, opts.Patterns, regexes, opts.IsRegex)
}

func isFileMatch(relPath string, globs []string) bool {
	hasGlobs := len(globs) > 0
	if !hasGlobs {
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
		hasPat := cleanPat != "" && strings.Contains(lowerContent, cleanPat)
		if hasPat {
			return true
		}
	}
	return false
}

func renderSearchResults(matches []SearchMatch, opts CacheSearchOptions) {
	hasMatches := len(matches) > 0
	if !hasMatches {
		fmt.Printf("\n%sNo matching lines found in cache.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s  === CACHE SEARCH RESULTS (%d matches, context ±%d lines) ===%s\n",
		constants.ColorCyan, len(matches), opts.LinesToShow, constants.ColorReset)
	renderMatchList(matches)
	fmt.Println()
}

func renderMatchList(matches []SearchMatch) {
	for _, m := range matches {
		hasContext := len(m.ContextLines) > 0
		if hasContext {
			renderMatchWithContext(m)
			continue
		}
		renderSimpleMatch(m)
	}
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
