// Package cmd — agy_group.go handles Antigravity project grouping and batch prompting.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var agyGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"grp", "groups"},
	Short:   "Manage and prompt Antigravity project groups",
}

var agyGroupAddCmd = &cobra.Command{
	Use:   "add [group-name] [project-id-or-path...]",
	Short: "Add project(s) to an Antigravity group",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: gitmap agy group add <group-name> <project-id-or-path...>")
		}

		appErr := addEcosystemGroup("agy", args[0], "", args[1:])
		if appErr != nil {
			return fmt.Errorf("add group: %w", appErr.Unwrap())
		}

		fmt.Printf("%s Added %d target(s) to AGY group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, len(args)-1, args[0])
		return nil
	},
}

var agyGroupLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all Antigravity project groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyGroupList()
	},
}

var agyGroupShowCmd = &cobra.Command{
	Use:   "show [group-name]",
	Short: "Show details of an Antigravity group",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap agy group show <group-name>")
		}

		return runAgyGroupShow(args[0])
	},
}

var agyGroupRmCmd = &cobra.Command{
	Use:     "rm [group-name] [target]",
	Aliases: []string{"remove", "del", "delete"},
	Short:   "Remove a target or an entire group",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap agy group rm <group-name> [target]")
		}

		if len(args) == 1 {
			return runDeleteAgyGroup(args[0])
		}

		if err := removeEcosystemGroupTarget("agy", args[0], args[1]); err != nil {
			return fmt.Errorf("remove target: %w", err.Unwrap())
		}

		fmt.Printf("%s Removed %q from AGY group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, args[1], args[0])
		return nil
	},
}

func runDeleteAgyGroup(name string) error {
	if err := deleteEcosystemGroup("agy", name); err != nil {
		return fmt.Errorf("delete group: %w", err.Unwrap())
	}

	fmt.Printf("%s Deleted AGY group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, name)
	return nil
}

func initAgyGroup() {
	agyGroupCmd.AddCommand(agyGroupAddCmd)
	agyGroupCmd.AddCommand(agyGroupLsCmd)
	agyGroupCmd.AddCommand(agyGroupShowCmd)
	agyGroupCmd.AddCommand(agyGroupRmCmd)
	agyGroupCmd.AddCommand(agyGroupExportCmd)
	agyGroupCmd.AddCommand(agyGroupImportCmd)
	agyGroupCmd.AddCommand(agyGroupPromptCmd)
}
