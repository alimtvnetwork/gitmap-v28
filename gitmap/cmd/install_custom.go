package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var customToolUrls = map[string]struct{ win, unix string }{
	"scripts-fixer": {
		win:  "https://raw.githubusercontent.com/alimtvnetwork/scripts-fixer-v20/main/install.ps1",
		unix: "https://raw.githubusercontent.com/alimtvnetwork/scripts-fixer-v20/main/install.sh",
	},
	"coding-guidelines": {
		win:  "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.ps1",
		unix: "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh",
	},
	"macro-ahk": {
		win:  "https://raw.githubusercontent.com/alimtvnetwork/macro-ahk-v55/main/scripts/download-extension.ps1",
		unix: "https://raw.githubusercontent.com/alimtvnetwork/macro-ahk-v55/main/scripts/install.sh",
	},
}

func runInstallCustomTool(tool string) {
	urls, ok := customToolUrls[tool]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown custom tool: %s\n", tool)
		os.Exit(1)
	}

	cmd := buildCustomToolCmd(urls.win, urls.unix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to install %s: %v\n", tool, err)
		os.Exit(1)
	}
}

func buildCustomToolCmd(winUrl, unixUrl string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		psCmd := fmt.Sprintf("irm %s | iex", winUrl)
		return exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	}
	shCmd := fmt.Sprintf("curl -fsSL %s | bash", unixUrl)
	return exec.Command("bash", "-c", shCmd)
}
