package macro

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var openLauncherFn = defaultOpenLauncher

// ParseOpenCommand checks if a macro step is an 'open' command.
func ParseOpenCommand(cmdText string) (bool, string) {
	trimmed := strings.TrimSpace(cmdText)
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return false, ""
	}
	if strings.ToLower(fields[0]) != "open" {
		return false, ""
	}
	if len(fields) == 1 {
		return true, "."
	}
	target := strings.TrimSpace(trimmed[len(fields[0]):])
	target = strings.Trim(target, "\"'")
	if len(target) == 0 {
		return true, "."
	}
	return true, target
}

func executeOpenStep(ctx context.Context, step MacroStep, cmdText, target, currentDir string, start time.Time, opts ExecOptions, idx int) (StepExecution, error) {
	err := openLauncherFn(ctx, target, currentDir)
	elapsed := time.Since(start)
	if err != nil {
		return handleStepFailure(step, cmdText, currentDir, elapsed, 1, err, opts, idx, []string{}, []string{err.Error()})
	}
	if !isStructuredOutput(opts) {
		fmt.Printf("%s✔ ok (%.1fs)%s\n", constants.ColorGreen, elapsed.Seconds(), constants.ColorReset)
	}
	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    cmdText,
		WorkingDir:     currentDir,
		Status:         "success",
		ExitCode:       0,
		ElapsedSeconds: elapsed.Seconds(),
		Logs:           []string{fmt.Sprintf("Opened %s", target)},
		ErrorLogs:      []string{},
	}, nil
}

func defaultOpenLauncher(ctx context.Context, target, currentDir string) error {
	if isChromeTarget(target) {
		return launchChrome(ctx)
	}
	if isURL, urlStr := parseURLTarget(target); isURL {
		return launchURL(ctx, urlStr)
	}
	resolvedPath := resolveTargetPath(target, currentDir)
	if isPathExists(resolvedPath) {
		return launchPath(ctx, resolvedPath)
	}
	return launchGeneric(ctx, target)
}

func isChromeTarget(target string) bool {
	lower := strings.ToLower(strings.TrimSpace(target))
	return lower == "chrome" || lower == "google-chrome" || lower == "google chrome" || lower == "chromium"
}

func parseURLTarget(target string) (bool, string) {
	trimmed := strings.TrimSpace(target)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return true, trimmed
	}
	if strings.Contains(trimmed, ".") && !strings.ContainsAny(trimmed, "/\\") && !strings.HasPrefix(trimmed, ".") {
		return true, "https://" + trimmed
	}
	if strings.HasPrefix(lower, "www.") {
		return true, "https://" + trimmed
	}
	return false, ""
}

func resolveTargetPath(target, currentDir string) string {
	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}
	return filepath.Clean(filepath.Join(currentDir, target))
}

func isPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func launchChrome(ctx context.Context) error {
	if runtime.GOOS == constants.OSWindows {
		return launchChromeWindows(ctx)
	}
	if runtime.GOOS == "darwin" {
		return exec.CommandContext(ctx, "open", "-a", "Google Chrome").Run()
	}
	return launchChromeLinux(ctx)
}

func findExistingChromeWindowsPath() (string, bool) {
	for _, p := range getChromeWindowsCandidates() {
		_, err := os.Stat(p)
		if err == nil {
			return p, true
		}
	}
	p, err := exec.LookPath("chrome.exe")
	if err == nil {
		return p, true
	}
	return "", false
}

func launchChromeWindows(ctx context.Context) error {
	p, hasChrome := findExistingChromeWindowsPath()
	if hasChrome {
		cmd := exec.CommandContext(ctx, p)
		return cmd.Start()
	}
	return exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", "Start-Process", "chrome").Run()
}

func getChromeWindowsCandidates() []string {
	return []string{
		filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
	}
}

func findExistingChromeLinuxPath() (string, bool) {
	for _, bin := range []string{"google-chrome", "google-chrome-stable", "chromium-browser", "chromium"} {
		p, err := exec.LookPath(bin)
		if err == nil {
			return p, true
		}
	}
	return "", false
}

func launchChromeLinux(ctx context.Context) error {
	p, hasChrome := findExistingChromeLinuxPath()
	if hasChrome {
		cmd := exec.CommandContext(ctx, p)
		return cmd.Start()
	}
	return exec.CommandContext(ctx, "xdg-open", "https://www.google.com").Start()
}

func launchURL(ctx context.Context, urlStr string) error {
	if runtime.GOOS == constants.OSWindows {
		return launchURLWindows(ctx, urlStr)
	}
	if runtime.GOOS == "darwin" {
		return exec.CommandContext(ctx, "open", urlStr).Run()
	}
	return exec.CommandContext(ctx, "xdg-open", urlStr).Start()
}

func launchURLWindows(ctx context.Context, urlStr string) error {
	cmd := exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", urlStr)
	err := cmd.Start()
	if err == nil {
		return nil
	}
	return exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", "Start-Process", fmt.Sprintf("'%s'", urlStr)).Run()
}

func launchPath(ctx context.Context, pathStr string) error {
	if runtime.GOOS == constants.OSWindows {
		cmd := exec.CommandContext(ctx, "explorer.exe", pathStr)
		return cmd.Start()
	}
	if runtime.GOOS == "darwin" {
		return exec.CommandContext(ctx, "open", pathStr).Run()
	}
	return exec.CommandContext(ctx, "xdg-open", pathStr).Start()
}

func launchGeneric(ctx context.Context, target string) error {
	if runtime.GOOS == constants.OSWindows {
		return exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", "Start-Process", target).Run()
	}
	if runtime.GOOS == "darwin" {
		return exec.CommandContext(ctx, "open", target).Run()
	}
	return exec.CommandContext(ctx, "xdg-open", target).Start()
}
