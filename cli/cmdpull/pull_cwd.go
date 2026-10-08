package cmdpull

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

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
	cleanDir := strings.TrimSpace(cwd)
	if cleanDir == "" {
		return apperror.NewValidationError("cannot execute git pull in empty directory")
	}

	gitArgs := append([]string{"pull"}, extraArgs...)
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s Running: git %s (cwd: %s)\n", constants.ColorCyan, arrow, constants.ColorReset, joinForLog(gitArgs), cleanDir)
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = cleanDir
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
