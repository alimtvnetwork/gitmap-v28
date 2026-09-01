// Package cmd — chrome_exec.go: cross-platform Chrome launcher supporting
// single and multi-profile URL opening with profile mapping.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type chromeLaunchTarget struct {
	Profile string
	URLs    []string
}

type chromeLaunchOptions struct {
	Targets    []chromeLaunchTarget
	IsIncog    bool
	IsNewWin   bool
	AppURL     string
	ExtraFlags []string
}

func findChromeBinaryPath() (string, error) {
	if runtime.GOOS == "windows" {
		return findChromeWindows()
	}
	if runtime.GOOS == "darwin" {
		return findChromeDarwin()
	}
	return findChromeLinux()
}

func findChromeWindows() (string, error) {
	candidates := []string{
		filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	if p, err := exec.LookPath("chrome.exe"); err == nil {
		return p, nil
	}
	return "", apperror.NewSimple("chrome executable not found on Windows", "E4101")
}

func findChromeDarwin() (string, error) {
	paths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		filepath.Join(os.Getenv("HOME"), "Applications", "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", apperror.NewSimple("Google Chrome not found in /Applications", "E4102")
}

func findChromeLinux() (string, error) {
	for _, name := range constants.ChromeProcessLinuxNames {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	return "", apperror.NewSimple("Chrome/Chromium executable not found in PATH", "E4103")
}

func runChromeOpen(args []string) error {
	checkHelp(constants.SubCmdChromeOpen, args)
	opts, err := parseChromeLaunchArgs(args)
	if err != nil {
		return err
	}
	bin, err := findChromeBinaryPath()
	if err != nil {
		return err
	}
	return executeChromeLaunches(bin, opts)
}

func parseChromeLaunchArgs(args []string) (chromeLaunchOptions, error) {
	var opts chromeLaunchOptions
	var positional []string
	var explicitProfile string

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--incognito" || a == "-incognito" {
			opts.IsIncog = true
		} else if a == "--new-window" {
			opts.IsNewWin = true
		} else if strings.HasPrefix(a, "--app=") {
			opts.AppURL = strings.TrimPrefix(a, "--app=")
		} else if strings.HasPrefix(a, "--profile=") {
			explicitProfile = strings.TrimPrefix(a, "--profile=")
		} else if (a == "-p" || a == "--profile") && i+1 < len(args) {
			i++
			explicitProfile = args[i]
		} else if strings.HasPrefix(a, "-") {
			opts.ExtraFlags = append(opts.ExtraFlags, a)
		} else {
			positional = append(positional, a)
		}
	}
	opts.Targets = resolveLaunchTargets(positional, explicitProfile)
	return opts, nil
}

func resolveLaunchTargets(positional []string, explicitProf string) []chromeLaunchTarget {
	if len(positional) == 0 && explicitProf == "" {
		return []chromeLaunchTarget{{Profile: "Default", URLs: []string{"chrome://newtab"}}}
	}
	if isMultiSegmentMapping(positional) {
		return parseMultiSegmentMapping(positional[0])
	}
	if len(positional) >= 2 && isProfileNameOrDir(positional[0]) {
		prof := positional[0]
		return []chromeLaunchTarget{{Profile: prof, URLs: positional[1:]}}
	}
	prof := explicitProf
	if prof == "" {
		prof = "Default"
	}
	return []chromeLaunchTarget{{Profile: prof, URLs: positional}}
}

func isMultiSegmentMapping(positional []string) bool {
	if len(positional) == 0 {
		return false
	}
	first := positional[0]
	return strings.Contains(first, "=") && (strings.Contains(first, "http://") || strings.Contains(first, "https://") || strings.Contains(first, "chrome://"))
}

func parseMultiSegmentMapping(arg string) []chromeLaunchTarget {
	var targets []chromeLaunchTarget
	segments := strings.Split(arg, ",")
	for _, seg := range segments {
		parts := strings.SplitN(seg, "=", 2)
		if len(parts) == 2 {
			targets = append(targets, chromeLaunchTarget{
				Profile: strings.TrimSpace(parts[0]),
				URLs:    []string{strings.TrimSpace(parts[1])},
			})
		}
	}
	return targets
}

func isProfileNameOrDir(name string) bool {
	_, hasDir := resolveChromeProfileDir(name)
	return hasDir
}

func executeChromeLaunches(bin string, opts chromeLaunchOptions) error {
	for _, tgt := range opts.Targets {
		if err := launchSingleTarget(bin, tgt, opts); err != nil {
			return err
		}
	}
	return nil
}

func launchSingleTarget(bin string, tgt chromeLaunchTarget, opts chromeLaunchOptions) error {
	dirName := tgt.Profile
	if resolved, hasDir := resolveChromeProfileDir(tgt.Profile); hasDir {
		dirName = filepath.Base(resolved)
	}
	cmdArgs := buildChromeCmdArgs(dirName, tgt.URLs, opts)
	cmd := exec.Command(bin, cmdArgs...)
	configureDetachedProcess(cmd)

	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("launch chrome profile %s", dirName))
	}
	displayName := chromeProfileDisplayName(dirName)
	urlsJoined := strings.Join(tgt.URLs, ", ")
	fmt.Printf("\033[1;92m✓ launched\033[0m  Chrome [\033[1m%s\033[0m / %q] → %s\n", dirName, displayName, urlsJoined)
	return nil
}

func buildChromeCmdArgs(dirName string, urls []string, opts chromeLaunchOptions) []string {
	var args []string
	if dirName != "" {
		args = append(args, fmt.Sprintf("--profile-directory=%s", dirName))
	}
	if opts.IsIncog {
		args = append(args, "--incognito")
	}
	if opts.IsNewWin {
		args = append(args, "--new-window")
	}
	if opts.AppURL != "" {
		args = append(args, fmt.Sprintf("--app=%s", opts.AppURL))
	}
	args = append(args, opts.ExtraFlags...)
	args = append(args, urls...)
	return args
}
