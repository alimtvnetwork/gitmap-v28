package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// syncRecordsToVSCodePM upserts every scanned repo into projects.json.
// noVSCodeSync skips the entire sync (mirrors the clone family).
// noAutoTags keeps the sync but suppresses marker-based tag detection
// so the user's existing tags remain authoritative.
func syncRecordsToVSCodePM(records []model.ScanRecord, isSkipVSCodeSync, isSkipAutoTags bool) {
	if isSkipVSCodeSync {
		return
	}

	pairs := buildScanPMPairs(records, isSkipAutoTags)
	syncClonedReposToVSCodePM(pairs, false)
}

// buildScanPMPairs maps ScanRecords to vscodepm.Pair tuples, applying
// the canonicalization + auto-tag rules used by the clone family so
// scan-discovered and clone-discovered rows are byte-identical.
func buildScanPMPairs(records []model.ScanRecord, isSkipAutoTags bool) []vscodepm.Pair {
	pairs := make([]vscodepm.Pair, 0, len(records))

	for _, rec := range records {
		canonical := canonicalizePMPath(rec.AbsolutePath)

		pair := vscodepm.Pair{
			RootPath: canonical,
			Name:     rec.RepoName,
		}

		if !isSkipAutoTags {
			pair.Tags = vscodepm.DetectTagsCustom(canonical)
		}

		pairs = append(pairs, pair)
	}

	return pairs
}
