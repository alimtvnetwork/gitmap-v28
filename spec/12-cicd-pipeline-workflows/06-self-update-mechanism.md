# Self-Update Mechanism

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Defines the CI/CD-relevant aspects of CLI self-update: what the release pipeline must produce and support for self-updating tools. The full client-side update implementation is in [`spec/14-update/`](../14-update/00-overview.md) — this document covers the **pipeline's responsibilities** toward enabling that flow.

---

## Core Problem

A running binary **cannot overwrite itself** on Windows. The update architecture uses a rename-first strategy to work around file locks. The CI/CD pipeline must produce artifacts that support both **source-based** and **binary-based** update paths.

---

## Placeholders

| Placeholder | Meaning | Example |
|-------------|---------|---------|
| `<binary>` | CLI binary name | `gitmap` |
| `<binary>.exe` | Windows binary with extension | `gitmap.exe` |
| `<deploy-dir>` | Directory where the binary is installed | `$env:LOCALAPPDATA\gitmap` |
| `<repo-root>` | Root of the source repository | `<repo-root>` |
| `<repo>` | GitHub repository path | `github.com/org/repo` |
| `<module>` | Go module path | `github.com/org/repo` |

---

## Two Update Strategies

The pipeline must support both strategies:

### Strategy 1: Source-Based (Build from Repo)

The user has the source repository locally. The CLI pulls latest, builds, and deploys using rename-first.

**Pipeline responsibility**: Embed the repository path via `-ldflags` at build time so the CLI can locate the source repo.

```yaml
- name: Build with embedded repo path
  run: |
    LDFLAGS="-s -w -X '<module>/constants.RepoPath=${{ github.workspace }}'"
    go build -ldflags "$LDFLAGS" -o dist/<binary> .
```

### Strategy 2: Binary-Based (Download Pre-Built)

No source repo available. The CLI downloads install scripts from the latest GitHub Release.

**Pipeline responsibility**: Produce version-pinned install scripts and publish them as release assets (see [04-install-script-generation.md](./04-install-script-generation.md)).

```
<binary> update
 │
 ├── Source repo found? → Source-Based (pull + build + rename-first deploy)
 └── No source repo?   → Binary-Based (download install script + execute)
```

---

## Deploy Path Resolution

The CLI resolves the deploy target with a 3-tier priority:

| Priority | Source | How |
|----------|--------|-----|
| 1 | CLI flag | `--deploy-path <path>` argument |
| 2 | Global PATH lookup | `os.Executable()` or `Get-Command` → resolve symlinks → extract directory |
| 3 | Config file default | JSON config or `$HOME/.local/bin` |

### Auto-Detection (Go)

```go
execPath, err := os.Executable()
if err != nil {
    return "", fmt.Errorf("cannot detect executable path: %w", err)
}

resolvedPath, err := filepath.EvalSymlinks(execPath)
if err != nil {
    return "", fmt.Errorf("cannot resolve symlinks: %w", err)
}

return filepath.Dir(resolvedPath), nil
```

### Nested Directory Detection

Many CLI deploy structures use `<binary>/<binary>.exe`. When detecting the PATH location, check if the parent directory name matches the binary name — if so, the deploy target is the **grandparent** directory.

> Full specification: [14-update/02-deploy-path-resolution.md](../14-update/02-deploy-path-resolution.md)

---

## Rename-First Deploy Strategy

### Why Rename-First

On Windows, a running process holds a file lock on its executable. You **cannot** delete or overwrite it. However, you **can** rename it.

### Flow

```
1. Rename running binary:  <binary>.exe  →  <binary>.exe.old
2. Copy new binary:        new/<binary>.exe  →  <binary>.exe (destination is now free)
3. On failure: rollback    <binary>.exe.old  →  <binary>.exe
```

### Implementation

