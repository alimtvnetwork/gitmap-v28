package clonenext

// Batch input handling for `gitmap cn` operating on multiple repos.
//
// Two entry points feed the same dispatcher:
//
//   - LoadBatchFromCSV: read a curated list of repo paths from a CSV file
//     (header optional; column lookup by name when a header is present).
//   - WalkBatchFromDir: scan one level under a directory and return every
//     subdirectory that is itself a git repo.
//
// Both return absolute paths in deterministic (lexicographic) order so that
// re-runs and report rows stay stable.
//
// CSV-parsing contract (hardened in v3.43.2):
//
//   - Line endings: LF, CRLF, and bare-CR are all accepted. Bare-CR is
//     pre-normalized to LF so Go's encoding/csv (which only splits on LF)
//     handles classic-Mac-style files saved by old tools.
//   - BOM: a leading UTF-8 BOM (Excel-on-Windows default) is stripped
//     before parsing so the first cell of row 0 still matches header
//     keywords like "repo".
//   - Columns: when a header is present, the path column is located by
//     name ("repo" | "path" | "repo_path", case-insensitive). Optional
//     columns like "note", "version", or "tag" may appear in any order
//     and are ignored. Headerless input continues to use column 0.
//   - Ragged rows: rows with fewer cells than the path column or with an
//     empty path cell are skipped silently — they are not data.

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ErrBatchEmpty is returned when neither the CSV nor the cwd walk yielded
// any candidate repos. Callers surface this as a soft warning, not a crash.
var ErrBatchEmpty = errors.New("clonenext: batch input contained no repos")

// utf8BOM is the three-byte UTF-8 byte-order mark.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// pathHeaderAliases is the set of header cell values that mark the path
// column. Lookup is case-insensitive; trimmed before comparison.
var pathHeaderAliases = map[string]struct{}{
	"repo":      {},
	"path":      {},
	"repo_path": {},
	"repopath":  {},
}

// LoadBatchFromCSV reads a CSV file and returns one absolute repo path per
// non-empty data row. See the package doc-comment for the full parsing
// contract (BOM, line endings, header detection, optional columns).
func LoadBatchFromCSV(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rows, err := readAllCSVRows(bytes.NewReader(normalizeCSVBytes(raw)))
	if err != nil {
		return nil, err
	}

	paths := extractPathColumn(rows)
	if len(paths) == 0 {
		return nil, ErrBatchEmpty
	}

	return absoluteAndSorted(paths), nil
}

// normalizeCSVBytes strips a leading UTF-8 BOM and converts bare-CR line
// endings to LF. CRLF is left intact so Excel/Notepad output keeps its
// line shape; encoding/csv handles CRLF natively.
func normalizeCSVBytes(in []byte) []byte {
	in = bytes.TrimPrefix(in, utf8BOM)

	return convertBareCRToLF(in)
}

// convertBareCRToLF replaces every '\r' that is NOT followed by '\n' with
// '\n'. Used for classic-Mac line endings that encoding/csv won't split.
func convertBareCRToLF(in []byte) []byte {
	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] != '\r' {
			out = append(out, in[i])

			continue
		}
		if i+1 < len(in) && in[i+1] == '\n' {
			// CRLF — keep both bytes; csv.Reader strips the CR.
			out = append(out, '\r', '\n')
			i++

			continue
		}
		out = append(out, '\n')
	}

	return out
}

// readAllCSVRows reads every row from r using the standard CSV parser with
// FieldsPerRecord disabled so ragged rows are tolerated.
func readAllCSVRows(r io.Reader) ([][]string, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	return cr.ReadAll()
}

// extractPathColumn returns the path cell from each data row, handling
// both header-present and header-absent inputs.
func extractPathColumn(rows [][]string) []string {
	if len(rows) == 0 {
		return nil
	}

	colIdx, startIdx := resolvePathColumn(rows)

	out := make([]string, 0, len(rows)-startIdx)
	for _, row := range rows[startIdx:] {
		cell := safeCell(row, colIdx)
		if len(cell) > 0 {
			out = append(out, cell)
		}
	}

	return out
}

// resolvePathColumn inspects row 0 and returns (column-index, first-data-row).
// If row 0 looks like a header, it is consumed and the path column is
// located by name. Otherwise column 0 is used and row 0 is data.
func resolvePathColumn(rows [][]string) (colIdx, startIdx int) {
	header := rows[0]
	if !looksLikeHeader(header) {
		return 0, 0
	}
	for i, cell := range header {
		key := strings.ToLower(strings.TrimSpace(cell))
		if _, ok := pathHeaderAliases[key]; ok {
			return i, 1
		}
	}

	// Header detected but no recognized path-column name — fall back to
	// column 0 and still skip the header row so labels aren't treated
	// as paths.
	return 0, 1
}

// looksLikeHeader returns true when the row contains at least one cell
// matching a known path-column alias.
func looksLikeHeader(row []string) bool {
	for _, cell := range row {
		key := strings.ToLower(strings.TrimSpace(cell))
		if _, ok := pathHeaderAliases[key]; ok {
			return true
		}
	}

	return false
}

// safeCell returns row[idx] trimmed, or "" when the row has fewer columns
// than expected. Prevents index-out-of-range on ragged data rows.
func safeCell(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[idx])
}

// WalkBatchFromDir returns every immediate subdirectory of root that is
// itself a git repository (i.e. contains a `.git` entry).
func WalkBatchFromDir(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(root, entry.Name())
		if isGitRepo(candidate) {
			repos = append(repos, candidate)
		}
	}

	if len(repos) == 0 {
		return nil, ErrBatchEmpty
	}

	return absoluteAndSorted(repos), nil
}

// (IsGitRepo / isGitRepo live in repodetect.go.)

// absoluteAndSorted resolves each input path to an absolute form and
// returns the result sorted lexicographically. Paths that fail to resolve
// are kept as-is so the caller can surface a meaningful per-repo error
// later instead of dropping rows silently.
func absoluteAndSorted(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		abs, err := filepath.Abs(p)
		if err == nil {
			out[i] = abs

			continue
		}
		out[i] = p
	}
	sort.Strings(out)

	return out
}

// (HasGitSubdir lives in repodetect.go.)
