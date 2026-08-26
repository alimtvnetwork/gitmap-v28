// Package cmd — prompt_remediation.go formats manual installation help.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func PrintPromptManualInstructions() {
	fmt.Println("\nTo install manually:")
	fmt.Printf("  • Unix/macOS: curl -sL %s | bash\n", constants.PromptArchitectBashURL)
	fmt.Printf("  • Windows:    irm %s | iex\n\n", constants.PromptArchitectPowerShellURL)
}
