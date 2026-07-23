// Package cmd — checkpoint resume for `cfrp` (and `cfr`) batches.
//
// Mirrors the `commit-in` checkpoint contract: each batch writes a
// state.json sidecar after every successful clone, so re-running the
// same command picks up where a crash left off instead of re-cloning
// already-processed entries.
//
// State file location: .gitmap/cfrp/<batch-id>/state.json
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CFRPCheckpoint captures the in-flight state of a cfrp batch.
type CFRPCheckpoint struct {
	BatchID   string    `json:"batch_id"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Total     int       `json:"total"`
	Done      []string  `json:"done"`
	Failed    []string  `json:"failed,omitempty"`
}

// CheckpointPath returns the on-disk sidecar path for a batch.
func CheckpointPath(repoRoot, batchID string) string {
	return filepath.Join(repoRoot, ".gitmap", "cfrp", batchID, "state.json")
}

// LoadCheckpoint reads the sidecar; returns (nil, nil) if missing.
func LoadCheckpoint(path string) (*CFRPCheckpoint, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var cp CFRPCheckpoint
	if err := json.Unmarshal(b, &cp); err != nil {
		return nil, fmt.Errorf("checkpoint %s: %w", path, err)
	}
	return &cp, nil
}

// SaveCheckpoint atomically writes the sidecar (write-then-rename).
func SaveCheckpoint(path string, cp *CFRPCheckpoint) error {
	cp.UpdatedAt = time.Now().UTC()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	b, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// FilterRemaining returns entries from `all` not present in cp.Done.
// Used at the start of a resumed batch to skip completed work.
func FilterRemaining(cp *CFRPCheckpoint, all []string) []string {
	if cp == nil || len(cp.Done) == 0 {
		return all
	}
	done := make(map[string]struct{}, len(cp.Done))
	for _, d := range cp.Done {
		done[d] = struct{}{}
	}
	out := make([]string, 0, len(all))
	for _, e := range all {
		if _, ok := done[e]; !ok {
			out = append(out, e)
		}
	}
	return out
}
