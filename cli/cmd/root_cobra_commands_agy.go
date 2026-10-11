package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/spf13/cobra"
)

func makeTopLevelRerunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rerun [target] [flags]",
		Aliases: []string{"rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue"},
		Short:   "Rerun active prompt(s) with IDE restart, image re-injection, and queued prefix check",
		RunE: func(c *cobra.Command, args []string) error {
			return cmdagy.RunRerunTopLevelCLI(args)
		},
	}
	cmd.Flags().AddFlagSet(cmdagy.GetAgyRerunCmd().Flags())
	cmd.ValidArgsFunction = cmdagy.GetAgyRerunCmd().ValidArgsFunction
	return cmd
}

func makeTopLevelSugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shutdown-until-green [command]",
		Aliases: []string{"sug", "shutdown-until"},
		Short:   "Monitor designated projects and trigger OS shutdown when all CI/CD pipelines turn green",
		RunE: func(c *cobra.Command, args []string) error {
			return cmdagy.RunSUGCLI(args)
		},
	}
	cmd.Flags().AddFlagSet(cmdagy.AgySUGCmd.Flags())
	cmd.ValidArgsFunction = cmdagy.AgySUGCmd.ValidArgsFunction
	return cmd
}

func makeTopLevelWPRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "watch-prompts-running [command]",
		Aliases: []string{"wpr", "watch-running-prompts"},
		Short:   "Watch, backup, auto-recover, and deploy running prompts",
		RunE: func(c *cobra.Command, args []string) error {
			return cmdagy.RunWPRCLI(args)
		},
	}
	cmd.Flags().AddFlagSet(cmdagy.AgyWPRCmd.Flags())
	cmd.ValidArgsFunction = cmdagy.AgyWPRCmd.ValidArgsFunction
	return cmd
}

func makeTopLevelRunningPromptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "running-prompts [command]",
		Aliases:           []string{"running-prompt", "rp-prompts", "backup-running-prompts", "restore-running-prompts"},
		Short:             "Manage running and queued Antigravity prompts",
		Run:               func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: cmdagy.RunningPromptsCmd.ValidArgsFunction,
	}
}

func makeTopLevelRunningProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "running-projects [ls]",
		Aliases:           []string{"runningprojects", "rp"},
		Short:             "List projects hosting active or queued Antigravity prompts",
		Run:               func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: cmdagy.AgyRunningProjectsCmd.ValidArgsFunction,
	}
	cmd.Flags().AddFlagSet(cmdagy.AgyRunningProjectsCmd.Flags())
	return cmd
}

func makeTopLevelFPUGCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "finish-prompts-until-green [targets...]",
		Aliases:           []string{"fpug"},
		Short:             "Monitor target projects until prompts finish and CI/CD pipelines turn green",
		Run:               func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: cmdagy.AgyFPUGCmd.ValidArgsFunction,
	}
	cmd.Flags().AddFlagSet(cmdagy.AgyFPUGCmd.Flags())
	return cmd
}
