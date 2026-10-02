package cmdpull

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

// pullOptions holds parsed pull flags.
type pullOptions struct {
	slug          string
	group         string
	all           bool
	verbose       bool
	stopOnFail    bool
	parallel      int
	workers       int
	hands         int
	isWWOH        bool
	onlyAvailable bool
	autoFix       bool
	yes           bool
	noFix         bool
	isRaw         bool
	useSSH        bool
	useHTTPS      bool
	showStatus    bool
	isJSON        bool
	isProbe       bool
}

// NormalizePullArgs converts positional pull-all and table arguments into flags.
func NormalizePullArgs(args []string) []string {
	normalized := make([]string, 0, len(args)+1)
	for i := 0; i < len(args); i++ {
		if isPullAllTableSeq(args, i) {
			normalized = append(normalized, "--all", "--status")
			i++
			continue
		}
		if isHandShortFlagWithVal(args, i) {
			normalized = append(normalized, "--hand", args[i+1])
			i++
			continue
		}
		normalized = appendNormalizedPullToken(normalized, args[i])
	}

	return reorderPullFlags(normalized)
}

func isHandShortFlagWithVal(args []string, i int) bool {
	if args[i] != "-h" || i+1 >= len(args) {
		return false
	}
	_, err := strconv.Atoi(args[i+1])
	return err == nil
}

func reorderPullFlags(args []string) []string {
	flags := make([]string, 0, len(args))
	pos := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isPullFlagTakingValue(arg) {
			flags = appendValuedFlag(flags, args, &i)
			continue
		}
		flags, pos = categorizePullArg(arg, flags, pos)
	}

	return append(flags, pos...)
}

func categorizePullArg(arg string, flags, pos []string) ([]string, []string) {
	if strings.HasPrefix(arg, "-") {
		return append(flags, arg), pos
	}

	return flags, append(pos, arg)
}

func appendValuedFlag(flags []string, args []string, idx *int) []string {
	flags = append(flags, args[*idx])
	if *idx+1 < len(args) {
		*idx++
		flags = append(flags, args[*idx])
	}
	return flags
}

func isPullFlagTakingValue(arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}

	return arg == "-g" || arg == "--group" || arg == "-p" || arg == "--parallel" ||
		arg == "-w" || arg == "--w" || arg == "--worker" || arg == "--workers" ||
		arg == "--hand" || arg == "--hands" || arg == "--h"
}

func isPullAllTableSeq(args []string, i int) bool {
	isAllToken := args[i] == "all" || args[i] == "--all" || args[i] == "pa" || args[i] == "ta" || args[i] == "pull-all"

	return isAllToken && i+1 < len(args) && args[i+1] == "table"
}

func appendNormalizedPullToken(normalized []string, token string) []string {
	if isTableToken(token) {
		return append(normalized, "--all", "--status")
	}
	if isSSHFleetToken(token) {
		return append(normalized, "--all", "--ssh")
	}
	if isAllToken(token) {
		return append(normalized, "--all")
	}

	return append(normalized, token)
}

func isTableToken(token string) bool {
	lower := strings.ToLower(token)
	return lower == "pat" || lower == "pull-all-table"
}

func isAllToken(token string) bool {
	lower := strings.ToLower(token)
	return lower == "all" || lower == "pa" || lower == "ta" || lower == "pull-all"
}

// runPull handles the "pull" subcommand.
func runPull(args []string) error {
	args = NormalizePullArgs(args)
	isPullAll := isPullAllInvocation(args)
	checkPullHelp(isPullAll, args)
	if handled, err := checkEfficientSubcommand(args); handled {
		return err
	}
	if isPullAll && hasSSHFleetFlag(args) {
		return handleSSHFleetPullAll(args)
	}

	return runPullStandardFlow(args, isPullAll)
}

func runPullStandardFlow(args []string, isPullAll bool) error {
	args = NormalizePullArgs(args)
	printHeaderUnlessJSON(isPullAll, args)
	requireOnline()
	if db, err := openDB(); err == nil {
		_, _ = OptimizeRedundantRepos(db, hasJSONArg(args))
		db.Close()
	}
	useSSH, useHTTPS, restArgs := ExtractTransportFlags(args)
	if isCWDTransportEligible(isPullAll, useSSH, useHTTPS) {
		return runPullCWDWithTransport(useSSH, useHTTPS, restArgs)
	}
	opts := resolveParsedPullOptions(restArgs, useSSH, useHTTPS)

	return executePullWithResolvedOptions(opts)
}

func printHeaderUnlessJSON(isPullAll bool, args []string) {
	if !hasJSONArg(args) {
		printPullInvocationHeader(isPullAll)
	}
}

func isCWDTransportEligible(isPullAll, useSSH, useHTTPS bool) bool {
	if isPullAll {
		return false
	}
	hasTransport := useSSH || useHTTPS
	if !hasTransport {
		return false
	}

	return isGitRepoCWD()
}

func checkPullHelp(isPullAll bool, args []string) {
	helpCmd := "pull"
	if isPullAll {
		helpCmd = "pull-all"
	}
	checkHelp(helpCmd, args)
}

func handleSSHFleetPullAll(args []string) error {
	cleanArgs := stripSSHFleetFlags(args)
	hasHook := RunRemoteSSHPullAllFleetFn != nil
	if !hasHook {
		return apperror.NewSimple("ssh fleet pull-all is not wired", "E_SSH_PULLALL")
	}

	return executeSSHFleetPullTask(cleanArgs)
}

func executeSSHFleetPullTask(cleanArgs []string) error {
	queueId, tDB := enqueueTaskQueue("pull-all-ssh", "fleet")
	defer closeTaskDB(tDB)
	updateTaskQueue(tDB, queueId, "running")

	errFleet := RunRemoteSSHPullAllFleetFn(cleanArgs)
	hasErr := errFleet != nil
	if hasErr {
		return handleSSHFleetFailure(tDB, queueId, errFleet)
	}

	updateTaskQueue(tDB, queueId, "completed")
	return nil
}

func handleSSHFleetFailure(tDB *store.TasksSplitDB, queueId string, errFleet error) error {
	store.LogInternalError("PULL_ALL_SSH", "SSH_DELEGATION_ERROR", errFleet.Error(), "", "")
	updateTaskQueue(tDB, queueId, "failed")

	return apperror.WrapSimple(errFleet, "SSH fleet pull-all failed")
}

func hasSSHFleetFlag(args []string) bool {
	for _, a := range args {
		if isSSHFleetToken(a) {
			return true
		}
	}
	return false
}

