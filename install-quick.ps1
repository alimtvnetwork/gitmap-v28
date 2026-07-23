<#
.SYNOPSIS
    Short interactive installer for gitmap on Windows.

.DESCRIPTION
    Prompts the user for an install drive/folder (with a sensible default),
    then delegates to the canonical gitmap/scripts/install.ps1 with that path.

    Versioned repo discovery: if the source repo URL ends with -v<N>, this
    script probes for higher-numbered sibling repos (-v<N+1>, -v<N+2>, ...)
    and delegates to the latest available one. See:
      spec/01-app/95-installer-script-find-latest-repo.md

    Run via one-liner:
      irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.ps1 | iex

    Or locally:
      ./install-quick.ps1
      ./install-quick.ps1 -InstallDir "E:\Tools\gitmap"
      ./install-quick.ps1 -NoDiscovery
      ./install-quick.ps1 -DiscoveryWindow 20
#>

param(
    [string]$InstallDir       = "",
    [string]$Version          = "",
    [switch]$NoDiscovery,
    [switch]$Interactive,
    [string]$LogFile          = "",
    # Legacy fail-fast knob (retained for back-compat). The canonical knob
    # per spec/07-generic-release/09 §6 is -DiscoveryWindow (default 20,
    # capped at 20 anonymous / 50 with $env:GITHUB_TOKEN).
    [int]$ProbeCeiling        = 30,
    [int]$DiscoveryWindow     = 20
)

$ErrorActionPreference = "Stop"
$ProgressPreference    = "SilentlyContinue"

$Repo          = "alimtvnetwork/gitmap-v27"
$InstallerUrl  = "https://raw.githubusercontent.com/$Repo/main/gitmap/scripts/install.ps1"
$DefaultDir    = "D:\gitmap"

# ---------------------------------------------------------------------------
# Logging + Invoke-Safe helper
# ---------------------------------------------------------------------------
# Why this exists: install-quick.ps1 is run via `irm | iex`, which means a
# raw exception aborts the pipeline with no breadcrumb. We wrap every IO /
# network / filesystem step in Invoke-Safe so we capture the failure, keep
# going where it's safe, and print a final summary pointing to a log file.

if ([string]::IsNullOrWhiteSpace($LogFile)) {
    $stamp   = (Get-Date).ToString("yyyyMMdd-HHmmss")
    $LogFile = Join-Path $env:TEMP "gitmap-install-quick-$stamp.log"
}
$script:InstallErrors = New-Object System.Collections.Generic.List[string]

function Write-Log([string]$message, [string]$level = "INFO") {
    $line = "[{0}] [{1}] {2}" -f (Get-Date -Format "yyyy-MM-dd HH:mm:ss"), $level, $message
    try { Add-Content -Path $LogFile -Value $line -Encoding UTF8 -ErrorAction SilentlyContinue } catch {}
}

function Invoke-Safe {
    param(
        [Parameter(Mandatory)][string]$Step,
        [Parameter(Mandatory)][scriptblock]$Action,
        [switch]$Fatal
    )
    Write-Log "BEGIN: $Step"
    try {
        $result = & $Action
        Write-Log "OK:    $Step"
        return $result
    } catch {
        $msg = "FAIL:  $Step :: $($_.Exception.Message)"
        Write-Log $msg "ERROR"
        Write-Log ($_.ScriptStackTrace) "ERROR"
        $script:InstallErrors.Add("$Step -> $($_.Exception.Message)")
        Write-Host "  [ERROR] $Step : $($_.Exception.Message)" -ForegroundColor Red
        if ($Fatal) { throw }
        return $null
    }
}

Write-Log "install-quick.ps1 started (Repo=$Repo, Interactive=$Interactive)"

# ---------------------------------------------------------------------------
# Versioned repo discovery (spec/01-app/95-installer-script-find-latest-repo.md)
# ---------------------------------------------------------------------------

function Split-RepoSuffix([string]$repo) {
    # Returns @{ Owner=...; Stem=...; N=<int> } or $null if no -v<N> suffix.
    if ($repo -match '^([^/]+)/(.+)-v(\d+)$') {
        return @{
            Owner = $Matches[1]
            Stem  = $Matches[2]
            N     = [int]$Matches[3]
        }
    }
    return $null
}

