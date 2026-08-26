# Consolidated: PowerShell Integration — Complete Reference

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for the cross-project reusable PowerShell runner. An AI reading only this file must be able to set up the runner, use all CLI flags, handle errors, configure WordPress plugin deployment, and integrate with CI/CD — in any Go + React project.

**Not project-specific** — usable by any Go + React fullstack project.

---

## Architecture

```
run.ps1 → powershell.json (config) → Project Folders
  Pipeline: Git Pull → Prerequisites → pnpm Install → Build → Copy → Run
```

---

## Pipeline Steps

| Step | Name | Description | Skip Flag |
|------|------|-------------|-----------|
| 1 | Git Pull | Sync latest changes | `-SkipPull` |
| 2 | Prerequisites | Check/install Go, Node.js, pnpm | Auto via winget |
| 3 | pnpm Install | Install dependencies with PnP | `-Force` clears & reinstalls |
| 4 | Frontend Build | Build React with pnpm | `-SkipBuild` |
| 5 | Copy Build | Copy dist to backend target dir | — |
| 6 | Start Backend | Create config if missing, run Go server | `-BuildOnly` to skip |

---

## CLI Flags — Complete Reference

| Short | Long | Type | Description |
|-------|------|------|-------------|
| `-h` | `-help` | Switch | Show help message and exit |
| `-b` | `-buildonly` | Switch | Build frontend only, don't start backend |
| `-s` | `-skipbuild` | Switch | Skip frontend build, only run backend |
| `-p` | `-skippull` | Switch | Skip git pull step |
| `-f` | `-force` | Switch | Force-clean build artifacts and pnpm folders before building |
| `-i` | `-install` | Switch | Install/update dependencies (frontend + backend), then exit |
| `-r` | `-rebuild` | Switch | Full reset: clean → install → build/run (install deferred until after clean) |
| `-fw` | `-openfirewall` | Switch | Add Windows Firewall rules (requires Admin) |
| `-u` | `-upload` | Switch | Upload default plugin to WordPress via upload-plugin-v2 |
| `-q` | `-qupload` | Switch | Upload plugin via QUpload API (upload-plugin-U-Q.ps1) |
| | `-u -q` | Combo | Upload Riseup Asia Uploader via QUpload API (shorthand) |
| `-ua` | `-uploadall` | Switch | ZIP all plugins (except QUpload) and upload each via QUpload API |
| | `-ua -xs 'slug'` | Combo | ZIP + upload all except named plugin(s), comma-separated |
| `-z` | `-zip` | Switch | ZIP default plugin (or specific via `-pp`) |
| `-za` | | Switch | ZIP all plugins in `wp-plugins/` with version numbers |
| `-zq` | `-zipqupload` | Switch | ZIP QUpload plugin |
| `-c` | `-clear` | Switch | Remove all existing ZIP files before zipping |
| | `-uas` | Switch | Upload ALL plugins to ALL configured sites (multi-site) |
| | `-uas -site 'name'` | Combo | Upload ALL plugins to a specific site by name |
| | `-uas -xs 'name'` | Combo | Upload ALL plugins to all sites EXCEPT named one(s) |
| `-t` | `-test` | Switch | Run Go backend tests and exit |
| `-pp` | `-pluginpath` | String | Override plugin folder path (use with `-u`, `-q`, `-ua`, or `-z`) |
| | `-site` | String | Target a specific site by name (use with `-uas`) |
| `-xs` | `-exclude` | String | Exclude site(s) by name, comma-separated (use with `-uas`) |
| `-d` | `-deploy` | Switch | Full deploy: git pull → upload all sites → plugin status → build & run |
| `-dbg` | `-debug` | Switch | Enable debug logging for upload |
| `-v` | `-verbose` | Switch | Show detailed debug output |

### Usage Examples

