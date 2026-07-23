package vscodepm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Pair is one (rootPath, name, paths, tags) tuple to upsert into
// projects.json. Both Paths (multi-root, v3.39.0+) and Tags (auto-derived,
// v3.40.0+) are UNIONed with whatever the user already has on disk so
// gitmap never silently removes a user-added path or tag.
type Pair struct {
	RootPath string
	Name     string
	Paths    []string
	Tags     []string
}

// Sync reconciles projects.json with the supplied DB-side pairs.
//
// Behavior:
//   - New rootPath -> append a default Entry with Paths = pair.Paths.
//   - Existing rootPath -> update Name. Paths becomes UNION(existing, pair.Paths).
//     Tags / Enabled / Profile are preserved (so user edits in the VS Code UI
//     survive untouched).
//   - Foreign entries (rootPath not in pairs) -> preserved verbatim.
//
// Writes are atomic: temp file in the same directory then os.Rename.
// Returns ErrUserDataMissing / ErrExtensionMissing when the path
// cannot be resolved — callers should treat those as soft skips.
func Sync(pairs []Pair) (SyncSummary, error) {
	return SyncMode(pairs, MergeModeUnion)
}

// SyncMode is Sync with an explicit merge strategy. Same path
// resolution and atomic-write semantics; only the per-entry tag
// reconciliation changes per the MergeMode dispatcher.
func SyncMode(pairs []Pair, mode MergeMode) (SyncSummary, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return SyncSummary{}, err
	}

	return SyncAtMode(path, pairs, mode)
}

// RenameByPath updates the Name field of the entry whose rootPath matches.
// Paths / Tags / Enabled / Profile are intentionally left alone.
// Returns true when an entry was actually renamed (false = no-op).
func RenameByPath(rootPath, newName string) (bool, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return false, err
	}

	entries, err := readEntries(path)
	if err != nil {
		return false, err
	}

	changed := false

	for i := range entries {
		if pathsEqual(entries[i].RootPath, rootPath) && entries[i].Name != newName {
			entries[i].Name = newName
			changed = true
		}
	}

	if !changed {
		return false, nil
	}

	return true, writeEntriesAtomic(path, entries)
}

// ListEntries returns every Entry currently persisted in projects.json
// without mutating the file. Used by `gitmap vscode-pm-sync` to walk
// the existing entries and rebuild Pair tuples for re-tagging. Missing
// projects.json => empty slice + nil error so callers can soft-report
// "nothing to do" instead of erroring out.
func ListEntries() ([]Entry, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return nil, err
	}

	return readEntries(path)
}

// readEntries returns the parsed entries. Missing file -> empty slice.
func readEntries(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Entry{}, nil
		}

		return nil, fmt.Errorf(constants.ErrVSCodePMReadFailed, path, err)
	}

	if len(data) == 0 {
		return []Entry{}, nil
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodePMParseFailed, path, err)
	}

	for i := range entries {
		entries[i] = ensureSlices(entries[i])
	}

	return entries, nil
}
