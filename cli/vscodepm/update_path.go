package vscodepm

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// UpdateRootPath updates the rootPath and optionally name in projects.json.
func UpdateRootPath(oldPath, newPath, newName string) error {
	path, err := ProjectsJSONPath()
	if err != nil {
		return err
	}

	return UpdateRootPathAt(path, oldPath, newPath, newName)
}

// UpdateRootPathAt modifies the specific projects.json file.
func UpdateRootPathAt(projectsFile, oldPath, newPath, newName string) error {
	entries, err := ReadEntries(projectsFile)
	if err != nil {
		return err
	}

	found := false
	for i := range entries {
		if !pathsEqual(entries[i].RootPath, oldPath) {
			continue
		}

		entries[i].RootPath = newPath
		if len(newName) > 0 {
			entries[i].Name = newName
		}

		found = true
		break
	}

	if !found {
		return apperror.NewNotFoundError(fmt.Sprintf("project not found with rootPath %s", oldPath))
	}

	return writeEntriesAtomic(projectsFile, entries)
}
