package cmdos

import (
	"fmt"
)

func runOSUpdateCommand(isUpgrade bool, args []string) error {
	if len(args) > 0 && isOSHelpArg(args[0]) {
		printOSUpdateHelp(isUpgrade)

		return nil
	}

	isDryRun := hasDryRunFlag(args)
	action := "Updating package indices"
	if isUpgrade {
		action = "Upgrading installed packages & tools"
	}

	fmt.Printf("▶ %s across detected package managers...\n", action)
	results := runAllUpdates(isUpgrade, isDryRun)
	printUpdateSummary(results)

	return nil
}

func hasDryRunFlag(args []string) bool {
	for _, a := range args {
		if a == "-n" || a == "--dry-run" {
			return true
		}
	}

	return false
}

func printUpdateSummary(results []UpdateResult) {
	if len(results) == 0 {
		fmt.Println("ℹ No supported package managers discovered on this system.")

		return
	}

	fmt.Println("▶ Summary of Package Manager Results:")
	for _, r := range results {
		status := "✔ OK"
		if !r.Success {
			status = "✖ FAILED"
		}
		fmt.Printf("  • %-10s : %s\n", r.Name, status)
	}
}

func printOSUpdateHelp(isUpgrade bool) {
	cmd := "update"
	desc := "Check and update package repository metadata"
	if isUpgrade {
		cmd = "upgrade"
		desc = "Perform full package upgrade across all detected managers (winget, apt, dnf, pacman, brew)"
	}

	fmt.Printf("Usage: gitmap os %s [flags]\n\n", cmd)
	fmt.Printf("Description:\n  %s\n\n", desc)
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run    Simulate update without modifying system")
	fmt.Println("  -h, --help       Show this help message")
}
