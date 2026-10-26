package cmdclone

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/cli/diag"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type ghRepoListing struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`
}

// ResolveRepoSlugForProfile resolves a repo slug specifically considering a target Git profile.
func ResolveRepoSlugForProfile(target, profileName string) string {
	if isFullURLOrPath(target) {
		return target
	}

	cleanTarget := strings.TrimSpace(target)
	hasSlash := strings.Contains(cleanTarget, "/")
	if len(profileName) > 0 && !hasSlash {
		return fmt.Sprintf("https://github.com/%s/%s", profileName, cleanTarget)
	}

	return ResolveRepoSlug(cleanTarget)
}

// ResolveRepoSlug resolves a bare repo name (e.g. "pwp-mobile") to a cloneable URL.
// If target is already a valid URL or path, it returns target unchanged.
func ResolveRepoSlug(target string) string {
	if isFullURLOrPath(target) {
		return target
	}

	cleanTarget := strings.TrimSpace(target)
	if cachedURL := resolveFromStore(cleanTarget); len(cachedURL) > 0 {
		return cachedURL
	}

	if ghURL := resolveFromGitHubCLI(cleanTarget); len(ghURL) > 0 {
		return ghURL
	}

	if isGitHubOwnerRepo(cleanTarget) {
		return "https://github.com/" + cleanTarget
	}

	return cleanTarget
}

func isFullURLOrPath(target string) bool {
	trimmed := strings.TrimSpace(target)
	hasHTTP := strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
	hasSSH := strings.HasPrefix(trimmed, "git@") || strings.HasPrefix(trimmed, "ssh://")
	hasFile := strings.HasPrefix(trimmed, "file://")
	if hasHTTP || hasSSH || hasFile {
		return true
	}

	return isLocalPath(trimmed)
}

func isLocalPath(target string) bool {
	if isRelativeOrHomePath(target) || isDriveLetterPath(target) {
		return true
	}

	info, err := os.Stat(target)

	return err == nil && info.IsDir()
}

func isRelativeOrHomePath(target string) bool {
	return strings.HasPrefix(target, "./") || strings.HasPrefix(target, ".\\") ||
		strings.HasPrefix(target, "../") || strings.HasPrefix(target, "..\\") ||
		strings.HasPrefix(target, "/") || strings.HasPrefix(target, "\\") ||
		strings.HasPrefix(target, "~")
}

func isDriveLetterPath(target string) bool {
	return len(target) >= 3 && target[1] == ':' && (target[2] == '/' || target[2] == '\\')
}

func isGitHubOwnerRepo(target string) bool {
	trimmed := strings.TrimSuffix(strings.TrimSpace(target), ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		return false
	}

	owner, repo := parts[0], parts[1]
	if len(owner) == 0 || len(repo) == 0 {
		return false
	}

	return isValidSlugPart(owner) && isValidSlugPart(repo)
}

func isValidSlugPart(s string) bool {
	for _, r := range s {
		if !isSlugChar(r) {
			return false
		}
	}

	return true
}

func isSlugChar(r rune) bool {
	isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	isDigit := r >= '0' && r <= '9'
	isSpecial := r == '-' || r == '_' || r == '.'

	return isLetter || isDigit || isSpecial
}

func buildSlugCandidates(slug string) []string {
	clean := strings.TrimSpace(slug)
	if slashIdx := strings.Index(clean, "/"); slashIdx > 0 {
		return buildScopedSlugCandidates(clean[:slashIdx], clean[slashIdx+1:])
	}

	return buildBareSlugCandidates(clean)
}

func buildScopedSlugCandidates(owner, repo string) []string {
	parsed := clonenext.ParseRepoName(repo)
	if parsed.HasVersion {
		return []string{owner + "/" + repo, owner + "/" + parsed.BaseName}
	}

	return []string{owner + "/" + repo, owner + "/" + parsed.BaseName + "-v1", owner + "/" + parsed.BaseName}
}

