<#
.SYNOPSIS
    Pull latest changes, run the coding-guidelines validator, then run all spec/docs linters.

.DESCRIPTION
    1. Performs `git pull`.
    2. Runs the Go-based coding-guidelines validator against -Path.
    3. Runs every linter under `linter-scripts/check-*.{py,sh}`, accumulating failures.
    Exits non-zero if the validator OR any linter fails. Use -d to skip steps 2-3,
    -SkipLinters to skip step 3, -LintersOnly to run only step 3.

.PARAMETER Path        Directory to scan in step 2 (default: src).
.PARAMETER Json        Output validator results as JSON.
.PARAMETER MaxLines    Max function body lines (default: 15).
.PARAMETER d           Skip validation (git pull only).
.PARAMETER SkipLinters Skip step 3.
.PARAMETER LintersOnly Run only step 3.

.EXAMPLE
    .\linter-scripts\run.ps1
    .\linter-scripts\run.ps1 -LintersOnly
    .\linter-scripts\run.ps1 -Path "cmd" -MaxLines 20
#>

param(
    [string]$Path = "src",
    [switch]$Json,
    [int]$MaxLines = 15,
    [switch]$d,
    [switch]$SkipLinters,
    [switch]$LintersOnly
)

# Do NOT use ErrorActionPreference=Stop here — we want to keep running linters
# even when one fails, so we can report a full summary.

function Write-Header { param([string]$Text) Write-Host "`n═══ $Text ═══" -ForegroundColor Cyan }

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$goFile    = Join-Path $scriptDir "validate-guidelines.go"

# ── Step 1: Git Pull ───────────────────────────────────────────────
if (-not $LintersOnly) {
    Write-Header "Step 1 — git pull"
    try {
        git pull
        if ($LASTEXITCODE -ne 0) {
            Write-Host "⚠️  git pull returned exit code $LASTEXITCODE" -ForegroundColor Yellow
        } else {
            Write-Host "✅ Repository up to date." -ForegroundColor Green
        }
    } catch {
        Write-Host "⚠️  git pull failed: $_" -ForegroundColor Yellow
    }
}

if ($d -and -not $LintersOnly) {
    Write-Host "`n⏭️  Skipping validation (-d flag)." -ForegroundColor Yellow
    exit 0
}

$validatorExit = 0

# ── Step 2: Go Validator ───────────────────────────────────────────
if (-not $LintersOnly) {
    Write-Header "Step 2 — Coding-guidelines validator"

    if (-not (Test-Path $goFile)) {
        Write-Host "❌ Cannot find $goFile" -ForegroundColor Red
        exit 1
    }
    try {
        $goVersion = & go version 2>&1
        Write-Host "Using $goVersion" -ForegroundColor Gray
    } catch {
        Write-Host "❌ Go is not installed or not in PATH (https://go.dev/dl/)." -ForegroundColor Red
        exit 1
    }

    $goArgs = @("run", $goFile, "--path", $Path, "--max-lines", $MaxLines)
    if ($Json) { $goArgs += "--json" }

    Write-Host "Scanning: $Path (max $MaxLines lines/function)`n" -ForegroundColor Gray
    & go @goArgs
    $validatorExit = $LASTEXITCODE
    if ($validatorExit -eq 0) {
        Write-Host "✅ Step 2 passed." -ForegroundColor Green
    } else {
        Write-Host "❌ Step 2 failed (exit=$validatorExit)." -ForegroundColor Red
    }
}

# ── Step 3: Linters ────────────────────────────────────────────────
$passed = @()
$failed = @()

function Invoke-Linter {
    param([string]$Label, [string]$Runner, [string]$Script)
    Write-Host "`n── $Label ──" -ForegroundColor Gray
    & $Runner (Join-Path $scriptDir $Script)
    $rc = $LASTEXITCODE
    if ($rc -eq 0) {
        $script:passed += $Label
    } else {
        $script:failed += "$Label (exit=$rc)"
    }
}

if (-not $SkipLinters) {
    Write-Header "Step 3 — Spec / docs linters"

    Invoke-Linter "tunable-constants"    "python3" "check-tunable-constants.py"
    Invoke-Linter "mws-error-codes"      "python3" "check-mws-error-codes.py"
    Invoke-Linter "function-lengths"     "python3" "check-function-lengths.py"
    Invoke-Linter "forbidden-strings"    "python3" "check-forbidden-strings.py"
    Invoke-Linter "placeholder-comments" "python3" "check-placeholder-comments.py"
    Invoke-Linter "memory-mirror-drift"  "python3" "check-memory-mirror-drift.py"
    Invoke-Linter "prompts-loaded"       "python3" "check-prompts-loaded.py"
    Invoke-Linter "readme-canonicals"    "python3" "check-readme-canonicals.py"
    Invoke-Linter "readme-install"       "python3" "check-readme-install-section.py"
    Invoke-Linter "root-readme"          "python3" "check-root-readme.py"
    Invoke-Linter "spec-cross-links"     "python3" "check-spec-cross-links.py"
    Invoke-Linter "spec-folder-refs"     "python3" "check-spec-folder-refs.py"
    Invoke-Linter "axios-version"        "bash"    "check-axios-version.sh"
    Invoke-Linter "forbidden-spec-paths" "bash"    "check-forbidden-spec-paths.sh"
    Invoke-Linter "runner-dispatch"      "bash"    "check-runner-dispatch-antipatterns.sh"
}

# ── Summary ────────────────────────────────────────────────────────
Write-Header "Summary"
Write-Host "Step 2 (validator): exit=$validatorExit"
if (-not $SkipLinters) {
    Write-Host ("Step 3 (linters): {0} passed, {1} failed" -f $passed.Count, $failed.Count)
    foreach ($item in $failed) {
        Write-Host "  ❌ $item" -ForegroundColor Red
    }
}

if ($validatorExit -ne 0 -or $failed.Count -gt 0) {
    exit 1
}
Write-Host "`n✅ All checks passed." -ForegroundColor Green
exit 0