```powershell
.\run.ps1                    # Full build and run
.\run.ps1 -f                 # Clean rebuild (clears pnpm store cache)
.\run.ps1 -s                 # Start backend only (skip build)
.\run.ps1 -b                 # Build only, don't start server
.\run.ps1 -p                 # Skip git pull
.\run.ps1 -i                 # Install/update all dependencies, then exit
.\run.ps1 -r                 # Full reset: clean + install + build/run
.\run.ps1 -fw                # Configure firewall (requires Admin)
.\run.ps1 -u                 # Upload default plugin (V2 uploader)
.\run.ps1 -q                 # Upload plugin via QUpload API
.\run.ps1 -ua                # ZIP + upload all plugins via QUpload API
.\run.ps1 -ua -c             # Clear old ZIPs, then ZIP + upload all
.\run.ps1 -uas               # Multi-site: upload all plugins to all enabled sites
.\run.ps1 -uas -site "Test"  # Multi-site: upload to specific site
.\run.ps1 -uas -xs "Test"    # Multi-site: exclude specific site
.\run.ps1 -t                 # Run Go backend tests
.\run.ps1 -d                 # Full deploy cycle
.\run.ps1 -v                 # Verbose output
```

---

## Configuration — `powershell.json`

### Standard Config (Recommended)

```json
{
  "$schema": "./spec/powershell-integration/schemas/powershell.schema.json",
  "version": "1.1.0",
  "projectName": "WP Plugin Publish",
  "rootDir": ".",
  "backendDir": "backend",
  "frontendDir": ".",
  "distDir": "dist",
  "targetDir": "backend/frontend/dist",
  "dataDir": "backend/data",
  "ports": [8080],
  "prerequisites": {
    "go": true,
    "node": true,
    "pnpm": true
  },
  "usePnp": true,
  "pnpmStorePath": "E:/.pnpm-store",
  "cleanPaths": [
    "node_modules", "dist", ".vite",
    ".pnp.cjs", ".pnp.loader.mjs",
    "backend/data/*.db"
  ],
  "buildCommand": "pnpm run build",
  "installCommand": "pnpm install",
  "runCommand": "go run cmd/server/main.go",
  "configFile": "config.json",
  "configExampleFile": "config.example.json"
}
```

### Minimal Config

```json
{
  "projectName": "My Project",
  "backendDir": "backend"
}
```

Uses all defaults: pnpm PnP with `.pnpm-store` in project root.

### All Config Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `projectName` | string | **required** | Display name for console output |
| `rootDir` | string | `"."` | Project root relative to script |
| `backendDir` | string | **required** | Path to Go backend directory |
| `frontendDir` | string | `"."` | Path to React frontend directory |
| `distDir` | string | `"dist"` | Frontend build output directory |
| `targetDir` | string | — | Where to copy built frontend for Go serving |
| `dataDir` | string | — | Data directory for SQLite databases |
| `ports` | int[] | `[8080]` | Ports for firewall rules |
| `prerequisites` | object | `{go:true, node:true, pnpm:true}` | Which prereqs to check/install |
| `usePnp` | boolean | `true` | Enable pnpm Plug'n'Play |
| `pnpmStorePath` | string | `"E:/.pnpm-store"` | pnpm shared store path |
| `cleanPaths` | string[] | `["node_modules","dist",".vite",".pnp.cjs"]` | Paths to remove on `-Force` |
| `buildCommand` | string | `"pnpm run build"` | Frontend build command |
| `installCommand` | string | `"pnpm install"` | Frontend install command |
| `runCommand` | string | `"go run cmd/server/main.go"` | Backend start command |
| `configFile` | string | `"config.json"` | Backend config file name |
| `configExampleFile` | string | `"config.example.json"` | Template config to copy |
| `requiredModules` | string[] | `[]` | NPM modules that trigger reinstall if missing |
| `env` | object | — | Environment variables to set before commands |
| `seedingDir` | string | — | Directory for seed data files |

---

## pnpm Plug'n'Play (PnP)

### Why pnpm PnP?

| Feature | npm | pnpm PnP |
|---------|-----|----------|
| Disk Usage | Full copy per project | Shared store, hard links |
| Install Speed | Moderate | Fast (cached) |
| node_modules | Required (~500MB+) | Not required |
| Deterministic | package-lock.json | pnpm-lock.yaml |

### Store Path Options

