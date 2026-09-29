package cmd

// runCommitCLI handles transparent git commit execution with help interception.
func runCommitCLI(args []string) error {
	checkHelp("commit", args)
	if len(args) == 0 {
		printHelpAndExit("commit", args)
		return nil
	}

	return runGitPassthrough(append([]string{"commit"}, args...))
}
