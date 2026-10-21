// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// triggerRemoteCloneIDERegistration invokes remote IDE registration for a successfully cloned repository.
func triggerRemoteCloneIDERegistration(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, out string) []string {
	repoName, remotePath := resolveClonedRepoNameAndPath(conn, opts, out)
	if repoName == "" {
		return nil
	}

	ides, err := RegisterRemoteIDEs(client, conn, repoName, remotePath, false)
	if err != nil {
		return nil
	}

	return ides
}

// resolveClonedRepoNameAndPath extracts repo name and remote repo path from output or clone options.
func resolveClonedRepoNameAndPath(conn db.SSHConnection, opts NodesCloneOptions, out string) (string, string) {
	if payload, hasJSON := cmdclone.ParseCloneJSONResponse(out); hasJSON && payload.RepoName != "" {
		remotePath := resolvePayloadPath(conn, opts, payload)

		return payload.RepoName, remotePath
	}

	return resolveFallbackRepoNameAndPath(conn, opts)
}

func resolvePayloadPath(conn db.SSHConnection, opts NodesCloneOptions, payload cmdclone.DirectCloneJSONResponse) string {
	if payload.Path != "" {
		return payload.Path
	}

	return resolveRemoteRepoPath(conn, opts, payload.RepoName)
}

func resolveFallbackRepoNameAndPath(conn db.SSHConnection, opts NodesCloneOptions) (string, string) {
	repoName := extractRepoNameFromArgs(opts.PassArgs)
	if repoName == "" {
		return "", ""
	}

	remotePath := resolveRemoteRepoPath(conn, opts, repoName)

	return repoName, remotePath
}

func resolveRemoteRepoPath(conn db.SSHConnection, opts NodesCloneOptions, repoName string) string {
	targetDir := resolveRemoteTargetDir(conn, opts)
	if isWindowsNode(conn) {
		return targetDir + `\` + repoName
	}

	return strings.TrimRight(targetDir, "/") + "/" + repoName
}

func extractRepoNameFromArgs(args []string) string {
	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed == "" || strings.HasPrefix(trimmed, "-") {
			continue
		}

		return cleanRepoNameToken(trimmed)
	}

	return ""
}

func cleanRepoNameToken(token string) string {
	token = strings.TrimSuffix(token, ".git")
	token = strings.TrimRight(token, "/\\")
	slashToken := filepath.ToSlash(token)
	parts := strings.Split(slashToken, "/")

	return parts[len(parts)-1]
}
