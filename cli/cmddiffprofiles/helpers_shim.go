package cmddiffprofiles

// profileExists checks if a profile name is in the list.
func profileExists(profiles []string, name string) bool {
	for _, p := range profiles {
		if p == name {
			return true
		}
	}

	return false
}
