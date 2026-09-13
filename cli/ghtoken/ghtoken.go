package ghtoken

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SourceType identifies which mechanism produced a token.
type SourceType string

const (
	SourceEnvGitHubToken  SourceType = "GITHUB_TOKEN env var"
	SourceEnvGHToken      SourceType = "GH_TOKEN env var"
	SourceWinRegistryUser SourceType = "Windows User Registry"
	SourceWinRegistrySys  SourceType = "Windows System Registry"
	SourceGitCredential   SourceType = "Git Credential Manager"
	SourceGhCLI           SourceType = "GitHub CLI (gh auth token)"
	SourceGhHostsFile     SourceType = "GitHub CLI config (hosts.yml)"
	SourceGitConfig       SourceType = "git config (github.token)"
	SourceNone            SourceType = ""
)

// ErrNoToken is returned when no source yields a usable token.
var ErrNoToken = errors.New("no GitHub token available (set GH_TOKEN, login via `gh auth login`, or use git credentials)")

// Resolve dynamically discovers and returns a GitHub token from the system.
func Resolve() (string, SourceType, error) {
	if tok, src, isDefined := resolveEnvToken(); isDefined {
		return tok, src, nil
	}

	if tok, src, isDefined := resolveSystemToken(); isDefined {
		return tok, src, nil
	}

	if tok, src, isDefined := resolveToolToken(); isDefined {
		return tok, src, nil
	}

	if tok, src, isDefined := resolveConfigToken(); isDefined {
		return tok, src, nil
	}

	return "", SourceNone, ErrNoToken
}

func resolveEnvToken() (string, SourceType, bool) {
	if t := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); len(t) > 0 {
		return t, SourceEnvGitHubToken, true
	}

	if t := strings.TrimSpace(os.Getenv("GH_TOKEN")); len(t) > 0 {
		return t, SourceEnvGHToken, true
	}

	return "", SourceNone, false
}

func resolveSystemToken() (string, SourceType, bool) {
	return tokenFromSystemRegistry()
}

func resolveToolToken() (string, SourceType, bool) {
	if t, isDefined := tokenFromGitCredential(); isDefined {
		return t, SourceGitCredential, true
	}

	if t, isDefined := tokenFromGhCLI(); isDefined {
		return t, SourceGhCLI, true
	}

	return "", SourceNone, false
}

func resolveConfigToken() (string, SourceType, bool) {
	if t, isDefined := tokenFromHostsFile(); isDefined {
		return t, SourceGhHostsFile, true
	}

	if t, isDefined := tokenFromGitConfig(); isDefined {
		return t, SourceGitConfig, true
	}

	return "", SourceNone, false
}

func tokenFromGhCLI() (string, bool) {
	if _, err := exec.LookPath("gh"); err != nil {
		return "", false
	}

	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", false
	}

	tok := strings.TrimSpace(string(out))
	isDefined := len(tok) > 0

	return tok, isDefined
}

func tokenFromGitCredential() (string, bool) {
	cmd := exec.Command("git", "credential", "fill")
	cmd.Stdin = strings.NewReader("protocol=https\nhost=github.com\n\n")

	out, err := cmd.Output()
	if err != nil {
		return "", false
	}

	return extractPasswordFromCredentialOutput(string(out))
}

func extractPasswordFromCredentialOutput(out string) (string, bool) {
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "password=") {
			tok := strings.TrimPrefix(trimmed, "password=")
			isDefined := len(tok) > 0

			return tok, isDefined
		}
	}

	return "", false
}

func tokenFromGitConfig() (string, bool) {
	out, err := exec.Command("git", "config", "--get", "github.token").Output()
	if err != nil {
		return "", false
	}

	tok := strings.TrimSpace(string(out))
	isDefined := len(tok) > 0

	return tok, isDefined
}

func tokenFromHostsFile() (string, bool) {
	path := getGhHostsFilePath()
	if len(path) == 0 {
		return "", false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	return extractOauthTokenFromYAML(string(data))
}

func getGhHostsFilePath() string {
	if appData := os.Getenv("APPDATA"); len(appData) > 0 {
		winPath := filepath.Join(appData, "GitHub CLI", "hosts.yml")
		if _, err := os.Stat(winPath); err == nil {
			return winPath
		}
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		unixPath := filepath.Join(homeDir, ".config", "gh", "hosts.yml")
		if _, err := os.Stat(unixPath); err == nil {
			return unixPath
		}
	}

	return ""
}

func extractOauthTokenFromYAML(content string) (string, bool) {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "oauth_token:") {
			tok := strings.TrimSpace(strings.TrimPrefix(trimmed, "oauth_token:"))
			tok = strings.Trim(tok, `"'`)
			isDefined := len(tok) > 0

			return tok, isDefined
		}
	}

	return "", false
}
