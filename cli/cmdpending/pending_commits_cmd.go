// Package cmd — pending_commits_cmd.go inspects and classifies uncommitted changes and unpushed commits.
package cmdpending

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// gitCmdExecutor defines a function that executes a git command in a specific directory.
type gitCmdExecutor func(dir string, args ...string) (string, error)

var defaultPendingCommitsGitExecutor gitCmdExecutor = func(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var currentPendingCommitsGitExecutor = defaultPendingCommitsGitExecutor

// pendingCommitsSSHExecutor defines an interface for executing commands over SSH.
type pendingCommitsSSHExecutor interface {
	Execute(conn db.SSHConnection, cmd, shell string) (string, error)
}

type defaultPendingCommitsSSHRunner struct{}

func (r defaultPendingCommitsSSHRunner) Execute(conn db.SSHConnection, cmd, shell string) (string, error) {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return "", errConnect
	}
	defer client.Close()

	out, errExec := secrets.RunCommand(client, cmd, shell)
	if errExec != nil && shell == "bash" {
		out, errExec = secrets.RunCommand(client, cmd, "sh")
	}
	return out, errExec
}

var currentPendingCommitsSSHExecutor pendingCommitsSSHExecutor = defaultPendingCommitsSSHRunner{}

// runPendingCommits is the CLI dispatch entry point.
// RunPendingCommits executes discovery of uncommitted and unpushed changes.
func RunPendingCommits(args []string) error {
	if isPendingCommitsHelpRequested(args) {
		printPendingCommitsHelp()
		return nil
	}

	opts := parsePendingCommitsOptions(args)
	if opts.IsSSH {
		return runFleetSSHPendingCommits(opts)
	}

	return runLocalPendingCommits(opts)
}

func isPendingCommitsHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

func parsePendingCommitsOptions(args []string) PendingCommitsOptions {
	opts := defaultPendingCommitsOptions()
	for i := 0; i < len(args); i++ {
		i = processPendingCommitArg(args, i, &opts)
	}
	normalizePendingCommitsOptions(&opts)
	return opts
}

func defaultPendingCommitsOptions() PendingCommitsOptions {
	return PendingCommitsOptions{
		SortMode:    "priority",
		DetailMode:  "summary",
		IsDirtyOnly: true,
	}
}

func normalizePendingCommitsOptions(opts *PendingCommitsOptions) {
	if opts.SortMode != "name" && opts.SortMode != "count" && opts.SortMode != "priority" {
		opts.SortMode = "priority"
	}
	if opts.DetailMode == "1" {
		opts.DetailMode = "1-by-1"
	}
}

func processPendingCommitArg(args []string, i int, opts *PendingCommitsOptions) int {
	arg := args[i]
	switch {
	case arg == "-s" || arg == "--sort" || strings.HasPrefix(arg, "--sort="):
		return parseSortFlag(args, i, opts)
	case arg == "-d" || arg == "--detail" || strings.HasPrefix(arg, "--detail="):
		return parseDetailFlag(args, i, opts)
	default:
		applyPendingBooleanOrTarget(arg, opts)
		return i
	}
}

func parseSortFlag(args []string, i int, opts *PendingCommitsOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--sort=") {
		opts.SortMode = strings.ToLower(strings.TrimPrefix(arg, "--sort="))
		return i
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		opts.SortMode = strings.ToLower(args[i+1])
		return i + 1
	}
	return i
}

func parseDetailFlag(args []string, i int, opts *PendingCommitsOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--detail=") {
		opts.DetailMode = strings.ToLower(strings.TrimPrefix(arg, "--detail="))
		return i
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		opts.DetailMode = strings.ToLower(args[i+1])
		return i + 1
	}
	opts.DetailMode = "all"
	return i
}

func applyPendingBooleanOrTarget(arg string, opts *PendingCommitsOptions) {
	if applyPendingScopeFlags(arg, opts) {
		return
	}
	applyPendingCacheAndTarget(arg, opts)
}

func applyPendingScopeFlags(arg string, opts *PendingCommitsOptions) bool {
	switch {
	case arg == "-j" || arg == "--json":
		opts.IsJSON = true
		return true
	case arg == "--ssh" || arg == "-ssh":
		opts.IsSSH = true
		return true
	case arg == "--dirty-only":
		opts.IsDirtyOnly, opts.IsAll = true, false
		return true
	case arg == "--all":
		opts.IsAll, opts.IsDirtyOnly = true, false
		return true
	}
	return false
}

