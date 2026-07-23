# Self-Update Mechanism

## Overview

A reusable pattern for CLI tools that update themselves from source,
solving the Windows file-lock problem where a running binary cannot
overwrite itself. This guide is framework-agnostic and applies to
any compiled CLI tool (Go, Rust, C#, etc.) deployed on Windows.

## The Problem

On Windows, a running `.exe` holds a file lock on its own binary.
If the update process tries to overwrite the file while the original
process is still running, the OS blocks the operation:

> "The process cannot access the file because it is being used by
> another process."

This does not occur on Linux/macOS, where a running binary can be
replaced on disk (the OS keeps the old inode until the process exits).

## Solution: Copy-and-Handoff

A three-layer approach that reliably bypasses file locks:

### Layer 1 — Copy and Re-launch

1. The running binary copies itself to a temp location:
   `%TEMP%\<toolname>-update-<pid>.exe` (or same directory as the active binary).
2. It launches the temp copy with a hidden worker command (e.g. `update-runner`)
   to indicate it's the delegated updater.
3. The parent **waits for the worker** using foreground/blocking execution.
   This keeps the terminal session stable. The handoff copy is a different
   file so there is no lock conflict with deploy — rename-first handles it.

```
# Pseudocode — applies to any compiled language
func runUpdate():
    tempPath = copyBinaryToTemp()
    runForeground(tempPath, ["update-runner"])  # Blocking — keeps terminal stable
    # Parent exits naturally after worker completes

func runUpdateRunner():
    # This is the worker — runs from the temp copy
    executeUpdate(repoPath)
```

### Layer 2 — Skip-if-Current + Delayed Rebuild

The temp copy orchestrates the actual update. Before rebuilding,
it checks whether there are any changes to apply:

1. **Capture current version** from the deployed binary.
2. **Pull latest source** and inspect the output.
3. **If "Already up to date"** — print a message and exit early.
   No rebuild, no deploy, no wasted time.
4. **Otherwise, wait 1–2 seconds** for the parent process to fully
   terminate and release all OS-level file handles.
5. **Run the build pipeline** (resolve deps, build, deploy).
6. **Compare versions** — warn if unchanged (version constant not
   bumped), confirm if different.

```
# Example: generated update script
$oldVersion = & $deployedBinary version 2>&1

$pullOutput = git pull 2>&1
if ($pullOutput -match "Already up to date") {
    print("No update needed — already running latest version")
    exit(0)
}

Start-Sleep -Seconds 1.2
& build-script -NoPull   # Already pulled above

$newVersion = & $newBinary version 2>&1
if ($oldVersion == $newVersion) {
    warn("Version unchanged — was the version constant bumped?")
} else {
    print("Updated: $oldVersion -> $newVersion")
}
```

### Layer 3 — Deploy with Rollback

The build pipeline's deploy step includes rollback safety:

1. **Backup** the existing binary before overwriting:
   `toolname.exe` → `toolname.exe.old`
2. **Attempt file copy** with retry loop (15–20 attempts × 500ms)
3. **On success** — leave the `.old` file in place for manual cleanup
4. **On failure after all retries** — restore the `.old` backup so the
   user still has a working binary

```
# Pseudocode — deploy with rollback
backup = destination + ".old"
if fileExists(destination):
    copy(destination, backup)

attempts = 0
success = false
while attempts < maxAttempts:
    try:
        copyFile(source, destination)
        success = true
        break
    catch fileLocked:
        log("Target in use, retrying...")
        sleep(500ms)
        attempts++

if success:
    log("Previous binary kept as .old (run cleanup command to remove)")
else:
    # Restore working binary from backup
    copy(backup, destination)
    fail("Deploy failed — previous version restored")
```

## Flow Diagram

```
User runs: <tool> update
   │
   ├─ Parent copies self → <tool>-update-<pid>.exe (same dir, fallback %TEMP%)
   ├─ Parent launches copy with: update-runner (foreground/blocking)
   ├─ Parent waits for worker to complete (terminal stays attached)
   │
   └─ Worker (update-runner) starts
      ├─ Captures current deployed version
      ├─ Runs: run.ps1 -Update (pull → build → deploy)
      │    ├─ Backs up existing binary as .old
      │    ├─ Deploys new binary (with retry)
      │    ├─ PATH sync: rename-first (.old), then copy new
      │    └─ On failure: restores .old backup
      ├─ Compares old vs new version
      ├─ Runs: <tool> changelog --latest
      ├─ Runs: <tool> update-cleanup (auto)
      │    ├─ Removes <tool>-update-*.exe
      │    └─ Removes *.old from deploy directory
      └─ Cleans up temp script
```

## Cleanup Command

Provide an explicit cleanup subcommand (e.g. `tool update-cleanup`)
that removes update artifacts:

1. **Temp update copies** — `%TEMP%\<tool>-update-*.exe` files left
   from previous update handoffs
2. **Old backup binaries** — `*.old` files in the deploy directory
   from rollback backups

### Auto-Cleanup at End of Update

The update process automatically invokes the cleanup command after
a successful build and deploy. This means the user doesn't need to
remember to clean up — it happens as part of the normal update cycle.

```
# At end of update script (after version comparison):
if newBinaryExists:
    run(newBinary, "update-cleanup")
```

### Manual Cleanup

The command is also available for ad-hoc use. This is useful when:
- A previous update was interrupted before cleanup ran
- The user wants to verify what artifacts exist
- Debugging update issues

```
# Pseudocode — cleanup command
func runUpdateCleanup():
    # Clean temp copies
    for file in glob("%TEMP%/<tool>-update-*.exe"):
        if file != currentExecutable:
            delete(file)
            print("Removed temp copy: " + basename(file))

    # Clean .old backups from deploy directory
    for file in glob(deployDir + "/*.old"):
        delete(file)
        print("Removed backup: " + basename(file))

    print("Cleanup complete")
```

### Why a Separate Command (Not Auto-Cleanup on Startup)?

- **Transparency** — the user sees exactly what's being deleted
- **Safety** — `.old` files serve as manual rollback if the new
  version has issues (user can rename `.old` back before cleanup runs)
- **Performance** — no filesystem scanning on every startup
- **Explicitness** — follows the principle of least surprise

## Prerequisites

- The **source repo path** must be available to the binary at runtime.
  Common approaches:
  - Embedded at build time via linker flags (e.g. Go `-ldflags`)
  - Stored in a config file next to the binary
  - Resolved from an environment variable
- A **build script** must exist at the known repo path.
- If the repo path is missing or invalid, the update command should
  print a clear error and exit.

## Optional Enhancements

Beyond the core pattern, consider these additional improvements:

### Checksum Verification

For tools that download pre-built binaries (rather than building
from source), verify integrity after download:

```
expected = readHashFile(buildOutput + ".sha256")
actual = sha256(newBinary)
if expected != actual:
    fail("Binary checksum mismatch — download may be corrupted")
```

For source-built tools, the version comparison serves a similar
validation purpose — if the binary runs and reports a version,
it's structurally valid.

### Exit Code Propagation

The temp copy should propagate the build script's exit code so
the calling process (or CI) can detect failures:

```
result = runBuildScript(scriptPath)
exit(result.exitCode)
```

### Alternative: Rename Trick (Windows)

Windows blocks *overwriting* a running `.exe` but allows *renaming*
it. This enables a simpler deploy strategy:

```
# Instead of retry-on-lock:
rename("tool.exe", "tool.exe.old")   # Works even while running
copy(newBinary, "tool.exe")          # No lock conflict
# On next startup: delete("tool.exe.old")
```

This is more deterministic than retry loops but requires the deploy
step to run from a different process than the one holding the lock.

## Error Diagnostics

When version verification fails after an update, the script must
print trace-level diagnostic output so the root cause is immediately
visible without manual inspection:

```
  Version before:   gitmap-v27 v2.66.0
  Version active:   gitmap-v27 v2.67.0
  Version deployed: unknown
  Active binary:    C:\Users\user\bin\gitmap.exe
  Deployed binary:  (not resolved)

  [WARN] Deployed binary could not be verified (not resolved or missing).
  [TRACE] activeAfter=gitmap-v27 v2.67.0  deployedAfter=unknown
  [HINT] Check that powershell.json 'deployPath' points to the correct directory
         and that the binary exists at: <path>
  [OK] Active PATH binary updated successfully: gitmap-v27 v2.67.0
```

### Verification Logic

The version check uses a three-branch decision:

| Condition | Result | Exit Code |
|-----------|--------|-----------|
| Active updated, deployed unknown | **Warning** — active PATH binary is valid, deployed path misconfigured | 0 (success) |
| Active unknown, or active ≠ deployed | **Failure** — real version mismatch or PATH not working | 1 (error) |
| Active = deployed | **Success** — both binaries match | 0 (success) |

This ensures that a missing or misconfigured `deployPath` in
`powershell.json` does not block an otherwise successful update.
The user sees a clear warning but the update completes.

### Required Trace Points

| Trace | When |
|-------|------|
| `deployedBinary: not resolved` | `$deployedBinary` is `$null` (config missing or `deployPath` unset) |
| `deployedBinary: path not found: <path>` | Config resolved but file doesn't exist at that path |
| `Get-Command gitmap-v27: not found in PATH` | Active binary not discoverable via PATH |
| `activeAfter=... deployedAfter=...` | Always printed on warning or mismatch |

## Error Handling

| Scenario | Behavior |
|----------|----------|
| No repo path configured | Print error, exit 1 |
| Repo path doesn't exist | Print error, exit 1 |
| Already up to date | Print message, exit 0 (no rebuild) |
| Build/compile fails | Script exits with error, backup remains |
| Deploy locked after retries | Restore backup, fail with clear message |
| Temp copy fails to launch | Print error, exit 1 |
| Version unchanged after update | Warn user (version constant not bumped) |
| Deployed binary unknown, active OK | **Warning** (not failure), exit 0 |
| Active binary unknown | Print trace, exit 1 |

## Platform Considerations

| Platform | File Lock Behavior | Self-Update Approach |
|----------|--------------------|----------------------|
| Windows | Binary locked while running | Copy-and-handoff (this pattern) |
| Linux | Binary replaceable on disk | Direct overwrite (simpler) |
| macOS | Binary replaceable on disk | Direct overwrite (simpler) |

On Linux/macOS, the copy-and-handoff pattern still works but is
unnecessary — a simple in-place replace suffices. Consider
platform-detecting which path to take if cross-platform support
is needed.

## Key Learnings

1. **`exit()` doesn't release locks instantly** — the OS may hold
   file handles briefly after process termination. Always add a
   delay before attempting to overwrite.
2. **Always add retry logic** for file operations on deployed
   binaries — even with the handoff, a small timing window exists.
3. **The handoff copy is a different file** — the parent's lock on
   the original binary does not conflict with the worker. Rename-first
   in `run.ps1` handles the locked active PATH binary.
4. **The parent MUST use `cmd.Run()` (foreground/blocking)** — keeps
   the terminal session stable. NEVER use `cmd.Start()` + `os.Exit(0)`
   which detaches and breaks the command line.
5. **Use rename-first, not copy-first** for PATH sync — Windows
   blocks overwrite of a running `.exe` but allows renaming it.
6. **Bump the version on every change** so the user can confirm
   the update actually applied.
7. **Always provide a rollback path** — if the update fails
   mid-deploy, the user should still have a working binary.
8. **Skip unnecessary rebuilds** — check for source changes before
   rebuilding to save time and avoid confusion.
9. **Compare versions before and after** — catch cases where the
   source changed but the version constant wasn't bumped.
10. **Log verbosely during update** — self-update failures are hard
    to debug without detailed logs of each step.
11. **Prefer explicit cleanup commands over auto-cleanup** — let
    the user decide when to remove artifacts; `.old` files serve as
    a manual rollback option until explicitly cleaned up.
12. **Never add `Read-Host` or interactive prompts** to generated
    scripts — they run in non-interactive PowerShell sessions.

## Cross-References (Generic Specifications)

This document is an application-level summary. The following generic,
tool-agnostic specs provide detailed breakdowns of each mechanism:

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Overall architecture | [01-self-update-overview.md](../08-generic-update/01-self-update-overview.md) | Platform behavior table, two update strategies (source vs binary), repo resolution tiers, version comparison |
| Deploy path resolution | [02-deploy-path-resolution.md](../08-generic-update/02-deploy-path-resolution.md) | 3-tier deploy target resolution (CLI flag → PATH → config) |
| Rename-first deploy | [03-rename-first-deploy.md](../08-generic-update/03-rename-first-deploy.md) | Full PowerShell + Bash implementations, rollback, PATH sync, retry reduction (20→5) |
| Build scripts | [04-build-scripts.md](../08-generic-update/04-build-scripts.md) | `run.ps1` / `run.sh` pipeline (pull → deps → build → deploy), config loading, ldflags |
| Handoff mechanism | [05-handoff-mechanism.md](../08-generic-update/05-handoff-mechanism.md) | Copy-and-handoff flow, worker launch, UTF-8 BOM, binary-based handoff (standalone updater) |
| Cleanup | [06-cleanup.md](../08-generic-update/06-cleanup.md) | Artifact inventory, `update-cleanup` command, `.old` lifecycle, auto vs manual cleanup |
| Repo path sync | [08-repo-path-sync.md](../08-generic-update/08-repo-path-sync.md) | Post-deploy DB sync of source repo path via `set-source-repo` command |

### Mapping: This Document → Generic Specs

| Section Here | Generic Equivalent |
|-------------|-------------------|
| Layer 1 (Copy and Re-launch) | `05-handoff-mechanism.md` §Solution: Copy-and-Handoff |
| Layer 2 (Skip-if-Current) | `01-self-update-overview.md` §Version Comparison |
| Layer 3 (Deploy with Rollback) | `03-rename-first-deploy.md` §Rollback |
| Cleanup Command | `06-cleanup.md` §Cleanup Command |
| Alternative: Rename Trick | `03-rename-first-deploy.md` §The Solution: Rename-First |
| Platform Considerations | `01-self-update-overview.md` §Platform Behavior |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
