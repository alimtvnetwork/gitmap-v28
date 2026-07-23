<#
.SYNOPSIS
  One-shot repo init: ensure repo is public, then rewrite stale version
  tokens via fix-repo. Both steps run; combined status is reported.

.DESCRIPTION
  Order (per spec/03-general/11-init-pipeline.md):
    1) visibility-change.ps1 -Visible pub -Yes  (no-op if already public)
    2) fix-repo.ps1 -All

  Failure policy: best-effort. Both steps run regardless of the first
  step's exit code; init.ps1 exits 0 only when both succeeded, else
  prints a combined report and exits with the first non-zero step code.

  -Yes is forwarded to visibility-change so the private->public
  confirmation never blocks. Pass -DryRun to preview both steps.

.EXAMPLE
  .\init.ps1
  .\init.ps1 -DryRun
#>

[CmdletBinding()]
param(
    [switch]$DryRun,
    [Alias('h')][switch]$Help
)

$ErrorActionPreference = 'Continue'

$Script:HereDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Show-InitHelp {
    @"
init.ps1 - run visibility-change (force public, auto-yes) then fix-repo --all.

Usage:
  .\init.ps1            # ensure public, then rewrite stale version tokens
  .\init.ps1 -DryRun    # preview both steps

Behavior:
  - Both steps always run (best-effort), even if the first fails.
  - Exit 0 only if both succeeded; otherwise exits with the first
    non-zero step exit code and prints a combined report.
"@ | Write-Host
}

function Invoke-Step {
    param([string]$Label, [string]$Script, [string[]]$ScriptArgs, [string]$GitmapSub, [string[]]$GitmapArgs)
    Write-Host ""
    $localScript = Join-Path $Script:HereDir $Script
    if (Test-Path -LiteralPath $localScript) {
        Write-Host ("==> [{0}] {1} {2}" -f $Label, $Script, ($ScriptArgs -join ' '))
        & $localScript @ScriptArgs
        return $LASTEXITCODE
    }
    $gitmap = Get-Command gitmap -ErrorAction SilentlyContinue
    if ($null -eq $gitmap) {
        Write-Host ("==> [{0}] SKIP — neither {1} nor 'gitmap' binary found on PATH" -f $Label, $Script) -ForegroundColor Yellow
        return 127
    }
    Write-Host ("==> [{0}] gitmap {1} {2}" -f $Label, $GitmapSub, ($GitmapArgs -join ' '))
    & gitmap $GitmapSub @GitmapArgs
    return $LASTEXITCODE
}

function Write-Summary {
    param([int]$VisRc, [int]$FixRc)
    Write-Host ""
    Write-Host "==> init summary"
    Write-Host ("    visibility-change : exit {0}" -f $VisRc)
    Write-Host ("    fix-repo          : exit {0}" -f $FixRc)
}

if ($Help) { Show-InitHelp; exit 0 }

$visScriptArgs = @('-Visible','pub','-Yes')
$visGitmapArgs = @('--yes')
if ($DryRun) { $visScriptArgs += '-DryRun'; $visGitmapArgs += '--dry-run' }
$visRc = Invoke-Step -Label 'visibility' -Script 'visibility-change.ps1' `
    -ScriptArgs $visScriptArgs -GitmapSub 'make-public' -GitmapArgs $visGitmapArgs

$fixScriptArgs = @('-All')
$fixGitmapArgs = @('--all')
if ($DryRun) { $fixScriptArgs += '-DryRun'; $fixGitmapArgs += '--dry-run' }
$fixRc = Invoke-Step -Label 'fix-repo' -Script 'fix-repo.ps1' `
    -ScriptArgs $fixScriptArgs -GitmapSub 'fix-repo' -GitmapArgs $fixGitmapArgs

Write-Summary -VisRc $visRc -FixRc $fixRc

if ($visRc -ne 0) { exit $visRc }
if ($fixRc -ne 0) { exit $fixRc }
exit 0