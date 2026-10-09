// Package cmdssh — secrets_resolver.go delegates to secretsresolver package.
package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// ResolveRepoSecretsRoot locates the canonical repo-secrets root directory.
func ResolveRepoSecretsRoot() string {
	return secrets.ResolveRepoSecretsRoot()
}

// ResolveRepoSecretsNodesPath locates gitmap-ssh-nodes.json or gitmap-ssh.json for a token or default.
func ResolveRepoSecretsNodesPath(token string) string {
	return secrets.ResolveRepoSecretsNodesPath(token)
}

// ResolveRepoSecretsManifest locates a clone or status manifest (gitmap.json) for a machine token.
func ResolveRepoSecretsManifest(token string, preferredFile string) string {
	return secrets.ResolveRepoSecretsManifest(token, preferredFile)
}
