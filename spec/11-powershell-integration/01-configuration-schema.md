# PowerShell Runner Configuration Schema

> **Version:** 2.1.0  
> **Updated:** 2026-02-04  
> **Status:** Active

---

## Overview

The `powershell.json` configuration file defines project-specific paths and settings for the PowerShell build runner. This allows a single generic script to work across multiple projects with pnpm Plug'n'Play support.

---

## Schema Definition

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "PowerShell Runner Configuration",
  "type": "object",
  "required": ["projectName", "backendDir"],
  "properties": {
    "$schema": {
      "type": "string",
      "description": "JSON Schema reference for validation"
    },
    "version": {
      "type": "string",
      "description": "Configuration/script version for tracking changes",
      "examples": ["1.0.0", "1.1.0"]
    },
    "projectName": {
      "type": "string",
      "description": "Display name for the project",
      "examples": ["WP Plugin Publish", "Spec Management", "My App"]
    },
    "rootDir": {
      "type": "string",
      "default": ".",
      "description": "Root directory of the project (relative to script location). Empty string or '.' means current directory."
    },
    "backendDir": {
      "type": "string",
      "description": "Path to Go backend directory",
      "examples": ["backend", "go-backend", "server"]
    },
    "frontendDir": {
      "type": "string",
      "default": ".",
      "description": "Path to React frontend directory (where package.json lives)"
    },
    "distDir": {
      "type": "string",
      "default": "dist",
      "description": "Frontend build output directory (relative to frontendDir)"
    },
    "targetDir": {
      "type": "string",
      "description": "Where to copy built frontend for serving by backend",
      "examples": ["backend/frontend/dist", "server/static"]
    },
    "dataDir": {
      "type": "string",
      "description": "Data directory for databases and storage",
      "examples": ["backend/data", "data"]
    },
    "ports": {
      "type": "array",
      "items": {"type": "integer"},
      "default": [8080],
      "description": "Ports to open in Windows Firewall"
    },
    "prerequisites": {
      "type": "object",
      "properties": {
        "go": {"type": "boolean", "default": true},
        "node": {"type": "boolean", "default": true},
        "pnpm": {"type": "boolean", "default": true}
      },
      "description": "Which prerequisites to check/install"
    },
    "usePnp": {
      "type": "boolean",
      "default": true,
      "description": "Enable pnpm Plug'n'Play mode for disk-efficient package management"
    },
    "pnpmStorePath": {
      "type": "string",
      "default": "E:/.pnpm-store",
      "description": "Path to pnpm store. Default: E:/.pnpm-store (shared drive for disk efficiency across projects)."
    },
    "cleanPaths": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Paths to remove on -Force clean build",
      "examples": [["node_modules", "dist", ".vite", ".pnp.cjs", "backend/data/*.db"]]
    },
    "buildCommand": {
      "type": "string",
      "default": "pnpm run build",
      "description": "Command to build frontend"
    },
    "installCommand": {
      "type": "string",
      "default": "pnpm install",
      "description": "Command to install dependencies"
    },
    "runCommand": {
      "type": "string",
      "default": "go run main.go",
      "description": "Command to start backend"
    },
    "configFile": {
      "type": "string",
      "default": "config.json",
      "description": "Backend config file name"
    },
    "configExampleFile": {
      "type": "string",
      "default": "config.example.json",
      "description": "Template config file to copy if config missing"
    },
    "requiredModules": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Node modules to verify exist after install (triggers reinstall if missing)",
      "examples": [["react", "react-dom", "lucide-react"]]
    },
    "env": {
      "type": "object",
      "additionalProperties": {"type": "string"},
      "description": "Environment variables to set before running",
      "examples": [{"NODE_ENV": "production"}]
    }
  }
}
```

---

## Example Configurations

### Minimal Configuration

```json
{
  "projectName": "My App",
  "backendDir": "backend"
}
```

Uses all defaults including pnpm PnP.

### Full Configuration (pnpm PnP)

```json
{
  "$schema": "./spec/powershell-integration/schemas/powershell.schema.json",
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
    "backend/data/*.db",
    "backend/data/*.db-shm",
    "backend/data/*.db-wal"
  ],
  "buildCommand": "pnpm run build",
  "installCommand": "pnpm install",
  "runCommand": "go run cmd/server/main.go",
  "configFile": "config.json",
  "configExampleFile": "config.example.json",
  "requiredModules": [
    "react",
    "react-dom",
    "lucide-react"
  ],
  "env": {
    "NODE_ENV": "production"
  }
}
```

### Multi-Project Shared Store

For multiple projects sharing a single pnpm store:

```json
{
  "projectName": "Project A",
  "backendDir": "backend",
  "usePnp": true,
  "pnpmStorePath": "E:/.pnpm-store"
}
```

All projects with the same `pnpmStorePath` share cached packages, saving significant disk space.

### Monorepo Configuration

```json
{
  "projectName": "Monorepo App",
  "rootDir": ".",
  "backendDir": "packages/server",
  "frontendDir": "packages/web",
  "distDir": "build",
  "targetDir": "packages/server/public",
  "dataDir": "packages/server/data",
  "ports": [3000, 3001],
  "usePnp": true,
  "pnpmStorePath": "E:/.pnpm-store",
  "buildCommand": "pnpm run build:web",
  "runCommand": "go run cmd/server/main.go"
}
```

---

## Path Resolution

All paths are resolved relative to the script location (`$MyInvocation.MyCommand.Path`).

**Key Rule:** If `rootDir` is empty or ".", the working directory is where `run.ps1` resides.

```
project-root/               ← Working directory
├── run.ps1                 ← Script location (RootDir base)
├── powershell.json         ← Config file
├── package.json            ← Frontend (frontendDir: ".")
├── .pnp.cjs                ← PnP resolution file
├── .pnpm-store/            ← pnpm store (pnpmStorePath)
│   └── v3/                 ← Cached packages
├── dist/                   ← Build output (distDir: "dist")
└── backend/                ← Backend (backendDir: "backend")
    ├── cmd/server/main.go
    ├── config.json
    ├── config.example.json
    ├── frontend/
    │   └── dist/           ← Target (targetDir)
    └── data/               ← Data (dataDir)
        └── *.db