```go
targetPath := filepath.Join(deployDir, binaryName)
oldPath := targetPath + ".old"

// Step 1: Rename running binary
if fileExists(targetPath) {
    _ = os.Remove(oldPath) // Remove previous .old if exists
    if err := os.Rename(targetPath, oldPath); err != nil {
        return fmt.Errorf("cannot rename running binary: %w", err)
    }
}

// Step 2: Copy new binary (retry up to 5 times)
for attempt := 1; attempt <= 5; attempt++ {
    if err := copyFile(newBinaryPath, targetPath); err == nil {
        break
    }
    if attempt == 5 {
        // Rollback: restore old binary
        _ = os.Rename(oldPath, targetPath)
        return fmt.Errorf("cannot deploy new binary after %d attempts", attempt)
    }
    time.Sleep(500 * time.Millisecond)
}

// Step 3: Clean up .old (best-effort — may be locked)
_ = os.Remove(oldPath)
```

### Platform Differences

| Platform | Rename Running Binary | Delete Running Binary | Notes |
|----------|----------------------|----------------------|-------|
| Windows | ✅ Works | ❌ File lock | Must rename first |
| Linux | ✅ Works | ✅ Works (inode) | Simpler but rename-first still works |
| macOS | ✅ Works | ✅ Works (inode) | Same as Linux |

### Retry Count

With rename-first, the destination is free after rename. Max **5 retries** with 500ms delay (not 20):

| Strategy | Max Retries | Delay | Total Wait |
|----------|-------------|-------|------------|
| Copy-only (old) | 20 | 500ms | 10 seconds |
| Rename-first (current) | 5 | 500ms | 2.5 seconds |

> Full specification: [14-update/03-rename-first-deploy.md](../14-update/03-rename-first-deploy.md)

---

## Handoff Mechanism (Windows Self-Replacement)

When the CLI updates **itself**, a copy-and-handoff is required because the running process holds a lock on its own file.

### Flow

```
<binary>.exe update
 │
 ├── 1. Copy self → <binary>-update-<PID>.exe (same dir or temp dir)
 ├── 2. Launch copy: <binary>-update-<PID>.exe update-runner
 ├── 3. Parent WAITS SYNCHRONOUSLY (cmd.Run(), NOT cmd.Start())
 │
 └── Worker (<binary>-update-<PID>.exe update-runner):
     ├── 4. Generate temp PowerShell script (UTF-8 BOM)
     ├── 5. Script calls run.ps1 -Update (pull, build, deploy)
     ├── 6. Deploy uses rename-first on the ORIGINAL binary
     ├── 7. Verify version
     └── 8. Clean up temp files
```

### Critical: Synchronous Wait

The parent process calls `cmd.Run()` (blocking), **NOT** `cmd.Start()` (non-blocking). This keeps the terminal attached so the user sees all build output.

```go
func launchWorker(copyPath, repoPath string) error {
    cmd := exec.Command(copyPath, "update-runner", "--repo-path", repoPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    // Synchronous — parent blocks until worker completes
    return cmd.Run()
}
```

### Linux/macOS — No Handoff

On Unix, `bash run.sh --update` directly — no copy, no worker.

### Binary-Based Handoff (No Source Repo)

When no source repo is available, the CLI delegates to a standalone `<binary>-updater` binary that downloads and executes `install.ps1` from the latest release.

> Full specification: [14-update/05-handoff-mechanism.md](../14-update/05-handoff-mechanism.md)

---

## Build Scripts

The pipeline supports cross-platform build scripts implementing a 4-step pipeline:

```
[1/4] Pull latest changes (git pull, branch check)
[2/4] Resolve dependencies (go mod tidy)
[3/4] Build binary (go build with ldflags)
[4/4] Deploy (rename-first to resolved target)
```

### `run.ps1` (Windows)

```powershell
$ErrorActionPreference = "Stop"

$Version = & go run . version 2>&1
$LdFlags = "-s -w -X '<module>/constants.Version=$Version'"
$OutputDir = "dist"
$BinaryName = "<binary>.exe"

New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
go build -ldflags $LdFlags -o "$OutputDir/$BinaryName" .

if ($LASTEXITCODE -ne 0) { throw "Build failed" }
Write-Host "Built $BinaryName ($Version)"
```

### `run.sh` (Linux/macOS)

