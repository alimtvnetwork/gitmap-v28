package cmd

import "fmt"

// PrintUsage displays the updater help text.
func PrintUsage() {
	fmt.Printf("gitmap-updater %s\n\n", Version)
	fmt.Println("Update gitmap via GitHub releases (no source repo required).")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitmap-updater check       Check if a newer version is available")
	fmt.Println("  gitmap-updater run         Download and install the latest version")
	fmt.Println("  gitmap-updater version     Show updater version")
	fmt.Println("  gitmap-updater help        Show this help")
	fmt.Println()
}