func isSSHFleetToken(token string) bool {
	low := strings.ToLower(token)
	return low == "--ssh" || low == "-ssh" || low == "--sh" || low == "-sh" || low == "ssh" || low == "pas" || low == "pull-all-ssh"
}

func stripSSHFleetFlags(args []string) []string {
	var clean []string
	for _, a := range args {
		if isSSHFleetToken(a) {
			continue
		}
		clean = append(clean, a)
	}
	return clean
}

func executePullWithResolvedOptions(opts pullOptions) error {
	if opts.isProbe {
		executePullProbeHook(opts)
	}
	if opts.verbose {
		initVerboseLog()
	}

	return dispatchPullExecution(opts)
}

func executePullProbeHook(opts pullOptions) {
	if RunSpecialRepoProbeOnPullFn == nil {
		return
	}
	if !opts.isJSON {
		fmt.Printf("\n  %s●%s %sProbing companion repositories (repo-secrets & repo-cache)...%s\n",
			constants.ColorCyan, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	}
	_ = RunSpecialRepoProbeOnPullFn("", opts.yes)
}

func hasJSONArg(args []string) bool {
	for _, a := range args {
		if strings.EqualFold(a, "--json") || strings.EqualFold(a, "-json") {
			return true
		}
	}

	return false
}

func resolveParsedPullOptions(args []string, useSSH, useHTTPS bool) pullOptions {
	opts := parsePullFlags(args)
	opts = applyTransportOptions(opts, useSSH, useHTTPS)
	if isPullAllTableRootCmd() {
		opts.all = true
		opts.showStatus = true
	}

	return opts
}

func applyTransportOptions(opts pullOptions, useSSH, useHTTPS bool) pullOptions {
	if useSSH {
		opts.useSSH = true
	}
	if useHTTPS {
		opts.useHTTPS = true
	}

	return opts
}

func isPullAllInvocation(args []string) bool {
	if isPullAllRootCmd() {
		return true
	}

	return hasPullAllArg(args)
}

func checkEfficientSubcommand(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if isEfficientTableSubcmd(args[0]) {
		return true, RunPullAllEfficient(args[1:], true, "pull "+args[0], isShortEfficientSubcmd(args[0]))
	}
	if isEfficientPullSubcmd(args[0]) {
		return true, RunPullAllEfficient(args[1:], false, "pull "+args[0], isShortEfficientSubcmd(args[0]))
	}

	return false, nil
}

func isEfficientTableSubcmd(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "all-efficient-table" || lower == "paet" || lower == "aet"
}

func isEfficientPullSubcmd(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "all-efficient" || lower == "ae" || lower == "pae" || lower == "pull-ae" || lower == "efficient"
}

func isShortEfficientSubcmd(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "ae" || lower == "pae" || lower == "pull-ae" || lower == "paet" || lower == "aet"
}

func isPullAllRootCmd() bool {
	if len(os.Args) <= 1 {
		return false
	}
	first := strings.ToLower(os.Args[1])

	return isPullAllToken(first)
}

func isPullAllTableRootCmd() bool {
	if len(os.Args) <= 1 {
		return false
	}
	first := strings.ToLower(os.Args[1])

	return first == "pull-all-table" || first == "pat"
}

func hasPullAllArg(args []string) bool {
	for _, a := range args {
		if isPullAllToken(strings.ToLower(a)) {
			return true
		}
	}

	return false
}

func isPullAllToken(token string) bool {
	return token == "--all" || token == "-all" || token == "-a" || token == "all" || token == "pa" || token == "ta" || token == "pull-all" || token == "pat" || token == "pull-all-table" || token == "pas" || token == "pull-all-ssh"
}

func printPullInvocationHeader(isPullAll bool) {
	cwd, _ := os.Getwd()
	cmdName := "pull"
	if isPullAll {
		cmdName = "pull-all"
	}
	fmt.Printf("\n  %s→%s %sgitmap %s%s %s(cwd: %s)%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, cmdName, constants.ColorReset,
		constants.ColorDim, cwd, constants.ColorReset)
}

func resolveSubArrow() string {
	if glyphs.Resolve() == glyphs.ModeSafe {
		return "->"
	}

	return "→"
}

func dispatchPullExecution(opts pullOptions) error {
	if isPullCWDEnabled(opts) {
		return executePullCWDWithNotice(opts.isRaw)
	}
	if ShouldFallbackToPullAll(opts) {
		announceNonGitFallbackUnlessJSON(opts.isJSON)
		opts.all = true
	}

	return runPullBatch(opts)
}

func executePullCWDWithNotice(isRaw bool) error {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %scwd is a git repo — running plain `git pull` here%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)

	return runPullCWD(isRaw)
}

func announceNonGitFallbackUnlessJSON(isJSON bool) {
	if !isJSON {
		AnnounceNonGitPullFallback()
	}
}

func handleEmptyBatchRecords(isJSON bool) error {
	if isJSON {
		return renderPullBatchJSONSummary(0, nil, 0)
	}

	return nil
}

func runPullBatch(opts pullOptions) error {
	records, isFound := resolvePullBatchRecords(opts)
	hasRecords := isFound && len(records) > 0
	if !hasRecords {
		return handleEmptyBatchRecords(opts.isJSON)
	}
	printResolvedPullReposUnlessJSON(len(records), opts.isJSON)
	filtered := applyPullAvailableFilter(records, opts.onlyAvailable)
	if isNoAvailablePullTargets(opts.onlyAvailable, len(filtered)) {
		return handleNoAvailablePullTargets(opts)
	}

	return executePullBatchLifecycle(filtered, opts)
}

func printResolvedPullReposUnlessJSON(count int, isJSON bool) {
	if !isJSON {
		printResolvedPullRepos(count)
	}
}

func isNoAvailablePullTargets(onlyAvailable bool, count int) bool {
	return onlyAvailable && count == 0
}

func handleNoAvailablePullTargets(opts pullOptions) error {
	if opts.isJSON {
		return renderPullBatchJSONSummary(0, nil, 0)
	}
	fmt.Print(constants.MsgPullNoAvailable)

	return nil
}

func printResolvedPullRepos(count int) {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s resolved %s%d%s repo(s) to pull\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorBold, count, constants.ColorReset)
}

func applyPullAvailableFilter(records []model.ScanRecord, isOnlyAvailable bool) []model.ScanRecord {
	if isOnlyAvailable {
		return filterByAvailableUpdates(records)
	}

	return records
}

