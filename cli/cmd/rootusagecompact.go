package cmd

import (
	"fmt"
)

// printUsageCompact outputs the minimal root help summary when gitmap
// is invoked with zero arguments, ensuring total output remains under 15 lines.
func printUsageCompact() {
	PrintBinaryLocations()

	printGitmapIdentityBlockShort()

	printCompactGuidanceFooter()
}

func printCompactGuidanceFooter() {
	fmt.Println("Tip: Run 'gitmap help' or 'gitmap -h' to see the full command catalog.")
	fmt.Println("     Run 'gitmap help <topic>' for command-specific documentation.")
}
