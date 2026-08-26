# Consolidated: Self-Update & App Update — Complete Reference

**Version:** 3.3.0  
**Updated:** 2026-04-22  
**Source Module:** [`spec/14-update/`](../14-update/00-overview.md)

---

## Purpose

This is the **standalone consolidated reference** for CLI self-update functionality and release distribution. An AI reading only this file must be able to implement the complete update system without consulting source specs.

### Source-Folder Coverage Map

The source module `spec/14-update/` contains 25 specs + diagrams. The table below confirms every one is covered.

| Source File | Section | Status |
|-------------|---------|--------|
| `01-self-update-overview.md` | §Core Problem, §Two Strategies | ✅ Full |
| `02-deploy-path-resolution.md` | §Deploy Path Resolution | ✅ Full |
| `03-rename-first-deploy.md` | §Rename-First Deploy | ✅ Full |
| `04-build-scripts.md` | §Build Scripts | ✅ Full |
| `05-handoff-mechanism.md` | §Handoff Mechanism | ✅ Full |
| `06-cleanup.md` | §Cleanup | ✅ Full |
| `07-console-safe-handoff.md` | §Console-Safe Handoff (added) | ✅ Full |
| `08-repo-path-sync.md` | §Repo Path Sync (added) | ✅ Full |
| `09-version-verification.md` | §Version Verification (added) | ✅ Full |
| `10-last-release-detection.md` | §Last-Release Detection (added) | ✅ Full |
| `11-windows-icon-embedding.md` | §Cross-ref to CI/CD spec | ✅ Reference |
| `12-code-signing.md` | §Cross-ref to CI/CD spec | ✅ Reference |
| `13-release-assets.md` | §Release Assets | ✅ Full |
| `14-checksums-verification.md` | §Checksums | ✅ Full |
| `15-release-versioning.md` | §Release Versioning | ✅ Full |
| `16-cross-compilation.md` | §Cross-Compilation | ✅ Full |
| `17-release-pipeline.md` | §Release Pipeline | ✅ Full |
| `18-install-scripts.md` | §Install Scripts | ✅ Full |
| `19-updater-binary.md` | §Updater Binary | ✅ Full |
| `20-network-requirements.md` | §Network Requirements | ✅ Full |
| `21-config-file.md` | §Config File | ✅ Full |
| `22-update-command-workflow.md` | §Update Command Workflow | ✅ Full |
| `23-install-script-version-probe.md` | §Install Script Version Probe (added) | ✅ Full |
| `24-update-check-mechanism/` | §Update Check Mechanism (added) | ✅ Full |
| `25-release-pinned-installer.md` | §Release-Pinned Installer (added) | ✅ Full |
| `diagrams/` | Diagram links | ✅ |

---

## Core Problem

A running binary **cannot overwrite itself** on Windows. The entire update architecture exists to work around this constraint using a rename-first strategy.

---

## Two Update Strategies

### Strategy 1: Source-Based (Build from Repo)

When the user has the source repository locally:

```
1. Resolve source repo path
2. Pull latest changes (git pull)
3. Resolve dependencies (go mod tidy)
4. Build binary (go build with ldflags)
5. Rename-first deploy
6. Verify version
7. Cleanup artifacts
```

### Strategy 2: Binary-Based (Download Pre-Built)

When no source repo is available:

```
1. Locate <binary>-updater on PATH
2. Execute: <binary>-updater --install-dir <path>
3. Updater queries GitHub API for latest version
4. Downloads install script
5. Executes install (handles download + verify + install)
```

The `<binary> update` command auto-detects which strategy to use based on whether a source repo path can be resolved.

---

## Deploy Path Resolution (3-Tier Priority)

| Priority | Source | When Used |
|----------|--------|-----------|
| 1 | CLI flag (`--deploy-path`) | Explicit override |
| 2 | Global PATH lookup | Binary already installed |
| 3 | Config file default | First-time install |

### Repo Path Resolution (5-Tier Priority)

| Priority | Method | Description |
|----------|--------|-------------|
| 1 | `--repo-path` flag | CLI flag passed by user |
| 2 | Embedded constant | Compiled into binary via `-ldflags -X` |
| 3 | Database lookup | Previously saved path from prior update |
| 4 | Interactive prompt | Ask user to provide the path |
| 5 | Updater fallback | Delegate to `<binary>-updater` (binary-based update) |

