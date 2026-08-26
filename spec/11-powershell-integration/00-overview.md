# PowerShell Integration for Project Runner

> **Spec Version:** 2.25.0  
> **Script Version:** 2.25.0  
> **Updated:** 2026-03-19  
> **Status:** Active  
> **Location:** `spec/powershell-integration/`  
> **Purpose:** Reusable PowerShell runner for Go backend + React frontend projects with pnpm PnP support

---

## Summary

This specification defines a **cross-project reusable** PowerShell integration pattern for building and running fullstack applications with Go backend and React frontend. The system uses a JSON configuration file (`powershell.json`) to define project-specific paths and settings.

**Key Features:**
- **pnpm Plug'n'Play (PnP)** - Disk-efficient package management with shared store
- **Relative Path Resolution** - All paths relative to script location (working directory)
- **Force Reinstall** - Clear caches and reset everything with `-Force` flag
- **Multi-Project Root Folder** - Shared pnpm store across Node.js projects

**This spec is NOT project-specific** — it can be used by:
- WP Plugin Publish
- Spec Management Software
- Any Go + React fullstack project

---

## User Stories

- As a developer, I want to run a single command to build and start my fullstack app
- As a developer, I want clean build options to reset everything when needed
- As a developer, I want the script to auto-install missing dependencies (Go, Node.js, pnpm)
- As a developer, I want to configure paths via JSON instead of editing the script
- As a developer, I want firewall rules configured automatically for development
- As a developer, I want pnpm PnP to save disk space across multiple projects

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     PowerShell Runner Architecture v2.0                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                  │
│   │   run.ps1    │───▶│ powershell.  │───▶│   Project    │                  │
│   │   (Script)   │    │ json config  │    │   Folders    │                  │
│   └──────────────┘    └──────────────┘    └──────────────┘                  │
│          │                   │                    │                          │
│          │                   ▼                    ▼                          │
│          │           ┌──────────────┐    ┌──────────────┐                   │
│          │           │  pnpm Store  │    │  Go Backend  │                   │
│          │           │  (Shared)    │    │  + React FE  │                   │
│          │           └──────────────┘    └──────────────┘                   │
│          ▼                                                                   │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                         Pipeline Steps                               │   │
│   │  1. Git Pull → 2. Prerequisites → 3. pnpm Install → 4. Build → 5. Run│  │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Pipeline Steps

| Step | Name | Description | Flags |
|------|------|-------------|-------|
| 1 | Git Pull | Sync latest changes | `-SkipPull` to skip |
| 2 | Prerequisites | Check/install Go, Node.js, pnpm | Auto-install via winget |
| 3 | pnpm Install | Install dependencies with PnP | `-Force` clears store & reinstalls |
| 4 | Frontend Build | Build React with pnpm | `-SkipBuild` to skip |
| 5 | Copy & Run | Copy dist, start Go server | `-BuildOnly` to skip run |

---

## Package Management: pnpm with Plug'n'Play

### Why pnpm PnP?

| Feature | npm | pnpm PnP |
|---------|-----|----------|
| Disk Usage | Full copy per project | Shared store, hard links |
| Install Speed | Moderate | Fast (cached) |
| node_modules | Required (~500MB+) | Not required |
| Deterministic | package-lock.json | pnpm-lock.yaml |

### Configuration

```json
{
  "usePnp": true,
  "pnpmStorePath": "E:/.pnpm-store"
}
```

- `usePnp: true` - Enable pnpm with PnP mode
- `pnpmStorePath` - Custom store location (relative to rootDir or absolute)

### Store Path Options

| Option | Path | Description |
|--------|------|-------------|
| **Default (Recommended)** | `E:/.pnpm-store` | Shared drive for all projects |
| **Relative (Isolated)** | `.pnpm-store` | Store in project root |
| **User Home** | `~/.pnpm-store` | Global store in user home |

---

## Folder Structure

```
spec/powershell-integration/
├── 00-overview.md               ← This file
├── 01-configuration-schema.md   ← JSON config format with pnpm options
├── 02-script-reference.md       ← CLI flags and functions
├── 03-integration-guide.md      ← How to add to any project
├── 04-error-codes.md            ← Exit codes (9500-9599)
├── 05-firewall-rules.md         ← Windows firewall setup
├── schemas/
│   └── powershell.schema.json   ← JSON Schema for validation
├── templates/
│   ├── run.ps1                  ← Main script template
│   └── powershell.json          ← Example config with pnpm
└── examples/
    └── server-client-project.json  ← Sample for server/client layout

spec/upload-scripts/              ← Related: WordPress plugin upload scripts
├── readme.md                    ← Upload pipeline overview
├── 01-upload-plugin-v1.md       ← V1: Basic single-file upload
├── 02-upload-plugin-v2.md       ← V2: Envelope-aware upload
├── 03-upload-plugin-v3.md       ← V3: Parallel multi-plugin deployment
├── 04-upload-plugin-custom.md   ← Custom path deployments
└── 05-configuration.md          ← Auth, headers, fallback config
```

---

## Quick Start

```powershell

# Full build and run (pnpm PnP enabled)

.\run.ps1

# Clean rebuild everything (clears pnpm store cache)

.\run.ps1 -Force

# Just start backend (skip frontend build)

.\run.ps1 -SkipBuild

# Build only (don't start server)

.\run.ps1 -BuildOnly

# Skip git pull + clean build

.\run.ps1 -SkipPull -Force

# Configure firewall (requires Admin)

.\run.ps1 -OpenFirewall

# Show help

.\run.ps1 -Help
```

