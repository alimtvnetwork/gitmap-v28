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
    -Strict:      after rewrite + gofmt, run `go test` on every
                  touched Go package; exit 9 on test failure.
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
. (Join-Path $Script:HereDir 'scripts/fix-repo/Repo-Identity.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/File-Scan.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/Rewrite-Engine.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/Config-Loader.ps1')
. (Join-Path $Script:HereDir 'scripts/fix-repo/Paired-Literal-Audit.ps1')

$Script:ExitOk              = 0
$Script:ExitNotARepo        = 2
$Script:ExitNoRemote        = 3
$Script:ExitNoVersionSuffix = 4
$Script:ExitBadVersion      = 5
$Script:ExitBadFlag         = 6
$Script:ExitWriteFailed     = 7
$Script:ExitBadConfig       = 8
# ExitTestsFailed mirrors the Go binary's FixRepoExitTestsFailed (see
# gitmap/constants/constants_fixrepo.go). Distinct from 7 (write-failed)
# so CI can tell "rewrite produced semantically broken code" apart from
# "the file system rejected the write".
$Script:ExitTestsFailed     = 9
# ExitPairedLiteral fires when the post-rewrite audit catches a
# {base}-v{Current} token next to a stale sibling literal of the
# previous version (typically in *_test.go fixtures).
# See .lovable/memory/issues/2026-05-02-fixrepo-paired-literal-desync.md.
$Script:ExitPairedLiteral   = 10

function Test-IsModeFlag { param([string]$A) return $A -in '-2','-3','-5','-all','-All','-ALL' }

function Resolve-Mode {
    param([string[]]$Args)
    $modes      = @()
    $dryRun     = $false
    $verbose    = $false
    $strict     = $false
    $configPath = $null
    $unknown    = @()
    $i = 0
    while ($i -lt $Args.Count) {
        $a = $Args[$i]
        if (Test-IsModeFlag $a) { $modes += $a; $i++; continue }
        if ($a -in '-DryRun','-dryrun','--dry-run')   { $dryRun  = $true;  $i++; continue }
        if ($a -in '-Verbose','-verbose','--verbose') { $verbose = $true;  $i++; continue }
        if ($a -in '-Strict','-strict','--strict')    { $strict  = $true;  $i++; continue }
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
    return @{ Mode = $mode; DryRun = $dryRun; Verbose = $verbose; Strict = $strict; ConfigPath = $configPath }
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
    Write-Host ("host:     {0}  owner={1}" -f $Identity.RepoHost, $Identity.Owner)
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
        return [pscustomobject]@{ Reps=0; Failed=$true; FullPath=$full }
    }
    if ($reps -gt 0 -and $Verbose) { Write-Host ("modified: {0} ({1} replacements)" -f $Rel, $reps) }
    return [pscustomobject]@{ Reps=$reps; Failed=$false; FullPath=$full }
}

function Invoke-RewriteSweep {
    param([string]$RepoRoot, [string]$Base, [int]$Current, [int[]]$Targets, [bool]$DryRun, [bool]$Verbose)
    $files = Get-TrackedFiles -RepoRoot $RepoRoot
    $scanned = 0; $changed = 0; $totalReps = 0; $failed = $false
    $goFiles = New-Object System.Collections.Generic.List[string]
    $changedFiles = New-Object System.Collections.Generic.List[string]
    foreach ($rel in $files) {
        $r = _Process-OneFile -RepoRoot $RepoRoot -Rel $rel -Base $Base -Current $Current -Targets $Targets -DryRun $DryRun -Verbose $Verbose
        if (-not $r) { continue }
        $scanned++
        if ($r.Failed) { $failed = $true; continue }
        if ($r.Reps -gt 0) {
            $changed++; $totalReps += $r.Reps
            $changedFiles.Add($r.FullPath) | Out-Null
            if ([System.IO.Path]::GetExtension($rel).ToLowerInvariant() -eq '.go') { $goFiles.Add($r.FullPath) | Out-Null }
        }
    }
    return [pscustomobject]@{ Scanned=$scanned; Changed=$changed; Reps=$totalReps; Failed=$failed; GoFiles=$goFiles; ChangedFiles=$changedFiles }
}

