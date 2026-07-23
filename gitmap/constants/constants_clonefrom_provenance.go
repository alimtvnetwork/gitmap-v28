package constants

// Provenance map for the JSON report envelope. Documents which
// pipeline stage populates each row-level field, so downstream
// consumers (audit dashboards, CI gates, debugging tools) can
// answer "where did this value come from?" without reading source.
//
// Stages — stable strings, treat as enum:
//
//	ProvenanceStageScan      — value flowed in from the user-provided
//	                           input file (JSON/CSV row); clone-from
//	                           never modifies it.
//	ProvenanceStageMapper    — value was derived inside clonefrom's
//	                           pre-execute resolver (e.g. Dest after
//	                           DeriveDest fallback when Row.Dest is
//	                           empty).
//	ProvenanceStageClonefrom — value was set during execute (Status,
//	                           Detail, DurationSeconds — all written
//	                           by executeRow / runGitClone).
//
// The map itself is a SLICE OF PAIRS (not a Go map) so the emit
// order is stable across runs without relying on encoding/json's
// alphabetical key sort. Field names match reportRowJSON's JSON
// tags one-for-one; a drift is caught by
// TestProvenance_CoversEveryReportField.
const (
	ProvenanceStageScan      = "scan"
	ProvenanceStageMapper    = "mapper"
	ProvenanceStageClonefrom = "clonefrom"
)

// CloneFromReportProvenanceField is one entry in the ordered
// provenance list. Public so tests and the JSON-Schema emitter can
// iterate it without re-deriving the mapping.
type CloneFromReportProvenanceField struct {
	Field string
	Stage string
}

// CloneFromReportProvenance is the canonical field → stage mapping
// embedded under `provenance` in every JSON report. Order matches
// reportRowJSON's struct field order so a reader scanning the
// envelope sees provenance entries in the same order as the row
// columns. Adding a new row field MUST add a corresponding entry
// here; the test guard fails fast on omission.
var CloneFromReportProvenance = []CloneFromReportProvenanceField{
	{Field: "url", Stage: ProvenanceStageScan},
	{Field: "dest", Stage: ProvenanceStageMapper},
	{Field: "branch", Stage: ProvenanceStageScan},
	{Field: "depth", Stage: ProvenanceStageScan},
	{Field: "status", Stage: ProvenanceStageClonefrom},
	{Field: "detail", Stage: ProvenanceStageClonefrom},
	{Field: "duration_seconds", Stage: ProvenanceStageClonefrom},
}