---

## Configuration File

Create `powershell.json` in project root:

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
    "node_modules",
    "dist",
    ".vite",
    ".pnp.cjs",
    ".pnp.loader.mjs",
    "backend/data/*.db"
  ],
  "buildCommand": "pnpm run build",
  "installCommand": "pnpm install",
  "runCommand": "go run cmd/server/main.go",
  "configFile": "config.json",
  "configExampleFile": "config.example.json"
}
```

---

## Features

### Auto-Install Dependencies

- **Go**: Installs via `winget install GoLang.Go` if missing
- **Node.js**: Installs via `winget install OpenJS.NodeJS.LTS` if missing
- **pnpm**: Installs via `npm install -g pnpm` if missing

### Force Clean Build

The `-Force` flag removes:
- `.pnp.cjs` and `.pnp.loader.mjs` files
- `node_modules/` directory (if exists)
- `dist/` directory
- `.vite/` cache
- SQLite databases (`*.db`, `*.db-shm`, `*.db-wal`)
- Prunes pnpm store cache

### Required .gitignore Entries

**IMPORTANT:** Add these entries to your `.gitignore` to exclude pnpm artifacts from version control:

```gitignore

# pnpm store (local cache)

.pnpm-store/

# pnpm PnP files (generated)

.pnp.cjs
.pnp.loader.mjs

# Build artifacts

dist/
.vite/
```

### pnpm Store Management

```powershell

# Check store status

pnpm store status

# Prune unused packages

pnpm store prune

# View store path

pnpm store path
```

### Firewall Configuration

The `-OpenFirewall` flag (requires Administrator):
- Creates inbound rules for configured ports
- Sets profile to Private and Domain
- Names rules consistently for easy management

---

## Path Resolution

All paths are resolved relative to the script location (`$MyInvocation.MyCommand.Path`).

```
project-root/           ← Working directory (where run.ps1 lives)
├── run.ps1             ← Script location (rootDir base)
├── powershell.json     ← Config file
├── package.json        ← Frontend (frontendDir: ".")
├── .pnp.cjs            ← PnP resolution (generated by pnpm)
├── .pnpm-store/        ← pnpm store (pnpmStorePath)
├── dist/               ← Build output (distDir: "dist")
└── backend/            ← Backend (backendDir: "backend")
    ├── cmd/server/main.go
    ├── config.json
    ├── config.example.json
    ├── frontend/
    │   └── dist/       ← Target (targetDir)
    └── data/           ← Data (dataDir)
        └── *.db
```

---

## Using in Projects

### For New Projects

1. Copy `templates/run.ps1` to project root
2. Create `powershell.json` with project-specific paths
3. Set `usePnp: true` and configure `pnpmStorePath`
4. Run `.\run.ps1 -Help` to verify

### Multi-Project Setup (Shared Store)

For multiple projects sharing a pnpm store:

```json
{
  "pnpmStorePath": "E:/.pnpm-store"
}
```

All projects pointing to the same store share cached packages.

---

## AI Handoff Instructions

To integrate this PowerShell runner into any project, share:

```
spec/powershell-integration/
```

Tell the AI:
> "Follow the spec at `spec/powershell-integration/` to add the PowerShell build runner. Create a `powershell.json` config for my project structure. Enable pnpm PnP for disk-efficient package management."

---

## Cross-References

| Document | Description |
|----------|-------------|
| [Configuration Schema](./01-configuration-schema.md) | JSON config format with pnpm options |
| [Script Reference](./02-script-reference.md) | CLI flags and functions |
| [Integration Guide](./03-integration-guide.md) | Step-by-step setup |
| [Error Codes](./04-error-codes.md) | Exit codes 9500-9599 |
| [Firewall Rules](./05-firewall-rules.md) | Windows firewall setup |
| Upload Scripts Spec | WordPress plugin upload scripts (V1, V2, V3) — *folder pending creation* |
| Upload V1 | Single-file upload via Invoke-RestMethod — *folder pending creation* |
| Upload V2 | Envelope-aware upload with unwrapping — *folder pending creation* |
| Upload V3 | Parallel multi-plugin deployment via Start-Job — *folder pending creation* |
| Upload Custom | Custom path deployments via `run.ps1 -u -pp` — *folder pending creation* |
| Upload Config | Authentication, headers, and fallback config — *folder pending creation* |

---

*This spec enables consistent, reproducible builds across all fullstack projects with optimized disk usage via pnpm PnP.*

---

## Verification

_Auto-generated section — see `spec/11-powershell-integration/97-acceptance-criteria.md` for the full criteria index._

### AC-PS-000: PowerShell integration conformance: Overview

**Given** Lint PowerShell scripts and modules in `scripts/` for naming, parameter binding, and error propagation.  
**When** Run the verification command shown below.  
**Then** Filenames are lowercase-kebab-case; functions are `Verb-Noun` PascalCase; `$ErrorActionPreference = 'Stop'` is set; no `Write-Host` for control flow.

**Verification command:**

```bash
pwsh -NoProfile -Command "Invoke-ScriptAnalyzer -Path scripts -Recurse -Severity Warning"
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
