// Package transport centralizes HTTPS/SSH URL classification so every
// consumer (mapper, probe, clone-from, reclone, scan) shares one rule
// set. Replaces the historically-duplicated `classifyTransport` in
// gitmap/mapper and `looksLikeSCP` in gitmap/clonefrom (#15 of the
// post-v6.53.0 improvements list).
//
// The buckets match `constants.ScanTransport*` so existing consumers
// can switch over without renaming values.
package transport

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Classify returns one of the three stable transport buckets:
// constants.ScanTransportSSH, ScanTransportHTTPS, ScanTransportOther.
// Whitespace-only URLs are treated as "other" so downstream
// column-count invariants stay intact.
func Classify(url string) string {
	trimmed := strings.TrimSpace(url)
	if strings.HasPrefix(trimmed, "ssh://") {
		return constants.ScanTransportSSH
	}
	if strings.HasPrefix(trimmed, "https://") {
		return constants.ScanTransportHTTPS
	}
	if IsSCPStyle(trimmed) {
		return constants.ScanTransportSSH
	}
	return constants.ScanTransportOther
}

// IsSCPStyle reports whether `url` is the `[user@]host:path` form
// that git accepts as an SSH remote (e.g. `git@github.com:o/r.git`).
func IsSCPStyle(url string) bool {
	if strings.Contains(url, "://") {
		return false
	}
	colon := strings.Index(url, ":")
	if colon <= 0 {
		return false
	}
	host := url[:colon]
	return !strings.ContainsAny(host, "/\\")
}
