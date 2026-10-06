package cmdlogin

// Run executes the gitmap login CLI entry point.
func Run(args []string) error {
	return runLogin(args)
}

// RunLogout executes the gitmap logout CLI entry point.
func RunLogout(args []string) error {
	return runLogout(args)
}
