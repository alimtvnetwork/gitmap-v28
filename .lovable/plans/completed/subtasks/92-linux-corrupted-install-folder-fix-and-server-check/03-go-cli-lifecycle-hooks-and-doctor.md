# Subtask 03: Go CLI Lifecycle Hooks & Doctor Integration

## Objective
Wire automated detection and cleanup into the core Gitmap CLI runtime:
1. In `gitmap/cmd/root.go`:
   - Add background/startup hook in `Run()`:
     - On Linux/macOS (`runtime.GOOS != "windows"`), perform a fast silent scan for corrupted directories.
     - Skip when running `version` / `-v` commands so scripted version checks remain completely clean.
     - Silently cleans empty corrupted directories or logs an informative notice to `os.Stderr`.
2. In `gitmap/cmd/selfinstall.go`:
   - In `runSelfInstallWorkflow()`, execute `CleanCorruptedDirs()` before staging scripts and immediately after completion.
   - Ensure `promptInstallDir()` validates user input against control characters and newlines.
3. In `gitmap/cmd/doctor.go` & `gitmap/cmd/doctor_run.go`:
   - Add `probeCorruptedInstallDirs()` health check.
   - Reports `[ok]` if no corrupted directories exist.
   - Reports `[fail]` if corrupted directories exist with remediation instructions.
   - Supports auto-remediation when running `gitmap doctor --fix`.

## Files Affected
- `gitmap/cmd/root.go`
- `gitmap/cmd/selfinstall.go`
- `gitmap/cmd/doctor.go`
- `gitmap/cmd/doctor_run.go`