func buildBareSlugCandidates(repo string) []string {
	parsed := clonenext.ParseRepoName(repo)
	if parsed.HasVersion {
		return []string{repo, parsed.BaseName}
	}

	return []string{repo, parsed.BaseName + "-v1", parsed.BaseName}
}

func resolveFromStore(slug string) string {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return ""
	}

	defer mainDB.Close()

	return resolveFromStoreWithDB(mainDB, slug)
}

func resolveFromStoreWithDB(db *store.DB, slug string) string {
	cleanSlug := strings.ToLower(slug)
	if url := findURLInRepoTable(db, cleanSlug); len(url) > 0 {
		return url
	}

	if url := findURLInOwnerRepoIndex(db, cleanSlug); len(url) > 0 {
		return url
	}

	if isAlimtvnetworkUser() || hasGitHubCredentials() {
		return findURLForAlimtvnetwork(db, cleanSlug)
	}

	return ""
}

func findURLInRepoTable(db *store.DB, cleanSlug string) string {
	repos, err := db.ListRepos()
	if err != nil {
		return ""
	}

	candidates := buildSlugCandidates(cleanSlug)
	for _, c := range candidates {
		if url := matchReposExactOrSuffix(repos, c); len(url) > 0 {
			return url
		}
	}

	if url := matchReposFuzzy(repos, cleanSlug); len(url) > 0 {
		return url
	}

	return matchAliasURL(db, repos, cleanSlug)
}

func matchReposExactOrSuffix(repos []model.ScanRecord, slug string) string {
	for _, r := range repos {
		if strings.EqualFold(r.RepoName, slug) || strings.EqualFold(r.Slug, slug) {
			return pickRepoURL(r)
		}
	}

	for _, r := range repos {
		if isRepoURLSuffixMatch(r, slug) {
			return pickRepoURL(r)
		}
	}

	return matchRepoByBaseName(repos, slug)
}

func matchRepoByBaseName(repos []model.ScanRecord, slug string) string {
	base := clonenext.ParseRepoName(slug).BaseName
	for _, r := range repos {
		if strings.EqualFold(clonenext.ParseRepoName(r.RepoName).BaseName, base) {
			return pickRepoURL(r)
		}
	}

	return ""
}

func isRepoURLSuffixMatch(r model.ScanRecord, slug string) bool {
	httpLower := strings.ToLower(r.HTTPSUrl)
	sshLower := strings.ToLower(r.SSHUrl)
	hasHTTPSuffix := strings.HasSuffix(httpLower, "/"+slug) || strings.HasSuffix(httpLower, "/"+slug+".git")
	hasSSHSuffix := strings.HasSuffix(sshLower, ":"+slug) || strings.HasSuffix(sshLower, ":"+slug+".git") ||
		strings.HasSuffix(sshLower, "/"+slug) || strings.HasSuffix(sshLower, "/"+slug+".git")

	return hasHTTPSuffix || hasSSHSuffix
}

func matchReposFuzzy(repos []model.ScanRecord, slug string) string {
	for _, r := range repos {
		if strings.HasPrefix(strings.ToLower(r.RepoName), slug) || strings.HasPrefix(strings.ToLower(r.Slug), slug) {
			return pickRepoURL(r)
		}
	}

	for _, r := range repos {
		if strings.Contains(strings.ToLower(r.RepoName), slug) || strings.Contains(strings.ToLower(r.Slug), slug) {
			return pickRepoURL(r)
		}
	}

	return ""
}

func matchAliasURL(db *store.DB, repos []model.ScanRecord, slug string) string {
	resolved, err := db.ResolveAlias(slug)
	if err == nil && len(resolved.Slug) > 0 {
		return matchRepoURLBySlug(repos, resolved.Slug)
	}

	return ""
}

func matchRepoURLBySlug(repos []model.ScanRecord, slug string) string {
	for _, r := range repos {
		if strings.EqualFold(r.Slug, slug) {
			return pickRepoURL(r)
		}
	}

	return ""
}

