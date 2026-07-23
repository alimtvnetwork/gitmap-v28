package clonenow

// Parser entry point: ParseFile dispatches on the explicit format
// override (cfg.Format, when non-empty) or otherwise on file
// extension, returns a fully-validated Plan.
//
// Three formats are supported because `gitmap scan` itself produces
// all three (.json, .csv, .txt) and clone-now's contract is "feed me
// any scan output -- I will round-trip it back into a working tree."
//
// JSON and CSV piggyback on formatter.ParseJSON / formatter.ParseCSV
// so the schema (RepoName, HTTPSUrl, SSHUrl, Branch, RelativePath, ...)
// is read by exactly one parser per format -- if scan ever evolves
// the schema, clone-now picks up the change automatically.

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/formatter"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ParseFile is the package's only public parser entry point.
// `format` may be "" (auto-detect from extension) or one of
// constants.CloneNowFormat{JSON,CSV,Text}. `mode` and `onExists`
// must already be validated by the caller -- ParseFile records
// them on the Plan but does not interpret them (that's the
// executor / renderer's job).
func ParseFile(path, format, mode, onExists string) (Plan, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Plan{}, fmt.Errorf(constants.ErrCloneNowAbsPath, path, err)
	}
	resolved := format
	if len(resolved) == 0 {
		detected, derr := detectFormat(abs)
		if derr != nil {
			return Plan{}, derr
		}
		resolved = detected
	}
	rows, err := parseByFormat(abs, resolved)
	if err != nil {
		return Plan{}, err
	}
	rows = dedupRows(rows)
	if len(rows) == 0 {
		return Plan{}, fmt.Errorf(constants.ErrCloneNowEmpty, abs)
	}

	return Plan{
		Source: abs, Format: resolved,
		Mode: mode, OnExists: onExists,
		Rows: rows,
	}, nil
}

// detectFormat picks json / csv / text from the lowercased file
// extension. Unknown extensions return a clear error so users get
// an actionable message ("supported extensions are .json, .csv,
// .txt") instead of a confusing "0 rows" failure from a fallback
// parser. Callers can still bypass extension-based dispatch by
// passing an explicit --format.
func detectFormat(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return constants.CloneNowFormatJSON, nil
	case ".csv":
		return constants.CloneNowFormatCSV, nil
	case ".txt":
		return constants.CloneNowFormatText, nil
	}

	return "", fmt.Errorf(constants.ErrCloneNowUnsupportedExt, ext, path)
}

// parseByFormat opens the file once and dispatches the io.Reader to
// the format-specific parser. Centralized so all three branches
// share identical open / close / error-wrapping behavior.
func parseByFormat(path, format string) ([]Row, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrCloneNowOpen, path, err)
	}
	defer f.Close()

	switch format {
	case constants.CloneNowFormatJSON:
		return parseJSONWithSchema(f)
	case constants.CloneNowFormatCSV:
		return parseCSVWithSchema(f)
	}

	return parseTextRows(f)
}

// parseJSONWithSchema slurps the JSON file once, runs the schema
// guard against the bytes, then hands the same bytes to the
// tolerant formatter.ParseJSON so we don't re-decode twice with
// divergent semantics.
func parseJSONWithSchema(f io.Reader) ([]Row, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrCloneNowJSONDecode, err)
	}
	if err := validateJSONSchema(data); err != nil {
		return nil, err
	}
	recs, err := formatter.ParseJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf(constants.ErrCloneNowJSONDecode, err)
	}

	return rowsFromRecords(recs), nil
}

// parseCSVWithSchema slurps the CSV file once, validates the header
// row, then re-reads the full payload through formatter.ParseCSV.
// Header validation must run BEFORE the tolerant parser because
// otherwise a wrong header would silently mis-map columns.
func parseCSVWithSchema(f io.Reader) ([]Row, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrCloneNowCSVRead, err)
	}
	if err := validateCSVSchema(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	recs, err := parseCSVByHeaderName(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf(constants.ErrCloneNowCSVRead, err)
	}
	if len(recs) == 0 {
		// Fall back to the legacy positional parser for full 8+ col
		// scan exports whose header names are not normalized cleanly.
		recs, err = formatter.ParseCSV(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf(constants.ErrCloneNowCSVRead, err)
		}
	}

	return rowsFromRecords(recs), nil
}

