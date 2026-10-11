package cmdpending

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func resolveAndInspectLocalReposWithDB(targetRepo string, cacheConn, backupConn *sql.DB, opts PendingCommitsOptions) ([]RepoPendingCommitRecord, error) {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		cwd = "."
	}
	records := ResolveWorkspaceRepositories(cwd)
	recordsToInspect, errTarget := filterWorkspaceRepositories(records, targetRepo)
	if errTarget != nil {
		return nil, errTarget
	}
	return inspectRepositoriesWithDB(recordsToInspect, cacheConn, backupConn, opts), nil
}

func inspectRepositoriesWithDB(records []model.ScanRecord, cacheConn, backupConn *sql.DB, opts PendingCommitsOptions) []RepoPendingCommitRecord {
	results := make([]RepoPendingCommitRecord, 0, len(records))
	for _, rec := range records {
		results = append(results, inspectSingleRepoWithDB(rec, cacheConn, backupConn, opts))
	}
	return results
}

func inspectSingleRepoWithDB(rec model.ScanRecord, cacheConn, backupConn *sql.DB, opts PendingCommitsOptions) RepoPendingCommitRecord {
	repoPath := resolveRelativePath(rec)
	dir := resolveRecordDir(rec)
	nowUnix := time.Now().Unix()

	if cachedRecord, isHit := tryGetCachedRecord(cacheConn, repoPath, dir, nowUnix, opts); isHit {
		return cachedRecord
	}

	liveRecord := inspectSingleRepoPendingCommits(rec)
	headSHA := queryRepoHeadSHA(dir)
	saveRepoToCache(cacheConn, repoPath, headSHA, liveRecord, nowUnix, opts.IsNoCache)
	_ = SaveBackupPendingStatus(backupConn, liveRecord, headSHA, nowUnix)

	return liveRecord
}

func tryGetCachedRecord(cacheConn *sql.DB, repoPath, dir string, nowUnix int64, opts PendingCommitsOptions) (RepoPendingCommitRecord, bool) {
	if opts.IsNoCache || opts.IsRefresh || cacheConn == nil {
		return RepoPendingCommitRecord{}, false
	}
	headSHA := queryRepoHeadSHA(dir)
	cached, hasHit, _ := GetCachedPendingStatus(cacheConn, repoPath, headSHA, nowUnix)
	if !hasHit || cached == nil {
		return RepoPendingCommitRecord{}, false
	}
	return buildCachedRepoPendingRecord(model.ScanRecord{AbsolutePath: dir}, dir, cached), true
}

func queryRepoHeadSHA(dir string) string {
	head, _ := currentPendingCommitsGitExecutor(dir, "rev-parse", "HEAD")
	return strings.TrimSpace(head)
}

func buildCachedRepoPendingRecord(rec model.ScanRecord, dir string, cached *PendingCommitCacheRecord) RepoPendingCommitRecord {
	branch := safeQueryRepoBranch(dir, rec.Branch)
	hasUpstream := checkRepoUpstream(dir)
	version := resolveRepoVersion(dir)
	shortVerBranch := formatShortVersionBranch(version, branch)
	out := assembleRepoPendingCommitRecord(
		rec, branch, version, shortVerBranch,
		0, cached.UncommittedCount, 0, cached.UncommittedCount,
		cached.UnpushedCount, hasUpstream, nil, nil,
	)
	out.IsDirty = cached.IsDirty
	out.HasUncommitted = cached.UncommittedCount > 0
	out.HasUnpushed = cached.UnpushedCount > 0
	out.IsClean = !out.IsDirty && !out.HasUnpushed
	if out.IsDirty {
		out.RemediationOptions = buildRepoRemediationOptions(out)
	}
	return out
}

func saveRepoToCache(conn *sql.DB, repoPath, headSHA string, rec RepoPendingCommitRecord, nowUnix int64, isNoCache bool) {
	if isNoCache || conn == nil {
		return
	}
	cacheRec := PendingCommitCacheRecord{
		RepoPath:         repoPath,
		HeadSHA:          headSHA,
		UncommittedCount: rec.TotalUncommitted,
		UnpushedCount:    rec.UnpushedCommitsCount,
		IsDirty:          rec.IsDirty,
		CachedAtUnix:     nowUnix,
	}
	_ = SaveCachedPendingStatus(conn, cacheRec)
}

func inspectSingleRepoPendingCommits(rec model.ScanRecord) RepoPendingCommitRecord {
	dir := resolveRecordDir(rec)
	rec.RepoName = resolveProperRepoName(dir, rec.RepoName)
	untracked, modified, staged, allFiles := queryRepoStatus(dir)
	totalUncommitted := len(allFiles)
	branch := safeQueryRepoBranch(dir, rec.Branch)
	hasUpstream := checkRepoUpstream(dir)
	unpushedCount, unpushedSHAs := resolveUnpushedCommits(dir, hasUpstream)
	version := resolveRepoVersion(dir)
	shortVerBranch := formatShortVersionBranch(version, branch)
	out := assembleRepoPendingCommitRecord(rec, branch, version, shortVerBranch, untracked, modified, staged, totalUncommitted, unpushedCount, hasUpstream, allFiles, unpushedSHAs)
	if out.IsDirty {
		out.RemediationOptions = buildRepoRemediationOptions(out)
	}
	return out
}

