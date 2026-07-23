# 05 — Handoff Mechanism (Windows)

## Purpose

Define how a CLI tool replaces its own binary on Windows, where the
running process holds a file lock that prevents self-overwrite.

---

## The Core Problem

When a user runs `<binary> update`:

1. The process `<binary>.exe` is running.
2. The update needs to replace `<binary>.exe` with a newer version.
3. Windows locks `<binary>.exe` — it cannot be overwritten or deleted.

Even with rename-first, the build script (`run.ps1`) is launched by
the running binary, which may still be alive when deploy happens.

---

## Solution: Three-Phase Copy-and-Handoff

The running binary creates a **temporary copy of itself**, launches
the copy as a **worker process**, and exits. The worker is a different
file, so it does not conflict with the original binary's lock. After
build & deploy, the worker hands off ONE MORE TIME — to the freshly
deployed binary — to perform cleanup. This is required because the
worker still holds a lock on its own handoff copy and cannot remove it.

### Flow

![Handoff Mechanism Flow](images/handoff-mechanism-flow.mmd)

```
<binary>.exe update
    │
    ├── Phase 1: Copy self → <binary>-update-<pid>.exe (temp dir or same dir)
    ├── Phase 1: Launch temp copy: <binary>-update-<pid>.exe update-runner
    ├── Phase 1: Parent waits (synchronous — cmd.Run, not cmd.Start)
    │
    └── Worker (<binary>-update-<pid>.exe update-runner):
        ├── Phase 2: Generate temp PowerShell script
        ├── Phase 2: Script calls run.ps1 (pull, build, deploy)
        ├── Phase 2: Deploy uses rename-first on the ORIGINAL binary
        ├── Phase 2: Verify version
        │
        └── Phase 3: Spawn deployed <binary>.exe update-cleanup (detached, ~2s delayed)
            └── Worker exits — its handoff file becomes unlocked
                └── Deployed binary removes handoff copy + .old backup
```

### Why the Parent Waits

The parent process (`<binary>.exe`) calls `cmd.Run()` (blocking), NOT
`cmd.Start()` (non-blocking). This keeps the terminal attached so the
user sees all build output. The parent holds the lock on its own file,
but the worker deploys to the ORIGINAL path — the lock conflict is
between parent and deploy, which rename-first resolves.

---

## Step 1 — Create Handoff Copy

```go
func createHandoffCopy() (string, error) {
    selfPath, err := os.Executable()
    if err != nil {
        return "", err
    }

    // Try same directory first, then fall back to temp dir
    selfDir := filepath.Dir(selfPath)
    copyName := fmt.Sprintf("<binary>-update-%d.exe", os.Getpid())
    copyPath := filepath.Join(selfDir, copyName)

    data, err := os.ReadFile(selfPath)
    if err != nil {
        return "", err
    }

    if err := os.WriteFile(copyPath, data, 0o755); err != nil {
        // Fallback to temp directory
        copyPath = filepath.Join(os.TempDir(), copyName)
        if err := os.WriteFile(copyPath, data, 0o755); err != nil {
            return "", err
        }
    }

    return copyPath, nil
}
```

### Same-Dir vs Temp-Dir

| Location | Advantage | Risk |
|----------|-----------|------|
| Same directory as binary | Simpler path resolution | May fail if directory is read-only |
| System temp directory | Always writable | Path may be long, some policies block execution from temp |

Try same-dir first, fall back to temp.

---

## Step 2 — Launch Worker

```go
func launchWorker(copyPath, repoPath string) error {
    args := []string{"update-runner", "--repo-path", repoPath}

    cmd := exec.Command(copyPath, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    // Synchronous — parent blocks until worker completes
    return cmd.Run()
}
```

**Critical**: Use `cmd.Run()` (blocking), not `cmd.Start()`.

---

## Step 3 — Worker Generates Script

The worker (running from the temp copy) generates a temporary
PowerShell script that calls the build script:

```go
func buildUpdateScript(repoPath, runPS1 string) string {
    return fmt.Sprintf(`
Set-Location '%s'

# Detect current deploy path
$deployDir = Split-Path (Get-Command <binary>).Source -Parent

# Show version before update
$beforeVersion = & <binary> version 2>&1
Write-Host "  Before: $beforeVersion"

# Run the build script with update flag
& '%s' -Update

# Show version after update
$afterVersion = & <binary> version 2>&1
Write-Host "  After:  $afterVersion"

if ($beforeVersion -eq $afterVersion) {
    Write-Host "  !! Version unchanged — source may already be up to date" -ForegroundColor Yellow
}
`, repoPath, runPS1)
}
```

### UTF-8 BOM

Write the script with a UTF-8 BOM so PowerShell correctly handles
Unicode characters (checkmarks, arrows, etc.):

```go
bom := []byte{0xEF, 0xBB, 0xBF}
tmpFile.Write(bom)
tmpFile.WriteString(script)
```

---

## Step 4 — Execute Script

```go
func runUpdateScript(scriptPath string) {
    cmd := exec.Command("powershell.exe",
        "-ExecutionPolicy", "Bypass",
        "-NoProfile", "-NoLogo",
        "-File", scriptPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    // Handle error...
}
```

### PowerShell Flags

| Flag | Purpose |
|------|---------|
| `-ExecutionPolicy Bypass` | Allow running unsigned scripts |
| `-NoProfile` | Skip loading user profile (faster, predictable) |
| `-NoLogo` | Suppress PowerShell banner |
| `-File` | Execute a script file (not a command string) |

---

## Linux/macOS — No Handoff Needed

On Unix systems, the binary can be replaced in-place. The update
command simply runs the build script directly:

```go
func executeUpdateUnix(repoPath string) {
    runSH := filepath.Join(repoPath, "run.sh")
    cmd := exec.Command("bash", runSH, "--update")
    cmd.Dir = repoPath
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    // Handle error...
}
```

No temporary copy, no worker process, no generated scripts.

---

## Binary-Based Handoff (No Source Repo)

When the source repo is not available, use a standalone updater binary:

```
<binary> update
    │
    ├── Source repo found? YES → source-based handoff (above)
    └── Source repo found? NO  → delegate to <binary>-updater
        │
        ├── Copy <binary>-updater → <binary>-updater-tmp-<pid>.exe
        ├── Launch: <binary>-updater-tmp-<pid>.exe update-worker <version>
        └── Worker:
            ├── Download install.ps1 from release assets
            ├── Execute install.ps1 (handles download + verify + install)
            ├── Verify version
            └── Clean up self
```

---

## Constraints

- The parent MUST wait synchronously (`cmd.Run()`) — never fire and
  forget.
- The worker runs from a DIFFERENT file than the original binary.
- Temporary scripts and copies are cleaned up after completion.
- The handoff copy name includes the PID to avoid collisions when
  multiple updates run simultaneously.
- On Unix, skip the handoff entirely — run `bash run.sh --update`
  directly.
- All output (stdout/stderr) must be piped to the terminal so the
  user sees progress.

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Self-update orchestration (Windows-safe handoff) |
| [03-self-update-mechanism.md](../03-general/03-self-update-mechanism.md) | Two-phase handoff, worker launch, foreground blocking |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
