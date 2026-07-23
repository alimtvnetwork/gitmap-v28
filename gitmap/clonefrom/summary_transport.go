package clonefrom

// summary_transport.go — single source of truth for the SSH-vs-HTTPS
// line printed by both RenderSummary (legacy) and
// RenderSummaryTerminal (enriched). Built on top of ClassifyScheme so
// the transport counts can never drift from the per-scheme tally
// shown in the terminal block. No new fields are read from Result —
// the URL on r.Row.URL is the only input, matching what the rest of
// the summary surface already consumes.

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TransportTally collapses ClassifyScheme's seven buckets into the
// three columns the user-facing summary reports: ssh, https, other.
// Returned in that order so call sites can format with a single
// printf without referencing field names.
func TransportTally(results []Result) (sshCount, httpsCount, otherCount int) {
	for _, r := range results {
		switch ClassifyScheme(r.Row.URL) {
		case constants.CloneFromSchemeSSH, constants.CloneFromSchemeSCP:
			sshCount++
		case constants.CloneFromSchemeHTTPS:
			httpsCount++
		default:
			otherCount++
		}
	}

	return sshCount, httpsCount, otherCount
}

// writeTransportLine emits the single shared transport line. Kept
// tiny so both renderers can call it without re-deriving the format.
func writeTransportLine(w io.Writer, results []Result) error {
	ssh, https, other := TransportTally(results)
	_, err := fmt.Fprintf(w, constants.CloneFromSummaryTransportFmt,
		ssh, https, other)

	return err
}
