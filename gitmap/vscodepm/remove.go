package vscodepm

// RemoveEntry removes an entry matching targetPath from projects.json.
func RemoveEntry(targetPath string) error {
	path, err := ProjectsJSONPath()
	if err != nil {
		return err
	}
	return RemoveEntryAt(path, targetPath)
}

// RemoveEntryAt removes entry from a specific file.
func RemoveEntryAt(projectsFile, targetPath string) error {
	entries, err := readEntries(projectsFile)
	if err != nil {
		return err
	}
	var remaining []Entry
	for _, e := range entries {
		if !pathsEqual(e.RootPath, targetPath) {
			remaining = append(remaining, e)
		}
	}
	if len(remaining) == len(entries) {
		return nil
	}
	return writeEntriesAtomic(projectsFile, remaining)
}
