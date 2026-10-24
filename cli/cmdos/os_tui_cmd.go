package cmdos

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// RunOSTUICommand launches the interactive OS configuration dashboard.
func RunOSTUICommand(args []string) error {
	if isTUIHelpRequested(args) {
		printTUIUsage()
		return nil
	}

	isDryRun := hasTUIDryRunFlag(args)
	model := NewOSTUIModel(isDryRun)
	prog := tea.NewProgram(model)

	_, err := prog.Run()
	return err
}

func isTUIHelpRequested(args []string) bool {
	return anyArgIn(args, "-h", "--help", "help")
}

func hasTUIDryRunFlag(args []string) bool {
	return anyArgIn(args, "--dry-run", "-n")
}

func printTUIUsage() {
	fmt.Println("Usage: gitmap os tui [flags]")
	fmt.Println()
	fmt.Println("Interactive terminal dashboard for OS tweaks, auto-login, display, DNS, and maintenance.")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --dry-run, -n    Simulate execution without applying real system changes")
	fmt.Println("  -h, --help       Show help for os tui")
	fmt.Println()
	fmt.Println("Keybindings:")
	fmt.Println("  Tab / ← / →      Switch tabs")
	fmt.Println("  1-6              Jump directly to tab 1 through 6")
	fmt.Println("  ↑ / ↓ (k / j)    Navigate options")
	fmt.Println("  Space            Toggle selection")
	fmt.Println("  a / n            Select all / Deselect all in active tab")
	fmt.Println("  Enter            Apply all selected options")
	fmt.Println("  q / Esc          Quit dashboard")
}
