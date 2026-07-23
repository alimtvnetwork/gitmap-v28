// Package cloner: LFS smudge failure recovery.
//
// Some public repos track binary assets with Git LFS but ship broken or
// missing LFS objects on the remote (404 on the LFS storage backend).
// Vanilla `git clone` then reports:
//
//	Errors logged to '.../.git/lfs/logs/<ts>.log'.
//	error: external filter 'git-lfs filter-process' failed
//	fatal: <path>: smudge filter lfs failed
//	warning: Clone succeeded, but checkout failed.
//
// Git actually did fetch every object; only the LFS smudge on checkout
// blew up. The user-facing symptom is a non-zero exit and a partially
// checked-out working tree. Retrying the clone with the well-known
// GIT_LFS_SKIP_SMUDGE=1 env var skips the smudge filter and lets the
// checkout complete. Pointer files stay as pointers, which is the
// correct behavior when the LFS objects are unavailable.
package cloner

import (
	"os"
	"os/exec"
	"strings"
)

// LFSSkipSmudgeEnv is the standard opt-out env var honored by git-lfs.
const LFSSkipSmudgeEnv = "GIT_LFS_SKIP_SMUDGE=1"

// LFSRetryNote is appended to CloneResult.Notes when the retry fires,
// so downstream reports show the recovery instead of silently masking
// a broken remote.
const LFSRetryNote = "lfs-skip-smudge-retry"

// isLFSSmudgeFailure returns true when the combined output of a failed
// `git clone` matches the git-lfs smudge failure signature.
func isLFSSmudgeFailure(out string) bool {
	lower := strings.ToLower(out)
	return strings.Contains(lower, "smudge filter lfs failed") ||
		strings.Contains(lower, "git-lfs filter-process' failed") ||
		strings.Contains(lower, "external filter 'git-lfs")
}

// retryCloneSkipSmudge reruns the same argv with GIT_LFS_SKIP_SMUDGE=1.
// Caller must have already removed the destination if git left it in a
// partial state; we do that here defensively so the retry cannot fail
// with "destination already exists".
func retryCloneSkipSmudge(bin string, args []string, dest string) ([]byte, error) {
	_ = os.RemoveAll(dest)
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), LFSSkipSmudgeEnv)
	return cmd.CombinedOutput()
}
