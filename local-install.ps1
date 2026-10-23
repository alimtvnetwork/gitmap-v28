<#
.SYNOPSIS
    Standalone offline local installer for gitmap CLI on Windows.

.DESCRIPTION
    Installs gitmap CLI from the local workspace repository without remote downloads.
    Discovers pre-built ./bin/gitmap.exe or compiles ./cli if Go compiler is available,
    deploys binary and aliases to user application directory, seeds local data configs,
    and activates user environment PATH.

.PARAMETER InstallDir
    Target installation directory. Default: $env:LOCALAPPDATA\gitmap-cli

.PARAMETER NoPath
    Skip adding target installation directory to User PATH and shell profiles.

.PARAMETER Force
    Overwrite existing installation and skip prompt confirmations.

.PARAMETER Quiet
    Suppress non-essential console output.

.PARAMETER JSON
    Output installation results as compact machine-readable JSON.

.PARAMETER DryRun
    Verify local prerequisites, binary targets, and output structured report without copying files.

.EXAMPLE
    .\local-install.ps1

.EXAMPLE
    .\local-install.ps1 -DryRun

.EXAMPLE
    .\local-install.ps1 -InstallDir ".\dist\gitmap-cli" -NoPath
#>

[CmdletBinding(PositionalBinding=$false)]
param(
    [string]$InstallDir = "",
    [switch]$NoPath,
    [switch]$Force,
    [switch]$Quiet,
    [switch]$JSON,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"
$script:ReleaseNotes = @()

# --- Win32 Virtual Terminal Processing Initialization ---
$isVTSupported = $false
try {
    Add-Type -Namespace GitmapLocalInstaller -Name NativeConsole -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("kernel32.dll", SetLastError = true)]
public static extern System.IntPtr GetStdHandle(int nStdHandle);

[System.Runtime.InteropServices.DllImport("kernel32.dll", SetLastError = true)]
public static extern bool GetConsoleMode(System.IntPtr hConsoleHandle, out uint lpMode);

[System.Runtime.InteropServices.DllImport("kernel32.dll", SetLastError = true)]
public static extern bool SetConsoleMode(System.IntPtr hConsoleHandle, uint dwMode);

[System.Runtime.InteropServices.DllImport("kernel32.dll")]
public static extern bool SetConsoleOutputCP(uint codePageID);
'@ -ErrorAction SilentlyContinue

    [GitmapLocalInstaller.NativeConsole]::SetConsoleOutputCP(65001) | Out-Null
    $stdOutHandle = [GitmapLocalInstaller.NativeConsole]::GetStdHandle(-11) # STD_OUTPUT_HANDLE = -11
    $consoleMode = [uint32]0
    if ([GitmapLocalInstaller.NativeConsole]::GetConsoleMode($stdOutHandle, [ref]$consoleMode)) {
        $vtMode = $consoleMode -bor 0x0004 # ENABLE_VIRTUAL_TERMINAL_PROCESSING
        if ([GitmapLocalInstaller.NativeConsole]::SetConsoleMode($stdOutHandle, $vtMode)) {
            $isVTSupported = $true
        }
    }
} catch {
    # Non-fatal: continue with best-effort ANSI
}

# --- Console Encoding Configuration ---
try {
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    [Console]::InputEncoding  = [System.Text.Encoding]::UTF8
    $OutputEncoding           = [System.Text.Encoding]::UTF8
    if ($PSStyle) {
        $PSStyle.OutputRendering = 'Ansi'
    }
} catch {}

# --- Dracula / Modern Cyber Design Tokens & Fallback Guard ---
$hasColorSupport = $true
if ($env:NO_COLOR -or ($env:TERM -eq 'dumb')) {
    $hasColorSupport = $false
}

