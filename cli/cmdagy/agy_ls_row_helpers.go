// Package cmdagy — agy_ls_row_helpers.go provides row helper utilities for agy ls.
package cmdagy

func resolveRowBranch(branch string) string {
	if branch == "" {
		return "—"
	}
	return branch
}

func resolveRowStatusAndPath(path string, isPinned bool) (string, string) {
	if path == "" {
		return "global", "—"
	}
	if isPinned {
		return "pinned", path
	}
	return "active", path
}

func getPinnedProjectsSet() map[string]bool {
	set := make(map[string]bool)
	store, err := loadPinnedProjectsStore()
	if err != nil || store == nil {
		return set
	}
	for _, p := range store.Projects {
		set[p.ID] = true
		set[p.Name] = true
	}
	return set
}
