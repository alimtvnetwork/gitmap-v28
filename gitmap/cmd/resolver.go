package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ResolveRepo finds a single repo matching target string.
func ResolveRepo(db *store.DB, target string) (*model.ScanRecord, error) {
	all, err := db.ListRepos()
	if err != nil {
		return nil, fmt.Errorf("list repos: %w", err)
	}
	t := strings.TrimSpace(target)
	if len(t) == 0 || t == "." {
		return resolveByPWD(all)
	}
	if rec := resolveByPath(t, all); rec != nil {
		return rec, nil
	}
	if rec := resolveByAlias(db, t, all); rec != nil {
		return rec, nil
	}
	if rec := resolveBySlug(t, all); rec != nil {
		return rec, nil
	}
	return nil, fmt.Errorf("no repository matched %q", target)
}

// ResolveMultiRepos resolves a slice of targets (including globs).
func ResolveMultiRepos(db *store.DB, targets []string) ([]model.ScanRecord, []string) {
	all, err := db.ListRepos()
	if err != nil {
		return nil, targets
	}
	var out []model.ScanRecord
	var missing []string
	seen := make(map[int64]bool)

	for _, t := range targets {
		hits := resolveOneMulti(db, t, all)
		if len(hits) == 0 {
			missing = append(missing, t)
			continue
		}
		for _, r := range hits {
			if !seen[r.ID] {
				seen[r.ID] = true
				out = append(out, r)
			}
		}
	}
	return out, missing
}

func resolveOneMulti(db *store.DB, target string, all []model.ScanRecord) []model.ScanRecord {
	t := strings.TrimSpace(target)
	if isGlob(t) {
		return resolveByGlob(t, all)
	}
	rec, err := ResolveRepo(db, t)
	if err == nil && rec != nil {
		return []model.ScanRecord{*rec}
	}
	return nil
}

func resolveBySlug(target string, all []model.ScanRecord) *model.ScanRecord {
	tLow := strings.ToLower(target)
	for _, r := range all {
		if strings.EqualFold(r.Slug, tLow) || strings.EqualFold(r.RepoName, tLow) {
			return &r
		}
	}
	return nil
}

// PrintRepoSuggestions queries the database for suggestions and prints them.
func PrintRepoSuggestions(db *store.DB, target string) {
	suggs, _ := db.GetRepoSuggestions(target)
	if len(suggs) > 0 {
		fmt.Fprintf(os.Stderr, "Did you mean:\n")
		for _, s := range suggs {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
	}
}

// resolveEndpointString tries to resolve a user-provided string to an absolute path.
func resolveEndpointString(raw string) string {
	lower := strings.ToLower(raw)
	for _, p := range []string{"https://", "http://", "ssh://", "git@"} {
		if strings.HasPrefix(lower, p) {
			return raw
		}
	}
	db, err := openDB()
	if err == nil && db != nil {
		defer db.Close()
		if rec, err := ResolveRepo(db, raw); err == nil && rec != nil {
			return rec.AbsolutePath
		}
	}
	return raw
}
