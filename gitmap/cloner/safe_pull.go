// Package cloner re-clones repos from structured files.
package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

var (
	unlinkOldRegex    = regexp.MustCompile(`(?i)unable to unlink old '([^']+)'`)
	unlinkPromptRegex = regexp.MustCompile(`(?i)unlink of file '([^']+)' failed`)
)

func cloneOrPullOne(rec model.ScanRecord, targetDir string, safePull bool) model.CloneResult {
	dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
	if safePull && isGitRepo(dest) {
		return safePullRepo(rec, dest)
	}

	return cloneOne(rec, targetDir)
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
	return safePullRepo(rec, repoDir)
}

func safePullRepo(rec model.ScanRecord, repoDir string) model.CloneResult {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-pull starting: %s → %s", rec.RepoName, repoDir)
	}

	var lastError string
	for attempt := 1; attempt <= constants.SafePullRetryAttempts; attempt++ {
		output, err := runGitPull(repoDir)
		if log != nil {
			log.Log("pull attempt %d/%d for %s: exit=%v output=%s",
				attempt, constants.SafePullRetryAttempts, rec.RepoName, err, trimOutput(output))
		}
		if err == nil {
			if log != nil {
				log.Log("safe-pull succeeded: %s (attempt %d)", rec.RepoName, attempt)
			}
			return model.CloneResult{Record: rec, Success: true}
		}

		cleared := clearReadOnlyAttrs(repoDir, output)
		if log != nil && cleared {
			log.Log("cleared read-only attributes for blocked files in %s", repoDir)
		}
		diagnosis := buildPullDiagnosis(repoDir, output)
		if log != nil {
			log.Log("diagnosis for %s: %s", rec.RepoName, diagnosis)
		}
		lastError = fmt.Sprintf(
			"safe-pull failed for %s (attempt %d/%d): repo=%q branch=%q: %v\n%s\nDiagnosis: %s",
			recordTag(rec),
			attempt,
			constants.SafePullRetryAttempts,
			repoDir,
			rec.Branch,
			err,
			trimOutput(output),
			diagnosis,
		)

		if attempt < constants.SafePullRetryAttempts {
			time.Sleep(time.Duration(constants.SafePullRetryDelayMS) * time.Millisecond)
		}
	}

	if log != nil {
		log.Log("safe-pull FAILED after all retries: %s — %s", rec.RepoName, lastError)
	}

	return model.CloneResult{Record: rec, Success: false, Error: lastError}
}

func runGitPull(repoDir string) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, constants.GitPull, constants.GitFFOnlyFlag)
	out, err := cmd.CombinedOutput()

	return string(out), err
}
