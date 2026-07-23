# misspell-local.ps1 — Windows-first wrapper around scripts/misspell-local.sh.
# Mirrors the CI `spell-check` job filters loaded from gitmap/data/config.json.
#
# Usage:
#   .\scripts\misspell-local.ps1                   # diff vs origin/main
#   .\scripts\misspell-local.ps1 -Base HEAD~1
#   .\scripts\misspell-local.ps1 -All
#   .\scripts\misspell-local.ps1 -Staged
#   .\scripts\misspell-local.ps1 -Files a.md,b.go
[CmdletBinding()]
param(
  [string]   $Base    = "",
  [switch]   $All,
  [switch]   $Staged,
  [string[]] $Files   = @()
)
$ErrorActionPreference = "Stop"
$bash = Get-Command bash -ErrorAction SilentlyContinue
if (-not $bash) {
  Write-Error "bash is required (install Git for Windows or WSL); the underlying script is scripts/misspell-local.sh"
  exit 2
}
$args = @("scripts/misspell-local.sh")
if ($All)    { $args += "--all" }
if ($Staged) { $args += "--staged" }
if ($Base)   { $args += @("--base", $Base) }
if ($Files.Count -gt 0) { $args += "--files"; $args += $Files }
& $bash.Source @args
exit $LASTEXITCODE