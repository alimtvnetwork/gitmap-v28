# PowerShell Build & Run Script — Generic Template
# Version: 1.3.0
# Generic template for Go backend + React frontend projects with pnpm PnP support
# Configure via powershell.json — see spec/powershell-integration/01-configuration-schema.md
#
# USAGE:
#   Copy this file and powershell.json to your project root.
#   Edit powershell.json with your project-specific paths and settings.
#   Run: .\run.ps1 -h for help.
#
# FLAGS:
#   -h   Show help
#   -b   Build frontend only (don't start backend)
#   -s   Skip frontend build (start backend only)
#   -p   Skip git pull
#   -f   Force clean build (remove caches, deps, databases)
#   -i   Install/update all dependencies (frontend + backend)
#   -r   Rebuild (combines -f + -i for complete clean reinstall)
#   -fw  Add Windows Firewall inbound rules (requires Admin)
#   -u   Upload plugin to WordPress via upload-plugin-v2.ps1
#   -pp  Override plugin path for upload (use with -u)
#   -d   Debug mode for upload (verbose request/response logging)
#   -v   Verbose debug output
#
# PIPELINE:
#   1. Git Pull → 2. Prerequisites → 3. pnpm Install → 4. Build → 5. Copy → 6. Run
#   Upload mode (-u): 1. Git Pull → 2. Prerequisites → Upload Plugin V2
#
# FEATURES:
#   - Auto-install Go, Node.js, pnpm via winget if missing
#   - pnpm PnP mode for disk-efficient package management
#   - Shared pnpm store across multiple projects
#   - Force clean mode: removes node_modules, dist, .vite, PnP artifacts, databases
#   - Windows Firewall rule management
#   - pnpm v10+ compatibility (auto --dangerously-allow-all-builds)
#   - Cross-drive store detection (falls back to isolated linker)
#   - Node v24+ detection (falls back to isolated linker for ESM compat)
#   - WordPress plugin upload integration (optional, via powershell.json)
#
# CONFIGURATION (powershell.json):
#   {
#     "projectName": "My App",
#     "rootDir": ".",
#     "backendDir": "backend",
#     "frontendDir": ".",
#     "distDir": "dist",
#     "targetDir": "backend/frontend/dist",
#     "dataDir": "backend/data",
#     "ports": [8080],
#     "prerequisites": { "go": true, "node": true, "pnpm": true },
#     "usePnp": true,
#     "pnpmStorePath": ".pnpm-store",
#     "cleanPaths": ["node_modules", "dist", ".vite"],
#     "buildCommand": "pnpm run build",
#     "installCommand": "pnpm install",
#     "runCommand": "go run cmd/server/main.go",
#     "configFile": "config.json",
#     "configExampleFile": "config.example.json",
#     "upload": {
#       "scriptPath": "wp-plugins/scripts/upload-plugin-v2.ps1",
#       "defaultPluginPath": "wp-plugins/riseup-asia-uploader",
#       "configPath": "wp-plugins/scripts/wp-plugin-config.json"
#     }
#   }
#
# UPLOAD INTEGRATION (optional):
#   Add the "upload" section to powershell.json to enable -u flag.
#   The upload script must exist at the configured scriptPath.
#   Example:
#     .\run.ps1 -u              # Upload default plugin
#     .\run.ps1 -u -d           # Upload with debug output
#     .\run.ps1 -u -pp "C:\custom-plugin"  # Upload custom path
#
# See spec/powershell-integration/ for full documentation.

param(
    [Alias('b')][switch]$buildonly,
    [Alias('s')][switch]$skipbuild,
    [Alias('p')][switch]$skippull,
    [Alias('f')][switch]$force,
    [Alias('i')][switch]$install,
    [Alias('r')][switch]$rebuild,
    [Alias('fw')][switch]$openfirewall,
    [Alias('u')][switch]$upload,
    [Alias('pp')][string]$pluginpath = "",
    [Alias('d')][switch]$debugmode,
    [Alias('h')][switch]$help,
    [Alias('v')][switch]$verbose
)

# -rebuild is a convenience flag that combines -force and -install
if ($rebuild) {
    $force = $true
    $install = $true
}

$ErrorActionPreference = "Stop"

# ============================================================================
# PATH RESOLUTION: Script location is the working directory
# ============================================================================
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
if ([string]::IsNullOrWhiteSpace($ScriptDir)) {
    $ScriptDir = Get-Location
}

