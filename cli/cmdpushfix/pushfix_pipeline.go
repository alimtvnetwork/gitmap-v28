package cmdpushfix

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// ExecutePushFix runs the complete 4-phase diagnosis, recovery, and push pipeline.
func ExecutePushFix(opts PushFixOptions) error {
	repoDir, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "cmdpushfix.Getwd")
	}

	PrintPushFixHeader(repoDir)
	state, probeErr := inspectGitState(repoDir, opts)
	if probeErr != nil {
		return probeErr
	}

	executeAuthRemediation(state, opts)
	return executePushDispatch(state, opts)
}

func inspectGitState(repoDir string, opts PushFixOptions) (*PushFixState, error) {
	PrintPhaseStep(1, "Pre-flight Inspection", "Checking git status, remotes & branches")
	branch := resolveGitBranch(repoDir, opts.Branch)
	remote := opts.Remote
	if len(remote) == 0 {
		remote = "origin"
	}

	remoteURL := resolveRemoteURL(repoDir, remote)
	unpushed, hasUpstream := countUnpushedCommits(repoDir)
	state := &PushFixState{
		RepoDir:       repoDir,
		RemoteName:    remote,
		RemoteURL:     remoteURL,
		Branch:        branch,
		UnpushedCount: unpushed,
		HasUpstream:   hasUpstream,
	}
	return state, nil
}

func executeAuthRemediation(state *PushFixState, opts PushFixOptions) {
	PrintPhaseStep(2, "Remote Auth Probe", "Verifying SSH & transport credentials")
	host := extractHostFromURL(state.RemoteURL)
	isSSHReady, _ := ProbeSSHConnection(host)

	PrintPhaseStep(3, "Self-Healing Engine", "Synchronizing keys, configs & transport")
	if !isSSHReady || opts.IsSSH {
		RemediateSSHKeyAndConfig(host)
	}

	if opts.IsSSH && strings.HasPrefix(state.RemoteURL, "http") {
		ConvertRemoteToSSH(state.RepoDir, state.RemoteName, state.RemoteURL)
	}
}

func executePushDispatch(state *PushFixState, opts PushFixOptions) error {
	PrintPhaseStep(4, "Push Execution", fmt.Sprintf("Pushing to %s/%s", state.RemoteName, state.Branch))
	pushArgs := buildPushArgs(state, opts)
	runErr, stderr := runSafePushCommand(state.RepoDir, pushArgs)
	if runErr == nil {
		PrintPushFixSuccess(state.RemoteName, state.Branch, resolveHeadSHA(state.RepoDir), opts.IsDryRun)
		return nil
	}

	if CheckIsNonFastForward(stderr) && retryPushAfterRebase(state.RepoDir, pushArgs) {
		PrintPushFixSuccess(state.RemoteName, state.Branch, resolveHeadSHA(state.RepoDir), opts.IsDryRun)
		return nil
	}

	return createReportedPushError(runErr)
}

func retryPushAfterRebase(repoDir string, pushArgs []string) bool {
	if rebaseErr := ExecuteAutoRebase(repoDir); rebaseErr != nil {
		return false
	}
	runRetryErr, _ := runSafePushCommand(repoDir, pushArgs)

	return runRetryErr == nil
}

func buildPushArgs(state *PushFixState, opts PushFixOptions) []string {
	args := []string{"push"}
	if opts.IsDryRun {
		args = append(args, "--dry-run")
	}
	if opts.IsForce {
		args = append(args, "--force-with-lease")
	}
	if !state.HasUpstream {
		args = append(args, "-u")
	}
	args = append(args, state.RemoteName, state.Branch)
	return args
}

func runSafePushCommand(repoDir string, args []string) (error, string) {
	cmd := exec.Command("git", append([]string{"-C", repoDir}, args...)...)
	cmd.Env = gitutil.BuildSafeGitEnv()
	var buf bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
	return cmd.Run(), buf.String()
}

func resolveGitBranch(repoDir, customBranch string) string {
	if len(customBranch) > 0 {
		return customBranch
	}
	out, err := exec.Command("git", "-C", repoDir, "branch", "--show-current").Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}
	return "main"
}

func resolveRemoteURL(repoDir, remote string) string {
	out, err := exec.Command("git", "-C", repoDir, "remote", "get-url", remote).Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func countUnpushedCommits(repoDir string) (int, bool) {
	out, err := exec.Command("git", "-C", repoDir, "rev-list", "--count", "@{u}..HEAD").Output()
	if err != nil {
		return 0, false
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return count, true
}

func resolveHeadSHA(repoDir string) string {
	out, err := exec.Command("git", "-C", repoDir, "rev-parse", "--short", "HEAD").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func extractHostFromURL(rawURL string) string {
	if strings.Contains(rawURL, "@") && strings.Contains(rawURL, ":") {
		parts := strings.SplitN(rawURL, "@", 2)
		hostPart := strings.SplitN(parts[1], ":", 2)
		return hostPart[0]
	}
	return "github.com"
}

func createReportedPushError(runErr error) error {
	appErr := apperror.WrapSimple(runErr, "cmdpushfix.ExecutePushFix")
	appErr.Type = apperror.ErrorTypeAbort
	if appErr.Ctx == nil {
		appErr.Ctx = make(map[string]any)
	}
	appErr.Ctx["reported"] = true
	return appErr
}