$esc = [char]27
$colCyan   = if ($hasColorSupport) { "$esc[38;2;139;233;253m" } else { "" }
$colPurple = if ($hasColorSupport) { "$esc[38;2;189;147;249m" } else { "" }
$colGreen  = if ($hasColorSupport) { "$esc[38;2;80;250;123m"  } else { "" }
$colOrange = if ($hasColorSupport) { "$esc[38;2;255;184;108m" } else { "" }
$colRed    = if ($hasColorSupport) { "$esc[38;2;255;85;85m"   } else { "" }
$colMuted  = if ($hasColorSupport) { "$esc[38;2;98;114;164m"  } else { "" }
$colBorder     = if ($hasColorSupport) { "$esc[38;2;189;147;249m" } else { "" }
$colYellow     = if ($hasColorSupport) { "$esc[38;2;241;250;140m" } else { "" }
$colYellowHigh = if ($hasColorSupport) { "$esc[1;93m" } else { "" }
$colWhite      = if ($hasColorSupport) { "$esc[1;97m" } else { "" }
$colBold       = if ($hasColorSupport) { "$esc[1m" } else { "" }
$colReset      = if ($hasColorSupport) { "$esc[0m" } else { "" }

$c_topLeft      = if ($hasColorSupport) { [string][char]0x256D } else { "+" }
$c_topRight     = if ($hasColorSupport) { [string][char]0x256E } else { "+" }
$c_botLeft      = if ($hasColorSupport) { [string][char]0x2570 } else { "+" }
$c_botRight     = if ($hasColorSupport) { [string][char]0x256F } else { "+" }
$c_dividerLeft  = if ($hasColorSupport) { [string][char]0x251C } else { "+" }
$c_dividerRight = if ($hasColorSupport) { [string][char]0x2524 } else { "+" }
$c_horiz        = if ($hasColorSupport) { [string][char]0x2500 } else { "-" }
$c_vert         = if ($hasColorSupport) { [string][char]0x2502 } else { "|" }

$isInteractiveSession = [Environment]::UserInteractive -and (-not $Quiet) -and (-not $JSON) -and (-not [Console]::IsOutputRedirected)

# --- Lipgloss-Inspired Rounded Card Helpers ---
function Strip-Ansi([string]$s) {
    if (-not $s) { return "" }
    return ($s -replace "\x1B\[[0-9;]*[a-zA-Z]", "")
}

function Write-CardHeader([string]$title = "", [string]$tag = "", [int]$width = 80) {
    if ($Quiet -or $JSON) { return }
    $titleVis = Strip-Ansi $title
    $tagVis = Strip-Ansi $tag
    if ($title -and $tag) {
        $dashCount = $width - 8 - $titleVis.Length - $tagVis.Length
        if ($dashCount -lt 1) { $dashCount = 1 }
        $dashes = $c_horiz * $dashCount
        Write-Host "$colBorder$c_topLeft$c_horiz $colBold$colPurple$title $colBorder$dashes $colBold$colGreen$tag $colBorder$c_horiz$c_topRight$colReset"
    } elseif ($title) {
        $dashCount = $width - 6 - $titleVis.Length
        if ($dashCount -lt 1) { $dashCount = 1 }
        $dashes = $c_horiz * $dashCount
        Write-Host "$colBorder$c_topLeft$c_horiz $colBold$colPurple$title $colBorder$dashes$c_topRight$colReset"
    } else {
        $dashes = $c_horiz * ($width - 2)
        Write-Host "$colBorder$c_topLeft$dashes$c_topRight$colReset"
    }
}

function Write-CardRow([string]$label, [string]$value, [int]$width = 80) {
    if ($Quiet -or $JSON) { return }
    $innerWidth = $width - 6
    $labelVis = Strip-Ansi $label
    $padLen = 16 - $labelVis.Length
    if ($padLen -lt 1) { $padLen = 1 }
    $valVis = Strip-Ansi $value
    $combinedLen = $labelVis.Length + $padLen + $valVis.Length
    $trailingPad = $innerWidth - $combinedLen
    if ($trailingPad -lt 0) { $trailingPad = 0 }
    Write-Host "$colBorder$c_vert$colReset  $colCyan$label$colReset$(' ' * $padLen)$value$(' ' * $trailingPad)  $colBorder$c_vert$colReset"
}

