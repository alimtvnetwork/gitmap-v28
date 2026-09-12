// Package cmdsetup provides setup subcommands for development environments,
// WordPress, Laravel, Ubuntu toolchains, and file permissions.
package cmdsetup

// RunSetup handles the root setup command logic.
func RunSetup(args []string) error {
	return runSetup(args)
}

// RunSetupPerms handles setup permissions command logic.
func RunSetupPerms(args []string) error {
	return runSetupPerms(args)
}

// RunSetupWordpress handles setup wordpress command logic.
func RunSetupWordpress(args []string) error {
	return runSetupWordpress(args)
}

// RunSetupLaravel handles setup laravel command logic.
func RunSetupLaravel(args []string) error {
	return runSetupLaravel(args)
}

// RunPrintPathSnippet executes the print path snippet command logic.
func RunPrintPathSnippet(args []string) {
	runPrintPathSnippet(args)
}

// WarnIfNoWrapper prints a stderr warning when cd is called without wrapper.
func WarnIfNoWrapper() {
	warnIfNoWrapper()
}

// ResolveSetupConfigPath resolves the setup config path with bundled fallback.
func ResolveSetupConfigPath(configPath string, hasConfig bool) string {
	return resolveSetupConfigPath(configPath, hasConfig)
}

// IsWrapperActive checks whether the shell wrapper is active in the environment.
func IsWrapperActive() bool {
	return isWrapperActive()
}