# Post-rewrite gofmt step. Mirrors the Go binary's runFixRepoGofmt
# (gitmap/cmd/fixrepo_gofmt.go) so script users get the same CI-clean
# output. Guarded by Get-Command so non-Go environments degrade
# gracefully. See .lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md.
function Invoke-PostRewriteGofmt {
    param([System.Collections.Generic.List[string]]$GoFiles, [bool]$DryRun)
    if ($DryRun)              { Write-Host 'gofmt:   skipped (dry-run)';      return $true }
    if ($GoFiles.Count -eq 0) { Write-Host 'gofmt:   no .go files modified';  return $true }
    if (-not (Get-Command gofmt -ErrorAction SilentlyContinue)) {
        [Console]::Error.WriteLine('fix-repo: WARN  gofmt not found on PATH; skipping post-rewrite formatting')
        return $true
    }
    $out = & gofmt -w @($GoFiles) 2>&1
    if ($LASTEXITCODE -ne 0) {
        [Console]::Error.WriteLine(("fix-repo: ERROR gofmt failed: exit {0}`n{1}" -f $LASTEXITCODE, ($out -join "`n")))
        return $false
    }
    Write-Host ("gofmt:   {0} .go file(s) reformatted" -f $GoFiles.Count)
    return $true
}

# Post-rewrite strict step. Mirrors the Go binary's runFixRepoStrict
# (gitmap/cmd/fixrepo_strict.go): when -Strict is supplied AND the
# rewrite touched at least one .go file AND `go` is on PATH, run
# `go test` against the unique set of touched packages. Returns $true
# on success / safe-skip; $false ONLY when go test itself fails.
# Caller maps $false to ExitTestsFailed (9). Off by default so non-Go
# repos and machines without a Go toolchain stay unaffected.
function Invoke-PostRewriteStrict {
    param([System.Collections.Generic.List[string]]$GoFiles, [bool]$DryRun, [bool]$Strict, [string]$RepoRoot)
    if (-not $Strict)         { return $true }
    if ($DryRun)              { Write-Host 'strict:  skipped (dry-run)';                    return $true }
    if ($GoFiles.Count -eq 0) { Write-Host 'strict:  no .go files modified; skipping go test'; return $true }
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        [Console]::Error.WriteLine('fix-repo: WARN  go not found on PATH; --strict skipped')
        return $true
    }
    # Derive ./<rel-dir> patterns. Forward slashes so the emitted
    # patterns are identical across Windows + Unix log readers (go
    # test accepts both, but downstream tooling prefers forward).
    $packages = New-Object 'System.Collections.Generic.HashSet[string]'
    foreach ($g in $GoFiles) {
        $rel = [System.IO.Path]::GetRelativePath($RepoRoot, $g) -replace '\\','/'
        if ($rel.StartsWith('../')) { continue }
        $dir = [System.IO.Path]::GetDirectoryName($rel) -replace '\\','/'
        if ([string]::IsNullOrEmpty($dir) -or $dir -eq '.') { [void]$packages.Add('.') }
        else                                                { [void]$packages.Add('./' + $dir) }
    }
    if ($packages.Count -eq 0) {
        Write-Host 'strict:  no Go packages derived from modified files; skipping go test'
        return $true
    }
    $sortedPkgs = @($packages) | Sort-Object
    Write-Host ("strict:  running go test on {0} package(s): {1}" -f $sortedPkgs.Count, ($sortedPkgs -join ' '))
    Push-Location $RepoRoot
    try {
        $out = & go test @($sortedPkgs) 2>&1
        if ($LASTEXITCODE -ne 0) {
            [Console]::Error.WriteLine(("fix-repo: ERROR strict mode: go test failed (E_TESTS_FAILED): exit {0}`n{1}" -f $LASTEXITCODE, ($out -join "`n")))
            return $false
        }
    } finally { Pop-Location }
    Write-Host ("strict:  go test passed ({0} package(s))" -f $sortedPkgs.Count)
    return $true
}

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
    return [pscustomobject]@{ Root=$root; RepoHost=$parsed.RepoHost; Owner=$parsed.Owner; Base=$split.Base; Current=$split.Version }
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

if (-not (Invoke-PostRewriteGofmt -GoFiles $result.GoFiles -DryRun $parsed.DryRun)) { $result.Failed = $true }

# Strict step runs AFTER gofmt so go test sees fully-formatted source.
# Tests-failed is a distinct exit code (9) so CI can branch on
# "rewrite produced semantically broken code" vs. write/IO failures.
if (-not (Invoke-PostRewriteStrict -GoFiles $result.GoFiles -DryRun $parsed.DryRun -Strict $parsed.Strict -RepoRoot $identity.Root)) {
    exit $Script:ExitTestsFailed
}

# Paired-literal audit runs LAST (after gofmt + strict) so the source
# being audited is the final, formatted, test-validated state. The
# audit is cheap (regex over *_test.go only) and gives a precise
# diagnostic for the desync class strict-mode would otherwise surface
# as a generic test failure.
if (-not (Invoke-PairedLiteralAudit -ChangedFiles $result.ChangedFiles -Base $identity.Base -Current $identity.Current -DryRun $parsed.DryRun)) {
    exit $Script:ExitPairedLiteral
}

if ($result.Failed) { exit $Script:ExitWriteFailed }
exit $Script:ExitOk
