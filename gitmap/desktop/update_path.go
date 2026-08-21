package desktop

// UpdateRepoPath informs GitHub Desktop of a repository relocation.
func UpdateRepoPath(oldPath, newPath string) error {
	cli := ResolveCLI()
	if cli == "" {
		return nil
	}
	return addOne(newPath, cli)
}