func executePullBatchLifecycle(records []model.ScanRecord, opts pullOptions) error {
	taskID, taskDB := beginPullTask(records)
	queueId, tDB := initLocalPullTaskQueue(opts.all)
	defer closePullBatchDBs(taskDB, tDB)

	applySelectiveTransport(records, opts)
	ignoreHandle := startAsyncIgnoreScan(records)
	bar, states, dur := runPullBatchWork(records, opts)
	if opts.isJSON {
		return finalizePullBatchJSON(taskDB, taskID, tDB, queueId, opts.all, len(records), states, dur)
	}

	return finalizePullBatchStandard(taskDB, taskID, tDB, queueId, opts.all, records, states, dur, opts, bar.Failed(), ignoreHandle)
}

func closePullBatchDBs(taskDB *store.DB, tDB *store.TasksSplitDB) {
	if taskDB != nil {
		_ = taskDB.Close()
	}
	closeTaskDB(tDB)
}

func initLocalPullTaskQueue(isAll bool) (string, *store.TasksSplitDB) {
	if !isAll {
		return "", nil
	}
	queueId, tDB := enqueueTaskQueue("pull-all", "local")
	updateTaskQueue(tDB, queueId, "running")

	return queueId, tDB
}

func applySelectiveTransport(records []model.ScanRecord, opts pullOptions) {
	if !opts.all {
		maybeApplyTransportToRecords(records, opts.useSSH, opts.useHTTPS)
	}
}

func runPullBatchWork(records []model.ScanRecord, opts pullOptions) (*PullProgressBar, []*PullRepoState, time.Duration) {
	bar, sortedStates, dur := runPullBatchExecution(records, opts)
	syncPullBatchTelemetry(records, sortedStates, dur, opts)

	return bar, sortedStates, dur
}

func finalizePullBatchJSON(taskDB *store.DB, taskID int64, tDB *store.TasksSplitDB, queueId string, isAll bool, total int, states []*PullRepoState, dur time.Duration) error {
	completePendingTask(taskDB, taskID)
	finalizePullTaskQueueJSON(isAll, tDB, queueId)

	return renderPullBatchJSONSummary(total, states, dur)
}

func finalizePullBatchStandard(taskDB *store.DB, taskID int64, tDB *store.TasksSplitDB, queueId string, isAll bool, records []model.ScanRecord, states []*PullRepoState, dur time.Duration, opts pullOptions, failedCount int, ignoreHandle *IgnoreScanHandle) error {
	renderPullBatchOutput(records, states, dur, opts)
	collectAndRemediateIgnoreIssues(ignoreHandle, opts)
	finalizePullTaskQueueOutput(isAll, tDB, queueId, failedCount)

	return finalizePullBatchTask(taskDB, taskID, failedCount)
}

func finalizePullTaskQueueJSON(isAll bool, tDB *store.TasksSplitDB, queueId string) {
	if isAll {
		updateTaskQueue(tDB, queueId, "completed")
	}
}

func finalizePullTaskQueueOutput(isAll bool, tDB *store.TasksSplitDB, queueId string, failedCount int) {
	if !isAll {
		return
	}
	status := "completed"
	if failedCount > 0 {
		status = "failed"
	}
	updateTaskQueue(tDB, queueId, status)
}

// IgnoreRepoIssue holds ignore configuration and index issues for a repository.
type IgnoreRepoIssue struct {
	RepoName          string
	RepoPath          string
	DuplicatePatterns []string
	TrackedPaths      []string
}

// HasIssues reports whether any ignore or tracking issues were detected.
func (i IgnoreRepoIssue) HasIssues() bool {
	return len(i.DuplicatePatterns) > 0 || len(i.TrackedPaths) > 0
}

// IgnoreScanHandle manages asynchronous ignore scanning across repositories.
type IgnoreScanHandle struct {
	done chan []IgnoreRepoIssue
}

func startAsyncIgnoreScan(records []model.ScanRecord) *IgnoreScanHandle {
	return StartThrottledAsyncIgnoreScan(records, resolvePullIgnoreTTL())
}

func resolvePullIgnoreTTL() time.Duration {
	s, err := store.OpenDefault()
	if err != nil {
		return 24 * time.Hour
	}
	defer s.Close()
	return config.GetGitIgnoreTTL(s)
}

func (h *IgnoreScanHandle) Collect() []IgnoreRepoIssue {
	if h == nil {
		return nil
	}
	return <-h.done
}

// DeduplicateIgnoreIssues eliminates duplicate issues by canonical repo path.
func DeduplicateIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue {
	if len(issues) <= 1 {
		return issues
	}
	seen := make(map[string]bool, len(issues))
	unique := make([]IgnoreRepoIssue, 0, len(issues))
	for _, issue := range issues {
		unique = appendUniqueIgnoreIssue(unique, issue, seen)
	}

	return unique
}

func appendUniqueIgnoreIssue(unique []IgnoreRepoIssue, issue IgnoreRepoIssue, seen map[string]bool) []IgnoreRepoIssue {
	key := resolveIgnoreIssueKey(issue)
	if seen[key] {
		return unique
	}
	seen[key] = true

	return append(unique, issue)
}

func resolveIgnoreIssueKey(issue IgnoreRepoIssue) string {
	key := CanonicalRepoPathKey(issue.RepoPath)
	if key != "" {
		return key
	}

	return strings.ToLower(strings.TrimSpace(issue.RepoName))
}

func sortIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue {
	deduped := DeduplicateIgnoreIssues(issues)
	sorted := make([]IgnoreRepoIssue, len(deduped))
	copy(sorted, deduped)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].RepoName) < strings.ToLower(sorted[j].RepoName)
	})

	return sorted
}

func inspectRepoForIgnoreIssues(repoDir, repoName string) IgnoreRepoIssue {
	issue := IgnoreRepoIssue{RepoPath: repoDir, RepoName: repoName}
	if !gitignoreagm.IsGitRepository(repoDir) {
		return issue
	}
	issue.DuplicatePatterns = findGitignoreDuplicatePatterns(repoDir)
	issue.TrackedPaths = findTrackedDefaultIgnorePaths(repoDir)
	return issue
}

func findGitignoreDuplicatePatterns(repoDir string) []string {
	data, err := os.ReadFile(filepath.Join(repoDir, ".gitignore"))
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	seen := make(map[string]bool, len(lines))
	dupSeen := make(map[string]bool)
	var duplicates []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if isIgnorablePatternLine(trimmed) {
			continue
		}
		norm := strings.TrimPrefix(trimmed, "/")
		if !seen[norm] {
			seen[norm] = true
			continue
		}
		if !dupSeen[norm] {
			dupSeen[norm] = true
			duplicates = append(duplicates, trimmed)
		}
	}
	return duplicates
}

