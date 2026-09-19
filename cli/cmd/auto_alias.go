package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var knownCompoundPrefixes = []string{
	"git", "web", "doc", "app", "cli", "dev", "code", "data",
	"auto", "work", "pipe", "test", "cloud", "task", "file", "net",
	"front", "back", "full", "mono", "micro", "quick", "desk", "sync",
	"view", "node", "host", "bit", "repo", "tool", "pack", "term",
}

var knownCompoundSuffixes = []string{
	"map", "hub", "lab", "top", "run", "flow", "line", "base",
	"stack", "end", "ops", "fix", "site", "pass", "auth", "store",
	"port", "box", "path", "desk", "queue", "link", "list", "view",
	"book", "board", "craft", "tree", "bucket", "runner", "server", "serve",
}

func init() {
	cmdinstall.GenerateAutoAliasFn = GenerateAutoAlias
	cmdinstall.PopulateRepoAliasesFn = func(db *store.DB) (int, error) {
		count, appErr := PopulateRepoAliasesWithCount(db)
		if appErr != nil {
			return 0, appErr
		}

		return count, nil
	}
}

// GenerateAutoAlias creates a concise, deterministic alias for a repository name.
func GenerateAutoAlias(raw string) string {
	name := normalizeAliasTarget(raw)
	if name == "" {
		return ""
	}
	if isPresentationRepo(raw, name) {
		return formatPresentationAlias(raw, name)
	}
	if isWpPackage(name) {
		return formatWpAlias(name)
	}

	return resolveGenericAlias(name)
}

func resolveGenericAlias(name string) string {
	if alias := generateMultiWordAcronym(name); alias != "" {
		return alias
	}

	return generateCompoundAcronym(name)
}

func normalizeAliasTarget(raw string) string {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimSuffix(clean, ".git")
	clean = filepath.Base(filepath.ToSlash(clean))

	return strings.ToLower(clean)
}

func isPresentationRepo(raw, name string) bool {
	rawLower := strings.ToLower(filepath.ToSlash(raw))
	if strings.Contains(rawLower, "presentation") || strings.Contains(rawLower, "presentations") {
		return true
	}

	return strings.HasPrefix(name, "prep-") || strings.HasPrefix(name, "prep_")
}

func formatPresentationAlias(raw, name string) string {
	tokens := extractPresentationTokens(raw, name)
	for _, tok := range tokens {
		if isPresentationIdentifier(tok) {
			return "prep-" + tok
		}
	}

	return "prep"
}

func extractPresentationTokens(raw, name string) []string {
	combined := name
	if strings.Contains(raw, "/") || strings.Contains(raw, "\\") {
		combined = filepath.ToSlash(raw)
	}

	return splitWordTokens(combined)
}

func isPresentationIdentifier(tok string) bool {
	if tok == "presentation" || tok == "presentations" || tok == "prep" {
		return false
	}
	if tok == "repo" || tok == "repos" || tok == "repository" || tok == "repositories" {
		return false
	}

	return !isVersionSuffix(tok) && len(tok) >= 2
}

func isWpPackage(name string) bool {
	if strings.HasPrefix(name, "wp-") || strings.HasPrefix(name, "wp_") {
		return true
	}

	return strings.HasPrefix(name, "wp ")
}

func formatWpAlias(name string) string {
	remainder := stripWpPrefix(name)
	if remainder == "" {
		return "wp"
	}
	if alias := generateMultiWordAcronym(remainder); alias != "" {
		return "wp-" + alias
	}
	if alias := generateCompoundAcronym(remainder); alias != "" {
		return "wp-" + alias
	}

	return "wp-" + remainder
}

func stripWpPrefix(name string) string {
	for _, prefix := range []string{"wp-", "wp_", "wp "} {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimSpace(name[len(prefix):])
		}
	}

	return name
}

func generateMultiWordAcronym(name string) string {
	words := filterMeaningfulWords(splitWordTokens(name))
	if len(words) <= 1 {
		return ""
	}
	var sb strings.Builder
	for _, w := range words {
		sb.WriteByte(w[0])
	}

	return sb.String()
}