# ============================================================================
# CONFIGURATION LOADING
# ============================================================================
$ConfigPath = Join-Path $ScriptDir "powershell.json"

if (-not (Test-Path $ConfigPath)) {
    Write-Host "ERROR: powershell.json not found at: $ConfigPath" -ForegroundColor Red
    Write-Host "Create a powershell.json configuration file in the script directory." -ForegroundColor Yellow
    Write-Host "See spec/powershell-integration/01-configuration-schema.md for format." -ForegroundColor Yellow
    exit 1
}

try {
    $Config = Get-Content $ConfigPath -Raw | ConvertFrom-Json
} catch {
    Write-Host "ERROR: Failed to parse powershell.json: $_" -ForegroundColor Red
    exit 1
}

# Resolve paths - handles both relative and absolute paths
function Resolve-RelativePath($Path) {
    if ([string]::IsNullOrWhiteSpace($Path) -or $Path -eq ".") {
        return $ScriptDir
    }
    if ($Path -match '^[A-Za-z]:' -or $Path -match '^\\\\') {
        return $Path -replace '/', '\'
    }
    return Join-Path $ScriptDir $Path
}

# Configuration with defaults
$ProjectName = if ($Config.projectName) { $Config.projectName } else { "Project" }
$RootDir = Resolve-RelativePath $Config.rootDir
$BackendDir = Resolve-RelativePath $Config.backendDir
$FrontendDir = Resolve-RelativePath $Config.frontendDir
$DistDir = if ($Config.distDir) { $Config.distDir } else { "dist" }
$TargetDir = if ($Config.targetDir) { Resolve-RelativePath $Config.targetDir } else { $null }
$DataDir = if ($Config.dataDir) { Resolve-RelativePath $Config.dataDir } else { $null }
$Ports = if ($Config.ports) { $Config.ports } else { @(8080) }
$BuildCommand = if ($Config.buildCommand) { $Config.buildCommand } else { "pnpm run build" }
$InstallCommand = if ($Config.installCommand) { $Config.installCommand } else { "pnpm install" }
$RunCommand = if ($Config.runCommand) { $Config.runCommand } else { "go run cmd/server/main.go" }
$CleanPaths = if ($Config.cleanPaths) { $Config.cleanPaths } else { @("node_modules", "dist", ".vite") }
$ConfigFile = if ($Config.configFile) { $Config.configFile } else { "config.json" }
$ConfigExampleFile = if ($Config.configExampleFile) { $Config.configExampleFile } else { "config.example.json" }
$RequiredModules = if ($Config.requiredModules) { $Config.requiredModules } else { @() }

# pnpm configuration
$PnpmStorePath = if ($Config.pnpmStorePath) { Resolve-RelativePath $Config.pnpmStorePath } else { $null }
$UsePnp = if ($null -ne $Config.usePnp) { $Config.usePnp } else { $true }

# Prerequisites
$CheckGo = if ($null -ne $Config.prerequisites -and $null -ne $Config.prerequisites.go) { $Config.prerequisites.go } else { $true }
$CheckNode = if ($null -ne $Config.prerequisites -and $null -ne $Config.prerequisites.node) { $Config.prerequisites.node } else { $true }
$CheckPnpm = if ($null -ne $Config.prerequisites -and $null -ne $Config.prerequisites.pnpm) { $Config.prerequisites.pnpm } else { $true }

# pnpm version-aware install behavior (pnpm v10+ blocks dependency build scripts by default)
$PnpmMajor = 0
$NodeMajor = 0
$EffectiveInstallCommand = $InstallCommand
$DidFrontendInstall = $false
$EffectiveNodeLinker = if ($UsePnp) { "pnp" } else { "isolated" }

$TotalStopwatch = [System.Diagnostics.Stopwatch]::StartNew()

