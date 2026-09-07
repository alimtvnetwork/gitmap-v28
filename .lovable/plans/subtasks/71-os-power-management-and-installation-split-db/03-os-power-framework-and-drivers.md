# Subtask 03: OS Power Management Framework & Pluggable Drivers

## Status
Completed

## Context & Objectives
Implement the pluggable, DRY power management framework in `gitmap/power/`:
1. **Interface Definition** (`gitmap/power/manager.go`):
   - `type Manager interface` with `Platform()`, `GetStatus()`, `SetNeverSleep()`, `SetTimeouts()`, `ApplySettings()`.
   - `NewManager() (Manager, error)` dynamically returns the OS driver based on `runtime.GOOS`.
2. **Data Types** (`gitmap/power/types.go`):
   - `Settings` struct: `Platform`, `DisplayTimeoutMinutes`, `SleepTimeoutMinutes`, `DiskTimeoutMinutes`, `IsNeverSleep`, `IsLockDisabled`, `IsActive`.
3. **Windows Driver** (`gitmap/power/manager_windows.go`):
   - Uses `powercfg.exe` (`/change monitor-timeout-ac`, `/change standby-timeout-ac`, query active schemes).
4. **Linux Driver** (`gitmap/power/manager_linux.go`):
   - 3-tier fallback: GNOME `gsettings` (`org.gnome.desktop.session idle-delay`, `sleep-inactive-ac-timeout`, `lock-enabled`), X11 `xset s off -dpms`, and `systemd-logind`.
5. **Darwin & Fallback Drivers** (`gitmap/power/manager_darwin.go`, `manager_fallback.go`):
   - Stubs for future macOS `pmset` and unsupported OSes.

## Verification Steps
- Unit tests in `gitmap/power/manager_test.go` verify factory instantiation and settings parsing.
