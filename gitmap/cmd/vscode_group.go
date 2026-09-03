// Package cmd — vscode_group.go implements group operations for VSCode and GitHub Desktop.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runVSCodeGroup(args []string) error {
	return dispatchEcosystemGroupCli("vscode", "VS Code", "targets", args)
}

func runGitHubDesktopGroup(args []string) error {
	return dispatchEcosystemGroupCli("githubdesktop", "GitHub Desktop", "repos", args)
}

func dispatchEcosystemGroupCli(ecosystem string, label string, unit string, args []string) error {
	if len(args) == 0 || args[0] == "ls" || args[0] == "list" {
		return runEcosystemGroupList(ecosystem, label)
	}

	if args[0] == "add" && len(args) >= 3 {
		return handleEcosystemGroupAdd(ecosystem, label, unit, args[1], args[2:])
	}

	if (args[0] == "rm" || args[0] == "delete") && len(args) >= 2 {
		return handleEcosystemGroupRm(ecosystem, label, args[1:])
	}

	return fmt.Errorf("usage: gitmap %s group [ls|add <group> <%s...>|rm <group> [%s]]", ecosystem, unit, unit)
}

func handleEcosystemGroupAdd(ecosystem string, label string, unit string, group string, targets []string) error {
	appErr := addEcosystemGroup(ecosystem, group, "", targets)
	if appErr != nil {
		return fmt.Errorf("add %s group: %w", label, appErr.Unwrap())
	}

	fmt.Printf("%s Added %d %s to %s group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, len(targets), unit, label, group)
	return nil
}

func handleEcosystemGroupRm(ecosystem string, label string, args []string) error {
	if len(args) == 1 {
		return handleEcosystemGroupDelete(ecosystem, label, args[0])
	}

	if err := removeEcosystemGroupTarget(ecosystem, args[0], args[1]); err != nil {
		return fmt.Errorf("remove %s target: %w", label, err.Unwrap())
	}

	fmt.Printf("%s Removed %q from %s group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, args[1], label, args[0])
	return nil
}

func handleEcosystemGroupDelete(ecosystem string, label string, group string) error {
	if err := deleteEcosystemGroup(ecosystem, group); err != nil {
		return fmt.Errorf("delete %s group: %w", label, err.Unwrap())
	}

	fmt.Printf("%s Deleted %s group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, label, group)
	return nil
}
