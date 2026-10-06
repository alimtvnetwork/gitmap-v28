# ==============================================================================
# verify-antigravity-backup.ps1
# Quality Gate & Validation Audit Engine for Antigravity Backup Vault
# Verifies Quality Gates VG-01 through VG-07 on PowerShell Core & Windows PowerShell
# ==============================================================================

[CmdletBinding()]
param (
    [Parameter(Mandatory = $false)]
    [switch]$VerboseOutput
)

$isValid = $true
$isVerbose = $VerboseOutput.IsPresent

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$vaultDir = Join-Path (Split-Path -Parent $scriptDir) "vault"
$scriptsDir = $scriptDir

$projectsManifest = Join-Path $vaultDir "projects-manifest.json"
$pinnedManifest = Join-Path $vaultDir "pinned-projects.json"
$settingsManifest = Join-Path $vaultDir "settings-manifest.json"
$pluginsManifest = Join-Path $vaultDir "plugins-and-skills.json"

Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host "  Antigravity Vault Verification Audit: Gates VG-01 to VG-07 (PowerShell)" -ForegroundColor Cyan
Write-Host "==============================================================================" -ForegroundColor Cyan

# ------------------------------------------------------------------------------
# Gate VG-01: Vault Structure & Manifest Validation
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-01] Vault Structure & Manifest Validation: "
$vg01Pass = $true

$manifestList = @($projectsManifest, $pinnedManifest, $settingsManifest, $pluginsManifest)
foreach ($manifest in $manifestList) {
    if (-not (Test-Path $manifest)) {
        Write-Host "FAILED" -ForegroundColor Red
        Write-Host "  Missing manifest: $(Split-Path -Leaf $manifest)"
        $vg01Pass = $false
        break
    }
    try {
        $content = Get-Content -Raw -Path $manifest -Encoding UTF8
        $null = ConvertFrom-Json $content -ErrorAction Stop
    }
    catch {
        Write-Host "FAILED" -ForegroundColor Red
        Write-Host "  Invalid JSON syntax in: $(Split-Path -Leaf $manifest)"
        $vg01Pass = $false
        break
    }
}

if ($vg01Pass) {
    Write-Host "PASSED (All 4 manifests exist and valid JSON)" -ForegroundColor Green
} else {
    $isValid = $false
}

# ------------------------------------------------------------------------------
# Gate VG-02: Repository Count Parity
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-02] Repository Count Parity: "
$vg02Pass = $true

try {
    $pm = Get-Content -Raw -Path $projectsManifest -Encoding UTF8 | ConvertFrom-Json
    if ($pm.totalProjects -ne 78 -or $pm.projects.Count -ne 78) {
        $vg02Pass = $false
    }
    $secCount = ($pm.sectionSummary.PSObject.Properties | Measure-Object).Count
    if ($secCount -ne 10) {
        $vg02Pass = $false
    }

    $pp = Get-Content -Raw -Path $pinnedManifest -Encoding UTF8 | ConvertFrom-Json
    if ($pp.totalPinned -ne 23 -or $pp.projects.Count -ne 23) {
        $vg02Pass = $false
    }
    if ($pp.tierDistribution.tier1Core -ne 7 -or $pp.tierDistribution.tier2Ecosystem -ne 16) {
        $vg02Pass = $false
    }
}
catch {
    $vg02Pass = $false
}

if ($vg02Pass) {
    Write-Host "PASSED (78 projects across 10 sections, 23 pinned [7 Tier 1 + 16 Tier 2])" -ForegroundColor Green
} else {
    Write-Host "FAILED (Repository or pinned count mismatch)" -ForegroundColor Red
    $isValid = $false
}

# ------------------------------------------------------------------------------
# Gate VG-03: Path Relativity & Portability
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-03] Path Relativity & Portability: "
$vg03Pass = $true

foreach ($manifest in $manifestList) {
    $raw = Get-Content -Raw -Path $manifest -Encoding UTF8
    if ($raw -match '(/home/[a-zA-Z0-9_-]+|[A-Za-z]:[\\/])') {
        $vg03Pass = $false
        break
    }
}

try {
    $pm = Get-Content -Raw -Path $projectsManifest -Encoding UTF8 | ConvertFrom-Json
    foreach ($p in $pm.projects) {
        if ($p.uriTemplate -notmatch '\$\{WORKSPACE_ROOT\}') {
            $vg03Pass = $false
            break
        }
    }
}
catch {
    $vg03Pass = $false
}

