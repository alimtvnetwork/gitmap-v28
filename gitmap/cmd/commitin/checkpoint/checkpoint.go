// Package checkpoint persists per-input progress for commit-in so a
// re-run after a crash or Ctrl-C skips the source SHAs that already
// produced a Created/Skipped outcome in the previous attempt.
//
// Item #3 of the post-v6.53.0 plan. The SQLite runlog already records
// every outcome, but re-walking + re-querying it at startup is slow
// for large repos. A tiny JSON sidecar gives O(1) "already done?"
// lookups at stage boundaries and survives DB corruption.
//
// File layout (one file per staged input):
//
//	<sourceRoot>/.gitmap/commit-in/state/<input-fingerprint>.json
//
// Schema is intentionally minimal — DoneShas is the source of truth;
// LastRunID + UpdatedAt are diagnostic only.
package checkpoint

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DirName is the subdirectory under commit-in that holds state files.
const DirName = "state"

// State is the on-disk JSON shape. Versioned so we can evolve the
// schema without breaking old workspaces.
type State struct {
	Version    int       `json:"version"`
	Input      string    `json:"input"`
	LastRunID  int64     `json:"lastRunId"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DoneShas   []string  `json:"doneShas"`
}

// File wraps the on-disk path plus an in-memory done-set with a mutex
// so concurrent processOneCommit calls (future parallel mode) stay
// safe.
type File struct {
	path string
	mu   sync.Mutex
	done map[string]struct{}
	st   State
}

// Open loads (or initializes) the checkpoint for a given staged input.
// stateDir is typically <Paths.CommitInRoot>/state.
func Open(stateDir, inputFingerprint string, runID int64) (*File, error) {
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(stateDir, inputFingerprint+".json")
	f := &File{path: path, done: map[string]struct{}{}, st: State{
		Version: 1, Input: inputFingerprint, LastRunID: runID, UpdatedAt: time.Now(),
	}}
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return f, nil
	}
	if err != nil {
		return nil, err
	}
	if jsonErr := json.Unmarshal(raw, &f.st); jsonErr != nil {
		// Corrupted state file is not fatal — start fresh.
		f.st = State{Version: 1, Input: inputFingerprint, LastRunID: runID}
		return f, nil
	}
	for _, sha := range f.st.DoneShas {
		f.done[sha] = struct{}{}
	}
	f.st.LastRunID = runID
	return f, nil
}

// IsDone reports whether sha was completed in a previous run.
func (f *File) IsDone(sha string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.done[sha]
	return ok
}

// MarkDone records sha as completed and atomically rewrites state.json.
// Crash-safety: write to <path>.tmp then rename.
func (f *File) MarkDone(sha string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.done[sha]; ok {
		return nil
	}
	f.done[sha] = struct{}{}
	f.st.DoneShas = append(f.st.DoneShas, sha)
	f.st.UpdatedAt = time.Now()
	return f.flushLocked()
}

func (f *File) flushLocked() error {
	raw, err := json.MarshalIndent(f.st, "", "  ")
	if err != nil {
		return err
	}
	tmp := f.path + ".tmp"
	if writeErr := os.WriteFile(tmp, raw, 0o644); writeErr != nil {
		return writeErr
	}
	return os.Rename(tmp, f.path)
}