func pickRepoURL(r model.ScanRecord) string {
	if len(r.HTTPSUrl) > 0 {
		return r.HTTPSUrl
	}

	return r.SSHUrl
}

func findURLInOwnerRepoIndex(db *store.DB, cleanSlug string) string {
	if db == nil || db.Conn() == nil {
		return ""
	}

	cleanTarget := extractRepoNameFromTarget(cleanSlug)
	candidates := buildBareSlugCandidates(cleanTarget)
	for _, c := range candidates {
		if url := queryOwnerIndexRow(db, c); len(url) > 0 {
			return url
		}
	}

	return ""
}

func queryOwnerIndexRow(db *store.DB, candidate string) string {
	const query = `SELECT Owner, RepoName FROM OwnerRepoNameIndex 
		WHERE RepoName = ? OR BaseName = ? 
		ORDER BY VersionNumber DESC LIMIT 1`

	var owner, repoName string
	err := db.Conn().QueryRow(query, candidate, candidate).Scan(&owner, &repoName)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("https://github.com/%s/%s", owner, repoName)
}

func findURLForAlimtvnetwork(db *store.DB, cleanSlug string) string {
	cleanTarget := extractRepoNameFromTarget(cleanSlug)
	candidates := buildBareSlugCandidates(cleanTarget)
	for _, c := range candidates {
		if url := matchAlimCandidateWithCache(db, c); len(url) > 0 {
			return url
		}
	}

	return ""
}

func matchAlimCandidateWithCache(db *store.DB, c string) string {
	scoped := "alimtvnetwork/" + c
	if url := findURLInRepoTable(db, scoped); len(url) > 0 {
		return url
	}

	if url := queryOwnerIndexForOwner(db, "alimtvnetwork", c); len(url) > 0 {
		return url
	}

	return probeAndCacheRemoteCandidate(db, "alimtvnetwork", c)
}

func probeAndCacheRemoteCandidate(db *store.DB, owner, candidate string) string {
	url := fmt.Sprintf("https://github.com/%s/%s", owner, candidate)
	if !probeRemoteGitURL(url) {
		return ""
	}

	cacheDiscoveredRepo(db, owner, candidate)

	return url
}

func probeRemoteGitURL(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--heads", url)
	err := cmd.Run()

	return err == nil
}

func cacheDiscoveredRepo(db *store.DB, owner, repo string) {
	if db == nil {
		return
	}

	_ = db.EnsureOwnerRepoNameIndex()
	_ = db.UpsertOwnerRepoNameIndex("github", owner, []string{repo}, time.Now())
}

func queryOwnerIndexForOwner(db *store.DB, owner, candidate string) string {
	const query = `SELECT Owner, RepoName FROM OwnerRepoNameIndex 
		WHERE Owner = ? AND (RepoName = ? OR BaseName = ?) 
		ORDER BY VersionNumber DESC LIMIT 1`

	var rowOwner, repoName string
	err := db.Conn().QueryRow(query, owner, candidate, candidate).Scan(&rowOwner, &repoName)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("https://github.com/%s/%s", rowOwner, repoName)
}

func isAlimtvnetworkUser() bool {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return false
	}

	if strings.EqualFold(cfg.Active, "alimtvnetwork") || strings.EqualFold(cfg.Default, "alimtvnetwork") {
		return true
	}

	return hasProfileNamed(cfg.Profiles, "alimtvnetwork")
}

func hasProfileNamed(profiles []model.GitProfile, name string) bool {
	for _, p := range profiles {
		if strings.EqualFold(p.Name, name) {
			return true
		}
	}

	return false
}

func hasGitHubCredentials() bool {
	if _, hasToken := GetGlobalAccessToken(); hasToken {
		return true
	}

	_, _, err := secrets.Resolve()

	return err == nil
}

