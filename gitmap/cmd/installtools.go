package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// installTool dispatches to the platform-specific installer.
func installTool(opts installOptions) {
	manager := resolvePackageManager(opts.Manager)
	installCmd := buildInstallCommand(manager, opts.Tool, opts.Version)

	versionLabel := opts.Version
	if versionLabel == "" {
		versionLabel = "latest"
	}

	// Step 1: Show install plan.
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", opts.Tool)
	fmt.Printf("  | Version: %s\n", versionLabel)
	fmt.Printf("  | Manager: %s\n", manager)
	fmt.Printf("  | Command: %s\n", strings.Join(installCmd, " "))
	fmt.Printf("  +--------------------------------------\n\n")

	if opts.DryRun {
		if manager == constants.PkgMgrApt {
			fmt.Printf(constants.MsgInstallDryCmd, "sudo apt-get update")
		}

		fmt.Printf(constants.MsgInstallDryCmd, strings.Join(installCmd, " "))

		return
	}

	totalSteps := 3
	step := 1

	// Step: Package index update (apt only).
	if manager == constants.PkgMgrApt {
		totalSteps = 4
		fmt.Printf("  [%d/%d] Updating package index...\n", step, totalSteps)
		runAptUpdate(opts.Verbose)
		step++
	}

	// Step: Install.
	fmt.Printf("  [%d/%d] Installing %s v%s via %s...\n", step, totalSteps, opts.Tool, versionLabel, manager)
	runInstallCommand(installCmd, opts)
	step++

	// Step: Verify.
	fmt.Printf("  [%d/%d] Verifying installation...\n", step, totalSteps)
	verifyInstallation(opts.Tool)
	step++

	// Step: Record in database.
	fmt.Printf("  [%d/%d] Recording installation...\n", step, totalSteps)
	recordInstallation(opts.Tool, manager)
}

// buildInstallCommand builds the install command for a given manager and tool.
func buildInstallCommand(manager, tool, version string) []string {
	pkgName := resolvePackageName(manager, tool)

	if manager == constants.PkgMgrChocolatey {
		return buildChocoCommand(pkgName, version)
	}
	if manager == constants.PkgMgrWinget {
		return buildWingetCommand(pkgName, version)
	}
	if manager == constants.PkgMgrApt {
		return buildAptCommand(pkgName, version)
	}
	if manager == constants.PkgMgrBrew {
		return buildBrewCommand(tool, pkgName)
	}
	if manager == constants.PkgMgrSnap {
		return buildSnapCommand(pkgName)
	}

	return buildChocoCommand(pkgName, version)
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
	return []string{"sudo", "snap", "install", pkg}
}

// isBrewCaskTool returns true for GUI apps needing --cask.
func isBrewCaskTool(tool string) bool {
	if tool == constants.ToolVSCode {
		return true
	}
	if tool == constants.ToolGitHubDesktop {
		return true
	}
	if tool == constants.ToolPowerShell {
		return true
	}
	if tool == constants.ToolDbeaver {
		return true
	}
	if tool == constants.ToolOBS {
		return true
	}

	return false
}

// runAptUpdate runs sudo apt-get update to refresh the package index.
func runAptUpdate(verbose bool) {
	fmt.Print(constants.MsgInstallAptUpdate)

	cmd := exec.Command("sudo", "apt-get", "update")

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInstallAptUpdateFailed, err)

		return
	}

	fmt.Print(constants.MsgInstallAptUpdateDone)
}

// runInstallCommand executes the install command and logs errors.
func runInstallCommand(args []string, opts installOptions) {
	cmd := exec.Command(args[0], args[1:]...)

	var output []byte
	var err error

	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	} else {
		output, err = cmd.CombinedOutput()
	}

	if err != nil {
		manager := resolvePackageManager(opts.Manager)
		logPath := writeInstallErrorLog(opts.Tool, manager, opts.Version, args, output, err)

		fmt.Fprintf(os.Stderr, constants.ErrInstallFailed, opts.Tool)

		versionLabel := opts.Version
		if versionLabel == "" {
			versionLabel = "latest"
		}

		fmt.Fprintf(os.Stderr, constants.ErrInstallFailedVersion, versionLabel)
		fmt.Fprintf(os.Stderr, constants.ErrInstallFailedManager, manager)
		fmt.Fprintf(os.Stderr, constants.ErrInstallFailedCmd, strings.Join(args, " "))
		fmt.Fprintf(os.Stderr, constants.ErrInstallFailedReason, err)

		if logPath != "" {
			fmt.Fprintf(os.Stderr, constants.ErrInstallFailedLog, logPath)
			fmt.Fprint(os.Stderr, constants.ErrInstallFailedHint)
		}

		os.Exit(1)
	}

	fmt.Printf("  ✓ %s install command completed successfully.\n", opts.Tool)
}

