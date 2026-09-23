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
	agmUpdateRemote   string
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
	Short:   "Install Antigravity Manager GUI / Tools",
	RunE:    runAgmInstallCmd,
}

var agmUpdateCmd = &cobra.Command{
	Use:     "update [flags]",
	Aliases: []string{"up", "u"},
	Short:   "Update Antigravity Manager GUI / Tools to latest",
	RunE:    runAgmUpdateCmd,
}

func init() {
	bindAgmInstallFlags()
	bindAgmUpdateFlags()
	AgmCmd.AddCommand(agmInstallCmd)
	AgmCmd.AddCommand(agmUpdateCmd)
}

func bindAgmInstallFlags() {
	agmInstallCmd.Flags().BoolVarP(&agmInstallDryRun, "dry-run", "n", false, "Simulate installation without downloading")
	agmInstallCmd.Flags().BoolVarP(&agmInstallYes, "yes", "y", false, "Automatic yes to prompts")
	agmInstallCmd.Flags().BoolVarP(&agmInstallVerbose, "verbose", "v", false, "Enable verbose output")
	agmInstallCmd.Flags().StringVar(&agmInstallVersion, "version", "", "Specific release version to install (e.g. 4.7.1)")
}

func bindAgmUpdateFlags() {
	agmUpdateCmd.Flags().BoolVarP(&agmInstallDryRun, "dry-run", "n", false, "Simulate update without downloading")
	agmUpdateCmd.Flags().BoolVarP(&agmInstallYes, "yes", "y", false, "Automatic yes to prompts")
	agmUpdateCmd.Flags().BoolVarP(&agmInstallVerbose, "verbose", "v", false, "Enable verbose output")
	agmUpdateCmd.Flags().StringVarP(&agmUpdateRemote, "remote", "r", "", "Remote node alias or IP to update")
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

func runAgmUpdateCmd(cmd *cobra.Command, args []string) error {
	if agmUpdateRemote != "" && RemoteAgmUpdateFn != nil {
		return RemoteAgmUpdateFn(agmUpdateRemote)
	}
	opts := buildAgmInstallOptions()

	return runUpdateAgManagerWithOpts(opts)
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