func isIgnorablePatternLine(trimmed string) bool {
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}

func findTrackedDefaultIgnorePaths(repoDir string) []string {
	paths := buildDefaultTrackedPathsToCheck(repoDir)
	var tracked []string
	for _, p := range paths {
		if isPathTrackedInGitIndex(repoDir, p) {
			tracked = append(tracked, p)
		}
	}
	return tracked
}

func buildDefaultTrackedPathsToCheck(repoDir string) []string {
	return []string{
		".gitmap/backup/",
		"antigravity-resume_task.json",
		".antigravity_resume_task.json",
		"antigravity_resume_task.json",
		".antigravity-resume_task.json",
	}
}

func isPathTrackedInGitIndex(repoDir, pathspec string) bool {
	return gitutil.ExecGitCheck(5*time.Second, repoDir, "-C", repoDir, "ls-files", "--error-unmatch", pathspec)
}

func collectAndRemediateIgnoreIssues(handle *IgnoreScanHandle, opts pullOptions) {
	if handle == nil {
		return
	}
	issues := handle.Collect()
	if len(issues) == 0 {
		return
	}
	deduped := DeduplicateIgnoreIssues(issues)
	isAutoYes := opts.yes || opts.autoFix
	handleIgnoreRemediation(deduped, isAutoYes, opts.isJSON)
}

func handleIgnoreRemediation(issues []IgnoreRepoIssue, isAutoYes, isJSON bool) {
	if len(issues) == 0 || isJSON {
		return
	}
	printIgnoreIssuesReport(issues)
	if isAutoYes {
		remediateAllIgnoreIssues(issues)
		return
	}
	if !isInteractiveTerminal() {
		markIssuesSkipped(issues)
		printNonInteractiveIgnoreNotice(len(issues))
		return
	}
	dispatchInteractiveIgnoreRemediation(issues)
}

func printIgnoreIssuesReport(issues []IgnoreRepoIssue) {
	fmt.Printf("\n  %s⚠%s %sDetected .gitignore issues in %d repository(ies):%s\n",
		constants.ColorYellow, constants.ColorReset,
		constants.ColorBold, len(issues), constants.ColorReset)
	for _, issue := range issues {
		printSingleRepoIgnoreIssue(issue)
	}
	fmt.Println()
}

func printSingleRepoIgnoreIssue(issue IgnoreRepoIssue) {
	fmt.Printf("    %s• %s%s\n", constants.ColorCyan, issue.RepoName, constants.ColorReset)
	for _, dup := range issue.DuplicatePatterns {
		fmt.Printf("      %s-%s Duplicate pattern in .gitignore: %s%s%s\n",
			constants.ColorDim, constants.ColorReset, constants.ColorYellow, dup, constants.ColorReset)
	}
	for _, tracked := range issue.TrackedPaths {
		fmt.Printf("      %s-%s Tracked in git index (should be ignored): %s%s%s\n",
			constants.ColorDim, constants.ColorReset, constants.ColorRed, tracked, constants.ColorReset)
	}
}

func isInteractiveTerminal() bool {
	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") != "" {
		return false
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func printNonInteractiveIgnoreNotice(count int) {
	fmt.Printf("  %sℹ Detected ignore issues in %d repository(ies). Run 'gitmap fix-ignore-all' to remediate.%s\n\n",
		constants.ColorCyan, count, constants.ColorReset)
}

func dispatchInteractiveIgnoreRemediation(issues []IgnoreRepoIssue) {
	choice := promptIgnoreRemediationChoice()
	if isChoiceAll(choice) {
		remediateAllIgnoreIssues(issues)
		return
	}
	if isChoiceSingle(choice) {
		remediateSingleRepoInteractive(issues)
		return
	}
	markIssuesSkipped(issues)
	fmt.Printf("  %s↷ Skipped ignore resolution.%s\n\n", constants.ColorDim, constants.ColorReset)
}

func markIssuesSkipped(issues []IgnoreRepoIssue) {
	for _, issue := range issues {
		_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
	}
}

func promptIgnoreRemediationChoice() string {
	fmt.Printf("  %s?%s Resolve all at once [y/all], one-by-one [s/single], or skip [n]?: ",
		constants.ColorCyan, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "n"
	}
	return strings.ToLower(strings.TrimSpace(line))
}

func isChoiceAll(choice string) bool {
	return choice == "y" || choice == "yes" || choice == "all" || choice == "a"
}

func isChoiceSingle(choice string) bool {
	return choice == "s" || choice == "single" || choice == "one" || choice == "1"
}

func remediateSingleRepoInteractive(issues []IgnoreRepoIssue) {
	reader := bufio.NewReader(os.Stdin)
	remediatedCount := 0
	for _, issue := range issues {
		if !promptSingleRepoRemediation(reader, issue.RepoName) {
			_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
			continue
		}
		if remediateSingleRepoIgnore(issue) {
			remediatedCount++
		}
	}
	fmt.Printf("\n  %s✓ Completed remediation across %d repository(ies).%s\n\n",
		constants.ColorGreen, remediatedCount, constants.ColorReset)
}

func promptSingleRepoRemediation(reader *bufio.Reader, repoName string) bool {
	fmt.Printf("  %s?%s [%s] Resolve ignore issues in %s? [Y/n]: ",
		constants.ColorCyan, constants.ColorReset, repoName, repoName)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "" || ans == "y" || ans == "yes"
}

func remediateAllIgnoreIssues(issues []IgnoreRepoIssue) {
	remediatedCount := 0
	for _, issue := range issues {
		if remediateSingleRepoIgnore(issue) {
			remediatedCount++
		}
	}
	fmt.Printf("\n  %s✓ Completed ignore remediation across %d repository(ies).%s\n\n",
		constants.ColorGreen, remediatedCount, constants.ColorReset)
}

func remediateSingleRepoIgnore(issue IgnoreRepoIssue) bool {
	wasUntracked := untrackRepoPaths(issue.RepoPath, issue.TrackedPaths)
	wasSanitized := sanitizeRepoGitignore(issue.RepoPath)
	if wasUntracked || wasSanitized {
		fmt.Printf("  %s✓%s [%s] Resolved ignore issues and sanitized .gitignore\n",
			constants.ColorGreen, constants.ColorReset, issue.RepoName)
		_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "clean", 1, 0)
		return true
	}
	_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
	return false
}

