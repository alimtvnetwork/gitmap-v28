package cmdlemon

import (
	"errors"
	"os/exec"
	"strings"
	"time"
)

// WhitelistAuthType specifies the authorization state.
type WhitelistAuthType string

const (
	AuthAllowedType, AuthDeniedType WhitelistAuthType = "ALLOWED", "DENIED"
	DefaultTimeout                  time.Duration     = 60 * time.Second
	MaxTimeout                      time.Duration     = 300 * time.Second
	MaxBufferBytes                                    = 64 * 1024
	dangerousChars                                    = ";|&><$`\r\n"
)

var allowedSubcommands = map[string]bool{
	"status": true, "pull": true, "pe": true,
	"version": true, "nodes": true, "agm": true,
	"agy": true, "space": true, "scan": true,
}

// IsSubcommandAllowed verifies if a subcommand is whitelisted.
func IsSubcommandAllowed(subcmd string) bool {
	return allowedSubcommands[strings.ToLower(strings.TrimSpace(subcmd))]
}

// IsSafeToken verifies a token has no shell injection characters.
func IsSafeToken(token string) bool {
	return !strings.ContainsAny(token, dangerousChars)
}

// ValidateTokens ensures all tokens are injection-free and starts with allowed subcmd.
func ValidateTokens(tokens []string) ([]string, bool) {
	clean := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		t := strings.TrimSpace(tok)
		if t != "" && IsSafeToken(t) {
			clean = append(clean, t)
		} else if t != "" {
			return nil, false
		}
	}
	if len(clean) == 0 {
		return nil, false
	}
	idx := 0
	if first := strings.ToLower(clean[0]); first == "gitmap" || first == "gitmap.exe" {
		idx = 1
	}
	if idx >= len(clean) || !IsSubcommandAllowed(strings.TrimPrefix(strings.ToLower(clean[idx]), "/")) {
		return nil, false
	}
	clean[idx] = strings.TrimPrefix(strings.ToLower(clean[idx]), "/")
	return clean[idx:], true
}

// ParseAndAuthorizeCommand tokenizes and authorizes a raw command.
func ParseAndAuthorizeCommand(raw string) ([]string, error) {
	tokens := strings.Fields(strings.TrimSpace(raw))
	args, isAllowed := ValidateTokens(tokens)
	if !isAllowed {
		return nil, errors.New("command not permitted by security whitelist")
	}
	return args, nil
}

// ResolveTimeout clamps timeout between DefaultTimeout and MaxTimeout.
func ResolveTimeout(req time.Duration) time.Duration {
	if req > MaxTimeout {
		return MaxTimeout
	} else if req > 0 {
		return req
	}
	return DefaultTimeout
}

// ResolveExitCode extracts the exit code from an execution error.
func ResolveExitCode(err error) int {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	} else if err != nil {
		return 1
	}
	return 0
}

// TruncateBuffer caps output buffers at maxLen bytes.
func TruncateBuffer(buf []byte, maxLen int) (string, bool) {
	if len(buf) > maxLen {
		return string(buf[:maxLen]) + "\n...[output truncated]", true
	}
	return string(buf), false
}
