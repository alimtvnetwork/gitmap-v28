// Package cmdssh — secrets_resolver.go delegates to secretsresolver package.
package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/secretsresolver"
)

// ResolveRepoSecretsRoot locates the canonical repo-secrets root directory.
func ResolveRepoSecretsRoot() string {
	return secretsresolver.ResolveRepoSecretsRoot()
}

// ResolveRepoSecretsNodesPath locates gitmap-ssh-nodes.json or gitmap-ssh.json for a token or default.
func ResolveRepoSecretsNodesPath(token string) string {
	return secretsresolver.ResolveRepoSecretsNodesPath(token)
}

// ResolveRepoSecretsManifest locates a clone or status manifest (gitmap.json) for a machine token.
func ResolveRepoSecretsManifest(token string, preferredFile string) string {
	return secretsresolver.ResolveRepoSecretsManifest(token, preferredFile)
}