# ============================================================================
# HELP
# ============================================================================
if ($help) {
    Write-Host ""
    Write-Host "$ProjectName - Build & Run Script" -ForegroundColor Cyan
    Write-Host ("=" * ($ProjectName.Length + 22)) -ForegroundColor Cyan
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\run.ps1 [flags]"
    Write-Host ""
    Write-Host "FLAGS:" -ForegroundColor Yellow
    Write-Host "  -h,  -help          Show this help message and exit"
    Write-Host "  -b,  -buildonly     Build frontend only, don't start the backend server"
    Write-Host "  -s,  -skipbuild     Skip frontend build, only run the backend server"
    Write-Host "  -p,  -skippull      Skip git pull step"
    Write-Host "  -f,  -force         Clean build: remove caches, dependencies, databases"
    Write-Host "  -i,  -install       Install/update dependencies (frontend + backend)"
    Write-Host "  -r,  -rebuild       Complete clean reinstall (combines -f + -i)"
    Write-Host "  -fw, -openfirewall  (Admin) Add Windows Firewall inbound rules"
    Write-Host "  -u,  -upload        Upload plugin to WordPress via upload-plugin-v2"
    Write-Host "  -pp, -pluginpath    Override plugin folder path (use with -u)"
    Write-Host "  -d,  -debugmode     Debug mode for upload (verbose HTTP logging)"
    Write-Host "  -v,  -verbose       Show detailed debug output"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\run.ps1              # Full build and run"
    Write-Host "  .\run.ps1 -i           # Install/update all dependencies"
    Write-Host "  .\run.ps1 -r           # Complete clean reinstall and build"
    Write-Host "  .\run.ps1 -f           # Clean rebuild everything"
    Write-Host "  .\run.ps1 -s           # Just start the backend (skip build)"
    Write-Host "  .\run.ps1 -b           # Build only, don't start server"
    Write-Host "  .\run.ps1 -p -f        # Clean build without git pull"
    Write-Host "  .\run.ps1 -u           # Upload default plugin to WordPress"
    Write-Host "  .\run.ps1 -u -d        # Upload with debug output"
    Write-Host "  .\run.ps1 -u -pp 'C:\path'  # Upload custom plugin path"
    Write-Host ""
    Write-Host "CONFIGURATION:" -ForegroundColor Yellow
    Write-Host "  Config file: $ConfigPath"
    Write-Host "  Project: $ProjectName"
    Write-Host "  Backend: $BackendDir"
    Write-Host "  Frontend: $FrontendDir"
    if ($PnpmStorePath) {
        Write-Host "  pnpm Store: $PnpmStorePath"
    }
    Write-Host ""
    exit 0
}

# ============================================================================
# BANNER
# ============================================================================
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  $ProjectName - Build & Run Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($verbose) {
    Write-Host "Configuration:" -ForegroundColor Gray
    Write-Host "  Script Dir: $ScriptDir" -ForegroundColor Gray
    Write-Host "  Root Dir: $RootDir" -ForegroundColor Gray
    Write-Host "  Backend Dir: $BackendDir" -ForegroundColor Gray
    Write-Host "  Frontend Dir: $FrontendDir" -ForegroundColor Gray
    Write-Host "  pnpm Store: $PnpmStorePath" -ForegroundColor Gray
    Write-Host ""
}

# ============================================================================
# UTILITY FUNCTIONS
# ============================================================================

function Format-ElapsedTime($Stopwatch) {
    $elapsed = $Stopwatch.Elapsed
    if ($elapsed.TotalMinutes -ge 1) {
        return "{0:N0}m {1:N1}s" -f [Math]::Floor($elapsed.TotalMinutes), $elapsed.Seconds
    } else {
        return "{0:N1}s" -f $elapsed.TotalSeconds
    }
}

function Test-Command($Command) {
    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = 'SilentlyContinue'
    try { 
        $result = Get-Command $Command -ErrorAction SilentlyContinue
        return $null -ne $result
    }
    catch { return $false }
    finally { $ErrorActionPreference = $oldPreference }
}

function Test-IsAdmin {
    try {
        $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
        $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
        return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    } catch {
        return $false
    }
}

function Refresh-Path {
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + 
                [System.Environment]::GetEnvironmentVariable("Path", "User")
}

function Get-PnpmMajorVersion([string]$Version) {
    try {
        $major = ($Version -split '\.')[0]
        return [int]$major
    } catch {
        return 0
    }
}

function Get-NodeMajorVersion([string]$Version) {
    try {
        $v = $Version.Trim()
        if ($v.StartsWith('v')) { $v = $v.Substring(1) }
        $major = ($v -split '\.')[0]
        return [int]$major
    } catch {
        return 0
    }
}

function Get-DriveRoot([string]$Path) {
    try {
        if ([string]::IsNullOrWhiteSpace($Path)) { return $null }
        $full = [System.IO.Path]::GetFullPath($Path)
        if ($full -match '^[A-Za-z]:') { return $full.Substring(0, 2).ToUpper() }
        return $null
    } catch {
        return $null
    }
}

