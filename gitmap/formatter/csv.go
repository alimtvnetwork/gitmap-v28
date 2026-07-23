package formatter

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// WriteCSV writes records to the given writer in CSV format.
//
// Records are validated first; per-issue warnings are emitted to the
// configured sink (default os.Stderr) but the write always proceeds.
// See validate.go for the warn-and-write policy.
func WriteCSV(w io.Writer, records []model.ScanRecord) error {
	issueCount := emitValidationWarnings(records)

	cw := csv.NewWriter(w)
	// Force CRLF line endings on every row (header + data) so the
	// output matches RFC 4180 and stays byte-identical between
	// Linux/macOS and Windows runs. Pinned by csvcrlf_contract_test.go.
	cw.UseCRLF = true
	err := cw.Write(constants.ScanCSVHeaders)
	if err != nil {
		return err
	}
	err = writeCSVRows(cw, records)
	if err != nil {
		return err
	}
	emitWriteSummary("csv", len(records), issueCount)

	return nil
}

// writeCSVRows writes each record as a CSV row and flushes.
func writeCSVRows(cw *csv.Writer, records []model.ScanRecord) error {
	for _, r := range records {
		err := writeCSVRow(cw, r)
		if err != nil {
			return err
		}
	}
	cw.Flush()

	return cw.Error()
}

// writeCSVRow converts a single record to a CSV row. Depth is
// rendered as a base-10 integer so it sorts numerically when the
// CSV is opened in a spreadsheet. New columns (repoId,
// discoveredUrl, transport) are appended at the end so legacy
// positional readers still resolve columns 0..9 unchanged.
func writeCSVRow(cw *csv.Writer, r model.ScanRecord) error {
	row := []string{
		r.RepoName, r.HTTPSUrl, r.SSHUrl, r.Branch, r.BranchSource,
		r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes,
		strconv.Itoa(r.Depth),
		r.RepoID, r.DiscoveredURL, r.Transport,
	}

	return cw.Write(row)
}

// ParseCSV reads records from a CSV reader.
func ParseCSV(reader io.Reader) ([]model.ScanRecord, error) {
	cr := csv.NewReader(reader)
	cr.FieldsPerRecord = -1 // tolerate legacy 8/9/10/12-col CSVs alongside current 13-col layout
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}

	return parseCSVRows(rows), nil
}

// parseCSVRows converts raw CSV rows (skipping header) into records.
// Supports five layouts (auto-detected by column count) so older
// CSVs keep round-tripping after each additive schema bump:
//
//   - legacy 8 cols : pre-branchSource layout.
//   - 9 cols        : pre-depth layout (branchSource present).
//   - 10 cols       : pre-repoId layout (depth present).
//   - 12 cols       : pre-transport layout (repoId, discoveredUrl).
//   - 13 cols       : current layout (transport appended).
func parseCSVRows(rows [][]string) []model.ScanRecord {
	records := make([]model.ScanRecord, 0, len(rows))
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) >= 8 {
			records = append(records, rowToRecord(row))
		}
	}

	return records
}

// rowToRecord maps a CSV row to a ScanRecord. Dispatches on column
// count to keep each branch under the 15-line function budget.
func rowToRecord(row []string) model.ScanRecord {
	if len(row) >= 9 {

		return rowToRecordWithSource(row)
	}

	return rowToRecordLegacy(row)
}

// rowToRecordWithSource handles the 9-col (no depth), 10-col
// (depth, no repoId), 12-col (no transport), and 13-col (current)
// layouts. Missing fields fall back to zero values so partially-
// populated CSVs still load.
func rowToRecordWithSource(row []string) model.ScanRecord {
	depth := 0
	if len(row) >= 10 {
		parsed, err := strconv.Atoi(row[9])
		if err == nil {
			depth = parsed
		}
	}
	repoID, discovered, transport := "", "", ""
	if len(row) >= 12 {
		repoID = row[10]
		discovered = row[11]
	}
	if len(row) >= 13 {
		transport = row[12]
	}

	return model.ScanRecord{
		RepoName: row[0], HTTPSUrl: row[1], SSHUrl: row[2],
		Branch: row[3], BranchSource: row[4],
		RelativePath: row[5], AbsolutePath: row[6],
		CloneInstruction: row[7], Notes: row[8],
		Depth:         depth,
		RepoID:        repoID,
		DiscoveredURL: discovered,
		Transport:     transport,
	}
}

// rowToRecordLegacy handles the pre-branchSource 8-col layout.
// BranchSource, Depth, RepoID, DiscoveredURL are left at zero values.
func rowToRecordLegacy(row []string) model.ScanRecord {
	notes := ""
	if len(row) > 7 {
		notes = row[7]
	}

	return model.ScanRecord{
		RepoName: row[0], HTTPSUrl: row[1], SSHUrl: row[2],
		Branch: row[3], RelativePath: row[4], AbsolutePath: row[5],
		CloneInstruction: row[6], Notes: notes,
	}
}
