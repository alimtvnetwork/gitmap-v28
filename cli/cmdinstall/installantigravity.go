package cmdinstall

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

func fallbackToAgyCli(opts installOptions, origErr error) error {
	fmt.Printf("  ⚠ Antigravity desktop IDE package is unavailable (%v)\n", origErr)
	fmt.Println("  → Falling back to Google Antigravity CLI (agy)...")
	if cliErr := runInstallAgyWithOpts(opts); cliErr == nil {
		recordAntigravityDesktopInstalled("antigravity-cli")
		return nil
	}
	reportVerificationFailure(constants.ToolAntigravity, "antigravity")

	return fmt.Errorf("antigravity desktop IDE installation failed: %w", origErr)
}

func performAntigravityDesktopInstall(opts installOptions) error {
	fmt.Println("Installing Google Antigravity Desktop IDE...")
	if err := installAntigravityDesktopPlatform(opts); err != nil {
		return fallbackToAgyCli(opts, err)
	}

	path, isFound := findInstalledAntigravityDesktopPath()
	if !isFound {
		return fallbackToAgyCli(opts, fmt.Errorf("desktop application binary not found"))
	}

	fmt.Printf(constants.ColorGreen+"✓"+constants.ColorReset+" Google Antigravity Desktop IDE installed: %s\n", path)
	recordAntigravityDesktopInstalled(path)

	return nil
}
