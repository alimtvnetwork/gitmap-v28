// Package cmd — cg_resolver.go resolves flexible repository specifiers (path, alias, ID, URL).
package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ResolveCGTarget resolves a single specifier (path, alias, ID, URL) to an absolute repo path.
func ResolveCGTarget(spec string) (string, bool) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return "", false
	}

	// 1. Check if it is a local filesystem path
	if info, err := os.Stat(trimmed); err == nil && info.IsDir() {
		abs, _ := filepath.Abs(trimmed)
		return abs, true
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		return "", false
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil {
		return "", false
	}

	// 2. Check if it matches a numeric ID
	if path, ok := matchRepoByID(repos, trimmed); ok {
		return path, true
	}

	// 3. Check alias / slug / repo name
	for _, r := range repos {
		if strings.EqualFold(r.Slug, trimmed) || strings.EqualFold(r.RepoName, trimmed) {
			return r.AbsolutePath, true
		}
	}

	// 4. Check clone URL or remote URL
	for _, r := range repos {
		if strings.EqualFold(r.HTTPSUrl, trimmed) || strings.EqualFold(r.SSHUrl, trimmed) || strings.EqualFold(r.DiscoveredURL, trimmed) {
			return r.AbsolutePath, true
		}
	}

	return "", false
}

func matchRepoByID(repos []model.ScanRecord, idStr string) (string, bool) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return "", false
	}
	for _, r := range repos {
		if r.ID == id {
			return r.AbsolutePath, true
		}
	}
	return "", false
}

// ResolveAllCGTargets resolves a slice of specifiers into absolute filesystem paths.
func ResolveAllCGTargets(specs []string) []string {
	var resolved []string
	seen := make(map[string]bool)

	for _, spec := range specs {
		path, ok := ResolveCGTarget(spec)
		if !ok || seen[path] {
			continue
		}
		seen[path] = true
		resolved = append(resolved, path)
	}
	return resolved
}
