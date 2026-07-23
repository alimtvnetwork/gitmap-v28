package mapper

// Transport classification for ScanRecord.Transport.
//
// Collapses the discovered remote URL into one of three stable
// buckets surfaced in the CSV / JSON output:
//
//   - "ssh"   — ssh://...  OR  scp-style  user@host:path
//   - "https" — https://...
//   - "other" — http://, git://, file://, empty string, anything
//               that doesn't match the two transports we filter on
//
// Rationale for the three-bucket collapse (vs the seven-bucket
// ClassifyScheme that clonefrom uses for its terminal summary): the
// CSV column is meant to make `awk -F, '$13=="ssh"'` filtering
// trivial, so we stick to the two transports the user actually
// chooses between in practice and lump the rest under "other".
// Parity with clonefrom.TransportTally is pinned by
// TestClassifyTransport_MatchesClonefromTally so the two views
// can never silently diverge.

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// classifyTransport buckets `url` into the three labels documented
// above. Empty / whitespace-only URLs land in "other" so a repo
// with no remote configured still gets a stable, filterable value
// (rather than an empty cell that breaks downstream column-count
// invariants).
func classifyTransport(url string) string {
	trimmed := strings.TrimSpace(url)
	if strings.HasPrefix(trimmed, "ssh://") {
		return constants.ScanTransportSSH
	}
	if strings.HasPrefix(trimmed, "https://") {
		return constants.ScanTransportHTTPS
	}
	if isSCPStyle(trimmed) {
		return constants.ScanTransportSSH
	}

	return constants.ScanTransportOther
}

// isSCPStyle reports whether `url` is the `[user@]host:path` form
// that git accepts as an SSH remote (e.g. `git@github.com:o/r.git`).
// Mirrors clonefrom.looksLikeSCP's rule: must contain a colon, must
// NOT contain `://` (that would be a scheme-shaped URL handled by
// the caller), and the segment before the colon must look like a
// host (no slashes). Kept tiny so classifyTransport stays one
// readable function.
func isSCPStyle(url string) bool {
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