func untrackRepoPaths(repoDir string, trackedPaths []string) bool {
	if len(trackedPaths) == 0 {
		return false
	}
	args := append([]string{"-C", repoDir, "rm", "--cached", "-r", "-f", "--ignore-unmatch", "--"}, trackedPaths...)
	_, err := gitutil.ExecGitWithTimeout(10*time.Second, repoDir, args...)
	if err != nil {
		return false
	}
	commitArgs := []string{"-C", repoDir, "commit", "-m", "chore(git): untrack ignored files from index"}
	_, _ = gitutil.ExecGitWithTimeout(10*time.Second, repoDir, commitArgs...)
	return true
}

func sanitizeRepoGitignore(repoDir string) bool {
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data))
	if !isModified {
		return false
	}
	if writeErr := os.WriteFile(ignorePath, []byte(cleaned), 0o644); writeErr != nil {
		return false
	}
	return commitSanitizedGitignore(repoDir)
}

func commitSanitizedGitignore(repoDir string) bool {
	addArgs := []string{"-C", repoDir, "add", ".gitignore"}
	_, errAdd := gitutil.ExecGitWithTimeout(10*time.Second, repoDir, addArgs...)
	if errAdd != nil {
		return false
	}
	commitArgs := []string{"-C", repoDir, "commit", "-m", "chore(git): sanitize .gitignore patterns"}
	_, errCommit := gitutil.ExecGitWithTimeout(10*time.Second, repoDir, commitArgs...)
	return errCommit == nil
}

func runPullBatchExecution(records []model.ScanRecord, opts pullOptions) (*PullProgressBar, []*PullRepoState, time.Duration) {
	bar := NewPullProgressBar(len(records), opts.isJSON, opts.stopOnFail)
	if !opts.isJSON {
		bar.Start()
	}
	startTime := time.Now()
	executePull(records, bar, opts)
	if !opts.isJSON {
		bar.Stop()
	}
	dur := time.Since(startTime)
	dedupedStates := DeduplicateRepoStates(bar.States())

	return bar, sortStatesAlphabetically(dedupedStates), dur
}

func renderPullBatchOutput(records []model.ScanRecord, sortedStates []*PullRepoState, dur time.Duration, opts pullOptions) {
	if isConcisePullOutput(opts.all, opts.showStatus) {
		renderConciseActiveResults(sortedStates, records)
		activeCount, upToDateCount := countActiveStates(sortedStates)
		printPullAllFastSummary(len(sortedStates), activeCount, upToDateCount, dur)
	} else {
		renderPullBatchResults(sortedStates)
	}
	handlePullRemediationForRecords(records, opts)
}

func isConcisePullOutput(isAll, isShowStatus bool) bool {
	if isShowStatus {
		return false
	}

	return isAll
}

func countActiveStates(states []*PullRepoState) (int, int) {
	activeCount := 0
	upToDateCount := 0
	for _, s := range states {
		if ResolveRepoStatusLabel(s.Changes) == "up-to-date" {
			upToDateCount++
		} else {
			activeCount++
		}
	}

	return activeCount, upToDateCount
}

func printPullAllFastSummary(pulledCount, activeCount, upToDateCount int, dur time.Duration) {
	if activeCount > 0 {
		fmt.Printf("\n  %s✓%s %sPull all complete:%s %d pulled (%d active, %d up-to-date) (%s)\n\n",
			constants.ColorGreen, constants.ColorReset,
			constants.ColorBold, constants.ColorReset,
			pulledCount, activeCount, upToDateCount, FormatPullDuration(dur))

		return
	}
	fmt.Printf("\n  %s✓%s %sPull all complete:%s %d pulled (%s)\n\n",
		constants.ColorGreen, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		pulledCount, FormatPullDuration(dur))
}

func syncPullBatchTelemetry(records []model.ScanRecord, states []*PullRepoState, dur time.Duration, opts pullOptions) {
	telemetry := buildPullSessionTelemetry(records, states, dur, opts.all)
	_ = RecordPullBatchSession(telemetry, states)
}

func buildPullSessionTelemetry(records []model.ScanRecord, states []*PullRepoState, dur time.Duration, isAll bool) PullSessionTelemetry {
	cwd, _ := os.Getwd()
	succ, fail := countBatchStateOutcomes(states)

	return PullSessionTelemetry{
		CommandType: resolveBatchCommandType(isAll), WorkingDir: cwd,
		TotalRepos: len(records), PulledRepos: len(states),
		SkippedRepos: len(records) - len(states), SuccessCount: succ, FailedCount: fail,
		IsEfficient: false, Duration: dur, GitMapVersion: constants.Version,
	}
}

func resolveBatchCommandType(isAll bool) string {
	if isAll {
		return "pull-all"
	}

	return "pull"
}

func countBatchStateOutcomes(states []*PullRepoState) (int, int) {
	successCount := 0
	failedCount := 0
	for _, s := range states {
		if s.ErrorMsg == "" && s.Step != PullStepTypeError {
			successCount++
		} else {
			failedCount++
		}
	}

	return successCount, failedCount
}

func maybeApplyTransportToRecords(records []model.ScanRecord, useSSH, useHTTPS bool) {
	hasTransport := useSSH || useHTTPS
	if !hasTransport {
		return
	}
	for _, rec := range records {
		applyTransportIfPathExists(rec.AbsolutePath, useSSH, useHTTPS)
	}
}

func applyTransportIfPathExists(path string, useSSH, useHTTPS bool) {
	hasPath := len(path) > 0
	if !hasPath {
		return
	}

	ApplyTransportFlag(path, useSSH, useHTTPS)
}

func sortStatesAlphabetically(states []*PullRepoState) []*PullRepoState {
	sorted := make([]*PullRepoState, len(states))
	copy(sorted, states)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].RepoName) < strings.ToLower(sorted[j].RepoName)
	})

	return sorted
}

func renderPullBatchResults(states []*PullRepoState) {
	var tableRows []model.PullTableRow
	for _, state := range states {
		tableRows = append(tableRows, buildPullTableRowFromState(state))
	}
	RenderPullBatchTable(tableRows)
}

func buildPullTableRowFromState(state *PullRepoState) model.PullTableRow {
	row := initBasePullTableRow(state)
	row.LatestBranch = gitutil.GetLatestRemoteBranch(state.RepoPath)
	row.PRStatus = gitutil.DetectPRStatusFast(state.RepoPath)
	row.PullStatus = resolveStateStatus(state)
	row.Release = resolveRepoRelease(state.RepoPath, row.LatestBranch)
	row.LastSHA = resolveRowSHA(state)

	return row
}

