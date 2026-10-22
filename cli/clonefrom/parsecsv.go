package clonefrom

// CSV parser half of ParseFile. Split from parse.go so each file
// stays under the 200-line per-file budget. Public surface is
// unchanged: callers go through ParseFile in parse.go which
// dispatches to parseCSV here based on file extension.

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// parseCSV expects a header row of `url,dest,branch,depth` (case-
// insensitive, only `url` required). Missing optional columns
// default to empty/zero. Extra columns past `depth` are ignored.
func parseCSV(r io.Reader) RowSliceWrap {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1 // tolerate ragged rows
	header, err := cr.Read()
	if err != nil {
		return result.FailSlice[Row](apperror.WrapSimple(err, "clone-from: read CSV header"))
	}

	idx := indexCSVHeader(header)
	if idx.url < 0 {
		return result.FailSlice[Row](apperror.WrapSimple(errors.New(constants.ErrCloneFromCSVNoURL), "parseCSV"))
	}

	return readCSVRows(cr, idx)
}

// readCSVRows is the inner loop split out so parseCSV stays under
// the function-length budget. Per-row failures are wrapped with the
// 1-indexed row number AND, when the failure is attributable to a
// single field, the offending column name — so an operator editing
// a 5,000-row spreadsheet can jump straight to the bad cell.
func readCSVRows(cr *csv.Reader, idx csvIndex) RowSliceWrap {
	var out []Row
	rowNum := 1 // header was row 1; first data row is 2
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		rowNum++
		if err != nil {
			return result.FailSlice[Row](apperror.WrapSimple(err, fmt.Sprintf("clone-from: CSV row %d", rowNum)))
		}

		row, col, err := csvRow(rec, idx)
		if err != nil {
			return result.FailSlice[Row](apperror.WrapSimple(wrapCSVRowErr(rowNum, col, err), "readCSVRows"))
		}

		out = append(out, row)
	}

	return result.OkSlice(out)
}

// wrapCSVRowErr picks the column-aware format when a column is
// known, falling back to the row-only format otherwise. Centralized
// so every caller produces identical wording.
func wrapCSVRowErr(rowNum int, col string, err error) error {
	if len(col) == 0 {
		return apperror.WrapSimple(err, fmt.Sprintf("clone-from: CSV row %d", rowNum))
	}

	return apperror.WrapSimple(err, fmt.Sprintf("clone-from: CSV row %d, column %q", rowNum, col))
}

// csvIndex maps logical column names to record indices. Negative
// value means "column absent" — only `url` is required.
type csvIndex struct{ url, dest, branch, depth, checkout int }

// indexCSVHeader walks the header row once and records each
// column's position. Headers are normalized through
// constants.CanonicalCSVColumn so common variations like "URL",
// "httpsURL", "relpath", and "https_url" all map to their canonical
// column. Unknown headers are ignored — extra spreadsheet columns
// must not break parsing.
func indexCSVHeader(header []string) csvIndex {
	idx := csvIndex{url: -1, dest: -1, branch: -1, depth: -1, checkout: -1}
	for i, name := range header {
		switch constants.CanonicalCSVColumn(name) {
		case constants.CSVColumnURL:
			idx.url = i
		case constants.CSVColumnDest:
			idx.dest = i
		case constants.CSVColumnBranch:
			idx.branch = i
		case constants.CSVColumnDepth:
			idx.depth = i
		case constants.CSVColumnCheckout:
			idx.checkout = i
		}
	}

	return idx
}

// csvRow extracts one Row from a parsed CSV record using the
// pre-computed column index. Returns the offending column name
// alongside the error so wrapCSVRowErr can name the bad cell.
// Returns "" for col when the failure is row-wide (e.g. dedup).
func csvRow(rec []string, idx csvIndex) (Row, string, error) {
	row := Row{
		URL:      strings.TrimSpace(get(rec, idx.url)),
		Dest:     strings.TrimSpace(get(rec, idx.dest)),
		Branch:   strings.TrimSpace(get(rec, idx.branch)),
		Checkout: strings.ToLower(strings.TrimSpace(get(rec, idx.checkout))),
	}

	depthStr := strings.TrimSpace(get(rec, idx.depth))
	row, col, err := applyDepth(row, depthStr)
	if err != nil {
		return row, col, err
	}

	if col, err := validateRowWithColumn(row); err != nil {
		return row, col, err
	}

	return row, "", nil
}

func applyDepth(row Row, depthStr string) (Row, string, error) {
	if len(depthStr) == 0 {
		return row, "", nil
	}

	d, err := strconv.Atoi(depthStr)
	if err != nil {
		return row, constants.CSVColumnDepth, apperror.WrapValidation(err, fmt.Sprintf(constants.ErrCloneFromBadDepth, depthStr))
	}

	row.Depth = d

	return row, "", nil
}

// get is a bounds-safe slice accessor. Returns empty string when
// the index is negative (column absent in header) or past the
// record's end (ragged row).
func get(rec []string, i int) string {
	if i < 0 || i >= len(rec) {
		return ""
	}

	return rec[i]
}