func resolveRecordDir(rec model.ScanRecord) string {
	if rec.AbsolutePath != "" {
		return rec.AbsolutePath
	}
	return rec.RelativePath
}

func checkRepoUpstream(dir string) bool {
	upstream, errUpstream := currentPendingCommitsGitExecutor(dir, "rev-parse", "--abbrev-ref", "@{u}")
	if errUpstream != nil {
		return false
	}
	clean := strings.TrimSpace(upstream)
	if clean == "" {
		return false
	}
	return !strings.Contains(strings.ToLower(clean), "fatal")
}

func assembleRepoPendingCommitRecord(
	rec model.ScanRecord, branch, version, shortVerBranch string,
	untracked, modified, staged, totalUncommitted, unpushedCount int,
	hasUpstream bool, pendingFiles, unpushedSHAs []string,
) RepoPendingCommitRecord {
	isDirty := totalUncommitted > 0
	out := RepoPendingCommitRecord{
		RepoName: rec.RepoName, RelativePath: resolveRelativePath(rec),
		CurrentBranch: branch, Version: version, ShortVersionBranch: shortVerBranch,
		IsDirty: isDirty, IsClean: evaluateIsClean(isDirty, unpushedCount > 0),
		HasUncommitted: isDirty, HasUnpushed: unpushedCount > 0, HasUpstream: hasUpstream,
		TotalUncommitted: totalUncommitted, UntrackedFilesCount: untracked,
		ModifiedFilesCount: modified, StagedFilesCount: staged,
		UnpushedCommitsCount: unpushedCount, PendingFiles: pendingFiles,
		UnpushedCommitSHAs: unpushedSHAs,
	}
	return out
}

func resolveRelativePath(rec model.ScanRecord) string {
	if rec.RelativePath != "" && !looksLikeHashOrOpaqueID(rec.RelativePath) && rec.RelativePath != "." {
		return rec.RelativePath
	}
	return rec.RepoName
}

func evaluateIsClean(isDirty, hasUnpushed bool) bool {
	if isDirty {
		return false
	}
	if hasUnpushed {
		return false
	}
	return true
}

func resolveRepoVersion(dir string) string {
	tag, errTag := currentPendingCommitsGitExecutor(dir, "describe", "--tags", "--abbrev=0")
	cleanTag := strings.TrimSpace(tag)
	if errTag == nil && cleanTag != "" && !strings.Contains(strings.ToLower(cleanTag), "fatal") {
		return cleanTag
	}
	return readFallbackVersion(dir)
}

func readFallbackVersion(dir string) string {
	v := readJsonVersion(filepath.Join(dir, ".gitmap", "version.json"))
	if v != "" {
		return ensureVersionVPrefix(v)
	}
	vRoot := readJsonVersion(filepath.Join(dir, "version.json"))
	if vRoot != "" {
		return ensureVersionVPrefix(vRoot)
	}
	return ""
}

func ensureVersionVPrefix(v string) string {
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}

func readJsonVersion(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}
	return extractVersionString(m)
}

func extractVersionString(m map[string]interface{}) string {
	if v, ok := m["Version"].(string); ok && len(v) > 0 {
		return v
	}
	if v, ok := m["version"].(string); ok && len(v) > 0 {
		return v
	}
	return ""
}

func formatShortVersionBranch(version, branch string) string {
	if version == "" {
		return truncateStringWithEllipsis(branch, 16)
	}
	combined := version + "/" + branch
	if utf8.RuneCountInString(combined) <= 16 {
		return combined
	}
	return abbreviateVersionBranch(version, branch)
}

func abbreviateVersionBranch(version, branch string) string {
	prefix := version + "/"
	prefixRunes := []rune(prefix)
	if len(prefixRunes) >= 15 {
		return truncateStringWithEllipsis(version+"/"+branch, 16)
	}

	avail := 16 - len(prefixRunes) - 1
	branchRunes := []rune(branch)
	if avail < len(branchRunes) {
		return string(prefixRunes) + string(branchRunes[:avail]) + "…"
	}

	return truncateStringWithEllipsis(version+"/"+branch, 16)
}

func truncateStringWithEllipsis(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 1 {
		return "…"
	}
	return string(runes[:maxRunes-1]) + "…"
}

func resolveUnpushedCommits(dir string, hasUpstream bool) (int, []string) {
	if !hasUpstream {
		return 0, nil
	}

	count := queryUnpushedCommitCount(dir)
	if count <= 0 {
		return 0, nil
	}

	shas := queryUnpushedCommitSHAs(dir)
	return count, shas
}

func queryUnpushedCommitCount(dir string) int {
	revCount, errCount := currentPendingCommitsGitExecutor(dir, "rev-list", "--count", "@{u}..HEAD")
	if errCount != nil {
		return 0
	}
	count, _ := strconv.Atoi(strings.TrimSpace(revCount))
	return count
}

func queryUnpushedCommitSHAs(dir string) []string {
	revLog, errLog := currentPendingCommitsGitExecutor(dir, "log", "--oneline", "-n", "10", "@{u}..HEAD")
	if errLog != nil {
		return nil
	}
	return parsePendingCommitLines(revLog)
}
