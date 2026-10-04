<#
.SYNOPSIS
    Installs Cursor IDE on Ubuntu Linux and synchronizes Cursor Project Manager.
.DESCRIPTION
    Downloads the official Cursor AppImage for Ubuntu Linux x86_64, creates launcher
    symlinks, registers a desktop entry, and synchronizes GitMap tracked projects into
    projects.json.
.EXAMPLE
    pwsh scripts/install-cursor-ubuntu.ps1
#>
[CmdletBinding()]
param(
    [string]$InstallDir = "",
    [switch]$UserMode
)

$ErrorActionPreference = "Stop"

Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "  GitMap — Ubuntu Cursor IDE Automated Provisioning & Sync (pwsh)" -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan

# 1. Platform & Architecture Verification
$arch = & uname -m 2>$null
if ($arch -ne "x86_64") {
    Write-Error "Cursor Linux binary requires x86_64 architecture (detected: $arch)."
    exit 1
}

$isRoot = ([Security.Principal.WindowsIdentity]::GetCurrent().Name -eq "root") -or ((id -u 2>$null) -eq 0)

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    if ($isRoot -and -not $UserMode) {
        $InstallDir = "/opt/cursor"
        $binLink = "/usr/local/bin/cursor"
        $desktopDir = "/usr/share/applications"
        $iconDir = "/usr/share/icons/hicolor/scalable/apps"
    } else {
        $homeDir = [System.Environment]::GetFolderPath([System.Environment+SpecialFolder]::UserProfile)
        $InstallDir = Join-Path $homeDir ".local/share/cursor"
        $binLink = Join-Path $homeDir ".local/bin/cursor"
        $desktopDir = Join-Path $homeDir ".local/share/applications"
        $iconDir = Join-Path $homeDir ".local/share/icons/hicolor/scalable/apps"
    }
} else {
    $binLink = "/usr/local/bin/cursor"
    $desktopDir = "/usr/share/applications"
    $iconDir = "/usr/share/icons/hicolor/scalable/apps"
}

# 2. Ensure directories exist
$dirsToCreate = @(
    $InstallDir,
    (Split-Path $binLink -Parent),
    $desktopDir,
    $iconDir
)

foreach ($d in $dirsToCreate) {
    if (-not (Test-Path $d)) {
        $null = New-Item -ItemType Directory -Path $d -Force
    }
}

# 3. Download Cursor AppImage
$appImagePath = Join-Path $InstallDir "cursor.appimage"
$downloadUrl = "https://downloader.cursor.sh/linux/appImage/x64"

Write-Host "[*] Downloading Cursor AppImage from official endpoint..." -ForegroundColor Cyan
if (Get-Command curl -ErrorAction SilentlyContinue) {
    & curl -fsSL $downloadUrl -o $appImagePath
} else {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $appImagePath
}

& chmod +x $appImagePath

# 4. Create launcher wrapper script
$wrapperPath = Join-Path $InstallDir "cursor"
$wrapperContent = @"
#!/usr/bin/env bash
exec "$appImagePath" --no-sandbox "`$@"
"@
Set-Content -Path $wrapperPath -Value $wrapperContent -NoNewline
& chmod +x $wrapperPath

if (Test-Path $binLink) {
    Remove-Item $binLink -Force
}
& ln -sf $wrapperPath $binLink
Write-Host "[+] Installed Cursor executable wrapper to: $binLink" -ForegroundColor Green

# 5. Write embedded SVG Icon
$iconPath = Join-Path $iconDir "cursor.svg"
$svgContent = @"
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="cursor-grad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#000000" />
      <stop offset="100%" stop-color="#1e1e1e" />
    </linearGradient>
  </defs>
  <rect width="512" height="512" rx="100" fill="url(#cursor-grad)"/>
  <path d="M120 100 L380 256 L240 280 L320 400 L260 430 L180 310 L120 370 Z" fill="#58a6ff"/>
</svg>
"@
Set-Content -Path $iconPath -Value $svgContent -NoNewline

# 6. Create Desktop Application Launcher
$desktopFilePath = Join-Path $desktopDir "cursor.desktop"
$desktopContent = @"
[Desktop Entry]
Version=1.0
Type=Application
Name=Cursor
GenericName=AI Code Editor
Comment=Cursor AI-first Code Editor
Exec=$binLink %F
Icon=$iconPath
Terminal=false
Categories=Development;IDE;TextEditor;
StartupWMClass=Cursor
MimeType=text/plain;inode/directory;
"@
Set-Content -Path $desktopFilePath -Value $desktopContent -NoNewline
& chmod +x $desktopFilePath
Write-Host "[+] Desktop application launcher registered: $desktopFilePath" -ForegroundColor Green

# 7. Project Synchronization into Cursor Project Manager (projects.json)
$homeDir = [System.Environment]::GetFolderPath([System.Environment+SpecialFolder]::UserProfile)
$cursorConfigDir = Join-Path $homeDir ".config/Cursor/User/globalStorage/alefragnani.project-manager"
if (-not (Test-Path $cursorConfigDir)) {
    $null = New-Item -ItemType Directory -Path $cursorConfigDir -Force
}
$projectsFile = Join-Path $cursorConfigDir "projects.json"

Write-Host "[*] Synchronizing Cursor Project Manager ($projectsFile)..." -ForegroundColor Cyan
if (Get-Command gitmap -ErrorAction SilentlyContinue) {
    Write-Host "[*] Running gitmap cursor sync..." -ForegroundColor Cyan
    & gitmap cursor sync
} else {
    if (-not (Test-Path $projectsFile)) {
        Set-Content -Path $projectsFile -Value "[]" -NoNewline
    }
    Write-Host "[+] Initialized Cursor Project Manager database: $projectsFile" -ForegroundColor Green
}

Write-Host "==================================================================" -ForegroundColor Green
Write-Host "  [SUCCESS] Cursor IDE provisioned and synchronized successfully!" -ForegroundColor Green
Write-Host "  Launch anytime by running: cursor [target-directory]" -ForegroundColor Green
Write-Host "==================================================================" -ForegroundColor Green