Each resolved repo path is **saved to the database** for future use.

### Special Cases

- Symlink resolution: Follow symlinks to find the real binary path
- Nested directory detection: `<binary>/<binary>.exe` pattern → resolve to parent
- `installed-dir` utility command: Print resolved install directory

---

## Rename-First Deploy — Complete Specification

Windows file lock workaround with rollback safety:

```
Step 1: Rename running <binary>.exe → <binary>.exe.old
Step 2: Copy new binary → <binary>.exe (destination now free)
Step 3: Verify new binary works (execute --version)
Step 4: If verify fails → rollback (rename .old back)
Step 5: If verify succeeds → .old is cleanup candidate
```

### Retry Logic

| Parameter | Value |
|-----------|-------|
| Max retries | 5 |
| Retry delay | 500ms |
| Retry applies to | Rename and copy operations |
| Previous approach | 20 retries with copy-only (replaced) |

### PATH Sync

When the deploy target differs from the PATH location, the script must sync:
1. Copy binary to deploy target
2. Update PATH to point to deploy target (if different from current PATH entry)

### Implementation Available In

| Platform | Script |
|----------|--------|
| PowerShell | `run.ps1 -Update` |
| Bash | `run.sh --update` |

---

## Build Scripts (Cross-Platform)

Both `run.ps1` and `run.sh` implement a 4-step pipeline:

```
[1/4] Pull latest changes
[2/4] Resolve dependencies (go mod tidy)
[3/4] Build binary (go build with ldflags)
[4/4] Deploy (rename-first)
```

### Key Features

| Feature | Detail |
|---------|--------|
| JSON config loading | Reads build config from `config.json` |
| Git pull | With conflict resolution (abort merge, force pull) |
| Source file validation | Checks `main.go` exists before building |
| Data folder copy | Copies `data/` to deploy directory |
| Colored logging | Consistent prefix scheme: `[1/4]`, `[OK]`, `[FAIL]` |
| -ldflags embedding | `Version`, `RepoPath`, `CommitSHA`, `BuildDate` |

---

## Handoff Mechanism (Windows Self-Replacement)

### Flow

```
1. Running binary detects Windows OS
2. Copies itself to <binary>-update-<PID>.exe (temp worker)
3. Launches worker with: <worker> update-runner --repo-path <path> --deploy-path <path>
4. Parent waits SYNCHRONOUSLY (cmd.Run(), NOT cmd.Start())
5. Worker generates temp PowerShell script
6. PowerShell script calls: run.ps1 -Update
7. run.ps1 performs rename-first deploy on the ORIGINAL binary path
8. Worker exits → parent resumes → prints result
```

### Critical Rules

- **Synchronous wait**: Parent uses `cmd.Run()` to block until worker completes
- **Output piped**: `cmd.Stdout = os.Stdout`, `cmd.Stderr = os.Stderr` — user sees all output
- **PID-based naming**: Worker binary includes parent PID to avoid collisions

### Unix

No handoff needed — `bash run.sh --update` directly replaces the binary.

### Binary-Based Fallback

When no source repo is found:

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
    return cmd.Run() == nil
}
```

---

## Updater Binary (`<binary>-updater`)

### Purpose

Separate Go module that downloads and installs updates when the main binary cannot replace itself.

### Module Structure

```
<binary>-updater/
├── go.mod                  # Independent module
├── main.go
├── cmd/root.go             # CLI command definitions
├── updater/
│   ├── download.go         # GitHub API + asset download
│   ├── install.go          # Rename-first deploy logic
│   └── verify.go           # Post-install version verification
├── version/version.go      # Embedded version constants
├── winres.json             # Windows icon/version resources
└── assets/icon-updater.png # Distinct icon from main binary
```

### CLI Interface

```
<binary>-updater [flags]
  --install-dir <path>    Target directory (required)
  --version <ver>         Specific version (default: latest)
  --repo <owner/repo>     GitHub repository (embedded at build time)
  --binary-name <name>    Name of binary to update (embedded)
  --skip-checksum         Skip SHA-256 verification (not recommended)
  --verbose               Enable detailed output
