package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runUninstall handles the "uninstall" command.
//
// Two modes, dispatched by whether a positional tool name is present:
//
//  1. `gitmap uninstall <tool> [flags]` — third-party tool uninstaller
//     (vscode, npp, …). The original behavior.
//  2. `gitmap uninstall [flags]` (no tool) — shortcut that hands off to
//     `gitmap self-uninstall`. Flags pass through verbatim, so e.g.
//     `gitmap uninstall --confirm --keep-data` works the same as
//     `gitmap self-uninstall --confirm --keep-data`.
func runUninstall(args []string) {
	checkHelp("uninstall", args)

	if !hasPositionalToolArg(args) {
		runSelfUninstall(args)

		return
	}

	fs := flag.NewFlagSet("uninstall", flag.ExitOnError)

	var dryRun, force, purge bool

	fs.BoolVar(&dryRun, constants.FlagUninstallDryRun, false, constants.FlagDescUninstallDryRun)
	fs.BoolVar(&force, constants.FlagUninstallForce, false, constants.FlagDescUninstallForce)
	fs.BoolVar(&purge, constants.FlagUninstallPurge, false, constants.FlagDescUninstallPurge)
	fs.Parse(args)

	tool := fs.Arg(0)
	if tool == "" {
		// Defensive — hasPositionalToolArg already filtered this case.
		runSelfUninstall(args)

		return
	}

	validateToolName(tool)

	if tool == constants.ToolCtx {
		runUninstallCtx()

		return
	}

	db, err := openDB()
	if err != nil {
		if !force {
			fmt.Fprintf(os.Stderr, constants.ErrUninstallNotFound, tool)
			os.Exit(1)
		}
	}

	if db != nil {
		defer db.Close()

		if !db.IsToolInstalled(tool) && !force {
			fmt.Fprintf(os.Stderr, constants.ErrUninstallNotFound, tool)
			os.Exit(1)
		}
	}

	if !force && !confirmUninstall(tool) {
		return
	}

	manager := resolveUninstallManager(db, tool)
	uninstallCmd := buildUninstallCommand(manager, tool, purge)

	if dryRun {
		fmt.Printf(constants.MsgUninstallDryCmd, strings.Join(uninstallCmd, " "))

		return
	}

	fmt.Printf(constants.MsgUninstallRemoving, tool)
	runInstallCommand(uninstallCmd, installOptions{Tool: tool, Verbose: true})

	if db != nil {
		if err := db.RemoveInstalledTool(tool); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrUninstallDBRemove, tool, err)
		}
	}

	fmt.Printf(constants.MsgUninstallSuccess, tool)
}

// confirmUninstall prompts the user for confirmation.
func confirmUninstall(tool string) bool {
	fmt.Printf(constants.MsgUninstallConfirm, tool)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// hasPositionalToolArg reports whether args contain at least one
// non-flag, non-flag-value token. Used to decide between the third-party
// tool uninstaller and the self-uninstall shortcut.
//
// Boolean flags supported by `uninstall` (--dry-run, --force, --purge)
// never consume a value, so the simple "starts with -" check is enough.
// Self-uninstall passthrough flags (--confirm, --keep-data, --keep-snippet,
// --shell-mode <v>) include one value-taking flag — we treat the value
// as a positional only if it does not look like a flag itself.
func hasPositionalToolArg(args []string) bool {
	skipNext := false
	for _, a := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if isFlagToken(a) {
			if a == "--shell-mode" || a == "-shell-mode" {
				skipNext = true
			}
			continue
		}

		return true
	}

	return false
}

// resolveUninstallManager determines which manager was used to install.
func resolveUninstallManager(db *store.DB, tool string) string {
	if db == nil {
		return resolvePackageManager("")
	}

	record, err := db.GetInstalledTool(tool)
	if err != nil || record.PackageManager == "" {
		return resolvePackageManager("")
	}

	return record.PackageManager
}

// buildUninstallCommand builds the uninstall command for a manager.
func buildUninstallCommand(manager, tool string, purge bool) []string {
	pkgName := resolvePackageName(manager, tool)

	if manager == constants.PkgMgrChocolatey {
		return buildChocoUninstall(pkgName, purge)
	}
	if manager == constants.PkgMgrWinget {
		return []string{"winget", "uninstall", pkgName}
	}
	if manager == constants.PkgMgrApt {
		return buildAptUninstall(pkgName, purge)
	}
	if manager == constants.PkgMgrBrew {
		return []string{"brew", "uninstall", pkgName}
	}
	if manager == constants.PkgMgrSnap {
		return []string{"sudo", "snap", "remove", pkgName}
	}

	return buildChocoUninstall(pkgName, purge)
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
