package cmdbrowse

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
)

// LaunchBrowserAppWindow attempts to launch targetURL in Chrome or Edge in --app mode,
// or opens in the system default browser if app mode is disabled or no browser binary is found.
// When GITMAP_BROWSER_MODE is "default", "browser", or "no-app", it directly opens using the default system browser.
func LaunchBrowserAppWindow(targetURL string) bool {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("GITMAP_BROWSER_MODE")))
	if mode == "default" || mode == "browser" || mode == "no-app" {
		return false
	}

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
