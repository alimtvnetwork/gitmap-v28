// Package model — scan_records.go provides loader for serialized ScanRecord data.
package model

import (
	"encoding/json"
	"os"
)

// LoadStatusRecords reads ScanRecords from a JSON file.
func LoadStatusRecords(path string) ([]ScanRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var records []ScanRecord
	err = json.Unmarshal(data, &records)

	return records, err
}
