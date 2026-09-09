package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	agyInstallDryRun  bool
	agyInstallYes     bool
	agyInstallVerbose bool
	agyInstallVersion string
)

var agyInstallCmd = &cobra.Command{
	Use:     "install [manager|cli|all] [flags]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity Manager GUI or Antigravity CLI",
	Long: `Install Antigravity tools:
  manager    - Install latest release of Antigravity Manager GUI (default)
  cli (agy)  - Install Antigravity CLI coding assistant
  all        - Install both Antigravity Manager and Antigravity CLI`,
	RunE: runAgyInstallCmd,
}

func init() {
	bindAgyInstallFlags()
}

func bindAgyInstallFlags() {
	agyInstallCmd.Flags().BoolVarP(&agyInstallDryRun, "dry-run", "n", false, "Simulate installation without downloading")
	agyInstallCmd.Flags().BoolVarP(&agyInstallYes, "yes", "y", false, "Automatic yes to prompts")
	agyInstallCmd.Flags().BoolVarP(&agyInstallVerbose, "verbose", "v", false, "Enable verbose output")
	agyInstallCmd.Flags().StringVar(&agyInstallVersion, "version", "", "Specific release version to install (e.g. 4.6.9)")
}

func buildAgyInstallOptions() installOptions {

	return installOptions{
		DryRun:  agyInstallDryRun,
		Yes:     agyInstallYes,
		Verbose: agyInstallVerbose,
		Version: agyInstallVersion,
	}
}

func runAgyInstallCmd(cmd *cobra.Command, args []string) error {
	opts := buildAgyInstallOptions()
	target := "manager"
	if len(args) > 0 {
		target = strings.ToLower(args[0])
	}

	return dispatchAgyInstallTarget(target, opts)
}

func dispatchAgyInstallTarget(target string, opts installOptions) error {
	switch target {
	case "manager", "ag-manager", "gui":

		return runInstallAgManagerWithOpts(opts)
	case "cli", "antigravity", "agy":

		return runInstallAntigravityWithOpts(opts)
	case "all", "both":

		return runInstallAgyAll(opts)
	default:
		fmt.Printf("Unknown target '%s'. Installing Antigravity Manager (default)...\n", target)

		return runInstallAgManagerWithOpts(opts)
	}
}

func runInstallAgyAll(opts installOptions) error {
	fmt.Println("=== [1/2] Installing Antigravity Manager ===")
	if err := runInstallAgManagerWithOpts(opts); err != nil {

		return err
	}
	fmt.Println("\n=== [2/2] Installing Antigravity CLI ===")

	return runInstallAntigravityWithOpts(opts)
}
