package usercontext

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

var (
	runCommandOutput = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		cmd := execCommandContext(ctx, name, args...)
		return cmd.CombinedOutput()
	}
	execCommandContext = exec.CommandContext
	lookPath           = exec.LookPath
	tokenResolver      = secrets.Resolve
)

const (
	defaultAuthTimeout = 2 * time.Second
	sourceGhCLI        = "gh-cli"
	fallbackUser       = "github-user"
)

// DetectGhAuth inspects GitHub CLI status and falls back to system tokens.
func DetectGhAuth() GhAuthInfo {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAuthTimeout)
	defer cancel()
	if !isGhInstalled() {
		return resolveWithoutGh(ctx)
	}
	info, isAuthed := probeGhCLI(ctx)
	if isAuthed {
		return info
	}
	return resolveTokenFallback(ctx)
}

func isGhInstalled() bool {
	_, err := lookPath("gh")
	return err == nil
}

func probeGhCLI(ctx context.Context) (GhAuthInfo, bool) {
	info, hasStatus := probeGhAuthStatus(ctx)
	if hasStatus {
		return info, true
	}
	return probeGhApiUser(ctx)
}

func probeGhAuthStatus(ctx context.Context) (GhAuthInfo, bool) {
	out, _ := runCommandOutput(ctx, "gh", "auth", "status")
	text := string(out)
	username := extractUsernameFromAuthStatus(text)
	if len(username) == 0 {
		return GhAuthInfo{}, false
	}
	return newGhAuthInfo(username, extractScopeFromAuthStatus(text)), true
}

func newGhAuthInfo(username, scope string) GhAuthInfo {
	return GhAuthInfo{
		Username: username,
		Status:   StatusAuthorized,
		Source:   sourceGhCLI,
		Scope:    scope,
	}
}

func extractUsernameFromAuthStatus(text string) string {
	lines := strings.Split(text, "\n")
	for _, raw := range lines {
		user := parseAuthStatusLine(raw)
		if len(user) > 0 {
			return user
		}
	}
	return ""
}

func parseAuthStatusLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if idx := strings.Index(trimmed, "account "); idx != -1 {
		return extractAccountName(trimmed[idx+8:])
	}
	if idx := strings.Index(trimmed, "as "); idx != -1 {
		return extractAccountName(trimmed[idx+3:])
	}
	return ""
}

func extractAccountName(raw string) string {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return ""
	}
	name := strings.Trim(fields[0], "()[]'\",")
	return strings.TrimSpace(name)
}

func extractScopeFromAuthStatus(text string) string {
	lines := strings.Split(text, "\n")
	for _, raw := range lines {
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "Token scopes:") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "Token scopes:"))
		}
	}
	return ""
}

func probeGhApiUser(ctx context.Context) (GhAuthInfo, bool) {
	out, err := runCommandOutput(ctx, "gh", "api", "user", "--jq", ".login")
	if err != nil {
		return GhAuthInfo{}, false
	}
	user := strings.TrimSpace(string(out))
	if len(user) == 0 {
		return GhAuthInfo{}, false
	}
	return newGhAuthInfo(user, ""), true
}

func resolveWithoutGh(ctx context.Context) GhAuthInfo {
	info, hasToken := resolveTokenFallbackCore(ctx)
	if hasToken {
		return info
	}
	return GhAuthInfo{
		Status: StatusToolMissing,
	}
}

func resolveTokenFallback(ctx context.Context) GhAuthInfo {
	info, hasToken := resolveTokenFallbackCore(ctx)
	if hasToken {
		return info
	}
	return GhAuthInfo{
		Status: StatusUnauthenticated,
	}
}

func resolveTokenFallbackCore(ctx context.Context) (GhAuthInfo, bool) {
	tok, src, err := tokenResolver()
	if err != nil || len(tok) == 0 {
		return GhAuthInfo{}, false
	}
	username := resolveFallbackUsername(ctx)
	return GhAuthInfo{
		Username: username,
		Status:   StatusAuthorized,
		Source:   string(src),
	}, true
}

func resolveFallbackUsername(ctx context.Context) string {
	out, err := runCommandOutput(ctx, "git", "config", "user.name")
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}
	return fallbackUser
}
