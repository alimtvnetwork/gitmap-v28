// Package cloner re-clones repos from structured files.
package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	clonerconstants "github.com/alimtvnetwork/gitmap-v28/cli/cloner/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

// Unlink-failure patterns, compiled lazily on first use.
const (
	unlinkOldPattern    = `(?i)unable to unlink old '([^']+)'`
	unlinkPromptPattern = `(?i)unlink of file '([^']+)' failed`
)

var (
	unlinkOldRegex    = lazyregex.New(unlinkOldPattern)
	unlinkPromptRegex = lazyregex.New(unlinkPromptPattern)
)

func cloneOrPullOne(rec model.ScanRecord, targetDir string, opts CloneOptions) model.CloneResult {
	dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
	dirExists := isPathExisting(dest)
	if dirExists && opts.IsMissingOnly {
		return model.CloneResult{Record: rec, IsSuccess: true, Notes: "skipped (existing directory)"}
	}
	res, isHandled := handleExistingTargetDir(rec, dest, dirExists, opts)
	if isHandled {
		return res
	}

	return cloneOne(rec, targetDir)
}

func isPathExisting(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func handleExistingTargetDir(rec model.ScanRecord, dest string, dirExists bool, opts CloneOptions) (model.CloneResult, bool) {
	if cleanErr := cleanDirIfRequested(dirExists, opts.IsClean, dest); cleanErr != nil {
		msg := fmt.Sprintf("failed to clean existing directory %q: %v", dest, cleanErr)
		return model.CloneResult{Record: rec, IsSuccess: false, Error: msg}, true
	}
	isStillExists := dirExists && !opts.IsClean
	if isStillExists && opts.IsSafePull && isGitRepo(dest) {
		rewriteRemoteToSSHIfNeeded(dest, opts.IsUseSSH)
		return safePullRepo(rec, dest), true
	}
	if isStillExists && !isGitRepo(dest) {
		msg := fmt.Sprintf("target directory %q exists but is not a git repository (conflict)", dest)
		return model.CloneResult{Record: rec, IsSuccess: false, Error: msg}, true
	}

	return model.CloneResult{}, false
}

func isGitRepo(path string) bool {
	return IsGitRepo(path)
}

// IsGitRepo checks whether the given path contains a .git directory.
func IsGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))

	return err == nil
}

// IsMissingRepo returns true when the path is not a valid git repository.
func IsMissingRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))

	return err != nil
}

// SafePullOne runs safe-pull on a single repo. Exported for use by the pull command.
func SafePullOne(rec model.ScanRecord, repoDir string) model.CloneResult {
	return SafePullOneWithProgress(rec, repoDir, nil)
}

// SafePullOneWithProgress runs safe-pull forwarding streaming lines to callback.
func SafePullOneWithProgress(rec model.ScanRecord, repoDir string, onProgress func(string)) model.CloneResult {
	return safePullRepoWithProgress(rec, repoDir, onProgress)
}

func safePullRepo(rec model.ScanRecord, repoDir string) model.CloneResult {
	return safePullRepoWithProgress(rec, repoDir, nil)
}

func safePullRepoWithProgress(rec model.ScanRecord, repoDir string, onProgress func(string)) model.CloneResult {
	logSafePullStart(rec, repoDir)

	if HasUnmergedConflicts(repoDir) {
		conflicts := ExtractConflictedFilesFromRepo(repoDir)
		conflictDesc := formatActiveConflictDesc(conflicts)
		err := apperror.WrapSimple(ErrMergeConflict, conflictDesc)

		return buildAttemptFailureResult(rec, repoDir, 1, conflictDesc, err)
	}

	var lastFail model.CloneResult
	for attempt := 1; attempt <= constants.SafePullRetryAttempts; attempt++ {
		res, isDone := executePullAttempt(rec, repoDir, attempt, onProgress)
		if isDone {
			return res
		}
		lastFail = res
		sleepRetryBackoff(attempt)
	}
	if lastFail.Error != "" {
		return lastFail
	}

	return buildFinalFailureResult(rec)
}

func logSafePullStart(rec model.ScanRecord, repoDir string) {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-pull starting: %s → %s", rec.RepoName, repoDir)
	}
}

func hasWincredmanFailure(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "wincredman") || strings.Contains(lower, "unable to persist credentials")
}

func retryPullWithSSHConversion(rec model.ScanRecord, repoDir string, attempt int, onProgress func(string)) (model.CloneResult, bool) {
	log := verbose.Get()
	if log != nil {
		log.Log("wincredman failure detected for %s; converting remote to SSH and retrying", rec.RepoName)
	}
	rewriteRemoteToSSHIfNeeded(repoDir, true)
	out, err := runGitPullWithProgress(repoDir, rec.Branch, onProgress)
	logPullAttempt(rec, attempt, out, err)
	if err == nil {
		return handleAttemptSuccess(rec, attempt, out), true
	}
	handleAttemptFailure(repoDir, out)

	return buildAttemptFailureResult(rec, repoDir, attempt, out, err), isTerminalPullError(err, out, repoDir)
}

func isTerminalPullError(err error, output, repoDir string) bool {
	return isMergeConflict(err, output) || HasUnmergedConflicts(repoDir)
}

func executePullAttempt(rec model.ScanRecord, repoDir string, attempt int, onProgress func(string)) (model.CloneResult, bool) {
	output, err := runGitPullWithProgress(repoDir, rec.Branch, onProgress)
	logPullAttempt(rec, attempt, output, err)
	if err == nil {
		return handleAttemptSuccess(rec, attempt, output), true
	}
	if hasWincredmanFailure(output) {
		return retryPullWithSSHConversion(rec, repoDir, attempt, onProgress)
	}
	handleAttemptFailure(repoDir, output)
	if isTerminalPullError(err, output, repoDir) {
		return buildAttemptFailureResult(rec, repoDir, attempt, output, err), true
	}

	return buildAttemptFailureResult(rec, repoDir, attempt, output, err), false
}

