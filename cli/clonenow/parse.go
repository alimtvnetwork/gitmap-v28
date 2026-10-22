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

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/formatter"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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
	var derr error
	if len(resolved) == 0 {
		resolved, derr = detectFormat(abs)
	}

	if derr != nil {
		return Plan{}, derr
	}

	rowsWrap := parseByFormat(abs, resolved)
	if rowsWrap.HasError() {
		return Plan{}, rowsWrap.Error()
	}
	rows := rowsWrap.Items()

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
func parseByFormat(path, format string) RowSliceWrap {
	f, err := os.Open(path)
	if err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "parseByFormat open"))
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

func parseJSONWithSchema(f io.Reader) RowSliceWrap {
	data, err := io.ReadAll(f)
	if err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "parseJSONWithSchema read"))
	}

	payload, _, extractErr := jsonx.ExtractPayload(data)
	if extractErr != nil {
		payload = data
	}

	if err := validateJSONSchema(payload); err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "parseJSONWithSchema validate"))
	}

	recs, err := formatter.ParseJSON(bytes.NewReader(payload))
	if err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "parseJSONWithSchema parse"))
	}

	return result.OkSlice(rowsFromRecords(recs))
}

func parseCSVWithSchema(f io.Reader) RowSliceWrap {
	data, err := io.ReadAll(f)
	if err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "parseCSVWithSchema read"))
	}

	if err := validateCSVSchema(bytes.NewReader(data)); err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "parseCSVWithSchema validate"))
	}

	recsWrap := parseCSVByHeaderName(bytes.NewReader(data))
	fallbackRes := fallbackParseCSV(data, recsWrap.Items())
	if fallbackRes.HasError() {
		return result.FailSlice[Row](fallbackRes.Err)
	}

	return result.OkSlice(rowsFromRecords(fallbackRes.Items()))
}

func fallbackParseCSV(data []byte, recs []model.ScanRecord) ScanRecordSliceWrap {
	if len(recs) > 0 {
		return result.OkSlice(recs)
	}

	parsed, parseErr := formatter.ParseCSV(bytes.NewReader(data))
	if parseErr != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(parseErr, "parseCSVWithSchema parse"))
	}

	return result.OkSlice(parsed)
}

func parseCSVByHeaderName(r io.Reader) ScanRecordSliceWrap {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.LazyQuotes = true
	rows, err := cr.ReadAll()
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "parseCSVByHeaderName read"))
	}

	if len(rows) < 2 {
		return result.OkSlice[model.ScanRecord](nil)
	}

	idx := buildHeaderIndex(rows[0])
	if !headerHasURL(idx) {
		return result.OkSlice[model.ScanRecord](nil)
	}

	out := make([]model.ScanRecord, 0, len(rows)-1)
	for _, row := range rows[1:] {
		out = append(out, recordFromIndexedRow(row, idx))
	}

	return result.OkSlice(out)
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
		i, isFound := idx[name]
		if !isFound || i >= len(row) {
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
		if idx, isSeen := seen[key]; isSeen {
			out[idx] = r

			continue
		}

		seen[key] = len(out)
		out = append(out, r)
	}

	return out
}
