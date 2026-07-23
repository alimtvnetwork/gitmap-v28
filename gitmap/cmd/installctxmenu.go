package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runVSCodeContextMenu adds VS Code to the Windows right-click context menu.
func runVSCodeContextMenu() {
	if runtime.GOOS != "windows" {
		fmt.Fprintf(os.Stderr, "  Error: vscode-ctx is only supported on Windows (current OS: %s)\n", runtime.GOOS)

		return
	}

	fmt.Println("  Adding VS Code to Windows context menu...")

	regCommands := [][]string{
		{"reg", "add", `HKCU\Software\Classes\Directory\Background\shell\VSCode`, "/ve", "/d", "Open with VS Code", "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\Background\shell\VSCode`, "/v", "Icon", "/d", `C:\Program Files\Microsoft VS Code\Code.exe`, "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\Background\shell\VSCode\command`, "/ve", "/d", `"C:\Program Files\Microsoft VS Code\Code.exe" "%V"`, "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\shell\VSCode`, "/ve", "/d", "Open with VS Code", "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\shell\VSCode`, "/v", "Icon", "/d", `C:\Program Files\Microsoft VS Code\Code.exe`, "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\shell\VSCode\command`, "/ve", "/d", `"C:\Program Files\Microsoft VS Code\Code.exe" "%V"`, "/f"},
	}

	runRegistryCommands("VS Code", regCommands)
}

// runPwshContextMenu adds PowerShell 7 to the Windows right-click context menu.
func runPwshContextMenu() {
	if runtime.GOOS != "windows" {
		fmt.Fprintf(os.Stderr, "  Error: pwsh-ctx is only supported on Windows (current OS: %s)\n", runtime.GOOS)

		return
	}

	fmt.Println("  Adding PowerShell 7 to Windows context menu...")

	regCommands := [][]string{
		{"reg", "add", `HKCU\Software\Classes\Directory\Background\shell\pwsh`, "/ve", "/d", "Open PowerShell 7 here", "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\Background\shell\pwsh`, "/v", "Icon", "/d", `C:\Program Files\PowerShell\7\pwsh.exe`, "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\Background\shell\pwsh\command`, "/ve", "/d", `"C:\Program Files\PowerShell\7\pwsh.exe" -NoExit -Command "Set-Location '%V'"`, "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\shell\pwsh`, "/ve", "/d", "Open PowerShell 7 here", "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\shell\pwsh`, "/v", "Icon", "/d", `C:\Program Files\PowerShell\7\pwsh.exe`, "/f"},
		{"reg", "add", `HKCU\Software\Classes\Directory\shell\pwsh\command`, "/ve", "/d", `"C:\Program Files\PowerShell\7\pwsh.exe" -NoExit -Command "Set-Location '%V'"`, "/f"},
	}

	runRegistryCommands("PowerShell 7", regCommands)
}

// runRegistryCommands executes a set of registry commands and reports results.
func runRegistryCommands(label string, commands [][]string) {
	success := 0

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "  ! Registry command failed: %v\n", err)

			continue
		}

		success++
	}

	fmt.Printf("  ✓ %s context menu entries added (%d/%d registry keys).\n", label, success, len(commands))
}

// runAllDevTools installs all core developer tools sequentially.
func runAllDevTools(opts installOptions) {
	allTools := []string{
		constants.ToolGit,
		constants.ToolGitLFS,
		constants.ToolGo,
		constants.ToolNodeJS,
		constants.ToolPnpm,
		constants.ToolPython,
		constants.ToolVSCode,
		constants.ToolGitHubDesktop,
		constants.ToolPowerShell,
		constants.ToolCPP,
		constants.ToolNpp,
	}

	fmt.Printf("\n  Installing %d core developer tools...\n\n", len(allTools))

	installed := 0
	skipped := 0

	for _, tool := range allTools {
		fmt.Printf("  --- %s ---\n", tool)

		existing := detectInstalledVersion(tool)
		if existing != "" {
			fmt.Printf("  ✓ %s already installed, skipping.\n\n", tool)
			skipped++

			continue
		}

		toolOpts := installOptions{
			Tool:    tool,
			Manager: opts.Manager,
			Version: "",
			Verbose: opts.Verbose,
			DryRun:  opts.DryRun,
			Check:   false,
			Yes:     true,
		}

		executeInstall(toolOpts)
		installed++
		fmt.Println()
	}

	fmt.Printf("\n  ✅ All dev tools: %d installed, %d already present.\n", installed, skipped)
}
