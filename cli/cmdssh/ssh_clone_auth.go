package cmdssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func isAuthFailure(out string, err error) bool {
	combined := strings.ToLower(out)
	if err != nil {
		combined += " " + strings.ToLower(err.Error())
	}
	return strings.Contains(combined, "permission denied") ||
		strings.Contains(combined, "publickey") ||
		strings.Contains(combined, "repository not found") ||
		strings.Contains(combined, "could not read username") ||
		strings.Contains(combined, "terminal prompts disabled") ||
		strings.Contains(combined, "authentication failed") ||
		strings.Contains(combined, "401") ||
		strings.Contains(combined, "403") ||
		strings.Contains(combined, "host key verification failed")
}

func extractGitHost(repoURL string) string {
	clean := strings.TrimPrefix(repoURL, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimPrefix(clean, "ssh://")
	clean = strings.TrimPrefix(clean, "git@")
	slashIdx := strings.Index(clean, "/")
	colonIdx := strings.Index(clean, ":")
	splitIdx := slashIdx
	if colonIdx >= 0 && (splitIdx < 0 || colonIdx < splitIdx) {
		splitIdx = colonIdx
	}
	if splitIdx > 0 {
		return clean[:splitIdx]
	}
	return "github.com"
}

func executeSelfHealingAuth(client *ssh.Client, c db.SSHConnection, repoURL, header string) error {
	fmt.Printf("  %s %s[auth-probe] Repository access denied on node %s. Deploying authentication keys...%s\n",
		header, constants.ColorYellow, c.Alias, constants.ColorReset)

	pubKey, keyPath, appErr := discoverPublicKey("")
	if appErr == nil && pubKey != "" {
		executeKeyInjection(client, header, c.OS, pubKey)
		fmt.Printf("  %s %s[auth-heal] Step 1/3: Deployed local public key (%s) to node %s%s\n",
			header, constants.ColorCyan, keyPath, c.Alias, constants.ColorReset)
	}

	gitHost := extractGitHost(repoURL)
	configureRemoteGitTrust(client, c.OS, gitHost)
	fmt.Printf("  %s %s[auth-heal] Step 2/3: Configured Git host '%s' trust & deployment keys on node %s%s\n",
		header, constants.ColorCyan, gitHost, c.Alias, constants.ColorReset)

	fmt.Printf("  %s %s[auth-heal] Step 3/3: Re-probing repository accessibility...%s\n",
		header, constants.ColorCyan, constants.ColorReset)
	if probeRemoteRepoAccess(client, c.OS, repoURL) {
		fmt.Printf("  %s %s[auth-heal] Authentication self-healed successfully on node %s!%s\n",
			header, constants.ColorGreen, c.Alias, constants.ColorReset)
		return nil
	}

	return apperror.NewExecutionError(fmt.Sprintf("repository access probe failed for %s on node %s after self-healing", repoURL, c.Alias))
}

func configureRemoteGitTrust(client *ssh.Client, osType, gitHost string) {
	shell := resolveRemoteShell(osType)
	var cmd string
	if isWindowsOS(osType) {
		cmd = fmt.Sprintf(`powershell -NoProfile -Command "$dir=Join-Path $env:USERPROFILE '.ssh'; if (-not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }; if (-not (Test-Path (Join-Path $dir 'id_ed25519'))) { ssh-keygen -t ed25519 -N '\"\"' -f (Join-Path $dir 'id_ed25519') -q }; ssh-keyscan -H %s 2>$null | Out-File -Append -Encoding ascii (Join-Path $dir 'known_hosts')"`, gitHost)
	} else {
		cmd = fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && (test -f ~/.ssh/id_ed25519 || ssh-keygen -t ed25519 -N '' -f ~/.ssh/id_ed25519 -q) && ssh-keyscan -H %s >> ~/.ssh/known_hosts 2>/dev/null", gitHost)
	}
	_, _ = crypto.RunCommand(client, cmd, shell)
}

func probeRemoteRepoAccess(client *ssh.Client, osType, repoURL string) bool {
	probeCmd := fmt.Sprintf("git ls-remote --exit-code -h '%s' HEAD", repoURL)
	shell := resolveRemoteShell(osType)
	out, err := crypto.RunCommand(client, probeCmd, shell)
	return err == nil && !isAuthFailure(out, nil)
}
