package cmdinstall

// CheckHelpFn delegates help checking to cmd package.
var CheckHelpFn func(cmd string, args []string)

// RemoteAgmUpdateFn delegates remote AGM update to cmdssh package.
var RemoteAgmUpdateFn func(target string) error

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

// RunUpdateAgManagerWithOpts updates Antigravity Manager with options.
func RunUpdateAgManagerWithOpts(opts installOptions) error {
	return runUpdateAgManagerWithOpts(opts)
}

// RunUpdateAgManager updates Antigravity Manager.
func RunUpdateAgManager() error {
	return runUpdateAgManagerWithOpts(installOptions{})
}

// RunInstallAgyWithOpts installs Antigravity CLI.
func RunInstallAgyWithOpts(opts installOptions) error {
	return runInstallAgyWithOpts(opts)
}

// RunInstallAntigravityWithOpts installs Antigravity IDE.
func RunInstallAntigravityWithOpts(opts installOptions) error {
	return runInstallAntigravityWithOpts(opts)
}

// RunInstallTar installs an archive package on Linux.
func RunInstallTar(args []string) error {
	return runInstallTar(args)
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

// ResolveToolAlias normalizes known tool aliases.
func ResolveToolAlias(tool string) string {
	return resolveToolAlias(tool)
}

// ResolvePackageName resolves the package name for a tool.
func ResolvePackageName(tool, manager string) string {
	return resolvePackageName(manager, tool)
}

// RunInstallCommand executes the package manager install command.
func RunInstallCommand(args []string, opts installOptions) error {
	return runInstallCommand(args, opts)
}

// CheckCustomStandaloneTool reports whether a tool is a standalone custom tool.
func CheckCustomStandaloneTool(tool string) bool {
	return IsCustomStandaloneTool(tool)
}

// UninstallCustomTool removes custom tool binaries and records.
func UninstallCustomTool(tool string, purge bool) error {
	return RunUninstallCustomTool(tool, purge)
}

// FindInstalledAntigravityDesktopPath returns the desktop executable path if installed.
func FindInstalledAntigravityDesktopPath() (string, bool) {
	return findInstalledAntigravityDesktopPath()
}

// ResolveToolBinaryPath resolves the binary path for a tool with fallback paths.
func ResolveToolBinaryPath(binary string) string {
	return resolveToolBinaryPath(binary)
}