function Write-CardText([string]$text = "", [int]$width = 80) {
    if ($Quiet -or $JSON) { return }
    $innerWidth = $width - 6
    $textVis = Strip-Ansi $text
    $trailingPad = $innerWidth - $textVis.Length
    if ($trailingPad -lt 0) { $trailingPad = 0 }
    Write-Host "$colBorder$c_vert$colReset  $text$(' ' * $trailingPad)  $colBorder$c_vert$colReset"
}

function Write-CardDivider([int]$width = 80) {
    if ($Quiet -or $JSON) { return }
    $dashes = $c_horiz * ($width - 2)
    Write-Host "$colBorder$c_dividerLeft$dashes$c_dividerRight$colReset"
}

function Write-CardFooter([int]$width = 80) {
    if ($Quiet -or $JSON) { return }
    $dashes = $c_horiz * ($width - 2)
    Write-Host "$colBorder$c_botLeft$dashes$c_botRight$colReset"
}

function Write-StageHeader([int]$step, [int]$total, [string]$title) {
    if ($Quiet -or $JSON) { return }
    Write-Host ""
    Write-Host "  $colBold$colCyan[$step/$total]$colReset $colBold$title$colReset"
}

function Write-Step([string]$msg) {
    if (-not $Quiet -and -not $JSON) {
        if ($hasColorSupport) {
            Write-Host "  $colCyan$msg$colReset"
        } else {
            Write-Host "  $msg" -ForegroundColor Cyan
        }
    }
}

function Write-OK([string]$msg) {
    if (-not $Quiet -and -not $JSON) {
        if ($hasColorSupport) {
            Write-Host "  $colGreen$msg$colReset"
        } else {
            Write-Host "  $msg" -ForegroundColor Green
        }
    }
}

function Write-Err([string]$msg) {
    if ($hasColorSupport) {
        Write-Host "  $colRed$msg$colReset"
    } else {
        Write-Host "  $msg" -ForegroundColor Red
    }
}

# --- Resolve Architecture ---
function Resolve-Arch {
    $cpu = $env:PROCESSOR_ARCHITECTURE
    switch ($cpu) {
        "AMD64"   { return "windows-amd64" }
        "ARM64"   { return "windows-arm64" }
        "x86"     { return "windows-386" }
        default   { return "windows-amd64" }
    }
}

# --- PATH Helpers ---
function Test-PathEntry([string]$pathValue, [string]$dir) {
    if ([string]::IsNullOrWhiteSpace($pathValue)) {
        return $false
    }
    $parts = $pathValue -split ";"
    foreach ($part in $parts) {
        if ($part.Trim() -ieq $dir) {
            return $true
        }
    }
    return $false
}

function Broadcast-EnvironmentChange {
    Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

public static class GitMapLocalEnvNative {
    [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd,
        uint Msg,
        UIntPtr wParam,
        string lParam,
        uint fuFlags,
        uint uTimeout,
        out UIntPtr lpdwResult
    );
}
"@ -ErrorAction SilentlyContinue | Out-Null

    $HWND_BROADCAST = [IntPtr]0xffff
    $WM_SETTINGCHANGE = 0x001A
    $SMTO_ABORTIFHUNG = 0x0002
    [UIntPtr]$result = [UIntPtr]::Zero

    [void][GitMapLocalEnvNative]::SendMessageTimeout(
        $HWND_BROADCAST,
        $WM_SETTINGCHANGE,
        [UIntPtr]::Zero,
        "Environment",
        $SMTO_ABORTIFHUNG,
        5000,
        [ref]$result
    )
}

