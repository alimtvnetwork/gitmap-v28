// Package ghtoken resolves a GitHub API token from the environment
// or, as a fallback, from the locally installed GitHub CLI (`gh`).
//
// Resolution order:
//  1. GITHUB_TOKEN env var (legacy, CI, scripts).
//  2. GH_TOKEN env var (gh CLI's own override).
//  3. `gh auth token` if `gh` is on PATH AND the user is authenticated.
//
// This means a developer who has `gh auth login`'d once never has to
// set GITHUB_TOKEN by hand — releases, asset uploads, and remote
// repo creation all "just work".
//
// All resolution is best-effort: when no source yields a token, an
// empty string + a typed reason are returned so callers can render
// a user-friendly message instead of a stack trace.
package ghtoken

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Source identifies which mechanism produced a token. Used for
// colorful "✓ Using token from <source>" log lines.
type Source string

const (
	SourceEnvGitHubToken Source = "GITHUB_TOKEN env var"
	SourceEnvGHToken     Source = "GH_TOKEN env var"
	SourceGhCLI          Source = "GitHub CLI (gh auth token)"
	SourceNone           Source = ""
)

// ErrNoToken is returned by Resolve when neither env vars nor the
// gh CLI yield a usable token. Wrapped so callers can branch on
// errors.Is without losing the human-readable detail.
var ErrNoToken = errors.New("no GitHub token available (set GITHUB_TOKEN or run `gh auth login`)")

// Resolve returns (token, source, err). On success err is nil and
// source identifies where the token came from. On failure token is
// "" and err wraps ErrNoToken with a one-line reason.
func Resolve() (string, Source, error) {
	if t := os.Getenv("GITHUB_TOKEN"); len(t) > 0 {
		return t, SourceEnvGitHubToken, nil
	}
	if t := os.Getenv("GH_TOKEN"); len(t) > 0 {
		return t, SourceEnvGHToken, nil
	}
	if t, ok := tokenFromGhCLI(); ok {
		return t, SourceGhCLI, nil
	}

	return "", SourceNone, ErrNoToken
}

// tokenFromGhCLI shells out to `gh auth token`. Returns ("", false)
// when gh is missing, not authenticated, or prints empty output.
// Never panics; never logs — pure helper.
func tokenFromGhCLI() (string, bool) {
	if _, err := exec.LookPath("gh"); err != nil {
		return "", false
	}
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", false
	}
	tok := strings.TrimSpace(string(out))
	if len(tok) == 0 {
		return "", false
	}

	return tok, true
}
