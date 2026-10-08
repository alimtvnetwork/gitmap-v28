package cmdpull

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// CanonicalRepoPathKey produces a normalized path key respecting host OS case sensitivity.
func CanonicalRepoPathKey(path string) string {
	return fsutil.CanonicalPathKey(path)
}

// deduplicatePullRecords filters duplicates from a slice of ScanRecord.
// It ensures each repository canonical path and slug is enqueued exactly once.
func deduplicatePullRecords(records []model.ScanRecord) []model.ScanRecord {
	if len(records) <= 1 {
		return records
	}
	seenPath := make(map[string]bool, len(records))
	seenSlug := make(map[string]bool, len(records))
	unique := make([]model.ScanRecord, 0, len(records))
	for _, rec := range records {
		unique = appendUniquePullRecord(unique, rec, seenPath, seenSlug)
	}

	return unique
}

func appendUniquePullRecord(unique []model.ScanRecord, rec model.ScanRecord, seenPath, seenSlug map[string]bool) []model.ScanRecord {
	pathKey := CanonicalRepoPathKey(rec.AbsolutePath)
	slugKey := resolveRecordSlugKey(rec)
	if isRecordDuplicate(pathKey, slugKey, seenPath, seenSlug) {
		return unique
	}
	markRecordSeen(pathKey, slugKey, seenPath, seenSlug)

	return append(unique, rec)
}

func resolveRecordSlugKey(rec model.ScanRecord) string {
	slug := strings.ToLower(strings.TrimSpace(rec.Slug))
	if slug != "" {
		return slug
	}

	return strings.ToLower(strings.TrimSpace(rec.RepoName))
}

func isRecordDuplicate(pathKey, slugKey string, seenPath, seenSlug map[string]bool) bool {
	if pathKey != "" && seenPath[pathKey] {
		return true
	}
	if slugKey != "" && seenSlug[slugKey] {
		return true
	}

	return false
}

func markRecordSeen(pathKey, slugKey string, seenPath, seenSlug map[string]bool) {
	if pathKey != "" {
		seenPath[pathKey] = true
	}
	if slugKey != "" {
		seenSlug[slugKey] = true
	}
}

func resolvePullTargets(slug, groupName string, all bool) []model.ScanRecord {
	raw := resolveRawPullTargets(slug, groupName, all)

	return deduplicatePullRecords(raw)
}

func resolveRawPullTargets(slug, groupName string, all bool) []model.ScanRecord {
	if HasAlias() {
		return resolveAliasRecord()
	}
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}
	if all {
		return optimizeAndLoadAllRecordsDB()
	}

	return resolveSlugTarget(slug)
}

func optimizeAndLoadAllRecordsDB() []model.ScanRecord {
	db, err := openDB()
	if err == nil {
		_, _ = OptimizeRedundantRepos(db, false)
		db.Close()
	}

	return loadAllRecordsDB()
}

func resolveAliasRecord() []model.ScanRecord {
	return []model.ScanRecord{{
		RepoName:     GetAliasSlug(),
		Slug:         GetAliasSlug(),
		AbsolutePath: GetAliasPath(),
	}}
}

func resolveSlugTarget(slug string) []model.ScanRecord {
	if len(slug) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrPullSlugRequired)

		return nil
	}

	return lookupBySlugDBFirst(slug)
}

func lookupBySlugDBFirst(slug string) []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		return lookupBySlugJSON(slug)
	}
	defer db.Close()
	repos, dbErr := db.FindBySlug(strings.ToLower(slug))
	foundRepos := dbErr == nil && len(repos) > 0
	if foundRepos {
		return repos
	}

	return lookupBySlugJSON(slug)
}

func lookupBySlugJSON(slug string) []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadJSONRecords(jsonPath)
	if err != nil {
		return nil
	}

	return findBySlug(records, slug)
}

func loadJSONRecords(path string) ([]model.ScanRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open json records")
	}
	defer file.Close()
	var records []model.ScanRecord
	err = json.NewDecoder(file).Decode(&records)
	if err != nil {
		return nil, apperror.WrapSimple(err, "decode json records")
	}

	return records, nil
}

func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	exact, partial := partitionBySlug(records, slug)
	if len(exact) > 0 {
		return exact
	}

	return partial
}

func partitionBySlug(records []model.ScanRecord, slug string) ([]model.ScanRecord, []model.ScanRecord) {
	var exact, partial []model.ScanRecord
	for _, r := range records {
		if strings.EqualFold(r.RepoName, slug) {
			exact = append(exact, r)
		} else if strings.Contains(strings.ToLower(r.RepoName), strings.ToLower(slug)) {
			partial = append(partial, r)
		}
	}

	return exact, partial
}

func pullOneRepo(rec model.ScanRecord) {
	fmt.Printf(constants.MsgPullStarting, rec.RepoName, rec.AbsolutePath)
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		fmt.Fprintf(os.Stderr, constants.ErrPullNotRepo, rec.AbsolutePath)

		return
	}
	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.IsSuccess {
		fmt.Printf(constants.MsgPullSuccess, rec.RepoName)
	} else {
		fmt.Fprintf(os.Stderr, constants.MsgPullFailed, rec.RepoName, result.Error)
	}
}

func findChildrenOfCWD(cwd string) []model.ScanRecord {
	all := loadAllRecordsDB()
	var children []model.ScanRecord
	prefix := ensureTrailingPathSep(cwd)
	for _, r := range all {
		if strings.HasPrefix(r.AbsolutePath, prefix) || r.AbsolutePath == cwd {
			children = append(children, r)
		}
	}

	return deduplicatePullRecords(children)
}

func ensureTrailingPathSep(path string) string {
	sep := string(os.PathSeparator)
	if strings.HasSuffix(path, sep) {
		return path
	}

	return path + sep
}
