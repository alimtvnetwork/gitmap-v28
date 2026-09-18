# 07 — Cleanup

## Purpose

Define how temporary artifacts from the update process are identified
and removed after a successful update.


## Flow Diagram

See [`images/cleanup-flow.mmd`](images/cleanup-flow.mmd)

---

## Artifacts That Need Cleanup

| Artifact | Created By | Location |
|----------|------------|----------|
| `<binary>-update-<pid>.exe` | Handoff copy (Step 1) | Same dir as binary or temp dir |
| `<binary>-update-*.ps1` | Generated PowerShell script | System temp dir |
| `<binary>.exe.old` | Rename-first deploy | Deploy directory |
| `<binary>-updater-tmp-<pid>.exe` | Standalone updater handoff | Same dir as updater |

---

## Cleanup Command

Provide an explicit cleanup command:

```
<binary> update-cleanup
```

This command scans the binary's directory and temp directory for
leftover artifacts and removes them.

### Implementation

```go
func runUpdateCleanup() {
    // 1. Clean up .old backups in the binary's directory
    selfDir := resolveInstalledDir()
    cleanGlob(selfDir, "*.old")

    // 2. Clean up handoff copies
    cleanGlob(selfDir, "<binary>-update-*")
    cleanGlob(os.TempDir(), "<binary>-update-*")

    // 3. Clean up generated scripts
    cleanGlob(os.TempDir(), "<binary>-update-*.ps1")
}

func cleanGlob(dir, pattern string) {
    matches, err := filepath.Glob(filepath.Join(dir, pattern))
    if err != nil {
        return
    }

    for _, match := range matches {
        if err := os.Remove(match); err != nil {
            fmt.Printf("  !! Could not remove %s: %v\n", filepath.Base(match), err)
            continue
        }
        fmt.Printf("  OK Removed %s\n", filepath.Base(match))
    }
}
```

---

## Automatic Cleanup (MANDATORY)

The update process **MUST** invoke `update-cleanup` automatically at
the end of every successful update. Cleanup is not optional and not
deferred to the user — it runs as the final step of the update cycle.

| Rule | Rationale |
|------|-----------|
| **MUST** auto-invoke `update-cleanup` after version verify passes | User never has to remember to clean up |
| **MUST** use the *new* (deployed) binary to run the cleanup | Old binary may be missing the cleanup subcommand |
| **MUST** ignore locked-file errors silently | Worker copy may still be alive — next run cleans it |
| **MUST NOT** fail the update because cleanup failed | Cleanup is best-effort, never blocking |

```

# At the very end of the update worker, after version verify:

if newBinaryExists and versionVerifyPassed:
    run(newBinaryPath, "update-cleanup")  # best-effort, ignore exit code
```

```go
func attemptAutoCleanup() {
    // Best-effort — ignore errors
    matches, _ := filepath.Glob(filepath.Join(selfDir, "<binary>-update-*"))
    for _, m := range matches {
        os.Remove(m) // Ignore error — may be locked
    }
}
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
| [11-build-deploy.md](../13-generic-cli/11-build-deploy.md) | `.old` artifact lifecycle during deploy |
| [11-self-update-mechanism.md](../12-cicd-pipeline-workflows/11-self-update-mechanism.md) | Post-update cleanup of handoff copies and `.old` files |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
