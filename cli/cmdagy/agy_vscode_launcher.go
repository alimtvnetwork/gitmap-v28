package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func launchVSCodeForLastProjects(projCount, promptLimit int) *apperror.AppError {
	projects := discoverRecentCommitProjects(projCount)
	tempFile, genErr := generateVSCodeInspectionFile(projects, promptLimit)
	if genErr != nil {
		return apperror.WrapSimple(genErr, "generate inspection file")
	}

	fmt.Printf("\n  %s● Generated prompt history file:%s %s\n", constants.ColorCyan, constants.ColorReset, tempFile)

	return executeVSCodeWindow(tempFile)
}

func executeVSCodeWindow(tempFile string) *apperror.AppError {
	launchErr := spawnNonAdminVSCode(tempFile)
	if launchErr != nil {
		fmt.Printf("  %s⚠ Could not launch VS Code directly: %v%s\n", constants.ColorYellow, launchErr, constants.ColorReset)
		fmt.Printf("    Open manually: code -n %s\n\n", tempFile)

		return nil
	}

	fmt.Printf("  %s✔ Launched new non-admin VS Code window (code -n).%s\n\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func spawnNonAdminVSCode(filePath string) error {
	bin := resolveVSCodeBinary()
	cmd := exec.Command(bin, "-n", filePath)
	cmd.Env = filterAdminEnv(os.Environ())
	configureBackgroundProcess(cmd)

	return cmd.Start()
}

func filterAdminEnv(env []string) []string {
	var filtered []string
	for _, e := range env {
		if !isAdminEnvEntry(e) {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

func isAdminEnvEntry(e string) bool {
	upper := strings.ToUpper(e)

	return strings.HasPrefix(upper, "ELEVATED=") ||
		strings.HasPrefix(upper, "__COMPAT_LAYER=")
}

func resolveVSCodeBinary() string {
	for _, cand := range getVSCodeCandidates() {
		path, err := exec.LookPath(cand)
		if err == nil {
			return path
		}
	}

	return "code"
}

func getVSCodeCandidates() []string {
	if runtime.GOOS == "windows" {
		return []string{"code.cmd", "code"}
	}

	return []string{"code"}
}
