# Update Command — Step-by-Step Workflow

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Source:** sibling reference implementation `cmd/update.go` and `cmd/updatecleanup.go`

---

## Purpose

Provide a complete, step-by-step walkthrough of the `<binary> update` and `<binary> update-cleanup` commands — derived directly from the sibling reference implementation. Any AI or engineer reading this document should be able to implement both commands without ambiguity.

---

## Command: `<binary> update`

### Overview

The `update` command replaces the currently running CLI binary with a newer version. It supports two strategies:

1. **Source-Based Update** — Pull latest code, build from source, deploy.
2. **Binary-Based Update** — Delegate to a standalone `<binary>-updater` binary that downloads a pre-built release.

The command automatically selects the appropriate strategy based on whether a source repository path can be resolved.

---

### Step-by-Step Flow

#### Step 1 — Network Check

```
requireOnline()
```

Verify internet connectivity before proceeding. If offline, print an error and exit immediately. This prevents wasted time on git pull or download failures.

#### Step 2 — Resolve Source Repository Path

```
repoPath := resolveRepoPath()
```

The repo path is resolved using a **priority cascade**:

| Priority | Method | Description |
|----------|--------|-------------|
| 1 | `--repo-path` flag | CLI flag passed by the user |
| 2 | Embedded constant | Compiled into the binary via `-ldflags -X` |
| 3 | Database lookup | Previously saved path from a prior update |
| 4 | Interactive prompt | Ask the user to provide the path |
| 5 | Updater fallback | Delegate to `<binary>-updater` (binary-based update) |

Each resolved path is **saved to the database** for future use (`saveRepoPathToDB()`).

**If all 4 source methods fail**, the command attempts updater fallback (Step 2b).

#### Step 2b — Updater Fallback (Binary-Based)

```go
func tryUpdaterFallback() bool {
    updaterPath, err := exec.LookPath("<binary>-updater")
    if err != nil {
        return false  // updater not found on PATH
    }

    cmd := exec.Command(updaterPath, "run")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        // Propagate exit code if available
        var exitErr *exec.ExitError
        if errors.As(err, &exitErr) {
            os.Exit(exitErr.ExitCode())
        }
        return false
    }

    return true
}
```

If `<binary>-updater` is found on PATH:
- Execute `<binary>-updater run` synchronously.
- Pipe all stdout/stderr to the terminal.
- On success, exit with code 0.
- On failure, propagate the exit code.

If `<binary>-updater` is not found, print error and exit.

#### Step 3 — Locate Current Binary

```go
selfPath, err := os.Executable()
```

Resolve the absolute path of the currently running binary. This is needed for the handoff copy.

#### Step 4 — Create Handoff Copy

```go
copyPath := createHandoffCopy(selfPath)
```

Create a temporary copy of the running binary to work around Windows file locks:

1. Generate a name: `<binary>-update-<PID>.exe` (Windows) or `<binary>-update-<PID>` (Unix).
2. Try copying to the **same directory** as the binary.
3. If that fails (e.g., read-only directory), copy to the **system temp directory**.
4. Set executable permissions on Unix (`chmod 0755`).

```
Naming convention:
  Windows: <binary>-update-12345.exe
  Unix:    <binary>-update-12345
```

#### Step 5 — Print Status

```
 [+] Self-update active
     Running:  C:\tools\<binary>\<binary>.exe
     Worker:   C:\tools\<binary>\<binary>-update-12345.exe
```

#### Step 6 — Launch Handoff Worker

```go
func launchHandoff(copyPath, repoPath string) {
    args := []string{"update-runner"}
    if hasFlag("--verbose") {
        args = append(args, "--verbose")
    }
    args = append(args, "--repo-path", repoPath)

    cmd := exec.Command(copyPath, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    err := cmd.Run()  // BLOCKING — parent waits
}
```

**Critical behaviors:**
- Uses `cmd.Run()` (blocking), NOT `cmd.Start()` (non-blocking).
- Terminal stays attached — user sees all build/deploy output.
- The `--verbose` flag is forwarded if present.
- On error, propagate the worker's exit code.

#### Step 7 — Worker Executes Update (`update-runner`)

The `update-runner` is a **hidden subcommand** — not shown in `--help`. It runs from the handoff copy and performs the actual update:

```go
func runUpdateRunner() {
    repoPath := resolveRepoPath()
    initRunnerVerbose()
    fmt.Printf(" Starting update...")
    fmt.Printf(" Repo: %s", repoPath)
    executeUpdate(repoPath)
}
```

The `executeUpdate()` function:

**On Windows:**
1. Generate a temporary PowerShell script with UTF-8 BOM.
2. Script calls `run.ps1 -Update` in the repo directory.
3. Execute via `powershell.exe -ExecutionPolicy Bypass -NoProfile -NoLogo -File <script>`.

**On Linux/macOS:**
1. Execute `bash run.sh --update` directly.
2. No temporary script needed.

Both paths invoke the build script which:
1. Pulls latest code (`git pull`).
2. Resolves dependencies (`go mod download`).
3. Builds the binary (`go build`).
4. Deploys using **rename-first** strategy (see `03-rename-first-deploy.md`).
5. Verifies the new version.

#### Step 8 — Version Verification

After deploy, compare versions:

```
 Before: v2.48.0
 After:  v2.49.0
```

If versions are identical:
```
 !! Version unchanged — source may already be up to date
```

#### Step 9 — Exit

The worker exits, the parent process receives the exit code and terminates.

### Visual Diagram

See [`diagrams/01-self-update-workflow.mmd`](./diagrams/01-self-update-workflow.mmd) for a complete Mermaid flowchart of the update command lifecycle.

---

## Command: `<binary> update-cleanup`

### Overview

Removes leftover temporary files from previous update operations. Should be run periodically or after a failed update.

