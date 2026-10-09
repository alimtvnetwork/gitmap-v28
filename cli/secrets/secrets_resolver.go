// Package secretsresolver provides path resolution for repo-secrets machine folders and manifests.
package secrets

import (
	"os"
	"path/filepath"
	"strings"
)

// Default filenames
const (
	DefaultSSHNodesJSONFile    = "gitmap-ssh-nodes.json"
	DefaultSSHNodesAltJSONFile = "gitmap-ssh.json"
	DefaultGitMapJSONFile      = "gitmap.json"
)

var knownMachineFolderMap = map[string]string{
	"w1":        "04-w1-machine",
	"w2":        "05-w2-machine",
	"w3":        "06-w3-machine",
	"w4":        "04-w4-machine",
	"01":        "01-gitmap",
	"01-gitmap": "01-gitmap",
	"gitmap":    "01-gitmap",
	"nodes":     "01-gitmap",
}

// ResolveRepoSecretsRoot locates the canonical repo-secrets root directory.
func ResolveRepoSecretsRoot() string {
	candidates := []string{
		`D:\work\repo-secrets`,
		filepath.Join("..", "repo-secrets"),
		"repo-secrets",
	}
	for _, c := range candidates {
		if isDirectoryPath(c) {
			return c
		}
	}
	return `D:\work\repo-secrets`
}

// ResolveRepoSecretsNodesPath locates gitmap-ssh-nodes.json or gitmap-ssh.json for a token or default.
func ResolveRepoSecretsNodesPath(token string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed != "" && isDirectPath(trimmed) {
		return trimmed
	}
	root := ResolveRepoSecretsRoot()
	if folderPath := resolveNonEmptyFolderToken(root, trimmed); folderPath != "" {
		return findNodesFileInDir(folderPath)
	}
	return findDefaultNodesPath(root)
}

func resolveNonEmptyFolderToken(root, trimmed string) string {
	if trimmed == "" {
		return ""
	}
	return resolveFolderByToken(root, trimmed)
}

// ResolveRepoSecretsManifest locates a clone or status manifest (gitmap.json) for a machine token.
func ResolveRepoSecretsManifest(token string, preferredFile string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed != "" && isDirectPath(trimmed) {
		return trimmed
	}
	root := ResolveRepoSecretsRoot()
	targetDir := resolveFolderByToken(root, trimmed)
	if targetDir == "" {
		return ""
	}
	if preferredPath := resolvePreferredManifestPath(targetDir, preferredFile); preferredPath != "" {
		return preferredPath
	}
	return findAnyManifestInDir(targetDir)
}

func resolvePreferredManifestPath(targetDir, preferredFile string) string {
	if preferredFile == "" {
		return ""
	}
	preferredPath := filepath.Join(targetDir, preferredFile)
	if fileExists(preferredPath) {
		return preferredPath
	}
	return ""
}

func isDirectPath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isDirectoryPath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func resolveFolderByToken(root, token string) string {
	low := strings.ToLower(token)
	if mapped, ok := knownMachineFolderMap[low]; ok && isDirectoryPath(filepath.Join(root, mapped)) {
		return filepath.Join(root, mapped)
	}
	directPath := filepath.Join(root, token)
	if isDirectoryPath(directPath) {
		return directPath
	}
	return scanSubdirMatchingToken(root, low)
}

func scanSubdirMatchingToken(root, low string) string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		nameLow := strings.ToLower(e.Name())
		if matchesMachineFolderName(nameLow, low) {
			return filepath.Join(root, e.Name())
		}
	}
	return ""
}

func matchesMachineFolderName(folderName, token string) bool {
	return strings.Contains(folderName, "-"+token+"-") ||
		strings.HasSuffix(folderName, "-"+token) ||
		strings.HasPrefix(folderName, token+"-") ||
		strings.Contains(folderName, token)
}

func findNodesFileInDir(dir string) string {
	candidates := []string{
		filepath.Join(dir, DefaultSSHNodesJSONFile),
		filepath.Join(dir, DefaultSSHNodesAltJSONFile),
	}
	for _, c := range candidates {
		if fileExists(c) {
			return c
		}
	}
	return filepath.Join(dir, DefaultSSHNodesJSONFile)
}

func findAnyManifestInDir(dir string) string {
	candidates := []string{
		filepath.Join(dir, DefaultGitMapJSONFile),
		filepath.Join(dir, "gitmap.csv"),
	}
	for _, c := range candidates {
		if fileExists(c) {
			return c
		}
	}
	return ""
}

func findDefaultNodesPath(root string) string {
	localCandidates := []string{
		DefaultSSHNodesJSONFile,
		DefaultSSHNodesAltJSONFile,
		filepath.Join(root, "01-gitmap", DefaultSSHNodesJSONFile),
		filepath.Join(root, "04-w1-machine", DefaultSSHNodesJSONFile),
		filepath.Join(root, DefaultSSHNodesJSONFile),
	}
	for _, c := range localCandidates {
		if fileExists(c) {
			return c
		}
	}
	return DefaultSSHNodesJSONFile
}
