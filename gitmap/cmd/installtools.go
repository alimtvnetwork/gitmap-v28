package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// installTool dispatches to the platform-specific installer.
func installTool(opts installOptions) {
	manager := resolvePackageManager(opts.Manager, opts.Tool)
	installCmd := buildInstallCommand(manager, opts.Tool, opts.Version)
	versionLabel := opts.Version
	if versionLabel == "" {
		versionLabel = "latest"
	}
	printInstallPlan(opts.Tool, versionLabel, manager, installCmd)
	if handleDryRunInstall(opts.DryRun, manager, installCmd) {
		return
	}
	executeInstallSteps(installCmd, opts, manager, versionLabel)
}

func printInstallPlan(tool, version, manager string, installCmd []string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: %s\n", version)
	fmt.Printf("  | Manager: %s\n", manager)
	fmt.Printf("  | Command: %s\n", strings.Join(installCmd, " "))
	fmt.Printf("  +--------------------------------------\n\n")
}

func handleDryRunInstall(dryRun bool, manager string, installCmd []string) bool {
	if !dryRun {
		return false
	}
	if manager == constants.PkgMgrApt {
		fmt.Printf(constants.MsgInstallDryCmd, "sudo apt-get update")
	}
	fmt.Printf(constants.MsgInstallDryCmd, strings.Join(installCmd, " "))
	return true
}

func executeInstallSteps(installCmd []string, opts installOptions, manager, versionLabel string) {
	totalSteps := 3
	step := 1
	if manager == constants.PkgMgrApt {
		totalSteps = 4
		fmt.Printf("  [%d/%d] Updating package index...\n", step, totalSteps)
		runAptUpdate(opts.Verbose)
		step++
	}
	runStepInstall(installCmd, opts, manager, versionLabel, step, totalSteps)
}

func runStepInstall(installCmd []string, opts installOptions, manager, versionLabel string, step, totalSteps int) {
	fmt.Printf("  [%d/%d] Installing %s v%s via %s...\n", step, totalSteps, opts.Tool, versionLabel, manager)
	runInstallCommand(installCmd, opts)
	step++
	fmt.Printf("  [%d/%d] Verifying installation...\n", step, totalSteps)
	verifyInstallation(opts.Tool)
	step++
	fmt.Printf("  [%d/%d] Recording installation...\n", step, totalSteps)
	recordInstallation(opts.Tool, manager)
}

// buildInstallCommand builds the install command for a given manager and tool.
func buildInstallCommand(manager, tool, version string) []string {
	pkg := resolvePackageName(manager, tool)
	if cmd, ok := buildUnixInstallCommand(manager, tool, pkg, version); ok {
		return cmd
	}
	return buildWindowsInstallCommand(manager, pkg, version)
}

func buildUnixInstallCommand(manager, tool, pkg, version string) ([]string, bool) {
	switch manager {
	case constants.PkgMgrApt:
		return buildAptCommand(pkg, version), true
	case constants.PkgMgrBrew:
		return buildBrewCommand(tool, pkg), true
	case constants.PkgMgrSnap:
		return buildSnapCommand(pkg), true
	}
	return nil, false
}

func buildWindowsInstallCommand(manager, pkg, version string) []string {
	if manager == constants.PkgMgrWinget {
		return buildWingetCommand(pkg, version)
	}
	return buildChocoCommand(pkg, version)
}

// buildChocoCommand builds a Chocolatey install command.
func buildChocoCommand(pkg, version string) []string {
	args := []string{"choco", "install", pkg, "-y", "--no-progress"}
	if version != "" {
		args = append(args, "--version", version)
	}
	return args
}

// buildWingetCommand builds a Winget install command.
func buildWingetCommand(pkg, version string) []string {
	args := []string{"winget", "install", pkg, "--accept-package-agreements", "--accept-source-agreements", "--silent"}
	if version != "" {
		args = append(args, "--version", version)
	}
	return args
}

