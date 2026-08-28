// Package cmd — cg_prompts.go executes Prompt Architect installation from cg command.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func runCGInstallPrompts(targetDirs []string, isDryRun bool) error {
	if len(targetDirs) == 0 {
		resolved, err := ResolvePromptTarget("")
		if err == nil {
			targetDirs = resolved
		}
	}

	if len(targetDirs) == 0 {
		fmt.Println("No target repositories found to install Prompt Architect.")
		return nil
	}

	fmt.Printf("→ Installing Prompt Architect v2 in %d repository(ies)...\n", len(targetDirs))
	var results []model.PromptInstallResult

	for _, dir := range targetDirs {
		name := filepath.Base(dir)
		fmt.Printf("  • %s (%s)... ", name, dir)
		res := ExecuteSinglePromptInstall(dir, isDryRun)
		results = append(results, res)
		if res.IsSuccess {
			fmt.Println("✓ Done")
		} else {
			fmt.Printf("✖ Failed: %s\n", res.Error)
		}
	}

	RenderPromptInstallSummary(results)
	ReportPromptFailures(results)
	return nil
}
