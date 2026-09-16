package cmdclone

// RunCloneOnlyMissing handles the "clone-only-missing" subcommand.
func RunCloneOnlyMissing(args []string) error {
	checkHelp("clone-only-missing", args)
	cf := parseCloneFlags(args)
	cf.MissingOnly = true
	cf.NoReplace = true

	return executeParsedClone(cf)
}
