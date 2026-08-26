<#
.SYNOPSIS
    Rewrite prior versioned-repo-name tokens to the current version.

.DESCRIPTION
    Detects the repo's base name and current version from the git
    remote URL, then rewrites prior `{Base}-v{N}` tokens to
    `{Base}-v{Current}` across all tracked text files.

    Default mode: replace the last 2 prior versions.
    -2 / -3 / -5: replace the last N prior versions (closed set).
    -All:         replace every prior version (1..Current-1).
    -DryRun:      report changes; do not write.
    -Verbose:     print every modified file path.
    -Config <p>:  path to JSON config (default: ./fix-repo.config.json) with
                  ignoreDirs and ignorePatterns arrays.

    Full normative spec: spec-authoring/22-fix-repo/01-spec.md

.EXAMPLE
    .\fix-repo.ps1
    .\fix-repo.ps1 -3 -DryRun
    .\fix-repo.ps1 -All -Verbose
#>

[CmdletBinding()]
param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$RawArgs
)

$ErrorActionPreference = 'Stop'

$Script:HereDir = Split-Path -Parent $MyInvocation.MyCommand.Path
. (Join-Path $Script:HereDir 'scripts/fix-repo/RepoIdentity.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/FileScan.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/Rewrite.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/Config.ps1')

$Script:ExitOk              = 0
$Script:ExitNotARepo        = 2
$Script:ExitNoRemote        = 3
$Script:ExitNoVersionSuffix = 4
$Script:ExitBadVersion      = 5
$Script:ExitBadFlag         = 6
$Script:ExitWriteFailed     = 7
$Script:ExitBadConfig       = 8

function Test-IsModeFlag { param([string]$A) return $A -in '-2','-3','-5','-all','-All','-ALL' }

# lint-allow: function-length reason="flat CLI flag-parser dispatch; further split harms readability" max=25
function Resolve-Mode {
    param([string[]]$Args)
    $modes      = @()
    $dryRun     = $false
    $verbose    = $false
    $configPath = $null
    $unknown    = @()
    $i = 0
    while ($i -lt $Args.Count) {
        $a = $Args[$i]
        if (Test-IsModeFlag $a) { $modes += $a; $i++; continue }
        if ($a -in '-DryRun','-dryrun')   { $dryRun  = $true;  $i++; continue }
        if ($a -in '-Verbose','-verbose') { $verbose = $true;  $i++; continue }
        if ($a -in '-Config','-config') {
            if ($i + 1 -ge $Args.Count) { return @{ Error = "-Config requires a path" } }
            $configPath = $Args[$i+1]; $i += 2; continue
        }
        if ($a -like '-Config=*' -or $a -like '-config=*') {
            $configPath = $a.Substring($a.IndexOf('=') + 1); $i++; continue
        }
        $unknown += $a; $i++
    }
    if ($modes.Count -gt 1) { return @{ Error = "multiple mode flags: $($modes -join ' ')" } }
    if ($unknown)           { return @{ Error = "unknown flag(s): $($unknown -join ' ')" } }
    $mode = if ($modes.Count -eq 1) { $modes[0].ToLowerInvariant() } else { '-2' }
    return @{ Mode = $mode; DryRun = $dryRun; Verbose = $verbose; ConfigPath = $configPath }
}

function Get-SpanFromMode {
    param([string]$Mode, [int]$Current)
    switch ($Mode) {
        '-2'   { return 2 }
        '-3'   { return 3 }
        '-5'   { return 5 }
        '-all' { return $Current - 1 }
    }
    return 2
}

function Write-Header {
    param($Identity, [int]$Current, [string]$Mode, [int[]]$Targets)
    Write-Host ("fix-repo  base={0}  current=v{1}  mode={2}" -f $Identity.Base, $Current, $Mode)
    $list = if ($Targets.Count -gt 0) { ($Targets | ForEach-Object { "v$_" }) -join ', ' } else { '(none)' }
    Write-Host ("targets:  {0}" -f $list)
    Write-Host ("host:     {0}  owner={1}" -f $Identity.Host, $Identity.Owner)
    Write-Host ''
}

function Write-Summary {
    param([int]$Scanned, [int]$Changed, [int]$Replacements, [bool]$DryRun)
    $modeLabel = if ($DryRun) { 'dry-run' } else { 'write' }
    Write-Host ''
    Write-Host ("scanned: {0} files" -f $Scanned)
    Write-Host ("changed: {0} files ({1} replacements)" -f $Changed, $Replacements)
    Write-Host ("mode:    {0}" -f $modeLabel)
}

function Test-IsScannableFile {
    param([string]$FullPath)
    if (Test-IsSkippablePath -FullPath $FullPath)     { return $false }
    if (Test-IsBinaryExtension -Path $FullPath)       { return $false }
    if (Test-HasNullByte -FullPath $FullPath)         { return $false }
    return $true
}

