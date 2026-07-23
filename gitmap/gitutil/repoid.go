// Package gitutil — repoid.go derives a stable, transport-neutral
// identifier from a Git remote URL.
//
// The identifier is the canonical "host/owner/repo" form, lowercased,
// with the .git suffix and trailing slashes stripped. Two URLs that
// point at the same repository over different transports (HTTPS vs
// SSH vs scp-style) collapse to the same string, which makes the
// value safe to use as a deterministic key for re-cloning, dedup,
// and cross-export joins.
//
// Examples:
//
//	https://github.com/acme/widget.git    -> github.com/acme/widget
//	git@github.com:acme/widget.git        -> github.com/acme/widget
//	ssh://git@github.com/acme/widget      -> github.com/acme/widget
//
// Unfamiliar URL shapes are returned lowercased + trimmed but
// otherwise untouched, so callers can still equality-compare without
// risking a silent false-match against a different repo.
package gitutil

import "strings"

// CanonicalRepoID collapses an https / ssh / scp-style git URL down
// to a "host/owner/repo" string suitable for equality comparison and
// stable cross-format identity. Empty input returns "".
func CanonicalRepoID(raw string) string {
	s := strings.TrimSpace(raw)
	if len(s) == 0 {
		return ""
	}
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, ".git")
	switch {
	case strings.HasPrefix(s, "https://"):
		s = strings.TrimPrefix(s, "https://")
	case strings.HasPrefix(s, "http://"):
		s = strings.TrimPrefix(s, "http://")
	case strings.HasPrefix(s, "ssh://"):
		s = strings.TrimPrefix(s, "ssh://")
		s = strings.TrimPrefix(s, "git@")
	default:
		if at := strings.Index(s, "@"); at >= 0 {
			// scp-style: git@host:owner/repo -> host/owner/repo
			s = s[at+1:]
			s = strings.Replace(s, ":", "/", 1)
		}
	}

	return strings.ToLower(s)
}
