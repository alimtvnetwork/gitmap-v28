package cmdinstall

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

const isDefaultFalse = false

var (
	agyInstallDryRun  bool
	agyInstallYes     bool
	agyInstallVerbose bool
	agyInstallVersion string
	agyInstallForce   bool
	agyInstallPrefix  string
)

var agyInstallCmd = &cobra.Command{
	Use:     "install [cli|ide|manager|all] [flags]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity CLI, Desktop IDE, or Manager",
	Long: `Install Antigravity:
  cli (agy)     - Install Antigravity CLI coding assistant (default)
  ide           - Install Google Antigravity Desktop IDE
  manager (agm) - Install latest release of Antigravity Manager GUI / Tools
  all           - Install Antigravity CLI, Desktop IDE, and Manager`,
	RunE: runAgyInstallCmd,
}

func init() {
	bindAgyInstallFlags()
	cmdagy.AgyCmd.AddCommand(agyInstallCmd)
}

func bindAgyInstallBasicFlags() {
	agyInstallCmd.Flags().BoolVarP(&agyInstallDryRun, "dry-run", "n", isDefaultFalse, "Simulate installation without downloading")
	agyInstallCmd.Flags().BoolVarP(&agyInstallYes, "yes", "y", isDefaultFalse, "Automatic yes to prompts")
	agyInstallCmd.Flags().BoolVarP(&agyInstallVerbose, "verbose", "v", isDefaultFalse, "Enable verbose output")
}

func bindAgyInstallAdvancedFlags() {
	agyInstallCmd.Flags().StringVar(&agyInstallVersion, "version", "", "Specific release version to install (e.g. 2.13.0)")
	agyInstallCmd.Flags().BoolVarP(&agyInstallForce, "force", "f", isDefaultFalse, "Force reinstallation even if already present")
	agyInstallCmd.Flags().StringVar(&agyInstallPrefix, "prefix", "", "Custom installation directory prefix")
}

func bindAgyInstallFlags() {
	bindAgyInstallBasicFlags()
	bindAgyInstallAdvancedFlags()
}

func buildAgyInstallOptions() installOptions {
	return installOptions{
		DryRun:  agyInstallDryRun,
		Yes:     agyInstallYes,
		Verbose: agyInstallVerbose,
		Version: agyInstallVersion,
		Force:   agyInstallForce,
		Prefix:  agyInstallPrefix,
	}
}

func resolveAgyTarget(args []string) string {
	if len(args) > 0 {
		return strings.ToLower(args[0])
	}

	return "cli"
}

func runAgyInstallCmd(cmd *cobra.Command, args []string) error {
	opts := buildAgyInstallOptions()
	target := resolveAgyTarget(args)

	return dispatchAgyInstallTarget(target, opts)
}

func isAgyManagerTarget(target string) bool {
	switch target {
	case "manager", "ag-manager", "agm", "gui", "tools":

		return true
	default:

		return false
	}
}

func isAgyAllTarget(target string) bool {
	switch target {
	case "all", "both":

		return true
	default:

		return false
	}
}

func dispatchAgyInstallTarget(target string, opts installOptions) error {
	if isAgyManagerTarget(target) {
		return runInstallAgManagerWithOpts(opts)
	}

	if isAgyAllTarget(target) {
		return runInstallAgyAll(opts)
	}

	return runInstallAntigravityWithOpts(opts)
}

func runInstallAgyAll(opts installOptions) error {
	fmt.Println("=== [1/2] Installing Antigravity Desktop IDE & CLI ===")

	if err := runInstallAntigravityWithOpts(opts); err != nil {
		return err
	}

	fmt.Println("\n=== [2/2] Installing Antigravity Manager / Tools ===")

	return runInstallAgManagerWithOpts(opts)
}