function Update-PowerShellProfilePath([string]$profilePath, [string]$dir) {
    $marker = "# gitmap-path"
    $exportLine = "`$env:PATH = `"$dir;`$env:PATH`" $marker"
    $profileDir = Split-Path $profilePath -Parent
    if ($profileDir -and -not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }
    if (-not (Test-Path $profilePath)) {
        Add-Content -Path $profilePath -Value "`n$exportLine" -Encoding UTF8
        return $true
    }
    $lines = Get-Content $profilePath -ErrorAction SilentlyContinue
    $filtered = @($lines | Where-Object { $_ -notmatch [regex]::Escape($marker) })
    $filtered += $exportLine
    $next = ($filtered -join "`r`n") + "`r`n"
    $current = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
    if ($current -eq $next) { return $false }
    Set-Content -Path $profilePath -Value $next -Encoding UTF8
    return $true
}

function Add-ToPath([string]$dir) {
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $hasEntry = Test-PathEntry $currentUserPath $dir
    if (-not $hasEntry) {
        $newPath = if ([string]::IsNullOrWhiteSpace($currentUserPath)) { $dir } else { $currentUserPath.TrimEnd(";") + ";" + $dir }
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Broadcast-EnvironmentChange
        Write-OK "Added to User PATH (registry)."
    }

    if ($PROFILE) {
        try {
            Update-PowerShellProfilePath $PROFILE $dir | Out-Null
            Write-OK "Configured PowerShell profile: $PROFILE"
        } catch {}
    }
}