func applyPendingCacheAndTarget(arg string, opts *PendingCommitsOptions) bool {
	switch {
	case arg == "--no-cache":
		opts.IsNoCache = true
		return true
	case arg == "--refresh":
		opts.IsRefresh = true
		return true
	case arg == "--backup" || arg == "--serve-backup":
		opts.IsBackup = true
		return true
	case !strings.HasPrefix(arg, "-"):
		opts.TargetRepo = arg
		return true
	}
	return false
}

func runLocalPendingCommits(opts PendingCommitsOptions) error {
	var cacheConn *sql.DB
	if !opts.IsNoCache {
		var errCache error
		cacheConn, errCache = OpenPendingCommitsCache("")
		if errCache == nil && cacheConn != nil {
			defer cacheConn.Close()
		}
	}

	backupConn, errBackup := OpenPendingCommitsBackup("")
	if errBackup == nil && backupConn != nil {
		defer backupConn.Close()
	}

	if opts.IsBackup {
		if errBackup != nil {
			return errBackup
		}
		if backupConn == nil {
			return apperror.NewSimple("backup database is not available", "E9004")
		}
		return serveBackupPendingCommits(backupConn, opts)
	}

	inspected, errInspect := resolveAndInspectLocalReposWithDB(opts.TargetRepo, cacheConn, backupConn, opts)
	if errInspect != nil {
		return errInspect
	}
	payload := buildPendingCommitsPayload(inspected, opts)
	if opts.IsJSON {
		return emitPendingCommitsJSON(payload)
	}
	renderLocalPendingCommitsTerminal(payload, len(inspected), inspected, opts)
	return nil
}

func serveBackupPendingCommits(backupConn *sql.DB, opts PendingCommitsOptions) error {
	backupRecords, errList := ListAllBackupPendingStatuses(backupConn)
	if errList != nil {
		return errList
	}
	inspected := convertBackupToPendingRecords(backupRecords, opts.TargetRepo)
	payload := buildPendingCommitsPayload(inspected, opts)
	if opts.IsJSON {
		return emitPendingCommitsJSON(payload)
	}
	renderLocalPendingCommitsTerminal(payload, len(inspected), inspected, opts)
	return nil
}

func convertBackupToPendingRecords(backupRecords []PendingCommitBackupRecord, targetRepo string) []RepoPendingCommitRecord {
	var inspected []RepoPendingCommitRecord
	for _, b := range backupRecords {
		if !isTargetRepoMatch(b, targetRepo) {
			continue
		}
		var item RepoPendingCommitRecord
		if b.PayloadJSON != "" && b.PayloadJSON != "{}" {
			if err := json.Unmarshal([]byte(b.PayloadJSON), &item); err == nil && item.RepoName != "" {
				if item.IsDirty && len(item.RemediationOptions) == 0 {
					item.RemediationOptions = buildRepoRemediationOptions(item)
				}
				inspected = append(inspected, item)
				continue
			}
		}
		item = RepoPendingCommitRecord{
			RepoName:             b.RepoName,
			RelativePath:         b.RepoPath,
			CurrentBranch:        b.CurrentBranch,
			Version:              b.Version,
			ShortVersionBranch:   b.ShortVersionBranch,
			IsDirty:              b.IsDirty,
			IsClean:              !b.IsDirty && b.UnpushedCount == 0,
			HasUncommitted:       b.UncommittedCount > 0,
			HasUnpushed:          b.UnpushedCount > 0,
			TotalUncommitted:     b.UncommittedCount,
			UnpushedCommitsCount: b.UnpushedCount,
		}
		if item.IsDirty {
			item.RemediationOptions = buildRepoRemediationOptions(item)
		}
		inspected = append(inspected, item)
	}
	return inspected
}

