# Configuration Pattern

## Overview

This document describes a generic pattern for CLI tool configuration:
load a JSON file, merge with CLI flags, and apply sensible defaults.

## Three-Layer Config

```
Defaults (hardcoded) → Config file (JSON) → CLI flags (highest priority)
```

| Layer | Source | Priority |
|-------|--------|----------|
| 1. Defaults | Constants in code | Lowest |
| 2. Config file | `./data/config.json` or `--config <path>` | Medium |
| 3. CLI flags | `--mode ssh`, `--output-path ./out` | Highest |

## Config File

### Location

- Default path: `./data/config.json` (relative to binary).
- Override: `--config <path>` flag.
- Missing file: use defaults silently (no error).

### Schema

Define a minimal, flat JSON structure:

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": "./toolname-output",
  "excludeDirs": [".cache", "node_modules", "vendor", ".venv"],
  "notes": ""
}
```

### Rules

- All field names are camelCase.
- Array fields default to empty `[]`, not `null`.
- String fields default to `""`, not `null`.
- No nested objects unless absolutely necessary.

## Merge Logic

```go
func LoadAndMerge(configPath, flagMode, flagOutput string) Config {
    cfg := loadDefaults()

    if fileExists(configPath) {
        fileCfg := loadFromFile(configPath)
        cfg = merge(cfg, fileCfg)
    }

    if flagMode != "" {
        cfg.Mode = flagMode
    }
    if flagOutput != "" {
        cfg.Output = flagOutput
    }

    return cfg
}
```

### Merge Rules

1. Load hardcoded defaults into a config struct.
2. If config file exists, overlay its values onto the struct.
3. Apply CLI flags on top — flags always win over everything.
4. If config file is missing, proceed with defaults (no error).

## Build/Deploy Config (PowerShell)

The same pattern applies to build scripts with a separate JSON file:

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "toolname.exe",
  "copyData": true
}
```

### Loading in PowerShell

```powershell
function Load-Config {
    $path = Join-Path $ProjectDir "powershell.json"
    if (Test-Path $path) {
        return Get-Content $path | ConvertFrom-Json
    }
    return @{ deployPath = "E:\bin-run"; buildOutput = "./bin" }
}
```

### Override via Flags

```powershell
param([string]$DeployPath = "")

$target = $Config.deployPath
if ($DeployPath.Length -gt 0) {
    $target = $DeployPath
}
```

## Key Principles

| Principle | Detail |
|-----------|--------|
| Never crash on missing config | Use defaults, warn if needed |
| Flags always win | Explicit user intent overrides everything |
| Config paths are relative to binary | Unless they're absolute |
| Default paths live in `constants` | `DefaultConfigPath`, `DefaultOutputDir` |
| Config struct mirrors JSON exactly | No transformation between file and struct |
