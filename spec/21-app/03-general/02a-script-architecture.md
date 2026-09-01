# Script Architecture & Configuration

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

## Single Entry Point

One script (`run.ps1`) at the repo root handles the full lifecycle.
It reads configuration from a JSON file and exposes behavior through
switch/string parameters.

```powershell
[CmdletBinding(PositionalBinding=$false)]
param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [string]$DeployPath = "",
    [switch]$Update,
    [switch]$R,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
)
```

## Step-Based Execution

Break the pipeline into numbered steps:

| Step | Action | Skippable |
|------|--------|-----------|
| 1/4 | Git pull | `-NoPull` |
| 2/4 | Resolve dependencies | No |
| 3/4 | Build binary | No |
| 4/4 | Deploy to target | `-NoDeploy` |

Each step is a dedicated function with clear responsibility.

## Configuration Pattern

### External JSON Config

Store build/deploy settings in a JSON file alongside the source:

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "toolname.exe",
  "copyData": true
}
```

### Config Loading

```powershell
function Load-Config {
    $configPath = Join-Path $ProjectDir "powershell.json"
    if (Test-Path $configPath) {
        return Get-Content $configPath | ConvertFrom-Json
    }
    # Return sensible defaults if file is missing
    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "toolname.exe"
        copyData    = $true
    }
}
```

### Rules

- CLI flags always override config file values.
- Missing config file is a warning, not an error.
- All paths in config are relative to the project root unless absolute.

## Error Handling

| Pattern | Implementation |
|---------|----------------|
| `$ErrorActionPreference = "Stop"` | Fail fast on uncaught errors |
| Check `$LASTEXITCODE` after external commands | Detect non-PowerShell failures |
| Print error details before `exit 1` | User sees what went wrong |
| Use `try/finally` with `Push-Location/Pop-Location` | Always restore working directory |

## Cross-References

- Generic spec: [04-build-scripts.md](../08-generic-update/04-build-scripts.md) §PowerShell, §Config Loading
- Generic spec: [03-rename-first-deploy.md](../08-generic-update/03-rename-first-deploy.md) §Rollback
