package cmdgit

// RunShow handles native git show execution without shell wrapping.
func RunShow(args []string) error {
	targetDir, remaining := ParseGitDirAndArgs(args)

	gitArgs := append([]string{"show"}, remaining...)

	return ExecGitInspection(targetDir, gitArgs...)
}
