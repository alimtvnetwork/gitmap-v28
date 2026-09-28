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

	printRepoSlugSuggestions(target)

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
		if strings.EqualFold(r.RepoName, cleanSlug) || strings.EqualFold(r.Slug, cleanSlug) {
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

	suggs, _ := mainDB.GetRepoSuggestions(target)
	if len(suggs) == 0 {
		suggs = findClosestRepoSuggestions(mainDB, target)
	}
	if len(suggs) > 0 {
		fmt.Fprintf(os.Stderr, "Repository %q not found. Did you mean:\n", target)
		for _, s := range suggs {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
	}
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