function Get-EffectivePnpmInstallCommand([string]$BaseCommand, [int]$Major) {
    $cmd = $BaseCommand
    if ($Major -ge 10 -and $cmd -match '(^|\s)pnpm\s+install(\s|$)' -and $cmd -notmatch 'dangerously-allow-all-builds') {
        $cmd = "$cmd --dangerously-allow-all-builds"
    }
    return $cmd
}

function Enable-PnpmPnpNodeOptions([string]$ProjectDir) {
    $pnpCjs = Join-Path $ProjectDir ".pnp.cjs"
    $pnpLoader = Join-Path $ProjectDir ".pnp.loader.mjs"
    $additions = @()

    if (Test-Path $pnpCjs) {
        if ([string]::IsNullOrWhiteSpace($env:NODE_OPTIONS) -or ($env:NODE_OPTIONS -notmatch [regex]::Escape($pnpCjs))) {
            $additions += "--require `"$pnpCjs`""
        }
    }

    if (Test-Path $pnpLoader) {
        if ([string]::IsNullOrWhiteSpace($env:NODE_OPTIONS) -or ($env:NODE_OPTIONS -notmatch [regex]::Escape($pnpLoader))) {
            $additions += "--experimental-loader `"$pnpLoader`""
        }
    }

    if ($additions.Count -gt 0) {
        $env:NODE_OPTIONS = (($env:NODE_OPTIONS + " " + ($additions -join " ")).Trim())
    }
}

# ============================================================================
# INSTALLATION FUNCTIONS
# ============================================================================

function Install-NodeJS {
    Write-Host "  Attempting to install Node.js via winget..." -ForegroundColor Yellow
    if (-not (Test-Command "winget")) {
        Write-Host "ERROR: winget not available. Install Node.js manually: https://nodejs.org/" -ForegroundColor Red
        exit 1
    }
    try {
        winget install OpenJS.NodeJS.LTS --accept-package-agreements --accept-source-agreements
        if ($LASTEXITCODE -ne 0) { throw "winget install failed" }
        Refresh-Path
        Write-Host "  ✓ Node.js installed successfully" -ForegroundColor Green
    } catch {
        Write-Host "ERROR: Failed to install Node.js. Install manually: https://nodejs.org/" -ForegroundColor Red
        exit 1
    }
}

function Install-Go {
    Write-Host "  Attempting to install Go via winget..." -ForegroundColor Yellow
    if (-not (Test-Command "winget")) {
        Write-Host "ERROR: winget not available. Install Go manually: https://go.dev/dl/" -ForegroundColor Red
        exit 1
    }
    try {
        winget install GoLang.Go --accept-package-agreements --accept-source-agreements
        if ($LASTEXITCODE -ne 0) { throw "winget install failed" }
        Refresh-Path
        Write-Host "  ✓ Go installed successfully" -ForegroundColor Green
    } catch {
        Write-Host "ERROR: Failed to install Go. Install manually: https://go.dev/dl/" -ForegroundColor Red
        exit 1
    }
}

function Install-Pnpm {
    Write-Host "  Installing pnpm globally..." -ForegroundColor Yellow
    try {
        npm install -g pnpm
        if ($LASTEXITCODE -ne 0) { throw "pnpm install failed" }
        Refresh-Path
        Write-Host "  ✓ pnpm installed successfully" -ForegroundColor Green
    } catch {
        Write-Host "ERROR: Failed to install pnpm. Run: npm install -g pnpm" -ForegroundColor Red
        exit 1
    }
}

function Configure-PnpmStore {
    $projectDrive = Get-DriveRoot $FrontendDir
    $storeDrive = Get-DriveRoot $PnpmStorePath
    $crossDrive = $false
    if ($projectDrive -and $storeDrive -and ($projectDrive -ne $storeDrive)) {
        $crossDrive = $true
    }

    $nodeLinker = "isolated"
    if ($UsePnp -and (-not $crossDrive) -and ($NodeMajor -lt 24)) {
        $nodeLinker = "pnp"
    }

    $script:EffectiveNodeLinker = $nodeLinker

    if ($UsePnp -and $nodeLinker -ne "pnp") {
        Write-Host "  NOTE: Falling back to node-linker=isolated (Node v$NodeMajor / cross-drive store)." -ForegroundColor Yellow
    }

    if ($PnpmStorePath) {
        Write-Host "  Configuring pnpm store: $PnpmStorePath" -ForegroundColor Gray
        if (-not (Test-Path $PnpmStorePath)) {
            New-Item -ItemType Directory -Path $PnpmStorePath -Force | Out-Null
        }
        pnpm config set --location=project store-dir $PnpmStorePath 2>$null
    }

    pnpm config set --location=project virtual-store-dir .pnpm 2>$null
    pnpm config set --location=project node-linker $nodeLinker 2>$null

    if ($nodeLinker -eq "pnp") {
        pnpm config set --location=project symlink false 2>$null
    } else {
        pnpm config set --location=project symlink true 2>$null
    }

    pnpm config set --location=project package-import-method auto 2>$null
}

