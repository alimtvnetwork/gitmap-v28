package vscodepm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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
	path, err := EnsureProjectsJSONPath()
	if err != nil {
		return SyncSummary{}, err
	}

	filtered := filterProductionPairs(pairs)
	if len(filtered) == 0 {
		return SyncSummary{}, nil
	}

	return SyncAtMode(path, filtered, mode)
}

func filterProductionPairs(pairs []Pair) []Pair {
	out := make([]Pair, 0, len(pairs))
	for _, p := range pairs {
		if IsDisallowedProjectPath(p.RootPath) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// RenameByPath updates the Name field of the entry whose rootPath matches.
// Paths / Tags / Enabled / Profile are intentionally left alone.
// Returns true when an entry was actually renamed (false = no-op).
func RenameByPath(rootPath, newName string) (bool, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return false, err
	}

	entriesRes := readEntries(path)
	if entriesRes.HasError() {
		return false, entriesRes.Err
	}
	entries := entriesRes.Items()

	isChanged := false

	for i := range entries {
		if pathsEqual(entries[i].RootPath, rootPath) && entries[i].Name != newName {
			entries[i].Name = newName
			isChanged = true
		}
	}

	if !isChanged {
		return false, nil
	}

	return true, writeEntriesAtomic(path, entries)
}

// ListEntriesWrap returns every Entry wrapped in EntrySliceWrap currently persisted in projects.json.
func ListEntriesWrap() EntrySliceWrap {
	path, err := ProjectsJSONPath()
	if err != nil {
		return result.FailSlice[Entry](err)
	}
	return ReadEntriesWrap(path)
}

// ReadEntriesWrap returns the parsed entries wrapped in EntrySliceWrap.
func ReadEntriesWrap(path string) EntrySliceWrap {
	return readEntries(path)
}

// ListEntries returns every Entry currently persisted in projects.json
// without mutating the file. Used by `gitmap vscode-pm-sync` to walk
// the existing entries and rebuild Pair tuples for re-tagging. Missing
// projects.json => empty slice + nil error so callers can soft-report
// "nothing to do" instead of erroring out.
func ListEntries() ([]Entry, error) {
	res := ListEntriesWrap()
	if res.HasError() {
		return nil, res.Err
	}
	return res.Items(), nil
}

// ReadEntries returns the parsed entries. Missing file -> empty slice.
func ReadEntries(path string) ([]Entry, error) {
	res := ReadEntriesWrap(path)
	if res.HasError() {
		return nil, res.Err
	}
	return res.Items(), nil
}

// readEntries returns the parsed entries. Missing file -> empty slice.
func readEntries(path string) EntrySliceWrap {
	data, err := os.ReadFile(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return result.OkSlice([]Entry{})
	}

	if err != nil {
		return result.FailSlice[Entry](apperror.WrapSimple(err, fmt.Sprintf("vscode: failed to read %s", path)))
	}

	if len(data) == 0 {
		return result.OkSlice([]Entry{})
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return result.FailSlice[Entry](apperror.WrapSimple(err, fmt.Sprintf("vscode: %s is not valid JSON (left untouched)", path)))
	}

	for i := range entries {
		entries[i] = ensureSlices(entries[i])
	}

	return result.OkSlice(entries)
}
