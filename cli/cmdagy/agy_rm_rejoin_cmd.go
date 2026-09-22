// Package cmdagy — agy_rm_rejoin_cmd.go defines Cobra commands for rm-rejoin-read and rm-rejoin-pin-read.
package cmdagy

import (
	"github.com/spf13/cobra"
)

var agyRmRejoinReadCmd = &cobra.Command{
	Use:     "rm-rejoin-read [target...]",
	Aliases: []string{"rrr", "rejoin-read"},
	Short:   "Purge Antigravity conversations, rejoin project, and execute read prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "help" {
			renderAgyRmRejoinHelp()
			return nil
		}
		return runAgyRmRejoin(args, false)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeAgyRmArgs(toComplete)
	},
}

var agyRmRejoinPinReadCmd = &cobra.Command{
	Use:     "rm-rejoin-pin-read [target...]",
	Aliases: []string{"rrpr", "rrbr", "rejoin-pin-read"},
	Short:   "Purge Antigravity conversations, rejoin project, pin it, and execute read prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "help" {
			renderAgyRmRejoinHelp()
			return nil
		}
		return runAgyRmRejoin(args, true)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeAgyRmArgs(toComplete)
	},
}

func initAgyRmRejoinCmd() {
	agyRmRejoinReadCmd.AddCommand(agyRmRejoinHelpCmd)
	agyRmRejoinPinReadCmd.AddCommand(agyRmRejoinHelpCmd)

	AgyCmd.AddCommand(agyRmRejoinReadCmd)
	AgyCmd.AddCommand(agyRmRejoinPinReadCmd)
}
