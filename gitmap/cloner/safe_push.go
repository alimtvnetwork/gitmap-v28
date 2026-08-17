package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// SafePushOne pushes a single repo safely.
func SafePushOne(rec model.ScanRecord, repoDir string) model.CloneResult {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-push starting: %s → %s", rec.RepoName, repoDir)
	}

	var lastError string
	for attempt := 1; attempt <= constants.SafePullRetryAttempts; attempt++ {
		output, err := runGitPush(repoDir)
		if log != nil {
			log.Log("push attempt %d/%d for %s: exit=%v output=%s",
				attempt, constants.SafePullRetryAttempts, rec.RepoName, err, trimOutput(output))
		}
		var successResult *model.CloneResult
		if err == nil {
			successResult = &model.CloneResult{Record: rec, IsSuccess: true}
		}
		if successResult != nil && strings.Contains(output, "Everything up-to-date") {
			successResult.Notes = "up-to-date"
		}
		if successResult != nil {
			return *successResult
		}

		isNFF := isNonFastForwardRejection(output)
		if isNFF && log != nil {
			log.Log("push rejected (non-fast-forward) for %s — auto-running `git pull --rebase`", rec.RepoName)
		}
		var pullOut string
		var pullErr error
		if isNFF {
			pullOut, pullErr = runGitPullRebase(repoDir)
		}
		if isNFF && log != nil {
			log.Log("pull rebase output for %s: exit=%v output=%s", rec.RepoName, pullErr, trimOutput(pullOut))
		}
		if isNFF && pullErr == nil {
			continue
		}
		if isNFF && pullErr != nil {
			lastError = fmt.Sprintf("auto pull --rebase failed: %v\n%s", pullErr, trimOutput(pullOut))
			break
		}

		lastError = fmt.Sprintf("push failed: %v\n%s", err, trimOutput(output))

		if attempt < constants.SafePullRetryAttempts {
			time.Sleep(time.Duration(constants.SafePullRetryDelayMS) * time.Millisecond)
		}
	}

	if log != nil {
		log.Log("safe-push FAILED: %s — %s", rec.RepoName, lastError)
	}

	return model.CloneResult{Record: rec, IsSuccess: false, Error: lastError}
}

func runGitPush(repoDir string) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, constants.CmdPush)
	cmd.Env = append(os.Environ(),
		constants.EnvGitTerminalPromptZero,
		constants.EnvGitAskpassEmpty,
		constants.EnvSSHAskpassEmpty,
		constants.EnvGitSSHCommandBatchYes)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runGitPullRebase(repoDir string) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, constants.GitPull, "--rebase")
	cmd.Env = append(os.Environ(),
		constants.EnvGitTerminalPromptZero,
		constants.EnvGitAskpassEmpty,
		constants.EnvSSHAskpassEmpty,
		constants.EnvGitSSHCommandBatchYes)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func isNonFastForwardRejection(stderr string) bool {
	lower := strings.ToLower(stderr)
	if !strings.Contains(lower, "[rejected]") && !strings.Contains(lower, "failed to push some refs") {
		return false
	}
	return strings.Contains(lower, "fetch first") || strings.Contains(lower, "non-fast-forward")
}
