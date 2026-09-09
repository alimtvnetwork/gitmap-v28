package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var agyInstallCmd = &cobra.Command{
	Use:     "install [manager|cli|all]",
	Aliases: []string{"in", "i"},
	Short:   "Install Antigravity Manager GUI or Antigravity CLI",
	Long: `Install Antigravity tools:
  manager    - Install latest release of Antigravity Manager GUI (default)
  cli (agy)  - Install Antigravity CLI coding assistant
  all        - Install both Antigravity Manager and Antigravity CLI`,
	RunE: runAgyInstallCmd,
}

func runAgyInstallCmd(cmd *cobra.Command, args []string) error {
	target := "manager"
	if len(args) > 0 {
		target = strings.ToLower(args[0])
	}

	return dispatchAgyInstallTarget(target)
}

func dispatchAgyInstallTarget(target string) error {
	switch target {
	case "manager", "ag-manager", "gui":

		return runInstallAgManager()
	case "cli", "antigravity", "agy":

		return runInstallAntigravity()
	case "all", "both":

		return runInstallAgyAll()
	default:
		fmt.Printf("Unknown target '%s'. Installing Antigravity Manager (default)...\n", target)

		return runInstallAgManager()
	}
}

func runInstallAgyAll() error {
	fmt.Println("=== [1/2] Installing Antigravity Manager ===")
	if err := runInstallAgManager(); err != nil {

		return err
	}
	fmt.Println("\n=== [2/2] Installing Antigravity CLI ===")

	return runInstallAntigravity()
}
