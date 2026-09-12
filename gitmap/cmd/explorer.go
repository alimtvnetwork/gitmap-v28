// Package cmd — explorer.go: cross-platform graphical desktop file manager launcher.
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
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func runExplorerCmd(args []string) error {
	target, err := resolveExplorerTarget(args)
	if err != nil {
		return err
	}

	return openExplorerAtTarget(target)
}

func openExplorerAtTarget(target string) error {
	isDir, statErr := checkExplorerTargetType(target)
	if statErr != nil {
		return statErr
	}

	if err := launchNativeExplorer(target, isDir); err != nil {
		return err
	}

	printExplorerSuccess(target)

	return nil
}

func printExplorerSuccess(target string) {
	fmt.Printf("  %s📂 Opened file explorer at: %s%s%s\n",
		constants.ColorCyan, constants.ColorWhite, target, constants.ColorReset)
}

func resolveExplorerTarget(args []string) (string, error) {
	raw := extractTargetArg(args)
	expanded := macro.ExpandPathAndEnv(raw)
	absPath, err := filepath.Abs(expanded)
	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf("resolve path %q", raw))
	}

	return absPath, nil
}

func extractTargetArg(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && a != "" {
			return a
		}
	}

	return "."
}

func checkExplorerTargetType(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, apperror.WrapSimple(err, fmt.Sprintf("path does not exist: %s", path))
	}

	return info.IsDir(), nil
}

func launchNativeExplorer(path string, isDir bool) error {
	cmd := buildExplorerCmd(path, isDir)
	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("launch file explorer for %q", path))
	}

	return nil
}

func buildExplorerCmd(path string, isDir bool) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return buildWindowsExplorerCmd(path, isDir)
	}

	if runtime.GOOS == "darwin" {
		return buildDarwinExplorerCmd(path, isDir)
	}

	return buildLinuxExplorerCmd(path, isDir)
}

func buildWindowsExplorerCmd(path string, isDir bool) *exec.Cmd {
	clean := filepath.Clean(path)
	if !isDir {
		return exec.Command("explorer.exe", "/select,"+clean)
	}

	return exec.Command("explorer.exe", clean)
}

func buildDarwinExplorerCmd(path string, isDir bool) *exec.Cmd {
	if !isDir {
		return exec.Command("open", "-R", path)
	}

	return exec.Command("open", path)
}

func buildLinuxExplorerCmd(path string, isDir bool) *exec.Cmd {
	target := path
	if !isDir {
		target = filepath.Dir(path)
	}

	return exec.Command("xdg-open", target)
}
