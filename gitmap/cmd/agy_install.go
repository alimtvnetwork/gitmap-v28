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
	Use:     "install [manager|ide|cli|all] [flags]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity Manager, Desktop IDE, or CLI",
	Long: `Install Antigravity tools:
  manager    - Install latest release of Antigravity Manager GUI (default)
  ide (app)  - Install Google Antigravity Desktop IDE
  cli (agy)  - Install Antigravity CLI coding assistant
  all        - Install Antigravity Manager, Desktop IDE, and CLI`,
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
	case "cli", "agy":
		return runInstallAgyWithOpts(opts)
	case "ide", "desktop", "app", "antigravity":
		return runInstallAntigravityWithOpts(opts)
	case "all", "both":
		return runInstallAgyAll(opts)
	default:
		fmt.Printf("Unknown target '%s'. Installing Antigravity Manager (default)...\n", target)

		return runInstallAgManagerWithOpts(opts)
	}
}

func runInstallAgyAll(opts installOptions) error {
	fmt.Println("=== [1/3] Installing Antigravity Manager ===")
	if err := runInstallAgManagerWithOpts(opts); err != nil {
		return err
	}

	fmt.Println("\n=== [2/3] Installing Antigravity Desktop IDE ===")
	if err := runInstallAntigravityWithOpts(opts); err != nil {
		return err
	}

	fmt.Println("\n=== [3/3] Installing Antigravity CLI ===")

	return runInstallAgyWithOpts(opts)
}
