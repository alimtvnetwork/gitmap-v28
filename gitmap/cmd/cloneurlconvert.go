package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ConvertURLToSSH rewrites a Git remote URL into its `git@host:owner/repo.git`
// SSH-shorthand form. Inputs that already look like SSH are normalized
// (`.git` suffix appended) but otherwise returned as-is. Inputs that are
// not a recognised Git URL shape are returned unchanged with ok=false so
// the caller can decide whether to abort or fall through.
//
// Supported input shapes:
//
//	https://host/owner/repo(.git)?(/)?
//	http://host/owner/repo(.git)?(/)?
//	ssh://git@host[:port]/owner/repo(.git)?
//	git@host:owner/repo(.git)?
//
// Spec: spec/01-app/110-clone-ssh-flag.md
func ConvertURLToSSH(url string) (string, bool) {
	trimmed := strings.TrimRight(strings.TrimSpace(url), "/\\")
	if trimmed == "" {
		return url, false
	}

	lower := strings.ToLower(trimmed)

	// SSH shorthand already.
	if strings.HasPrefix(lower, "git@") {
		return ensureGitSuffix(trimmed), true
	}

	// ssh:// scheme — strip scheme + optional `git@` userinfo, then split
	// host[:port]/path into host + path so we can re-emit shorthand.
	if strings.HasPrefix(lower, "ssh://") {
		return sshSchemeToShorthand(trimmed)
	}

	// https:// or http://.
	if strings.HasPrefix(lower, constants.PrefixHTTPS) || strings.HasPrefix(lower, "http://") {
		return httpsToSSHShorthand(trimmed)
	}

	return url, false
}

// ConvertURLToHTTPS rewrites a Git remote URL into its
// `https://host/owner/repo.git` form. Symmetric counterpart to
// ConvertURLToSSH; intended for callers that want to force HTTPS even
// when the source manifest captured an SSH URL.
func ConvertURLToHTTPS(url string) (string, bool) {
	trimmed := strings.TrimRight(strings.TrimSpace(url), "/\\")
	if trimmed == "" {
		return url, false
	}

	lower := strings.ToLower(trimmed)

	if strings.HasPrefix(lower, constants.PrefixHTTPS) {
		return ensureGitSuffix(trimmed), true
	}
	if strings.HasPrefix(lower, "http://") {
		return ensureGitSuffix(constants.PrefixHTTPS + trimmed[len("http://"):]), true
	}
	if strings.HasPrefix(lower, "git@") {
		return shorthandToHTTPS(trimmed)
	}
	if strings.HasPrefix(lower, "ssh://") {
		return sshSchemeToHTTPS(trimmed)
	}

	return url, false
}

// ensureGitSuffix appends `.git` to a URL/shorthand when it isn't
// already present. Keeps query strings and fragments out of scope —
// Git clone URLs in this project never carry them.
func ensureGitSuffix(s string) string {
	if strings.HasSuffix(s, ".git") {
		return s
	}

	return s + ".git"
}

// httpsToSSHShorthand converts `https://host/owner/repo` to
// `git@host:owner/repo.git`. Returns ok=false when the path component
// is missing (`https://host` with no owner/repo).
func httpsToSSHShorthand(httpsURL string) (string, bool) {
	rest := strings.TrimPrefix(httpsURL, constants.PrefixHTTPS)
	rest = strings.TrimPrefix(rest, "http://")

	slash := strings.Index(rest, "/")
	if slash <= 0 || slash == len(rest)-1 {
		return httpsURL, false
	}

	host := rest[:slash]
	path := strings.TrimPrefix(rest[slash+1:], "/")
	if path == "" {
		return httpsURL, false
	}

	return ensureGitSuffix("git@" + host + ":" + path), true
}

// shorthandToHTTPS converts `git@host:owner/repo(.git)?` to
// `https://host/owner/repo.git`.
func shorthandToHTTPS(shorthand string) (string, bool) {
	at := strings.Index(shorthand, "@")
	colon := strings.Index(shorthand[at+1:], ":")
	if at < 0 || colon <= 0 {
		return shorthand, false
	}

	host := shorthand[at+1 : at+1+colon]
	path := shorthand[at+1+colon+1:]
	if host == "" || path == "" {
		return shorthand, false
	}

	return ensureGitSuffix(constants.PrefixHTTPS + host + "/" + path), true
}

// sshSchemeToShorthand converts `ssh://git@host[:port]/owner/repo` to
// the shorter `git@host:owner/repo.git` form. Port hints are dropped —
// SSH shorthand does not carry them; users that need a non-default port
// should set it in `~/.ssh/config` or stick with the explicit scheme.
func sshSchemeToShorthand(sshURL string) (string, bool) {
	rest := strings.TrimPrefix(sshURL, "ssh://")
	rest = strings.TrimPrefix(rest, "SSH://")

	if at := strings.Index(rest, "@"); at >= 0 {
		rest = rest[at+1:]
	}

	slash := strings.Index(rest, "/")
	if slash <= 0 || slash == len(rest)-1 {
		return sshURL, false
	}

	hostPart := rest[:slash]
	path := strings.TrimPrefix(rest[slash+1:], "/")

	if colon := strings.Index(hostPart, ":"); colon > 0 {
		hostPart = hostPart[:colon]
	}
	if hostPart == "" || path == "" {
		return sshURL, false
	}

	return ensureGitSuffix("git@" + hostPart + ":" + path), true
}

// sshSchemeToHTTPS converts `ssh://git@host[:port]/owner/repo` to
// `https://host/owner/repo.git`. Reuses sshSchemeToShorthand then
// flips shorthand → HTTPS so both code paths share one parser.
func sshSchemeToHTTPS(sshURL string) (string, bool) {
	short, ok := sshSchemeToShorthand(sshURL)
	if !ok {
		return sshURL, false
	}

	return shorthandToHTTPS(short)
}
