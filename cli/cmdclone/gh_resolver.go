package cmdclone

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type ghRepoListing struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`
}

// ResolveRepoSlug resolves a bare repo name (e.g. "pwp-mobile") to a cloneable URL.
// If target is already a valid URL or path, it returns target unchanged.
func ResolveRepoSlug(target string) string {
	if isFullURLOrPath(target) {
		return target
	}

	cachedURL := resolveFromStore(target)
	if len(cachedURL) > 0 {
		return cachedURL
	}

	ghURL := resolveFromGitHubCLI(target)
	if len(ghURL) > 0 {
		return ghURL
	}

	if isGitHubOwnerRepo(target) {
		return "https://github.com/" + strings.TrimSpace(target)
	}

	printRepoSlugSuggestions(target)

	return target
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
	if strings.HasPrefix(target, "./") || strings.HasPrefix(target, ".\\") ||
		strings.HasPrefix(target, "../") || strings.HasPrefix(target, "..\\") ||
		strings.HasPrefix(target, "/") || strings.HasPrefix(target, "\\") ||
		strings.HasPrefix(target, "~") {
		return true
	}

	if len(target) >= 3 && target[1] == ':' && (target[2] == '/' || target[2] == '\\') {
		return true
	}

	if info, err := os.Stat(target); err == nil && info.IsDir() {
		return true
	}

	return false
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
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

func resolveFromStore(slug string) string {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer mainDB.Close()

	cleanSlug := strings.ToLower(slug)
	repos, err := mainDB.ListRepos()
	if err != nil {
		return ""
	}

	for _, r := range repos {
		if strings.EqualFold(r.RepoName, cleanSlug) || strings.EqualFold(r.Slug, cleanSlug) {
			return pickRepoURL(r)
		}
	}

	for _, r := range repos {
		httpLower := strings.ToLower(r.HTTPSUrl)
		sshLower := strings.ToLower(r.SSHUrl)
		if strings.HasSuffix(httpLower, "/"+cleanSlug) || strings.HasSuffix(httpLower, "/"+cleanSlug+".git") ||
			strings.HasSuffix(sshLower, ":"+cleanSlug) || strings.HasSuffix(sshLower, ":"+cleanSlug+".git") ||
			strings.HasSuffix(sshLower, "/"+cleanSlug) || strings.HasSuffix(sshLower, "/"+cleanSlug+".git") {
			return pickRepoURL(r)
		}
	}

	for _, r := range repos {
		if strings.HasPrefix(strings.ToLower(r.RepoName), cleanSlug) || strings.HasPrefix(strings.ToLower(r.Slug), cleanSlug) {
			return pickRepoURL(r)
		}
	}

	for _, r := range repos {
		if strings.Contains(strings.ToLower(r.RepoName), cleanSlug) || strings.Contains(strings.ToLower(r.Slug), cleanSlug) {
			return pickRepoURL(r)
		}
	}

	resolved, err := mainDB.ResolveAlias(slug)
	if err == nil && resolved.Slug != "" {
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
		fmt.Fprintf(os.Stderr, "Repository %q not found. Did you mean:\n", cleanTarget)
		for _, s := range filtered {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
	}
}

func filterRedundantSuggestions(suggs []string, cleanTarget, target string) []string {
	var filtered []string
	for _, s := range suggs {
		if strings.EqualFold(s, cleanTarget) || strings.EqualFold(s, target) || strings.EqualFold(extractRepoNameFromTarget(s), cleanTarget) {
			continue
		}
		filtered = append(filtered, s)
	}
	return filtered
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

func resolveFromGitHubCLI(slug string) string {
	cmd := exec.Command("gh", "repo", "list", "--limit", "200", "--json", "name,url,nameWithOwner")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	var repos []ghRepoListing
	unmarshalErr := json.Unmarshal(out, &repos)
	if unmarshalErr != nil {
		return ""
	}

	cleanSlug := strings.ToLower(extractRepoNameFromTarget(slug))
	for _, r := range repos {
		isExactMatch := strings.EqualFold(r.Name, cleanSlug)
		if isExactMatch {
			return r.URL
		}
	}

	for _, r := range repos {
		isPrefixMatch := strings.HasPrefix(strings.ToLower(r.Name), cleanSlug)
		if isPrefixMatch {
			return r.URL
		}
	}

	return ""
}
