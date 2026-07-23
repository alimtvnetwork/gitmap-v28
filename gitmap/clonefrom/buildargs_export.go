package clonefrom

// buildargs_export.go — exported wrapper over buildGitArgs so the
// cmd package's --verify-cmd-faithful checker can compare the
// rendered `cmd:` line against the EXACT argv this executor would
// hand to exec.Command. See clonenow/buildargs_export.go for the
// full rationale (single-source-of-truth for git-clone argv).

// BuildGitArgs returns the argv (excluding the "git" binary) that
// Execute would pass to exec.Command for the given row.
func BuildGitArgs(r Row, dest string) []string {
	return buildGitArgs(r, dest)
}