func isTargetRepoMatch(b PendingCommitBackupRecord, targetRepo string) bool {
	if targetRepo == "" || strings.EqualFold(targetRepo, "all") {
		return true
	}
	return strings.EqualFold(b.RepoName, targetRepo) || strings.EqualFold(b.RepoPath, targetRepo)
}

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

	if !opts.IsNoCache && !opts.IsRefresh && cacheConn != nil {
		headSHA := queryRepoHeadSHA(dir)
		cached, hasHit, _ := GetCachedPendingStatus(cacheConn, repoPath, headSHA, nowUnix)
		if hasHit && cached != nil {
			return buildCachedRepoPendingRecord(rec, dir, cached)
		}
	}

	liveRecord := inspectSingleRepoPendingCommits(rec)
	headSHA := queryRepoHeadSHA(dir)
	saveRepoToCache(cacheConn, repoPath, headSHA, liveRecord, nowUnix, opts.IsNoCache)
	_ = SaveBackupPendingStatus(backupConn, liveRecord, headSHA, nowUnix)

	return liveRecord
}

func queryRepoHeadSHA(dir string) string {
	head, _ := currentPendingCommitsGitExecutor(dir, "rev-parse", "HEAD")
	return strings.TrimSpace(head)
}

func buildCachedRepoPendingRecord(rec model.ScanRecord, dir string, cached *PendingCommitCacheRecord) RepoPendingCommitRecord {
	branch := queryRepoBranch(dir, rec.Branch)
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

func resolveAndInspectLocalRepos(targetRepo string) ([]RepoPendingCommitRecord, error) {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		cwd = "."
	}
	records := ResolveWorkspaceRepositories(cwd)
	recordsToInspect, errTarget := filterWorkspaceRepositories(records, targetRepo)
	if errTarget != nil {
		return nil, errTarget
	}
	return inspectRepositoriesPendingCommits(recordsToInspect), nil
}

func buildPendingCommitsPayload(inspected []RepoPendingCommitRecord, opts PendingCommitsOptions) PendingCommitsPayload {
	dirtyCount, totalUncommitted, totalUnpushed := aggregatePendingCounts(inspected)
	displayed := filterDisplayedPendingCommits(inspected, opts.IsDirtyOnly)
	sortPendingCommits(displayed, opts.SortMode)
	return PendingCommitsPayload{
		Timestamp:             time.Now().UTC(),
		TotalReposScanned:     len(inspected),
		TotalDirtyRepos:       dirtyCount,
		TotalUncommittedFiles: totalUncommitted,
		TotalUnpushedCommits:  totalUnpushed,
		SortMode:              opts.SortMode,
		DetailMode:            opts.DetailMode,
		Repositories:          displayed,
	}
}

func aggregatePendingCounts(inspected []RepoPendingCommitRecord) (int, int, int) {
	dirtyCount, totalUncommitted, totalUnpushed := 0, 0, 0
	for _, item := range inspected {
		if item.IsDirty {
			dirtyCount++
		}
		totalUncommitted += item.TotalUncommitted
		totalUnpushed += item.UnpushedCommitsCount
	}
	return dirtyCount, totalUncommitted, totalUnpushed
}

func filterDisplayedPendingCommits(inspected []RepoPendingCommitRecord, isDirtyOnly bool) []RepoPendingCommitRecord {
	var displayed []RepoPendingCommitRecord
	for _, item := range inspected {
		if isPendingCommitDisplayed(item, isDirtyOnly) {
			displayed = append(displayed, item)
		}
	}
	return displayed
}

func renderLocalPendingCommitsTerminal(payload PendingCommitsPayload, totalScanned int, inspected []RepoPendingCommitRecord, opts PendingCommitsOptions) {
	renderPendingCommitsTerminal(payload, totalScanned, inspected)
	if opts.DetailMode == "all" {
		renderPendingCommitsDetailsAll(payload.Repositories)
	}
	if opts.DetailMode == "1-by-1" {
		renderPendingCommitsInteractive(payload.Repositories)
	}
}

func filterWorkspaceRepositories(records []model.ScanRecord, targetRepo string) ([]model.ScanRecord, error) {
	if targetRepo == "" || strings.EqualFold(targetRepo, "all") {
		return records, nil
	}

	var matched []model.ScanRecord
	for _, r := range records {
		if strings.EqualFold(r.RepoName, targetRepo) ||
			strings.EqualFold(r.Slug, targetRepo) ||
			strings.EqualFold(r.RelativePath, targetRepo) {
			matched = append(matched, r)
		}
	}
	if len(matched) == 0 {
		return nil, apperror.NewSimple(fmt.Sprintf("repository '%s' not found", targetRepo), "E9002")
	}
	return matched, nil
}

