// Package cmdagy — agy_running_prompts_cmd.go routes and dispatches running-prompts CLI commands and flags.
package cmdagy

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunningPromptsCmd is the root command for running-prompts operations.
var RunningPromptsCmd = &cobra.Command{
	Use:     "running-prompts",
	Aliases: []string{"running-prompt", "rp-prompts"},
	Short:   "Manage running and queued Antigravity prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		RenderRunningPromptsHelp()
		return nil
	},
}

// BackupRunningPromptsCmd provides direct access to backup-running-prompts.
var BackupRunningPromptsCmd = &cobra.Command{
	Use:     "backup-running-prompts",
	Aliases: []string{"backup-running-prompt", "brp"},
	Short:   "Snapshot running and queued prompts to Split-DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		isJSON, _ := cmd.Flags().GetBool("json")
		isSSH, _ := cmd.Flags().GetBool("ssh")
		return RunRunningPromptsBackup(file, isJSON, isSSH)
	},
}

// RestoreRunningPromptsCmd provides direct access to restore-running-prompts.
var RestoreRunningPromptsCmd = &cobra.Command{
	Use:     "restore-running-prompts",
	Aliases: []string{"restore-running-prompt", "rrp"},
	Short:   "Restore backed-up prompts into project queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := extractRestoreOptions(cmd)
		return RunRunningPromptsRestore(opts)
	},
}

func extractRestoreOptions(cmd *cobra.Command) store.RestoreOptions {
	isKeep, _ := cmd.Flags().GetBool("keep")
	isJSON, _ := cmd.Flags().GetBool("json")
	isSSH, _ := cmd.Flags().GetBool("ssh")
	file, _ := cmd.Flags().GetString("file")
	return store.RestoreOptions{
		IsKeep:     isKeep,
		IsJSON:     isJSON,
		IsSSH:      isSSH,
		TargetFile: file,
	}
}

func initAgyRunningPromptsCommands() {
	setupBackupRunningPromptsCmd()
	setupRestoreRunningPromptsCmd()
	setupRunningPromptsSubcommands()
	AgyCmd.AddCommand(RunningPromptsCmd)
	AgyCmd.AddCommand(BackupRunningPromptsCmd)
	AgyCmd.AddCommand(RestoreRunningPromptsCmd)
}

func setupBackupRunningPromptsCmd() {
	BackupRunningPromptsCmd.Flags().StringP("file", "f", "", "Target backup database path")
	BackupRunningPromptsCmd.Flags().Bool("json", false, "Output summary in JSON format")
	BackupRunningPromptsCmd.Flags().Bool("ssh", false, "Snapshot prompts across cluster SSH fleet")
	_ = cobra.MarkFlagFilename(BackupRunningPromptsCmd.Flags(), "file")
	lsCmd := makeBackupLsCmd()
	helpCmd := makeHelpCmd(RenderBackupHelp)
	BackupRunningPromptsCmd.AddCommand(lsCmd, helpCmd)
}

func setupRestoreRunningPromptsCmd() {
	RestoreRunningPromptsCmd.Flags().BoolP("keep", "k", false, "Exclude from retention auto-pruning")
	RestoreRunningPromptsCmd.Flags().Bool("json", false, "Output restoration in JSON format")
	RestoreRunningPromptsCmd.Flags().Bool("ssh", false, "Restore prompts across cluster SSH fleet")
	RestoreRunningPromptsCmd.Flags().StringP("file", "f", "", "Source database path")
	_ = cobra.MarkFlagFilename(RestoreRunningPromptsCmd.Flags(), "file")
	helpCmd := makeHelpCmd(RenderRestoreHelp)
	RestoreRunningPromptsCmd.AddCommand(helpCmd)
}

func setupRunningPromptsSubcommands() {
	backupCmd := makeRunningPromptsBackupCmd()
	restoreCmd := makeRunningPromptsRestoreCmd()
	cleanCmd := makeRunningPromptsCleanCmd()
	lsCmd := makeRunningPromptsLsCmd()
	exportCmd := makeRunningPromptsExportCmd()
	importCmd := makeRunningPromptsImportCmd()
	helpCmd := makeHelpCmd(RenderRunningPromptsHelp)
	RunningPromptsCmd.AddCommand(backupCmd, restoreCmd, cleanCmd, lsCmd, exportCmd, importCmd, helpCmd)
	RunningPromptsCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{
				"ls\tInspect and list running and queued prompts",
				"backup\tSnapshot running and queued prompts to Split-DB",
				"restore\tRestore backed-up prompts into project queue",
				"clean\tPrune expired restored entries or force cleanup",
				"export\tExport prompts snapshot to file",
				"import\tImport prompts from file into queue",
				"help\tShow command help and synopsis",
			}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

func makeHelpCmd(renderFn func()) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Show command help and synopsis",
		Run: func(cmd *cobra.Command, args []string) {
			renderFn()
		},
	}
}

func makeRunningPromptsBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Snapshot running and queued prompts to Split-DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			isJSON, _ := cmd.Flags().GetBool("json")
			isSSH, _ := cmd.Flags().GetBool("ssh")
			return RunRunningPromptsBackup(file, isJSON, isSSH)
		},
	}
	cmd.Flags().StringP("file", "f", "", "Target backup database path")
	cmd.Flags().Bool("json", false, "Output summary in JSON format")
	cmd.Flags().Bool("ssh", false, "Snapshot prompts across cluster SSH fleet")
	cmd.AddCommand(makeBackupLsCmd(), makeHelpCmd(RenderBackupHelp))
	return cmd
}

func makeBackupLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List backup batches recorded in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			isJSON, _ := cmd.Flags().GetBool("json")
			isSSH, _ := cmd.Flags().GetBool("ssh")
			return RunRunningPromptsBackupLs(file, isJSON, isSSH)
		},
	}
	cmd.Flags().StringP("file", "f", "", "Target backup database path")
	cmd.Flags().Bool("json", false, "Output in JSON format")
	cmd.Flags().Bool("ssh", false, "List backup batches across cluster SSH fleet")
	return cmd
}

func makeRunningPromptsRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore backed-up prompts into project queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := extractRestoreOptions(cmd)
			return RunRunningPromptsRestore(opts)
		},
	}
	cmd.Flags().BoolP("keep", "k", false, "Exclude from retention auto-pruning")
	cmd.Flags().Bool("json", false, "Output restoration in JSON format")
	cmd.Flags().Bool("ssh", false, "Restore prompts across cluster SSH fleet")
	cmd.Flags().StringP("file", "f", "", "Source database path")
	cmd.AddCommand(makeHelpCmd(RenderRestoreHelp))
	return cmd
}

func makeRunningPromptsCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Prune expired restored entries or force cleanup",
		RunE: func(cmd *cobra.Command, args []string) error {
			isForce, _ := cmd.Flags().GetBool("force")
			file, _ := cmd.Flags().GetString("file")
			return RunRunningPromptsClean(file, isForce)
		},
	}
	cmd.Flags().Bool("force", false, "Force delete all backup batches")
	cmd.Flags().StringP("file", "f", "", "Target backup database path")
	return cmd
}

func makeRunningPromptsLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Inspect and list running and queued prompts",
		RunE:  executeRunningPromptsLsCmd,
	}
	setupRunningPromptsLsFlags(cmd)
	return cmd
}

func executeRunningPromptsLsCmd(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")
	wc := resolveWordCountFlag(cmd)
	isFull, _ := cmd.Flags().GetBool("full")
	isJSON, _ := cmd.Flags().GetBool("json")
	isSSH, _ := cmd.Flags().GetBool("ssh")
	return RunRunningPromptsLs(limit, wc, isFull, isJSON, isSSH)
}

func setupRunningPromptsLsFlags(cmd *cobra.Command) {
	cmd.Flags().IntP("limit", "l", 8, "Maximum prompts to display")
	cmd.Flags().Int("wordcount", 100, "Maximum word count for prompt preview")
	cmd.Flags().Int("wc", 100, "Alias for wordcount")
	cmd.Flags().Bool("full", false, "Display full prompt text without truncation")
	cmd.Flags().Bool("json", false, "Output in JSON format")
	cmd.Flags().Bool("ssh", false, "Inspect prompts across cluster SSH fleet")
}

func resolveWordCountFlag(cmd *cobra.Command) int {
	if cmd.Flags().Changed("wc") {
		wc, _ := cmd.Flags().GetInt("wc")
		return wc
	}
	wc, _ := cmd.Flags().GetInt("wordcount")
	return wc
}

func makeRunningPromptsExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export running prompts to SQLite .db or .json",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			wc, _ := cmd.Flags().GetInt("wc")
			return RunRunningPromptsExport(file, wc)
		},
	}
	cmd.Flags().StringP("file", "f", "gitmap-running-prompts.db", "Export target file path (.db or .json)")
	cmd.Flags().Int("wc", 0, "Word count truncation limit (0 for full)")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "file")
	return cmd
}

func makeRunningPromptsImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import running prompts from SQLite .db or .json into queues",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			wc, _ := cmd.Flags().GetInt("wc")
			return RunRunningPromptsImport(file, wc)
		},
	}
	cmd.Flags().StringP("file", "f", "gitmap-running-prompts.db", "Import source file path (.db or .json)")
	cmd.Flags().Int("wc", 0, "Word count truncation limit (0 for full)")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "file")
	return cmd
}