function Test-RepoExists([string]$url) {
    try {
        $resp = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 5 `
            -UseBasicParsing -ErrorAction Stop
        return ($resp.StatusCode -eq 200)
    } catch {
        return $false
    }
}

# Per spec/07-generic-release/09-generic-install-script-behavior.md §4.1:
# Probe -v<N+1>..-v<N+window> CONCURRENTLY (max 20, or 50 if
# $env:GITHUB_TOKEN is set). Pick max(M) where HEAD returned 200.
# Gaps are tolerated (no fail-fast on first MISS).
function Resolve-EffectiveRepo([string]$repo, [int]$window) {
    $parts = Split-RepoSuffix $repo
    if ($null -eq $parts) {
        Write-Host "  [discovery] no -v<N> suffix on '$repo'; installing baseline as-is"
        return $repo
    }

    $owner    = $parts.Owner
    $stem     = $parts.Stem
    $baseline = $parts.N

    # Concurrency cap: 20 anonymous, 50 if GITHUB_TOKEN supplied.
    $maxConcurrency = 20
    if (-not [string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)) {
        $maxConcurrency = 50
    }
    if ($window -gt $maxConcurrency) { $window = $maxConcurrency }

    Write-Host "  [discovery] baseline: $owner/$stem-v$baseline"
    Write-Host "  [discovery] window: $window (parallel HEAD, max-hit-wins, gap-tolerant)"

    # Build the candidate list.
    $candidates = @()
    for ($m = $baseline + 1; $m -le ($baseline + $window); $m++) {
        $candidates += [pscustomobject]@{
            M   = $m
            Url = "https://github.com/$owner/$stem-v$m"
        }
    }

    # Use a runspace pool so this works on Windows PowerShell 5.1 (no
    # ForEach-Object -Parallel) AND on PowerShell 7+. Concurrency is
    # bounded by $window which we already capped at 20 (or 50 with token).
    $pool = [runspacefactory]::CreateRunspacePool(1, $window)
    $pool.Open()

    $jobs = @()
    foreach ($c in $candidates) {
        $ps = [powershell]::Create().AddScript({
            param($url, $m)
            try {
                $resp = Invoke-WebRequest -Uri $url -Method Head `
                    -TimeoutSec 5 -UseBasicParsing -ErrorAction Stop
                if ($resp.StatusCode -eq 200) {
                    [pscustomobject]@{ M = $m; Url = $url; Hit = $true }
                } else {
                    [pscustomobject]@{ M = $m; Url = $url; Hit = $false }
                }
            } catch {
                [pscustomobject]@{ M = $m; Url = $url; Hit = $false }
            }
        }).AddArgument($c.Url).AddArgument($c.M)
        $ps.RunspacePool = $pool
        $jobs += [pscustomobject]@{ PS = $ps; Handle = $ps.BeginInvoke(); M = $c.M; Url = $c.Url }
    }

    # Collect all results — wait for everything (no early break, gaps allowed).
    $effective = $baseline
    foreach ($j in $jobs) {
        try {
            $result = $j.PS.EndInvoke($j.Handle)
            $r = $result | Select-Object -First 1
            if ($r -and $r.Hit) {
                Write-Host "  [discovery] HEAD $($r.Url) ... HIT"
                if ($r.M -gt $effective) { $effective = $r.M }
            } else {
                Write-Host "  [discovery] HEAD $($j.Url) ... MISS"
            }
        } catch {
            Write-Host "  [discovery] HEAD $($j.Url) ... MISS (error: $_)"
        } finally {
            $j.PS.Dispose()
        }
    }

    $pool.Close()
    $pool.Dispose()

    if ($effective -eq $baseline) {
        Write-Host "  [discovery] no higher version found; using baseline -v$baseline"
        return $repo
    }

    Write-Host "  [discovery] effective: $owner/$stem-v$effective (was -v$baseline)"
    return "$owner/$stem-v$effective"
}

function Invoke-DelegatedInstaller([string]$effectiveRepo, [string]$installDir, [string]$version, [int]$ceiling, [int]$window) {
    $delegatedUrl = "https://raw.githubusercontent.com/$effectiveRepo/main/install-quick.ps1"
    Write-Host "  [discovery] delegating to $delegatedUrl"

    $env:INSTALLER_DELEGATED = "1"
    try {
        $script = (Invoke-WebRequest -Uri $delegatedUrl -UseBasicParsing -TimeoutSec 15).Content
    } catch {
        Write-Host "  [discovery] [WARN] could not fetch delegated installer: $_" -ForegroundColor Yellow
        Write-Host "  [discovery] falling back to baseline installer" -ForegroundColor Yellow
        Remove-Item Env:INSTALLER_DELEGATED -ErrorAction SilentlyContinue
        return $false
    }

    $block = [ScriptBlock]::Create($script)

    $passArgs = @{
        ProbeCeiling    = $ceiling
        DiscoveryWindow = $window
    }
    if (-not [string]::IsNullOrWhiteSpace($installDir)) { $passArgs.InstallDir = $installDir }
    if (-not [string]::IsNullOrWhiteSpace($version))    { $passArgs.Version    = $version }

    & $block @passArgs
    return $true
}