func resolveRowSHA(state *PullRepoState) string {
	sha := state.NewSHA
	if len(sha) == 0 {
		sha = state.OldSHA
	}
	if len(sha) == 0 && len(state.RepoPath) > 0 {
		sha = gitutil.GetLastCommitSHA(state.RepoPath)
	}

	return ShortenSHA(sha)
}

func resolveRepoRelease(repoPath, latestBranch string) string {
	hasPath := len(repoPath) > 0
	if !hasPath {
		return "-"
	}

	rel := queryReleaseIdentifier(repoPath, latestBranch)
	if len(rel) > 0 {
		return rel
	}

	return "-"
}

func queryReleaseIdentifier(repoPath, latestBranch string) string {
	tag := gitutil.GetLatestTag(repoPath)
	if len(tag) > 0 {
		return tag
	}

	if rel := extractReleaseFromBranch(latestBranch); len(rel) > 0 {
		return rel
	}

	return readRepoManifestVersion(repoPath)
}

func extractReleaseFromBranch(branch string) string {
	lower := strings.ToLower(branch)
	if strings.HasPrefix(lower, "release/") {
		return strings.TrimPrefix(branch, "release/")
	}
	if strings.HasPrefix(lower, "v") && isSemverLike(branch) {
		return branch
	}

	return ""
}

func isSemverLike(s string) bool {
	raw := strings.TrimPrefix(s, "v")
	parts := strings.Split(raw, ".")
	hasEnoughParts := len(parts) >= 2
	if !hasEnoughParts {
		return false
	}

	_, err := strconv.Atoi(parts[0])

	return err == nil
}

func readRepoManifestVersion(repoPath string) string {
	v := readJsonVersion(filepath.Join(repoPath, "version.json"))
	if len(v) > 0 {
		return ensureVPrefix(v)
	}
	pkgV := readJsonVersion(filepath.Join(repoPath, "package.json"))
	if len(pkgV) > 0 {
		return ensureVPrefix(pkgV)
	}

	return ""
}

func readJsonVersion(filePath string) string {
	data, err := os.ReadFile(filePath)
	hasData := err == nil && len(data) > 0
	if !hasData {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}

	return extractVersionField(m)
}

func extractVersionField(m map[string]interface{}) string {
	if v, ok := m["Version"].(string); ok && len(v) > 0 {
		return v
	}
	if v, ok := m["version"].(string); ok && len(v) > 0 {
		return v
	}

	return ""
}

func ensureVPrefix(v string) string {
	trimmed := strings.TrimSpace(v)
	hasV := strings.HasPrefix(trimmed, "v")
	if hasV {
		return trimmed
	}

	return "v" + trimmed
}

func initBasePullTableRow(state *PullRepoState) model.PullTableRow {
	return model.PullTableRow{
		RepoName:    state.RepoName,
		Branch:      state.Branch,
		LastSHA:     ShortenSHA(state.NewSHA),
		CommitRange: state.CommitRange,
		Changes:     state.Changes,
		Duration:    FormatPullDuration(state.Duration),
		IsDirty:     state.IsDirty,
		Reason:      state.ErrorMsg,
	}
}

func resolveStateStatus(state *PullRepoState) string {
	if state.IsDirty {
		return "DIRTY"
	}

	return string(state.Step)
}

func handlePullRemediationForRecords(records []model.ScanRecord, opts pullOptions) {
	var remItems []RemediationItem
	for _, rec := range records {
		diag := gitutil.InspectDirtyState(rec.AbsolutePath)
		if diag.IsDirty {
			remItems = append(remItems, buildRemediationItem(rec, diag))
		}
	}
	handlePullRemediation(remItems, opts)
}

func buildRemediationItem(rec model.ScanRecord, diag gitutil.DirtyDiagnosis) RemediationItem {
	recipes := gitutil.GenerateRemediationRecipes(rec.AbsolutePath, diag)

	return RemediationItem{
		RepoName:      rec.RepoName,
		RepoPath:      rec.AbsolutePath,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
		Files:         diag.AllFiles,
	}
}

func finalizePullBatchTask(taskDB *store.DB, taskID int64, failCount int) error {
	if failCount > 0 {
		errMsg := fmt.Sprintf("pull batch finished with %d failure(s)", failCount)
		failPendingTask(taskDB, taskID, errMsg)
		fmt.Fprintf(os.Stderr, "\n  %s%s%s\n\n", constants.ColorRed, errMsg, constants.ColorReset)
		cliexit.Exit(1)

		return nil
	}
	completePendingTask(taskDB, taskID)

	return nil
}

func isPullCWDEnabled(opts pullOptions) bool {
	if hasExplicitPullTarget(opts) {
		return false
	}

	return isGitRepoCWD()
}

func isGitRepoCWD() bool {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

// runPullCWD streams pull in CWD, using progress bar unless isRaw is true.
func runPullCWD(isRaw ...bool) error {
	var pullErr error
	if len(isRaw) > 0 && isRaw[0] {
		pullErr = runPullCWDWithTransport(false, false, nil)
	} else {
		pullErr = runPullCWDTracked()
	}
	handleCWDIgnoreChecks()

	return pullErr
}

func handleCWDIgnoreChecks() {
	cwd, err := os.Getwd()
	if err != nil || !gitignoreagm.IsGitRepository(cwd) {
		return
	}
	if isCWDIgnoreCheckRecent(cwd) {
		return
	}
	executeCWDIgnoreAudit(cwd)
}

func isCWDIgnoreCheckRecent(cwd string) bool {
	ttl := resolvePullIgnoreTTL()
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return false
	}
	defer db.Close()
	return db.IsCheckRecent(cwd, ttl)
}

func executeCWDIgnoreAudit(cwd string) {
	start := time.Now()
	name := filepath.Base(cwd)
	issue := inspectRepoForIgnoreIssues(cwd, name)
	status := "clean"
	if issue.HasIssues() {
		status = "has_issues"
	}
	_ = store.RecordIgnoreCheckResult(cwd, name, status, 0, time.Since(start))
	if issue.HasIssues() {
		handleIgnoreRemediation([]IgnoreRepoIssue{issue}, false, false)
	}
}