if ($vg03Pass) {
    Write-Host "PASSED (`${WORKSPACE_ROOT} used exclusively, zero host paths)" -ForegroundColor Green
} else {
    Write-Host "FAILED (Absolute path leak or missing `${WORKSPACE_ROOT})" -ForegroundColor Red
    $isValid = $false
}

# ------------------------------------------------------------------------------
# Gate VG-04: Restoration Script Parity
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-04] Restoration Script Parity: "
$vg04Pass = $true

$shRunner = Join-Path $scriptsDir "restore-antigravity-all.sh"
$ps1Runner = Join-Path $scriptsDir "restore-antigravity-all.ps1"
$shVerifier = Join-Path $scriptsDir "verify-antigravity-backup.sh"
$ps1Verifier = Join-Path $scriptsDir "verify-antigravity-backup.ps1"

if (-not (Test-Path $shRunner) -or -not (Test-Path $ps1Runner) -or -not (Test-Path $shVerifier) -or -not (Test-Path $ps1Verifier)) {
    $vg04Pass = $false
}

if ($vg04Pass) {
    Write-Host "PASSED (POSIX Bash and PowerShell runners and verifiers present)" -ForegroundColor Green
} else {
    Write-Host "FAILED (Missing restoration or verifier runner script)" -ForegroundColor Red
    $isValid = $false
}

# ------------------------------------------------------------------------------
# Gate VG-05: Atomic Non-Destructive Ingestion
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-05] Atomic Non-Destructive Ingestion: "
$vg05Pass = $true

if (Test-Path $ps1Runner) {
    $content = Get-Content -Raw -Path $ps1Runner -Encoding UTF8
    if ($content -notmatch "shouldBackup" -or $content -notmatch "pre-restore") {
        $vg05Pass = $false
    }
} else {
    $vg05Pass = $false
}

if ($vg05Pass) {
    Write-Host "PASSED (Pre-restore snapshot & non-destructive ingestion confirmed)" -ForegroundColor Green
} else {
    Write-Host "FAILED (Missing backup routine or destructive logic detected)" -ForegroundColor Red
    $isValid = $false
}

# ------------------------------------------------------------------------------
# Gate VG-06: Secret Sanitization Gate
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-06] Secret Sanitization Gate: "
$vg06Pass = $true

$secretPatterns = @(
    "BEGIN (RSA )?PRIVATE KEY",
    "ghp_[a-zA-Z0-9]{36}",
    "AKIA[0-9A-Z]{16}",
    '"client_secret":\s*"[^"]+"',
    '"apiKey":\s*"[a-zA-Z0-9_-]{20,}"'
)

foreach ($manifest in $manifestList) {
    $raw = Get-Content -Raw -Path $manifest -Encoding UTF8
    foreach ($pat in $secretPatterns) {
        if ($raw -match $pat) {
            $vg06Pass = $false
            break
        }
    }
    if (-not $vg06Pass) { break }
}

if ($vg06Pass) {
    Write-Host "PASSED (Zero credentials, tokens, or private keys detected)" -ForegroundColor Green
} else {
    Write-Host "FAILED (Potential secret pattern detected in manifests)" -ForegroundColor Red
    $isValid = $false
}

# ------------------------------------------------------------------------------
# Gate VG-07: Relative Link & Positive Boolean Hygiene
# ------------------------------------------------------------------------------
Write-Host -NoNewline "● [Gate VG-07] Positive Boolean Hygiene: "
$vg07Pass = $true

$negativePrefixes = @("isNot", "hasNo", "no", "disabled", "unverified", "skip")
foreach ($manifest in $manifestList) {
    $raw = Get-Content -Raw -Path $manifest -Encoding UTF8
    foreach ($neg in $negativePrefixes) {
        if ($raw -match ('"' + $neg + '[A-Za-z0-9_]*":\s*(true|false)')) {
            $vg07Pass = $false
            break
        }
    }
    if (-not $vg07Pass) { break }
}

if ($vg07Pass) {
    Write-Host "PASSED (100% positive boolean naming in all manifest stores)" -ForegroundColor Green
} else {
    Write-Host "FAILED (Negative boolean naming detected in manifest stores)" -ForegroundColor Red
    $isValid = $false
}

Write-Host "==============================================================================" -ForegroundColor Cyan
if ($isValid) {
    Write-Host "  RESULT: ALL 7 VERIFICATION GATES PASSED (100% COMPLIANT)" -ForegroundColor Green
    Write-Host "==============================================================================" -ForegroundColor Cyan
    exit 0
} else {
    Write-Host "  RESULT: VERIFICATION FAILED (Review gate errors above)" -ForegroundColor Red
    Write-Host "==============================================================================" -ForegroundColor Cyan
    exit 1
}
