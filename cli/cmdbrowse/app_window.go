package cmdbrowse

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
)

// LaunchBrowserAppWindow attempts to launch targetURL in Chrome or Edge in --app mode.
// Returns true if a browser is found and successfully started.
func LaunchBrowserAppWindow(targetURL string) bool {
	bin := findAppBrowserBinary()
	if bin == "" {
		return false
	}

	cmd := exec.Command(bin, "--app="+targetURL, "--new-window")
	cmdchrome.ConfigureDetachedProcess(cmd)

	return cmd.Start() == nil
}

func findAppBrowserBinary() string {
	if chromeBin, err := cmdchrome.FindChromeBinaryPath(); err == nil && chromeBin != "" {
		return chromeBin
	}

	return findEdgeBinaryPath()
}

func findEdgeBinaryPath() string {
	switch runtime.GOOS {
	case "windows":
		return findEdgeWindows()
	case "darwin":
		return findEdgeDarwin()
	default:
		return findEdgeLinux()
	}
}

func findEdgeWindows() string {
	candidates := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "Application", "msedge.exe"),
	}

	for _, cand := range candidates {
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}

	if p, err := exec.LookPath("msedge.exe"); err == nil {
		return p
	}

	return ""
}

func findEdgeDarwin() string {
	candidates := []string{
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		filepath.Join(os.Getenv("HOME"), "Applications", "Microsoft Edge.app", "Contents", "MacOS", "Microsoft Edge"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func findEdgeLinux() string {
	for _, name := range []string{"microsoft-edge-stable", "microsoft-edge", "msedge"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}

	return ""
}