```

### Update Flow

```
1. Resolve target version (flag → GitHub API latest)
2. Construct download URL
3. Download binary archive + checksums.txt
4. Verify SHA-256 (unless --skip-checksum)
5. Rename-first deploy to install-dir
6. Verify: execute <binary> version
7. Print result
```

---

## Cleanup

Post-update artifact removal:

| Artifact | Location | Created By |
|----------|----------|-----------|
| `<binary>-update-<PID>.exe` | Same dir or temp dir | Handoff mechanism |
| `<binary>-update-*.ps1` | Temp dir | PowerShell script generation |
| `<binary>.exe.old` | Deploy directory | Rename-first strategy |

### Rules

- **Explicit `update-cleanup` subcommand** for manual cleanup
- **Auto-cleanup** after successful update (best-effort, never fails the update)
- **`.old` files are the rollback safety net** — never deleted during deploy, only after verification
- Cleanup scans both the binary directory and temp directory

---

## Release Assets

### Naming Convention

```
<binary>-<os>-<arch>.<ext>
```

| OS | Extension |
|----|-----------|
| Windows | `.zip` |
| Linux/macOS | `.tar.gz` |

### Standard Release (9 files)

| # | File | Description |
|---|------|-------------|
| 1-6 | `<binary>-{os}-{arch}.{ext}` | 6 compressed archives |
| 7 | `checksums.txt` | SHA-256 checksums of all archives |
| 8 | `install.ps1` | PowerShell installer |
| 9 | `install.sh` | Bash installer |

**Archive contents:** Single flat binary — no nested directories.

---

## Checksums & Verification

| Rule | Detail |
|------|--------|
| Algorithm | SHA-256 exclusively — never MD5/SHA-1 |
| Source | Generated from **compressed archives**, not raw binaries |
| Requirement | **Mandatory** — never optional or skippable |
| macOS fallback | `shasum -a 256` when `sha256sum` unavailable |
| Format | `sha256hash  filename` (two-space separator) |
| TOCTOU | Checksumming archives (not raw files) prevents time-of-check issues |

---

## Release Versioning

### Version Resolution (3-Tier)

| Priority | Source | Example |
|----------|--------|---------|
| 1 | Explicit argument | `--version 1.3.0` |
| 2 | Bump flag | `--bump minor` |
| 3 | Current source constant | Read from `version.go` |

### Rules

- SemVer 2.0.0 with `v` prefix normalization (`v1.2.0` → `1.2.0`)
- Atomic version sync: source constant + CHANGELOG + metadata in same commit
- Release branch strategy: `main → release/x.y.z → tag → merge back`
- Optional `latest.json` for programmatic version queries

```json
{
  "Version": "1.3.0",
  "ReleasedAt": "2026-04-16T12:00:00Z",
  "DownloadUrl": "https://github.com/.../releases/download/v1.3.0/"
}
```

---

## Cross-Compilation

### 6-Target Build Matrix

| OS | Arch | Binary Name |
|----|------|-------------|
| windows | amd64 | `<binary>-windows-amd64.exe` |
| windows | arm64 | `<binary>-windows-arm64.exe` |
| linux | amd64 | `<binary>-linux-amd64` |
| linux | arm64 | `<binary>-linux-arm64` |
| darwin | amd64 | `<binary>-darwin-amd64` |
| darwin | arm64 | `<binary>-darwin-arm64` |

### Build Rules

| Rule | Detail |
|------|--------|
| `CGO_ENABLED=0` | Fully static binaries, no C dependencies |
| `-ldflags -X` | Embeds Version, RepoPath, CommitSHA, BuildDate |
| `.exe` extension | Mandatory for Windows only |
| Output directory | All binaries to `dist/` — single source of truth |
| Never rebuild | All packaging/publishing operates on existing artifacts |

---

## Network Requirements

| Requirement | Detail |
|-------------|--------|
| Connectivity check | `requireOnline()` before any network operation |
| GitHub API | Used for version queries and release downloads |
| Timeout | 30s for API calls, 5m for binary downloads |
| Retry | 3 attempts with exponential backoff for downloads |

---

## Config File

Update configuration stored in a config file:

```json
{
  "RepoPath": "/home/user/projects/my-cli",
  "DeployPath": "/usr/local/bin",
  "LastUpdateCheck": "2026-04-16T12:00:00Z",
  "AutoUpdate": false
}
```

---

## Update Command Workflow Summary

### `<binary> update`

```
1. requireOnline() — verify internet
2. resolveRepoPath() — 5-tier priority cascade
3. If source found → source-based update
4. If no source → try updater fallback
5. If no updater → error and exit
```

### `<binary> update-cleanup`

```
1. Scan binary dir for *.old files
2. Scan temp dir for <binary>-update-* files
3. Delete all found artifacts
4. Report results
```

---

## Key Rules

1. ALL self-update content belongs in `spec/14-update/`
2. CI/CD pipeline specs (GitHub Actions) belong in `spec/12-cicd-pipeline-workflows/`
3. Rename-first is mandatory on Windows — never attempt to overwrite a running binary
4. Updates are always synchronous — user sees all output in the terminal
5. Never leave the system without a working binary — always rollback on failure
6. Binaries are built exactly once — no stage may trigger a rebuild
7. Checksum verification is mandatory — never optional or skippable
8. `.old` files are rollback safety — never delete during deploy

---

## File Inventory

| # | File | Description |
|---|------|-------------|
| 01 | `01-self-update-overview.md` | Problem statement, strategies, command flow |
| 02 | `02-deploy-path-resolution.md` | 3-tier deploy target resolution |
| 03 | `03-rename-first-deploy.md` | Rename-first with retry and rollback |
| 04 | `04-build-scripts.md` | Cross-platform build scripts |
| 05 | `05-handoff-mechanism.md` | Windows self-replacement flow |
| 06 | `06-cleanup.md` | Post-update artifact removal |
| 07 | `07-release-assets.md` | Asset naming and packaging |
| 08 | `08-checksums-verification.md` | SHA-256 generation and verification |
| 09 | `09-release-versioning.md` | Version resolution, tagging, changelog |
| 10 | `10-cross-compilation.md` | 6-target static build matrix |
| 11 | `11-release-pipeline.md` | End-to-end CI/CD release workflow |
| 12 | `12-install-scripts.md` | One-liner installer generation |
| 13 | `13-updater-binary.md` | Standalone updater binary architecture |
| 14 | `14-network-requirements.md` | Connectivity checks, timeouts, retry |
| 15 | `15-config-file.md` | Update configuration storage |
| 16 | `16-update-command-workflow.md` | Step-by-step update/cleanup commands |

---

## Console-Safe Handoff

**Source:** `14-update/07-console-safe-handoff.md`

On Windows, when the running CLI hands off to a child updater process, the parent's stdout/stderr handles must be **detached** before the child inherits them — otherwise console output races or hangs. Required steps:

1. Flush all parent buffers (`os.Stdout.Sync()`).
2. Spawn child with `syscall.SysProcAttr{HideWindow: false, CreationFlags: CREATE_NEW_CONSOLE}`.
3. Parent exits **immediately** with code 0; do not wait on child.
4. Child writes status to a pre-agreed temp log file the parent can poll on next launch.

On macOS / Linux this is a no-op — `exec` semantics handle handoff cleanly.

---

## Repo Path Sync

**Source:** `14-update/08-repo-path-sync.md`

When using Strategy 1 (source-based update), the resolved source-repo path is **persisted to the config file** so subsequent updates skip path discovery. Sync rules:

- Path is stored as an **absolute path** (resolve symlinks at write time).
- If the path becomes invalid (deleted, moved), the next update falls back to Strategy 2 (binary-based) and clears the stale entry.
- Multiple repos are not supported — the most recently used wins.

---

## Version Verification

**Source:** `14-update/09-version-verification.md`

After deploy completes, the new binary is invoked with `--version` and the output is matched against the expected version string. Mismatch triggers automatic rollback:

1. Restore the renamed `<binary>.old` to `<binary>`.
2. Delete the failed deploy artifact.
3. Exit with code 2 + log the version mismatch (expected vs got).

Verification timeout is **5 seconds**; longer indicates the new binary is hung and triggers rollback.

---

## Last-Release Detection

**Source:** `14-update/10-last-release-detection.md`

The CLI caches the last-checked release version in `~/.config/<app>/last-release.json` with a 24-hour TTL:

```json
{
  "Version": "v1.5.2",
  "PublishedAt": "2026-04-20T08:00:00Z",
  "CheckedAt": "2026-04-22T10:15:00Z",
  "AssetsByTarget": { "linux/amd64": "...", "windows/amd64": "..." }
}
```

Within the TTL, `<binary> update --check` returns instantly without hitting GitHub. After TTL expiry, the cache is refreshed transparently.

---

## Install Script Version Probe

**Source:** `14-update/23-install-script-version-probe.md`

The generated install script has a `--probe` flag that prints the version it would install **without actually installing**. Used by the updater to confirm the install script matches the intended target version before execution. Probe must:

- Exit 0 with the version string on stdout.
- Never write to the filesystem.
- Complete in under 2 seconds (no network calls beyond the version manifest fetch).

---

## Update Check Mechanism

**Source:** `14-update/24-update-check-mechanism/`

The proactive update-check runs **once per launch**, in a background goroutine, with the following rules:

| Rule | Value |
|------|-------|
| Frequency cap | Max 1 check per 24 hours per user |
| Network timeout | 3 seconds — fail silently on timeout |
| Notification | Subtle banner on next interactive command, never blocking |
| Disable flag | `<binary> config set update.checkOnLaunch false` |
| CI detection | Skip check entirely when `CI=true` |

The check never auto-installs — user must explicitly run `<binary> update`.

---

## Release-Pinned Installer

**Source:** `14-update/25-release-pinned-installer.md`

For supply-chain integrity, the README install one-liner pins:

1. The release tag (`v1.5.2`, never `latest`).
2. The install script SHA-256.
3. The binary checksum file SHA-256.

Example:

```bash
curl -fsSL https://github.com/org/app/releases/download/v1.5.2/install.sh \
  | sha256sum -c <(echo "abc123...  -") \
  | bash
