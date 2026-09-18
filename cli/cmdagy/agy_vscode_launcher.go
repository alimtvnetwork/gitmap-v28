package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func launchVSCodeForLastProjects(projCount, promptLimit int) error {
	projects := discoverRecentCommitProjects(projCount)
	tempFile, genErr := generateVSCodeInspectionFile(projects, promptLimit)
	if genErr != nil {
		return genErr
	}
	fmt.Printf("\n  %s● Generated prompt history file:%s %s\n", constants.ColorCyan, constants.ColorReset, tempFile)

	return executeVSCodeWindow(tempFile)
}

func executeVSCodeWindow(tempFile string) error {
	if launchErr := spawnNonAdminVSCode(tempFile); launchErr != nil {
		fmt.Printf("  %s⚠ Could not launch VS Code directly: %v%s\n", constants.ColorYellow, launchErr, constants.ColorReset)
		fmt.Printf("    Open manually: code -n %s\n\n", tempFile)
		return nil
	}
	fmt.Printf("  %s✔ Launched new non-admin VS Code window (code -n).%s\n\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func spawnNonAdminVSCode(filePath string) error {
	cmd := exec.Command("code", "-n", filePath)
	cmd.Env = filterAdminEnv(os.Environ())

	return cmd.Start()
}

func filterAdminEnv(env []string) []string {
	var filtered []string
	for _, e := range env {
		if !strings.HasPrefix(strings.ToUpper(e), "ELEVATED=") {
			filtered = append(filtered, e)
		}
	}

	return filtered
}
