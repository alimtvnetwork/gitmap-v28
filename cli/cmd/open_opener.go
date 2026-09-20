// Package cmd — open_opener.go: native OS opener command execution.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func isHeadlessLinux() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	hasEnvDisplay := os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
	if hasEnvDisplay {
		return false
	}
	hasSocket := hasLocalX11Socket()
	if hasSocket {
		return false
	}
	return true
}

func hasLocalX11Socket() bool {
	fi, err := os.Stat("/tmp/.X11-unix/X0")
	return err == nil && fi != nil
}

func configureLinuxDisplay(cmd *exec.Cmd) {
	if runtime.GOOS != "linux" {
		return
	}
	hasEnvDisplay := os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
	if hasEnvDisplay {
		return
	}
	if hasLocalX11Socket() {
		cmd.Env = append(os.Environ(), "DISPLAY=:0")
		attachUserXAuthority(cmd)
	}
}

func attachUserXAuthority(cmd *exec.Cmd) {
	if os.Getenv("XAUTHORITY") != "" {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	xauth := filepath.Join(home, ".Xauthority")
	if fi, statErr := os.Stat(xauth); statErr == nil && fi != nil {
		cmd.Env = append(cmd.Env, "XAUTHORITY="+xauth)
	}
}

func handleHeadlessOpen(target string) error {
	fmt.Printf("No graphical display ($DISPLAY) detected.\nTarget: %s\n", target)
	return nil
}

func isDisplayError(msg string) bool {
	low := strings.ToLower(msg)
	hasDisplayErr := strings.Contains(low, "missing x server") ||
		strings.Contains(low, "no display") ||
		strings.Contains(low, "cannot open display") ||
		strings.Contains(low, "failed to initialize")
	return hasDisplayErr
}

func tryStartBrowserCandidate(cand, target string) bool {
	path, err := exec.LookPath(cand)
	if err != nil {
		return false
	}
	cmd := exec.Command(path, target)
	configureLinuxDisplay(cmd)
	configureDetachedProcess(cmd)
	startErr := cmd.Start()
	return startErr == nil
}

func tryDirectBrowser(target string) error {
	candidates := []string{"google-chrome", "chromium-browser", "chromium", "firefox", "sensible-browser"}
	for _, cand := range candidates {
		isStarted := tryStartBrowserCandidate(cand, target)
		if isStarted {
			return nil
		}
	}
	return handleHeadlessOpen(target)
}

func handleOpenerFailure(target, errStr string) error {
	if isDisplayError(errStr) {
		return handleHeadlessOpen(target)
	}
	if isWebTarget(target) {
		return tryDirectBrowser(target)
	}
	return handleHeadlessOpen(target)
}

func launchNativeOpener(target string) error {
	if isHeadlessLinux() {
		return handleHeadlessOpen(target)
	}

	cmd := buildOpenerCommand(target)
	configureLinuxDisplay(cmd)
	configureDetachedProcess(cmd)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return handleOpenerFailure(target, err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return handleOpenerFailure(target, errBuf.String())
	}

	return nil
}

func buildOpenerCommand(target string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	}

	if runtime.GOOS == "darwin" {
		return exec.Command("open", target)
	}

	return exec.Command("xdg-open", target)
}