```

Pinned installers are regenerated per release. Rolling-tag installers (`latest`) are explicitly forbidden — they break reproducibility and enable supply-chain attacks.

---

## Source-File Coverage Update

The original inventory below mapped source files 01–16 only. The added sections above cover source files **07-console-safe-handoff**, **08-repo-path-sync**, **09-version-verification**, **10-last-release-detection**, **23-install-script-version-probe**, **24-update-check-mechanism/**, and **25-release-pinned-installer**.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Full Source | `../14-update/00-overview.md` |
| CI/CD Pipeline Workflows | `../12-cicd-pipeline-workflows/00-overview.md` |
| Install Script Generation | `../12-cicd-pipeline-workflows/04-install-script-generation.md` |
| Code Signing | `../12-cicd-pipeline-workflows/05-code-signing.md` |
| Installation Flow | `../12-cicd-pipeline-workflows/08-installation-flow.md` |

---

*Consolidated self-update & app update — v3.2.0 — 2026-04-16*

---

## §15 Install-Script Version Probe — Runtime Contract (Verbatim)

This section is the **exact runtime contract** for the 20-repo version probe. A blind AI re-implementing or modifying the probe must use these exact identifiers — host integrations and CI tooling depend on them.

### 15.1 Identity-Resolution Regex (Verbatim)

```
^https?://[^/]+/(?<owner>[^/]+)/(?<base>[A-Za-z0-9._-]+?)-v(?<ver>\d+)/[^/]+/install\.(ps1|sh)(\?.*)?$
```

**Capture groups:**
- `owner` — GitHub org/user
- `base` — repo prefix (e.g., `movie-cli`)
- `ver` — current numbered version (digits only)

If the URL does not match, the probe logs `version probe disabled (no self-identity)` and falls back to the local install.

### 15.2 Environment Variable Contract

| Variable | Type | Default | Purpose |
|----------|------|---------|---------|
| `INSTALL_PROBE_OWNER` | string | (none) | Override owner from URL regex (Step 2 of identity ladder) |
| `INSTALL_PROBE_BASE` | string | (none) | Override base name |
| `INSTALL_PROBE_VERSION` | int | (none) | Override current version |
| `INSTALL_PROBE_HANDOFF_DEPTH` | int | `0` | Recursion guard (max `2` — depth `3` aborts) |
| `PROBE_OWNER_DEFAULT` | string | (compiled-in) | Built-in fallback for Step 3 |
| `PROBE_BASE_DEFAULT` | string | (compiled-in) | Built-in fallback for Step 3 |

**Names are case-sensitive and immutable.** Do not abbreviate or rename — host CI systems and support scripts grep for these literal strings.

### 15.3 Compiled-In Constants (Verbatim)

| Constant | Value | Source File | Purpose |
|----------|-------|-------------|---------|
| `PROBE_VERSION_FALLBACK` | `14` | `install.sh` line 65, `install.ps1` line ~83 | Fail-open `currentVersion` when identity ladder fully fails |
| Probe range | `currentVersion+1 .. currentVersion+20` | both scripts | 20 candidates |
| HEAD timeout | `2` seconds per request | both scripts | Per-request limit |
| Total settle window | `4` seconds (timeout × 2) | both scripts | Maximum wait before deciding |
| Success HTTP codes | `200`, `301`, `302` | both scripts | Counts as "responder" |
| Probe URL template | `https://raw.githubusercontent.com/<owner>/<base>-v<N>/main/install.<ext>` | both scripts | Per-candidate URL |
| Max handoff depth | `3` | both scripts | Aborts with `Probe loop guard triggered — aborting.` |

