package cmdagy

import "github.com/spf13/cobra"

var agyActiveCmd = &cobra.Command{
	Use:     "active",
	Aliases: []string{"running", "act"},
	Short:   "List Antigravity conversations with active running prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleAgyActivePrompts(args)
	},
}

var agyQueuesCmd = &cobra.Command{
	Use:     "queues",
	Aliases: []string{"qs", "queue-status"},
	Short:   "Inspect queued prompts across all workspace queues",
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleAgyQueuesInspector(args)
	},
}

func init() {
	agyActiveCmd.Flags().Bool("json", false, "Output in JSON format")
	agyQueuesCmd.Flags().Bool("json", false, "Output in JSON format")
	AgyCmd.AddCommand(agyActiveCmd)
	AgyCmd.AddCommand(agyQueuesCmd)
}
