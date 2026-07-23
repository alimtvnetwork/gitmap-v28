<#
.SYNOPSIS
  Toggle (or set) the visibility of the current repo on GitHub / GitLab.

.DESCRIPTION
  Detects provider + owner/repo from `git remote get-url origin`, reads
  current visibility via `gh` / `glab`, then either toggles it or sets
  it to the value passed via -Visible.

  Confirmation is prompted only when going private -> public (skip with -Yes).
  Use -DryRun to preview without API calls.

  Mirrors the Go-native `gitmap make-public` / `gitmap make-private`
  commands so CI and shell users have a script-only fallback.

.EXAMPLE
  .\visibility-change.ps1                          # toggle
  .\visibility-change.ps1 -Visible pub             # force public (with confirm)
  .\visibility-change.ps1 -Visible pri             # force private
  .\visibility-change.ps1 -Visible pub -Yes -DryRun
#>

[CmdletBinding()]
param(
    [string]$Visible,
    [switch]$Yes,
    [switch]$DryRun,
    [Alias('h')][switch]$Help
)

$ErrorActionPreference = 'Stop'

$Script:HereDir = Split-Path -Parent $MyInvocation.MyCommand.Path
. (Join-Path $Script:HereDir 'visibility-change/Provider.ps1')
. (Join-Path $Script:HereDir 'visibility-change/Apply.ps1')

$Script:ExitOk           = 0
$Script:ExitNotARepo     = 2
$Script:ExitNoOrigin     = 3
$Script:ExitBadProvider  = 4
$Script:ExitAuthFailed   = 5
$Script:ExitBadFlag      = 6
$Script:ExitConfirmReq   = 7
$Script:ExitVerifyFailed = 8

function Show-Help {
    @"
visibility-change.ps1 - toggle/set GitHub/GitLab repo visibility.

Usage:
  .\visibility-change.ps1                  # toggle current visibility
  .\visibility-change.ps1 -Visible pub     # force public
  .\visibility-change.ps1 -Visible pri     # force private
  .\visibility-change.ps1 -Yes             # skip private->public prompt
  .\visibility-change.ps1 -DryRun          # preview, no API call

Env:
  VISIBILITY_GITLAB_HOSTS  comma-separated allow-list of self-hosted GitLab hosts
"@ | Write-Host
}

function Resolve-TargetValue {
    param([string]$Raw)
    if (-not $Raw) { return '' }
    $v = $Raw.ToLowerInvariant()
    if ($v -in 'pub','public')  { return 'public' }
    if ($v -in 'pri','private') { return 'private' }
    return $null
}

function Test-IsRepoRoot {
    & git rev-parse --show-toplevel 2>$null | Out-Null
    return $LASTEXITCODE -eq 0
}

function Resolve-RequiredCli {
    param([string]$Provider)
    if ($Provider -eq 'github') { return 'gh' }
    return 'glab'
}

function Write-Err { param([string]$Msg) [Console]::Error.WriteLine($Msg) }

function Get-OriginUrlOrDie {
    if (-not (Test-IsRepoRoot)) { Write-Err 'visibility-change: ERROR not a git repository'; exit $Script:ExitNotARepo }
    $url = Get-OriginUrl
    if (-not $url) { Write-Err 'visibility-change: ERROR no origin remote'; exit $Script:ExitNoOrigin }
    return $url
}

function Resolve-Context {
    $url = Get-OriginUrlOrDie
    $provider = Resolve-Provider -Url $url
    if (-not $provider) { Write-Err ("visibility-change: ERROR unsupported host in '{0}'" -f $url); exit $Script:ExitBadProvider }
    $slug = Resolve-OwnerRepo -Url $url
    if (-not $slug) { Write-Err ("visibility-change: ERROR cannot parse owner/repo from '{0}'" -f $url); exit $Script:ExitBadProvider }
    return [pscustomobject]@{ Url=$url; Provider=$provider; Slug=$slug }
}

function Assert-CliReady {
    param([string]$Cli)
    if (Test-CliAvailable -Name $Cli) { return }
    Write-Err ("visibility-change: ERROR '{0}' not found on PATH (install: https://cli.github.com / https://gitlab.com/gitlab-org/cli)" -f $Cli)
    exit $Script:ExitAuthFailed
}

function Resolve-NextTarget {
    param([string]$Forced, [string]$Current)
    if ($Forced) { return $Forced }
    if ($Current -eq 'public') { return 'private' }
    return 'public'
}

function Invoke-MaybeConfirm {
    param([string]$Current, [string]$Target, [bool]$YesFlag, [string]$Slug, [string]$Provider)
    if ($Target -ne 'public') { return }
    if ($Current -ne 'private') { return }
    if ($YesFlag) { return }
    $ok = Confirm-PublicChange -Slug $Slug -Provider $Provider
    if ($ok) { return }
    Write-Err 'visibility-change: ERROR confirmation required (pass -Yes for non-interactive)'
    exit $Script:ExitConfirmReq
}

# -- Main ----------------------------------------------------------------
if ($Help) { Show-Help; exit $Script:ExitOk }

$forced = Resolve-TargetValue -Raw $Visible
if ($null -eq $forced) {
    Write-Err ("visibility-change: ERROR bad -Visible value '{0}' (use pub|public|pri|private)" -f $Visible)
    exit $Script:ExitBadFlag
}

$ctx = Resolve-Context
Assert-CliReady -Cli (Resolve-RequiredCli -Provider $ctx.Provider)

$current = Get-CurrentVisibility -Provider $ctx.Provider -Slug $ctx.Slug
if (-not $current) {
    Write-Err 'visibility-change: ERROR cannot read current visibility (auth?)'
    exit $Script:ExitAuthFailed
}

$target = Resolve-NextTarget -Forced $forced -Current $current

if ($current -eq $target) {
    Write-Host ("visibility: already {0} ({1})" -f $current, $ctx.Provider)
    exit $Script:ExitOk
}

Invoke-MaybeConfirm -Current $current -Target $target -YesFlag $Yes.IsPresent -Slug $ctx.Slug -Provider $ctx.Provider

if ($DryRun) {
    Write-Host ("[dry-run] visibility: {0} -> {1} ({2})" -f $current, $target, $ctx.Provider)
    exit $Script:ExitOk
}

$rc = Invoke-VisibilityApply -Provider $ctx.Provider -Slug $ctx.Slug -Target $target
if ($rc -ne 0) {
    Write-Err ("visibility-change: ERROR apply failed (exit {0})" -f $rc)
    exit $Script:ExitAuthFailed
}

if (-not (Test-VisibilityMatches -Provider $ctx.Provider -Slug $ctx.Slug -Target $target)) {
    Write-Err 'visibility-change: ERROR verification failed (visibility did not change)'
    exit $Script:ExitVerifyFailed
}

Write-Host ("visibility: {0} -> {1} ({2})" -f $current, $target, $ctx.Provider)
exit $Script:ExitOk
