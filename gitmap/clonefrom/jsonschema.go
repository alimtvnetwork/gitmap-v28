package clonefrom

// JSON-Schema emit surface. Backs `gitmap clone-from --emit-schema=<kind>`.
//
// Two kinds are supported, both draft 2020-12:
//
//   - "report" — the on-disk envelope written by WriteReportJSON
//     (.gitmap/clone-from-report-<unixts>.json). Schema mirrors
//     reportEnvelopeJSON / reportRowJSON / transportTallyJSON in
//     summary.go; if those types change, EmitReportSchema MUST be
//     updated and TestEmitReportSchema_PinnedKeys will fail until
//     it is.
//
//   - "input" — the array of scan records accepted by
//     `gitmap clone <file>` (clone-now path). Property set is sourced
//     from clonenow.KnownScanFields() so a new accepted field auto-
//     appears in the schema; required = at least one of the URL
//     fields returned by clonenow.RequiredScanURLFields().
//
// The schemas are emitted as pretty-printed JSON (2-space indent,
// trailing newline) so the output is reviewable as-is and stable
// across runs (encoding/json sorts struct fields by declaration
// order, and our maps below are built from sorted slices).

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// EmitSchema dispatches on kind and returns the pretty-printed JSON
// Schema bytes. Unknown kinds return an error wrapping the user-
// facing message so the CLI layer can print and exit 2 verbatim.
func EmitSchema(kind string) ([]byte, error) {
	switch kind {
	case constants.EmitSchemaKindReport:
		return EmitReportSchema()
	case constants.EmitSchemaKindInput:
		return EmitInputSchema()
	default:
		return nil, fmt.Errorf(constants.MsgCloneFromEmitSchemaUnknown, kind)
	}
}

func reportStatusEnum() map[string]any {
	return enumSchema("Outcome bucket.", []string{
		constants.CloneFromStatusOK,
		constants.CloneFromStatusSkipped,
		constants.CloneFromStatusFailed,
	})
}

func reportRowProps() map[string]any {
	return orderedProps(
		kv("url", strSchema("Clone source URL (verbatim from input).")),
		kv("dest", strSchema("Resolved destination directory.")),
		kv("branch", strSchema("Pinned branch; empty = remote HEAD.")),
		kv("depth", intSchema("Shallow-clone depth; 0 = full history.")),
		kv("status", reportStatusEnum()),
		kv("detail", strSchema("Status-specific context (skip reason / git stderr).")),
		kv("duration_seconds", numSchema("Wall-clock seconds for this row.")),
	)
}

func reportTransportProps() map[string]any {
	return orderedProps(
		kv("ssh", intSchema("SSH-classified URL count.")),
		kv("https", intSchema("HTTPS-classified URL count.")),
		kv("other", intSchema("Catch-all URL count.")),
	)
}

func reportProvenanceItem() map[string]any {
	stages := []string{
		constants.ProvenanceStageScan,
		constants.ProvenanceStageMapper,
		constants.ProvenanceStageClonefrom,
	}
	props := orderedProps(
		kv("field", strSchema("Row-level field name (matches a key under rows[].).")),
		kv("stage", enumSchema("Pipeline stage that populates the field.", stages)),
	)

	return objectSchema(props, []string{"field", "stage"}, "One field-origin record.")
}

func reportRowsSchema() map[string]any {
	req := []string{"url", "dest", "branch", "depth", "status", "detail", "duration_seconds"}
	rowObj := objectSchema(reportRowProps(), req, "One result per planned clone.")
	desc := "Always an array (never null) so consumers can treat it unconditionally."

	return arraySchema(rowObj, desc)
}

func reportRootProps() map[string]any {
	provDesc := "Envelope-level field-origin map (one entry per row field, " +
		"shared by all rows). Order matches rows[] column order."

	return orderedProps(
		kv("schemaVersion", constIntSchema(constants.CloneFromReportSchemaVersion,
			"Pinned envelope version; bumped only on shape changes.")),
		kv("transport", objectSchema(reportTransportProps(),
			[]string{"ssh", "https", "other"},
			"Per-mode tally matching the terminal summary line.")),
		kv("provenance", arraySchema(reportProvenanceItem(), provDesc)),
		kv("rows", reportRowsSchema()),
	)
}

// EmitReportSchema returns the JSON Schema for the clone-from JSON
// report envelope. The schemaVersion `const` is wired to
// constants.CloneFromReportSchemaVersion so a bump there is auto-
// reflected here (and the writer emits a matching number).
func EmitReportSchema() ([]byte, error) {
	root := rootSchema(constants.CloneFromSchemaIDReport,
		"gitmap clone-from JSON report envelope.",
		reportRootProps(), []string{"schemaVersion", "transport", "provenance", "rows"})

	return marshalSchema(root)
}

func buildInputItemSchema() map[string]any {
	itemProps := make([]propKV, 0, len(clonenow.KnownScanFields()))
	for _, fieldName := range clonenow.KnownScanFields() {
		itemProps = append(itemProps, kv(fieldName, scanFieldSchema(fieldName)))
	}
	desc := "One scan record. At least one of httpsUrl / sshUrl must be present " +
		"(enforced by anyOf, not by required, so either field alone satisfies it)."
	item := objectSchema(orderedProps(itemProps...), nil, desc)
	item["additionalProperties"] = false
	item["anyOf"] = anyOfRequired(clonenow.RequiredScanURLFields())

	return item
}

func buildInputRootSchema() map[string]any {
	root := rootSchema(constants.CloneFromSchemaIDInput,
		"gitmap clone-now input: top-level array of scan records.",
		nil, nil)
	root["type"] = "array"
	root["items"] = buildInputItemSchema()
	delete(root, "properties")
	delete(root, "required")

	return root
}

// EmitInputSchema returns the JSON Schema for the clone-now input
// (top-level array of scan records). Field set is sourced from
// clonenow.KnownScanFields() so additions auto-propagate.
func EmitInputSchema() ([]byte, error) {
	return marshalSchema(buildInputRootSchema())
}

// marshalSchema renders the schema as 2-space-indented JSON with a
// trailing newline (POSIX text-file convention). EscapeHTML off so
// `&` in $id stays readable.
func marshalSchema(schemaObj any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(schemaObj); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
