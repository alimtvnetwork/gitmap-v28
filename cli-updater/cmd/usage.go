package cmd

import "fmt"

// PrintUsage displays the updater help text.
func PrintUsage() {
	fmt.Printf("cli-updater %s\n\n", Version)
	fmt.Println("Update gitmap via GitHub releases (no source repo required).")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cli-updater check       Check if a newer version is available")
	fmt.Println("  cli-updater run         Download and install the latest version")
	fmt.Println("  cli-updater version     Show updater version")
	fmt.Println("  cli-updater help        Show this help")
	fmt.Println()
}