func runPullCWDTracked() error {
	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "getwd")
	}

	state := executeCWDTrackedPull(cwd)
	row := buildPullTableRowFromState(state)
	RenderPullBatchTable([]model.PullTableRow{row})
	_ = RecordSinglePullSession(cwd, state, "pull")

	return nil
}

func executeCWDTrackedPull(cwd string) *PullRepoState {
	rec := model.ScanRecord{RepoName: filepath.Base(cwd), AbsolutePath: cwd}
	bar := NewPullProgressBar(1, false, false)
	bar.Start()
	state := ExecuteTrackedPull(rec, bar)
	bar.Stop()

	return state
}

func runPullCWDWithTransport(useSSH, useHTTPS bool, extraArgs []string) error {
	cwd, _ := os.Getwd()
	if !isGitRepoCWD() {
		return handleNonGitPull(cwd, extraArgs)
	}
	if !applyTransportSafe(cwd, useSSH, useHTTPS) {
		return nil
	}

	return executeGitPullCommand(cwd, extraArgs)
}

func applyTransportSafe(cwd string, useSSH, useHTTPS bool) bool {
	if _, _, _, err := ApplyTransportFlag(cwd, useSSH, useHTTPS); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		cliexit.HandleGeneralError(apperror.WrapSimple(err, "apply transport flag"))

		return false
	}

	return true
}

func buildGitPullEnv() []string {
	return gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)
}

func executeGitPullCommand(cwd string, extraArgs []string) error {
	gitArgs := append([]string{"pull"}, extraArgs...)
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s Running: git %s (cwd: %s)\n", constants.ColorCyan, arrow, constants.ColorReset, joinForLog(gitArgs), cwd)
	cmd := exec.Command("git", gitArgs...)
	cmd.Env = buildGitPullEnv()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return handleGitExecResult(cmd.Run())
}

func handleGitExecResult(err error) error {
	if err == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		cliexit.HandleError(apperror.WrapSimple(exitErr, "git pull"), exitErr.ExitCode())
		return nil
	}
	fmt.Fprintf(os.Stderr, "git pull failed: %v\n", err)
	cliexit.HandleGeneralError(apperror.WrapSimple(err, "git pull failed"))

	return nil
}

func ExtractTransportFlags(args []string) (bool, bool, []string) {
	var useSSH, useHTTPS bool
	rest := make([]string, 0, len(args))
	for _, a := range args {
		useSSH, useHTTPS, rest = classifyTransportToken(a, useSSH, useHTTPS, rest)
	}

	return useSSH, useHTTPS, rest
}

func classifyTransportToken(a string, useSSH, useHTTPS bool, rest []string) (bool, bool, []string) {
	switch strings.ToLower(a) {
	case "--ssh", "-ssh", "--sh", "-sh", "ssh":
		return true, useHTTPS, rest
	case "--https", "-https", "--ht", "-ht", "https":
		return useSSH, true, rest
	default:
		return useSSH, useHTTPS, append(rest, a)
	}
}

func beginPullTask(records []model.ScanRecord) (int64, *store.DB) {
	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine working directory: %v\n", wdErr)
	}

	cmdArgs := buildCommandArgs(append([]string{"pull"}, os.Args[2:]...))
	targetPath := workDir
	if len(records) == 1 {
		targetPath = records[0].AbsolutePath
	}

	return createPendingTask(constants.TaskTypePull, targetPath, workDir, "pull", cmdArgs)
}

func executePull(records []model.ScanRecord, bar *PullProgressBar, opts pullOptions) {
	if opts.parallel < 0 || opts.workers < 0 || opts.hands < 0 {
		cliexit.HandleError(apperror.NewSimple("invalid concurrency", "E9000"), 1)
	}
	opts = clampPullConcurrencyForSSH(opts)
	if opts.parallel > 1 {
		runPullParallel(records, bar, opts.parallel)
		return
	}
	runSerialPull(records, bar)
}

func clampPullConcurrencyForSSH(opts pullOptions) pullOptions {
	isSSH := cloneconcurrency.IsSSHSession() || opts.useSSH
	if isSSH {
		opts.workers = 1
		opts.hands = 1
		opts.parallel = 1
		return opts
	}
	effectiveWorkers := opts.workers
	if effectiveWorkers <= 0 {
		effectiveWorkers = opts.parallel
	}
	workers, hands := cloneconcurrency.ResolveWorkerHands(effectiveWorkers, opts.hands, opts.isWWOH, false)
	opts.workers = workers
	opts.hands = hands
	opts.parallel = workers * hands
	return opts
}

func runSerialPull(records []model.ScanRecord, bar *PullProgressBar) {
	for _, rec := range records {
		if bar != nil && bar.IsStopped() {
			break
		}
		ExecuteTrackedPull(rec, bar)
	}
}

func handlePullRemediation(remItems []RemediationItem, opts pullOptions) {
	if len(remItems) == 0 || isConcisePullOutput(opts.all, opts.showStatus) {
		return
	}
	dispatchRemediationSummary(remItems, opts)
}

func dispatchRemediationSummary(remItems []RemediationItem, opts pullOptions) {
	if opts.noFix {
		PrintRemediationSummaryNoPrompt(remItems)
		return
	}
	if opts.yes || opts.autoFix {
		PrintRemediationSummaryAutoFix(remItems)
		return
	}
	PrintRemediationSummary(remItems)
}

type pullFlagHolders struct {
	vFlag, aFlag, sFlag, oFlag, fixFlag, yFlag, noFixFlag, rawFlag *bool
	sshFlag, httpsFlag, statusFlag, jsonFlag, probeFlag            *bool
	wwohFlag                                                       *bool
	gFlag                                                          *string
	pFlag, wFlag, handFlag                                         *int
}

func initPullFlagSet() (*flag.FlagSet, *pullFlagHolders) {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	h := &pullFlagHolders{}
	registerPullCoreFlags(fs, h)
	registerPullRemediationFlags(fs, h)
	registerPullOutputFlags(fs, h)

	return fs, h
}

func registerPullCoreFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.vFlag = fs.Bool("verbose", false, constants.FlagDescVerbose)
	h.gFlag = fs.String("group", "", constants.FlagDescGroup)
	h.aFlag = fs.Bool("all", false, constants.FlagDescAll)
	h.sFlag = fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	h.pFlag = fs.Int("parallel", 0, constants.FlagDescPullParallel)
	h.oFlag = fs.Bool("only-available", false, constants.FlagDescPullOnlyAvailable)
	h.rawFlag = fs.Bool("raw", false, constants.FlagDescPullRaw)
	h.sshFlag = fs.Bool("ssh", false, "Pull using SSH transport")
	h.httpsFlag = fs.Bool("https", false, "Pull using HTTPS transport")
	fs.StringVar(h.gFlag, "g", "", constants.FlagDescGroup)
	fs.IntVar(h.pFlag, "p", 0, constants.FlagDescPullParallel)

	h.wFlag = fs.Int("w", 0, "Worker pool size")
	fs.IntVar(h.wFlag, "worker", 0, "Worker pool size")
	fs.IntVar(h.wFlag, "workers", 0, "Worker pool size")

	h.handFlag = fs.Int("hand", 0, "Hands per worker")
	fs.IntVar(h.handFlag, "hands", 0, "Hands per worker")
	fs.IntVar(h.handFlag, "h", 0, "Hands per worker")

	h.wwohFlag = fs.Bool("wwoh", false, "Worker with one hand")
	fs.BoolVar(h.wwohFlag, "worker-with-one-hand", false, "Worker with one hand")
}

func registerPullRemediationFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.fixFlag = fs.Bool("fix", false, "Auto-remediate dirty repos")
	h.yFlag = fs.Bool("yes", false, "Remediate without prompt")
	h.noFixFlag = fs.Bool("no-fix", false, "Skip remediation prompt")
	fs.BoolVar(h.yFlag, "y", false, "Remediate without prompt")
}

func registerPullOutputFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.statusFlag = fs.Bool("status", false, constants.FlagDescPullStatus)
	fs.BoolVar(h.statusFlag, "table", false, constants.FlagDescPullStatus)
	fs.BoolVar(h.statusFlag, "status-table", false, constants.FlagDescPullStatus)
	h.jsonFlag = fs.Bool("json", false, constants.FlagDescPullJSON)
	h.probeFlag = fs.Bool("probe", false, constants.FlagDescPullProbe)
	fs.BoolVar(h.probeFlag, "probe-repos", false, constants.FlagDescPullProbe)
}

func buildPullOptions(h *pullFlagHolders) pullOptions {
	opts := pullOptions{
		group: *h.gFlag, all: *h.aFlag, verbose: *h.vFlag,
		stopOnFail: *h.sFlag, parallel: *h.pFlag,
		workers: *h.wFlag, hands: *h.handFlag, isWWOH: *h.wwohFlag,
		onlyAvailable: *h.oFlag,
		autoFix:       *h.fixFlag, yes: *h.yFlag, noFix: *h.noFixFlag,
		isRaw: *h.rawFlag, useSSH: *h.sshFlag, useHTTPS: *h.httpsFlag,
		showStatus: *h.statusFlag, isJSON: *h.jsonFlag, isProbe: *h.probeFlag,
	}

	return opts
}

func parsePullFlags(args []string) pullOptions {
	fs, h := initPullFlagSet()
	fs.Parse(args)
	opts := buildPullOptions(h)
	if fs.NArg() > 0 {
		opts.slug = fs.Arg(0)
	}

	return opts
}

func initVerboseLog() {
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}
	log.Close()
}

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
	slugLower := strings.ToLower(slug)
	exact, partial := partitionBySlug(records, slugLower)
	if len(exact) > 0 {
		return exact
	}

	return partial
}

func partitionBySlug(records []model.ScanRecord, slugLower string) ([]model.ScanRecord, []model.ScanRecord) {
	var exact, partial []model.ScanRecord
	for _, r := range records {
		nameLower := strings.ToLower(r.RepoName)
		if nameLower == slugLower {
			exact = append(exact, r)
		} else if strings.Contains(nameLower, slugLower) {
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

func resolveExplicitPullTargets(opts pullOptions) ([]model.ScanRecord, bool) {
	records := resolvePullTargets(opts.slug, opts.group, opts.all)
	if len(records) == 0 {
		warnPullTargetNotFoundUnlessJSON(opts)

		return nil, false
	}

	return records, true
}

func warnPullTargetNotFoundUnlessJSON(opts pullOptions) {
	if !opts.isJSON {
		handlePullTargetNotFound(opts)
	}
}

func resolvePullBatchRecords(opts pullOptions) ([]model.ScanRecord, bool) {
	if hasExplicitPullTarget(opts) {
		return resolveExplicitPullTargets(opts)
	}
	records := discoverCWDPullRecords()
	if len(records) == 0 {
		handleNoTrackedReposInCWD()

		return nil, false
	}

	return deduplicatePullRecords(records), true
}

func hasExplicitPullTarget(opts pullOptions) bool {
	return opts.slug != "" || opts.group != "" || opts.all || HasAlias()
}

func discoverCWDPullRecords() []model.ScanRecord {
	cwd, _ := os.Getwd()
	records := ResolvePullDirectoryTargets(cwd)
	if len(records) == 0 {
		return findChildrenOfCWD(cwd)
	}

	return records
}

func handleNoTrackedReposInCWD() {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %snothing to pull: no tracked repositories found in or under this directory.%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func handleNonGitPull(cwd string, extraArgs []string) error {
	childRepos, err := fsutil.DiscoverChildGitRepos(cwd)
	if err == nil && len(childRepos) > 0 {
		return pullDiscoveredChildren(cwd, childRepos, extraArgs)
	}
	fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap pull` inside a repo)")
	cliexit.HandleValidationError(apperror.NewValidationError("not a git repository (run `gitmap pull` inside a repo)"))

	return nil
}

func pullDiscoveredChildren(cwd string, childRepos []string, extraArgs []string) error {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s Discovered %d child repositories in %s for pull:\n",
		constants.ColorCyan, arrow, constants.ColorReset, len(childRepos), cwd)
	for _, r := range childRepos {
		pullSingleDiscoveredChild(r, extraArgs)
	}

	return nil
}

func pullSingleDiscoveredChild(repoPath string, extraArgs []string) {
	fmt.Printf("      • %s\n", filepath.Base(repoPath))
	gitArgs := append([]string{"-C", repoPath, "pull"}, extraArgs...)
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func enqueueTaskQueue(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())
	tasksDB, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDB.Conn(), "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'pull', ?, ?, 'pending', ?, ?)", queueId, action, target, now, now)
	return queueId, tasksDB
}

func updateTaskQueue(db *store.TasksSplitDB, queueId, status string) {
	hasDB := db != nil
	if !hasDB {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(db.Conn(), "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?", status, now, queueId)
}

func closeTaskDB(db *store.TasksSplitDB) {
	hasDB := db != nil
	if hasDB {
		_ = db.Close()
	}
}
