# build-stamp.ps1 — Windows companion to scripts/build-stamp.sh.
#
# Same purpose: print the exact commit SHA, branch, declared version,
# and a fingerprint of the two historically-problematic cmd/ files
# (updaterepo.go + updatedebugwindows.go) so a stale checkout is
# obvious in the first lines of the build log. See scripts/build-stamp.sh
# for the rationale and the v3.92.0 / v3.113.0 / v3.114.0 history that
# motivated this guard.
#
# Output goes to the host so it survives PowerShell's stream redirection
# quirks. Failures are non-fatal unless -Strict is passed.
#
# Usage:
#   pwsh scripts/build-stamp.ps1
#   pwsh scripts/build-stamp.ps1 -Strict
#
# Called from run.ps1 immediately before `go build`.

[CmdletBinding()]
param(
    [switch]$Strict
)

$ErrorActionPreference = 'Continue'

$StampScriptVersion = '1.0.0'
$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot '..')
$ConstantsFile  = Join-Path $RepoRoot 'gitmap\constants\constants.go'
$UpdateRepoFile = Join-Path $RepoRoot 'gitmap\cmd\updaterepo.go'
$UpdateDebugFile = Join-Path $RepoRoot 'gitmap\cmd\updatedebugwindows.go'

function Probe-Git {
    param([string[]]$Args)

    $git = Get-Command git -ErrorAction SilentlyContinue
    if (-not $git) {
        if ($Strict) {
            Write-Error 'build-stamp: git not found in PATH (strict mode)'
            exit 1
        }
        return '(unknown - git not in PATH)'
    }

    $out = & git -C $RepoRoot @Args 2>$null
    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($out)) {
        return '(unknown)'
    }
    return ($out -join "`n").Trim()
}

function Probe-ConstantsVersion {
    if (-not (Test-Path $ConstantsFile)) {
        return '(unknown - constants.go missing)'
    }
    $line = Select-String -Path $ConstantsFile -Pattern '^const Version = ' | Select-Object -First 1
    if (-not $line) { return '(unknown - pattern miss)' }
    if ($line.Line -match '"([^"]+)"') { return $Matches[1] }
    return '(unknown - parse miss)'
}

function Fingerprint-File {
    param([string]$Label, [string]$Path)

    if (-not (Test-Path $Path)) {
        return ('  {0,-22} (missing - {1})' -f $Label, ($Path.Replace($RepoRoot.Path + '\', '')))
    }
    $sha = (Get-FileHash -Path $Path -Algorithm SHA256).Hash.Substring(0, 12).ToLower()
    $lines = (Get-Content -Path $Path).Count
    $rel = $Path.Replace($RepoRoot.Path + '\', '')
    return ('  {0,-22} sha256:{1}  lines:{2}  {3}' -f $Label, $sha, $lines, $rel)
}

function Detect-RedeclRisk {
    if (-not (Test-Path $UpdateRepoFile) -or -not (Test-Path $UpdateDebugFile)) {
        return '  redecl-risk-check       skipped (one or both source files missing)'
    }
    $pattern = '^func (fileExists|fileExistsLoose)\('
    $repoHits  = (Select-String -Path $UpdateRepoFile  -Pattern $pattern).Count
    $debugHits = (Select-String -Path $UpdateDebugFile -Pattern $pattern).Count

    if ($repoHits -gt 0 -and $debugHits -gt 0) {
        $msg = @(
            '  redecl-risk-check       FAIL - fileExists/fileExistsLoose declared in both files'
            '                           (this checkout predates the v3.113.0 fsutil migration)'
            '                           expected fix: git pull origin main'
        ) -join "`n"
        if ($Strict) {
            Write-Error 'build-stamp: redeclaration risk detected in strict mode'
            Write-Host $msg
            exit 1
        }
        return $msg
    }
    return '  redecl-risk-check       ok (no local fileExists* in cmd/ - fsutil migration present)'
}

$commit       = Probe-Git -Args @('rev-parse', 'HEAD')
$short        = Probe-Git -Args @('rev-parse', '--short=10', 'HEAD')
$branch       = Probe-Git -Args @('rev-parse', '--abbrev-ref', 'HEAD')
$describe     = Probe-Git -Args @('describe', '--tags', '--always', '--dirty')
$commitDate   = Probe-Git -Args @('log', '-1', '--format=%cI')
$commitSubject = Probe-Git -Args @('log', '-1', '--format=%s')
$declaredVer  = Probe-ConstantsVersion

Write-Host ''
Write-Host "=== gitmap build-stamp v$StampScriptVersion ===================================="
Write-Host 'Provenance for stale-checkout detection. Compare these against the SHA'
Write-Host "and version you expected to build - if they don't match, run"
Write-Host "'git pull origin main' before debugging the build error."
Write-Host ''
Write-Host 'git'
Write-Host ('  commit                  {0}' -f $commit)
Write-Host ('  short                   {0}' -f $short)
Write-Host ('  branch                  {0}' -f $branch)
Write-Host ('  describe                {0}' -f $describe)
Write-Host ('  commit-date             {0}' -f $commitDate)
Write-Host ('  commit-subject          {0}' -f $commitSubject)
Write-Host ''
Write-Host 'source'
Write-Host ('  declared-version        {0}' -f $declaredVer)
Write-Host (Fingerprint-File -Label 'constants.go'          -Path $ConstantsFile)
Write-Host (Fingerprint-File -Label 'updaterepo.go'         -Path $UpdateRepoFile)
Write-Host (Fingerprint-File -Label 'updatedebugwindows.go' -Path $UpdateDebugFile)
Write-Host ''
Write-Host 'guards'
Write-Host (Detect-RedeclRisk)
Write-Host '====================================================================='
Write-Host ''