func logPullAttempt(rec model.ScanRecord, attempt int, output string, err error) {
	log := verbose.Get()
	if log != nil {
		log.Log("pull attempt %d/%d for %s: exit=%v output=%s",
			attempt, constants.SafePullRetryAttempts, rec.RepoName, err, trimOutput(output))
	}
}

func handleAttemptSuccess(rec model.ScanRecord, attempt int, output string) model.CloneResult {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-pull succeeded: %s (attempt %d)", rec.RepoName, attempt)
	}
	res := model.CloneResult{Record: rec, IsSuccess: true}
	if strings.Contains(output, "[auto-merged]") || strings.Contains(output, "Merge made by") {
		res.Notes = "auto-merged"
	} else if strings.Contains(output, "Already up to date.") {
		res.Notes = "up-to-date"
	}

	return res
}

func handleAttemptFailure(repoDir, output string) {
	log := verbose.Get()
	cleared := clearReadOnlyAttrs(repoDir, output)
	if log != nil && cleared {
		log.Log("cleared read-only attributes for blocked files in %s", repoDir)
	}
}

func buildAttemptFailureResult(rec model.ScanRecord, repoDir string, attempt int, output string, err error) model.CloneResult {
	diagnosis := buildPullDiagnosis(repoDir, output)
	log := verbose.Get()
	if log != nil {
		log.Log("diagnosis for %s: %s", rec.RepoName, diagnosis)
	}
	msg := formatPullAttemptError(rec, repoDir, attempt, output, diagnosis, err)

	return model.CloneResult{Record: rec, IsSuccess: false, Error: msg}
}

func formatPullAttemptError(rec model.ScanRecord, repoDir string, attempt int, output, diagnosis string, err error) string {
	return fmt.Sprintf(
		"safe-pull failed for %s (attempt %d/%d): repo=%q branch=%q: %v\n%s\nDiagnosis: %s",
		recordTag(rec), attempt, constants.SafePullRetryAttempts, repoDir, rec.Branch, err, trimOutput(output), diagnosis,
	)
}

func sleepRetryBackoff(attempt int) {
	if attempt < constants.SafePullRetryAttempts {
		time.Sleep(time.Duration(constants.SafePullRetryDelayMS) * time.Millisecond)
	}
}

func buildFinalFailureResult(rec model.ScanRecord) model.CloneResult {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-pull FAILED after all retries: %s", rec.RepoName)
	}

	return model.CloneResult{Record: rec, IsSuccess: false, Error: "safe-pull failed after all retries"}
}

func buildSafePullEnv() []string {
	return gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)
}

func runGitPullWithProgress(repoDir, branch string, onProgress SafePullProgressFunc) (string, error) {
	out, err := execGitPullFF(repoDir, onProgress)
	if err == nil {
		return out, nil
	}
	if isDivergedOutput(out) {
		return attemptAutoMergePull(repoDir, branch, onProgress)
	}

	return out, err
}

func execGitPullFF(repoDir string, onProgress SafePullProgressFunc) (string, error) {
	args := append([]string{constants.GitDirFlag, repoDir}, clonerconstants.GitPullFastForwardArgs...)
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Env = buildSafePullEnv()
	stream := newProgressStreamWriter(onProgress)
	cmd.Stdout = stream
	cmd.Stderr = stream
	err := cmd.Run()
	stream.flush()

	return stream.output(), err
}

// IsDivergedOutput returns true if git output indicates diverged local and remote branches.
func IsDivergedOutput(output string) bool {
	return isDivergedOutput(output)
}

func isDivergedOutput(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "diverg") ||
		strings.Contains(lower, "not possible to fast-forward") ||
		strings.Contains(lower, "cannot fast-forward") ||
		strings.Contains(lower, "reconcile divergent") ||
		strings.Contains(lower, "need to specify how to reconcile") ||
		strings.Contains(lower, "non-fast-forward") ||
		strings.Contains(lower, "refusing to merge unrelated histories")
}

func executeAutoMergeCmd(dir string, progress SafePullProgressFunc) (string, error) {
	args := append([]string{constants.GitDirFlag, dir}, clonerconstants.GitPullAutoMergeArgs...)
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Env = buildSafePullEnv()
	stream := newProgressStreamWriter(progress)
	cmd.Stdout = stream
	cmd.Stderr = stream
	err := cmd.Run()
	stream.flush()

	return stream.output(), err
}

func attemptAutoMergePull(dir, branch string, progress SafePullProgressFunc) (string, error) {
	out, err := executeAutoMergeCmd(dir, progress)
	if err == nil && !isMergeConflictOutput(out) {
		return "[auto-merged]\n" + out, nil
	}
	SafeAbortMerge(dir)

	return out, ErrMergeConflict
}

func cleanDirIfRequested(isDirExists, isClean bool, dest string) *apperror.AppError {
	if !isDirExists || !isClean {
		return nil
	}

	if err := os.RemoveAll(dest); err != nil {
		return apperror.WrapSimple(err, "remove all")
	}

	return nil
}

func formatActiveConflictDesc(conflicts []string) string {
	if len(conflicts) == 0 {
		return "active unmerged conflict detected"
	}

	return fmt.Sprintf("active unmerged conflict detected in %s", strings.Join(conflicts, ", "))
}