function Ensure-FirewallRules {
    param([int[]]$PortList = @(8080))

    if (-not (Test-IsAdmin)) {
        Write-Host "  WARNING: -OpenFirewall requires Administrator." -ForegroundColor Yellow
        return
    }

    if (-not (Test-Command "New-NetFirewallRule")) {
        Write-Host "  WARNING: New-NetFirewallRule not available." -ForegroundColor Yellow
        return
    }

    foreach ($p in $PortList) {
        $ruleName = "$ProjectName (Backend) TCP $p"
        $existing = Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue
        if ($null -eq $existing) {
            New-NetFirewallRule `
                -DisplayName $ruleName `
                -Direction Inbound `
                -Action Allow `
                -Protocol TCP `
                -LocalPort $p `
                -Profile Private,Domain `
                | Out-Null
            Write-Host "  ✓ Firewall rule added: $ruleName" -ForegroundColor Green
        } else {
            Write-Host "  ✓ Firewall rule exists: $ruleName" -ForegroundColor Green
        }
    }
}

# ============================================================================
# STEP TRACKING
# ============================================================================
$StepTimes = @{}

# Set environment variables from config
if ($Config.env) {
    foreach ($key in $Config.env.PSObject.Properties.Name) {
        [System.Environment]::SetEnvironmentVariable($key, $Config.env.$key, "Process")
    }
}

# ============================================================================
# STEP 1: GIT PULL
# ============================================================================
$stepWatch = [System.Diagnostics.Stopwatch]::StartNew()
if (-not $skippull) {
    Write-Host "[1/5] Pulling latest changes from git..." -ForegroundColor Yellow
    Push-Location $RootDir
    try {
        if (Test-Path ".git") {
            git pull 2>&1 | Out-Host
            if ($LASTEXITCODE -ne 0) {
                Write-Host "  WARNING: git pull failed, continuing..." -ForegroundColor Yellow
            } else {
                Write-Host "  ✓ Git pull complete" -ForegroundColor Green
            }
        } else {
            Write-Host "  Skipping (not a git repository)" -ForegroundColor Gray
        }
    }
    finally { Pop-Location }
} else {
    Write-Host "[1/5] Skipping git pull (-p)" -ForegroundColor Gray
}
$stepWatch.Stop()
$StepTimes["Git Pull"] = $stepWatch.Elapsed
Write-Host "  ⏱ $(Format-ElapsedTime $stepWatch)" -ForegroundColor DarkGray
Write-Host ""

# ============================================================================
# STEP 2: PREREQUISITES
# ============================================================================
$stepWatch = [System.Diagnostics.Stopwatch]::StartNew()
Write-Host "[2/5] Checking prerequisites..." -ForegroundColor Yellow

if ($CheckGo) {
    if (-not (Test-Command "go")) { Install-Go }
    $goVersion = (go version 2>&1) -replace 'go version ', ''
    Write-Host "  ✓ Go found: $goVersion" -ForegroundColor Green
}

if ($CheckNode) {
    if (-not (Test-Command "node")) { Install-NodeJS }
    $nodeVersion = node --version 2>&1
    Write-Host "  ✓ Node.js found: $nodeVersion" -ForegroundColor Green
    $NodeMajor = Get-NodeMajorVersion $nodeVersion
}

if ($CheckPnpm) {
    if (-not (Test-Command "pnpm")) { Install-Pnpm }
    $pnpmVersion = pnpm --version 2>&1
    Write-Host "  ✓ pnpm found: $pnpmVersion" -ForegroundColor Green
    $PnpmMajor = Get-PnpmMajorVersion $pnpmVersion
    $EffectiveInstallCommand = Get-EffectivePnpmInstallCommand $InstallCommand $PnpmMajor
    Configure-PnpmStore
}

$stepWatch.Stop()
$StepTimes["Prerequisites"] = $stepWatch.Elapsed
Write-Host "  ⏱ $(Format-ElapsedTime $stepWatch)" -ForegroundColor DarkGray
Write-Host ""

# ============================================================================
# UPLOAD MODE (-u): Upload plugin to WordPress via upload-plugin-v2.ps1
# This section is optional — only runs when -u flag is passed.
# Requires "upload" section in powershell.json with:
#   scriptPath:        Path to upload-plugin-v2.ps1
#   defaultPluginPath: Default plugin folder to upload
#   configPath:        Path to wp-plugin-config.json (credentials)
# ============================================================================
if ($upload) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  Upload Mode (-u)" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    
    # Resolve upload config from powershell.json
    $uploadConfig = $Config.upload
    if (-not $uploadConfig) {
        Write-Host "ERROR: No 'upload' section in powershell.json." -ForegroundColor Red
        Write-Host "Add this to powershell.json:" -ForegroundColor Yellow
        Write-Host '  "upload": {' -ForegroundColor Gray
        Write-Host '    "scriptPath": "path/to/upload-plugin-v2.ps1",' -ForegroundColor Gray
        Write-Host '    "defaultPluginPath": "path/to/plugin-folder",' -ForegroundColor Gray
        Write-Host '    "configPath": "path/to/wp-plugin-config.json"' -ForegroundColor Gray
        Write-Host '  }' -ForegroundColor Gray
        exit 1
    }
    
    $uploadScript = Resolve-RelativePath $uploadConfig.scriptPath
    $defaultPlugin = Resolve-RelativePath $uploadConfig.defaultPluginPath
    $wpConfig = Resolve-RelativePath $uploadConfig.configPath
    
    if (-not (Test-Path $uploadScript)) {
        Write-Host "ERROR: Upload script not found: $uploadScript" -ForegroundColor Red
        exit 1
    }
    
    # Determine plugin path (CLI override or default from config)
    $targetPlugin = if ($pluginpath -ne "") { $pluginpath } else { $defaultPlugin }
    
    if (-not (Test-Path $targetPlugin)) {
        Write-Host "ERROR: Plugin folder not found: $targetPlugin" -ForegroundColor Red
        exit 1
    }
    
    $pluginName = Split-Path $targetPlugin -Leaf
    Write-Host "  Plugin: $pluginName" -ForegroundColor White
    Write-Host "  Path:   $targetPlugin" -ForegroundColor Gray
    
    # Read config to show site URL
    if (Test-Path $wpConfig) {
        try {
            $wpConfigData = Get-Content $wpConfig -Raw | ConvertFrom-Json
            Write-Host "  Site:   $($wpConfigData.wordPressSiteURL)" -ForegroundColor Gray
        } catch {}
    }
    Write-Host ""
    
    # Build JSON config for V2 script
    if (Test-Path $wpConfig) {
        $configContent = Get-Content $wpConfig -Raw | ConvertFrom-Json
        $configContent.pluginFolderPath = $targetPlugin
        $jsonConfig = $configContent | ConvertTo-Json -Compress
        Write-Host "Parsing inline JSON config..." -ForegroundColor Gray
    } else {
        Write-Host "ERROR: Config file not found: $wpConfig" -ForegroundColor Red
        exit 1
    }
    
    # Build V2 arguments
    $v2Args = @("-JsonConfig", $jsonConfig)
    if ($debugmode) { $v2Args += "-DebugMode" }
    
    & $uploadScript @v2Args
    exit $LASTEXITCODE
}

# ============================================================================
# INSTALL MODE (-i): Install dependencies for both frontend and backend
# ============================================================================
if ($install) {
    $stepWatch = [System.Diagnostics.Stopwatch]::StartNew()
    Write-Host "[INSTALL] Installing/updating all dependencies..." -ForegroundColor Cyan
    Write-Host ""
    
    # Frontend
    if ($rebuild) {
        Write-Host "  [Frontend] Rebuild mode: deferring until after force-clean..." -ForegroundColor Yellow
    } else {
        Write-Host "  [Frontend] Running pnpm install..." -ForegroundColor Yellow
        Push-Location $FrontendDir
        try {
            Configure-PnpmStore
            Invoke-Expression $EffectiveInstallCommand
            if ($LASTEXITCODE -ne 0) { throw "pnpm install failed" }
            $DidFrontendInstall = $true
            Write-Host "  ✓ Frontend dependencies installed" -ForegroundColor Green
        }
        finally { Pop-Location }
    }
    
    # Backend
    Write-Host ""
    Write-Host "  [Backend] Running go mod tidy && go mod download..." -ForegroundColor Yellow
    Push-Location $BackendDir
    try {
        go mod tidy
        if ($LASTEXITCODE -ne 0) { throw "go mod tidy failed" }
        go mod download
        if ($LASTEXITCODE -ne 0) { throw "go mod download failed" }
        Write-Host "  ✓ Backend dependencies installed" -ForegroundColor Green
    }
    finally { Pop-Location }
    
    $stepWatch.Stop()
    $StepTimes["Install Dependencies"] = $stepWatch.Elapsed
    
    if (-not $rebuild) {
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "  Dependencies installed!" -ForegroundColor Cyan
        Write-Host "  Time: $(Format-ElapsedTime $stepWatch)" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        exit 0
    }
}

# ============================================================================
# STEP 3: FRONTEND BUILD
# ============================================================================
$stepWatch = [System.Diagnostics.Stopwatch]::StartNew()
if (-not $skipbuild) {
    Write-Host "[3/5] Building React frontend..." -ForegroundColor Yellow
    
    Push-Location $FrontendDir
    try {
        # Force clean
        if ($force) {
            Write-Host "  FORCE MODE: Cleaning build artifacts..." -ForegroundColor Magenta
            
            foreach ($cleanPath in $CleanPaths) {
                if ($cleanPath -match '\*') {
                    $resolvedPath = Resolve-RelativePath ($cleanPath -replace '\*.*$', '')
                    $pattern = $cleanPath -replace '^.*[\\/]', ''
                    if (Test-Path $resolvedPath) {
                        Get-ChildItem -Path $resolvedPath -Filter $pattern -ErrorAction SilentlyContinue | ForEach-Object {
                            Write-Host "  Removing: $($_.Name)..." -ForegroundColor Gray
                            Remove-Item -Force -Recurse $_.FullName -ErrorAction SilentlyContinue
                        }
                    }
                } else {
                    $resolvedPath = Resolve-RelativePath $cleanPath
                    if (Test-Path $resolvedPath) {
                        Write-Host "  Removing: $cleanPath..." -ForegroundColor Gray
                        Remove-Item -Recurse -Force $resolvedPath -ErrorAction SilentlyContinue
                    }
                }
            }

            # Always clean pnpm artifacts
            foreach ($extraPath in @("node_modules", ".pnpm", ".pnp.cjs", ".pnp.loader.mjs", ".pnp.data.json")) {
                if (Test-Path $extraPath) {
                    Write-Host "  Removing: $extraPath..." -ForegroundColor Gray
                    Remove-Item -Recurse -Force $extraPath -ErrorAction SilentlyContinue
                }
            }
            
            if ($CheckPnpm) {
                Write-Host "  Clearing pnpm cache..." -ForegroundColor Gray
                pnpm store prune 2>&1 | Out-Null
            }

            # Clean backend runtime data if dataDir configured
            if ($DataDir) {
                foreach ($rtPath in @(
                    (Join-Path $DataDir "sessions"),
                    (Join-Path $DataDir "request-sessions"),
                    (Join-Path $DataDir "errors")
                )) {
                    if (Test-Path $rtPath) {
                        Write-Host "  Removing: $rtPath..." -ForegroundColor Gray
                        Remove-Item -Recurse -Force $rtPath -ErrorAction SilentlyContinue
                    }
                }
                foreach ($logFile in @("log.txt", "error.log.txt")) {
                    $logPath = Join-Path $DataDir $logFile
                    if (Test-Path $logPath) {
                        Remove-Item -Force $logPath -ErrorAction SilentlyContinue
                    }
                }
            }
            
            Write-Host "  ✓ Clean complete" -ForegroundColor Magenta
        }
        
        # Check if install needed
        $depsPresent = if ($EffectiveNodeLinker -eq "pnp") { (Test-Path ".pnp.cjs") } else { (Test-Path "node_modules") }
        $NeedsInstall = $install -or (-not $depsPresent)
        
        if (-not $NeedsInstall -and $EffectiveNodeLinker -ne "pnp" -and $RequiredModules.Count -gt 0) {
            foreach ($m in $RequiredModules) {
                if (-not (Test-Path (Join-Path "node_modules" $m))) {
                    $NeedsInstall = $true
                    break
                }
            }
        }

        if ($NeedsInstall -or $force) {
            Write-Host "  Installing dependencies..." -ForegroundColor Gray
            Invoke-Expression $EffectiveInstallCommand
            if ($LASTEXITCODE -ne 0) { throw "pnpm install failed" }
        }
        
        # Build
        Write-Host "  Running: $BuildCommand" -ForegroundColor Gray
        $oldNodeOptions = $env:NODE_OPTIONS
        try {
            if ($EffectiveNodeLinker -eq "pnp") {
                Enable-PnpmPnpNodeOptions -ProjectDir (Get-Location)
            }
            Invoke-Expression $BuildCommand
            if ($LASTEXITCODE -ne 0) { throw "Build failed" }
        }
        finally { $env:NODE_OPTIONS = $oldNodeOptions }
        
        Write-Host "  ✓ Frontend built successfully" -ForegroundColor Green
    }
    finally { Pop-Location }
    
    $stepWatch.Stop()
    $StepTimes["Frontend Build"] = $stepWatch.Elapsed
    Write-Host "  ⏱ $(Format-ElapsedTime $stepWatch)" -ForegroundColor DarkGray
    Write-Host ""
    
    # STEP 4: COPY BUILD TO BACKEND
    $stepWatch = [System.Diagnostics.Stopwatch]::StartNew()
    if ($TargetDir) {
        Write-Host "[4/5] Copying build to backend..." -ForegroundColor Yellow
        $SourceDist = Join-Path $FrontendDir $DistDir
        if (-not (Test-Path $SourceDist)) {
            Write-Host "  WARNING: Build output not found: $SourceDist" -ForegroundColor Yellow
        } else {
            $TargetParent = Split-Path -Parent $TargetDir
            if (-not (Test-Path $TargetParent)) {
                New-Item -ItemType Directory -Path $TargetParent -Force | Out-Null
            }
            if (Test-Path $TargetDir) { Remove-Item -Recurse -Force $TargetDir }
            Copy-Item -Recurse $SourceDist $TargetDir
            Write-Host "  ✓ Copied to: $TargetDir" -ForegroundColor Green
        }
    } else {
        Write-Host "[4/5] Skipping copy (no targetDir)" -ForegroundColor Gray
    }
    $stepWatch.Stop()
    $StepTimes["Copy Build"] = $stepWatch.Elapsed
} else {
    Write-Host "[3/5] Skipping frontend build (-s)" -ForegroundColor Gray
    Write-Host "[4/5] Skipping copy" -ForegroundColor Gray
    $stepWatch.Stop()
    $StepTimes["Frontend Build"] = [TimeSpan]::Zero
    $StepTimes["Copy Build"] = [TimeSpan]::Zero
}
Write-Host ""

# BUILD ONLY EXIT
if ($buildonly) {
    $TotalStopwatch.Stop()
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  Build complete! (-b mode)" -ForegroundColor Cyan
    Write-Host "  Total time: $(Format-ElapsedTime $TotalStopwatch)" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Build Summary:" -ForegroundColor Yellow
    foreach ($step in $StepTimes.GetEnumerator()) {
        $time = "{0:N1}s" -f $step.Value.TotalSeconds
        Write-Host "  $($step.Key): $time" -ForegroundColor Gray
    }
    exit 0
}

# ============================================================================
# STEP 5: START BACKEND
# ============================================================================
Write-Host "[5/5] Starting Go backend..." -ForegroundColor Yellow

Push-Location $BackendDir
try {
    $BackendConfigPath = Join-Path $BackendDir $ConfigFile
    $BackendConfigExample = Join-Path $BackendDir $ConfigExampleFile
    
    if (-not (Test-Path $BackendConfigPath)) {
        if (Test-Path $BackendConfigExample) {
            Write-Host "  Creating $ConfigFile from $ConfigExampleFile..." -ForegroundColor Gray
            Copy-Item $BackendConfigExample $BackendConfigPath
        } else {
            Write-Host "  WARNING: No $ConfigFile found" -ForegroundColor Yellow
        }
    }
    
    if ($DataDir -and -not (Test-Path $DataDir)) {
        New-Item -ItemType Directory -Path $DataDir -Force | Out-Null
    }

    if ($openfirewall) {
        Ensure-FirewallRules -PortList $Ports
    }
    
    $TotalStopwatch.Stop()
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  $ProjectName starting..." -ForegroundColor Cyan
    Write-Host "  Open: http://localhost:$($Ports[0])" -ForegroundColor Cyan
    Write-Host "  Press Ctrl+C to stop" -ForegroundColor Cyan
    Write-Host "  Build time: $(Format-ElapsedTime $TotalStopwatch)" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    
    Invoke-Expression $RunCommand
}
finally { Pop-Location }
