# Build, Deploy & Run

## Overview

The project uses a single PowerShell script (`run.ps1`) at the repo root
to pull, build, deploy, and optionally run the gitmap CLI.
Build configuration lives in `gitmap/powershell.json`.

## Build Script — `run.ps1`

| Step | Description |
|------|-------------|
| 1. Git Pull | Pulls latest changes from remote |
| 2. Resolve Deps | Runs `go mod tidy` in `gitmap/` |
| 2b. Win Resources | Runs `go-winres make` to embed icon + metadata (Windows only, optional) |
| 3. Build | Compiles binary to `./bin/gitmap.exe` |
| 3b. Version | Runs the built binary with `version` and prints result |
| 4. Deploy | Copies binary + `data/` to deploy target (with retry on lock) |

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-NoPull` | Skip `git pull` | pull enabled |
| `-NoDeploy` | Skip deploy step | deploy enabled |
| `-DeployPath <dir>` | Override deploy directory | from `powershell.json` |
| `-Update` | Update mode: runs full pipeline (pull, build, deploy, sync) with post-update validation and cleanup | off |
| `-R` | Switch - run gitmap after build | off |
| *(trailing args)* | All args after `-R` are forwarded to gitmap | `scan <parent-folder>` |

### Examples

```powershell
# Full pipeline: pull, build, deploy
.\run.ps1

# Build only, no pull or deploy
.\run.ps1 -NoPull -NoDeploy

# Build and scan parent folder
.\run.ps1 -R scan

# Build and scan specific folder with SSH mode
.\run.ps1 -R scan D:\repos --mode ssh

# Build and clone from JSON
.\run.ps1 -R clone .\.gitmap/output\gitmap.json --target-dir .\restored

# Build and clone with GitHub Desktop registration
.\run.ps1 -R clone .\gitmap.json --github-desktop

# Deploy to custom path
.\run.ps1 -DeployPath "D:\tools"
```

## Configuration — `gitmap/powershell.json`

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "gitmap.exe",
  "copyData": true
}
```

| Field | Description | Default |
|-------|-------------|---------|
| `deployPath` | Directory where binary is deployed | `E:\bin-run` |
| `buildOutput` | Local build output directory | `./bin` |
| `binaryName` | Name of the compiled binary | `gitmap.exe` |
| `copyData` | Whether to copy `data/` alongside binary | `true` |

## Build Output

After a successful build, the `./bin/` directory contains:

```
bin/
├── gitmap.exe
└── data/
    └── config.json
```

## Deploy Target Resolution

The deploy target is resolved with a 3-tier priority:

| Priority | Source | Description |
|----------|--------|-------------|
| 1 | `-DeployPath` flag | Explicit CLI override — always wins |
| 2 | Global PATH lookup | If `gitmap` is already on PATH, deploy to its detected install directory |
| 3 | `powershell.json` | Fall back to `deployPath` from config (default `E:\bin-run`) |

This means **first-time installs** use the `powershell.json` default, but
**subsequent builds** automatically detect where gitmap is running from and
deploy there — no manual path configuration needed.

## Deploy Structure

The deploy target uses a nested `gitmap/` subfolder:

```
<deploy-target>\
└── gitmap\
    ├── gitmap.exe
    └── data\
        └── config.json
```