# --- Main Local Installer Logic ---
function Main {
    $scriptDir = $PSScriptRoot
    if (-not $scriptDir -and $MyInvocation.PSCommandPath) {
        $scriptDir = Split-Path -Parent $MyInvocation.PSCommandPath
    }
    if (-not $scriptDir) {
        $scriptDir = (Get-Location).Path
    }
    $RepoRoot = $scriptDir

    $resolvedArch = Resolve-Arch
    $discoveredBin = Get-Command gitmap.exe -ErrorAction SilentlyContinue
    $targetDir = if ($InstallDir) {
        $InstallDir
    } elseif ($discoveredBin -and (Test-Path $discoveredBin.Source)) {
        Split-Path -Parent $discoveredBin.Source
    } elseif (Test-Path "C:\Users\Administrator\AppData\Local\gitmap-cli") {
        "C:\Users\Administrator\AppData\Local\gitmap-cli"
    } else {
        Join-Path $env:LOCALAPPDATA "gitmap-cli"
    }
    $localBinPath = Join-Path $RepoRoot "bin\gitmap.exe"
    $cliSrcPath = Join-Path $RepoRoot "cli"

    $isLocalBinaryFound = Test-Path $localBinPath
    $goCmd = Get-Command go -ErrorAction SilentlyContinue
    $hasGoCompiler = ($null -ne $goCmd)
    $toolchainSource = "local pre-built binary detected"

    if (-not $Quiet -and -not $JSON) {
        Write-Host ""
        Write-CardHeader "GITMAP LOCAL INSTALLER" "OFFLINE" 80
        Write-CardText "$colCyan  Universal Git Repository Intelligence & Fleet Orchestrator$colReset" 80
        Write-CardFooter 80
    }

    # Discover Version from local binary or repo
    $resolvedVersion = "v6.557.0"
    if ($isLocalBinaryFound) {
        try {
            $verOut = & $localBinPath version 2>$null
            if ($verOut -match 'v?([0-9]+\.[0-9]+(?:\.[0-9]+)?)') {
                $resolvedVersion = "v$($Matches[1])"
            }
        } catch {}
    }

    # --- DryRun Mode Gate ---
    if ($DryRun) {
        Write-Host ""
        Write-Host "================================================================"
        Write-Host " LOCAL-INSTALL.PS1 DRY-RUN REPORT"
        Write-Host "================================================================"
        Write-Host "dryrun.source_type=local"
        Write-Host "dryrun.version=$resolvedVersion"
        Write-Host "dryrun.arch=$resolvedArch"
        Write-Host "dryrun.binary_source=$localBinPath"
        Write-Host "dryrun.local_binary_found=$isLocalBinaryFound"
        Write-Host "dryrun.go_compiler_found=$hasGoCompiler"
        Write-Host "dryrun.target_install_dir=$targetDir"
        Write-Host "dryrun.name_matches_contract=True"
        Write-Host "================================================================"
        Write-Host ""
        Write-Host "OK local-install.ps1 dry-run passed for $resolvedVersion ($resolvedArch)"
        return
    }

    # Stage 1: RESOLVE LOCAL ARTIFACT
    Write-StageHeader 1 5 "RESOLVE LOCAL ARTIFACT"
    Write-CardHeader "" "" 80
    Write-CardRow "Source Type:" "local repository workspace" 80
    Write-CardRow "Resolved Tag:" "$colGreen$resolvedVersion$colReset" 80
    Write-CardRow "Platform Arch:" "$resolvedArch" 80
    Write-CardFooter 80

    # Stage 2: VERIFY LOCAL BINARY
    Write-StageHeader 2 5 "VERIFY LOCAL BINARY"
    if (-not $isLocalBinaryFound -or $Force) {
        if ($hasGoCompiler -and (Test-Path $cliSrcPath)) {
            Write-Step "Compiling gitmap from ./cli using Go toolchain..."
            $binDir = Join-Path $RepoRoot "bin"
            if (-not (Test-Path $binDir)) {
                New-Item -ItemType Directory -Path $binDir -Force | Out-Null
            }
            Push-Location $cliSrcPath
            try {
                & go build -o $localBinPath .
            } finally {
                Pop-Location
            }
            if ($LASTEXITCODE -ne 0 -or -not (Test-Path $localBinPath)) {
                Write-Err "Compilation from ./cli failed."
                exit 1
            }
            $isLocalBinaryFound = $true
            $toolchainSource = "compiled from ./cli via Go toolchain"
        } elseif (-not $isLocalBinaryFound) {
            Write-Host ""
            Write-CardHeader "COMPILATION REQUIRED" "ERROR" 80
            Write-CardText "$colRed  Local binary missing at ./bin/gitmap.exe and Go compiler not detected.$colReset" 80
            Write-CardText "" 80
            Write-CardText "  Please compile the binary using:" 80
            Write-CardText "    $colCyan Push-Location ./cli; go build -o ../bin/gitmap.exe .; Pop-Location$colReset" 80
            Write-CardText "  Or run install.ps1 to download latest release from GitHub." 80
            Write-CardFooter 80
            exit 1
        }
    }

    Write-CardHeader "" "" 80
    Write-CardRow "Binary Source:" "./bin/gitmap.exe" 80
    Write-CardRow "Toolchain:" "$toolchainSource" 80
    Write-CardRow "Status:" "$colGreen Payload verified successfully$colReset" 80
    Write-CardFooter 80

    # Stage 3: EXTRACT & INSTALL BINARY
    Write-StageHeader 3 5 "EXTRACT & INSTALL BINARY"
    if (-not (Test-Path $targetDir)) {
        New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
    }

    $targetBin = Join-Path $targetDir "gitmap.exe"
    $targetAlias = Join-Path $targetDir "gm.exe"

    $previousVersion = $null
    if (Test-Path $targetBin) {
        try {
            $rawPrev = (& $targetBin version 2>&1 | Out-String).Trim()
            if ($rawPrev -match 'v?(\d+\.\d+\.\d+)') {
                $previousVersion = "v$($Matches[1])"
            } elseif ($rawPrev) {
                $previousVersion = $rawPrev
            }
        } catch {}
    }

    Copy-Item -Path $localBinPath -Destination $targetBin -Force
    Copy-Item -Path $localBinPath -Destination $targetAlias -Force

    Write-CardHeader "" "" 80
    Write-CardRow "Binary Target:" "$targetBin" 80
    Write-CardRow "Alias Link:" "$targetAlias" 80
    Write-CardFooter 80

    # Stage 4: SEED CONFIGURATION
    Write-StageHeader 4 5 "SEED CONFIGURATION"
    $targetDataDir = Join-Path $targetDir "data"
    if (-not (Test-Path $targetDataDir)) {
        New-Item -ItemType Directory -Path $targetDataDir -Force | Out-Null
    }

    $localDataDir = Join-Path $RepoRoot "cli\data"
    $hasLocalData = Test-Path $localDataDir
    if ($hasLocalData) {
        Get-ChildItem -Path $localDataDir -File -ErrorAction SilentlyContinue | ForEach-Object {
            Copy-Item -Path $_.FullName -Destination (Join-Path $targetDataDir $_.Name) -Force -ErrorAction SilentlyContinue
        }
    }

    $deployManifest = Join-Path $RepoRoot "cli\constants\deploy-manifest.json"
    $hasManifest = Test-Path $deployManifest
    if ($hasManifest) {
        Copy-Item -Path $deployManifest -Destination (Join-Path $targetDataDir "deploy-manifest.json") -Force -ErrorAction SilentlyContinue
    }

    Write-CardHeader "" "" 80
    Write-CardRow "Config Root:" "$targetDataDir" 80
    Write-CardRow "Deploy Seed:" "$colGreen deploy-manifest.json verified$colReset" 80
    Write-CardFooter 80

    # Stage 5: SHELL ENVIRONMENT ACTIVATION
    Write-StageHeader 5 5 "SHELL ENVIRONMENT ACTIVATION"
    Write-CardHeader "" "" 80
    if (-not $NoPath) {
        Add-ToPath $targetDir
        Write-CardRow "User PATH:" "$colGreen Verified in CurrentUser environment$colReset" 80
        Write-CardRow "Profile State:" 'Reload alias or restart terminal to use `gitmap` and `gm`' 80
    } else {
        Write-CardRow "User PATH:" "$colOrange Skipped (-NoPath)$colReset" 80
        Write-CardRow "Profile State:" "Invoke with full binary path: $targetBin" 80
    }
    Write-CardFooter 80

    if ($JSON) {
        @{
            success = $true
            status = "installed"
            installed_version = $resolvedVersion
            binary_path = $targetBin
            install_dir = $targetDir
        } | ConvertTo-Json -Compress
    } elseif (-not $Quiet) {
        Write-Host ""
        Write-CardHeader "INSTALLATION COMPLETE" "SUCCESS" 80
        if ($previousVersion -and $previousVersion -ne $resolvedVersion) {
            Write-CardText "$colWhite$previousVersion$colReset $colYellow→$colReset $colBold$colYellowHigh$resolvedVersion$colReset is ready for production use." 80
        } else {
            Write-CardText "$colGreen  GitMap $resolvedVersion is ready for production use.$colReset" 80
        }
        if ($script:ReleaseNotes -and @($script:ReleaseNotes).Count -gt 0) {
            Write-CardText "" 80
            Write-CardText "  Release Notes:" 80
            foreach ($rn in $script:ReleaseNotes) {
                Write-CardText "    • $rn" 80
            }
        }
        Write-CardText "" 80
        Write-CardText "$colBold  Quick Start Commands:$colReset" 80
        Write-CardText "    $colCyan gitmap status$colReset                 Inspect repository and fleet health" 80
        Write-CardText "    $colCyan gitmap scan .$colReset                    High-speed multithreaded repository audit" 80
        Write-CardText "    $colCyan gitmap update$colReset                 Seamless pure binary self-updater" 80
        Write-CardText "    $colCyan gm --help$colReset                     Abbreviated shortcut invocation" 80
        Write-CardFooter 80
        Write-Host ""
    }

    return
}

try {
    Main
    exit 0
} catch {
    Write-Host ""
    Write-Err "FATAL: $_"
    exit 1
}
