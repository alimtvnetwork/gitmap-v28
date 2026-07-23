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
