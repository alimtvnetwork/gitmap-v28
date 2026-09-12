package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runInstallAntigravityWithOpts(opts installOptions) error {
	path, isFound := findInstalledAntigravityDesktopPath()
	if isFound {
		fmt.Printf("  ✓ Google Antigravity Desktop IDE is already installed (%s)\n", path)

		return nil
	}

	if opts.DryRun {
		fmt.Println("  [dry-run] Would download and install Google Antigravity Desktop IDE")

		return nil
	}

	return performAntigravityDesktopInstall(opts)
}

func performAntigravityDesktopInstall(opts installOptions) error {
	fmt.Println("Installing Google Antigravity Desktop IDE...")
	if err := installAntigravityDesktopPlatform(opts); err != nil {
		reportVerificationFailure(constants.ToolAntigravity, "antigravity")

		return fmt.Errorf("antigravity desktop IDE installation failed: %w", err)
	}

	path, isFound := findInstalledAntigravityDesktopPath()
	if isFound {
		fmt.Printf(constants.ColorGreen+"✓"+constants.ColorReset+" Google Antigravity Desktop IDE installed: %s\n", path)
		recordAntigravityDesktopInstalled(path)

		return nil
	}

	reportVerificationFailure(constants.ToolAntigravity, "antigravity")

	return fmt.Errorf("antigravity desktop IDE installation verification failed")
}