// buildAptCommand builds an apt install command.
func buildAptCommand(pkg, version string) []string {
	target := pkg
	if version != "" {
		target = pkg + "=" + version
	}
	return []string{"sudo", "apt", "install", "-y", target}
}

// buildBrewCommand builds a Homebrew install command.
func buildBrewCommand(tool, pkg string) []string {
	if isBrewCaskTool(tool) {
		return []string{"brew", "install", "--cask", pkg}
	}
	return []string{"brew", "install", pkg}
}

// buildSnapCommand builds a Snap install command.
func buildSnapCommand(pkg string) []string {
	if pkg == "code" {
		return []string{"sudo", "snap", "install", pkg, "--classic"}
	}

	return []string{"sudo", "snap", "install", pkg}
}

// isBrewCaskTool returns true for GUI apps needing --cask.
func isBrewCaskTool(tool string) bool {
	switch tool {
	case constants.ToolVSCode, constants.ToolGitHubDesktop, constants.ToolPowerShell,
		constants.ToolDbeaver, constants.ToolOBS:
		return true
	default:
		return false
	}
}

// runAptUpdate runs sudo apt-get update to refresh the package index.
func runAptUpdate(verbose bool) {
	fmt.Print(constants.MsgInstallAptUpdate)
	cmd := exec.Command("sudo", "apt-get", "update")
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInstallAptUpdateFailed, err)
		return
	}
	fmt.Print(constants.MsgInstallAptUpdateDone)
}

// runInstallCommand executes the install command and logs errors.
func runInstallCommand(args []string, opts installOptions) {
	output, err := execInstallCommand(args, opts.Verbose)
	if err != nil {
		handleInstallError(args, opts, output, err)
		return
	}
	fmt.Printf("  ✓ %s install command completed successfully.\n", opts.Tool)
}

func execInstallCommand(args []string, verbose bool) ([]byte, error) {
	cmd := exec.Command(args[0], args[1:]...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return nil, cmd.Run()
	}
	return cmd.CombinedOutput()
}

func handleInstallError(args []string, opts installOptions, output []byte, err error) {
	manager := resolvePackageManager(opts.Manager, opts.Tool)
	logPath := writeInstallErrorLog(opts.Tool, manager, opts.Version, args, output, err)
	printInstallFailureDetails(opts.Tool, manager, opts.Version, args, err, logPath)
	os.Exit(1)
}

func printInstallFailureDetails(tool, manager, version string, args []string, err error, logPath string) {
	versionLabel := version
	if versionLabel == "" {
		versionLabel = "latest"
	}
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailed, tool)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedVersion, versionLabel)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedManager, manager)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedCmd, strings.Join(args, " "))
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedReason, err)
	if logPath != "" {
		fmt.Fprintf(os.Stderr, constants.ErrInstallFailedLog, logPath)
		fmt.Fprint(os.Stderr, constants.ErrInstallFailedHint)
	}
}

// writeInstallErrorLog writes detailed error information to a log file.
func writeInstallErrorLog(tool, manager, version string, args []string, output []byte, installErr error) string {
	logDir := constants.InstallLogDir
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not create log directory %s: %v\n", logDir, err)
		return ""
	}
	logPath := filepath.Join(logDir, fmt.Sprintf("%s-error-%s.log", tool, time.Now().Format("2006-01-02_15-04-05")))
	content := buildInstallErrorLogContent(tool, manager, version, args, output, installErr)
	if err := os.WriteFile(logPath, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not write error log to %s: %v\n", logPath, err)
		return ""
	}
	return logPath
}

