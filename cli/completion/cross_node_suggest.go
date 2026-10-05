package completion

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CrossNodeCandidate represents a repository suggestion target on a local or remote node.
type CrossNodeCandidate struct {
	NodeAlias, RepoSlug, FullPath string
	IsRemote                      bool
}

var (
	LocalCandidateSupplier      func() []string
	RemoteCandidateSupplier     func() []CrossNodeCandidate
	CollectCrossNodeSuggestions = GetCrossNodeRepoSuggestions
)

// GetCrossNodeRepoSuggestions returns local and remote node repository suggestions matching prefix.
func GetCrossNodeRepoSuggestions(prefix string) []string {
	slugs := queryLocalSlugs()
	cands := queryRemoteCandidates(slugs)
	items := append([]string{}, slugs...)
	for _, rc := range cands {
		items = append(items, formatCandidate(rc))
	}
	return filterAndDedupe(items, prefix)
}

func queryLocalSlugs() []string {
	if LocalCandidateSupplier != nil {
		return LocalCandidateSupplier()
	}
	if DynamicRepoSupplier != nil {
		return DynamicRepoSupplier()
	}
	return queryFallbackSlugs()
}

func queryFallbackSlugs() []string {
	sdb, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer sdb.Close()
	repos, _ := sdb.ListRepos()
	out := make([]string, len(repos))
	for i, r := range repos {
		out[i] = r.Slug
	}
	return out
}

func queryRemoteCandidates(localSlugs []string) []CrossNodeCandidate {
	if RemoteCandidateSupplier != nil {
		return RemoteCandidateSupplier()
	}
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil || len(conns) == 0 {
		return nil
	}
	return buildNodeCandidates(conns, localSlugs)
}

func buildNodeCandidates(conns []db.SSHConnection, slugs []string) []CrossNodeCandidate {
	var out []CrossNodeCandidate
	for _, c := range conns {
		a := c.Alias
		if a == "" {
			a = c.IPAddress
		}
		out = append(out, CrossNodeCandidate{NodeAlias: a, RepoSlug: "work", IsRemote: true})
		for _, s := range slugs {
			out = append(out, CrossNodeCandidate{NodeAlias: a, RepoSlug: s, IsRemote: true})
		}
	}
	return out
}

func formatCandidate(c CrossNodeCandidate) string {
	val := c.RepoSlug
	if val == "" {
		val = c.FullPath
	}
	if c.IsRemote && c.NodeAlias != "" {
		return c.NodeAlias + ":" + val
	}
	return val
}

func filterAndDedupe(items []string, prefix string) []string {
	seen := make(map[string]bool, len(items))
	var out []string
	for _, item := range items {
		hasMatch := prefix == "" || strings.HasPrefix(item, prefix)
		if item != "" && !seen[item] && hasMatch {
			seen[item] = true
			out = append(out, item)
		}
	}
	sort.Strings(out)
	return out
}