// writeInstallErrorLog writes detailed error information to a log file.
func writeInstallErrorLog(tool, manager, version string, args []string, output []byte, installErr error) string {
	logDir := constants.InstallLogDir

	err := os.MkdirAll(logDir, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not create log directory %s: %v\n", logDir, err)

		return ""
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s-error-%s.log", tool, timestamp)
	logPath := filepath.Join(logDir, filename)

	versionLabel := version
	if versionLabel == "" {
		versionLabel = "latest"
	}

	var sb strings.Builder

	sb.WriteString("gitmap install error log\n")
	sb.WriteString("========================\n\n")
	sb.WriteString(fmt.Sprintf("Tool:            %s\n", tool))
	sb.WriteString(fmt.Sprintf("Version:         %s\n", versionLabel))
	sb.WriteString(fmt.Sprintf("Package Manager: %s\n", manager))
	sb.WriteString(fmt.Sprintf("Command:         %s\n", strings.Join(args, " ")))
	sb.WriteString(fmt.Sprintf("Timestamp:       %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Error:           %v\n", installErr))
	sb.WriteString("\n--- Installer Output ---\n\n")

	if len(output) > 0 {
		sb.Write(output)
	} else {
		sb.WriteString("(no output captured — verbose mode pipes directly to stdout/stderr)\n")
	}

	err = os.WriteFile(logPath, []byte(sb.String()), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not write error log to %s: %v\n", logPath, err)

		return ""
	}

	return logPath
}

// resolvePackageName maps tool name to package ID for a manager.
func resolvePackageName(manager, tool string) string {
	if manager == constants.PkgMgrWinget {
		return resolveWingetPackage(tool)
	}
	if manager == constants.PkgMgrApt {
		return resolveAptPackage(tool)
	}
	if manager == constants.PkgMgrBrew {
		return resolveBrewPackage(tool)
	}
	if manager == constants.PkgMgrSnap {
		return resolveSnapPackage(tool)
	}

	return resolveChocoPackage(tool)
}

// resolveChocoPackage maps tool names to Chocolatey package IDs.
func resolveChocoPackage(tool string) string {
	chocoMap := map[string]string{
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

	pkg, exists := chocoMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveWingetPackage maps tool names to Winget package IDs.
func resolveWingetPackage(tool string) string {
	wingetMap := map[string]string{
		constants.ToolVSCode:        constants.WingetPkgVSCode,
		constants.ToolPowerShell:    constants.WingetPkgPowerShell,
		constants.ToolDbeaver:       constants.WingetPkgDbeaver,
		constants.ToolOBS:           constants.WingetPkgOBS,
		constants.ToolStickyNotes:   constants.WingetPkgStickyNotes,
		constants.ToolGitHubDesktop: constants.WingetPkgGitHubDesktop,
	}

	pkg, exists := wingetMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveAptPackage maps tool names to apt package IDs.
func resolveAptPackage(tool string) string {
	aptMap := map[string]string{
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

	pkg, exists := aptMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveBrewPackage maps tool names to Homebrew package IDs.
func resolveBrewPackage(tool string) string {
	brewMap := map[string]string{
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

	pkg, exists := brewMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveSnapPackage maps tool names to Snap package IDs.
func resolveSnapPackage(tool string) string {
	snapMap := map[string]string{
		constants.ToolCouchDB: constants.SnapPkgCouchDB,
		constants.ToolRedis:   constants.SnapPkgRedis,
	}

	pkg, exists := snapMap[tool]
	if exists {
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

	err = db.SaveInstalledTool(tool, version, manager)
	if err != nil {
		return
	}

	if version != "" {
		fmt.Printf(constants.MsgInstallRecorded, tool, version)
	}
}
