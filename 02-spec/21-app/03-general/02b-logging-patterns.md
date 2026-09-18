# Logging Patterns

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

## Semantic Logging Functions

Use color-coded helper functions for consistent output:

| Function | Color | Prefix | Use Case |
|----------|-------|--------|----------|
| `Write-Step` | Magenta | `[N/M]` | Step headers |
| `Write-Success` | Green | `OK` | Successful operations |
| `Write-Info` | Cyan/Gray | `->` | Informational messages |
| `Write-Warn` | Yellow | `!!` | Non-fatal warnings |
| `Write-Fail` | Red | `XX` | Errors before exit |

## Banner

Display an ASCII banner at script start for visual identity:

```powershell
function Show-Banner {
    Write-Host "  +======================================+"
    Write-Host "  |         toolname builder             |"
    Write-Host "  +======================================+"
}
```
