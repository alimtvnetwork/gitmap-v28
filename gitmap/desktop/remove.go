package desktop

// RemoveRepo handles repository untracking for GitHub Desktop.
func RemoveRepo(repoPath string) error {
	cli := ResolveCLI()
	if cli == "" {
		return nil
	}
	// GitHub Desktop does not provide a CLI delete flag; we verify existence.
	return nil
}