func resolveFromGitHubCLI(slug string) string {
	repos, err := fetchGitHubCLIRepos()
	if err != nil || len(repos) == 0 {
		return ""
	}

	candidates := buildBareSlugCandidates(extractRepoNameFromTarget(slug))
	for _, c := range candidates {
		if url := matchGitHubCLIRepo(repos, c); len(url) > 0 {
			return url
		}
	}

	return ""
}

func fetchGitHubCLIRepos() ([]ghRepoListing, error) {
	cmd := exec.Command("gh", "repo", "list", "--limit", "200", "--json", "name,url,nameWithOwner")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var repos []ghRepoListing
	if err := json.Unmarshal(out, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func matchGitHubCLIRepo(repos []ghRepoListing, candidate string) string {
	clean := strings.ToLower(candidate)
	for _, r := range repos {
		if strings.EqualFold(r.Name, clean) {
			return r.URL
		}
	}

	for _, r := range repos {
		if strings.HasPrefix(strings.ToLower(r.Name), clean) {
			return r.URL
		}
	}

	return ""
}

func printRepoSlugSuggestions(target string) {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return
	}

	defer mainDB.Close()

	cleanTarget := extractRepoNameFromTarget(target)
	suggs, _ := mainDB.GetRepoSuggestions(cleanTarget)
	if len(suggs) == 0 {
		suggs = findClosestRepoSuggestions(mainDB, cleanTarget)
	}

	filtered := filterRedundantSuggestions(suggs, cleanTarget, target)
	if len(filtered) > 0 {
		renderCatppuccinRepoSuggestions(cleanTarget, filtered)
	}
}

func renderCatppuccinRepoSuggestions(target string, suggestions []string) {
	group := diag.SuggestionGroup{
		Title:       "Repository Suggestions",
		Reason:      formatSuggestionReason(target),
		InputToken:  target,
		Suggestions: buildDiagSuggestions(suggestions),
	}
	if err := diag.RenderBox(os.Stderr, group); err != nil {
		printPlainRepoSuggestions(target, suggestions)
	}
}

func formatSuggestionReason(target string) string {
	if len(target) == 0 || target == "available" {
		return "Available repositories in local catalog:"
	}

	return fmt.Sprintf("Repository %q not found. Did you mean:", target)
}

func buildDiagSuggestions(suggestions []string) []diag.Suggestion {
	items := make([]diag.Suggestion, 0, len(suggestions))
	for _, s := range suggestions {
		items = append(items, diag.Suggestion{
			Command:     fmt.Sprintf("gitmap clone %s", s),
			Description: fmt.Sprintf("Clone %s", s),
			Confidence:  0.8,
		})
	}

	return items
}

func printPlainRepoSuggestions(target string, suggestions []string) {
	if len(target) > 0 && target != "available" {
		fmt.Fprintf(os.Stderr, "Repository %q not found. Did you mean:\n", target)
	} else {
		fmt.Fprintln(os.Stderr, "Available repositories:")
	}

	for _, s := range suggestions {
		fmt.Fprintf(os.Stderr, "  %s\n", s)
	}
}

func filterRedundantSuggestions(suggs []string, cleanTarget, target string) []string {
	var filtered []string
	for _, s := range suggs {
		if isRedundantSuggestion(s, cleanTarget, target) {
			continue
		}

		filtered = append(filtered, s)
	}

	return filtered
}

func isRedundantSuggestion(s, cleanTarget, target string) bool {
	return strings.EqualFold(s, cleanTarget) ||
		strings.EqualFold(s, target) ||
		strings.EqualFold(extractRepoNameFromTarget(s), cleanTarget)
}

func extractRepoNameFromTarget(target string) string {
	clean := strings.TrimSpace(target)
	clean = strings.TrimSuffix(clean, ".git")
	clean = strings.TrimRight(clean, "/\\")
	if idx := strings.LastIndexAny(clean, "/:"); idx != -1 {
		clean = clean[idx+1:]
	}

	return clean
}
