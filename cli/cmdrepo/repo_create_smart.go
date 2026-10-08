package cmdrepo

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var stdinInteractiveHook func() bool
var promptReaderHook func() (string, error)

func extractOwnerFromURL(raw string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(raw), ".git")
	parts := strings.Split(clean, "/")
	if len(parts) < 2 {
		return ""
	}

	owner := parts[len(parts)-2]
	if idx := strings.LastIndex(owner, ":"); idx >= 0 {
		owner = owner[idx+1:]
	}

	return owner
}

func resolveOwnerFromGitConfig() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "config", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return extractOwnerFromURL(strings.TrimSpace(string(out)))
}

func resolveOwnerFromGH() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "api", "user", "-q", ".login")
	out, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}

	return ""
}

func resolveGitHubOwner() string {
	if user := os.Getenv("GITHUB_USER"); user != "" {
		return user
	}

	if owner := resolveOwnerFromGitConfig(); owner != "" {
		return owner
	}

	return resolveOwnerFromGH()
}

func probeLsRemoteURL(remoteURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--exit-code", remoteURL, "HEAD")
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if cmd.Run() == nil {
		return true
	}

	cmdHeads := exec.CommandContext(ctx, "git", "ls-remote", "--heads", remoteURL)
	cmdHeads.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmdHeads.CombinedOutput()

	return err == nil && len(strings.TrimSpace(string(out))) > 0
}

func probeViaGitLsRemote(slug string) bool {
	if strings.Contains(slug, "/") {
		return probeLsRemoteURL(fmt.Sprintf("https://github.com/%s.git", slug))
	}

	if owner := resolveGitHubOwner(); owner != "" {
		return probeLsRemoteURL(fmt.Sprintf("https://github.com/%s/%s.git", owner, slug))
	}

	return false
}

func probeViaGhRepoView(slug string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "repo", "view", slug, "--json", "name")

	return cmd.Run() == nil
}

func probeRemoteRepoExists(slug string) (bool, error) {
	if slug == "" {
		return false, nil
	}

	if probeViaGitLsRemote(slug) {
		return true, nil
	}

	if probeViaGhRepoView(slug) {
		return true, nil
	}

	return false, nil
}

func isRepoCollisionOutput(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "already exists")
}

func isTerminalStdinInteractive() bool {
	if !cmddb.IsInteractiveStdin() {
		return false
	}

	return term.IsTerminal(int(os.Stdin.Fd()))
}

func isStdinInteractive() bool {
	if stdinInteractiveHook != nil {
		return stdinInteractiveHook()
	}

	return isTerminalStdinInteractive()
}

func readInteractiveCloneAnswer(slug string) (string, error) {
	if promptReaderHook != nil {
		return promptReaderHook()
	}

	fmt.Printf("Remote repository '%s' already exists. Do you want to clone it? [Y/n]: ", slug)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	trimmed := strings.TrimSpace(line)
	if err != nil && len(trimmed) == 0 {
		return "", err
	}

	return trimmed, nil
}

func isCloneConfirmedFromAnswer(ans string) bool {
	lower := strings.ToLower(strings.TrimSpace(ans))

	return lower == "" || lower == "y" || lower == "yes"
}

func hasCloneOnExistsFlag() bool {
	return hasArgFlag(os.Args, "--clone-on-exists")
}

func hasForceFlag() bool {
	return hasArgFlag(os.Args, "--force") || hasArgFlag(os.Args, "-f")
}

func handleRemoteCollision(absDir, slug string) (string, error) {
	if isStdinInteractive() {
		return handleInteractiveCollision(absDir, slug)
	}

	return handleHeadlessCollision(absDir, slug)
}

func handleInteractiveCollision(absDir, slug string) (string, error) {
	if hasCloneOnExistsFlag() {
		return dispatchCloneFallback(absDir, slug)
	}

	ans, err := readInteractiveCloneAnswer(slug)
	if err != nil {
		return "", apperror.WrapSimple(err, "read confirmation:")
	}

	if isCloneConfirmedFromAnswer(ans) {
		return dispatchCloneFallback(absDir, slug)
	}

	fmt.Printf("  %sRepository creation aborted: remote repository '%s' already exists.%s\n", constants.ColorYellow, slug, constants.ColorReset)

	return "", apperror.NewSimple("repository creation aborted: remote already exists", "E1078")
}

func handleHeadlessCollision(absDir, slug string) (string, error) {
	if hasCloneOnExistsFlag() {
		fmt.Printf("  %sHeadless gate: --clone-on-exists active, cloning existing repo...%s\n", constants.ColorCyan, constants.ColorReset)

		return dispatchCloneFallback(absDir, slug)
	}

	if hasForceFlag() {
		fmt.Printf("  %sHeadless gate: --force active, overwriting existing remote repository...%s\n", constants.ColorYellow, constants.ColorReset)

		return executeForcePushRemote(absDir, slug)
	}

	return "", apperror.NewSimple("remote repository already exists and stdin is non-interactive; use --clone-on-exists or --force", "E1078")
}

func prepareCloneTargetDir(absDir string) {
	cwd, err := os.Getwd()
	isCurrentDir := err == nil && filepath.Clean(absDir) == filepath.Clean(cwd)
	if isCurrentDir {
		_ = os.RemoveAll(filepath.Join(absDir, ".git"))

		return
	}

	_ = os.RemoveAll(absDir)
}

func runCloneWithCmdClone(slug, absDir string) error {
	return cmdclone.RunClone([]string{slug, absDir})
}

func runCloneWithGH(slug, absDir string) error {
	cmd := exec.Command("gh", "repo", "clone", slug, absDir)

	return cmd.Run()
}

func runCloneWithGit(sshURL, httpURL, absDir string) error {
	cmdSSH := exec.Command("git", "clone", sshURL, absDir)
	if err := cmdSSH.Run(); err == nil {
		return nil
	}

	cmdHTTP := exec.Command("git", "clone", httpURL, absDir)

	return cmdHTTP.Run()
}

func dispatchCloneFallback(absDir, slug string) (string, error) {
	prepareCloneTargetDir(absDir)
	remoteURL := fmt.Sprintf("https://github.com/%s", slug)
	sshURL := fmt.Sprintf("git@github.com:%s.git", slug)

	if err := runCloneWithCmdClone(slug, absDir); err == nil {
		return remoteURL, nil
	}

	if err := runCloneWithGH(slug, absDir); err == nil {
		return remoteURL, nil
	}

	if err := runCloneWithGit(sshURL, remoteURL, absDir); err == nil {
		return remoteURL, nil
	}

	return "", apperror.NewSimple("failed to clone remote repository "+slug, "E1078")
}

func executeForcePushRemote(absDir, slug string) (string, error) {
	sshURL := fmt.Sprintf("git@github.com:%s.git", slug)
	cmdSet := exec.Command("git", "-C", absDir, "remote", "set-url", "origin", sshURL)
	if err := cmdSet.Run(); err != nil {
		cmdAdd := exec.Command("git", "-C", absDir, "remote", "add", "origin", sshURL)
		_ = cmdAdd.Run()
	}

	cmdPush := exec.Command("git", "-C", absDir, "push", "-u", "origin", "main", "--force")
	out, err := cmdPush.CombinedOutput()
	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf("git push --force failed: %s", string(out)))
	}

	return fmt.Sprintf("https://github.com/%s", slug), nil
}
