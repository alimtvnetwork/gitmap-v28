package cmdworkdir

var (
	CheckHelpFn func(command string, args []string)
)

// RunWorkDir handles gitmap workdir CLI commands.
func RunWorkDir(args []string) error {
	return runWorkDir(args)
}

// IsWorkDirKeyword checks if a string is a workdir keyword.
func IsWorkDirKeyword(name string) bool {
	return isWorkDirKeyword(name)
}