The `<deploy-target>\gitmap\` directory must be on the system `PATH` so
the user can run `gitmap` from any terminal.

## Rename-First Deploy Strategy

When the target binary already exists, the deploy step uses a
**rename-first** strategy to avoid file-lock failures (especially on
Windows, where a running `.exe` cannot be overwritten):

1. **Rename** the existing binary to `<binary>.old` (Windows allows
   renaming a running executable).
2. **Copy** the newly built binary into the now-free destination path.
3. If the copy fails after retries, **rollback** by renaming `.old`
   back to the original name.

```
existing gitmap.exe  →  gitmap.exe.old   (rename — succeeds even if locked)
new build bin/gitmap.exe  →  gitmap.exe  (copy — destination is free)
```

The `.old` file is left in place and cleaned up by
`gitmap update-cleanup`. On Linux/macOS, `mv` is used instead of
`Rename-Item`, providing identical behavior.

## Embedded Repo Path

The build step embeds the **absolute path of the source repo** into the
binary via Go `-ldflags`:

```powershell
$ldflags = "-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/constants.RepoPath=$absRepoRoot'"
go build -ldflags $ldflags -o $outPath .
```

This enables the `gitmap update` command to locate the source repo and
trigger a self-update without the user needing to know where the repo lives.

## `-R` Flag Behavior

`-R` is a **switch** parameter. All remaining positional arguments after it
are captured via `[Parameter(ValueFromRemainingArguments)]` into `$RunArgs`
and forwarded directly to the gitmap binary.

```powershell
param(
    [switch]$R,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
)
```

- If `-R` is used with no trailing arguments, it defaults to `scan <parent-folder>`.
- `-R` runs after build and deploy steps complete.

### Path Resolution

Relative path arguments (e.g., `..`, `../..`, `./projects`) are
automatically resolved to **absolute paths** before being passed to the
gitmap binary. Resolution uses `Resolve-Path` with a fallback to
`[System.IO.Path]::GetFullPath()` for paths that don't yet exist.

```powershell
# User runs:
.\run.ps1 -R scan "../.."

# Script resolves "../.." to absolute, e.g.:
# gitmap scan D:\wp-work
```

### RUN Context Logging

Before executing gitmap, the script prints diagnostic context:

```
  [RUN] Executing gitmap
  ──────────────────────────────────────────────────
  → Runner CWD: D:\wp-work\riseup-asia\gitmap-v28
  → Repo root: D:\wp-work\riseup-asia\gitmap-v28
  → Command: gitmap scan D:\wp-work
  → Scan target: D:\wp-work
  ──────────────────────────────────────────────────
```

| Line | Description |
|------|-------------|
| Runner CWD | Current working directory of the PowerShell session |
| Repo root | Root of the gitmap-v28 project |
| Command | Full command being executed |
| Scan target | Resolved absolute path passed to `scan` (shown only for scan commands) |

## Deploy Target

The deploy target is resolved via the 3-tier priority described in
**Deploy Target Resolution** above. The resolved directory contains a
`gitmap/` subfolder with the binary and data. That subfolder must be on
the system `PATH` so the tool can be run from any terminal.

## Logging

The script uses colored, step-numbered output:

- **Magenta** — step headers (`[1/4]`, `[2/4]`, etc.)
- **Green** — success messages (OK)
- **Cyan** — informational messages (->)
- **Yellow** — warnings (!!)
- **Red** — errors (XX)

## Version Display

After a successful build, the script immediately runs the new binary
with `version` and prints the result:

```
  -> Version: gitmap v1.1.2
