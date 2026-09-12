// Package cmdprompt — prompt_dry_run.go simulates prompt installation without making changes.
package cmdprompt

import "fmt"

func SimulatePromptInstallation(targetDir string) {
	fmt.Printf("  [DRY-RUN] Would install Prompt Architect v2 in: %s\n", targetDir)
}