| Option | Path | Description |
|--------|------|-------------|
| Shared (Recommended) | `E:/.pnpm-store` | Shared drive for all projects |
| Relative (Isolated) | `.pnpm-store` | Store in project root |
| User Home | `~/.pnpm-store` | Global store in user home |

### Multi-Project Disk Savings

| Scenario | npm | pnpm (shared) |
|----------|-----|---------------|
| 1 project | ~500MB | ~300MB |
| 5 projects | ~2.5GB | ~400MB |
| **Savings** | — | **~80%** |

### Node 24+ / Cross-Drive Fallback

If `usePnp` is enabled, the runner falls back to `node-linker=isolated` when:
- Node.js major version is **24+**
- The pnpm store is on a **different drive** than the project

This avoids `ERR_MODULE_NOT_FOUND` failures (e.g., Vite failing to resolve `esbuild`).

### pnpm v10+ Build Scripts

pnpm v10+ blocks dependency build scripts by default. The runner auto-appends `--dangerously-allow-all-builds` when pnpm v10+ is detected, ensuring native deps like esbuild/@swc work for Vite.

---

## Functions Reference

### Format-ElapsedTime

Formats a Stopwatch elapsed time. Output: `2.3s` or `1m 45.2s`.

```powershell
function Format-ElapsedTime($Stopwatch) {
    $elapsed = $Stopwatch.Elapsed
    if ($elapsed.TotalMinutes -ge 1) {
        return "{0:N0}m {1:N1}s" -f [Math]::Floor($elapsed.TotalMinutes), $elapsed.Seconds
    } else {
        return "{0:N1}s" -f $elapsed.TotalSeconds
    }
}
```

### Test-Command

Checks if a command exists in PATH.

```powershell
function Test-Command($Command) {
    try { if (Get-Command $Command) { return $true } }
    catch { return $false }
}
```

### Test-IsAdmin

Checks if running with Administrator privileges.

```powershell
function Test-IsAdmin {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}
```

### Install-Pnpm / Install-NodeJS / Install-Go

Auto-install via npm or winget with PATH refresh:

```powershell
function Install-Pnpm {
    npm install -g pnpm
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + 
                [System.Environment]::GetEnvironmentVariable("Path","User")
}

function Install-NodeJS {
    winget install OpenJS.NodeJS.LTS --accept-package-agreements --accept-source-agreements
    # Refresh PATH
}

function Install-Go {
    winget install GoLang.Go --accept-package-agreements --accept-source-agreements
    # Refresh PATH
}
```

### Ensure-FirewallRules

Creates Windows Firewall inbound rules (requires Admin):