func splitWordTokens(name string) []string {
	clean := strings.Map(replaceDelimiterWithSpace, name)
	rawWords := strings.Fields(clean)
	var words []string
	for _, w := range rawWords {
		if len(w) > 0 {
			words = append(words, w)
		}
	}

	return words
}

func replaceDelimiterWithSpace(r rune) rune {
	if r == '-' || r == '_' || r == '/' || r == '\\' || r == '.' {
		return ' '
	}

	return r
}

func filterMeaningfulWords(words []string) []string {
	var filtered []string
	for _, w := range words {
		if !isVersionSuffix(w) && len(w) > 0 {
			filtered = append(filtered, w)
		}
	}

	return filtered
}

func isVersionSuffix(tok string) bool {
	if len(tok) < 2 || (tok[0] != 'v' && tok[0] != 'V') {
		return false
	}
	for i := 1; i < len(tok); i++ {
		if tok[i] < '0' || tok[i] > '9' {
			return false
		}
	}

	return true
}

func generateCompoundAcronym(name string) string {
	if camelAcronym := splitCamelCaseAcronym(name); camelAcronym != "" {
		return camelAcronym
	}

	return findCompoundAcronym(name)
}

func splitCamelCaseAcronym(name string) string {
	var initials []byte
	for i := 0; i < len(name); i++ {
		if name[i] >= 'A' && name[i] <= 'Z' {
			initials = append(initials, name[i]+32)
		}
	}
	if len(initials) >= 2 {
		return string(initials)
	}

	return ""
}

func findCompoundAcronym(name string) string {
	for _, p := range knownCompoundPrefixes {
		if strings.HasPrefix(name, p) && len(name) > len(p) {
			suffix := name[len(p):]
			if isKnownCompoundSuffix(suffix) {
				return string([]byte{p[0], suffix[0]})
			}
		}
	}

	return ""
}

func isKnownCompoundSuffix(suffix string) bool {
	for _, s := range knownCompoundSuffixes {
		if suffix == s || strings.HasPrefix(suffix, s) {
			return true
		}
	}

	return len(suffix) >= 3
}

// ResolveUniqueAlias resolves collision by appending a numeric suffix if needed.
func ResolveUniqueAlias(baseAlias string, isTaken func(string) bool) string {
	if baseAlias == "" {
		return ""
	}
	if !isTaken(baseAlias) {
		return baseAlias
	}
	for i := 2; i <= 99; i++ {
		candidate := fmt.Sprintf("%s%d", baseAlias, i)
		if !isTaken(candidate) {
			return candidate
		}
	}

	return ""
}

// PopulateRepoAliases scans unaliased repos and populates auto-generated aliases.
func PopulateRepoAliases(db *store.DB) *apperror.AppError {
	_, err := PopulateRepoAliasesWithCount(db)

	return err
}

// PopulateRepoAliasesWithCount populates aliases and returns the count of added aliases.
func PopulateRepoAliasesWithCount(db *store.DB) (int, *apperror.AppError) {
	repos, err := db.ListUnaliasedRepos()
	if err != nil {
		return 0, apperror.WrapSimple(err, "list unaliased repos")
	}

	return countPopulatedAliases(db, repos), nil
}

func countPopulatedAliases(db *store.DB, repos []store.UnaliasedRepo) int {
	created := 0
	for _, r := range repos {
		if populateSingleRepoAlias(db, r) {
			created++
		}
	}

	return created
}

func populateSingleRepoAlias(db *store.DB, r store.UnaliasedRepo) bool {
	alias := GenerateAutoAlias(r.RepoName)
	if alias == "" {
		return false
	}
	unique := ResolveUniqueAlias(alias, db.AliasExists)
	if unique == "" {
		return false
	}
	_, err := db.CreateAliasWithDetails(unique, r.ID, true, "auto")

	return err == nil
}
