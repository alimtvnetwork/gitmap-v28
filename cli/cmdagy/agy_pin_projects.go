// Package cmdagy — agy_pin_projects.go defines Cobra commands for pinned Antigravity projects.
package cmdagy

import (
	"github.com/spf13/cobra"
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
		if len(args) == 0 {
			return runAgyPinProjectsLs()
		}
		if args[0] == "help" {
			renderAgyPinsHelp()
			return nil
		}
		if args[0] == "ls" || args[0] == "list" {
			return runAgyPinProjectsLs()
		}
		return runAgyPinProjectsAdd(args)
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
	Use:     "add [target...]",
	Aliases: []string{"pin", "set"},
	Short:   "Pin one or more Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinProjectsAdd(args)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeAgyRmArgs(toComplete)
	},
}

var agyPinProjectsRmCmd = &cobra.Command{
	Use:     "rm [target...]",
	Aliases: []string{"remove", "del", "delete", "unpin"},
	Short:   "Unpin one or more Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinProjectsRm(args)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeAgyRmArgs(toComplete)
	},
}

func initAgyPinProjects() {
	agyPinProjectsLsCmd.Flags().BoolVar(&agyPinProjectsJSON, "json", false, "Output results as JSON")
	agyPinProjectsRmCmd.Flags().BoolVarP(&agyPinProjectsRmAll, "all", "a", false, "Unpin all projects")

	agyPinProjectsCmd.AddCommand(agyPinProjectsLsCmd)
	agyPinProjectsCmd.AddCommand(agyPinProjectsAddCmd)
	agyPinProjectsCmd.AddCommand(agyPinProjectsRmCmd)
	agyPinProjectsCmd.AddCommand(agyPinProjectsEditCmd)
	agyPinProjectsCmd.AddCommand(agyPinProjectsHelpCmd)
}