```powershell
function Ensure-FirewallRules {
    param([int[]]$Ports = @(8080))
    foreach ($p in $Ports) {
        $ruleName = "$ProjectName (Go Backend) TCP $p"
        $existing = Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue
        if ($null -eq $existing) {
            New-NetFirewallRule -DisplayName $ruleName -Direction Inbound `
                -Action Allow -Protocol TCP -LocalPort $p -Profile Private,Domain | Out-Null
        }
    }
}
```

---

## Pipeline Steps — Detail

### Step 1: Git Pull

- Skipped if `-SkipPull`
- Warns but continues if git pull fails
- Skips if not a git repository

### Step 2: Prerequisites Check

```powershell
if ($config.prerequisites.go -and -not (Test-Command "go")) { Install-Go }
if ($config.prerequisites.node -and -not (Test-Command "node")) { Install-NodeJS }
if ($config.prerequisites.pnpm -and -not (Test-Command "pnpm")) { Install-Pnpm }
```

### Step 3: pnpm Install (PnP Mode)

- Configures store path from `pnpmStorePath`
- `-Force` removes: `node_modules`, `.pnpm`, `.pnp.cjs`, `.pnp.loader.mjs`, `.pnp.data.json`, then `pnpm store prune`
- `-Force` also cleans backend runtime data: `sessions/`, `request-sessions/`, `errors/`, `log.txt`, `error.log.txt`
- `-rebuild` defers install until AFTER force-clean
- Install detection: PnP checks `.pnp.cjs`, isolated checks `node_modules`
- `-i` and `-r` flags always trigger install even if deps appear present

### Step 4: Frontend Build

- Runs `buildCommand` (default: `pnpm run build`)
- PnP loader options handled automatically when `node-linker=pnp` is active

### Step 5: Copy Build

- Removes old target directory
- Copies `distDir` → `targetDir`

### Step 6: Start Backend

- Creates `configFile` from `configExampleFile` if missing
- Creates `data/` directory if missing
- Runs `runCommand` (default: `go run cmd/server/main.go`)

---

## Force Clean Build (`-Force`)

Removes:
- `.pnp.cjs` and `.pnp.loader.mjs` files
- `node_modules/` directory
- `dist/` directory
- `.vite/` cache
- SQLite databases (`*.db`, `*.db-shm`, `*.db-wal`)
- Backend runtime data (sessions, error logs)
- Prunes pnpm store cache

---

## Error Codes — Complete Reference

### Exit Codes (Quick)

| Code | Name | Description |
|------|------|-------------|
| 0 | SUCCESS | All steps completed |
| 1 | ERR_PREREQUISITES | Prerequisites installation failed |
| 2 | ERR_NPM_INSTALL | Install failed |
| 3 | ERR_NPM_BUILD | Build failed |
| 4 | ERR_GO_RUN | Go backend failed |
| 5 | ERR_CONFIG_MISSING | powershell.json not found |
| 6 | ERR_CONFIG_INVALID | JSON parse error |
| 7 | ERR_PATH_NOT_FOUND | Configured path doesn't exist |
| 8 | ERR_COPY_FAILED | Build copy failed |
| 9 | ERR_GIT_FAILED | Git pull failed critically |
| 10 | ERR_FIREWALL | Firewall configuration failed |

### Detailed Error Codes (9500–9599 range)

#### 9500 — Configuration Errors

| Code | Name | Message |
|------|------|---------|
| 9500 | ERR_CONFIG_NOT_FOUND | powershell.json not found in project root |
| 9501 | ERR_CONFIG_PARSE | Failed to parse powershell.json: {details} |
| 9502 | ERR_CONFIG_MISSING_FIELD | Required field '{field}' missing |
| 9503 | ERR_CONFIG_INVALID_PATH | Path '{path}' does not exist |
| 9504 | ERR_CONFIG_INVALID_PORT | Port must be 1–65535 |

#### 9510 — Prerequisites Errors

| Code | Name | Message |
|------|------|---------|
| 9510 | ERR_WINGET_NOT_FOUND | winget not available for auto-install |
| 9511 | ERR_GO_INSTALL_FAILED | Failed to install Go via winget |
| 9512 | ERR_NODE_INSTALL_FAILED | Failed to install Node.js via winget |
| 9513 | ERR_GO_NOT_IN_PATH | Go installed but not in PATH — restart required |
| 9514 | ERR_NPM_NOT_IN_PATH | npm installed but not in PATH — restart required |

#### 9520 — Build Errors

| Code | Name | Message |
|------|------|---------|
| 9520 | ERR_NPM_INSTALL_FAILED | Install failed with exit code {code} |
| 9521 | ERR_NPM_BUILD_FAILED | Build failed with exit code {code} |
| 9522 | ERR_DIST_NOT_CREATED | Build completed but dist not found |
| 9523 | ERR_COPY_DIST_FAILED | Failed to copy dist to target |
| 9524 | ERR_CLEAN_FAILED | Failed to remove {path} during force clean |

#### 9530 — Backend Errors

| Code | Name | Message |
|------|------|---------|
| 9530 | ERR_BACKEND_DIR_NOT_FOUND | Backend directory not found |
| 9531 | ERR_MAIN_GO_NOT_FOUND | main.go not found |
| 9532 | ERR_GO_BUILD_FAILED | go build failed with exit code {code} |
| 9533 | ERR_GO_RUN_FAILED | go run failed with exit code {code} |
| 9534 | ERR_CONFIG_COPY_FAILED | Failed to copy config.example.json |

#### 9540 — Firewall Errors

| Code | Name | Message |
|------|------|---------|
| 9540 | ERR_NOT_ADMIN | Requires Administrator privileges |
| 9541 | ERR_FIREWALL_CMDLET | New-NetFirewallRule not available |
| 9542 | ERR_FIREWALL_RULE_FAILED | Failed to create rule for port {port} |

#### 9550 — Git Errors

| Code | Name | Message |
|------|------|---------|
| 9550 | ERR_NOT_GIT_REPO | Not a git repository (skipping pull) |
| 9551 | ERR_GIT_PULL_FAILED | git pull failed with exit code {code} |
| 9552 | ERR_GIT_CONFLICT | git pull aborted due to conflicts |

### Error Handling Pattern

```powershell
$ErrorActionPreference = "Stop"
try {
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR [9521]: npm build failed with exit code $LASTEXITCODE" -ForegroundColor Red
        exit 3
    }
} catch {
    Write-Host "ERROR [9521]: $($_.Exception.Message)" -ForegroundColor Red
    exit 3
}
```

### Graceful Degradation (Non-Critical)

```powershell
git pull
if ($LASTEXITCODE -ne 0) {
    Write-Host "WARNING [9551]: git pull failed, continuing..." -ForegroundColor Yellow
}
```

### Console Output Levels

| Level | Color | Usage |
|-------|-------|-------|
| Info | Cyan | Step headers, status |
| Success | Green | Completion messages |
| Warning | Yellow | Non-fatal issues |
| Error | Red | Fatal errors |
| Debug | Gray | Timing, paths |

---

## Firewall Configuration

### Automatic (`-OpenFirewall`, requires Admin)

Creates inbound rules for all ports in config:

| Property | Value |
|----------|-------|
| DisplayName | `{ProjectName} (Go Backend) TCP {port}` |
| Direction | Inbound |
| Action | Allow |
| Protocol | TCP |
| Profile | Private, Domain (never Public) |

### Profile Security

| Profile | When Active | Recommendation |
|---------|-------------|----------------|
| Domain | Corporate network | ✅ Allow |
| Private | Home/trusted | ✅ Allow |
| Public | Coffee shop, airport | ❌ Block |

### Verification

```powershell
Get-NetFirewallRule -DisplayName "LLM Runner*"
Test-NetConnection -ComputerName YOUR_IP -Port 8080
```

### Removal

```powershell
Get-NetFirewallRule -DisplayName "LLM Runner*" | Remove-NetFirewallRule
```

---

## WordPress Plugin Management

### Config — `wpPlugins` section

```json
{
  "wpPlugins": {
    "defaultUploader": "riseup-asia-uploader",
    "defaultQUploader": "qupload",
    "pluginsDir": "wp-plugins",
    "skipPlugins": [],
    "phpCheckSkipFolders": ["vendor"],
    "sites": [...],
    "plugins": {
      "my-plugin": {
        "name": "My Plugin",
        "path": "wp-plugins/my-plugin",
        "mainFile": "my-plugin.php",
        "autoUpload": false
      }
    }
  }
}
```

### Multi-Site Deployment (`-uas`)

Upload plugins to multiple WordPress sites from a single command:

```json
{
  "wpPlugins": {
    "sites": [
      {
        "name": "Test V1",
        "url": "https://testv1.example.com",
        "enabled": true,
        "credentials": [
          {
            "appName": "deployer",
            "usernameBase64": "<base64>",
            "passwordBase64": "<base64>",
            "isDefault": true
          }
        ]
      }
    ]
  }
}
```

### Credential Encoding

```powershell
[Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("my-username"))
```

### Multi-Site Flow

1. Validates `wpPlugins.sites` configuration
2. Lists all configured sites (shows enabled/disabled)
3. Filters by `-site` name or uses all enabled sites
4. Discovers uploadable plugins (respects `skipPlugins`)
5. ZIPs all plugins once
6. For each target site: resolves default credential → decodes Base64 → uploads via QUpload API
7. Prints multi-site summary with per-site per-plugin results

---

## PHP Known Issues

### QUpload Syntax Validator False Positives

`token_get_all($content, TOKEN_PARSE)` produces false positives for certain PHP patterns:

| Pattern | Error | Fix |
|---------|-------|-----|
| `is_array($var)` | `unexpected token "array"` | `gettype($var) === PhpNativeType::PhpArray->value` |
| `array()` long syntax | `unexpected token "array"` | Use `[]` short syntax |
| `= array()` default | `unexpected token "array"` | Use `= []` |

### Template Detection Skip Logic

Files skipped (not validated) when they start with `<?php` AND contain `?>`. Pure PHP files (no closing tag) are always validated.

### PhpNativeType Enum

| Case | Value | Replaces |
|------|-------|----------|
| `PhpArray` | `'array'` | `is_array()` |
| `PhpString` | `'string'` | `is_string()` |
| `PhpInteger` | `'integer'` | `is_int()` |
| `PhpDouble` | `'double'` | `is_float()` |
| `PhpBoolean` | `'boolean'` | `is_bool()` |
| `PhpObject` | `'object'` | `is_object()` |
| `PhpNull` | `'NULL'` | `is_null()` |

---

## Path Resolution

All paths resolved relative to script location (`$MyInvocation.MyCommand.Path`):

```
project-root/
├── run.ps1              ← Script location
├── powershell.json      ← Config
├── package.json         ← Frontend
├── pnpm-lock.yaml       ← Lockfile (commit this!)
├── .pnp.cjs             ← PnP resolution (generated, gitignored)
├── .pnpm-store/         ← pnpm store (generated, gitignored)
├── dist/                ← Build output
├── src/                 ← React source
└── backend/
    ├── cmd/server/main.go
    ├── config.json
    ├── config.example.json
    └── frontend/dist/   ← Target (copied dist)