func isPendingCommitDisplayed(item RepoPendingCommitRecord, isDirtyOnly bool) bool {
	if !isDirtyOnly {
		return true
	}
	return item.IsDirty || item.HasUnpushed
}

func ResolveWorkspaceRepositories(cwd string) []model.ScanRecord {
	targets := cmdpull.ResolvePullDirectoryTargets(cwd)
	if len(targets) > 0 {
		return targets
	}

	gitPath := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		base := filepath.Base(cwd)
		return []model.ScanRecord{
			{
				RepoName:     base,
				Slug:         base,
				AbsolutePath: cwd,
				RelativePath: ".",
			},
		}
	}

	return nil
}

func inspectRepositoriesPendingCommits(records []model.ScanRecord) []RepoPendingCommitRecord {
	results := make([]RepoPendingCommitRecord, 0, len(records))
	for _, rec := range records {
		results = append(results, inspectSingleRepoPendingCommits(rec))
	}
	return results
}

func inspectSingleRepoPendingCommits(rec model.ScanRecord) RepoPendingCommitRecord {
	dir := resolveRecordDir(rec)
	untracked, modified, staged, allFiles := queryRepoStatus(dir)
	totalUncommitted := len(allFiles)
	branch := queryRepoBranch(dir, rec.Branch)
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

func queryRepoStatus(dir string) (int, int, int, []string) {
	porcelain, _ := currentPendingCommitsGitExecutor(dir, "status", "--porcelain")
	return ParsePorcelainStatusLines(porcelain)
}

func queryRepoBranch(dir, fallback string) string {
	branch, _ := currentPendingCommitsGitExecutor(dir, "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if branch != "" {
		return branch
	}
	if fallback != "" {
		return fallback
	}
	return "main"
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
	if rec.RelativePath != "" {
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
	if errTag == nil {
		cleanTag := strings.TrimSpace(tag)
		if cleanTag != "" && !strings.Contains(strings.ToLower(cleanTag), "fatal") {
			return cleanTag
		}
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
	if len(prefixRunes) < 15 {
		avail := 16 - len(prefixRunes) - 1
		branchRunes := []rune(branch)
		if avail < len(branchRunes) {
			return string(prefixRunes) + string(branchRunes[:avail]) + "…"
		}
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

func buildRepoRemediationOptions(rec RepoPendingCommitRecord) []RemediationOption {
	if !rec.IsDirty {
		return []RemediationOption{}
	}
	targetPath := rec.RelativePath
	if targetPath == "" {
		targetPath = rec.RepoName
	}
	return createRemediationOptionsPair(targetPath)
}

func createRemediationOptionsPair(targetPath string) []RemediationOption {
	opt1 := RemediationOption{
		OptionNumber: 1,
		Label:        "Commit & Push WIP",
		Command:      fmt.Sprintf(`git -C "%s" add -A && git commit -m "wip: save changes" && git push`, targetPath),
	}
	opt2 := RemediationOption{
		OptionNumber: 2,
		Label:        "Stash WIP",
		Command:      fmt.Sprintf(`git -C "%s" stash -u`, targetPath),
	}
	return []RemediationOption{opt1, opt2}
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

func ParsePorcelainStatusLines(output string) (int, int, int, []string) {
	var untracked, modified, staged int
	var pendingFiles []string

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		pendingFiles = append(pendingFiles, strings.TrimSpace(line))
		if strings.HasPrefix(line, "??") {
			untracked++
			continue
		}
		stagedDelta, modifiedDelta := classifyStatusChars(line)
		staged += stagedDelta
		modified += modifiedDelta
	}

	return untracked, modified, staged, pendingFiles
}

func classifyStatusChars(line string) (int, int) {
	if len(line) < 2 {
		return 0, 0
	}
	stagedDelta := 0
	modifiedDelta := 0
	idxChar := line[0]
	wtChar := line[1]
	if idxChar != ' ' && idxChar != '?' {
		stagedDelta = 1
	}
	if wtChar != ' ' && wtChar != '?' {
		modifiedDelta = 1
	}
	return stagedDelta, modifiedDelta
}

func parsePendingCommitLines(output string) []string {
	var shas []string
	lines := strings.Split(output, "\n")
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t != "" {
			shas = append(shas, t)
		}
	}
	return shas
}

func sortPendingCommits(records []RepoPendingCommitRecord, sortMode string) {
	switch sortMode {
	case "count":
		sort.SliceStable(records, func(i, j int) bool {
			countI := resolvePendingSortCount(records[i])
			countJ := resolvePendingSortCount(records[j])
			if countI != countJ {
				return countI > countJ
			}
			return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
		})
	case "name":
		sort.SliceStable(records, func(i, j int) bool {
			return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
		})
	default: // "priority"
		sort.SliceStable(records, func(i, j int) bool {
			pI := calculatePriorityRank(records[i])
			pJ := calculatePriorityRank(records[j])
			if pI != pJ {
				return pI < pJ
			}
			return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
		})
	}
}

func calculatePriorityRank(r RepoPendingCommitRecord) int {
	if r.IsDirty && r.HasUnpushed {
		return 1
	}
	if r.IsDirty {
		return 2
	}
	if r.HasUnpushed {
		return 3
	}
	return 4
}

func resolvePendingSortCount(r RepoPendingCommitRecord) int {
	total := r.TotalUncommitted
	if total == 0 {
		total = r.UntrackedFilesCount + r.ModifiedFilesCount + r.StagedFilesCount
	}
	return total + r.UnpushedCommitsCount
}

func emitPendingCommitsJSON(payload PendingCommitsPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal pending commits json")
	}
	fmt.Println(string(b))
	return nil
}

func renderPendingCommitsTerminal(payload PendingCommitsPayload, totalScanned int, allInspected []RepoPendingCommitRecord) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	renderTableBanner(border, reset)
	renderTableSummary(payload, border, reset)
	renderTableHeader(border, reset)
	renderTableRows(payload.Repositories, border, reset)
	renderCleanSummaryIfNeeded(payload, totalScanned, allInspected, border, reset)
	renderTableFooter(border, reset)
}

func renderTableBanner(border, reset string) {
	bold := constants.ColorBold
	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s                    GITMAP PENDING COMMITS SUMMARY                            %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableSummary(payload PendingCommitsPayload, border, reset string) {
	summaryText := fmt.Sprintf(" Scanned: %d repos │ Dirty: %d repos │ Uncommitted: %d files │ Unpushed: %d",
		payload.TotalReposScanned, payload.TotalDirtyRepos, payload.TotalUncommittedFiles, payload.TotalUnpushedCommits)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(summaryText, 78), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableHeader(border, reset string) {
	headerText := fmt.Sprintf(" %-22s  %-16s  %12s  %10s  %-9s",
		"REPOSITORY", "VER/BRANCH", "UNCOMMITTED", "UNPUSHED", "STATUS")
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(headerText, 78), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableRows(repos []RepoPendingCommitRecord, border, reset string) {
	for _, rec := range repos {
		renderPendingRepoRow(rec, border, reset)
		if rec.IsDirty {
			renderTreeRemediationHints(rec, border, reset)
		}
	}
}

func renderPendingRepoRow(rec RepoPendingCommitRecord, border, reset string) {
	repoDisplay := truncateStringWithEllipsis(rec.RepoName, 22)
	branchDisplay := rec.ShortVersionBranch
	if branchDisplay == "" {
		branchDisplay = rec.CurrentBranch
	}
	branchDisplay = truncateStringWithEllipsis(branchDisplay, 16)
	statusText := resolveStatusLabel(rec, reset)
	fmt.Printf("  %s│%s %-22s  %-16s  %12d  %10d  %s%s│%s\n",
		border, reset,
		padRight(repoDisplay, 22), padRight(branchDisplay, 16),
		rec.TotalUncommitted, rec.UnpushedCommitsCount,
		statusText, border, reset)
}

func resolveStatusLabel(rec RepoPendingCommitRecord, reset string) string {
	if rec.IsDirty {
		return constants.ColorYellow + "● PEND   " + reset
	}
	if rec.HasUnpushed {
		return constants.ColorYellow + "● PEND   " + reset
	}
	return constants.ColorGreen + "○ CLEAN  " + reset
}

func renderTreeRemediationHints(rec RepoPendingCommitRecord, border, reset string) {
	totalOpts := len(rec.RemediationOptions)
	for idx, opt := range rec.RemediationOptions {
		prefix := "├──"
		if idx == totalOpts-1 {
			prefix = "└──"
		}
		rawLine := fmt.Sprintf("   %s Option %d: %s", prefix, opt.OptionNumber, opt.Command)
		formattedLine := formatTreeLine(rawLine)
		fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedLine, border, reset)
	}
}

func formatTreeLine(raw string) string {
	runes := []rune(raw)
	if len(runes) >= 78 {
		return raw
	}
	return raw + strings.Repeat(" ", 78-len(runes))
}

func renderCleanSummaryIfNeeded(payload PendingCommitsPayload, totalScanned int, allInspected []RepoPendingCommitRecord, border, reset string) {
	cleanCount := countCleanRepositories(allInspected)
	if payload.DetailMode != "summary" {
		return
	}
	if cleanCount <= 0 {
		return
	}
	if len(payload.Repositories) >= totalScanned {
		return
	}
	renderCleanSummaryRow(cleanCount, border, reset)
}

func countCleanRepositories(records []RepoPendingCommitRecord) int {
	count := 0
	for _, rec := range records {
		if rec.IsClean {
			count++
		}
	}
	return count
}

func renderCleanSummaryRow(cleanCount int, border, reset string) {
	cleanLabel := fmt.Sprintf("%d clean repos", cleanCount)
	cleanStatus := constants.ColorGreen + "○ CLEAN  " + reset
	fmt.Printf("  %s│%s %-22s  %-16s  %12d  %10d  %s%s│%s\n",
		border, reset,
		padRight(cleanLabel, 22), padRight("-", 16),
		0, 0,
		cleanStatus, border, reset)
}

func renderTableFooter(border, reset string) {
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	footerMsg := ` Fleet Remediation: gitmap cpar "wip: save changes"`
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(footerMsg, 78), border, reset)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func padRight(s string, width int) string {
	rc := utf8.RuneCountInString(s)
	if rc >= width {
		return s
	}
	return s + strings.Repeat(" ", width-rc)
}

func renderPendingCommitsDetailsAll(records []RepoPendingCommitRecord) {
	fmt.Println()
	fmt.Printf("  %sDetailed Pending Breakdown:%s\n", constants.ColorBold, constants.ColorReset)
	for _, rec := range records {
		if rec.IsClean {
			continue
		}
		fmt.Printf("\n  %s● %s%s (%s)\n", constants.ColorCyan, rec.RepoName, constants.ColorReset, rec.CurrentBranch)
		if len(rec.PendingFiles) > 0 {
			fmt.Printf("    %sModified Files (%d):%s\n", constants.ColorYellow, len(rec.PendingFiles), constants.ColorReset)
			for _, f := range rec.PendingFiles {
				fmt.Printf("      %s\n", f)
			}
		}
		if len(rec.UnpushedCommitSHAs) > 0 {
			fmt.Printf("    %sUnpushed Commits (%d):%s\n", constants.ColorYellow, len(rec.UnpushedCommitSHAs), constants.ColorReset)
			for _, sha := range rec.UnpushedCommitSHAs {
				fmt.Printf("      %s\n", sha)
			}
		}
	}
}

func renderPendingCommitsInteractive(records []RepoPendingCommitRecord) {
	var dirtyList []RepoPendingCommitRecord
	for _, r := range records {
		if !r.IsClean {
			dirtyList = append(dirtyList, r)
		}
	}

	if len(dirtyList) == 0 {
		fmt.Println("  All repositories are clean. Nothing to step through.")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for idx, rec := range dirtyList {
		fmt.Printf("\n  [%d/%d] Repository: %s%s%s (branch: %s)\n",
			idx+1, len(dirtyList), constants.ColorCyan, rec.RepoName, constants.ColorReset, rec.CurrentBranch)
		if len(rec.PendingFiles) > 0 {
			fmt.Printf("    Files (%d):\n", len(rec.PendingFiles))
			for _, f := range rec.PendingFiles {
				fmt.Printf("      %s\n", f)
			}
		}
		if len(rec.UnpushedCommitSHAs) > 0 {
			fmt.Printf("    Unpushed Commits (%d):\n", len(rec.UnpushedCommitSHAs))
			for _, sha := range rec.UnpushedCommitSHAs {
				fmt.Printf("      %s\n", sha)
			}
		}

		if idx >= len(dirtyList)-1 {
			continue
		}
		if isDetailPromptTerminated(scanner) {
			break
		}
	}
}

func isDetailPromptTerminated(scanner *bufio.Scanner) bool {
	fmt.Printf("  Press [Enter] for next, or 'q' to quit: ")
	if !scanner.Scan() {
		return true
	}
	input := strings.TrimSpace(scanner.Text())
	return strings.EqualFold(input, "q") || strings.EqualFold(input, "quit")
}

func runFleetSSHPendingCommits(opts PendingCommitsOptions) error {
	conns, errFetch := cmdssh.FetchAllSSHConnections()
	if errFetch != nil {
		return apperror.WrapSimple(errFetch, "fetch fleet ssh connections")
	}

	if len(conns) == 0 {
		return handleEmptyFleet(opts.IsJSON)
	}

	results := make([]NodePendingCommitsRecord, len(conns))
	var wg sync.WaitGroup

	for i, conn := range conns {
		wg.Add(1)
		go func(idx int, c db.SSHConnection) {
			defer wg.Done()
			results[idx] = querySingleNodePendingCommits(c, opts)
		}(i, conn)
	}

	wg.Wait()

	onlineCount := 0
	for _, r := range results {
		if r.IsOnline {
			onlineCount++
		}
	}

	fleetPayload := NodesPendingCommitsPayload{
		Timestamp:    time.Now().UTC(),
		TotalNodes:   len(conns),
		OnlineNodes:  onlineCount,
		FleetResults: results,
	}

	if opts.IsJSON {
		return emitPendingCommitsFleetJSON(fleetPayload)
	}

	renderPendingCommitsFleetTerminal(fleetPayload)
	return nil
}

func handleEmptyFleet(isJSON bool) error {
	if isJSON {
		return emitPendingCommitsFleetJSON(NodesPendingCommitsPayload{
			Timestamp:    time.Now().UTC(),
			TotalNodes:   0,
			OnlineNodes:  0,
			FleetResults: []NodePendingCommitsRecord{},
		})
	}
	fmt.Printf("  %sNo registered SSH cluster nodes found.%s\n", constants.ColorYellow, constants.ColorReset)
	return nil
}

func querySingleNodePendingCommits(conn db.SSHConnection, opts PendingCommitsOptions) NodePendingCommitsRecord {
	start := time.Now()
	rec := NodePendingCommitsRecord{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OSType:    conn.OS,
	}

	isWin := strings.EqualFold(conn.OS, "windows")
	cmdStr := "gitmap pending-commits --json"
	if opts.SortMode != "" && opts.SortMode != "priority" {
		cmdStr += " --sort=" + opts.SortMode
	}
	if opts.DetailMode == "all" {
		cmdStr += " --detail=all"
	}
	if opts.IsAll {
		cmdStr += " --all"
	}

	shell := "bash"
	if isWin {
		cmdStr = fmt.Sprintf(`powershell.exe -NoProfile -Command "%s"`, cmdStr)
		shell = "ps"
	}

	out, errExec := currentPendingCommitsSSHExecutor.Execute(conn, cmdStr, shell)
	rec.LatencyMs = time.Since(start).Milliseconds()

	if errExec != nil {
		return handleSSHExecFailure(conn, errExec, rec.LatencyMs)
	}

	rec.IsOnline = true
	rec.IsSuccess = true

	cleanJSON := extractJSONPayload(out)
	var nodePayload PendingCommitsPayload
	if errUnmarshal := json.Unmarshal([]byte(cleanJSON), &nodePayload); errUnmarshal == nil {
		rec.Payload = &nodePayload
	}

	return rec
}

func handleSSHExecFailure(conn db.SSHConnection, errExec error, latencyMs int64) NodePendingCommitsRecord {
	isOffline := isOfflineErrorString(errExec.Error())
	return NodePendingCommitsRecord{
		NodeAlias:    conn.Alias,
		Host:         conn.IPAddress,
		OSType:       conn.OS,
		IsOnline:     !isOffline,
		IsSuccess:    false,
		LatencyMs:    latencyMs,
		ErrorMessage: errExec.Error(),
	}
}

func isOfflineErrorString(errStr string) bool {
	low := strings.ToLower(errStr)
	return strings.Contains(low, "network unreachable") ||
		strings.Contains(low, "connection refused") ||
		strings.Contains(low, "timeout") ||
		strings.Contains(low, "i/o timeout") ||
		strings.Contains(low, "offline")
}

func extractJSONPayload(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

func emitPendingCommitsFleetJSON(payload NodesPendingCommitsPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal fleet pending commits json")
	}
	fmt.Println(string(b))
	return nil
}

func renderPendingCommitsFleetTerminal(fleet NodesPendingCommitsPayload) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold

	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s               GITMAP FLEET PENDING COMMITS (SSH CLUSTER)                     %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s Total Nodes: %-6d │ Online Nodes: %-6d                                   %s│%s\n",
		border, reset, fleet.TotalNodes, fleet.OnlineNodes, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s %-16s %-15s %-7s %-7s %-8s %-10s %7s %s│%s\n",
		border, reset,
		"NODE", "HOST", "STATUS", "DIRTY", "UNCOMMIT", "UNPUSHED", "LATENCY",
		border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)

	for _, node := range fleet.FleetResults {
		statusLabel, dirtyStr, uncommitStr, unpushedStr := formatFleetNodeRow(node, reset)

		fmt.Printf("  %s│%s %-16s %-15s %-7s %-7s %-8s %-10s %5dms %s│%s\n",
			border, reset,
			node.NodeAlias,
			node.Host,
			statusLabel,
			dirtyStr,
			uncommitStr,
			unpushedStr,
			node.LatencyMs,
			border, reset)
	}

	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func formatFleetNodeRow(node NodePendingCommitsRecord, reset string) (string, string, string, string) {
	if !node.IsOnline {
		return constants.ColorRed + "OFFLINE" + reset, "-", "-", "-"
	}

	if node.IsSuccess && node.Payload != nil {
		statusLabel := constants.ColorGreen + "ONLINE " + reset
		dirtyStr := fmt.Sprintf("%d", node.Payload.TotalDirtyRepos)
		uncommitStr := fmt.Sprintf("%d", node.Payload.TotalUncommittedFiles)
		unpushedStr := fmt.Sprintf("%d", node.Payload.TotalUnpushedCommits)
		return statusLabel, dirtyStr, uncommitStr, unpushedStr
	}

	return constants.ColorYellow + "ERROR  " + reset, "-", "-", "-"
}

func printPendingCommitsHelp() {
	cyan := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold
	yellow := constants.ColorYellow

	fmt.Printf(`
%s┌──────────────────────────────────────────────────────────────────────────────┐
│                      GITMAP PENDING COMMITS HELP                             │
└──────────────────────────────────────────────────────────────────────────────┘%s

%sDESCRIPTION:%s
  Scans local workspace repositories or fleet cluster nodes over SSH to identify
  uncommitted working tree modifications and unpushed commits ahead of origin.

%sUSAGE:%s
  gitmap pending-commits [repoName|all] [flags]
  gitmap pc [repoName|all] [flags]

%sFLAGS:%s
  -s, --sort <mode>       Sort order: priority | name | count (default: priority)
  -d, --detail <mode>     Detail view: summary | all | 1 | none (default: summary)
      --ssh, -ssh         Aggregate pending commits across all registered cluster SSH nodes
  -j, --json              Output machine-readable JSON telemetry
      --dirty-only        Display only repositories with changes (default: true)
      --all               Include completely clean repositories in output
      --no-cache          Bypass status cache and perform live inspection
      --refresh           Force refresh cached pending commit status
      --backup            Serve status from secondary persistent backup DB
      --serve-backup      Serve status from secondary persistent backup DB (alias)
  -h, --help              Display this help menu

%sEXAMPLES:%s
  %sgitmap pc%s
      Summary overview of all dirty repositories in workspace.

  %sgitmap pc --sort=count --detail=all%s
      Detailed file diff breakdown sorted by total modification count.

  %sgitmap pc -ssh --json%s
      Aggregate pending commits across entire SSH fleet as JSON.
`,
		cyan, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		yellow, reset,
		yellow, reset,
		yellow, reset,
	)
}