# lint-allow: function-length reason="flat per-file IO pipeline"
function _Process-OneFile {
    param([string]$RepoRoot, [string]$Rel, [string]$Base, [int]$Current, [int[]]$Targets, [bool]$DryRun, [bool]$Verbose)
    $full = Join-Path $RepoRoot $Rel
    if (-not (Test-Path -LiteralPath $full -PathType Leaf)) { return $null }
    if (Test-IsIgnoredPath -RelPath $Rel)                   { return $null }
    if (-not (Test-IsScannableFile -FullPath $full))        { return $null }
    try {
        $reps = Invoke-FileRewrite -FullPath $full -Base $Base -Targets $Targets -Current $Current -DryRun $DryRun
    } catch {
        Write-Host ("fix-repo: ERROR write failed for {0}: {1}" -f $Rel, $_.Exception.Message)
        return [pscustomobject]@{ Reps=0; Failed=$true }
    }
    if ($reps -gt 0 -and $Verbose) { Write-Host ("modified: {0} ({1} replacements)" -f $Rel, $reps) }
    return [pscustomobject]@{ Reps=$reps; Failed=$false }
}

# lint-allow: function-length reason="flat sweep dispatcher"
function Invoke-RewriteSweep {
    param([string]$RepoRoot, [string]$Base, [int]$Current, [int[]]$Targets, [bool]$DryRun, [bool]$Verbose)
    $files = Get-TrackedFiles -RepoRoot $RepoRoot
    $scanned = 0; $changed = 0; $totalReps = 0; $failed = $false
    foreach ($rel in $files) {
        $r = _Process-OneFile -RepoRoot $RepoRoot -Rel $rel -Base $Base -Current $Current -Targets $Targets -DryRun $DryRun -Verbose $Verbose
        if (-not $r) { continue }
        $scanned++
        if ($r.Failed) { $failed = $true; continue }
        if ($r.Reps -gt 0) { $changed++; $totalReps += $r.Reps }
    }
    return [pscustomobject]@{ Scanned=$scanned; Changed=$changed; Reps=$totalReps; Failed=$failed }
}

# lint-allow: function-length reason="flat identity-source fallback chain"
function Resolve-Identity {
    $root = Get-RepoRoot
    if (-not $root) { Write-Host "fix-repo: ERROR not a git repository (E_NOT_A_REPO)"; exit $Script:ExitNotARepo }
    $url = Get-RemoteUrl
    if (-not $url) { Write-Host "fix-repo: ERROR no remote URL found (E_NO_REMOTE)"; exit $Script:ExitNoRemote }
    $parsed = ConvertFrom-RemoteUrl -Url $url
    if (-not $parsed) { Write-Host ("fix-repo: ERROR cannot parse remote URL '{0}'" -f $url); exit $Script:ExitNoRemote }
    $split = Split-RepoVersion -RepoFull $parsed.Repo
    if (-not $split) {
        Write-Host ("fix-repo: ERROR no -vN suffix on repo name '{0}' (E_NO_VERSION_SUFFIX)" -f $parsed.Repo)
        exit $Script:ExitNoVersionSuffix
    }
    if ($split.Version -lt 1) { Write-Host "fix-repo: ERROR version <= 0 (E_BAD_VERSION)"; exit $Script:ExitBadVersion }
    return [pscustomobject]@{ Root=$root; Host=$parsed.Host; Owner=$parsed.Owner; Base=$split.Base; Current=$split.Version }
}

# ── Main ──────────────────────────────────────────────────────────────
$argList = if ($RawArgs) { $RawArgs } else { @() }
$parsed  = Resolve-Mode -Args $argList
if ($parsed.Error) {
    Write-Host ("fix-repo: ERROR {0} (E_BAD_FLAG)" -f $parsed.Error)
    exit $Script:ExitBadFlag
}

$identity = Resolve-Identity
try {
    Import-FixRepoConfig -ConfigPath $parsed.ConfigPath -RepoRoot $identity.Root
} catch {
    Write-Host ("fix-repo: ERROR {0} (E_BAD_CONFIG)" -f $_.Exception.Message)
    exit $Script:ExitBadConfig
}
$span     = Get-SpanFromMode -Mode $parsed.Mode -Current $identity.Current
$targets  = @(Get-TargetVersions -Current $identity.Current -Span $span)

Write-Header -Identity $identity -Current $identity.Current -Mode $parsed.Mode -Targets $targets

if ($targets.Count -eq 0) {
    Write-Summary -Scanned 0 -Changed 0 -Replacements 0 -DryRun $parsed.DryRun
    Write-Host 'fix-repo: nothing to replace'
    exit $Script:ExitOk
}

$result = Invoke-RewriteSweep -RepoRoot $identity.Root -Base $identity.Base -Current $identity.Current `
    -Targets $targets -DryRun $parsed.DryRun -Verbose $parsed.Verbose

Write-Summary -Scanned $result.Scanned -Changed $result.Changed -Replacements $result.Reps -DryRun $parsed.DryRun

if ($result.Failed) { exit $Script:ExitWriteFailed }
exit $Script:ExitOk
