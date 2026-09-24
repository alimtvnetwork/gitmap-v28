package cmdssh

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// MatchesOSToken checks if a connection's OS matches the given OS token.
func MatchesOSToken(nodeOS string, token string) bool {
	tok := strings.ToLower(strings.TrimSpace(token))
	osLower := strings.ToLower(strings.TrimSpace(nodeOS))
	if tok == "" {
		return false
	}

	switch tok {
	case "win", "windows":
		return osLower == "windows" || osLower == "win"
	case "unix":
		return osLower != "windows" && osLower != "win"
	case "linux", "ubuntu", "debian", "centos", "rhel", "alpine", "arch", "fedora":
		return osLower == "linux" || (osLower == "" && !isWindowsOS(nodeOS))
	case "darwin", "mac", "macos", "osx":
		return osLower == "darwin" || osLower == "mac" || osLower == "macos"
	default:
		return strings.EqualFold(osLower, tok)
	}
}

// IsConnectionOSExcluded reports whether nodeOS matches any token in exceptOSRaw.
func IsConnectionOSExcluded(nodeOS string, exceptOSRaw string) bool {
	tokens := splitExceptTokens(exceptOSRaw)
	for _, tok := range tokens {
		if MatchesOSToken(nodeOS, tok) {
			return true
		}
	}
	return false
}

// IsConnectionOSIncluded reports whether nodeOS matches any token in targetOSRaw.
func IsConnectionOSIncluded(nodeOS string, targetOSRaw string) bool {
	tokens := splitExceptTokens(targetOSRaw)
	if len(tokens) == 0 {
		return true
	}
	for _, tok := range tokens {
		if MatchesOSToken(nodeOS, tok) {
			return true
		}
	}
	return false
}

// FilterSSHConnectionsByOS filters connections by target OS and exclusion OS tokens.
func FilterSSHConnectionsByOS(conns []db.SSHConnection, targetOS, exceptOS string) []db.SSHConnection {
	if targetOS == "" && exceptOS == "" {
		return conns
	}
	var out []db.SSHConnection
	for _, c := range conns {
		if exceptOS != "" && IsConnectionOSExcluded(c.OS, exceptOS) {
			continue
		}
		if targetOS != "" && !IsConnectionOSIncluded(c.OS, targetOS) {
			continue
		}
		out = append(out, c)
	}
	return out
}
