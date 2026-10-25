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
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isPullCWDEnabled(opts pullOptions) bool {
	if opts.all || opts.group != "" {
		return false
	}
	if opts.slug == "." || opts.slug == "cwd" || opts.slug == "" {
		return isGitRepoCWD()
	}

	return false
}

func isGitRepoCWD() bool {
	cwd, err := os.Getwd()
	if err != nil || gitutil.IsWorkspaceContainerDirectory(cwd) || !gitutil.IsValidGitRepository(cwd) {
		return false
	}

	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

// runPullCWD streams pull in CWD, using progress bar unless isRaw is true.
func runPullCWD(isRaw ...bool) error {
	cwd, _ := os.Getwd()
	if gitutil.IsWorkspaceContainerDirectory(cwd) {
		return handleWorkspaceContainerPull(cwd)
	}

	if !gitutil.IsValidGitRepository(cwd) {
		return handleNonGitPull(cwd, nil)
	}

	raw := len(isRaw) > 0 && isRaw[0]
	served, err := checkCWDCacheOrProceed(raw)
	if served {
		return err
	}

	var pullErr error
	if raw {
		pullErr = runPullCWDWithTransport(false, false, nil)
	} else {
		pullErr = runPullCWDTracked()
	}
	handleCWDIgnoreChecks()

	return pullErr
}

func handleWorkspaceContainerPull(cwd string) error {
	base := filepath.Base(cwd)
	fmt.Fprintf(os.Stderr, "  %s[!] Cannot pull current directory: '%s' is a workspace container, not a Git repository.%s\n",
		constants.ColorYellow, base, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      Run 'gitmap pull' to pull all child repositories, or 'cd <repo>' first.\n\n")

	return apperror.NewValidationError("cannot pull workspace container directory: " + base)
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
	if state != nil {
		_ = RecordSinglePullSession(cwd, state, "pull")
	}

	return nil
}

func executeCWDTrackedPull(cwd string) *PullRepoState {
	targetDir := resolveCWDRepoDir(cwd)
	if targetDir == "" {
		_ = handleNonGitPull(cwd, nil)
		return nil
	}

	rec := model.ScanRecord{RepoName: filepath.Base(targetDir), AbsolutePath: targetDir}
	printPullProgressNotice()
	state := ExecuteTrackedPull(rec, nil)
	clearPullProgressNotice()
	RenderSingleRepoPullSummary(state, targetDir)

	return state
}

func printPullProgressNotice() {
	if resolveTerminalMode() {
		fmt.Print("  ↓ pulling remote changes...")
		return
	}

	fmt.Println("  ↓ pulling remote changes...")
}

func clearPullProgressNotice() {
	if resolveTerminalMode() {
		fmt.Print("\r\033[K")
	}
}

func resolveCWDRepoDir(cwd string) string {
	if gitutil.IsWorkspaceContainerDirectory(cwd) {
		return ""
	}

	top := getGitRepoCWDTopLevel()
	if top != "" && gitutil.IsValidGitRepository(top) && !gitutil.IsWorkspaceContainerDirectory(top) {
		return top
	}

	if gitutil.IsValidGitRepository(cwd) {
		return cwd
	}

	return ""
}

func getGitRepoCWDTopLevel() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}

	return filepath.Clean(strings.TrimSpace(string(out)))
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
	records := loadAllRecordsDB()
	if len(records) == 0 {
		handleNoTrackedReposInCWD()

		return nil, false
	}

	return deduplicatePullRecords(records), true
}

func hasExplicitPullTarget(opts pullOptions) bool {
	return opts.slug != "" || opts.group != "" || opts.all || HasAlias()
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
