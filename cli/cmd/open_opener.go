// Package cmd — open_opener.go: native OS opener command execution.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func launchNativeOpener(target string) error {
	if isHeadlessLinux() {
		return handleHeadlessOpen(target)
	}

	cmd := buildOpenerCommand(target)
	configureLinuxDisplay(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return handleHeadlessOpen(target)
	}

	if err := cmd.Wait(); err != nil {
		return handleHeadlessOpen(target)
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
