package cmdlemon

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon/enums"

// WhitelistAuthType aliases the strongly-typed enum from the enums package.
type WhitelistAuthType = enums.WhitelistAuthType

const (
	AuthAllowedType = enums.AuthAllowedType
	AuthDeniedType  = enums.AuthDeniedType
)

const (
	DefaultTimeout, MaxTimeout     = 60 * time.Second, 300 * time.Second
	MaxBufferBytes, dangerousChars = 64 * 1024, ";|&><$`\r\n"
)

var allowedSubcommands = map[string]bool{
	"status": true, "pull": true, "pe": true, "version": true, "nodes": true,
	"agm": true, "agy": true, "space": true, "scan": true,
}

func IsSubcommandAllowed(subcmd string) bool {
	return allowedSubcommands[strings.ToLower(strings.TrimSpace(subcmd))]
}

func IsSafeToken(token string) bool {
	return !strings.ContainsAny(token, dangerousChars)
}

func ValidateTokens(tokens []string) ([]string, bool) {
	clean := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		t := strings.TrimSpace(tok)
		if t == "" {
			continue
		}
		if !IsSafeToken(t) {
			return nil, false
		}
		clean = append(clean, t)
	}

	if len(clean) == 0 {
		return nil, false
	}

	idx := 0
	if strutil.EqualFoldAny(clean[0], "gitmap", "gitmap.exe") {
		idx = 1
	}
	if idx >= len(clean) || !IsSubcommandAllowed(strings.TrimPrefix(strings.ToLower(clean[idx]), "/")) {
		return nil, false
	}

	clean[idx] = strings.TrimPrefix(strings.ToLower(clean[idx]), "/")

	return clean[idx:], true
}

func ParseAndAuthorizeCommand(raw string) ([]string, error) {
	tokens := strings.Fields(strings.TrimSpace(raw))
	args, isAllowed := ValidateTokens(tokens)
	if !isAllowed {
		return nil, errors.New("command not permitted by security whitelist")
	}

	return args, nil
}

func ResolveTimeout(req time.Duration) time.Duration {
	if req > MaxTimeout {
		return MaxTimeout
	}
	if req > 0 {
		return req
	}

	return DefaultTimeout
}

func ResolveExitCode(err error) int {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	if err != nil {
		return 1
	}

	return 0
}

func TruncateBuffer(buf []byte, maxLen int) (string, bool) {
	if len(buf) > maxLen {
		return string(buf[:maxLen]) + "\n...[output truncated]", true
	}
	return string(buf), false
}
