package cmdinstall

// CheckHelpFn delegates help checking to cmd package.
var CheckHelpFn func(cmd string, args []string)

// InstallOptions is the alias for installOptions.
type InstallOptions = installOptions

// RunInstall executes the install command logic.
func RunInstall(args []string) error {
	return runInstall(args)
}

// RunInstalledDir executes the installed-dir command logic.
func RunInstalledDir() error {
	return runInstalledDir()
}

// InstallTool installs the specified tool.
func InstallTool(opts installOptions) {
	installTool(opts)
}

// RunInstallAdd executes the install add command logic.
func RunInstallAdd(args []string) error {
	return runInstallAdd(args)
}

// RunInstallChromeLinux installs Chrome on Linux.
func RunInstallChromeLinux(opts installOptions) error {
	return runInstallChromeLinux(opts)
}

// RunInstallAgManagerWithOpts installs Antigravity Manager.
func RunInstallAgManagerWithOpts(opts installOptions) error {
	return runInstallAgManagerWithOpts(opts)
}

// RunInstallAgyWithOpts installs Antigravity CLI.
func RunInstallAgyWithOpts(opts installOptions) error {
	return runInstallAgyWithOpts(opts)
}

// RunInstallAntigravityWithOpts installs Antigravity IDE.
func RunInstallAntigravityWithOpts(opts installOptions) error {
	return runInstallAntigravityWithOpts(opts)
}

// ResolvePowerShellBinaryWithLookPath resolves PowerShell binary with a custom lookPath.
func ResolvePowerShellBinaryWithLookPath(lookPath func(string) (string, error)) string {
	return resolvePowerShellBinaryWithLookPath(lookPath)
}

// GetInstalledVersion returns the installed version of a binary.
func GetInstalledVersion(binary string) string {
	return getInstalledVersion(binary)
}

// ResolveProfileTree resolves a profile composition by name.
func ResolveProfileTree(slug string) (ProfileComposition, bool) {
	return resolveProfileTree(slug)
}

// PrintProfileTree prints the tree representation of a profile.
func PrintProfileTree(prof ProfileComposition) {
	printProfileTree(prof)
}

// PrintProfileInstallSummary prints the profile install summary.
func PrintProfileInstallSummary(slug string) {
	printProfileInstallSummary(slug)
}

// ValidateToolName validates the given tool name.
func ValidateToolName(tool string) {
	validateToolName(tool)
}

// RunUninstallCtx uninstalls desktop context menu entries.
func RunUninstallCtx() {
	runUninstallCtx()
}

// ResolvePackageManager detects the system package manager.
func ResolvePackageManager(override, tool string) string {
	return resolvePackageManager(override, tool)
}

// ResolvePackageName resolves the package name for a tool.
func ResolvePackageName(tool, manager string) string {
	return resolvePackageName(tool, manager)
}

// RunInstallCommand executes the package manager install command.
func RunInstallCommand(args []string, opts installOptions) error {
	return runInstallCommand(args, opts)
}
