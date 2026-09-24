package cmdssh

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
)

func resolveCloneTargetAndPath(opts sshCloneOptions) (string, string, string, error) {
	rawRepo, err := resolveRawRepo(opts.Repo)
	if err != nil {
		return "", "", "", err
	}

	repoURL := expandRepoTokenToURL(rawRepo)
	repoURL = applyCloneProtocol(repoURL, opts.IsSSH, opts.IsHTTPS)
	repoName := resolveRepoName(repoURL)
	return repoURL, repoName, opts.DestPath, nil
}

func resolveRawRepo(rawRepo string) (string, error) {
	if rawRepo != "" && rawRepo != "." && rawRepo != "git" {
		return rawRepo, nil
	}
	curURL, err := detectCurrentRepoURL()
	if err != nil {
		return "", apperror.NewValidationError("missing repository name or URL; not in a git repo")
	}
	return curURL, nil
}

func applyCloneProtocol(repoURL string, isSSH, isHTTPS bool) string {
	if isSSH {
		return convertSSHProtocol(repoURL)
	}
	if isHTTPS {
		return convertHTTPSProtocol(repoURL)
	}
	return repoURL
}

func convertSSHProtocol(repoURL string) string {
	sshURL, isConverted := cmdclone.ConvertURLToSSH(repoURL)
	if isConverted {
		return sshURL
	}
	return repoURL
}

func convertHTTPSProtocol(repoURL string) string {
	httpsURL, isConverted := cmdclone.ConvertURLToHTTPS(repoURL)
	if isConverted {
		return httpsURL
	}
	return repoURL
}

func expandRepoTokenToURL(token string) string {
	low := strings.ToLower(token)
	if strings.HasPrefix(low, "http://") || strings.HasPrefix(low, "https://") ||
		strings.HasPrefix(low, "git@") || strings.HasPrefix(low, "ssh://") {
		return token
	}

	if strings.Contains(token, "/") {
		return "https://github.com/" + strings.TrimSuffix(token, ".git") + ".git"
	}

	org := detectDefaultGitOrg()
	return "https://github.com/" + org + "/" + strings.TrimSuffix(token, ".git") + ".git"
}

func detectCurrentRepoURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return "", apperror.NewValidationError("empty origin url")
	}
	return trimmed, nil
}

func detectDefaultGitOrg() string {
	curURL, err := detectCurrentRepoURL()
	if err != nil || curURL == "" {
		return "alimtvnetwork"
	}
	parts := strings.Split(strings.TrimSuffix(curURL, ".git"), "/")
	if len(parts) < 2 {
		return "alimtvnetwork"
	}
	return parts[len(parts)-2]
}

func resolveRepoName(repoURL string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(repoURL), ".git")
	clean = strings.TrimRight(clean, "/\\")
	colonIdx := strings.LastIndex(clean, ":")
	slashIdx := strings.LastIndex(clean, "/")
	splitIdx := slashIdx
	if colonIdx > splitIdx {
		splitIdx = colonIdx
	}
	if splitIdx >= 0 && splitIdx < len(clean)-1 {
		return clean[splitIdx+1:]
	}
	return filepath.Base(clean)
}

func resolveRemoteDestPathForNode(repoName, rawDest, osType string) string {
	isWin := isWindowsOS(osType)
	if rawDest == "" || rawDest == "git" {
		return resolveDefaultGitPath(repoName, isWin)
	}

	expanded := ExpandUniversalPath(rawDest, osType)
	if isDirDestination(rawDest) {
		return joinRemotePath(expanded, repoName, isWin)
	}
	return expanded
}

func resolveDefaultGitPath(repoName string, isWin bool) string {
	if isWin {
		return "~\\git\\" + repoName
	}
	return "~/git/" + repoName
}
