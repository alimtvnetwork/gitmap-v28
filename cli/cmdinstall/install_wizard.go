package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	wizardTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9")).MarginBottom(1)
	wizardSecStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b")).Underline(true)
)

func runInstallWizard() error {
	printWizardHeader()
	installed := loadInstalledLookup()
	printWizardToolsGrid(installed)
	printInstallProfilesSection(installed)
	printWizardContextMenuStatus()
	printWizardActions()

	return nil
}

func printWizardHeader() {
	title := fmt.Sprintf("GitMap Developer Tools & Profiles Install Wizard (v%s)", constants.Version)
	fmt.Println(wizardTitleStyle.Render(title))
}

func printWizardToolsGrid(installed map[string]string) {
	fmt.Println(wizardSecStyle.Render("Core Developer Tools Status"))
	tools := []string{
		constants.ToolGit,
		constants.ToolGo,
		constants.ToolNodeJS,
		constants.ToolVSCode,
		constants.ToolBeyondCompare,
		constants.ToolAntigravity,
		constants.ToolPowerShell,
		constants.ToolPython,
	}
	for _, tool := range tools {
		printToolRow(tool, installed)
	}

	fmt.Println()
}

func printWizardContextMenuStatus() {
	fmt.Println(wizardSecStyle.Render("OS Context Menu Integration"))
	if isContextMenuActive() {
		fmt.Printf("  %s %s\n", installedDot, "Active (Registered in OS file explorer)")
		fmt.Println("    To reinstall: gitmap install context-menu")
		fmt.Println("    To uninstall: gitmap uninstall context-menu")
		fmt.Println()

		return
	}

	fmt.Printf("  %s %s\n", missingDot, "Inactive (Not registered in OS file explorer)")
	fmt.Println("    To install:   gitmap install context-menu")
	fmt.Println()
}

func printWizardActions() {
	fmt.Println(wizardSecStyle.Render("Wizard Quick Commands"))
	fmt.Println("  • Install specific tool:    gitmap install <tool>")
	fmt.Println("  • Install profile:          gitmap install profile <name>")
	fmt.Println("  • Preview profile tree:     gitmap install profile <name> --tree")
	fmt.Println("  • Manage context menu:      gitmap context-menu [install|uninstall|status]")
	fmt.Println("  • Clean context removal:    gitmap uninstall context-menu")
	fmt.Println()
}

func isContextMenuActive() bool {
	if runtime.GOOS == "windows" {
		return isWindowsCtxActive()
	}
	if runtime.GOOS == "linux" {
		return isLinuxCtxActive()
	}
	if runtime.GOOS == "darwin" {
		return isDarwinCtxActive()
	}

	return false
}

func isWindowsCtxActive() bool {
	cmd := exec.Command("reg", "query", constants.CtxRootKeyDirectory)

	return cmd.Run() == nil
}

func isLinuxCtxActive() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	nautilusPath := filepath.Join(home, constants.CtxLinuxNautilusRel)
	if fi, err := os.Stat(nautilusPath); err == nil && fi.IsDir() {
		return true
	}
	dolphinPath := filepath.Join(home, constants.CtxLinuxDolphinRel, constants.CtxLinuxDolphinFile)
	_, dolphinErr := os.Stat(dolphinPath)

	return dolphinErr == nil
}

func isDarwinCtxActive() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	p := filepath.Join(home, constants.CtxMacServicesRel)
	fi, err := os.Stat(p)

	return err == nil && fi.IsDir()
}
