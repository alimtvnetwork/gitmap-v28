# Smoke test (Windows): verify a freshly-installed gitmap reports the
# expected version.
#
# Modes:
#   source   Build gitmap from the current checkout into a tempdir, then
#            run `<tempdir>\gitmap.exe version` and assert it matches
#            v$EXPECTED. Used by ci.yml on every PR — no release dependency.
#
#   release  Run gitmap/scripts/install.ps1 against a published GitHub
#            release with a pinned -Version, then run the installed binary
#            and assert. Used by release.yml after the release is cut.
#
# Reads $env:EXPECTED (e.g. "4.1.0"). Falls back to constants.Version.
#
# Exits 0 on success, non-zero with diagnostic on failure.

[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [ValidateSet('source', 'release')]
    [string]$Mode = 'source'
)

$ErrorActionPreference = 'Stop'

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot '..\..')
$expected = $env:EXPECTED

function Get-DeployManifestSmokeConfig {
    $defaults = @{
        AppSubdir        = 'gitmap-cli'
        BinaryNameWin    = 'gitmap.exe'
        LegacyAppSubdirs = @('gitmap')
    }

    $manifestPath = Join-Path $repoRoot 'gitmap\constants\deploy-manifest.json'
    if (-not (Test-Path $manifestPath)) {
        return $defaults
    }

    try {
        $manifest = Get-Content $manifestPath -Raw | ConvertFrom-Json
        return @{
            AppSubdir        = if ($manifest.appSubdir) { [string]$manifest.appSubdir } else { $defaults.AppSubdir }
            BinaryNameWin    = if ($manifest.binaryName.windows) { [string]$manifest.binaryName.windows } else { $defaults.BinaryNameWin }
            LegacyAppSubdirs = if ($manifest.legacyAppSubdirs) { @($manifest.legacyAppSubdirs | ForEach-Object { [string]$_ }) } else { $defaults.LegacyAppSubdirs }
        }
    }
    catch {
        return $defaults
    }
}

$deployConfig = Get-DeployManifestSmokeConfig
if (-not $expected) {
    $constantsPath = Join-Path $repoRoot 'gitmap\constants\constants.go'
    $line = Select-String -Path $constantsPath -Pattern '^const Version' | Select-Object -First 1
    if (-not $line) {
        Write-Error '::error::Could not determine expected version from constants.go'
        exit 2
    }
    $expected = ($line.Line -split '"')[1]
}
$expected = $expected -replace '^v', ''

$work = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "gitmap-smoke-$(Get-Random)") -Force
try {
    Write-Host "▶ Smoke mode:    $Mode"
    Write-Host "▶ Expected:      v$expected"
    Write-Host "▶ Workdir:       $work"

    $bin = $null
    switch ($Mode) {
        'source' {
            Write-Host "▶ Building gitmap from source into $work"
            Push-Location (Join-Path $repoRoot 'gitmap')
            try {
                $binPath = Join-Path $work 'gitmap.exe'
                & go build -o $binPath .
                if ($LASTEXITCODE -ne 0) {
                    Write-Error "::error::go build failed (exit $LASTEXITCODE)"
                    exit 3
                }
                $bin = $binPath
            } finally {
                Pop-Location
            }
        }
        'release' {
            $dest = Join-Path $work 'install'
            New-Item -ItemType Directory -Path $dest -Force | Out-Null
            Write-Host "▶ Running install.ps1 -Version v$expected -NoDiscovery"
            $installer = Join-Path $repoRoot 'gitmap\scripts\install.ps1'
            & $installer -Version "v$expected" -InstallDir $dest -NoPath -NoDiscovery
            if ($LASTEXITCODE -ne 0) {
                Write-Error "::error::install.ps1 failed (exit $LASTEXITCODE)"
                exit 3
            }
            # Resolve from the deploy manifest so the smoke test tracks the
            # installer's canonical layout instead of stale hardcoded paths.
            $candidates = @(
                (Join-Path (Join-Path $dest $deployConfig.AppSubdir) $deployConfig.BinaryNameWin),
                (Join-Path $dest $deployConfig.BinaryNameWin)
            )
            foreach ($legacyAppSubdir in $deployConfig.LegacyAppSubdirs) {
                $candidates += (Join-Path (Join-Path $dest $legacyAppSubdir) $deployConfig.BinaryNameWin)
            }
            $bin = $candidates | Where-Object { Test-Path $_ -PathType Leaf } | Select-Object -First 1
            if (-not $bin) {
                Write-Host "▶ Searching for $($deployConfig.BinaryNameWin) under $dest"
                $found = Get-ChildItem -Path $dest -Recurse -Filter $deployConfig.BinaryNameWin -File -ErrorAction SilentlyContinue | Select-Object -First 1
                if ($found) { $bin = $found.FullName }
            }
            if (-not $bin) {
                Write-Error "::error::Could not locate installed gitmap.exe under $dest"
                Get-ChildItem -Path $dest -Recurse -Depth 4 -ErrorAction SilentlyContinue | ForEach-Object { Write-Host $_.FullName }
                exit 3
            }
            $shim = Join-Path (Split-Path $bin -Parent) 'gitmap.ps1'
            if (-not (Test-Path $shim -PathType Leaf)) {
                Write-Error "::error::Installed release is missing gitmap.ps1 beside gitmap.exe"
                exit 3
            }
            Write-Host "▶ Located binary: $bin"
        }
    }

    if (-not (Test-Path $bin)) {
        Write-Error "::error::Binary not found at $bin"
        exit 3
    }

    $actual = (& $bin version 2>&1 | ForEach-Object { $_.ToString() } | Where-Object { $_ -match '^gitmap v[0-9]' } | Select-Object -First 1)
    if ($actual) { $actual = $actual.Trim() }
    Write-Host "▶ Actual output: $actual"

    $expectedLine = "gitmap v$expected"
    if ($actual -ne $expectedLine) {
        Write-Error "::error::Version mismatch`n  expected: $expectedLine`n  actual:   $actual"
        exit 4
    }

    Write-Host "✅ Installer smoke test passed: $actual"
} finally {
    Remove-Item -Recurse -Force $work -ErrorAction SilentlyContinue
}
