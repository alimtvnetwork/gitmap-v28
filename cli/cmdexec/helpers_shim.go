package cmdexec

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"strings"
)

// appendSummaryPart conditionally appends a colored summary segment.
func appendSummaryPart(parts []string, count int, color, format string) []string {
	if count == 0 {
		return parts
	}

	if len(color) > 0 {
		colored := fmt.Sprintf("%s"+format+"%s", color, count, constants.ColorReset)

		return append(parts, colored)
	}

	return append(parts, fmt.Sprintf(format, count))
}

func buildCommandArgs(args []string) string {
	if len(args) <= 1 {
		return ""
	}

	return strings.Join(args[1:], " ")
}

// closeTaskDB closes a *store.DB handle returned by createPendingTask
// when it is non-nil.
func closeTaskDB(db *store.DB) {
	if db == nil {
		return
	}

	_ = db.Close()
}

func loadAllRecordsDB() []model.ScanRecord {
	if LoadAllRecordsDBFn != nil {
		return deduplicatePullRecords(LoadAllRecordsDBFn())
	}

	return nil
}

// canonicalRepoPathKey produces a normalized path key respecting host OS case sensitivity.
func canonicalRepoPathKey(path string) string {
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
	pathKey := canonicalRepoPathKey(rec.AbsolutePath)
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