```

This provides immediate confirmation that the build produced the
expected version.

## Deploy Retry

The deploy step retries the `Copy-Item` up to 20 times with a 500ms
delay between attempts if the target binary is locked by another process.
This handles the case where `gitmap update` may still be releasing its
file handle when deploy starts.

## Self-Update Flow (`gitmap update`)

1. `gitmap update` detects the active `gitmap` executable currently resolved by `PATH`.
2. It creates a handoff copy beside that active binary (same directory), such as `gitmap-update-<pid>.exe` (fallback to `%TEMP%` if locked).
3. It launches the handoff copy with the hidden `update-runner` command using **foreground/blocking** execution (`cmd.Run()`). The parent waits for the worker to complete so the terminal session stays stable. This is safe because the handoff copy is a different file — the parent's lock on the original binary is resolved by rename-first sync.
4. The handoff copy (`update-runner`) resolves the repo path and runs `run.ps1 -Update` from the repo root.
5. `run.ps1 -Update` performs the full pipeline: pull -> build -> deploy.
6. PATH sync uses **rename-first** strategy: renames active binary to `.old`, copies deployed binary to active path. Falls back to copy-retry loop (20 x 500ms) only if rename fails.
7. The updater prints executable-derived version comparison (`before` vs `after`) using `gitmap version`.
8. It runs `gitmap changelog --latest` using the updated binary.
9. It runs `gitmap update-cleanup` to remove temporary handoff and `.old` artifacts.

### Critical Rules

- Parent MUST use `cmd.Run()` (foreground/blocking), NEVER `cmd.Start()` + `os.Exit(0)` (async detach breaks terminal).
- PATH sync MUST use rename-first in update mode. Copy-overwrite fails on Windows when any process holds the binary.
- Generated scripts MUST NOT contain `Read-Host` or interactive prompts.

### Minimum Confirmation Output

- Active version before update
- Deployed version after update
- Final active version after sync (must match deployed)
- Last released version (from binary, `latest.json`, or git tag)
- Latest changelog entries from updated binary

## Last Release Detection — `Get-LastRelease.ps1`

A standalone PowerShell script at `gitmap/scripts/Get-LastRelease.ps1`
resolves and displays the latest released version. It is invoked
automatically by both `run.ps1` (after "All done!") and the generated
update script (in the version-verify block).

### Resolution Order

The script uses a three-tier fallback strategy:

| Priority | Source | Method |
|----------|--------|--------|
| 1 | Binary | `gitmap list-versions --limit 1` — parses first `vX.Y.Z` from output |
| 2 | JSON | `.gitmap/release/latest.json` — reads `tag` or `version` field |
| 3 | Git tag | `git tag --list "v*" --sort=-version:refname` — first stable `vX.Y.Z` |

### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `-BinaryPath` | string | `""` | Path to gitmap binary; falls back to `Get-Command gitmap` |
| `-RepoRoot` | string | `""` | Repo root for `latest.json` lookup; falls back to CWD |
| `-Label` | string | `"Last release"` | Display label prefix |

### Output

```
  Last release:    v2.24.0 (binary)
```

The parenthetical suffix indicates which source resolved the version:
`binary`, `latest.json`, or `git tag`. If all sources fail, prints
`unknown`.

### Integration Points

**`run.ps1`** — called after the final "All done!" message:

```powershell
$lastReleaseScript = Join-Path (Join-Path (Join-Path $RepoRoot "gitmap") "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    & $lastReleaseScript -BinaryPath $changelogBinaryPath -RepoRoot $RepoRoot
}
```

**Update script (`constants_update.go`)** — embedded in the
`UpdatePSVerify` section between the version lines and the
active/deployed match check:

```powershell
$lastReleaseScript = Join-Path (Join-Path (Join-Path "<repoPath>" "gitmap") "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    & $lastReleaseScript -BinaryPath $activeBinary -RepoRoot "<repoPath>"
}
```

### Design Decisions

- **Separate file** keeps `run.ps1` lean and allows reuse from any
  context (manual invocation, CI, update scripts).
- **Three-tier fallback** ensures a result even when the binary is
  unavailable (fresh clone) or when `.gitmap/release/` metadata hasn't been
  generated yet.
- **No error exits** — the script always succeeds; missing data simply
  shows `unknown`.

## Cross-References (Generic Specifications)

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Build pipeline | [04-build-scripts.md](../08-generic-update/04-build-scripts.md) | `run.ps1` / `run.sh` full pipeline, config loading, ldflags |
| Deploy strategy | [03-rename-first-deploy.md](../08-generic-update/03-rename-first-deploy.md) | Rename-first flow, rollback, PATH sync, retry reduction |
| Deploy path resolution | [02-deploy-path-resolution.md](../08-generic-update/02-deploy-path-resolution.md) | 3-tier deploy target resolution (CLI → PATH → config) |
| Self-update overview | [01-self-update-overview.md](../08-generic-update/01-self-update-overview.md) | Platform behavior, update strategies, version comparison |
| Handoff mechanism | [05-handoff-mechanism.md](../08-generic-update/05-handoff-mechanism.md) | Copy-and-handoff, worker launch, foreground blocking |
| Cleanup | [06-cleanup.md](../08-generic-update/06-cleanup.md) | `.old` lifecycle, `update-cleanup` command |
| PowerShell patterns | [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Script architecture, config, logging, deploy, self-update |
| Self-update mechanism | [03-self-update-mechanism.md](../03-general/03-self-update-mechanism.md) | Three-layer approach, skip-if-current, error diagnostics |
| Icon embedding | [04-windows-icon-embedding.md](../03-general/04-windows-icon-embedding.md) | go-winres integration, `.syso` generation |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
