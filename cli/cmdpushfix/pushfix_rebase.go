package cmdpushfix

import (
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// CheckIsNonFastForward detects non-fast-forward rejection from git stderr.
func CheckIsNonFastForward(stderr string) bool {
	lower := strings.ToLower(stderr)
	isRejected := strings.Contains(lower, "[rejected]") || strings.Contains(lower, "failed to push some refs")
	if !isRejected {
		return false
	}
	return strings.Contains(lower, "fetch first") || strings.Contains(lower, "non-fast-forward")
}

// ExecuteAutoRebase runs git pull --rebase with safe anti-hang environment.
func ExecuteAutoRebase(repoDir string) error {
	PrintRemediationApplied("Diverged remote detected — auto-running `git pull --rebase`")
	cmd := exec.Command("git", "-C", repoDir, "pull", "--rebase")
	cmd.Env = gitutil.BuildSafeGitEnv()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