# ---------------------------------------------------------------------------
# Discovery: only run when not already delegated and not opted out.
# ---------------------------------------------------------------------------

$alreadyDelegated = ($env:INSTALLER_DELEGATED -eq "1")

if ($alreadyDelegated) {
    Write-Host "  [discovery] INSTALLER_DELEGATED=1; skipping discovery (loop guard)"
} elseif ($NoDiscovery) {
    Write-Host "  [discovery] -NoDiscovery set; skipping probe"
} elseif (-not [string]::IsNullOrWhiteSpace($Version)) {
    # Strict-tag contract (spec/07-generic-release/09-generic-install-script-behavior.md §3):
    # An explicit -Version pins the install to that exact release.
    # MUST NOT probe -v<N+i> sibling repos. MUST NOT call releases/latest.
    # MUST NOT fall back to main on failure. The canonical installer
    # downstream enforces the same contract on the asset download path.
    Write-Host "  [strict] -Version $Version pinned; skipping repo probe (no fallback)"
} else {
    $effective = Resolve-EffectiveRepo $Repo $DiscoveryWindow
    if ($effective -ne $Repo) {
        $delegated = Invoke-DelegatedInstaller $effective $InstallDir $Version $ProbeCeiling $DiscoveryWindow
        if ($delegated) { return }
        # If delegation failed we fall through and install baseline.
    }
}

# ---------------------------------------------------------------------------
# Baseline install flow (unchanged behavior).
# ---------------------------------------------------------------------------

function Read-InstallDir([string]$default) {
    Write-Host ""
    Write-Host "  gitmap quick installer" -ForegroundColor Cyan
    Write-Host "  ---------------------" -ForegroundColor DarkGray
    Write-Host "  Choose install folder. Press Enter to accept the default." -ForegroundColor Gray
    Write-Host "  Default: $default" -ForegroundColor DarkGray

    $answer = Read-Host "  Install path"
    if ([string]::IsNullOrWhiteSpace($answer)) { return $default }
    return $answer.Trim('"').Trim()
}

function Save-DeployPath([string]$dir) {
    # Persist the chosen install path so `gitmap install scripts` and
    # `run.ps1` pick the same drive/folder automatically.
    try {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
        }
        $cfgPath = Join-Path $dir "powershell.json"
        $cfg = [ordered]@{
            deployPath  = $dir
            buildOutput = "./bin"
            binaryName  = "gitmap.exe"
            goSource    = "./gitmap"
            copyData    = $true
        }
        ($cfg | ConvertTo-Json) | Set-Content -Path $cfgPath -Encoding UTF8
        Write-Host "  Saved deployPath -> $cfgPath" -ForegroundColor DarkGray
    } catch {
        Write-Host "  [WARN] Could not save powershell.json: $_" -ForegroundColor Yellow
    }
}

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    if ($Interactive) {
        $InstallDir = Read-InstallDir $DefaultDir
    } else {
        $InstallDir = $DefaultDir
        Write-Host "  [info] Using default install dir: $InstallDir (pass -Interactive to choose)" -ForegroundColor DarkGray
    }
}

Write-Host ""
Write-Host "  Installing gitmap to: $InstallDir" -ForegroundColor Green
Write-Host "  Log file: $LogFile" -ForegroundColor DarkGray
Write-Host ""

Invoke-Safe -Step "Save deploy path ($InstallDir)" -Action { Save-DeployPath $InstallDir }

$script = Invoke-Safe -Step "Download canonical installer ($InstallerUrl)" -Fatal -Action {
    (Invoke-WebRequest -Uri $InstallerUrl -UseBasicParsing).Content
}
$block  = [ScriptBlock]::Create($script)

Invoke-Safe -Step "Run canonical install.ps1" -Action {
    if ($Version -ne "") {
        & $block -InstallDir $InstallDir -Version $Version
    } else {
        & $block -InstallDir $InstallDir
    }
}

# ---------------------------------------------------------------------------
# Final summary
# ---------------------------------------------------------------------------
Write-Host ""
if ($script:InstallErrors.Count -eq 0) {
    Write-Host "  [OK] gitmap install-quick completed with no errors." -ForegroundColor Green
    Write-Log  "install-quick.ps1 finished OK"
} else {
    Write-Host "  [SUMMARY] $($script:InstallErrors.Count) error(s) occurred during install:" -ForegroundColor Yellow
    foreach ($e in $script:InstallErrors) {
        Write-Host "    - $e" -ForegroundColor Red
    }
    Write-Host ""
    Write-Host "  Full log written to: $LogFile" -ForegroundColor Yellow
    Write-Log  "install-quick.ps1 finished with $($script:InstallErrors.Count) error(s)"
}
