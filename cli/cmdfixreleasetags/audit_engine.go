// Package cmdfixreleasetags provides tag discovery, GitHub release inspection,
// asset integrity checking, and CI/CD workflow auditing.
package cmdfixreleasetags

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isChecksumFileName(name string) bool {
	lower := strings.ToLower(name)

	return strings.HasSuffix(lower, ".sha256") ||
		strings.HasSuffix(lower, ".sha512") ||
		strings.HasSuffix(lower, ".md5") ||
		strings.Contains(lower, "checksum") ||
		strings.Contains(lower, "sha256sums")
}

func isBinaryArchiveName(name string) bool {
	lower := strings.ToLower(name)

	return strings.HasSuffix(lower, ".zip") ||
		strings.HasSuffix(lower, ".tar.gz") ||
		strings.HasSuffix(lower, ".tgz") ||
		strings.HasSuffix(lower, ".7z") ||
		strings.HasSuffix(lower, ".exe") ||
		strings.HasSuffix(lower, ".tar.xz")
}

// InspectAssets analyzes uploaded release assets for integrity and completeness.
func InspectAssets(rawAssets []ghAssetItem) ([]ReleaseAssetInfo, bool, bool, bool) {
	var assets []ReleaseAssetInfo
	hasChecksum := false
	hasBinary := false
	hasCorrupted := false

	for _, a := range rawAssets {
		isChecksum := isChecksumFileName(a.Name)
		isBinary := isBinaryArchiveName(a.Name)
		if isChecksum {
			hasChecksum = true
		}
		if isBinary {
			hasBinary = true
		}

		isCorrupt := a.Size <= 0
		if isCorrupt {
			hasCorrupted = true
		}

		assets = append(assets, ReleaseAssetInfo{
			Name:              a.Name,
			Size:              a.Size,
			State:             a.State,
			IsValid:           !isCorrupt,
			HasBinaryArchives: isBinary,
		})
	}

	return assets, hasChecksum, hasBinary, hasCorrupted
}

func hasFailingRun(runs []CIWorkflowRunInfo) bool {
	for _, r := range runs {
		conc := strings.ToLower(strings.TrimSpace(r.Conclusion))
		if conc == "failure" || strings.HasPrefix(conc, "cancel") || conc == "timed_out" || conc == "startup_failure" {
			return true
		}
	}

	return false
}

// DetermineAuditReason classifies the release state based on the 5 standard criteria.
func DetermineAuditReason(hasRelease bool, isDraft bool, assetCount int, hasChecksum bool, hasCorrupted bool, runs []CIWorkflowRunInfo) AuditReason {
	if !hasRelease {
		return ReasonOrphanTag
	}

	if isDraft {
		return ReasonDraftRelease
	}

	if assetCount == 0 || !hasChecksum {
		return ReasonMissingAssets
	}

	if hasCorrupted {
		return ReasonCorruptedAssets
	}

	if hasFailingRun(runs) {
		return ReasonCICDFailed
	}

	return ReasonHealthy
}

func resolveAuditInputs(opts AuditFilterOptions) (string, time.Duration, CommandExecutor) {
	activeVersion := constants.Version
	if opts.ActiveVersionOverride != "" {
		activeVersion = opts.ActiveVersionOverride
	}

	graceDuration := opts.GracePeriod
	if graceDuration <= 0 {
		graceDuration = 15 * time.Minute
	}

	executor := opts.CommandExecutor
	if executor == nil {
		executor = DefaultCommandExecutor{}
	}

	return activeVersion, graceDuration, executor
}

func buildReleaseMap(items []ghReleaseItem) map[string]ghReleaseItem {
	m := make(map[string]ghReleaseItem, len(items))
	for _, item := range items {
		m[item.TagName] = item
	}

	return m
}

func collectAllTagNames(localTags []string, remoteTags map[string]bool, releaseItems []ghReleaseItem) []string {
	tagSet := make(map[string]bool)
	for _, t := range localTags {
		tagSet[t] = true
	}
	for t := range remoteTags {
		tagSet[t] = true
	}
	for _, r := range releaseItems {
		tagSet[r.TagName] = true
	}

	var allTags []string
	for t := range tagSet {
		allTags = append(allTags, t)
	}

	SortTagsDescending(allTags)

	return allTags
}

