package cmdinstall

import (
	"context"
	"strings"

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
}

var agmInstallCmd = &cobra.Command{
	Use:     "install [flags]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity Manager GUI / Tools (latest GitHub release)",
	RunE:    runAgmInstallCmd,
}

func init() {
	bindAgmInstallFlags()
	AgmCmd.AddCommand(agmInstallCmd)
}

func bindAgmInstallFlags() {
	agmInstallCmd.Flags().BoolVarP(&agmInstallDryRun, "dry-run", "n", false, "Simulate installation without downloading")
	agmInstallCmd.Flags().BoolVarP(&agmInstallYes, "yes", "y", false, "Automatic yes to prompts")
	agmInstallCmd.Flags().BoolVarP(&agmInstallVerbose, "verbose", "v", false, "Enable verbose output")
	agmInstallCmd.Flags().StringVar(&agmInstallVersion, "version", "", "Specific release version to install (e.g. 4.7.1)")
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

	return runInstallAgManagerWithOpts(opts)
}

// DispatchAgm routes CLI arguments to agm commands.
func DispatchAgm(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) > 0 && isAgmCommand(args[0]) {
		args = args[1:]
	}

	AgmCmd.SetArgs(args)

	return AgmCmd.ExecuteContext(ctx)
}

func isAgmCommand(arg string) bool {
	low := strings.ToLower(arg)

	return low == "agm" || low == "ag-manager" || low == "antigravity-manager"
}
