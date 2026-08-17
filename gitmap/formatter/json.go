package formatter

import (
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// WriteJSON writes records to the given writer as a JSON array.
//
// Records are validated first; per-issue warnings are emitted to the
// configured sink (default os.Stderr) but the write always proceeds.
// See validate.go for the warn-and-write policy.
func WriteJSON(w io.Writer, records []model.ScanRecord) error {
	issueCount := emitValidationWarnings(records)

	enc := json.NewEncoder(w)
	enc.SetIndent("", constants.JSONIndent)
	err := enc.Encode(records)
	if err != nil {
		return err
	}
	emitWriteSummary("json", len(records), issueCount)

	return nil
}

// WriteJSONCompact writes a minified JSON containing only essential fields.
func WriteJSONCompact(w io.Writer, records []model.ScanRecord) error {
	issueCount := emitValidationWarnings(records)

	// Essential fields only
	type compactRecord struct {
		RepoName     string `json:"repoName"`
		HTTPSUrl     string `json:"httpsUrl,omitempty"`
		SSHUrl       string `json:"sshUrl,omitempty"`
		Branch       string `json:"branch,omitempty"`
		RelativePath string `json:"relativePath"`
	}

	compacts := make([]compactRecord, 0, len(records))
	for _, r := range records {
		compacts = append(compacts, compactRecord{
			RepoName:     r.RepoName,
			HTTPSUrl:     r.HTTPSUrl,
			SSHUrl:       r.SSHUrl,
			Branch:       r.Branch,
			RelativePath: r.RelativePath,
		})
	}

	enc := json.NewEncoder(w)
	// No SetIndent for minified output
	err := enc.Encode(compacts)
	if err != nil {
		return err
	}
	emitWriteSummary("json", len(records), issueCount)

	return nil
}

// ParseJSON reads records from a JSON reader.
func ParseJSON(reader io.Reader) ([]model.ScanRecord, error) {
	var records []model.ScanRecord
	dec := json.NewDecoder(reader)
	err := dec.Decode(&records)
	if err != nil {
		return nil, err
	}

	return records, nil
}
