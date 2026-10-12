package cmdinstall

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const isDefaultFalse = false

var (
	agyInstallDryRun       bool
	agyInstallYes          bool
	agyInstallVerbose      bool
	agyInstallVersion      string
	agyInstallForce        bool
	agyInstallPrefix       string
	agyInstallDownloadMust bool
)

var agyInstallCmd = &cobra.Command{
	Use:     "install [cli|ide|manager|all] [flags]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity CLI, Desktop IDE, or Manager",
	Long: `Install Antigravity:
  cli (agy)     - Install Antigravity CLI coding assistant (default)
  ide           - Install Google Antigravity Desktop IDE
  manager (agm) - Install latest release of Antigravity Manager GUI / Tools
  all           - Install Antigravity CLI, Desktop IDE, and Manager

Examples:
  gitmap agy install ide
  gitmap agy install ide --force
  gitmap agy update ide`,
	RunE: runAgyInstallCmd,
}

var agyUpdateCmd = &cobra.Command{
	Use:     "update [cli|ide|manager|all] [flags]",
	Aliases: []string{"up"},
	Short:   "Update Antigravity CLI, Desktop IDE, or Manager",
	Long: `Update Antigravity components:
  cli (agy)     - Update Antigravity CLI coding assistant (default)
  ide           - Update Google Antigravity Desktop IDE
  manager (agm) - Update latest release of Antigravity Manager GUI / Tools
  all           - Update Antigravity CLI, Desktop IDE, and Manager

Examples:
  gitmap agy update ide
  gitmap agy update manager
  gitmap agy update all`,
	RunE: runAgyUpdateCmd,
}

func init() {
	bindAgyInstallFlags()
	bindAgyUpdateFlags()
}

// RegisterAgyInstallSubcommand attaches the agy install and update subcommands to the given
// parent command. Called by the owning agy package during its init so that
// cmdinstall does not import cmdagy (import cycle).
func RegisterAgyInstallSubcommand(parent *cobra.Command) {
	parent.AddCommand(agyInstallCmd)
	parent.AddCommand(agyUpdateCmd)
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
	bindAgyInstallDownloadFlags()
}

func bindAgyInstallDownloadFlags() {
	agyInstallCmd.Flags().BoolVar(&agyInstallDownloadMust, "download-must", isDefaultFalse, "Bypass download cache")
	agyInstallCmd.Flags().BoolVar(&agyInstallDownloadMust, "force-download", isDefaultFalse, "Force re-download")
	agyInstallCmd.Flags().BoolVar(&agyInstallDownloadMust, "redownload", isDefaultFalse, "Force re-download")
}

func bindAgyInstallFlags() {
	bindAgyInstallBasicFlags()
	bindAgyInstallAdvancedFlags()
}

func bindAgyUpdateFlags() {
	agyUpdateCmd.Flags().BoolVarP(&agyInstallDryRun, "dry-run", "n", isDefaultFalse, "Simulate update without downloading")
	agyUpdateCmd.Flags().BoolVarP(&agyInstallYes, "yes", "y", isDefaultFalse, "Automatic yes to prompts")
	agyUpdateCmd.Flags().BoolVarP(&agyInstallVerbose, "verbose", "v", isDefaultFalse, "Enable verbose output")
	agyUpdateCmd.Flags().StringVar(&agyInstallVersion, "version", "", "Specific release version to update to (e.g. 2.19.1)")
	agyUpdateCmd.Flags().BoolVarP(&agyInstallForce, "force", "f", isDefaultFalse, "Force reinstallation even if already present")
	agyUpdateCmd.Flags().StringVar(&agyInstallPrefix, "prefix", "", "Custom installation directory prefix")
	bindAgyUpdateDownloadFlags()
}

func bindAgyUpdateDownloadFlags() {
	agyUpdateCmd.Flags().BoolVar(&agyInstallDownloadMust, "download-must", isDefaultFalse, "Bypass download cache")
	agyUpdateCmd.Flags().BoolVar(&agyInstallDownloadMust, "force-download", isDefaultFalse, "Force re-download")
	agyUpdateCmd.Flags().BoolVar(&agyInstallDownloadMust, "redownload", isDefaultFalse, "Force re-download")
}

func buildAgyInstallOptions() installOptions {
	return installOptions{
		DryRun:         agyInstallDryRun,
		Yes:            agyInstallYes,
		Verbose:        agyInstallVerbose,
		Version:        agyInstallVersion,
		Force:          agyInstallForce,
		Prefix:         agyInstallPrefix,
		IsDownloadMust: agyInstallDownloadMust,
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
	if len(args) > 0 && (strings.EqualFold(args[0], "update") || strings.EqualFold(args[0], "up")) {
		opts.Force = true
		args = args[1:]
	}
	target := resolveAgyTarget(args)

	return dispatchAgyInstallTarget(target, opts)
}

func runAgyUpdateCmd(cmd *cobra.Command, args []string) error {
	opts := buildAgyInstallOptions()
	opts.Force = true
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
