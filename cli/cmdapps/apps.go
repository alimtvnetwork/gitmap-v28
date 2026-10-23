package cmdapps

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ListApps scans and discovers installed applications across the system.
func ListApps(opts ListOptions) (AppListResponse, error) {
	apps := make([]InstalledApp, 0, 32)
	if runtime.GOOS == "linux" {
		apps = append(apps, loadLinuxApps(opts)...)
	} else if runtime.GOOS == "windows" {
		apps = append(apps, loadWindowsApps(opts)...)
	}

	filtered := filterApps(apps, opts)
	resp := AppListResponse{
		Success:  true,
		Count:    len(filtered),
		Platform: runtime.GOOS,
		Data:     filtered,
	}

	return resp, nil
}

func loadLinuxApps(opts ListOptions) []InstalledApp {
	apps, err := scanLinuxApps(opts)
	if err == nil {
		return apps
	}

	return nil
}

func loadWindowsApps(opts ListOptions) []InstalledApp {
	apps, err := scanWindowsApps(opts)
	if err == nil {
		return apps
	}

	return nil
}

// UninstallApp performs surgical application uninstallation and cache synchronization.
func UninstallApp(target string, opts UninstallOptions) (AppUninstallResponse, error) {
	startTime := time.Now()
	resp := AppUninstallResponse{
		Success:      false,
		AppID:        target,
		Purged:       opts.IsPurge,
		RemovedFiles: make([]string, 0),
		CachesReset:  make([]string, 0),
	}

	if target == "" {
		resp.Error = "target application identifier cannot be empty"
		return resp, apperror.NewValidationError(resp.Error)
	}

	if opts.IsDryRun {
		resp.Success = true
		resp.DurationMs = time.Since(startTime).Milliseconds()
		return resp, nil
	}

	var hasExecutionErr error
	switch runtime.GOOS {
	case "linux":
		resp, hasExecutionErr = executeLinuxUninstall(target, opts, resp)
	case "windows":
		resp, hasExecutionErr = executeWindowsUninstall(target, opts, resp)
	default:
		resp.Error = fmt.Sprintf("unsupported platform for uninstaller: %s", runtime.GOOS)
		return resp, apperror.NewExecutionError(resp.Error)
	}

	resp.DurationMs = time.Since(startTime).Milliseconds()
	if hasExecutionErr != nil {
		resp.Success = false
		resp.Error = hasExecutionErr.Error()
		return resp, hasExecutionErr
	}

	resp.Success = true
	return resp, nil
}

func matchesFilter(app InstalledApp, normalizedFilter string) bool {
	if normalizedFilter == "" {
		return true
	}
	hasIdMatch := strings.Contains(strings.ToLower(app.ID), normalizedFilter)
	hasNameMatch := strings.Contains(strings.ToLower(app.Name), normalizedFilter)
	hasPkgMatch := strings.Contains(strings.ToLower(app.Package), normalizedFilter)
	hasExecMatch := strings.Contains(strings.ToLower(app.Exec), normalizedFilter)

	return hasIdMatch || hasNameMatch || hasPkgMatch || hasExecMatch
}

func filterApps(apps []InstalledApp, opts ListOptions) []InstalledApp {
	filtered := make([]InstalledApp, 0, len(apps))
	normalizedFilter := strings.ToLower(strings.TrimSpace(opts.FilterStr))

	for _, app := range apps {
		if opts.IsSystem && app.Scope != ScopeSystem {
			continue
		}
		if opts.IsUser && app.Scope != ScopeUser {
			continue
		}
		if !matchesFilter(app, normalizedFilter) {
			continue
		}

		filtered = append(filtered, app)
	}

	return filtered
}