### 15.4 Skip Flags (CLI Arguments)

| Flag (Bash) | Flag (PowerShell) | Effect |
|-------------|-------------------|--------|
| `-n` | `-n` | Skip probe entirely (alias) |
| `--no-latest` | `-NoLatest` | Skip probe (long form) |
| `--no-probe` | `-NoProbe` | Skip probe (alternate long form) |

All three resolve to the same internal `should_skip_probe()` / `Test-ShouldSkipProbe` early-exit.

### 15.5 Identity Ladder (4 Steps, Fail-Open)

```
1. Parse $0 / $MyInvocation URL via §15.1 regex   → if match, set owner/base/ver
2. Read INSTALL_PROBE_OWNER/_BASE/_VERSION env    → override any/all from Step 1
3. Read PROBE_OWNER_DEFAULT/_BASE_DEFAULT consts  → fill any still-missing field
4. If currentVersion still unset → currentVersion = PROBE_VERSION_FALLBACK (14)
   If owner OR base still unset  → log "version probe disabled" and skip
```

### 15.6 Middle-Out Dispatch Order

Although all 20 HEAD requests fire in parallel (ordering is functionally moot), the dispatch sequence is **middle-out** for two reasons documented in the source spec:
1. **Connection pooling fairness** — first-dispatched requests get sockets first under degraded parallelism.
2. **Result-scan optimization** — descending sort of responders short-circuits at the first match.

