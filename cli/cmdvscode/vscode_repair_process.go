package cmdvscode

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/lipgloss"
)

var (
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8"))
	successMark = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("✓")
	warnMark    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c")).Render("!")
	bulletMark  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Render("•")
)

func terminateVSCodeProcesses() int {
	if runtime.GOOS == "windows" {
		return terminateWindowsProcesses()
	}

	return terminateUnixProcesses()
}

func terminateWindowsProcesses() int {
	targets := []string{"Code.exe", "inno_updater.exe", "code-tunnel.exe"}
	terminated := 0
	for _, target := range targets {
		if killWindowsProcess(target) {
			terminated++
		}
	}

	return terminated
}

func killWindowsProcess(name string) bool {
	cmd := exec.Command("taskkill", "/F", "/IM", name)
	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Printf("    %s Terminated running process: %s\n", bulletMark, subtleStyle.Render(name))
		return true
	}

	_ = out
	return false
}

func terminateUnixProcesses() int {
	targets := []string{"code"}
	terminated := 0
	for _, target := range targets {
		if killUnixProcess(target) {
			terminated++
		}
	}

	return terminated
}

func killUnixProcess(name string) bool {
	cmd := exec.Command("pkill", "-f", name)
	if err := cmd.Run(); err == nil {
		fmt.Printf("    %s Terminated running process: %s\n", bulletMark, subtleStyle.Render(name))
		return true
	}

	return false
}
