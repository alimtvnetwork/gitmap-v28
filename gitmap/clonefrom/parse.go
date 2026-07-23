package clonefrom

// Parser entry point: ParseFile dispatches on file extension and
// returns a fully-validated Plan. JSON and CSV are the only formats
// supported — picking format from extension (rather than from a
// `--format` flag) keeps the CLI ergonomic for shell scripts that
// typically know whether they wrote `.json` or `.csv`.
//
// Validation happens in two passes:
//
//   1. Per-row syntax (URL non-empty, depth non-negative).
//   2. Cross-row dedup (same URL+dest combination collapses to one
//      Row; later rows win for branch/depth so users can override
//      a default by re-listing the URL).
//
// Errors point at line/row numbers (1-indexed for CSV including the
// header) so users can grep their input file directly.

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ParseFile is the package's only public parser entry point.
// Returns a fully-validated Plan or a wrapped error explaining
// where parsing failed (file open, format mismatch, row N invalid).
func ParseFile(path string) (Plan, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Plan{}, fmt.Errorf(constants.ErrCloneFromAbsPath, path, err)
	}
	f, err := os.Open(abs)
	if err != nil {
		return Plan{}, fmt.Errorf(constants.ErrCloneFromOpen, abs, err)
	}
	defer f.Close()

	format := detectFormat(abs)
	rows, err := parseByFormat(f, format)
	if err != nil {
		return Plan{}, err
	}
	rows = dedupRows(rows)

	return Plan{Source: abs, Format: format, Rows: rows}, nil
}

// detectFormat picks json vs csv from the lowercased file
// extension. Anything else falls back to csv (the more permissive
// format) — a malformed csv produces a row-level error users can
// fix; a "format=unknown" hard error would just confuse them.
func detectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		return "json"
	}

	return "csv"
}

// parseByFormat is the dispatch helper. Kept tiny so future
// formats (yaml? toml?) are a one-line addition.
func parseByFormat(r io.Reader, format string) ([]Row, error) {
	if format == "json" {
		return parseJSON(r)
	}

	return parseCSV(r)
}

// parseJSON expects a top-level array of objects. Each object's
// fields map 1:1 to Row fields with lowercase keys: url, dest,
// branch, depth. Unknown fields are tolerated (forward-compat:
// future schema additions don't break old gitmap binaries).
func parseJSON(r io.Reader) ([]Row, error) {
	var raw []map[string]any
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf(constants.ErrCloneFromJSONDecode, err)
	}
	out := make([]Row, 0, len(raw))
	for i, obj := range raw {
		row, err := jsonRow(obj)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrCloneFromJSONRow, i+1, err)
		}
		out = append(out, row)
	}

	return out, nil
}

// jsonRow extracts one Row from a parsed JSON object. Centralized
// so parseJSON stays under the per-function budget.
func jsonRow(obj map[string]any) (Row, error) {
	url, _ := obj["url"].(string)
	dest, _ := obj["dest"].(string)
	branch, _ := obj["branch"].(string)
	checkout, _ := obj["checkout"].(string)
	depth := 0
	if d, ok := obj["depth"].(float64); ok {
		depth = int(d)
	}
	row := Row{URL: strings.TrimSpace(url), Dest: strings.TrimSpace(dest),
		Branch: strings.TrimSpace(branch), Depth: depth,
		Checkout: strings.ToLower(strings.TrimSpace(checkout))}

	return row, validateRow(row)
}

// CSV parsing (parseCSV, csvIndex, indexCSVHeader, csvRow, get)
// lives in parsecsv.go to keep this file under the per-file budget.