```

---

## Required .gitignore Entries

```gitignore
.pnpm-store/
.pnp.cjs
.pnp.loader.mjs
dist/
.vite/
```

---

## CI/CD Integration

### GitHub Actions

```yaml
- name: Setup pnpm
  uses: pnpm/action-setup@v2
  with:
    version: 8

- name: Build Frontend
  shell: pwsh
  run: .\run.ps1 -BuildOnly -SkipPull
```

### Azure DevOps

```yaml
- task: PowerShell@2
  inputs:
    filePath: 'run.ps1'
    arguments: '-BuildOnly -SkipPull'
```

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| pnpm not found after install | Restart PowerShell or refresh PATH manually |
| PnP module not found errors | `.\run.ps1 -Force` |
| pnpm install fails (network/corrupt) | `pnpm store prune && pnpm cache clean && .\run.ps1 -Force` |
| Firewall rules not applied | Run as Administrator |
| Build takes too long | `.\run.ps1 -SkipBuild` (skip if no frontend changes) |

---

## Setup for New Projects

1. Copy `templates/run.ps1` to project root
2. Create `powershell.json` with project-specific paths
3. Set `usePnp: true` and configure `pnpmStorePath`
4. Update `.gitignore` with required entries
5. Run `.\run.ps1 -Help` to verify

---

## Timing Output Example

```
========================================
  WP Plugin Publish - Build & Run Script
========================================

[1/5] Pulling latest changes from git...
  ✓ Git pull complete
  ⏱ 1.2s

[2/5] Checking prerequisites...
  ✓ Go found: go version go1.21.0 windows/amd64
  ✓ Node.js found: v20.10.0
  ✓ pnpm found: 8.12.0
  ⏱ 0.3s

[3/5] Installing dependencies (pnpm PnP)...
  Store path: .pnpm-store
  ✓ Dependencies installed
  ⏱ 5.2s

[4/5] Building React frontend...
  ✓ Frontend built successfully
  ⏱ 12.5s

[5/5] Starting Go backend...
========================================
  WP Plugin Publish starting...
  Open: http://localhost:8080
  Press Ctrl+C to stop
  Build time: 19.2s
========================================
```

---

## Prerequisites

| Requirement | Minimum Version | Auto-Install |
|-------------|-----------------|--------------|
| Windows | 10/11 or Server 2019+ | N/A |
| PowerShell | 5.1 or 7+ | N/A |
| winget | Latest | N/A |
| Go | 1.21+ | ✅ via winget |
| Node.js | 18+ LTS | ✅ via winget |
| pnpm | 8+ | ✅ via npm |

---

*Consolidated PowerShell integration — v3.2.0 — 2026-04-16*
