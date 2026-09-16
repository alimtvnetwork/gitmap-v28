package cloner

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func rewriteRemoteToSSHIfNeeded(repoDir string, isUseSSH bool) {
	if !isUseSSH {
		return
	}
	url := queryRemoteOriginURL(repoDir)
	sshURL := toSSHURL(url)
	if sshURL != url && sshURL != "" {
		setRemoteOriginURL(repoDir, sshURL)
	}
}

func queryRemoteOriginURL(repoDir string) string {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func setRemoteOriginURL(repoDir, sshURL string) {
	setCmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, "remote", "set-url", "origin", sshURL)
	_ = setCmd.Run()
}

func toSSHURL(raw string) string {
	if !strings.HasPrefix(raw, constants.PrefixHTTPS) {
		return raw
	}
	trimmed := strings.TrimPrefix(raw, constants.PrefixHTTPS)
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 2 {
		return "git@" + parts[0] + ":" + parts[1]
	}

	return raw
}