// parseCSVByHeaderName reads a CSV using the header row as the
// authoritative column-name → ScanRecord-field mapping. Tolerates
// any subset of known fields (minimum: at least one URL column),
// in any order, and survives BOM / quoted / padded header cells via
// normalizeHeaderName. Returns nil records on a header that maps no
// fields so the caller can fall back to the positional parser.
func parseCSVByHeaderName(r io.Reader) ([]model.ScanRecord, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.LazyQuotes = true
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	idx := buildHeaderIndex(rows[0])
	if !headerHasURL(idx) {
		return nil, nil
	}
	out := make([]model.ScanRecord, 0, len(rows)-1)
	for _, row := range rows[1:] {
		out = append(out, recordFromIndexedRow(row, idx))
	}

	return out, nil
}

// buildHeaderIndex maps each known field name to its column index
// in the header row. Unknown / empty cells are ignored.
func buildHeaderIndex(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, col := range header {
		name := normalizeHeaderName(col)
		if len(name) > 0 && knownScanFields[name] {
			idx[name] = i
		}
	}

	return idx
}

// headerHasURL reports whether the indexed header includes at least
// one URL column, matching validateCSVHeader's contract.
func headerHasURL(idx map[string]int) bool {
	_, hasHTTPS := idx["httpsUrl"]
	_, hasSSH := idx["sshUrl"]

	return hasHTTPS || hasSSH
}

// recordFromIndexedRow assembles a ScanRecord by pulling each known
// field from its mapped column. Missing columns yield zero values.
func recordFromIndexedRow(row []string, idx map[string]int) model.ScanRecord {
	get := func(name string) string {
		i, ok := idx[name]
		if !ok || i >= len(row) {
			return ""
		}
		return row[i]
	}
	depth := 0
	if d, err := strconv.Atoi(get("depth")); err == nil {
		depth = d
	}
	return model.ScanRecord{
		RepoName: get("repoName"), HTTPSUrl: get("httpsUrl"), SSHUrl: get("sshUrl"),
		Branch: get("branch"), BranchSource: get("branchSource"),
		RelativePath: get("relativePath"), AbsolutePath: get("absolutePath"),
		CloneInstruction: get("cloneInstruction"), Notes: get("notes"),
		Depth:         depth,
		RepoID:        get("repoId"),
		DiscoveredURL: get("discoveredUrl"),
		Transport:     get("transport"),
	}
}

// rowsFromRecords lifts a slice of scan records into clone-now Rows,
// dropping any record that has neither URL set (truly nothing to
// clone) and back-filling RelativePath from the URL basename when
// the record didn't carry one (older scan exports may omit it).
func rowsFromRecords(recs []model.ScanRecord) []Row {
	out := make([]Row, 0, len(recs))
	for _, rec := range recs {
		if len(rec.HTTPSUrl) == 0 && len(rec.SSHUrl) == 0 {
			continue
		}
		dest := rec.RelativePath
		if len(dest) == 0 {
			dest = deriveDestFromRecord(rec)
		}
		out = append(out, Row{
			RepoName:     rec.RepoName,
			HTTPSUrl:     strings.TrimSpace(rec.HTTPSUrl),
			SSHUrl:       strings.TrimSpace(rec.SSHUrl),
			Branch:       strings.TrimSpace(rec.Branch),
			RelativePath: strings.TrimSpace(dest),
		})
	}

	return out
}

// deriveDestFromRecord picks the best non-empty URL on the record
// and returns its basename (sans .git). Used as the fallback when a
// scan record came in without a recorded RelativePath.
func deriveDestFromRecord(rec model.ScanRecord) string {
	url := rec.HTTPSUrl
	if len(url) == 0 {
		url = rec.SSHUrl
	}

	return DeriveDest(url)
}

// dedupRows collapses rows that share the same destination path.
// Later rows win so users can override an earlier row by re-listing
// the same dest -- mirrors clonefrom's "later wins" semantics so the
// two commands feel uniform when a user switches between them.
func dedupRows(rows []Row) []Row {
	seen := make(map[string]int, len(rows))
	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		key := r.RelativePath
		if idx, ok := seen[key]; ok {
			out[idx] = r

			continue
		}
		seen[key] = len(out)
		out = append(out, r)
	}

	return out
}
