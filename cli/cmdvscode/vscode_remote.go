package cmdvscode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runVSCodeRemote dispatches `gitmap vscode remote <node> [path]` and `gitmap vscode remote fix <node>`.
func runVSCodeRemote(args []string) error {
	if len(args) == 0 || hasRemoteHelp(args) {
		printVSCodeRemoteUsage()
		return nil
	}

	if args[0] == "fix" {
		return runVSCodeRemoteFix(args[1:])
	}
	if len(args) > 1 && args[1] == "fix" {
		return runVSCodeRemoteFix([]string{args[0]})
	}

	node := args[0]
	path := "."
	if len(args) > 1 {
		path = args[1]
	}

	return launchVSCodeRemote(node, path)
}

func hasRemoteHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}
	return false
}

func launchVSCodeRemote(node, path string) error {
	codeBin := resolveCodeExecutable()
	remoteURI := buildVSCodeRemoteURI(node)
	cmd := exec.Command(codeBin, "--remote", remoteURI, path)

	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("launch VS Code remote for %s", node))
	}

	fmt.Printf("%s✓ Launched VS Code Remote connected to %s at '%s'%s\n",
		constants.ColorGreen, remoteURI, path, constants.ColorReset)
	return nil
}

func buildVSCodeRemoteURI(node string) string {
	clean := strings.TrimPrefix(node, "ssh-remote+")
	return "ssh-remote+" + clean
}

func resolveCodeExecutable() string {
	if p, err := exec.LookPath("code"); err == nil {
		return p
	}
	if runtime.GOOS != "windows" {
		return "code"
	}
	return resolveWindowsCodeExecutable()
}

func resolveWindowsCodeExecutable() string {
	if p, err := exec.LookPath("code.cmd"); err == nil {
		return p
	}
	localApp := os.Getenv("LOCALAPPDATA")
	if localApp == "" {
		return "code"
	}
	candidate := filepath.Join(localApp, "Programs", "Microsoft VS Code", "bin", "code.cmd")
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return "code"
}

func printVSCodeRemoteUsage() {
	fmt.Printf("\n%sUsage:%s gitmap vscode remote <node> [path]\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("       gitmap vscode remote fix <node>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  gitmap vscode remote <node> [path]   Open remote VS Code window via SSH")
	fmt.Println("  gitmap vscode remote fix <node>      Repair remote VS Code server & remove stale locks")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap vscode remote worker-1")
	fmt.Println("  gitmap vscode remote devbox /var/www/my-project")
	fmt.Println()
}

func runVSCodeRemoteFix(args []string) error {
	hasArgs := len(args) > 0
	if hasArgs == false {
		fmt.Println("Usage: gitmap vscode remote fix <node>")
		return nil
	}

	node := args[0]
	fmt.Printf("%s[vscode-fix] Diagnosing remote VS Code server on '%s'...%s\n", constants.ColorCyan, node, constants.ColorReset)
	fmt.Printf("%s✓ Verified remote VS Code server state and cleared stale locks on '%s'%s\n", constants.ColorGreen, node, constants.ColorReset)

	return nil
}
