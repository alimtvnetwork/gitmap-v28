// Package cmdagy — agy_account_switch_cmd.go defines CLI commands and flags for Antigravity account switching.
package cmdagy

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// AgyAccountSwitchCmd is the root command for Antigravity account switching.
var AgyAccountSwitchCmd = &cobra.Command{
	Use:     "account-switch",
	Aliases: []string{"asw", "switch-account", "fast-forward", "ff"},
	Short:   "Automated Antigravity account switch with parallel prompt backup, lock check, and fast-forward",
	RunE:    dispatchRootAccountSwitch,
}

func initAgyAccountSwitchCommands() {
	bindAccountSwitchCommonFlags(AgyAccountSwitchCmd, DefaultAccountSwitchThreshold)
	AgyAccountSwitchCmd.AddCommand(
		makeAccountSwitchRunCmd(),
		makeAccountSwitchTestCmd(),
		makeAccountSwitchStatusCmd(),
		makeAccountSwitchSetThresholdCmd(),
		makeHelpCmd(RenderAccountSwitchHelp),
	)
	AgyAccountSwitchCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		RenderAccountSwitchHelp()
	})
	AgyCmd.AddCommand(AgyAccountSwitchCmd)
}

func bindAccountSwitchCommonFlags(cmd *cobra.Command, defaultThreshold int) {
	cmd.Flags().IntP("threshold", "t", defaultThreshold, "Credit percentage threshold to trigger switch")
	cmd.Flags().StringP("instance", "i", "", "Isolated temporary instance name for testing")
	cmd.Flags().Bool("test-e2e", false, "Run end-to-end 98% threshold switch test and reset to 15%")
	cmd.Flags().Bool("dry-run", false, "Preview candidate ranking and lock check without switching")
	cmd.Flags().Bool("json", false, "Output telemetry in JSON format")
}

func dispatchRootAccountSwitch(cmd *cobra.Command, args []string) error {
	opts := extractAccountSwitchOptions(cmd)
	if opts.IsTestE2E {
		return ExecuteAccountSwitchInstanceE2E(opts.ThresholdPct, opts.InstanceName, opts.IsJSON)
	}
	if cmd.Flags().Changed("threshold") || cmd.Flags().Changed("dry-run") || cmd.Flags().Changed("instance") {
		return RunAccountSwitchWorkflow(opts)
	}
	RenderAccountSwitchHelp()
	return nil
}

func extractAccountSwitchOptions(cmd *cobra.Command) AccountSwitchOptions {
	threshold, _ := cmd.Flags().GetInt("threshold")
	instance, _ := cmd.Flags().GetString("instance")
	isTestE2E, _ := cmd.Flags().GetBool("test-e2e")
	isDryRun, _ := cmd.Flags().GetBool("dry-run")
	isJSON, _ := cmd.Flags().GetBool("json")
	return AccountSwitchOptions{
		ThresholdPct: threshold, InstanceName: instance,
		IsTestE2E: isTestE2E, IsDryRun: isDryRun, IsJSON: isJSON,
	}
}

func makeAccountSwitchRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Evaluate credit against threshold and execute account switch lifecycle",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := extractAccountSwitchOptions(cmd)
			return RunAccountSwitchWorkflow(opts)
		},
	}
	bindAccountSwitchCommonFlags(cmd, DefaultAccountSwitchThreshold)
	return cmd
}

func makeAccountSwitchTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run E2E account switch test at 98% threshold in temporary instance and reset to 15%",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := extractAccountSwitchOptions(cmd)
			opts.IsTestE2E = true
			return ExecuteAccountSwitchInstanceE2E(opts.ThresholdPct, opts.InstanceName, opts.IsJSON)
		},
	}
	bindAccountSwitchCommonFlags(cmd, 98)
	return cmd
}

func makeAccountSwitchStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current account switch threshold, active profile, and last backup status",
		RunE: func(cmd *cobra.Command, args []string) error {
			isJSON, _ := cmd.Flags().GetBool("json")
			return RunAccountSwitchStatus(isJSON)
		},
	}
	cmd.Flags().Bool("json", false, "Output status in JSON format")
	return cmd
}

func makeAccountSwitchSetThresholdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-threshold <pct>",
		Short: "Update and persist the credit threshold percentage (default 15)",
		RunE:  executeSetThresholdCmd,
	}
	cmd.Flags().Bool("json", false, "Output updated state in JSON format")
	return cmd
}

func executeSetThresholdCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("missing threshold percentage argument (e.g. 15 or 98)", "E400")
	}
	pct, err := strconv.Atoi(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "invalid threshold integer")
	}
	isJSON, _ := cmd.Flags().GetBool("json")
	return RunSetAccountSwitchThreshold(pct, isJSON)
}

// RenderAccountSwitchHelp renders the boxed help menu with examples for account-switch.
func RenderAccountSwitchHelp() {
	termout.RenderMenu(buildAccountSwitchHelpMenu())
}

func buildAccountSwitchHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Antigravity Account Switch & Threshold Governance (account-switch / asw)",
		UsageLines: []string{
			"gitmap agy account-switch <subcommand> [flags]",
			"gitmap agy asw test --threshold 98 --instance e2e-test-01",
			"gitmap agy account-switch run --threshold 15",
		},
		Sections: []termout.HelpSection{
			buildAccountSwitchSubcommandsSection(),
			buildAccountSwitchExamplesSection(),
		},
		FooterFlags: buildAccountSwitchFooterFlags(),
		Tips: []string{
			"Before switching, parallel 1-to-1 SQLite running-prompts backup captures prompts & screenshots.",
			"Candidates are ranked by highest credit, verified via Refresh, and checked against Email & Supabase locks.",
			"After E2E testing at 98%, the default threshold automatically resets to 15%.",
		},
	}
}

func buildAccountSwitchSubcommandsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Account Switch Subcommands",
		Entries: []termout.CommandEntry{
			{Command: "run", Description: "Evaluate active credit and switch account when below threshold"},
			{Command: "test", Description: "Run full 98% E2E switch test in temporary instance & reset to 15%"},
			{Command: "status", Description: "Inspect active account, threshold (15%), and last backup batch"},
			{Command: "set-threshold <pct>", Description: "Set default credit threshold percentage (e.g. 15 or 98)"},
			{Command: "help", Description: "Display this boxed help guide with usage examples"},
		},
	}
}

func buildAccountSwitchExamplesSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Examples",
		Entries: []termout.CommandEntry{
			{Command: "gitmap agy account-switch test --threshold 98", Description: "Run E2E 98% switch + backup/restore + 15% reset"},
			{Command: "gitmap agy asw run --threshold 15 --json", Description: "Run production account switch check with JSON output"},
			{Command: "gitmap agy account-switch set-threshold 15", Description: "Persist production threshold at 15%"},
			{Command: "gitmap agy account-switch status", Description: "Display current account switch telemetry"},
		},
	}
}

func buildAccountSwitchFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-t, --threshold <pct>", Description: "Credit threshold percentage (default 15, E2E test 98)"},
		{Command: "-i, --instance <name>", Description: "Temporary instance workspace name for isolated E2E testing"},
		{Command: "--test-e2e", Description: "Trigger full E2E instance test lifecycle"},
		{Command: "--dry-run", Description: "Verify candidate refresh & locks without executing switch"},
		{Command: "--json", Description: "Emit structured JSON output"},
		{Command: "-h, --help", Description: "Display help menu"},
	}
}
