package cmdclone

import (
	"encoding/json"
	"os/exec"
	"strings"

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
	hasCachedURL := len(cachedURL) > 0
	if hasCachedURL {
		return cachedURL
	}

	ghURL := resolveFromGitHubCLI(target)
	hasGhURL := len(ghURL) > 0
	if hasGhURL {
		return ghURL
	}

	return target
}

func isFullURLOrPath(target string) bool {
	trimmed := strings.TrimSpace(target)
	hasHTTP := strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
	hasSSH := strings.HasPrefix(trimmed, "git@") || strings.HasPrefix(trimmed, "ssh://")
	hasFile := strings.HasPrefix(trimmed, "file://")
	hasSlash := strings.Contains(trimmed, "/") || strings.Contains(trimmed, "\\")

	return hasHTTP || hasSSH || hasFile || hasSlash
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
		isMatchingName := strings.EqualFold(r.RepoName, cleanSlug) || strings.EqualFold(r.Slug, cleanSlug)
		if isMatchingName {
			if len(r.HTTPSUrl) > 0 {
				return r.HTTPSUrl
			}
			if len(r.SSHUrl) > 0 {
				return r.SSHUrl
			}
		}
	}

	return ""
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

	cleanSlug := strings.ToLower(slug)
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
