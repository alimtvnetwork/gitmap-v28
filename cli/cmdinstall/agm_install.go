package cmdinstall

import (
	"context"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
	"github.com/spf13/cobra"
)

var (
	agmInstallDryRun  bool
	agmInstallYes     bool
	agmInstallVerbose bool
	agmInstallVersion string
)

// AgmCmd manages Antigravity Manager GUI and tools.
var AgmCmd = &cobra.Command{
	Use:   "agm",
	Short: "Antigravity Manager GUI and Tools management",
	Long:  "Manage and install Antigravity Manager GUI desktop application and tools.",
	Run: func(cmd *cobra.Command, args []string) {
		printAgmHelp()
	},
}

var agmInstallCmd = &cobra.Command{
	Use:     "install [flags]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity Manager GUI / Tools",
	RunE:    runAgmInstallCmd,
}

// RunAGMVersionTagsLSFn delegates AGM version listing to cmd package.
var RunAGMVersionTagsLSFn func() error

var agmVersionCmd = &cobra.Command{
	Use:     "version [ls|list]",
	Aliases: []string{"versions", "tags", "ls"},
	Short:   "List all available Antigravity Manager release tags from GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		if RunAGMVersionTagsLSFn != nil {
			return RunAGMVersionTagsLSFn()
		}
		return nil
	},
}

func init() {
	bindAgmInstallFlags()
	AgmCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printAgmHelp()
	})
	AgmCmd.AddCommand(agmInstallCmd)
	AgmCmd.AddCommand(agmVersionCmd)
	initAgmUninstallCmd()
}

func printAgmHelp() {
	termhelp.RenderMenu(buildAgmHelpMenu())
}

func buildAgmHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Antigravity Manager (AGM) Installation & Fleet Update Suite",
		UsageLines: []string{
			"gitmap agm <command> [flags]",
			"gitmap agm update [--version <tag>] [--ssh]",
		},
		Sections: []termhelp.HelpSection{
			buildAgmCommandsSection(),
			buildAgmExamplesSection(),
		},
		FooterFlags: buildAgmFooterFlags(),
		Tips: []string{
			"Run 'gitmap agm update' to upgrade to the latest release with cache-busted install.ps1.",
			"Pin a specific release using 'gitmap agm update --version 4.89.0'.",
			"Combine with 'gitmap agy running-prompts backup' and 'gitmap agy account-switch test --threshold 98'.",
		},
	}
}

func buildAgmCommandsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Core Commands",
		Entries: []termhelp.CommandEntry{
			{Command: "install", Description: "Install Antigravity Manager (agm-alim) desktop application"},
			{Command: "update", Description: "Update Antigravity Manager to the latest or pinned GitHub release"},
			{Command: "update-all", Description: "Update Antigravity Manager across all SSH fleet nodes"},
			{Command: "version ls", Description: "List available Antigravity Manager release tags from GitHub"},
			{Command: "uninstall", Description: "Uninstall Antigravity Manager (use --all for full purge)"},
			{Command: "uninstall-all", Description: "Full purge of Antigravity Manager configuration and shortcuts"},
		},
	}
}

func buildAgmExamplesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Real-World Examples",
		Entries: []termhelp.CommandEntry{
			{Command: "gitmap agm update", Description: "Update local Antigravity Manager to the latest GitHub release"},
			{Command: "gitmap agm update --version 4.89.0", Description: "Install or update to pinned release v4.89.0"},
			{Command: "gitmap agm update --ssh", Description: "Update Antigravity Manager across all remote SSH nodes"},
			{Command: "gitmap agy running-prompts backup", Description: "Snapshot active & queued prompts before AGM maintenance"},
			{Command: "gitmap agy running-prompts restore", Description: "Restore backed-up prompts after AGM account switch"},
			{Command: "gitmap agy account-switch test --threshold 98", Description: "Run live E2E account-switch verification at 98% threshold"},
		},
	}
}

func buildAgmFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--version <tag>", Description: "Specific release version to install or update (e.g. 4.89.0)"},
		{Command: "-s, --ssh", Description: "Execute update across remote SSH cluster fleet"},
		{Command: "-r, --remote <node>", Description: "Target a specific remote node alias or IP"},
		{Command: "-n, --dry-run", Description: "Simulate installation or update without downloading"},
		{Command: "-y, --yes", Description: "Automatic yes to prompts"},
		{Command: "-h, --help", Description: "Display this Antigravity Manager help menu"},
	}
}

func bindAgmInstallFlags() {
	agmInstallCmd.Flags().BoolVarP(&agmInstallDryRun, "dry-run", "n", false, "Simulate installation without downloading")
	agmInstallCmd.Flags().BoolVarP(&agmInstallYes, "yes", "y", false, "Automatic yes to prompts")
	agmInstallCmd.Flags().BoolVarP(&agmInstallVerbose, "verbose", "v", false, "Enable verbose output")
	agmInstallCmd.Flags().StringVar(&agmInstallVersion, "version", "", "Specific release version to install (e.g. 4.89.0)")
}

func buildAgmInstallOptions() installOptions {
	return installOptions{
		DryRun:  agmInstallDryRun,
		Yes:     agmInstallYes,
		Verbose: agmInstallVerbose,
		Version: agmInstallVersion,
	}
}

func runAgmInstallCmd(cmd *cobra.Command, args []string) error {
	opts := buildAgmInstallOptions()
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && opts.Version == "" {
			opts.Version = strings.TrimPrefix(a, "v")
		}
	}

	return runInstallAgManagerWithOpts(opts)
}

// DispatchAgm routes CLI arguments to agm commands.
func DispatchAgm(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) > 0 && isAgmCommand(args[0]) {
		args = args[1:]
	}
	if len(args) == 0 || isAgmHelpArg(args[0]) {
		printAgmHelp()
		return nil
	}
	AgmCmd.SetArgs(args)
	return AgmCmd.ExecuteContext(ctx)
}

func isAgmHelpArg(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))
	return low == "help" || low == "-h" || low == "--help"
}

func isAgmCommand(arg string) bool {
	low := strings.ToLower(arg)

	return low == "agm" || low == "ag-manager" || low == "antigravity-manager"
}
