# PowerShell Integration Guide

> **Spec Version:** 2.2.0
> **Script Version:** run.ps1 1.2.0, upload-plugin-v2.ps1 2.1.0
> **Updated:** 2026-02-10
> **Status:** Active
> **Purpose:** Step-by-step guide to integrate the PowerShell runner with pnpm PnP into any project

---

## Prerequisites

| Requirement | Minimum Version | Auto-Install |
|-------------|-----------------|--------------|
| Windows | 10/11 or Server 2019+ | N/A |
| PowerShell | 5.1 or 7+ | N/A |
| winget | Latest | N/A |
| Go | 1.21+ | ✅ Yes |
| Node.js | 18+ LTS | ✅ Yes |
| pnpm | 8+ | ✅ Yes |

---

## Step 1: Copy Template Files

Copy from spec templates to your project:

```powershell

# From this spec folder

Copy-Item "02-spec/powershell-integration/templates/run.ps1" "YOUR_PROJECT/run.ps1"
```

Or use the existing `run.ps1` if already in your project root.

---

## Step 2: Create Configuration

Create `powershell.json` in your project root:

### Minimal Config

```json
{
  "projectName": "My Project",
  "backendDir": "backend"
}
```

Uses all defaults including pnpm PnP with `.pnpm-store` in project root.

### Standard Config (Recommended)

```json
{
  "$schema": "./02-spec/powershell-integration/schemas/powershell.schema.json",
  "version": "1.1.0",
  "projectName": "My Project",
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
  "pnpmStorePath": ".pnpm-store",
  "cleanPaths": [
    "node_modules",
    "dist",
    ".vite",
    ".pnp.cjs",
    ".pnp.loader.mjs",
    ".pnp.data.json"
  ],
  "buildCommand": "pnpm run build",
  "installCommand": "pnpm install",
  "runCommand": "go run cmd/server/main.go"
}
```

---

## Step 3: Update .gitignore

**IMPORTANT:** Add these entries to your `.gitignore` to exclude pnpm artifacts:

```gitignore

# pnpm store (local cache - do not commit)

.pnpm-store/

# pnpm PnP files (generated on install)

.pnp.cjs
.pnp.loader.mjs

# Build artifacts

dist/
.vite/
```

---

## Step 4: Verify Folder Structure

Ensure your project matches this structure:

```
your-project/
├── .gitignore              ← Must exclude .pnpm-store, .pnp.cjs
├── run.ps1                 ← PowerShell script
├── powershell.json         ← Configuration
├── package.json            ← Frontend dependencies
├── pnpm-lock.yaml          ← pnpm lockfile (commit this!)
├── .pnp.cjs                ← PnP resolution (generated, gitignored)
├── .pnpm-store/            ← pnpm store (generated, gitignored)
├── src/                    ← React source
├── backend/                ← Go backend
│   ├── cmd/server/main.go
│   ├── config.json         ← (created automatically)
│   ├── config.example.json ← Template config
│   └── frontend/
│       └── dist/           ← Build output copied here
└── dist/                   ← pnpm build output (gitignored)
```

---

## Step 5: First Run

```powershell

# Navigate to project

cd YOUR_PROJECT

# Run with help to verify

.\run.ps1 -Help

# Full build and run

.\run.ps1
```

**Expected Output:**

```
========================================
  My Project - Build & Run Script
========================================

[1/5] Pulling latest changes from git...
  ✓ Git pull complete
  ⏱ 1.2s

[2/5] Checking prerequisites...
  ✓ Go found: go version go1.21.0
  ✓ Node.js found: v20.10.0
  ✓ pnpm found: 8.12.0
  ⏱ 0.3s

[3/5] Installing dependencies (pnpm PnP)...
  Store path: .pnpm-store
  ✓ Dependencies installed
  ⏱ 8.2s

[4/5] Building React frontend...
  ✓ Frontend built successfully
  ⏱ 15.2s

[5/5] Starting Go backend...
========================================
  My Project starting...
  Open: http://localhost:8080
  Press Ctrl+C to stop
========================================
```

---

## Step 6: Configure Firewall (Optional)

For network access, run as Administrator:

```powershell

# Right-click PowerShell → Run as Administrator

.\run.ps1 -OpenFirewall
```

This creates inbound rules for your configured ports.

---

## pnpm PnP Configuration

### Setting Up pnpm Store

The store path can be:

