// Package cmd — cg_version_cmd.go handles the CLI command for installing version.json.
package cmd

import (
	"fmt"
	"path/filepath"
)

// runCGInstallVersionJSON handles `gitmap cg install-version-json [targets...] [--version=<ver>] [--all]`.
func runCGInstallVersionJSON(targetDirs []string, initialVersion string, isDryRun bool) error {
	if len(targetDirs) == 0 {
		resolved, err := ResolvePromptTarget("")
		if err == nil {
			targetDirs = resolved
		}
	}

	if len(targetDirs) == 0 {
		fmt.Println("No target repositories found to install version.json.")
		return nil
	}

	cfg := DefaultVersionInstallConfig(initialVersion)
	fmt.Printf("→ Installing canonical version.json (v%s) in %d repository(ies)...\n", cfg.InitialVersion, len(targetDirs))

	var successCount, failCount int
	for _, dir := range targetDirs {
		name := filepath.Base(dir)
		fmt.Printf("  • %s (%s)... ", name, dir)
		err := InstallVersionJSON(dir, cfg, isDryRun)
		if err == nil {
			fmt.Println("✓ Done")
			successCount++
		} else {
			fmt.Printf("✖ Failed: %v\n", err)
			failCount++
		}
	}

	fmt.Printf("\nversion.json Installation Summary: %d Success | %d Failed\n", successCount, failCount)
	return nil
}
