package vscodepm

import (
	"fmt"
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
	entries, err := readEntries(projectsFile)
	if err != nil {
		return err
	}
	found := false
	for i := range entries {
		if pathsEqual(entries[i].RootPath, oldPath) {
			entries[i].RootPath = newPath
			if len(newName) > 0 {
				entries[i].Name = newName
			}
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("project not found with rootPath %s", oldPath)
	}
	return writeEntriesAtomic(projectsFile, entries)
}
