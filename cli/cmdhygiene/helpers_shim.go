package cmdhygiene

// isGitRepo is the unexported alias kept for back-compat with existing
// call sites inside this package. New callers should use IsGitRepo.
func isGitRepo(path string) bool {
	return IsGitRepo(path)
}
