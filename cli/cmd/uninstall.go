package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runUninstall handles the "uninstall" command.
func runUninstall(args []string) error {
	checkHelp("uninstall", args)
	if !hasPositionalToolArg(args) {
		runSelfUninstall(args)
		return nil
	}

	fs := flag.NewFlagSet("uninstall", flag.ExitOnError)
	var dryRun, force, purge bool
	fs.BoolVar(&dryRun, constants.FlagUninstallDryRun, false, constants.FlagDescUninstallDryRun)
	fs.BoolVar(&force, constants.FlagUninstallForce, false, constants.FlagDescUninstallForce)
	fs.BoolVar(&purge, constants.FlagUninstallPurge, false, constants.FlagDescUninstallPurge)
	fs.Parse(reorderFlagsBeforeArgs(args))

	tool := fs.Arg(0)
	if tool == "" || tool == "gitmap" || tool == "gitmap-cli" {
		runSelfUninstall(args)
		return nil
	}

	validateToolName(tool)
	canonical := resolveToolAlias(tool)
	if tool == constants.ToolCtx || canonical == constants.ToolCtx {
		runUninstallCtx()
		return nil
	}

	db, _ := openDB()
	if db != nil {
		defer db.Close()
	}

	if err := checkUninstallEligibility(db, tool, canonical, force); err != nil {
		return err
	}

	if !force && !confirmUninstall(tool) {
		return nil
	}

	return dispatchUninstallExecution(db, tool, canonical, dryRun, purge)
}

func executeStandardUninstall(db *store.DB, tool string, dryRun, purge bool) error {
	manager := resolveUninstallManager(db, tool)
	uninstallCmd := buildUninstallCommand(manager, tool, purge)
	if dryRun {
		fmt.Printf(constants.MsgUninstallDryCmd, strings.Join(uninstallCmd, " "))
		return nil
	}

	fmt.Printf(constants.MsgUninstallRemoving, tool)
	runInstallCommand(uninstallCmd, installOptions{Tool: tool, Verbose: true})

	if db != nil {
		if errRemove := db.RemoveInstalledTool(tool); errRemove != nil {
			fmt.Fprintf(os.Stderr, constants.ErrUninstallDBRemove, tool, errRemove)
		}
	}

	fmt.Printf(constants.MsgUninstallSuccess, tool)
	return nil
}

// confirmUninstall prompts the user for confirmation.
func confirmUninstall(tool string) bool {
	fmt.Printf(constants.MsgUninstallConfirm, tool)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// hasPositionalToolArg reports whether args contain at least one non-flag token.
func hasPositionalToolArg(args []string) bool {
	skipNext := false
	for _, a := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if isFlagToken(a) && (a == "--shell-mode" || a == "-shell-mode") {
			skipNext = true
			continue
		}
		if isFlagToken(a) {
			continue
		}
		return true
	}
	return false
}

// resolveUninstallManager determines which manager was used to install.
func resolveUninstallManager(db *store.DB, tool string) string {
	if db == nil {
		return resolvePackageManager("", tool)
	}
	record, err := db.GetInstalledTool(tool)
	if err != nil || record.PackageManager == "" {
		return resolvePackageManager("", tool)
	}
	return record.PackageManager
}

// buildUninstallCommand builds the uninstall command for a manager.
func buildUninstallCommand(manager, tool string, purge bool) []string {
	pkgName := resolvePackageName(tool, manager)
	switch manager {
	case constants.PkgMgrChocolatey:
		return buildChocoUninstall(pkgName, purge)
	case constants.PkgMgrWinget:
		return []string{"winget", "uninstall", pkgName}
	case constants.PkgMgrApt:
		return buildAptUninstall(pkgName, purge)
	case constants.PkgMgrBrew:
		return []string{"brew", "uninstall", pkgName}
	case constants.PkgMgrSnap:
		return []string{"sudo", "snap", "remove", pkgName}
	default:
		return buildChocoUninstall(pkgName, purge)
	}
}

// buildChocoUninstall builds a Chocolatey uninstall command.
func buildChocoUninstall(pkg string, purge bool) []string {
	args := []string{"choco", "uninstall", pkg, "-y"}
	if purge {
		args = append(args, "-x")
	}
	return args
}

// buildAptUninstall builds an apt uninstall command.
func buildAptUninstall(pkg string, purge bool) []string {
	action := "remove"
	if purge {
		action = "purge"
	}
	return []string{"sudo", "apt", action, "-y", pkg}
}
