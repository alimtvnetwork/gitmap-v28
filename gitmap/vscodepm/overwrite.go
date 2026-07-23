package vscodepm

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// OverwritePaths sets the Paths field of the entry matching rootPath to
// EXACTLY the supplied slice (no union). Used by `gitmap code paths rm` so
// removed paths actually disappear from projects.json instead of being
// re-merged back in by the regular Sync union semantics.
//
// If the entry does not exist yet, it is created with the supplied paths.
// Other fields (Tags / Enabled / Profile) are preserved on existing rows.
func OverwritePaths(rootPath, name string, paths []string) error {
	path, err := ProjectsJSONPath()
	if err != nil {
		return err
	}

	entries, err := readEntries(path)
	if err != nil {
		return err
	}

	if paths == nil {
		paths = []string{}
	}

	found := false
	for i := range entries {
		if pathsEqual(entries[i].RootPath, rootPath) {
			entries[i].Name = name
			entries[i].Paths = paths
			found = true

			break
		}
	}

	if !found {
		entries = append(entries, newEntry(rootPath, name, paths, nil))
	}

	if err := writeEntriesAtomic(path, entries); err != nil {
		return fmt.Errorf("%s: %w", constants.VSCodePMProjectsFile, err)
	}

	return nil
}
