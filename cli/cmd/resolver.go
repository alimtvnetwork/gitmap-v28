package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResolveRepo finds a single repo matching target string across DB and JSON records.
func ResolveRepo(db *store.DB, target string) (*model.ScanRecord, error) {
	all := loadUnifiedCandidates(db)
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

	if rec := resolveByTarget(t, all); rec != nil {
		return rec, nil
	}

	return nil, apperror.NewNotFoundError(fmt.Sprintf("no repository matched %q", target))
}

// ResolveMultiRepos resolves a slice of targets (including globs) across DB and JSON pools.
func ResolveMultiRepos(db *store.DB, targets []string) ([]model.ScanRecord, []string) {
	all := loadUnifiedCandidates(db)
	var out []model.ScanRecord
	var missing []string
	seen := make(map[string]bool)

	for _, t := range targets {
		hits := resolveOneMulti(db, t, all)
		if len(hits) == 0 {
			missing = append(missing, t)
			continue
		}

		out = appendUniqueRecords(out, hits, seen)
	}

	return out, missing
}

func loadUnifiedCandidates(db *store.DB) []model.ScanRecord {
	dbRepos := loadDbCandidates(db)
	jsonRepos := loadJSONCandidates()

	return mergeCandidateRepos(dbRepos, jsonRepos)
}

func loadDbCandidates(db *store.DB) []model.ScanRecord {
	if db == nil {
		return nil
	}

	repos, err := db.ListRepos()
	if err != nil {
		return nil
	}

	return repos
}

func loadJSONCandidates() []model.ScanRecord {
	paths := []string{
		filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile),
		filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
	}
	for _, p := range paths {
		if records, isOk := tryLoadJSONPath(p); isOk {
			return records
		}
	}

	return nil
}

func tryLoadJSONPath(path string) ([]model.ScanRecord, bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, false
	}

	records, err := model.LoadStatusRecords(path)
	if err != nil || len(records) == 0 {
		return nil, false
	}

	return records, true
}

func mergeCandidateRepos(primary, secondary []model.ScanRecord) []model.ScanRecord {
	seen := make(map[string]bool)
	var out []model.ScanRecord

	for _, r := range primary {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	for _, r := range secondary {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	return out
}

func appendUniqueRecords(out, hits []model.ScanRecord, seen map[string]bool) []model.ScanRecord {
	for _, r := range hits {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	return out
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
	return resolveByTarget(target, all)
}

func resolveByTarget(target string, all []model.ScanRecord) *model.ScanRecord {
	for _, r := range all {
		if matchesRepoTarget(r, target) {
			return &r
		}
	}

	return nil
}

func matchesRepoTarget(r model.ScanRecord, target string) bool {
	t := strings.TrimSpace(target)
	if len(t) == 0 {
		return false
	}

	if matchesRepoNameOrSlug(r, t) {
		return true
	}

	return matchesRepoPath(r, t)
}

func matchesRepoNameOrSlug(r model.ScanRecord, t string) bool {
	if strings.EqualFold(r.Slug, t) || strings.EqualFold(r.RepoName, t) {
		return true
	}

	slugBase := filepath.Base(r.Slug)
	repoBase := filepath.Base(r.RepoName)

	return strings.EqualFold(slugBase, t) || strings.EqualFold(repoBase, t)
}

func matchesRepoPath(r model.ScanRecord, t string) bool {
	if fsutil.EqualPaths(r.AbsolutePath, t) {
		return true
	}

	dirBase := filepath.Base(fsutil.NormalizeSlashes(r.AbsolutePath))
	if strings.EqualFold(dirBase, t) {
		return true
	}

	abs, err := filepath.Abs(t)

	return err == nil && fsutil.EqualPaths(r.AbsolutePath, abs)
}

// PrintRepoSuggestions queries the database for suggestions and prints them.
func PrintRepoSuggestions(db *store.DB, target string) {
	if db == nil {
		return
	}

	suggs, _ := db.GetRepoSuggestions(target)
	if len(suggs) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "Did you mean:\n")
	for _, s := range suggs {
		fmt.Fprintf(os.Stderr, "  %s\n", s)
	}
}

// resolveEndpointString tries to resolve a user-provided string to an absolute path.
func resolveEndpointString(raw string) string {
	if isNetworkEndpoint(raw) {
		return raw
	}

	db, err := openDB()
	if err != nil || db == nil {
		return raw
	}

	defer db.Close()

	rec, err := ResolveRepo(db, raw)
	if err == nil && rec != nil {
		return rec.AbsolutePath
	}

	return raw
}

func isNetworkEndpoint(raw string) bool {
	lower := strings.ToLower(raw)
	prefixes := []string{"https://", "http://", "ssh://", "git@"}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}

	return false
}
