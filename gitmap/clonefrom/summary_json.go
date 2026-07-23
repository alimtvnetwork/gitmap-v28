package clonefrom

// JSON-report emit surface for `gitmap clone-from --execute`.
// Split from summary.go so each file stays under the 200-line cap
// (mem://style/code-constraints, item 3).
//
// Envelope shape (schemaVersion 3):
//
//	{
//	  "schemaVersion": 3,
//	  "transport": { "ssh": N, "https": N, "other": N },
//	  "provenance": [ {"field": "url", "stage": "scan"}, ... ],
//	  "rows": [ {row...}, ... ]
//	}
//
// Provenance is ENVELOPE-LEVEL (one entry per row field, not per
// row): every row in the same file shares the same field-origin
// contract, so duplicating the map per-row would bloat the file
// without adding signal. The authoritative mapping lives in
// constants.CloneFromReportProvenance.

import (
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// reportRowJSON is the on-disk JSON shape per result. Field names
// (and JSON tags) mirror the CSV column set 1:1 so consumers can
// flip between formats without a schema delta. Keep this type
// private; the only path to a JSON report is writeReportRowsJSON.
type reportRowJSON struct {
	URL             string  `json:"url"`
	Dest            string  `json:"dest"`
	Branch          string  `json:"branch"`
	Depth           int     `json:"depth"`
	Status          string  `json:"status"`
	Detail          string  `json:"detail"`
	DurationSeconds float64 `json:"duration_seconds"`
}

// transportTallyJSON mirrors the terminal `transport: N ssh, N https,
// N other` line so JSON consumers see the SAME counts without having
// to re-derive them from row URLs. Always emitted (zeros included)
// so the envelope shape is unconditional.
type transportTallyJSON struct {
	SSH   int `json:"ssh"`
	HTTPS int `json:"https"`
	Other int `json:"other"`
}

// provenanceEntryJSON is one field-origin record under the envelope's
// `provenance` array. Slice-of-objects (not a map) so the on-disk
// order is stable regardless of encoding/json's map-key sort
// behavior; readers can scan it in row-column order.
type provenanceEntryJSON struct {
	Field string `json:"field"`
	Stage string `json:"stage"`
}

type reportEnvelopeJSON struct {
	SchemaVersion int                   `json:"schemaVersion"`
	Transport     transportTallyJSON    `json:"transport"`
	Provenance    []provenanceEntryJSON `json:"provenance"`
	Rows          []reportRowJSON       `json:"rows"`
}

// writeReportRowsJSON emits the result set as a versioned JSON
// envelope. Rows is always serialized as `[]` (never `null`) for
// the empty case so jq pipelines can treat it unconditionally.
// Trailing newline matches POSIX text-file convention.
func writeReportRowsJSON(w io.Writer, results []Result) error {
	rows := buildReportRows(results)
	ssh, https, other := TransportTally(results)
	envelope := reportEnvelopeJSON{
		SchemaVersion: constants.CloneFromReportSchemaVersion,
		Transport:     transportTallyJSON{SSH: ssh, HTTPS: https, Other: other},
		Provenance:    buildProvenanceEntries(),
		Rows:          rows,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	return enc.Encode(envelope)
}

// buildReportRows fans Result entries into JSON-tagged row structs.
// Split out so writeReportRowsJSON stays inside the 15-line budget.
func buildReportRows(results []Result) []reportRowJSON {
	rows := make([]reportRowJSON, 0, len(results))
	for _, r := range results {
		rows = append(rows, reportRowJSON{
			URL: r.Row.URL, Dest: r.Dest, Branch: r.Row.Branch,
			Depth: r.Row.Depth, Status: r.Status, Detail: r.Detail,
			DurationSeconds: r.Duration.Seconds(),
		})
	}

	return rows
}

// buildProvenanceEntries lifts constants.CloneFromReportProvenance
// into the JSON-tagged shape. Order preserved verbatim so the
// on-disk file is byte-stable across runs.
func buildProvenanceEntries() []provenanceEntryJSON {
	out := make([]provenanceEntryJSON, 0, len(constants.CloneFromReportProvenance))
	for _, p := range constants.CloneFromReportProvenance {
		out = append(out, provenanceEntryJSON{Field: p.Field, Stage: p.Stage})
	}

	return out
}