### Step-by-Step Flow

#### Step 1 — Print Start Message

```
 Cleaning up update artifacts...
```

#### Step 2 — Clean Temp Copies

Remove handoff binaries from previous updates:

1. Build glob patterns:
   - `<TEMP_DIR>/<binary>-update-*` (system temp directory)
   - `<BINARY_DIR>/<binary>-update-*` (same directory as binary)
2. For each match:
   - Normalize path (`filepath.Clean`).
   - Skip if already processed (deduplication via `seen` map).
   - Skip if it's the currently running binary.
   - Delete the file.
   - Print: `  [OK] Removed <filename>`.

#### Step 3 — Clean Old Backups

Remove `.old` backup files created by rename-first deploy:

1. Read the deploy path from `powershell.json` config:
   ```
   <repo-root>/.<binary>/powershell.json → "DeployPath" field
   ```
2. Build glob pattern: `<deploy-path>/<binary>/*.old`
3. Delete all matches.
4. Print: `  [OK] Removed <filename>`.

#### Step 4 — Print Summary

```
 ✓ Cleaned 3 artifact(s)
```

Or if nothing to clean:

```
 Nothing to clean up.
```

### Visual Diagram

See [`diagrams/02-update-cleanup-workflow.mmd`](./diagrams/02-update-cleanup-workflow.mmd) for a complete Mermaid flowchart of the cleanup command lifecycle.

---

## Complete Decision Tree

```
<binary> update
 │
 ├─ requireOnline() ── FAIL → "No internet" → EXIT 1
 │
 ├─ resolveRepoPath()
 │   ├─ --repo-path flag?     → YES → use it, save to DB
 │   ├─ Embedded constant?    → YES → use it, save to DB
 │   ├─ Database lookup?      → YES → use it
 │   ├─ Interactive prompt?   → YES → use it, save to DB
 │   └─ All failed?
 │       ├─ <binary>-updater on PATH? → YES → delegate, EXIT
 │       └─ NO → "No repo path" → EXIT 1
 │
 ├─ os.Executable() → selfPath
 │
 ├─ createHandoffCopy(selfPath)
 │   ├─ Try same directory → SUCCESS → copyPath
 │   └─ Try temp directory → SUCCESS → copyPath
 │                         → FAIL    → EXIT 1
 │
 ├─ Print status (running binary, worker path)
 │
 └─ launchHandoff(copyPath, repoPath)
     │
     └─ Worker (update-runner):
         ├─ Platform?
         │   ├─ Windows → Generate PS1 script → run.ps1 -Update
         │   └─ Unix    → bash run.sh --update
         │
         ├─ Build script:
         │   ├─ git pull
         │   ├─ go mod download
         │   ├─ go build
         │   ├─ rename-first deploy
         │   └─ verify version
         │
         └─ Print before/after version → EXIT
```

---

## Error Handling Matrix

| Error | Source | Response |
|-------|--------|----------|
| No internet | `requireOnline()` | Print error, exit 1 |
| Cannot find binary path | `os.Executable()` | Print error, exit 1 |
| Handoff copy failed (both locations) | `createHandoffCopy()` | Print error, exit 1 |
| No repo path (all methods exhausted) | `resolveRepoPath()` | Print error, exit 1 |
| Updater not found (fallback) | `tryUpdaterFallback()` | Print error, exit 1 |
| Worker process failed | `launchHandoff()` | Propagate exit code |
| Build script failed | `executeUpdate()` | Print error, exit 1 — original binary untouched |
| Deploy failed (rename-first) | `run.ps1` / `run.sh` | Rollback from `.old` backup |
| Version unchanged after update | Version check | Warn (yellow), exit 0 |

---

## Key Implementation Notes

### 1. Synchronous Execution is Mandatory

The parent process MUST use `cmd.Run()` (blocking). Using `cmd.Start()` would detach the terminal and the user would lose all output. This is the most common implementation mistake.

### 2. PID-Based Naming Prevents Collisions

The handoff copy includes `os.Getpid()` in the filename, preventing conflicts when multiple terminal sessions run `update` simultaneously.

### 3. UTF-8 BOM for PowerShell Scripts

Windows-generated PowerShell scripts must include the UTF-8 BOM (`0xEF, 0xBB, 0xBF`) prefix so that status icons (✓, ⚠, etc.) render correctly in all terminal emulators.

### 4. Exit Code Propagation

When the worker process fails, the parent extracts and propagates the exact exit code:

```go
var exitErr *exec.ExitError
if errors.As(err, &exitErr) {
    os.Exit(exitErr.ExitCode())
}
```

This ensures CI/CD pipelines and scripts can detect update failures.

### 5. Verbose Mode Forwarding

The `--verbose` flag is forwarded from the parent to the worker so that detailed logging persists across the handoff boundary.

---

## Cross-References

- [Self-Update Overview](./01-self-update-overview.md) — High-level architecture and strategy selection
- [Deploy Path Resolution](./02-deploy-path-resolution.md) — How the deploy target directory is determined
- [Rename-First Deploy](./03-rename-first-deploy.md) — File replacement strategy used during deploy
- [Build Scripts](./04-build-scripts.md) — `run.ps1` / `run.sh` invoked by the worker
- [Handoff Mechanism](./05-handoff-mechanism.md) — Deep-dive on the copy-and-handoff pattern
- [Cleanup](./06-cleanup.md) — Post-update artifact removal lifecycle
- [Updater Binary](./19-updater-binary.md) — Standalone updater for binary-based updates
- [Network Requirements](./20-network-requirements.md) — HTTP client and connectivity checks
- [Configuration File](./21-config-file.md) — Where repo path and deploy defaults are stored

---

*Update command workflow — v3.2.0 — 2026-04-13*
