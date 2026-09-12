//go:build windows

// Package cmdagy — agy_clean_cache_windows.go provides Windows-specific path and process discovery.
package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// DiscoverCacheTargets discovers all Antigravity cache directories on Windows.
func DiscoverCacheTargets(includeTemp bool) []AgyCacheTarget {
	var targets []AgyCacheTarget
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")
	homeDir, _ := os.UserHomeDir()

	paths := make([]struct {
		p    string
		desc string
	}, 0)

	if appData != "" {
		baseAgy := filepath.Join(appData, "Antigravity")
		paths = append(paths,
			struct{ p, desc string }{filepath.Join(baseAgy, "Cache"), "Chromium network & disk cache"},
			struct{ p, desc string }{filepath.Join(baseAgy, "Code Cache"), "Compiled JS code cache"},
			struct{ p, desc string }{filepath.Join(baseAgy, "GPUCache"), "GPU shader cache"},
			struct{ p, desc string }{filepath.Join(baseAgy, "blob_storage"), "Local blob storage cache"},
			struct{ p, desc string }{filepath.Join(baseAgy, "logs"), "Application runtime logs"},
		)

		matches, _ := filepath.Glob(filepath.Join(baseAgy, "Dawn*Cache"))
		for _, m := range matches {
			paths = append(paths, struct{ p, desc string }{m, "Dawn WebGPU shader cache"})
		}
	}

	if localAppData != "" {
		paths = append(paths, struct{ p, desc string }{
			filepath.Join(localAppData, "antigravity-updater"), "Auto-updater package cache",
		})
	}

	if homeDir != "" {
		geminiAgy := filepath.Join(homeDir, ".gemini", "antigravity")
		paths = append(paths,
			struct{ p, desc string }{filepath.Join(geminiAgy, "crashes"), "Crash dump logs"},
			struct{ p, desc string }{filepath.Join(geminiAgy, "scratch"), "Temporary scratch files"},
		)
	}

	if includeTemp {
		paths = append(paths, struct{ p, desc string }{os.TempDir(), "User temp directory"})
	}

	for _, item := range paths {
		target := buildTargetInfo(item.p, item.desc)
		targets = append(targets, target)
	}

	return targets
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

// DiscoverTargetProcesses finds running Antigravity and WebView2 processes via tasklist on Windows.
func DiscoverTargetProcesses() ([]AgyProcessInfo, error) {
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseTasklistCSV(string(out), os.Getpid()), nil
}

// TerminateProcesses kills the specified processes via taskkill on Windows.
func TerminateProcesses(procs []AgyProcessInfo) (int, []string) {
	var killed int
	var warnings []string

	for _, p := range procs {
		cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(p.PID))
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		runErr := cmd.Run()
		if runErr != nil {
			warnings = append(warnings, fmt.Sprintf("kill PID %d (%s): %v", p.PID, p.Name, runErr))
			continue
		}

		killed++
	}

	if killed > 0 {
		time.Sleep(600 * time.Millisecond)
	}

	return killed, warnings
}
