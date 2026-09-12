//go:build !windows

// Package cmdagy — agy_clean_cache_unix.go provides Unix (Linux/macOS) path and process discovery.
package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

// DiscoverCacheTargets discovers all Antigravity cache directories on Linux and macOS.
func DiscoverCacheTargets(includeTemp bool) []AgyCacheTarget {
	var targets []AgyCacheTarget
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		return targets
	}

	paths := make([]struct {
		p    string
		desc string
	}, 0)

	if runtime.GOOS == "darwin" {
		paths = append(paths, getDarwinPaths(homeDir)...)
	} else {
		paths = append(paths, getLinuxPaths(homeDir)...)
	}

	geminiAgy := filepath.Join(homeDir, ".gemini", "antigravity")
	paths = append(paths,
		struct{ p, desc string }{filepath.Join(geminiAgy, "crashes"), "Crash dump logs"},
		struct{ p, desc string }{filepath.Join(geminiAgy, "scratch"), "Temporary scratch files"},
	)

	if includeTemp {
		paths = append(paths, struct{ p, desc string }{os.TempDir(), "System temp directory"})
	}

	for _, item := range paths {
		targets = append(targets, buildTargetInfo(item.p, item.desc))
	}

	return targets
}

func getDarwinPaths(homeDir string) []struct{ p, desc string } {
	appSupport := filepath.Join(homeDir, "Library", "Application Support", "Antigravity")
	caches := filepath.Join(homeDir, "Library", "Caches")

	items := []struct{ p, desc string }{
		{filepath.Join(appSupport, "Cache"), "Chromium disk cache"},
		{filepath.Join(appSupport, "Code Cache"), "Compiled JS code cache"},
		{filepath.Join(appSupport, "GPUCache"), "GPU shader cache"},
		{filepath.Join(appSupport, "blob_storage"), "Blob storage cache"},
		{filepath.Join(appSupport, "logs"), "Application logs"},
		{filepath.Join(caches, "Antigravity"), "macOS app cache"},
		{filepath.Join(caches, "antigravity-updater"), "Auto-updater package cache"},
	}

	matches, _ := filepath.Glob(filepath.Join(appSupport, "Dawn*Cache"))
	for _, m := range matches {
		items = append(items, struct{ p, desc string }{m, "Dawn WebGPU shader cache"})
	}

	return items
}

func getLinuxPaths(homeDir string) []struct{ p, desc string } {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(homeDir, ".config")
	}

	baseAgy := filepath.Join(configDir, "Antigravity")
	items := []struct{ p, desc string }{
		{filepath.Join(baseAgy, "Cache"), "Chromium disk cache"},
		{filepath.Join(baseAgy, "Code Cache"), "Compiled JS code cache"},
		{filepath.Join(baseAgy, "GPUCache"), "GPU shader cache"},
		{filepath.Join(baseAgy, "blob_storage"), "Blob storage cache"},
		{filepath.Join(baseAgy, "logs"), "Application logs"},
		{filepath.Join(configDir, "antigravity-updater"), "Auto-updater package cache"},
	}

	matches, _ := filepath.Glob(filepath.Join(baseAgy, "Dawn*Cache"))
	for _, m := range matches {
		items = append(items, struct{ p, desc string }{m, "Dawn WebGPU shader cache"})
	}

	return items
}

func buildTargetInfo(p, desc string) AgyCacheTarget {
	info, err := os.Stat(p)
	exists := err == nil && info.IsDir()
	size := int64(0)
	files := 0

	if exists {
		size, files, _ = CalculateDirStats(p)
	}

	return AgyCacheTarget{
		Path:        p,
		Description: desc,
		Exists:      exists,
		SizeBytes:   size,
		FileCount:   files,
	}
}

func DiscoverTargetProcesses() ([]AgyProcessInfo, error) {
	cmd := exec.Command("ps", "-eo", "pid,comm")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parsePsOutput(string(out), os.Getpid()), nil
}

// TerminateProcesses kills the specified processes via SIGKILL on Unix.
func TerminateProcesses(procs []AgyProcessInfo) (int, []string) {
	var killed int
	var warnings []string

	for _, p := range procs {
		err := syscall.Kill(p.PID, syscall.SIGKILL)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("kill PID %d (%s): %v", p.PID, p.Name, err))
			continue
		}

		killed++
	}

	if killed > 0 {
		time.Sleep(600 * time.Millisecond)
	}

	return killed, warnings
}
