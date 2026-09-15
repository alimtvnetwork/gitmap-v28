// Package cloner re-clones repos from structured files.
package cloner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

var (
	unlinkOldRegex    = regexp.MustCompile(`(?i)unable to unlink old '([^']+)'`)
	unlinkPromptRegex = regexp.MustCompile(`(?i)unlink of file '([^']+)' failed`)
)

type progressStreamWriter struct {
	mu         sync.Mutex
	buf        bytes.Buffer
	lineBuf    []byte
	onProgress func(string)
}

func newProgressStreamWriter(onProgress func(string)) *progressStreamWriter {
	return &progressStreamWriter{
		onProgress: onProgress,
	}
}

func (w *progressStreamWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf.Write(p)
	if w.onProgress != nil {
		w.processChunkLocked(p)
	}

	return len(p), nil
}

func (w *progressStreamWriter) processChunkLocked(p []byte) {
	for _, b := range p {
		isDelim := b == '\r' || b == '\n'
		if isDelim {
			w.flushLineLocked()
			continue
		}
		w.lineBuf = append(w.lineBuf, b)
	}
}

func (w *progressStreamWriter) flushLineLocked() {
	if len(w.lineBuf) == 0 {
		return
	}
	line := string(w.lineBuf)
	w.lineBuf = w.lineBuf[:0]
	w.onProgress(line)
}

func (w *progressStreamWriter) flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.flushLineLocked()
}

func (w *progressStreamWriter) output() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.String()
}

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
	for attempt := 1; attempt <= constants.SafePullRetryAttempts; attempt++ {
		res, isDone := executePullAttempt(rec, repoDir, attempt, onProgress)
		if isDone {
			return res
		}
		sleepRetryBackoff(attempt)
	}

	return buildFinalFailureResult(rec)
}

func logSafePullStart(rec model.ScanRecord, repoDir string) {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-pull starting: %s → %s", rec.RepoName, repoDir)
	}
}

func executePullAttempt(rec model.ScanRecord, repoDir string, attempt int, onProgress func(string)) (model.CloneResult, bool) {
	output, err := runGitPullWithProgress(repoDir, onProgress)
	logPullAttempt(rec, attempt, output, err)
	if err == nil {
		return handleAttemptSuccess(rec, attempt, output), true
	}
	handleAttemptFailure(repoDir, output)

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
	if strings.Contains(output, "Already up to date.") {
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

func runGitPullWithProgress(repoDir string, onProgress func(string)) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, constants.GitPull, "--progress", constants.GitFFOnlyFlag, "--autostash")
	cmd.Env = append(os.Environ(),
		constants.EnvGitTerminalPromptZero,
		constants.EnvGitAskpassEmpty,
		constants.EnvSSHAskpassEmpty,
		constants.EnvGitSSHCommandBatchYes)
	stream := newProgressStreamWriter(onProgress)
	cmd.Stdout = stream
	cmd.Stderr = stream
	err := cmd.Run()
	stream.flush()

	return stream.output(), err
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