Example for `currentVersion = 24`:
```
Dispatch order: 35, 34, 36, 33, 37, 32, 38, 31, 39, 30, 40, 29, 41, 28, 42, 27, 43, 26, 44, 25
```

### 15.7 Handoff Mechanism

When the probe finds a higher version `N > currentVersion`:

```bash
export INSTALL_PROBE_HANDOFF_DEPTH=$((depth + 1))
exec curl -fsSL "https://raw.githubusercontent.com/<owner>/<base>-v${N}/main/install.sh" | bash -s -- "$@"
```

The depth counter is **incremented before re-exec** so the child sees `depth+1`. If depth reaches `3`, the child aborts with exit code `1` and the message `Probe loop guard triggered — aborting.`

### 15.8 Logging Contract

Support scripts grep for these exact strings — do not change wording:

| Phase | Exact Log Line |
|-------|----------------|
| Probe disabled | `version probe disabled (no self-identity)` |
| Probe skipped (flag) | `version probe skipped via flag` |
| Probe found newer | `→ found higher version v<N>, handing off` |
| Probe found nothing | `no higher version found, installing local v<currentVersion>` |
| Loop guard triggered | `Probe loop guard triggered — aborting.` |

### 15.9 Golden Rule

> **The probe is a courtesy, not a gate.** Any error in the probe path falls back to "install whatever I am right now."

A blind AI must **never** add a hard-fail path inside the probe — every failure mode must degrade gracefully to local install.

---

*Install-Script Probe Runtime Contract added — v3.4.0 — 2026-04-22*