```

---

## pnpm Store Path Options

| Option | Configuration | Description |
|--------|---------------|-------------|
| **Default (Recommended)** | `"E:/.pnpm-store"` | Shared drive for all projects (default) |
| **Relative (Isolated)** | `".pnpm-store"` | Store in project root, isolated per project |
| **User Home** | `"~/.pnpm-store"` | Global store shared across all projects |

### Disk Space Savings

With a shared store, common packages like `react`, `typescript`, and `vite` are stored once and hard-linked to all projects. Typical savings: **50-70%** of `node_modules` size.

---

## Environment Variables

The config can set environment variables before running:

```json
{
  "env": {
    "NODE_ENV": "production",
    "VITE_API_URL": "http://localhost:8080"
  }
}
```

Script sets these before build/run:
```powershell
foreach ($key in $config.env.PSObject.Properties.Name) {
    [System.Environment]::SetEnvironmentVariable($key, $config.env.$key, "Process")
}
```

---

## Validation

The script validates config on load:

```powershell
function Validate-Config($config) {
    if (-not $config.projectName) {
        throw "powershell.json: projectName is required"
    }
    if (-not $config.backendDir) {
        throw "powershell.json: backendDir is required"
    }
    if (-not (Test-Path (Join-Path $RootDir $config.backendDir))) {
        throw "powershell.json: backendDir '$($config.backendDir)' not found"
    }
    if ($config.usePnp -and -not (Test-Command "pnpm")) {
        Write-Host "Installing pnpm..." -ForegroundColor Yellow
        npm install -g pnpm
    }
}
```

---

## Cross-References

- [Overview](./00-overview.md) - Architecture and quick start
- [Script Reference](./02-script-reference.md) - How the script uses this config
- [Integration Guide](./03-integration-guide.md) - Setup instructions
