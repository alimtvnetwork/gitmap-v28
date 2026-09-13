package cmdinstall

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

var (
	agyInstallDryRun  bool
	agyInstallYes     bool
	agyInstallVerbose bool
	agyInstallVersion string
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
	target := "cli"
	if len(args) > 0 {
		target = strings.ToLower(args[0])
	}

	return dispatchAgyInstallTarget(target, opts)
}

func dispatchAgyInstallTarget(target string, opts installOptions) error {
	switch target {
	case "cli", "agy", "agy-cli", "antigravity-cli":
		return runInstallAgyWithOpts(opts)
	case "ide", "desktop", "app", "antigravity", "antigravity-ide":
		return runInstallAntigravityWithOpts(opts)
	case "manager", "ag-manager", "agm", "gui", "tools":
		return runInstallAgManagerWithOpts(opts)
	case "all", "both":
		return runInstallAgyAll(opts)
	default:
		fmt.Printf("Unknown target '%s'. Installing Antigravity CLI (default)...\n", target)

		return runInstallAgyWithOpts(opts)
	}
}

func runInstallAgyAll(opts installOptions) error {
	fmt.Println("=== [1/3] Installing Antigravity CLI ===")
	if err := runInstallAgyWithOpts(opts); err != nil {
		return err
	}

	fmt.Println("\n=== [2/3] Installing Antigravity Desktop IDE ===")
	if err := runInstallAntigravityWithOpts(opts); err != nil {
		return err
	}

	fmt.Println("\n=== [3/3] Installing Antigravity Manager / Tools ===")

	return runInstallAgManagerWithOpts(opts)
}
