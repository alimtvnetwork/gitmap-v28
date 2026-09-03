// Package cmd — agy_pin_projects.go defines Cobra commands for pinned Antigravity projects.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyPinProjectsJSON  bool
	agyPinProjectsRmAll bool
)

var agyPinProjectsCmd = &cobra.Command{
	Use:     "pin-projects",
	Aliases: []string{"pin-project", "pinned-projects", "pinned", "pins"},
	Short:   "Manage pinned Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinProjectsLs()
	},
}

var agyPinProjectsLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all pinned Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinProjectsLs()
	},
}

var agyPinProjectsAddCmd = &cobra.Command{
	Use:     "add [project-id-or-path...]",
	Aliases: []string{"pin", "set"},
	Short:   "Pin one or more Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinProjectsAdd(args)
	},
}

var agyPinProjectsRmCmd = &cobra.Command{
	Use:     "rm [project-id-or-path...]",
	Aliases: []string{"remove", "del", "delete", "unpin"},
	Short:   "Unpin one or more Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinProjectsRm(args)
	},
}

func runAgyPinProjectsLs() error {
	store, loadErr := loadPinnedProjectsStore()

	if loadErr != nil {
		return fmt.Errorf("load pinned projects: %w", loadErr.Unwrap())
	}

	if agyPinProjectsJSON {
		return outputPinnedProjectsJSON(store.Projects)
	}

	renderPinnedProjectsTable(store.Projects)

	return nil
}

func runAgyPinProjectsAdd(args []string) error {
	targets := args

	if len(targets) == 0 {
		targets = []string{"."}
	}

	for _, t := range targets {
		pinned, addErr := addPinnedProjectTarget(t)

		if addErr != nil {
			return fmt.Errorf("pin project %q: %w", t, addErr.Unwrap())
		}

		fmt.Printf("%s Pinned project: %s%s%s (%s)\n",
			constants.ColorGreen+"✓"+constants.ColorReset,
			constants.ColorCyan, pinned.Name, constants.ColorReset,
			pinned.Path)
	}

	return nil
}

func runAgyPinProjectsRm(args []string) error {
	if agyPinProjectsRmAll {
		return handleAgyPinProjectsClearAll()
	}

	if len(args) == 0 {
		return fmt.Errorf("requires project ID, name, path, or --all")
	}

	return handleAgyPinProjectsRemoveTargets(args)
}

func handleAgyPinProjectsClearAll() error {
	count, clearErr := clearAllPinnedProjects()

	if clearErr != nil {
		return fmt.Errorf("clear pinned projects: %w", clearErr.Unwrap())
	}

	fmt.Printf("%s Unpinned all (%d) projects.\n", constants.ColorGreen+"✓"+constants.ColorReset, count)

	return nil
}

func handleAgyPinProjectsRemoveTargets(targets []string) error {
	for _, t := range targets {
		removed, rmErr := removePinnedProjectTarget(t)

		if rmErr != nil {
			return fmt.Errorf("unpin project %q: %w", t, rmErr.Unwrap())
		}

		fmt.Printf("%s Unpinned project: %s%s%s (%s)\n",
			constants.ColorGreen+"✓"+constants.ColorReset,
			constants.ColorCyan, removed.Name, constants.ColorReset,
			removed.Path)
	}

	return nil
}

func initAgyPinProjects() {
	agyPinProjectsLsCmd.Flags().BoolVar(&agyPinProjectsJSON, "json", false, "Output results as JSON")
	agyPinProjectsRmCmd.Flags().BoolVarP(&agyPinProjectsRmAll, "all", "a", false, "Unpin all projects")

	agyPinProjectsCmd.AddCommand(agyPinProjectsLsCmd)
	agyPinProjectsCmd.AddCommand(agyPinProjectsAddCmd)
	agyPinProjectsCmd.AddCommand(agyPinProjectsRmCmd)
}