func buildInstallErrorLogContent(tool, manager, version string, args []string, output []byte, installErr error) string {
	versionLabel := version
	if versionLabel == "" {
		versionLabel = "latest"
	}
	header := fmt.Sprintf("gitmap install error log\n========================\n\nTool:            %s\nVersion:         %s\nPackage Manager: %s\nCommand:         %s\nTimestamp:       %s\nError:           %v\n\n--- Installer Output ---\n\n",
		tool, versionLabel, manager, strings.Join(args, " "), time.Now().Format(time.RFC3339), installErr)
	if len(output) > 0 {
		return header + string(output)
	}
	return header + "(no output captured — verbose mode pipes directly to stdout/stderr)\n"
}

// resolvePackageName maps tool name to package ID for a manager.
func resolvePackageName(manager, tool string) string {
	switch manager {
	case constants.PkgMgrWinget:
		return resolveWingetPackage(tool)
	case constants.PkgMgrApt:
		return resolveAptPackage(tool)
	case constants.PkgMgrBrew:
		return resolveBrewPackage(tool)
	case constants.PkgMgrSnap:
		return resolveSnapPackage(tool)
	default:
		return resolveChocoPackage(tool)
	}
}

var chocoPackageMap = map[string]string{
	constants.ToolVSCode:        constants.ChocoPkgVSCode,
	constants.ToolNodeJS:        constants.ChocoPkgNodeJS,
	constants.ToolYarn:          constants.ChocoPkgYarn,
	constants.ToolBun:           constants.ChocoPkgBun,
	constants.ToolPnpm:          constants.ChocoPkgPnpm,
	constants.ToolPython:        constants.ChocoPkgPython,
	constants.ToolGo:            constants.ChocoPkgGo,
	constants.ToolGit:           constants.ChocoPkgGit,
	constants.ToolGitLFS:        constants.ChocoPkgGitLFS,
	constants.ToolGHCLI:         constants.ChocoPkgGHCLI,
	constants.ToolGitHubDesktop: constants.ChocoPkgGitHubDesktop,
	constants.ToolCPP:           constants.ChocoPkgCPP,
	constants.ToolPHP:           constants.ChocoPkgPHP,
	constants.ToolPowerShell:    constants.ChocoPkgPowerShell,
	constants.ToolMySQL:         constants.ChocoPkgMySQL,
	constants.ToolMariaDB:       constants.ChocoPkgMariaDB,
	constants.ToolPostgreSQL:    constants.ChocoPkgPostgreSQL,
	constants.ToolSQLite:        constants.ChocoPkgSQLite,
	constants.ToolMongoDB:       constants.ChocoPkgMongoDB,
	constants.ToolCouchDB:       constants.ChocoPkgCouchDB,
	constants.ToolRedis:         constants.ChocoPkgRedis,
	constants.ToolNeo4j:         constants.ChocoPkgNeo4j,
	constants.ToolElasticsearch: constants.ChocoPkgElasticsearch,
	constants.ToolDuckDB:        constants.ChocoPkgDuckDB,
	constants.ToolNpp:           constants.ChocoPkgNpp,
	constants.ToolNppInstall:    constants.ChocoPkgNpp,
	constants.ToolDbeaver:       constants.ChocoPkgDbeaver,
	constants.ToolOBS:           constants.ChocoPkgOBS,
}

// resolveChocoPackage maps tool names to Chocolatey package IDs.
func resolveChocoPackage(tool string) string {
	if pkg, exists := chocoPackageMap[tool]; exists {
		return pkg
	}
	return tool
}

var wingetPackageMap = map[string]string{
	constants.ToolVSCode:        constants.WingetPkgVSCode,
	constants.ToolPowerShell:    constants.WingetPkgPowerShell,
	constants.ToolDbeaver:       constants.WingetPkgDbeaver,
	constants.ToolOBS:           constants.WingetPkgOBS,
	constants.ToolStickyNotes:   constants.WingetPkgStickyNotes,
	constants.ToolGitHubDesktop: constants.WingetPkgGitHubDesktop,
}

// resolveWingetPackage maps tool names to Winget package IDs.
func resolveWingetPackage(tool string) string {
	if pkg, exists := wingetPackageMap[tool]; exists {
		return pkg
	}
	return tool
}

