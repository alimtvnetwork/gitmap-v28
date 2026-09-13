package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

// IPSnapshotRecord represents a persisted IP snapshot in SQLite or JSON fallback.
type IPSnapshotRecord struct {
	IPSnapshotId  int64  `json:"ipSnapshotId"`
	InterfaceName string `json:"interfaceName"`
	IP            string `json:"ip"`
	Netmask       string `json:"netmask"`
	Gateway       string `json:"gateway"`
	DNS           string `json:"dns"`
	IsDHCP        bool   `json:"isDHCP"`
	Timestamp     int64  `json:"timestamp"`
	Notes         string `json:"notes,omitempty"`
	Comments      string `json:"comments,omitempty"`
}

// GetIPSnapshotFallbackPath returns the canonical path to the JSON fallback file.
func GetIPSnapshotFallbackPath() string {
	dir := tempdir.RepoTempDir()

	return filepath.Join(dir, "ip-snapshot-last.json")
}

// SaveIPSnapshotJSON persists the snapshot record to disk in the repo temp dir.
func SaveIPSnapshotJSON(record *IPSnapshotRecord) *apperror.AppError {
	if record == nil {
		return apperror.NewSimple("store.SaveIPSnapshotJSON", "E_NIL_RECORD")
	}

	payload, err := marshalSnapshotJSON(record)
	if err != nil {
		return err
	}

	return writeSnapshotJSONFile(payload)
}

func marshalSnapshotJSON(record *IPSnapshotRecord) ([]byte, *apperror.AppError) {
	payload, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.marshalSnapshotJSON")
	}

	return payload, nil
}

func writeSnapshotJSONFile(payload []byte) *apperror.AppError {
	filePath := GetIPSnapshotFallbackPath()
	if writeErr := os.WriteFile(filePath, payload, 0o600); writeErr != nil {
		return apperror.WrapSimple(writeErr, "store.writeSnapshotJSONFile")
	}

	return nil
}

// LoadLatestIPSnapshotJSON loads the latest snapshot from the JSON fallback file.
func LoadLatestIPSnapshotJSON() (*IPSnapshotRecord, *apperror.AppError) {
	filePath := GetIPSnapshotFallbackPath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.LoadLatestIPSnapshotJSON.ReadFile")
	}

	return unmarshalSnapshotJSON(data)
}

func unmarshalSnapshotJSON(data []byte) (*IPSnapshotRecord, *apperror.AppError) {
	var record IPSnapshotRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, apperror.WrapSimple(err, "store.unmarshalSnapshotJSON")
	}

	return &record, nil
}
