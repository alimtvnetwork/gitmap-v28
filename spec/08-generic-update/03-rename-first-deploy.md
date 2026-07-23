# 03 — Rename-First Deploy Strategy

## Purpose

Define the file-replacement strategy that works on all platforms,
including Windows where running executables cannot be overwritten.

---

## Flow Diagram

See [`images/rename-first-deploy-flow.mmd`](images/rename-first-deploy-flow.mmd)

---

## The Problem

On Windows, copying a file over a running `.exe` fails with:

```
The process cannot access the file '<binary>.exe' because it is being
used by another process.
```

This happens because:
- The OS holds a file lock on any running executable.
- `Copy-Item` / `cp` tries to open the destination for writing, which
  the lock prevents.

---

## The Solution: Rename-First

Windows allows **renaming** a running executable (the OS tracks the
process by handle, not by filename). So instead of overwriting:

```
Step 1: Rename existing <binary>.exe → <binary>.exe.old
Step 2: Copy new binary → <binary>.exe (destination is now free)
```

After the rename, the original process continues running from the
renamed file. The new binary occupies the original filename, ready
for the next invocation.

---

## Implementation

### PowerShell

```powershell
$destFile = Join-Path $appDir $Config.binaryName
$backupFile = "$destFile.old"
$hasBackup = $false
$deploySuccess = $false

if (Test-Path $destFile) {
    # Rename-first: move the locked binary out of the way
    try {
        if (Test-Path $backupFile) {
            Remove-Item $backupFile -Force -ErrorAction SilentlyContinue
        }
        Rename-Item $destFile $backupFile -Force -ErrorAction Stop
        $hasBackup = $true
        Write-Info "Renamed existing binary to $($Config.binaryName).old"
    } catch {
        Write-Warn "Rename-first failed: $_"
        # Fallback: try a copy-based backup
        try {
            Copy-Item $destFile $backupFile -Force -ErrorAction Stop
            $hasBackup = $true
        } catch {
            Write-Warn "Could not create backup: $_"
        }
    }
}

# Copy new binary — destination is free after rename
$maxAttempts = 5
$attempt = 1
while ($true) {
    try {
        Copy-Item $BinaryPath $destFile -Force -ErrorAction Stop
        $deploySuccess = $true
        break
    } catch {
        if ($attempt -ge $maxAttempts) { throw }
        Write-Warn "Target still locked; retrying ($attempt/$maxAttempts)..."
        Start-Sleep -Milliseconds 500
        $attempt++
    }
}
```

### Bash

```bash
dest_file="$app_dir/$BINARY_NAME"
backup_file="${dest_file}.old"
has_backup=false
deploy_success=false

if [[ -f "$dest_file" ]]; then
    # Rename-first: move the existing binary out of the way
    rm -f "$backup_file" 2>/dev/null || true
    if mv "$dest_file" "$backup_file" 2>/dev/null; then
        has_backup=true
        echo "  -> Renamed existing binary to ${BINARY_NAME}.old"
    else
        # Fallback: copy-based backup
        if cp "$dest_file" "$backup_file" 2>/dev/null; then
            has_backup=true
        fi
    fi
fi

# Copy new binary — destination is free after rename
max_attempts=5
attempt=1
while [[ $attempt -le $max_attempts ]]; do
    if cp "$BINARY_PATH" "$dest_file" 2>/dev/null; then
        deploy_success=true
        break
    fi
    echo "  !! Target still locked; retrying ($attempt/$max_attempts)..."
    sleep 1
    attempt=$((attempt + 1))
done
```

---

## Rollback

If the copy fails after all retries, restore the `.old` backup:

### PowerShell

```powershell
if ($hasBackup -and (Test-Path $backupFile) -and (-not (Test-Path $destFile))) {
    try {
        Rename-Item $backupFile $destFile -Force -ErrorAction Stop
        Write-Success "Rollback complete - previous version restored"
    } catch {
        Write-Fail "Rollback also failed: $_"
    }
}
```

### Bash

```bash
if [[ "$has_backup" == "true" ]] && [[ -f "$backup_file" ]] && [[ ! -f "$dest_file" ]]; then
    mv "$backup_file" "$dest_file" 2>/dev/null && \
        echo "  OK Rollback complete" || \
        echo "  XX Rollback also failed"
fi
```

**Important**: Use `Rename-Item` / `mv` for rollback too — not
`Copy-Item` / `cp`. If the reason for failure was a file lock, a copy
would also fail for the same reason.

---

## PATH Sync

When the deploy target directory differs from the PATH binary location,
a separate sync step is needed after deploy:

```
Scenario:
  Deploy target:  C:\Users\user\AppData\Local\<binary>\<binary>.exe
  PATH points to: D:\tools\<binary>\<binary>.exe

Solution:
  After deploying to the target, also copy/sync to the PATH location.
```

PATH sync also uses rename-first:

```powershell
$activeBinary = (Get-Command <binary>).Source
$activeBackup = "$activeBinary.old"

# Rename the PATH binary
Rename-Item $activeBinary $activeBackup -Force
# Copy the newly deployed binary to the PATH location
Copy-Item $deployedBinary $activeBinary -Force
```

---

## Why Not Just Copy with Retries?

The old approach (retry `Copy-Item` up to 20 times) fails when:

- The binary is the **currently running process** — it will never
  unlock until the process exits.
- Another process (e.g., an IDE, antivirus) holds a handle on the file.

Rename-first succeeds immediately because Windows does not block
renames on locked files. The retry loop after rename is only a safety
net for edge cases (antivirus scanning the destination path, etc.)
and typically needs only 1 attempt.

---

## Retry Count

With rename-first, reduce the retry count from 20 to **5**. The
destination is free after rename, so retries are for rare edge cases:

| Strategy | Max Retries | Delay | Total Wait |
|----------|-------------|-------|------------|
| Copy-only (old) | 20 | 500ms | 10 seconds |
| Rename-first (new) | 5 | 500ms | 2.5 seconds |

---

## Constraints

- Always try rename before falling back to copy-based backup.
- The `.old` file is left in place — it is cleaned up by a separate
  `update-cleanup` command (see [06-cleanup.md](06-cleanup.md)).
- Rollback uses rename/move, not copy.
- If rename fails AND copy-backup fails, log the error but continue
  attempting the deploy (the destination might still be writable).
- On Linux/macOS, `mv` is the rename equivalent and works identically.

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Retry-on-lock deploy and PATH sync |
| [03-self-update-mechanism.md](../03-general/03-self-update-mechanism.md) | Rename-first during self-update, rollback on failure |
| [09-build-deploy.md](../01-app/09-build-deploy.md) | gitmap-v27 rename-first deploy with retry |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