var aptPackageMap = map[string]string{
	constants.ToolNodeJS:        constants.AptPkgNodeJS,
	constants.ToolPython:        constants.AptPkgPython,
	constants.ToolGo:            constants.AptPkgGo,
	constants.ToolGit:           constants.AptPkgGit,
	constants.ToolGitLFS:        constants.AptPkgGitLFS,
	constants.ToolCPP:           constants.AptPkgCPP,
	constants.ToolPHP:           constants.AptPkgPHP,
	constants.ToolMySQL:         constants.AptPkgMySQL,
	constants.ToolMariaDB:       constants.AptPkgMariaDB,
	constants.ToolPostgreSQL:    constants.AptPkgPostgreSQL,
	constants.ToolSQLite:        constants.AptPkgSQLite,
	constants.ToolMongoDB:       constants.AptPkgMongoDB,
	constants.ToolCouchDB:       constants.AptPkgCouchDB,
	constants.ToolRedis:         constants.AptPkgRedis,
	constants.ToolCassandra:     constants.AptPkgCassandra,
	constants.ToolElasticsearch: constants.AptPkgElasticsearch,
}

// resolveAptPackage maps tool names to apt package IDs.
func resolveAptPackage(tool string) string {
	if pkg, exists := aptPackageMap[tool]; exists {
		return pkg
	}
	return tool
}

var brewPackageMap = map[string]string{
	constants.ToolNodeJS:        constants.BrewPkgNodeJS,
	constants.ToolPython:        constants.BrewPkgPython,
	constants.ToolGo:            constants.BrewPkgGo,
	constants.ToolGit:           constants.BrewPkgGit,
	constants.ToolGitLFS:        constants.BrewPkgGitLFS,
	constants.ToolGHCLI:         constants.BrewPkgGHCLI,
	constants.ToolCPP:           constants.BrewPkgCPP,
	constants.ToolPHP:           constants.BrewPkgPHP,
	constants.ToolMySQL:         constants.BrewPkgMySQL,
	constants.ToolMariaDB:       constants.BrewPkgMariaDB,
	constants.ToolPostgreSQL:    constants.BrewPkgPostgreSQL,
	constants.ToolSQLite:        constants.BrewPkgSQLite,
	constants.ToolMongoDB:       constants.BrewPkgMongoDB,
	constants.ToolCouchDB:       constants.BrewPkgCouchDB,
	constants.ToolRedis:         constants.BrewPkgRedis,
	constants.ToolNeo4j:         constants.BrewPkgNeo4j,
	constants.ToolElasticsearch: constants.BrewPkgElasticsearch,
	constants.ToolDuckDB:        constants.BrewPkgDuckDB,
	constants.ToolDbeaver:       constants.BrewPkgDbeaver,
	constants.ToolOBS:           constants.BrewPkgOBS,
}

// resolveBrewPackage maps tool names to Homebrew package IDs.
func resolveBrewPackage(tool string) string {
	if pkg, exists := brewPackageMap[tool]; exists {
		return pkg
	}
	return tool
}

var snapPackageMap = map[string]string{
	constants.ToolCouchDB: constants.SnapPkgCouchDB,
	constants.ToolRedis:   constants.SnapPkgRedis,
	constants.ToolVSCode:  "code",
}

// resolveSnapPackage maps tool names to Snap package IDs.
func resolveSnapPackage(tool string) string {
	if pkg, exists := snapPackageMap[tool]; exists {
		return pkg
	}
	return tool
}

// recordInstallation saves the install record to the database.
func recordInstallation(tool, manager string) {
	version := detectInstalledVersion(tool)
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()
	if err := db.SaveInstalledTool(tool, version, manager); err == nil && version != "" {
		fmt.Printf(constants.MsgInstallRecorded, tool, version)
	}
}