```bash
#!/usr/bin/env bash
set -euo pipefail

VERSION="$(go run . version 2>&1)"
LDFLAGS="-s -w -X '<module>/constants.Version=$VERSION'"
OUTPUT_DIR="dist"
BINARY_NAME="<binary>"

mkdir -p "$OUTPUT_DIR"
CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/$BINARY_NAME" .
echo "Built $BINARY_NAME ($VERSION)"
```

> Full specification: [14-update/04-build-scripts.md](../14-update/04-build-scripts.md)

---

## Cleanup

After a successful update, clean up residual artifacts:

| Artifact | Location | When to Clean |
|----------|----------|---------------|
| `<binary>.exe.old` | `<deploy-dir>/` | After new binary is confirmed working |
| `<binary>-update-<PID>.exe` | Same dir or temp dir | After worker completes |
| `<binary>-update-*.ps1` | Temp dir | After script executes |

### Cleanup Subcommand

The CLI should provide an explicit cleanup command:

```
<binary> update-cleanup
```

This scans the binary's directory and temp directory for leftover artifacts and removes them.

### Automatic Cleanup

- **After successful update**: best-effort auto-cleanup (never fails the update).
- **On startup**: check for stale `.old` files and remove them silently.
- **`.old` files are the rollback safety net** — never delete them during deploy, only after success.

```go
func cleanupStaleFiles(deployDir, binaryName string) {
    patterns := []string{
        filepath.Join(deployDir, binaryName+".old"),
        filepath.Join(deployDir, "<binary>-update-*"),
    }
    for _, pattern := range patterns {
        matches, _ := filepath.Glob(pattern)
        for _, m := range matches {
            _ = os.Remove(m) // Best-effort, ignore errors
        }
    }
}
```

> Full specification: [14-update/06-cleanup.md](../14-update/06-cleanup.md)

---

## Version String Normalization

All version comparisons must normalize the `v` prefix:

```go
func normalizeVersion(v string) string {
    v = strings.TrimSpace(v)
    v = strings.TrimPrefix(v, "v")
    return v
}

// "v1.2.0" and "1.2.0" are equal after normalization
```

After deploying, compare versions:

```bash
old_version="1.2.0"
new_version=$(<binary> version)

if [[ "$new_version" == "$old_version" ]]; then
    echo " !! Warning: version unchanged after update ($old_version)"
fi
```

---

## Constraints

- **Rename-first is mandatory on Windows** — never attempt to overwrite a running binary
- **Rollback on failure** — if copy fails after all retries, restore the `.old` binary
- **Best-effort cleanup** — `.old` deletion may fail if locked; handle gracefully
- **Cross-platform** — rename-first works on all platforms, use it universally
- **No elevated permissions** — install to user-local directories, not system directories
- **Version verification** — after update, run `<binary> version` to confirm
- **Synchronous handoff** — use `cmd.Run()`, never fire-and-forget
- **Updates are synchronous** — user sees all output in the same terminal session
- **Never leave the system without a working binary** — always rollback on failure

---

## Cross-References

- [Install Script Generation](./04-install-script-generation.md) — Install scripts used by binary-based update
- [Code Signing](./05-code-signing.md) — Signed binaries distributed via update
- [Go Binary Release Pipeline](./02-go-binary-deploy/02-release-pipeline.md) — Release pipeline that produces the binaries
- [Shared Conventions](./01-shared-conventions.md) — Version resolution patterns
- **Full Client-Side Specs:**
  - [Self-Update Overview](../14-update/01-self-update-overview.md) — Problem statement, strategies, command flow
  - [Deploy Path Resolution](../14-update/02-deploy-path-resolution.md) — 3-tier deploy target resolution
  - [Rename-First Deploy](../14-update/03-rename-first-deploy.md) — Rename-first with retry and rollback
  - [Build Scripts](../14-update/04-build-scripts.md) — Cross-platform build scripts
  - [Handoff Mechanism](../14-update/05-handoff-mechanism.md) — Windows self-replacement flow
  - [Cleanup](../14-update/06-cleanup.md) — Post-update artifact removal

---

*Self-update mechanism — v3.2.0 — 2026-04-10*
