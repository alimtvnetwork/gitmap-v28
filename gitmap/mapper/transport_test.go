package mapper

// Tests for classifyTransport. Pin the three documented buckets
// AND verify lockstep parity with clonefrom.TransportTally so the
// CSV column and the clone-from terminal summary can never silently
// disagree on what counts as ssh / https / other.

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestClassifyTransport_Buckets pins one URL per documented case
// across all three buckets so a future refactor can't silently drop
// or reroute a transport. SSH covers both ssh:// and scp-style.
func TestClassifyTransport_Buckets(t *testing.T) {
	cases := []struct {
		url, want string
	}{
		{"https://example.com/o/r.git", constants.ScanTransportHTTPS},
		{"ssh://git@example.com/o/r.git", constants.ScanTransportSSH},
		{"git@example.com:o/r.git", constants.ScanTransportSSH},
		{"http://example.com/o/r.git", constants.ScanTransportOther},
		{"git://example.com/o/r.git", constants.ScanTransportOther},
		{"file:///srv/repos/r.git", constants.ScanTransportOther},
		{"", constants.ScanTransportOther},
		{"   ", constants.ScanTransportOther},
		{"   https://example.com/o/r.git  ", constants.ScanTransportHTTPS},
	}
	for _, tc := range cases {
		if got := classifyTransport(tc.url); got != tc.want {
			t.Errorf("classifyTransport(%q) = %q, want %q",
				tc.url, got, tc.want)
		}
	}
}

// TestClassifyTransport_MatchesClonefromTally is the parity guard:
// for a representative URL set, the CSV-side three-bucket collapse
// must match the bucket clonefrom.TransportTally would assign. If
// these two ever drift the user filters their CSV one way and reads
// the terminal summary another.
func TestClassifyTransport_MatchesClonefromTally(t *testing.T) {
	urls := []string{
		"https://example.com/o/r.git",
		"http://example.com/o/r.git",
		"ssh://git@example.com/o/r.git",
		"git@example.com:o/r.git",
		"git://example.com/o/r.git",
		"file:///srv/repos/r.git",
		"",
		"not-a-url",
	}
	for _, u := range urls {
		mine := classifyTransport(u)
		theirs := tallyBucketFor(u)
		if mine != theirs {
			t.Errorf("transport drift on %q: mapper=%q clonefrom=%q",
				u, mine, theirs)
		}
	}
}

// tallyBucketFor reproduces the three-bucket collapse that
// clonefrom.TransportTally applies, using ONE result per call so we
// can compare URL-by-URL. Kept here (rather than reused from
// clonefrom) so a regression in either place is attributable.
func tallyBucketFor(url string) string {
	results := []clonefrom.Result{{Row: clonefrom.Row{URL: url}}}
	ssh, https, other := clonefrom.TransportTally(results)
	switch {
	case ssh == 1:
		return constants.ScanTransportSSH
	case https == 1:
		return constants.ScanTransportHTTPS
	case other == 1:
		return constants.ScanTransportOther
	}

	return "<unreachable: tally returned no bucket>"
}