func buildAuditRecord(tag string, repoPath string, hasLocal bool, hasRemote bool, rel *ghReleaseItem, exec CommandExecutor) ReleaseTagAuditRecord {
	rec := ReleaseTagAuditRecord{
		Tag:          tag,
		HasLocalTag:  hasLocal,
		HasRemoteTag: hasRemote,
	}

	if rel == nil {
		rec.AuditReason = ReasonOrphanTag

		return rec
	}

	rec.HasGitHubRelease = true
	rec.IsDraft = rel.IsDraft
	assets, hasChecksum, hasBinary, hasCorrupt := InspectAssets(rel.Assets)
	rec.Assets = assets
	rec.AssetCount = len(assets)
	rec.HasAssets = len(assets) > 0
	rec.HasChecksumFile = hasChecksum
	rec.HasBinaryArchives = hasBinary

	runs, _ := FetchTagWorkflows(repoPath, tag, "", exec)
	rec.WorkflowRuns = runs
	rec.AuditReason = DetermineAuditReason(true, rec.IsDraft, rec.AssetCount, hasChecksum, hasCorrupt, runs)

	return rec
}

func finalizeAuditRecords(records []ReleaseTagAuditRecord, activeVersion string, latestHealthy string, grace time.Duration) {
	for i := range records {
		rec := &records[i]
		EvaluateSafety(rec, activeVersion, latestHealthy, grace)
		if rec.ProtectionStatus != StatusEligible {
			rec.IsProtected = true
			rec.IsEligibleForDeletion = false
			continue
		}

		if rec.AuditReason != ReasonHealthy {
			rec.IsEligibleForDeletion = true
			rec.IsProtected = false
			continue
		}

		rec.IsEligibleForDeletion = false
		rec.IsProtected = false
	}
}

func buildAuditSummary(records []ReleaseTagAuditRecord) AuditSummary {
	var summary AuditSummary
	summary.TotalTagsChecked = len(records)

	for _, r := range records {
		countAuditCategory(&summary, r)
	}

	return summary
}

func countAuditCategory(summary *AuditSummary, r ReleaseTagAuditRecord) {
	if r.IsProtected {
		summary.ProtectedTagsCount++
	}
	if r.IsEligibleForDeletion {
		summary.EligibleDeletionsCount++
	}

	switch r.AuditReason {
	case ReasonHealthy:
		summary.HealthyReleasesCount++
	case ReasonOrphanTag:
		summary.OrphanTagsCount++
	case ReasonDraftRelease:
		summary.DraftReleasesCount++
	case ReasonMissingAssets:
		summary.MissingAssetsCount++
	case ReasonCorruptedAssets:
		summary.CorruptedAssetsCount++
	case ReasonCICDFailed:
		summary.CICDFailedCount++
	}
}

// AuditReleaseTags performs the comprehensive audit across Git and GitHub release tags.
func AuditReleaseTags(repoPath string, opts AuditFilterOptions) (*AuditReport, error) {
	activeVer, grace, exec := resolveAuditInputs(opts)

	localTags, _ := CollectLocalTags(repoPath, exec)
	remoteTags, _ := CollectRemoteTags(repoPath, exec)
	releases, _ := FetchGitHubReleases(repoPath, exec)

	releaseMap := buildReleaseMap(releases)
	allTags := collectAllTagNames(localTags, remoteTags, releases)

	var records []ReleaseTagAuditRecord
	for _, tag := range allTags {
		var relPtr *ghReleaseItem
		if item, ok := releaseMap[tag]; ok {
			relPtr = &item
		}
		rec := buildAuditRecord(tag, repoPath, containsTag(localTags, tag), remoteTags[tag], relPtr, exec)
		records = append(records, rec)
	}

	latestHealthy := FindLatestHealthyTag(records)
	finalizeAuditRecords(records, activeVer, latestHealthy, grace)
	summary := buildAuditSummary(records)

	return &AuditReport{
		RepoPath:         repoPath,
		ActiveVersion:    activeVer,
		LatestHealthyTag: latestHealthy,
		Summary:          summary,
		Records:          records,
	}, nil
}

func containsTag(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}

	return false
}
