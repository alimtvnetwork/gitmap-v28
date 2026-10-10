// Package cmdfixreleasetags provides tag discovery, GitHub release inspection,
// asset integrity checking, and CI/CD workflow auditing.
package cmdfixreleasetags

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type ghReleaseItem struct {
	TagName   string        `json:"tagName"`
	IsDraft   bool          `json:"isDraft"`
	CreatedAt string        `json:"createdAt"`
	Assets    []ghAssetItem `json:"assets"`
}

type ghAssetItem struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	State string `json:"state"`
}

type ghRunItem struct {
	DatabaseId uint64 `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	CreatedAt  string `json:"createdAt"`
}

// CollectLocalTags reads local git tags via git tag -l.
func CollectLocalTags(repoPath string, executor CommandExecutor) ([]string, error) {
	out, err := executor.Run(repoPath, "git", "tag", "-l")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var tags []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	return tags, nil
}

// CollectRemoteTags reads remote git tags via git ls-remote --tags origin.
func CollectRemoteTags(repoPath string, executor CommandExecutor) (map[string]bool, error) {
	out, err := executor.Run(repoPath, "git", "ls-remote", "--tags", "origin")
	if err != nil {
		return nil, err
	}

	tags := make(map[string]bool)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		processRemoteTagLine(line, tags)
	}

	return tags, nil
}

func processRemoteTagLine(line string, tags map[string]bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return
	}

	parts := strings.Fields(trimmed)
	if len(parts) < 2 || strings.HasSuffix(parts[1], "^{}") {
		return
	}

	const prefix = "refs/tags/"
	if strings.HasPrefix(parts[1], prefix) {
		tags[strings.TrimPrefix(parts[1], prefix)] = true
	}
}

// FetchGitHubReleases retrieves repository releases via gh release list.
func FetchGitHubReleases(repoPath string, executor CommandExecutor) ([]ghReleaseItem, error) {
	if os.Getenv("GITMAP_MOCK_GH") == "1" {
		return nil, nil
	}

	out, err := executor.Run(repoPath, "gh", "release", "list", "--limit", "100", "--json", "tagName,isDraft,createdAt,assets")
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" || trimmed == "[]" {
		return nil, nil
	}

	var items []ghReleaseItem
	err = json.Unmarshal([]byte(trimmed), &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FetchTagWorkflows retrieves CI/CD workflow runs for a specific tag or commit.
func FetchTagWorkflows(repoPath, tag, commitSha string, executor CommandExecutor) ([]CIWorkflowRunInfo, error) {
	if os.Getenv("GITMAP_MOCK_GH") == "1" {
		return nil, nil
	}

	args := buildWorkflowQueryArgs(tag, commitSha)
	out, err := executor.Run(repoPath, "gh", args...)
	if err != nil {
		return nil, err
	}

	return parseWorkflowRunItems(out)
}

func buildWorkflowQueryArgs(tag, commitSha string) []string {
	args := []string{"run", "list", "--limit", "10", "--json", "databaseId,name,status,conclusion,createdAt"}
	if commitSha != "" {
		return append(args, "--commit", commitSha)
	}

	if tag != "" {
		return append(args, "--branch", tag)
	}

	return args
}

func parseWorkflowRunItems(data []byte) ([]CIWorkflowRunInfo, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "[]" {
		return nil, nil
	}

	var rawItems []ghRunItem
	err := json.Unmarshal([]byte(trimmed), &rawItems)
	if err != nil {
		return nil, err
	}

	var runs []CIWorkflowRunInfo
	for _, r := range rawItems {
		parsedTime, _ := time.Parse(time.RFC3339, r.CreatedAt)
		isSucc := strings.EqualFold(r.Conclusion, "success")
		isInProg := isRunInProgress(CIWorkflowRunInfo{Status: r.Status})
		runs = append(runs, CIWorkflowRunInfo{
			RunId:        r.DatabaseId,
			WorkflowName: r.Name,
			Status:       r.Status,
			Conclusion:   r.Conclusion,
			CreatedAt:    parsedTime,
			IsSuccessful: isSucc,
			IsInProgress: isInProg,
		})
	}

	return runs, nil
}

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
		if conc == "failure" || conc == "cancelled" || conc == "timed_out" || conc == "startup_failure" {
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

func splitSemverInts(v string) [3]int {
	norm := NormalizeVersionTag(v)
	var nums [3]int
	parts := strings.SplitN(norm, ".", 3)
	for i, p := range parts {
		if i >= 3 {
			break
		}

		n := 0
		for _, r := range p {
			if r < '0' || r > '9' {
				break
			}
			n = n*10 + int(r-'0')
		}
		nums[i] = n
	}

	return nums
}

// CompareSemver compares two semantic versions. Returns 1 if a > b, -1 if a < b, 0 if equal.
func CompareSemver(a, b string) int {
	pa := splitSemverInts(a)
	pb := splitSemverInts(b)

	for i := 0; i < 3; i++ {
		if pa[i] > pb[i] {
			return 1
		}
		if pa[i] < pb[i] {
			return -1
		}
	}

	return 0
}

// SortTagsDescending sorts tag strings in descending semantic version order.
func SortTagsDescending(tags []string) {
	sort.Slice(tags, func(i, j int) bool {
		cmp := CompareSemver(tags[i], tags[j])
		if cmp != 0 {
			return cmp > 0
		}

		return tags[i] > tags[j]
	})
}

// FindLatestHealthyTag locates the highest semver release that is healthy.
func FindLatestHealthyTag(records []ReleaseTagAuditRecord) string {
	var healthyTags []string
	for _, rec := range records {
		if rec.AuditReason == ReasonHealthy {
			healthyTags = append(healthyTags, rec.Tag)
		}
	}

	if len(healthyTags) == 0 {
		return ""
	}

	SortTagsDescending(healthyTags)

	return healthyTags[0]
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
