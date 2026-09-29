package cmdpipeline

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// QueryRepoSuggestionsFromDB queries repository suggestions matching partial string.
func QueryRepoSuggestionsFromDB(identifier string) []string {
	sdb, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer sdb.Close()
	conn := sdb.Conn()

	if candidates := fetchCandidateSlugs(conn, identifier); len(candidates) > 0 {
		return candidates
	}
	return fetchFuzzyCandidates(conn, identifier)
}

func fetchCandidateSlugs(conn *sql.DB, identifier string) []string {
	pattern := "%" + identifier + "%"
	query := "SELECT Slug FROM Repo WHERE Slug LIKE ? OR RepoName LIKE ? ORDER BY UpdatedAt DESC LIMIT 5"
	rows, err := conn.Query(query, pattern, pattern)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var suggs []string
	seen := make(map[string]bool)
	for rows.Next() {
		var s string
		if rows.Scan(&s) == nil && s != "" && !seen[s] {
			seen[s] = true
			suggs = append(suggs, s)
		}
	}
	return suggs
}

func fetchFuzzyCandidates(conn *sql.DB, identifier string) []string {
	rows, err := conn.Query("SELECT Slug, RepoName FROM Repo")
	if err != nil {
		return nil
	}
	defer rows.Close()

	clean := strings.ToLower(strings.TrimSpace(identifier))
	var suggs []string
	seen := make(map[string]bool)
	for rows.Next() {
		var slug, repoName string
		if rows.Scan(&slug, &repoName) != nil || seen[slug] {
			continue
		}
		if !isRepoFuzzyMatch(clean, strings.ToLower(repoName), strings.ToLower(slug)) {
			continue
		}
		seen[slug] = true
		suggs = append(suggs, slug)
		if len(suggs) >= 5 {
			break
		}
	}
	return suggs
}

func isRepoFuzzyMatch(clean, name, slug string) bool {
	base := strings.Split(name, "-v")[0]
	if strings.Contains(name, clean) || (len(base) > 2 && strings.Contains(clean, base)) {
		return true
	}
	return levenshteinDistance(clean, name) <= 4 || levenshteinDistance(clean, slug) <= 4 || levenshteinDistance(clean, base) <= 3
}

// PrintRepoTargetNotFoundDiagnostic prints formatted suggestions when target repo does not exist.
func PrintRepoTargetNotFoundDiagnostic(target string, suggestions []string) {
	fmt.Printf("\n  %s✖ Repository not found:%s %q\n\n", constants.ColorRed, constants.ColorReset, target)
	if len(suggestions) > 0 {
		fmt.Printf("  %sDid you mean:%s\n", constants.ColorCyan, constants.ColorReset)
		for _, s := range suggestions {
			fmt.Printf("    • %sgitmap pe %s%s\n", constants.ColorYellow, s, constants.ColorReset)
		}
		fmt.Println()
		return
	}
	fmt.Printf("  %sUse 'gitmap list' or 'gitmap scan' to discover indexed repositories.%s\n\n",
		constants.ColorDim, constants.ColorReset)
}