| Type | Configuration | Use Case |
|------|---------------|----------|
| **Relative** | `".pnpm-store"` | Single project, isolated |
| **Shared (Recommended)** | `"D:/dev/.pnpm-store"` | Multiple projects, shared packages |
| **User Home** | `"~/.pnpm-store"` | Global store for all projects |

### Multi-Project Setup

For maximum disk savings across multiple projects:

1. Choose a central store location: `D:/dev/.pnpm-store`
2. Configure all projects to use it:

```json
{
  "pnpmStorePath": "D:/dev/.pnpm-store"
}
```

3. Common packages (React, TypeScript, Vite) are stored once and hard-linked.

### Disk Space Comparison

| Scenario | node_modules | pnpm Store |
|----------|--------------|------------|
| 1 project | ~500MB | ~300MB |
| 5 projects (npm) | ~2.5GB | N/A |
| 5 projects (pnpm shared) | N/A | ~400MB |
| **Savings** | - | **~80%** |

---

## Force Reinstall

The `-Force` flag performs a complete clean:

```powershell
.\run.ps1 -Force
```

**What it removes:**

- `.pnp.cjs` and `.pnp.loader.mjs` (PnP resolution files)
- `node_modules/` directory (if exists)
- `dist/` directory
- `.vite/` cache
- SQLite databases in `dataDir`
- Prunes pnpm store (removes unused packages)

**When to use:**

- After major dependency changes
- When builds fail mysteriously
- Periodically to clean up unused packages

---

## Customization

### Custom Build Command

If your project uses different build tools:

```json
{
  "buildCommand": "pnpm run build:prod",
  "runCommand": "go run cmd/server/main.go"
}
```

### Custom Clean Paths

For projects with additional caches:

```json
{
  "cleanPaths": [
    "node_modules",
    "dist",
    ".vite",
    ".pnp.cjs",
    ".pnp.loader.mjs",
    ".next",
    "backend/data/*.db",
    "tmp/"
  ]
}
```

### Required Modules Check

Ensure critical modules are installed:

```json
{
  "requiredModules": [
    "react",
    "react-dom",
    "lucide-react",
    "@tanstack/react-query"
  ]
}
```

If any are missing after install, triggers a reinstall.

### Monorepo Setup

For monorepo projects:

```json
{
  "projectName": "Monorepo App",
  "frontendDir": "packages/web",
  "backendDir": "packages/api",
  "distDir": "build",
  "targetDir": "packages/api/static",
  "usePnp": true,
  "pnpmStorePath": ".pnpm-store"
}
```

---

## Troubleshooting

### pnpm Not Found After Install

**Problem:** pnpm installed but not in PATH

**Solution:**
```powershell

# Restart PowerShell or manually refresh PATH

$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
```

### PnP Resolution Errors

**Problem:** Module not found errors after switching to PnP

**Solution:**
```powershell

# Force reinstall

.\run.ps1 -Force
```

### pnpm Install Fails

**Problem:** Network issues or corrupt cache

**Solution:**
```powershell

# Clear pnpm cache and store

pnpm store prune
pnpm cache clean

# Then force rebuild

.\run.ps1 -Force
```

### Firewall Rules Not Applied

**Problem:** Not running as Administrator

**Solution:**
```powershell

# Right-click PowerShell → Run as Administrator

# Then run:

.\run.ps1 -OpenFirewall
```

### Build Takes Too Long

**Problem:** Cold build with no cache

**Solution:**
```powershell

# Skip build if no frontend changes

.\run.ps1 -SkipBuild

# Or use incremental builds

.\run.ps1  # (without -Force)
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

## AI Handoff Checklist

When asking an AI to integrate this PowerShell runner:

1. ✅ Share `02-spec/powershell-integration/` spec folder
2. ✅ Provide current project structure
3. ✅ Specify port requirements
4. ✅ List any custom build commands
5. ✅ Indicate pnpm store path preference

**Example Prompt:**

> "Integrate the PowerShell runner from spec `02-spec/powershell-integration/` into this project. The backend is in `backend/` and frontend in root. Use port 8080. Enable pnpm PnP with a shared store at `D:/dev/.pnpm-store`."

---

## Cross-References

- [Overview](./01-index.md) - Architecture and quick start
- [Configuration Schema](./02-configuration-schema.md) - JSON config details
- [Script Reference](./04-script-reference.md) - All CLI flags
