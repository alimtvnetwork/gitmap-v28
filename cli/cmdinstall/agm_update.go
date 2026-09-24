package cmdinstall

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	agmUpdateRemote string
	agmUpdateSSH    bool
	agmUpdateExcept string
)

var agmUpdateCmd = &cobra.Command{
	Use:     "update [flags] [ssh]",
	Aliases: []string{"up", "u"},
	Short:   "Update Antigravity Manager GUI / Tools to latest",
	RunE:    runAgmUpdateCmd,
}

var agmUpdateAllCmd = &cobra.Command{
	Use:     "update-all [flags] [ssh]",
	Aliases: []string{"update-all-nodes", "updateallnodes", "update-nodes", "updateall", "up-all"},
	Short:   "Update Antigravity Manager across all SSH nodes or local",
	RunE:    runAgmUpdateAllCmd,
}

func init() {
	bindAgmUpdateFlags()
	AgmCmd.AddCommand(agmUpdateCmd)
	AgmCmd.AddCommand(agmUpdateAllCmd)
}

func bindAgmUpdateFlags() {
	bindCommonUpdateFlags(agmUpdateCmd)
	bindCommonUpdateFlags(agmUpdateAllCmd)
}

func bindCommonUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&agmInstallDryRun, "dry-run", "n", false, "Simulate update without downloading")
	cmd.Flags().BoolVarP(&agmInstallYes, "yes", "y", false, "Automatic yes to prompts")
	cmd.Flags().BoolVarP(&agmInstallVerbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().StringVar(&agmInstallVersion, "version", "", "Specific release version to update to (e.g. 4.7.1)")
	cmd.Flags().StringVarP(&agmUpdateRemote, "remote", "r", "", "Remote node alias or IP to update")
	cmd.Flags().BoolVarP(&agmUpdateSSH, "ssh", "s", false, "Update remote fleet machines via SSH")
	cmd.Flags().StringVarP(&agmUpdateExcept, "except", "e", "", "Exclude machines by ID, alias, or IP")
	cmd.Flags().StringVar(&agmUpdateExcept, "excep", "", "Exclude machines by ID, alias, or IP")
}

func runAgmUpdateCmd(cmd *cobra.Command, args []string) error {
	hasSSH := agmUpdateSSH || hasSSHArg(args)
	hasAll := hasAllArg(args)
	if (hasSSH || agmUpdateRemote != "" || hasAll) && RemoteAgmUpdateFleetFn != nil {
		target := resolveAgmUpdateTarget(agmUpdateRemote, hasAll)
		return RemoteAgmUpdateFleetFn(target, agmUpdateExcept)
	}
	opts := buildAgmInstallOptions()
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && a != "ssh" && a != "all" && opts.Version == "" {
			opts.Version = strings.TrimPrefix(a, "v")
		}
	}
	return runUpdateAgManagerWithOpts(opts)
}

func resolveAgmUpdateTarget(target string, hasAll bool) string {
	if target == "" || hasAll {
		return "all-nodes"
	}

	return target
}

func runAgmUpdateAllCmd(cmd *cobra.Command, args []string) error {
	if RemoteAgmUpdateFleetFn != nil {
		return RemoteAgmUpdateFleetFn("all-nodes", agmUpdateExcept)
	}
	opts := buildAgmInstallOptions()
	return runUpdateAgManagerWithOpts(opts)
}

func hasSSHArg(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "ssh" || low == "--ssh" || low == "-s" {
			return true
		}
	}
	return false
}

func hasAllArg(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "all" || low == "all-nodes" || low == "update-all" {
			return true
		}
	}
	return false
}
