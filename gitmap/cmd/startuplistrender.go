package cmd

// Renderers for `gitmap startup-list --format=...`. Four formats:
// table (default human-readable), json (pretty array, indent
// configurable via --json-indent), jsonl (one minified object per
// line), csv (RFC4180 with header). All encoders take an io.Writer
// so contract tests can capture bytes into a buffer; the CLI
// dispatcher passes os.Stdout.
//
// JSON / JSONL encoding goes through gitmap/stablejson rather than
// encoding/json directly: stablejson builds each object key-by-key
// in caller-declared order and CANNOT be reordered by a future Go
// release. Pretty 2-space output is byte-identical to the legacy
// Encoder.SetIndent("", "  "), so existing golden fixtures pass
// unchanged. See gitmap/stablejson/stablejson.go for full rationale.

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// renderStartupList dispatches to the per-format encoder. The
// `jsonIndent` parameter is only consulted for `--format=json`;
// other formats ignore it (jsonl is line-oriented and minified by
// design, csv has no nesting, table is human-prose). This keeps
// shell scripts that always pass `--json-indent=N` regardless of
// format from breaking when they switch to a non-json sink.
func renderStartupList(format string, jsonIndent int, dir string, entries []startup.Entry) error {
	switch format {
	case constants.OutputJSON:

		return encodeStartupListJSONIndent(os.Stdout, entries, jsonIndent)
	case constants.StartupListFormatJSONL:

		return encodeStartupListJSONL(os.Stdout, entries)
	case constants.OutputCSV:

		return encodeStartupListCSV(os.Stdout, entries)
	default:

		return renderStartupListTable(os.Stdout, dir, entries)
	}
}

// renderStartupListTable is the legacy human-readable rendering —
// extracted verbatim from the original runStartupList so behavior
// for users who don't pass --format is byte-identical.
//
// Takes an io.Writer so the table contract test in
// startuplisttable_contract_test.go can capture the rendered bytes
// into a buffer instead of redirecting os.Stdout. Returns an error
// (always nil today) so the dispatcher signature stays uniform with
// the JSON/CSV encoders and a future writer error can be surfaced
// without changing the dispatch site.
func renderStartupListTable(w io.Writer, dir string, entries []startup.Entry) error {
	fmt.Fprintf(w, constants.MsgStartupListHeader, dir)
	if len(entries) == 0 {
		fmt.Fprint(w, constants.MsgStartupListEmpty)

		return nil
	}
	for _, e := range entries {
		fmt.Fprintf(w, constants.MsgStartupListRow, e.Name, renderExec(e.Exec))
	}
	fmt.Fprintf(w, constants.MsgStartupListFooter, len(entries))

	return nil
}

// startupListJSONFields names — single source of truth for both the
// on-the-wire field labels and their order. Centralized as constants
// so any consumer-facing rename or reorder is one diff to review.
//
// CONTRACT: the names AND order here are pinned by
// gitmap/cmd/startuplistjson_contract_test.go (golden bytes for
// empty + single, structural key-order check for multi-entry).
// Reordering, renaming, or adding fields will fail those tests —
// intentional changes require regenerating the golden fixtures and
// bumping the consumer-facing changelog.
const (
	startupListJSONKeyName = "name"
	startupListJSONKeyPath = "path"
	startupListJSONKeyExec = "exec"
)

// encodeStartupListJSON writes a JSON array to w using stablejson
// at the long-standing 2-space-indent default. Thin wrapper around
// encodeStartupListJSONIndent kept so existing contract tests (and
// any future caller that doesn't care about indent) don't have to
// thread a width parameter through every call site.
//
// Empty input still encodes as `[]\n` (NOT `null`) so jq pipelines
// that do `length` work without conditionals.
func encodeStartupListJSON(w io.Writer, entries []startup.Entry) error {
	return encodeStartupListJSONIndent(w, entries, constants.StartupListJSONIndentDefault)
}

// encodeStartupListJSONIndent writes a JSON array with caller-
// controlled per-level indent width. `jsonIndent==0` emits a
// single-line minified `[{"k":v}]\n` (matches `jq -c` framing);
// any positive N emits a pretty-printed document with N spaces per
// level. Key order is identical across every indent value — the
// contract test in startuplistjson_indent_contract_test.go pins
// this by re-parsing each variant and comparing the key sequence.
//
// stablejson handles the `[]\n` empty case identically across
// indent values, so jq pipelines that do `length` keep working
// regardless of which indent the user chose.
func encodeStartupListJSONIndent(w io.Writer, entries []startup.Entry, jsonIndent int) error {
	indent := indentSpaces(jsonIndent)

	return stablejson.WriteArrayIndent(w, buildStartupListJSONItems(entries), indent)
}

// indentSpaces converts the integer --json-indent value into the
// per-level prefix string stablejson expects. 0 → empty (minified
// branch); N>0 → N spaces. Centralized so a future "tabs" extension
// (e.g. --json-indent=tab) lands in exactly one place.
func indentSpaces(n int) string {
	if n <= 0 {

		return ""
	}
	out := make([]byte, n)
	for i := range out {
		out[i] = ' '
	}

	return string(out)
}

// encodeStartupListJSONL writes one compact JSON object per line in
// the same key order as encodeStartupListJSON. Empty input writes
// zero bytes (NOT a stray `\n`) so `wc -l` of the stream equals the
// entry count exactly. Going through buildStartupListJSONItems keeps
// the JSON and JSONL field-order contracts byte-locked together —
// any future add/remove/rename of a column lands in both formats in
// one diff and the contract test below catches drift.
func encodeStartupListJSONL(w io.Writer, entries []startup.Entry) error {
	return stablejson.WriteJSONLines(w, buildStartupListJSONItems(entries))
}

// buildStartupListJSONItems is the single source of (field name,
// field order, value) for both --format=json and --format=jsonl.
// Centralized so a column rename or reorder is one diff, not two.
func buildStartupListJSONItems(entries []startup.Entry) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(entries))
	for _, e := range entries {
		items = append(items, []stablejson.Field{
			{Key: startupListJSONKeyName, Value: e.Name},
			{Key: startupListJSONKeyPath, Value: e.Path},
			{Key: startupListJSONKeyExec, Value: e.Exec},
		})
	}

	return items
}

// encodeStartupListCSV writes a header row followed by one row per
// entry. encoding/csv handles quoting of values containing commas,
// quotes, or newlines — important because Exec lines can include
// shell quoting and embedded spaces.
func encodeStartupListCSV(w io.Writer, entries []startup.Entry) error {
	cw := csv.NewWriter(w)
	// CRLF for cross-platform byte-identical output (RFC 4180).
	// Pinned by gitmap/cmd/csvcrlf_contract_test.go.
	cw.UseCRLF = true
	if err := cw.Write([]string{"name", "path", "exec"}); err != nil {

		return err
	}
	for _, e := range entries {
		if err := cw.Write([]string{e.Name, e.Path, e.Exec}); err != nil {

			return err
		}
	}
	cw.Flush()

	return cw.Error()
}

// renderExec keeps long Exec lines from making the list table noisy
// — falls back to a placeholder when the .desktop file omits Exec.
func renderExec(exec string) string {
	if len(exec) == 0 {

		return "(no Exec line)"
	}

	return exec
}
