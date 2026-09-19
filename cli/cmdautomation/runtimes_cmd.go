package cmdautomation

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var runtimesCmd = &cobra.Command{
	Use:     "runtimes [list|refresh]",
	Aliases: []string{"runtime", "rt"},
	Short:   "Discover, probe, and inspect polyglot language runtimes",
	RunE:    runRuntimesCmd,
}

func runRuntimesCmd(cmd *cobra.Command, args []string) error {
	action := resolveRuntimeAction(args)
	if action == "refresh" {
		return handleRuntimesRefresh()
	}
	return handleRuntimesList()
}

func resolveRuntimeAction(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "list"
}

func handleRuntimesList() error {
	res := ListRuntimes()
	if res.IsFailure() {
		return res.Err
	}
	renderRuntimesTable(res.Value)
	return nil
}

func handleRuntimesRefresh() error {
	res := RefreshRuntimes()
	if res.IsFailure() {
		return res.Err
	}
	fmt.Printf("\n%s✔ Refreshed polyglot runtime cache in SQLite.%s\n", constants.ColorGreen, constants.ColorReset)
	renderRuntimesTable(res.Value)
	return nil
}

func renderRuntimesTable(runtimes []RuntimeInfo) {
	fmt.Printf("\n%s[Discovered Polyglot Runtimes]%s\n", constants.ColorBold, constants.ColorReset)
	if len(runtimes) == 0 {
		fmt.Printf("  %sNo language runtimes discovered or cached.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	printRuntimesHeader()
	for _, rt := range runtimes {
		printRuntimeRow(rt)
	}
	fmt.Println()
}

func printRuntimesHeader() {
	fmt.Printf("%-12s | %-18s | %-10s | %s\n", "Runtime", "Version", "Status", "Cached Path")
	fmt.Println("-------------|--------------------|------------|---------------------------------------------")
}

func printRuntimeRow(rt RuntimeInfo) {
	statusColor := getRuntimeStatusColor(rt.Status)
	statusText := fmt.Sprintf("%s%-10s%s", statusColor, rt.Status, constants.ColorReset)
	fmt.Printf("%-12s | %-18s | %s | %s\n", rt.Name, rt.Version, statusText, rt.BinaryPath)
}

func getRuntimeStatusColor(status string) string {
	if status == "active" {
		return constants.ColorGreen
	}
	if status == "missing" {
		return constants.ColorRed
	}
	return constants.ColorYellow
}

func init() {
	AutomationCmd.AddCommand(runtimesCmd)
}
