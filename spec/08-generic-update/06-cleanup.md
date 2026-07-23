# 06 — Cleanup

## Purpose

Define how temporary artifacts from the update process are identified
and removed after a successful update.


## Flow Diagram

See [`images/cleanup-flow.mmd`](images/cleanup-flow.mmd)

---

## Artifacts That Need Cleanup

| Artifact | Created By | Location | Lock-prone on Windows |
|----------|------------|----------|------------------------|
| `<binary>-update-<pid>.exe` | Handoff copy (Phase 1) | Same dir as binary or temp dir | YES — held by the handoff process |
| `<binary>-update-*.ps1` | Generated PowerShell script | System temp dir | NO |
| `<binary>.exe.old` | Rename-first deploy (Phase 2) | Deploy directory | YES — briefly held by AV/Explorer |
| `<binary>-updater-tmp-<pid>.exe` | Standalone updater handoff | Same dir as updater | YES |

---

## Why Cleanup CANNOT Run Inside the Handoff Process

A naive design would call `update-cleanup` at the end of the Phase 2
worker (the running handoff copy). This fails on Windows:

```
Cleaning up update artifacts...
Error: could not remove cleanup artifact at .\gitmap-update-4692.exe: Access is denied.
Error: could not remove cleanup artifact at .\gitmap.exe.old: Access is denied.
```

The handoff process holds an exclusive lock on its own binary
(`gitmap-update-<pid>.exe`), so `os.Remove` returns `Access is denied`.
The `.old` backup is also frequently locked at this exact moment by AV
scanners or Explorer indexing, because it was just renamed.

The fix is **Phase 3** — handing off cleanup to the freshly deployed
binary (a different file with no shared lock).

---

## Phase 3 Cleanup Handoff

After Phase 2 completes, the still-running handoff copy spawns the
freshly deployed binary detached, with a small delay so its own
process can exit and release the file lock first.

### Windows

```go
func scheduleDeployedCleanupHandoff() {
    deployed := resolveDeployedBinaryPath()
    if deployed == "" {
        return
    }
    cmdLine := fmt.Sprintf(
        `ping 127.0.0.1 -n 3 >nul & start "" /B "%s" update-cleanup`,
        deployed,
    )
    cmd := exec.Command("cmd.exe", "/C", cmdLine)
    _ = cmd.Start() // detached, fire-and-forget
}
```

- `ping 127.0.0.1 -n 3 >nul` sleeps ~2s using only built-in `cmd`.
- `start "" /B` detaches without opening a new window.
- The handoff process exits immediately after `cmd.Start()`. By the
  time the deployed binary's `update-cleanup` actually runs, the
  handoff file is unlocked.

### Unix

No lock conflicts exist — the deployed binary's `update-cleanup` is
invoked inline (or skipped entirely if the handoff copy IS the
deployed binary).

```go
cmd := exec.Command(deployed, "update-cleanup")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
_ = cmd.Run()
```

---

## Build Script Cleanup

The build script (`run.ps1` / `run.sh`) should call `update-cleanup`
at the end of an update run:

### PowerShell

```powershell
if ($Update) {
    Write-Info "Running update cleanup"
    & $binaryPath update-cleanup
}
```

### Bash

```bash
if [[ "$UPDATE" == "true" ]]; then
    echo "  -> Running update cleanup"
    "$BINARY_PATH" update-cleanup || true
fi
```

---

## `.old` File Lifecycle

The `.old` backup files from rename-first deploy follow this lifecycle:

```
Deploy:
  <binary>.exe → <binary>.exe.old    (renamed)
  new binary   → <binary>.exe        (copied)

After successful update:
  <binary>.exe.old → deleted          (by update-cleanup)

After failed update:
  <binary>.exe.old → <binary>.exe    (rollback — rename back)
```

**Important**: Never delete `.old` files during deploy. They are the
rollback safety net. Only delete them after the update is confirmed
successful (version check passes).

---

## Temp Directory Hygiene

The system temp directory can accumulate artifacts if updates are
interrupted. The cleanup command should scan for patterns:

| Pattern | Artifact |
|---------|----------|
| `<binary>-update-*.exe` | Handoff copies |
| `<binary>-update-*.ps1` | Generated scripts |
| `<binary>-install-*.ps1` | Downloaded install scripts |

Use conservative matching — only delete files that match the exact
prefix pattern for your tool.

---

## When to Run Cleanup

| Trigger | Automatic | Manual |
|---------|-----------|--------|
| After successful update | ✅ Best-effort | |
| User runs `update-cleanup` | | ✅ Full scan |
| Before starting a new update | ✅ Clean old copies | |
| On `doctor` command | ✅ Warn about stale files | |

---

## Constraints

- Cleanup must never fail the parent operation (update, build, etc.).
- Use `os.Remove()`, not recursive deletion — only target specific files.
- Match patterns conservatively to avoid deleting unrelated files.
- Never delete the active binary or its data directory.
- Log every file removed so the user has visibility.
- On Windows, some files may be locked — skip them without error.

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | `.old` artifact lifecycle during deploy |
| [03-self-update-mechanism.md](../03-general/03-self-update-mechanism.md) | Post-update cleanup of handoff and `.old` files |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
