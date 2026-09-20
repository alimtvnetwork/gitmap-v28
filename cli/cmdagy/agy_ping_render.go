package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderPingReport(report AgyPingReport) {
	fmt.Printf("\n%s● Antigravity IDE Ping & Health Report%s\n", constants.ColorCyan, constants.ColorReset)
	renderPingExecutable(report.Executable)
	renderPingProcess(report.Process)
	renderPingFilesystem(report.Filesystem)
	renderPingWorkspace(report.Workspace)
	renderPingQueue(report.Queue)
	renderPingOverall(report.IsHealthy)
	fmt.Println()
}

func renderPingJSON(report AgyPingReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal ping report json")
	}
	fmt.Println(string(data))

	return nil
}

func renderPingExecutable(chk AgyPingExecutableCheck) {
	if chk.IsFound {
		fmt.Printf("  IDE Executable: %sFOUND%s (%s)\n", constants.ColorGreen, constants.ColorReset, chk.Path)

		return
	}
	fmt.Printf("  IDE Executable: %sNOT FOUND%s\n", constants.ColorYellow, constants.ColorReset)
}

func renderPingProcess(chk AgyPingProcessCheck) {
	if chk.IsRunning {
		fmt.Printf("  IDE Process:    %sRUNNING%s (PID: %d, %s)\n", constants.ColorGreen, constants.ColorReset, chk.PID, chk.Name)

		return
	}
	fmt.Printf("  IDE Process:    %sSTOPPED%s (Offline filesystem mode supported)\n", constants.ColorYellow, constants.ColorReset)
}

func renderPingFilesystem(chk AgyFilesystemHealth) {
	if chk.IsAccessible {
		fmt.Printf("  Filesystem:     %sACCESSIBLE%s\n", constants.ColorGreen, constants.ColorReset)
		fmt.Printf("    • Brain Logs: %s\n", chk.BrainDir)
		fmt.Printf("    • Projects:   %s\n", chk.ProjectsDir)

		return
	}
	fmt.Printf("  Filesystem:     %sINACCESSIBLE%s (%s)\n", constants.ColorRed, constants.ColorReset, chk.Error)
}

func renderPingWorkspace(chk AgyPingWorkspaceCheck) {
	statusColor := constants.ColorGreen
	if chk.ConvStatus == "running" {
		statusColor = constants.ColorYellow
	}
	fmt.Printf("  Workspace:      %s\n", chk.TargetWorkspace)
	if len(chk.ConvID) > 0 {
		fmt.Printf("    • Status:     %s%s%s (ID: %s, %d steps)\n", statusColor, strings.ToUpper(chk.ConvStatus), constants.ColorReset, chk.ConvID, chk.StepCount)

		return
	}
	fmt.Printf("    • Status:     %s%s%s (No active conversation)\n", statusColor, strings.ToUpper(chk.ConvStatus), constants.ColorReset)
}

func renderPingQueue(chk AgyPingQueueCheck) {
	if chk.HasQueue {
		fmt.Printf("  Prompt Queue:   %d active, %d queued\n", chk.ActivePromptCount, chk.QueuedCount)

		return
	}
	fmt.Printf("  Prompt Queue:   %sUNAVAILABLE%s (%s)\n", constants.ColorYellow, constants.ColorReset, chk.Error)
}

func renderPingOverall(isHealthy bool) {
	if isHealthy {
		fmt.Printf("  Overall Health: %sHEALTHY%s\n", constants.ColorGreen, constants.ColorReset)

		return
	}
	fmt.Printf("  Overall Health: %sDEGRADED%s\n", constants.ColorRed, constants.ColorReset)
}
