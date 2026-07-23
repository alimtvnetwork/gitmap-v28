<#
.SYNOPSIS
    Build, deploy, and run gitmap CLI from the repo root.
.DESCRIPTION
    Pulls latest code, resolves Go dependencies, builds the binary
    into ./bin, copies data folder, deploys to a target directory,
    and optionally runs gitmap with any arguments.
.EXAMPLES
    .\run.ps1                                    # pull, build, deploy
    .\run.ps1 -NoPull                            # skip git pull
    .\run.ps1 -ForcePull                         # discard local changes + pull (no prompt)
    .\run.ps1 -NoDeploy                          # skip deploy step
    .\run.ps1 -uninstall                         # run uninstall-quick.ps1 -Yes and exit
    .\run.ps1 -reinstall                         # uninstall, then re-run run.ps1 with no args
    .\run.ps1 -R scan                            # build + scan parent folder
    .\run.ps1 -R scan D:\repos                   # build + scan specific path
    .\run.ps1 -R scan D:\repos --mode ssh        # build + scan with flags
    .\run.ps1 -R clone .\gitmap-output\gitmap.json --target-dir .\restored
    .\run.ps1 -R help                            # build + show help
    .\run.ps1 -NoPull -NoDeploy -R scan          # just build and scan
    .\run.ps1 -t                                 # run all unit tests with reports
.NOTES
    Configuration is read from gitmap/powershell.json.
    -R accepts ALL gitmap CLI arguments after it (scan, clone, help, flags, paths).
    If -R is used with no arguments, it defaults to: scan <parent folder>
    -t runs all Go unit tests and writes reports to gitmap/data/unit-test-reports/.
    -ForcePull automatically discards local changes and removes untracked files
    before pulling. Useful for CI or unattended builds.
#>

[CmdletBinding(PositionalBinding=$false)]
param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [switch]$NoSetup,
    [switch]$ForcePull,
    [string]$DeployPath = "",
    [Alias("d")]
    [switch]$Deploy,
    [switch]$Update,
    [switch]$R,
    [Alias("t")]
    [switch]$Test,
    [Alias("u")]
    [switch]$Uninstall,
    [Alias("ri")]
    [switch]$Reinstall,
    [switch]$DebugRepoDetect,
    [switch]$Quiet,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
)

# Honor env-var bridge from `gitmap update --quiet` (or callers that set it).
if ($env:GITMAP_QUIET -eq "1") { $Quiet = $true }

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$GitMapDir = Join-Path $RepoRoot "gitmap"

# -- Logging helpers -------------------------------------------
function Write-Step {
    param([string]$Step, [string]$Message)
    Write-Host ""
    Write-Host "  [$Step] " -ForegroundColor Magenta -NoNewline
    Write-Host $Message -ForegroundColor White
    Write-Host ("  " + ("-" * 50)) -ForegroundColor DarkGray
}

function Write-Success {
    param([string]$Message)
    Write-Host "  OK " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Green
}

function Write-Info {
    param([string]$Message)
    Write-Host "  -> " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor Gray
}

function Write-Warn {
    param([string]$Message)
    Write-Host "  !! " -ForegroundColor Yellow -NoNewline
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Fail {
    param([string]$Message)
    Write-Host "  XX " -ForegroundColor Red -NoNewline
    Write-Host $Message -ForegroundColor Red
}

# -- Error reporting (JSONL) -----------------------------------
# When run from `gitmap update --report-errors json`, env vars
# GITMAP_REPORT_ERRORS=json and GITMAP_REPORT_ERRORS_FILE=<path>
# are set. Each non-fatal failure appends one JSON object per line.
function Write-ReportError {
    param(
        [string]$Stage,
        [string]$Command,
        [int]$ExitCode,
        [string]$Message,
        [hashtable]$Paths = @{}
    )
    $fmt = $env:GITMAP_REPORT_ERRORS
    $file = $env:GITMAP_REPORT_ERRORS_FILE
    if (($fmt -ne "json") -or [string]::IsNullOrWhiteSpace($file)) {
        return
    }
    try {
        $entry = [ordered]@{
            timestamp = (Get-Date).ToUniversalTime().ToString("o")
            stage     = $Stage
            command   = $Command
            exitCode  = $ExitCode
            cwd       = (Get-Location).Path
            message   = $Message
            paths     = $Paths
            os        = "windows"
        }
        $line = ($entry | ConvertTo-Json -Compress -Depth 5)
        Add-Content -Path $file -Value $line -Encoding UTF8 -ErrorAction Stop
    } catch {
        Write-Host "  [WARN] Could not write report-errors entry: $_" -ForegroundColor Yellow
    }
}

# -- Repo-detect debug -----------------------------------------
# Emits structured diagnostics describing why docs auto-build was skipped
# or executed. Active when -DebugRepoDetect is passed OR
# $env:GITMAP_DEBUG_REPO_DETECT is set (typically by `gitmap update
# --debug-repo-detect`). Output is also mirrored to the report file when
# --report-errors json is active, with stage="repo-detect" and level=info.
function Test-DebugRepoDetect {
    if ($script:DebugRepoDetect) { return $true }
    return ($env:GITMAP_DEBUG_REPO_DETECT -eq "1")
}

function Write-RepoDetect {
    param(
        [string]$Check,
        [string]$Result,
        [string]$Detail = ""
    )
    if (-not (Test-DebugRepoDetect)) { return }
    Write-Host "  [DETECT] " -ForegroundColor DarkCyan -NoNewline
    Write-Host ("{0,-28} = " -f $Check) -ForegroundColor Cyan -NoNewline
    Write-Host $Result -ForegroundColor White -NoNewline
    if ($Detail.Length -gt 0) {
        Write-Host "  ($Detail)" -ForegroundColor DarkGray
    } else {
        Write-Host ""
    }
    # Mirror to JSONL report file when active.
    Write-ReportError -Stage "repo-detect" -Command $Check -ExitCode 0 `
        -Message $Result -Paths @{ detail = $Detail; level = "info" }
}

function Write-RepoDetectSnippet {
    param([string]$Title, [string]$Path, [int]$MaxLines = 6)
    if (-not (Test-DebugRepoDetect)) { return }
    Write-Host "  [DETECT] $Title :" -ForegroundColor Cyan
    if (-not (Test-Path $Path)) {
        Write-Host "    (file not found: $Path)" -ForegroundColor DarkGray
        return
    }
    try {
        $lines = Get-Content -Path $Path -TotalCount $MaxLines -ErrorAction Stop
        foreach ($l in $lines) {
            Write-Host "    $l" -ForegroundColor DarkGray
        }
    } catch {
        Write-Host "    (could not read: $_)" -ForegroundColor DarkGray
    }
}

# -- npm wrapper ----------------------------------------------
# Invoke-NpmQuiet runs an npm subcommand (e.g. install, run build) with
# $ErrorActionPreference temporarily relaxed to 'Continue' so that npm's
# stderr progress chatter is NOT promoted to a terminating
# NativeCommandError under the script-wide 'Stop' preference. The previous
# preference is ALWAYS restored via finally{}, even on exceptions.
#
# When $Quiet (or $env:GITMAP_QUIET=1) is active, all npm output is
# discarded and only a single "[npm] <cmd> exit=<code>" line is logged.
# Otherwise output streams to the terminal as usual. The function returns
# the npm exit code so callers can branch on it.
function Invoke-NpmQuiet {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string[]]$NpmArgs
    )

    $prevEAP = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    $exitCode = 0
    try {
        if ($script:Quiet) {
            & npm @NpmArgs *>&1 | Out-Null
        } else {
            # Stream npm output to the host without letting it enter this
            # function's output pipeline. Without Out-Host, every stdout line
            # from npm becomes part of the function's return value, so callers
            # using `return $exitCode` would get an Object[] (npm lines + the
            # int) instead of a single Int32 — which then breaks any param
            # typed as [int]ExitCode downstream.
            & npm @NpmArgs 2>&1 | Out-Host
        }
        $exitCode = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $prevEAP
    }

    if ($null -eq $exitCode) { $exitCode = 0 }

    if ($script:Quiet) {
        Write-Host ("  [npm] {0} exit={1}" -f ($NpmArgs -join ' '), $exitCode) -ForegroundColor DarkGray
    }

    # Force scalar [int] return so the caller never receives an array even
    # if some upstream change reintroduces stray pipeline output.
    return [int]$exitCode
}

# -- Banner ----------------------------------------------------
function Show-Banner {
    Write-Host ""
    Write-Host "  +======================================+" -ForegroundColor DarkCyan
    Write-Host "  |         " -ForegroundColor DarkCyan -NoNewline
    Write-Host "gitmap builder" -ForegroundColor Cyan -NoNewline
    Write-Host "              |" -ForegroundColor DarkCyan
    Write-Host "  +======================================+" -ForegroundColor DarkCyan
    Write-Host ""
}

# -- Load deploy manifest (single source of truth) -------------
# Mirrors run.sh's load_deploy_manifest. Reads
# gitmap/constants/deploy-manifest.json so AppSubdir / LegacyAppSubdirs
# aren't hardcoded across run.ps1, run.sh, install.sh, and Go constants.
# Renaming the deploy folder ONLY requires editing that JSON file.
$script:AppSubdir = "gitmap-cli"
$script:LegacyAppSubdirs = @("gitmap")
function Get-DeployManifest {
    $manifestPath = Join-Path $GitMapDir "constants/deploy-manifest.json"
    if (-not (Test-Path $manifestPath)) {
        Write-Warn "deploy-manifest.json not found at $manifestPath - using defaults"
        return
    }
    try {
        $manifest = Get-Content $manifestPath -Raw | ConvertFrom-Json
        if ($manifest.appSubdir) {
            $script:AppSubdir = $manifest.appSubdir
        }
        if ($manifest.legacyAppSubdirs) {
            $script:LegacyAppSubdirs = @($manifest.legacyAppSubdirs)
        }
    } catch {
        Write-Warn "Failed to parse deploy-manifest.json: $_"
    }
}

# Test-KnownAppSubdir returns $true if $Name matches AppSubdir or any legacy entry.
function Test-KnownAppSubdir {
    param([string]$Name)
    if ($Name -eq $script:AppSubdir) { return $true }
    foreach ($legacy in $script:LegacyAppSubdirs) {
        if ($Name -eq $legacy) { return $true }
    }
    return $false
}

# -- Load config -----------------------------------------------
function Load-Config {
    $configPath = Join-Path $GitMapDir "powershell.json"
    if (Test-Path $configPath) {
        Write-Info "Config loaded from powershell.json"

        return Get-Content $configPath | ConvertFrom-Json
    }
    Write-Warn "No powershell.json found, using defaults"

    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "gitmap.exe"
        copyData    = $true
    }
}

# -- Ensure main branch ----------------------------------------
function Ensure-MainBranch {
    Push-Location $RepoRoot
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $currentBranch = (git rev-parse --abbrev-ref HEAD 2>&1).Trim()
        $ErrorActionPreference = $prevPref

        if ($currentBranch -ne "main") {
            Write-Warn "Currently on branch '$currentBranch', switching to main..."
            $ErrorActionPreference = "Continue"
            $checkoutOutput = git checkout main 2>&1
            $checkoutExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($checkoutExit -ne 0) {
                Write-Fail "Failed to switch to main branch"
                foreach ($line in $checkoutOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Switched to main branch"
        }
    } finally {
        Pop-Location
    }
}

# -- Git pull --------------------------------------------------
function Invoke-GitPull {
    Write-Step "1/4" "Pulling latest changes"

    Ensure-MainBranch

    Push-Location $RepoRoot
    try {
        # Temporarily allow stderr output from git without throwing NativeCommandError.
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $output = git pull 2>&1
        $pullExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        foreach ($line in $output) {
            $text = "$line".Trim()
            if ($text.Length -gt 0) {
                Write-Info $text
            }
        }

        if ($pullExit -ne 0) {
            $outputText = ($output | ForEach-Object { "$_" }) -join "`n"
            $hasConflict = $outputText -match "Your local changes" -or
                           $outputText -match "overwritten by merge" -or
                           $outputText -match "not possible because you have unmerged" -or
                           $outputText -match "Please commit your changes or stash them"

            if ($hasConflict) {
                if ($ForcePull) {
                    Write-Warn "Force-pull: discarding local changes and removing untracked files..."
                    $prevPref = $ErrorActionPreference
                    $ErrorActionPreference = "Continue"

                    $resetOutput = git checkout -- . 2>&1
                    $resetExit = $LASTEXITCODE
                    if ($resetExit -ne 0) {
                        Write-Fail "Git checkout failed"
                        $ErrorActionPreference = $prevPref
                        exit 1
                    }
                    Write-Success "Local changes discarded"

                    $cleanOutput = git clean -fd 2>&1
                    $cleanExit = $LASTEXITCODE
                    $ErrorActionPreference = $prevPref

                    if ($cleanExit -ne 0) {
                        Write-Fail "Git clean failed"
                        exit 1
                    }

                    $cleanedFiles = @($cleanOutput | ForEach-Object { "$_".Trim() } | Where-Object { $_.Length -gt 0 })
                    if ($cleanedFiles.Count -gt 0) {
                        Write-Success "Removed $($cleanedFiles.Count) untracked file(s)"
                    }

                    Retry-GitPull
                } else {
                    Resolve-PullConflict
                }
            } else {
                Write-Fail "Git pull failed (exit code $pullExit)"
                exit 1
            }
        } else {
            Write-Success "Pull complete"
        }
    } finally {
        Pop-Location
    }
}

# -- Resolve pull conflict with local changes ------------------
function Resolve-PullConflict {
    Write-Warn "Git pull failed due to local changes"
    Write-Host ""
    Write-Host "  Choose how to proceed:" -ForegroundColor Yellow
    Write-Host "    [S] Stash changes (save for later, then pull)" -ForegroundColor Cyan
    Write-Host "    [D] Discard changes (reset working tree, then pull)" -ForegroundColor Cyan
    Write-Host "    [C] Clean all (discard changes + remove untracked files, then pull)" -ForegroundColor Cyan
    Write-Host "    [Q] Quit (abort without changes)" -ForegroundColor Cyan
    Write-Host ""

    $choice = Read-Host "  Enter choice (S/D/C/Q)"

    switch ($choice.ToUpper()) {
        "S" {
            Write-Info "Stashing local changes..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $stashOutput = git stash push -m "auto-stash before run.ps1 pull" 2>&1
            $stashExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($stashExit -ne 0) {
                Write-Fail "Git stash failed"
                foreach ($line in $stashOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Changes stashed"
            Write-Info "Run 'git stash pop' later to restore your changes"

            Retry-GitPull
        }
        "D" {
            Write-Warn "Discarding all local changes..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $resetOutput = git checkout -- . 2>&1
            $resetExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($resetExit -ne 0) {
                Write-Fail "Git checkout failed"
                foreach ($line in $resetOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Local changes discarded"

            Retry-GitPull
        }
        "C" {
            Write-Warn "Discarding all local changes and removing untracked files..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"

            $resetOutput = git checkout -- . 2>&1
            $resetExit = $LASTEXITCODE

            if ($resetExit -ne 0) {
                Write-Fail "Git checkout failed"
                foreach ($line in $resetOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                $ErrorActionPreference = $prevPref
                exit 1
            }
            Write-Success "Local changes discarded"

            $cleanOutput = git clean -fd 2>&1
            $cleanExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($cleanExit -ne 0) {
                Write-Fail "Git clean failed"
                foreach ($line in $cleanOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }

            $cleanedFiles = @($cleanOutput | ForEach-Object { "$_".Trim() } | Where-Object { $_.Length -gt 0 })
            if ($cleanedFiles.Count -gt 0) {
                foreach ($line in $cleanedFiles) {
                    Write-Info $line
                }
                Write-Success "Removed $($cleanedFiles.Count) untracked file(s)"
            } else {
                Write-Info "No untracked files to remove"
            }

            Retry-GitPull
        }
        default {
            Write-Info "Aborted by user"
            exit 0
        }
    }
}

# -- Retry git pull after stash/discard -----------------------
function Retry-GitPull {
    Write-Info "Retrying git pull..."
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $retryOutput = git pull 2>&1
    $retryExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref

    foreach ($line in $retryOutput) {
        $text = "$line".Trim()
        if ($text.Length -gt 0) {
            Write-Info $text
        }
    }

    if ($retryExit -ne 0) {
        Write-Fail "Git pull failed again (exit code $retryExit)"
        exit 1
    }

    Write-Success "Pull complete"
}

# -- Toolchain preflight --------------------------------------
function Test-GoToolchain {
    $goCmd = Get-Command go -ErrorAction SilentlyContinue
    if ($goCmd) { return }
    Write-Fail "Go toolchain not found on PATH"
    Write-Host ""
    Write-Host "  The 'go' command is required to build gitmap, but it is not" -ForegroundColor Yellow
    Write-Host "  available in this PowerShell session." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  Fix options:" -ForegroundColor Cyan
    Write-Host "    1. Install Go from https://go.dev/dl/  (use the v1.24.x MSI)"
    Write-Host "    2. If Go is already installed, add it to PATH for this session:"
    Write-Host "         `$env:Path = 'C:\Program Files\Go\bin;' + `$env:Path"
    Write-Host "    3. Then re-open PowerShell (or run .\run.ps1 again)."
    Write-Host ""
    Write-Host "  Verify with:  go version" -ForegroundColor Cyan
    Write-Host ""
    exit 1
}

# -- Resolve dependencies -------------------------------------
function Resolve-Dependencies {
    Write-Step "2/4" "Resolving Go dependencies"
    Push-Location $GitMapDir
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $tidyOutput = go mod tidy 2>&1
        $tidyExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($tidyExit -ne 0) {
            Write-Fail "go mod tidy failed"
            foreach ($line in $tidyOutput) {
                Write-Host "  $line" -ForegroundColor Red
            }
            exit 1
        }
        Write-Success "Dependencies resolved"
    } finally {
        Pop-Location
    }
}

# -- Pre-build validation --------------------------------------
function Test-SourceFiles {
    Write-Info "Validating source files..."

    $requiredFiles = @(
        "main.go",
        "go.mod",
        "cmd/root.go",
        "cmd/scan.go",
        "cmd/clone.go",
        "cmd/update.go",
        "cmd/pull.go",
        "cmd/rescan.go",
        "cmd/desktopsync.go",
        "constants/constants.go",
        "config/config.go",
        "scanner/scanner.go",
        "mapper/mapper.go",
        "model/record.go",
        "formatter/csv.go",
        "formatter/json.go",
        "formatter/terminal.go",
        "formatter/text.go",
        "formatter/structure.go",
        "formatter/clonescript.go",
        "formatter/directclone.go",
        "formatter/desktopscript.go",
        "cloner/cloner.go",
        "cloner/safe_pull.go",
        "gitutil/gitutil.go",
        "desktop/desktop.go",
        "verbose/verbose.go",
        "setup/setup.go",
        "cmd/setup.go",
        "cmd/status.go",
        "cmd/exec.go",
        "cmd/release.go",
        "cmd/releasebranch.go",
        "cmd/releasepending.go",
        "cmd/changelog.go",
        "cmd/doctor.go",
        "release/semver.go",
        "release/metadata.go",
        "release/gitops.go",
        "release/github.go",
        "release/changelog.go",
        "release/workflow.go"
    )

    $missing = @()
    foreach ($file in $requiredFiles) {
        $fullPath = Join-Path $GitMapDir $file
        if (-not (Test-Path $fullPath)) {
            $missing += $file
        }
    }

    if ($missing.Count -gt 0) {
        Write-Fail "Missing source files ($($missing.Count)):"
        foreach ($f in $missing) {
            Write-Host "  - $f" -ForegroundColor Red
        }
        exit 1
    }

    Write-Success "All $($requiredFiles.Count) source files present"
}

# -- Build binary ----------------------------------------------
function Build-Binary {
    param($Config)

    # Step 2b/4: Embed Windows icon + version metadata via go-winres
    $winresDir = Join-Path $GitMapDir "winres"
    $winresJson = Join-Path $winresDir "winres.json"
    if (Test-Path $winresJson) {
        Write-Step "2b/4" "Embedding Windows icon (go-winres)"
        $goWinres = Get-Command go-winres -ErrorAction SilentlyContinue
        if (-not $goWinres) {
            Write-Info "go-winres not found; installing pinned v0.3.3"
            # Native commands (go) write progress like "go: downloading ..." to stderr.
            # With $ErrorActionPreference='Stop' that stderr is promoted to a terminating
            # RemoteException even on success. Locally relax the preference and rely on
            # $LASTEXITCODE — the only reliable success signal for native binaries.
            $prevEAP = $ErrorActionPreference
            $ErrorActionPreference = 'Continue'
            try {
                $installOutput = & go install github.com/tc-hib/go-winres@v0.3.3 2>&1
                $installExit = $LASTEXITCODE
            } finally {
                $ErrorActionPreference = $prevEAP
            }
            if ($installExit -ne 0) {
                Write-Warn "go install go-winres failed (exit $installExit); binary will have no icon"
                foreach ($line in $installOutput) {
                    $text = "$line".Trim()
                    if ($text.Length -gt 0) { Write-Host "  $text" -ForegroundColor Yellow }
                }
            } else {
                $goWinres = Get-Command go-winres -ErrorAction SilentlyContinue
            }
        }

        if ($goWinres) {
            Push-Location $GitMapDir
            try {
                $cleanVersion = "0.0.0.0"
                $constantsDir = Join-Path $GitMapDir "constants"
                $constantsFile = Join-Path $constantsDir "constants.go"
                if (Test-Path $constantsFile) {
                    $verMatch = Select-String -Path $constantsFile -Pattern 'const\s+Version\s*=\s*"([^"]+)"' | Select-Object -First 1
                    if ($verMatch) { $cleanVersion = ($verMatch.Matches[0].Groups[1].Value) -replace '^v', '' }
                }
                if ([string]::IsNullOrWhiteSpace($cleanVersion)) { $cleanVersion = "0.0.0.0" }

                $prevEAP = $ErrorActionPreference
                $ErrorActionPreference = 'Continue'
                try {
                    $winresOutput = & go-winres make --product-version $cleanVersion --file-version $cleanVersion 2>&1
                    $winresExit = $LASTEXITCODE
                } finally {
                    $ErrorActionPreference = $prevEAP
                }
                if ($winresExit -ne 0) {
                    Write-Warn "go-winres make failed (exit $winresExit); continuing without embedded icon"
                    foreach ($line in $winresOutput) {
                        $text = "$line".Trim()
                        if ($text.Length -gt 0) { Write-Host "  $text" -ForegroundColor Yellow }
                    }
                } else {
                    Write-Info "Generated rsrc_windows_*.syso (icon + manifest + version)"
                }
            } finally {
                Pop-Location
            }
        }
    }

    Write-Step "3/4" "Building $($Config.binaryName)"
    Test-SourceFiles

    $binDir  = Join-Path $RepoRoot $Config.buildOutput
    $outPath = Join-Path $binDir $Config.binaryName

    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
        Write-Info "Created bin directory"
    }

    Push-Location $GitMapDir
    try {
        $absRepoRoot = (Resolve-Path $RepoRoot).Path

        # Build-time identity injection — stamps the *source* repo metadata
        # into the binary so the `gitmap binary` footer block never falls
        # back to probing the user's CWD (fixed in v5.60.0).
        Push-Location $absRepoRoot
        try {
            $buildCommit = (& git rev-parse HEAD 2>$null) -as [string]
            $buildBranch = (& git rev-parse --abbrev-ref HEAD 2>$null) -as [string]
            $buildRepo   = (& git config --get remote.origin.url 2>$null) -as [string]
        } finally { Pop-Location }
        if ($null -eq $buildCommit) { $buildCommit = "" }
        if ($null -eq $buildBranch) { $buildBranch = "" }
        if ($null -eq $buildRepo)   { $buildRepo   = "" }
        $buildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

        $ldflags = "-X 'github.com/alimtvnetwork/gitmap-v27/gitmap/constants.RepoPath=$absRepoRoot'" +
                   " -X 'github.com/alimtvnetwork/gitmap-v27/gitmap/cmd.BuildCommit=$buildCommit'" +
                   " -X 'github.com/alimtvnetwork/gitmap-v27/gitmap/cmd.BuildBranch=$buildBranch'" +
                   " -X 'github.com/alimtvnetwork/gitmap-v27/gitmap/cmd.BuildRepo=$buildRepo'" +
                   " -X 'github.com/alimtvnetwork/gitmap-v27/gitmap/cmd.BuildDate=$buildDate'"

        # Pre-build provenance stamp — prints commit SHA, branch, declared
        # version, and a fingerprint of the historically-problematic cmd/
        # files so a stale checkout is obvious in the build log before
        # `go build` runs. Non-fatal: stamp failures never block the build.
        $stampScript = Join-Path $RepoRoot 'scripts\build-stamp.ps1'
        if (Test-Path $stampScript) {
            try { & $stampScript } catch { Write-Warn "build-stamp failed: $_" }
        }

        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $buildOutput = go build -ldflags $ldflags -o $outPath . 2>&1
        $buildExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($buildExit -ne 0) {
            Write-Fail "Go build failed"
            foreach ($line in $buildOutput) {
                $text = "$line".Trim()
                if ($text.Length -gt 0) {
                    Write-Host "  $text" -ForegroundColor Red
                }
            }
            exit 1
        }
    } finally {
        Pop-Location
    }

    if ($Config.copyData) {
        Copy-DataFolder -BinDir $binDir
    }

    $size = (Get-Item $outPath).Length / 1MB
    Write-Success ("Binary built ({0:N2} MB) -> $outPath" -f $size)

    return $outPath
}

# -- Copy data folder -----------------------------------------
function Copy-DataFolder {
    param($BinDir)

    $dataSource = Join-Path $GitMapDir "data"
    $dataDest   = Join-Path $BinDir "data"

    if (Test-Path $dataSource) {
        if (Test-Path $dataDest) {
            Remove-Item $dataDest -Recurse -Force
        }
        Copy-Item $dataSource $dataDest -Recurse
        Write-Info "Copied data folder to bin"
    }
}

function Sync-DeployDataFolder {
    param(
        [string]$SourceDir,
        [string]$DestDir
    )

    $hasSourceDir = Test-Path $SourceDir
    if ($hasSourceDir -eq $false) {
        return
    }

    $hasDestDir = Test-Path $DestDir
    if ($hasDestDir -eq $false) {
        Copy-Item $SourceDir $DestDir -Recurse
        Write-Info "Copied data folder to gitmap app directory"
        return
    }

    $files = Get-ChildItem -Path $SourceDir -File
    foreach ($file in $files) {
        $targetFile = Join-Path $DestDir $file.Name
        Copy-Item $file.FullName $targetFile -Force
    }

    Write-Info "Synced static data files to gitmap app directory"
}

# -- Copy docs-site to deploy directory -----------------------
# Required for `gitmap help-dashboard` (hd) which resolves docs-site/
# relative to the binary directory. Without this, `gitmap hd` fails with:
#   "Docs site directory not found at <deploy>/docs-site"
#
# Source resolution order (first hit wins):
#   1. <repo>/docs-site/dist/   — legacy layout with a dedicated subdir
#   2. <repo>/docs-site/        — legacy layout, source only (no prebuilt dist)
#   3. <repo>/dist/             — current layout where the repo root IS the
#                                 Vite docs app (no docs-site/ subdir)
#   4. Auto-build at repo root  — if package.json has a `build` script and
#                                 npm is on PATH, run it and use <repo>/dist/.
#   5. Warn — `gitmap hd` will fail until docs are built.
function Copy-DocsSite {
    param($AppDir)

    $docsDest    = Join-Path $AppDir "docs-site"
    $legacyDir   = Join-Path $RepoRoot "docs-site"
    $legacyDist  = Join-Path $legacyDir "dist"
    $rootDist    = Join-Path $RepoRoot "dist"
    $rootPkg     = Join-Path $RepoRoot "package.json"
    $gitmapMain  = Join-Path $GitMapDir "main.go"
    $nodeModules = Join-Path $RepoRoot "node_modules"

    # Repo-detect diagnostics (active under -DebugRepoDetect or env var).
    Write-RepoDetect -Check "RepoRoot"          -Result $RepoRoot
    Write-RepoDetect -Check "GitMapDir"         -Result $GitMapDir
    Write-RepoDetect -Check "gitmap/main.go"    -Result $(if (Test-Path $gitmapMain) { "present" } else { "missing" }) -Detail $gitmapMain
    Write-RepoDetect -Check "package.json"      -Result $(if (Test-Path $rootPkg) { "present" } else { "missing" })   -Detail $rootPkg
    Write-RepoDetect -Check "node_modules/"     -Result $(if (Test-Path $nodeModules) { "present" } else { "missing" })
    Write-RepoDetect -Check "docs-site/dist/"   -Result $(if (Test-Path $legacyDist) { "present" } else { "missing" })
    Write-RepoDetect -Check "dist/ (root)"      -Result $(if (Test-Path $rootDist)  { "present" } else { "missing" })
    $npmCmd = Get-Command npm -ErrorAction SilentlyContinue
    Write-RepoDetect -Check "npm on PATH"       -Result $(if ($npmCmd) { "yes" } else { "no" }) -Detail $(if ($npmCmd) { $npmCmd.Source } else { "" })
    Write-RepoDetectSnippet -Title "package.json (first 6 lines)" -Path $rootPkg

    # 1. Legacy <repo>/docs-site/dist/
    if (Test-Path $legacyDist) {
        Write-RepoDetect -Check "decision" -Result "use-prebuilt-legacy" -Detail $legacyDist
        $distDest = Join-Path $docsDest "dist"
        if (Test-Path $distDest) { Remove-Item $distDest -Recurse -Force }
        New-Item -ItemType Directory -Path $docsDest -Force | Out-Null
        Copy-Item $legacyDist $distDest -Recurse
        Write-Info "Copied docs-site/dist to gitmap app directory"
        return
    }

    # 3. Current <repo>/dist/ (root-level Vite app)
    if (Test-Path $rootDist) {
        Write-RepoDetect -Check "decision" -Result "use-prebuilt-root" -Detail $rootDist
        $distDest = Join-Path $docsDest "dist"
        if (Test-Path $distDest) { Remove-Item $distDest -Recurse -Force }
        New-Item -ItemType Directory -Path $docsDest -Force | Out-Null
        Copy-Item $rootDist $distDest -Recurse
        Write-Info "Copied root dist/ to gitmap app docs-site/dist"
        return
    }

    # 4. Auto-build the root Vite app if package.json + npm available
    if ((Test-Path $rootPkg) -and $npmCmd) {
        $pkgRaw = Get-Content $rootPkg -Raw
        $hasBuild = ($pkgRaw -match '"build"\s*:')
        $hasVite  = ($pkgRaw -match '"vite"\s*:')
        Write-RepoDetect -Check "package.json:build" -Result $(if ($hasBuild) { "found" } else { "missing" })
        Write-RepoDetect -Check "package.json:vite"  -Result $(if ($hasVite)  { "found" } else { "missing" })
        if ($hasBuild) {
            Write-RepoDetect -Check "decision" -Result "auto-build" -Detail "npm run build at $RepoRoot"
            Push-Location $RepoRoot
            try {
                $nodeModules = Join-Path $RepoRoot "node_modules"
                $viteBin = Join-Path $nodeModules ".bin\vite.cmd"
                if (-not (Test-Path $nodeModules) -or -not (Test-Path $viteBin)) {
                    Write-Info "Installing docs dependencies (npm install) at repo root..."
                    $installExit = [int](Invoke-NpmQuiet -NpmArgs @('install','--no-audit','--no-fund','--silent'))
                    if ($installExit -ne 0) {
                        Write-Warn "npm install failed - skipping docs build"
                        Write-ReportError -Stage "docs-npm-install" `
                            -Command "npm install --no-audit --no-fund --silent" `
                            -ExitCode $installExit `
                            -Message "npm install failed at repo root; docs build skipped" `
                            -Paths @{ repoRoot = $RepoRoot; packageJson = $rootPkg }
                        Pop-Location
                        return
                    }
                }
                Write-Info "Auto-building docs (npm run build) at repo root..."
                $buildExit = [int](Invoke-NpmQuiet -NpmArgs @('run','build'))
                if ($buildExit -eq 0 -and (Test-Path $rootDist)) {
                    $distDest = Join-Path $docsDest "dist"
                    if (Test-Path $distDest) { Remove-Item $distDest -Recurse -Force }
                    New-Item -ItemType Directory -Path $docsDest -Force | Out-Null
                    Copy-Item $rootDist $distDest -Recurse
                    Write-Info "Built and copied docs to gitmap app docs-site/dist"
                    return
                }
            } finally {
                Pop-Location
            }
            Write-Warn "Auto-build failed - 'gitmap hd' will fail"
            Write-ReportError -Stage "docs-npm-build" `
                -Command "npm run build" `
                -ExitCode ([int]$buildExit) `
                -Message "npm run build did not produce dist/ output" `
                -Paths @{ repoRoot = $RepoRoot; expectedDist = $rootDist; packageJson = $rootPkg }
            return
        } else {
            Write-RepoDetect -Check "decision" -Result "skip-no-build-script" -Detail "package.json has no `"build`" entry"
        }
    } else {
        $skipReason = if (-not (Test-Path $rootPkg)) { "no package.json" } elseif (-not $npmCmd) { "npm not on PATH" } else { "unknown" }
        Write-RepoDetect -Check "decision" -Result "skip-not-a-vite-repo" -Detail $skipReason
    }

    # 2. Legacy <repo>/docs-site/ source-only (npm-dev fallback)
    if (Test-Path $legacyDir) {
        # fall through to existing source-copy block below
    } else {
        Write-RepoDetect -Check "decision" -Result "no-docs-source"
        Write-Warn "No docs found (checked docs-site/dist, docs-site/, dist/) - 'gitmap hd' will fail"
        return
    }

    # 2. Legacy <repo>/docs-site/ source-only — npm-dev fallback (no prebuilt dist).
    # Copy everything except node_modules to keep the deploy lean.
    Write-RepoDetect -Check "decision" -Result "use-legacy-source" -Detail $legacyDir
    if (Test-Path $docsDest) {
        Remove-Item $docsDest -Recurse -Force
    }
    New-Item -ItemType Directory -Path $docsDest -Force | Out-Null
    Get-ChildItem -Path $legacyDir -Force | Where-Object { $_.Name -ne "node_modules" } | ForEach-Object {
        Copy-Item $_.FullName -Destination $docsDest -Recurse -Force
    }
    Write-Warn "No prebuilt dist/ found - copied docs-site/ source only (run 'npm run build' for static mode)"
}

# -- Resolve deploy target -------------------------------------
# Priority: 1) -DeployPath flag  2) globally installed gitmap location  3) powershell.json default
function Resolve-DeployTarget {
    param($Config, $OverridePath)

    # 1) Explicit CLI override always wins
    if ($OverridePath.Length -gt 0) {
        Write-Info "Deploy target: CLI override -> $OverridePath"

        return $OverridePath
    }

    # 2) If gitmap is already on PATH, deploy to its parent directory
    $activeCmd = Get-Command gitmap -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($activeCmd) {
        $activePath = $activeCmd.Source
        if ($activePath -and (Test-Path $activePath)) {
            $resolvedActive = (Resolve-Path $activePath).Path
            $activeDir = Split-Path $resolvedActive -Parent
            $activeDirName = Split-Path $activeDir -Leaf

            # The binary lives in <deploy-target>/$AppSubdir/gitmap.exe (or
            # any legacy folder name in $LegacyAppSubdirs from the manifest).
            # Either way the deploy target is the parent of that subfolder.
            if (Test-KnownAppSubdir $activeDirName) {
                $deployTarget = Split-Path $activeDir -Parent
                Write-Info "Deploy target: detected from PATH -> $deployTarget"

                return $deployTarget
            }

            # Binary is directly in a folder (not nested under gitmap/)
            # Deploy target = that folder's parent so we create gitmap/ there
            $deployTarget = Split-Path $activeDir -Parent
            Write-Info "Deploy target: detected from PATH -> $deployTarget"

            return $deployTarget
        }
    }

    # 3) Fall back to powershell.json default
    Write-Info "Deploy target: powershell.json default -> $($Config.deployPath)"

    return $Config.deployPath
}

# -- Deploy to target directory --------------------------------
function Deploy-Binary {
    param($Config, $BinaryPath, $OverridePath)

    Write-Step "4/4" "Deploying"

    $target = Resolve-DeployTarget -Config $Config -OverridePath $OverridePath

    Write-Info "Target: $target"

    if (-not (Test-Path $target)) {
        New-Item -ItemType Directory -Path $target -Force | Out-Null
        Write-Info "Created deploy directory"
    }

    # Migrate any legacy unwrapped layout (DFD-3) BEFORE we resolve $appDir.
    Repair-DeployLayout -DeployTarget $target -BinaryName $Config.binaryName

    # Deploy into nested $AppSubdir/ subfolder (DFD-1). Folder name comes
    # from gitmap/constants/deploy-manifest.json (single source of truth).
    $appDir = Join-Path $target $script:AppSubdir
    if (-not (Test-Path $appDir)) {
        New-Item -ItemType Directory -Path $appDir -Force | Out-Null
        Write-Info "Created $($script:AppSubdir) app directory"
    }

    # Pre-deploy cleanup (DFD-6) — runs BEFORE the new binary is copied
    # so a locked .old file can't block the deploy.
    Invoke-DeployCleanup -DeployTarget $target -AppDir $appDir -BinaryName $Config.binaryName

    $destFile = Join-Path $appDir $Config.binaryName
    $backupFile = "$destFile.old"
    $hasBackup = $false
    $deploySuccess = $false

    if (Test-Path $destFile) {
        # Rename-first strategy: Windows allows renaming a running binary
        # but not overwriting it. Rename to .old, then copy the new one.
        try {
            if (Test-Path $backupFile) {
                Remove-Item $backupFile -Force -ErrorAction SilentlyContinue
            }
            Rename-Item $destFile $backupFile -Force -ErrorAction Stop
            $hasBackup = $true
            Write-Info "Renamed existing binary to $($Config.binaryName).old (rename-first)"
        } catch {
            Write-Warn "Rename-first failed: $_"
            # Fallback: try a backup copy instead
            try {
                Copy-Item $destFile $backupFile -Force -ErrorAction Stop
                $hasBackup = $true
                Write-Info "Backed up existing binary to $($Config.binaryName).old"
            } catch {
                Write-Warn "Could not create backup: $_"
            }
        }
    }

    # Copy new binary — after rename-first, the destination is free
    $maxAttempts = 5
    $attempt = 1
    while ($true) {
        try {
            Copy-Item $BinaryPath $destFile -Force -ErrorAction Stop
            $deploySuccess = $true
            break
        } catch {
            if ($attempt -ge $maxAttempts) {
                # Restore backup on failure
                if ($hasBackup -and (Test-Path $backupFile) -and (-not (Test-Path $destFile))) {
                    Write-Warn "Deploy failed - restoring previous binary from backup"
                    try {
                        Rename-Item $backupFile $destFile -Force -ErrorAction Stop
                        Write-Success "Rollback complete - previous version restored"
                    } catch {
                        Write-Fail "Rollback also failed: $_"
                    }
                }
                throw
            }
            Write-Warn "Target still locked; retrying ($attempt/$maxAttempts)..."
            Start-Sleep -Milliseconds 500
            $attempt++
        }
    }

    # Post-deploy: remove the .old immediately now that the new binary is in place.
    if ($hasBackup -and $deploySuccess -and (Test-Path $backupFile)) {
        try {
            Remove-Item $backupFile -Force -ErrorAction Stop
            Write-Info "Removed .old backup ($($Config.binaryName).old)"
        } catch {
            Write-Warn "Could not remove .old backup: $_ (will be cleaned on next deploy)"
        }
    }

    $binDir   = Split-Path $BinaryPath -Parent
    $dataDir  = Join-Path $binDir "data"
    $dataDest = Join-Path $appDir "data"
    Sync-DeployDataFolder -SourceDir $dataDir -DestDir $dataDest

    Copy-DocsSite -AppDir $appDir

    Write-Success "Deployed to $appDir"

    # Register the app folder on user PATH and refresh the current session (DFD-4, DFD-5).
    Register-OnPath -AppDir $appDir

    # Sync source repo path in DB so "gitmap update" uses this repo location
    $syncBinary = $destFile
    if (-not (Test-Path $syncBinary)) { $syncBinary = $BinaryPath }
    if (Test-Path $syncBinary) {
        try {
            & $syncBinary set-source-repo $RepoRoot 2>&1 | Out-Null
            Write-Info "Source repo path synced to DB: $RepoRoot"
        } catch {
            Write-Warn "Could not sync source repo path: $_"
        }
    }
}

# -- Repair legacy unwrapped layout (DFD-3) --------------------
# Two migrations happen here, in priority order. Both are idempotent.
#
#   1) v3.6.0 rename: <DeployTarget>\gitmap\ -> <DeployTarget>\gitmap-cli\
#      Triggered when an old install lives in the legacy "gitmap" subfolder
#      and the new "gitmap-cli" folder does not yet exist (or is empty).
#      The whole folder (binary + data + docs) is moved, then the empty
#      legacy folder is removed.
#
#   2) Pre-v3.6.0 unwrapped layout: <DeployTarget>\gitmap.exe (top-level)
#      gets nested under <DeployTarget>\gitmap-cli\ for DFD-1 compatibility.
function Repair-DeployLayout {
    param(
        [string]$DeployTarget,
        [string]$BinaryName
    )

    $appDir = Join-Path $DeployTarget $script:AppSubdir
    $newBinary = Join-Path $appDir $BinaryName

    # Migration 1: rename any legacy app folder to $AppSubdir.
    foreach ($legacy in $script:LegacyAppSubdirs) {
        $legacySubdir = Join-Path $DeployTarget $legacy
        if ($legacySubdir -eq $appDir) { continue }
        $legacySubBinary = Join-Path $legacySubdir $BinaryName

        if ((Test-Path $legacySubBinary) -and (-not (Test-Path $newBinary))) {
            Write-Info "Layout: migrating legacy '$legacySubdir' -> '$appDir'"
            try {
                if (-not (Test-Path $appDir)) {
                    Move-Item -Path $legacySubdir -Destination $appDir -Force -ErrorAction Stop
                    Write-Info "Layout: rename complete"
                }
            } catch {
                Write-Warn "Layout: rename failed ($_); leaving legacy folder in place"
            }
        } elseif ((Test-Path $legacySubdir) -and (Test-Path $newBinary)) {
            # Both exist after a previous half-migration — drop the empty legacy folder.
            try {
                $remaining = Get-ChildItem -Path $legacySubdir -Force -ErrorAction SilentlyContinue
                if (-not $remaining -or $remaining.Count -eq 0) {
                    Remove-Item $legacySubdir -Force -Recurse -ErrorAction Stop
                    Write-Info "Layout: removed empty legacy folder $legacySubdir"
                }
            } catch {
                Write-Warn "Layout: could not remove legacy folder $legacySubdir : $_"
            }
        }
    }

    # Migration 2: top-level unwrapped binary -> gitmap-cli\.
    $legacyBinary = Join-Path $DeployTarget $BinaryName
    $wrappedBinary = Join-Path $appDir $BinaryName

    if (-not (Test-Path $legacyBinary)) {
        Write-Info "Layout: OK (no legacy binary at $DeployTarget)"
        return
    }
    if (Test-Path $wrappedBinary) {
        # Both exist — wrapped wins; the legacy copy is leftover. Remove it.
        try {
            Remove-Item $legacyBinary -Force -ErrorAction Stop
            Write-Info "Layout: removed leftover legacy binary at $legacyBinary"
        } catch {
            Write-Warn "Layout: could not remove legacy binary $legacyBinary : $_"
        }
        return
    }

    Write-Info "Layout: migrating legacy unwrapped install -> $appDir"
    if (-not (Test-Path $appDir)) {
        New-Item -ItemType Directory -Path $appDir -Force | Out-Null
    }

    foreach ($name in @($BinaryName, "data", "CHANGELOG.md", "docs")) {
        $src = Join-Path $DeployTarget $name
        $dst = Join-Path $appDir $name
        if (-not (Test-Path $src)) { continue }
        if (Test-Path $dst) {
            Write-Info "Layout: $name already inside gitmap/, skipping move"
            continue
        }
        try {
            Move-Item -Path $src -Destination $dst -Force -ErrorAction Stop
            Write-Info "Layout: moved $name -> gitmap\$name"
        } catch {
            Write-Warn "Layout: could not move $name : $_"
        }
    }
}

# -- Pre-deploy cleanup (DFD-6, DFD-7) -------------------------
# Removes prior-deploy artifacts before the new binary is copied:
#   *.old, *-update-*.exe, updater-tmp-*.exe, temp *-update-*.ps1,
#   *.gitmap-tmp-* swap dirs, and the obsolete drive-root shim.
function Invoke-DeployCleanup {
    param(
        [string]$DeployTarget,
        [string]$AppDir,
        [string]$BinaryName
    )

    $removed = 0
    $scanDirs = @($DeployTarget, $AppDir) | Where-Object { Test-Path $_ } | Select-Object -Unique
    $patterns = @("*.old", "$($BinaryName.Replace('.exe',''))-update-*.exe",
                  "$($BinaryName.Replace('.exe',''))-update-*", "updater-tmp-*.exe")

    foreach ($dir in $scanDirs) {
        foreach ($pat in $patterns) {
            $matches = Get-ChildItem -Path $dir -Filter $pat -File -ErrorAction SilentlyContinue
            foreach ($f in $matches) {
                try {
                    Remove-Item $f.FullName -Force -ErrorAction Stop
                    Write-Info "[cleanup] removed $($f.FullName)"
                    $removed++
                } catch {
                    Write-Warn "[cleanup] could not remove $($f.FullName): $_"
                }
            }
        }
    }

    # Temp dir scripts
    $tempScripts = Get-ChildItem -Path $env:TEMP -Filter "$($BinaryName.Replace('.exe',''))-update-*.ps1" -File -ErrorAction SilentlyContinue
    foreach ($f in $tempScripts) {
        try {
            Remove-Item $f.FullName -Force -ErrorAction Stop
            Write-Info "[cleanup] removed temp script $($f.FullName)"
            $removed++
        } catch {
            Write-Warn "[cleanup] could not remove $($f.FullName): $_"
        }
    }

    # *.gitmap-tmp-* swap directories left by interrupted clones
    $tmpParents = @($DeployTarget) | Where-Object { Test-Path $_ }
    foreach ($parent in $tmpParents) {
        $swaps = Get-ChildItem -Path $parent -Directory -Filter "*.gitmap-tmp-*" -ErrorAction SilentlyContinue
        foreach ($d in $swaps) {
            try {
                Remove-Item $d.FullName -Recurse -Force -ErrorAction Stop
                Write-Info "[cleanup] removed swap dir $($d.FullName)"
                $removed++
            } catch {
                Write-Warn "[cleanup] could not remove $($d.FullName): $_"
            }
        }
    }

    # Drive-root shim from v2.90.0 (DFD-7)
    Remove-DriveRootShim -DeployTarget $DeployTarget -BinaryName $BinaryName | ForEach-Object { $removed += $_ }

    if ($removed -gt 0) {
        Write-Success "[cleanup] removed $removed artifact(s)"
    } else {
        Write-Info "[cleanup] nothing to clean"
    }
}

# -- Remove obsolete drive-root shim (DFD-7) -------------------
# In v2.90.0 we wrote <drive>:\gitmap.exe as a forwarding shim. That
# pattern is now removed. Detect and delete it if present, but never
# touch a binary that lives inside a gitmap\ folder.
function Remove-DriveRootShim {
    param(
        [string]$DeployTarget,
        [string]$BinaryName
    )

    $drive = [System.IO.Path]::GetPathRoot($DeployTarget)
    if ([string]::IsNullOrWhiteSpace($drive) -or $drive -notmatch '^[A-Za-z]:\\$') {
        return 0
    }

    $shimPath = Join-Path $drive $BinaryName
    if (-not (Test-Path $shimPath)) { return 0 }

    $shimDir = Split-Path $shimPath -Parent
    # Safety: only remove if it sits at the literal drive root and not inside an app subdir.
    $shimDirName = Split-Path $shimDir -Leaf
    if (Test-KnownAppSubdir $shimDirName) {
        return 0
    }

    $size = (Get-Item $shimPath).Length
    if ($size -gt 5MB) {
        Write-Warn "[cleanup] skipping drive-root $shimPath (size $size > 5MB; likely unrelated)"
        return 0
    }

    try {
        Remove-Item $shimPath -Force -ErrorAction Stop
        Write-Info "[cleanup] removed obsolete drive-root shim $shimPath"
        return 1
    } catch {
        Write-Warn "[cleanup] could not remove drive-root shim $shimPath : $_"
        return 0
    }
}

# -- Register on user PATH + refresh current session (DFD-4) ---
function Register-OnPath {
    param([string]$AppDir)

    if (-not (Test-Path $AppDir)) {
        Write-Warn "PATH: skipping (app dir does not exist: $AppDir)"
        return
    }

    $resolved = (Resolve-Path $AppDir).Path.TrimEnd('\')

    # Refresh current session unconditionally so the binary is callable now.
    $currentEntries = ($env:Path -split ';') | Where-Object { $_ -and ($_.TrimEnd('\') -ieq $resolved) }
    if (-not $currentEntries) {
        $env:Path = "$env:Path;$resolved"
        Write-Info "PATH: appended to current session -> $resolved"
    } else {
        Write-Info "PATH: already in current session"
    }

    # Persist to user-scope PATH (idempotent).
    try {
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        if (-not $userPath) { $userPath = "" }
        $userEntries = ($userPath -split ';') | Where-Object { $_ -and ($_.TrimEnd('\') -ieq $resolved) }
        if ($userEntries) {
            Write-Info "PATH: already in user PATH"
            return
        }
        $newUserPath = if ($userPath.Length -gt 0) { "$userPath;$resolved" } else { $resolved }
        [Environment]::SetEnvironmentVariable('Path', $newUserPath, 'User')
        Write-Success "PATH: persisted to user PATH -> $resolved"
        Write-Info "PATH: open a NEW shell, or run '. `$PROFILE' in PowerShell, to pick it up everywhere"
    } catch {
        Write-Warn "PATH: could not persist to user PATH: $_"
    }
}

# -- Remove a directory entry from user PATH (DFD-8) -----------
function Remove-FromUserPath {
    param([string]$DirToRemove)

    if ([string]::IsNullOrWhiteSpace($DirToRemove)) { return }
    $normalized = $DirToRemove.TrimEnd('\')

    $sessionParts = ($env:Path -split ';') | Where-Object { $_ -and ($_.TrimEnd('\') -ine $normalized) }
    $env:Path = ($sessionParts -join ';')

    try {
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        if (-not $userPath) { return }
        $userParts = ($userPath -split ';') | Where-Object { $_ -and ($_.TrimEnd('\') -ine $normalized) }
        $newUserPath = ($userParts -join ';')
        if ($newUserPath -ne $userPath) {
            [Environment]::SetEnvironmentVariable('Path', $newUserPath, 'User')
            Write-Success "PATH: removed stale entry -> $normalized"
        }
    } catch {
        Write-Warn "PATH: could not edit user PATH: $_"
    }
}

# -- Migrate a stale active gitmap binary off-target (DFD-8) ---
# When the binary on PATH lives outside the resolved deploy target, the
# old behavior was to copy the new build *into* the stale location,
# preserving the wrong path forever. Instead: delete the stale binary,
# remove its parent dir if it is now empty, and strip its directory
# from the user PATH.
function Migrate-StaleActiveBinary {
    param(
        [string]$StaleBinaryPath,
        [string]$DeployedAppDir,
        [string]$BinaryName
    )

    if (-not (Test-Path $StaleBinaryPath)) { return }

    $staleDir = Split-Path $StaleBinaryPath -Parent
    $staleResolved = (Resolve-Path $staleDir).Path.TrimEnd('\')
    $deployResolved = (Resolve-Path $DeployedAppDir).Path.TrimEnd('\')
    if ($staleResolved -ieq $deployResolved) { return }

    Write-Warn "PATH: stale active binary detected -> $StaleBinaryPath"
    Write-Info "PATH: migrating away from stale location"

    foreach ($pat in @($BinaryName, "$BinaryName.old", "$($BinaryName.Replace('.exe',''))-update-*.exe")) {
        $hits = Get-ChildItem -Path $staleDir -Filter $pat -File -ErrorAction SilentlyContinue
        foreach ($f in $hits) {
            try {
                Remove-Item $f.FullName -Force -ErrorAction Stop
                Write-Info "[cleanup] removed stale $($f.FullName)"
            } catch {
                Write-Warn "[cleanup] could not remove $($f.FullName): $_"
            }
        }
    }

    # Walk upward removing now-empty gitmap-owned dirs.
    $cursor = $staleDir
    for ($i = 0; $i -lt 3; $i++) {
        if (-not (Test-Path $cursor)) { break }
        $remaining = Get-ChildItem -Path $cursor -Force -ErrorAction SilentlyContinue
        if ($remaining -and $remaining.Count -gt 0) { break }
        try {
            Remove-Item $cursor -Force -Recurse -ErrorAction Stop
            Write-Info "[cleanup] removed empty stale dir $cursor"
        } catch {
            Write-Warn "[cleanup] could not remove $cursor : $_"
            break
        }
        $cursor = Split-Path $cursor -Parent
    }

    Remove-FromUserPath -DirToRemove $staleResolved
    $staleParent = Split-Path $staleResolved -Parent
    if ($staleParent -and ($staleParent.TrimEnd('\') -ne ([System.IO.Path]::GetPathRoot($staleParent).TrimEnd('\')))) {
        Remove-FromUserPath -DirToRemove $staleParent
    }
}

# -- Persist resolved deploy target back to powershell.json (DFD-9)
function Sync-ConfigDeployPath {
    param([string]$EffectiveDeployTarget)

    $configPath = Join-Path $GitMapDir "powershell.json"
    if (-not (Test-Path $configPath)) { return }

    try {
        $raw = Get-Content $configPath -Raw
        $cfg = $raw | ConvertFrom-Json
        $existing = "$($cfg.deployPath)".TrimEnd('\')
        $resolved = $EffectiveDeployTarget.TrimEnd('\')
        if ($existing -ieq $resolved) { return }

        $cfg.deployPath = $resolved
        ($cfg | ConvertTo-Json -Depth 10) | Set-Content -Path $configPath -Encoding UTF8
        Write-Success "Config: powershell.json deployPath updated -> $resolved"
    } catch {
        Write-Warn "Config: could not update powershell.json: $_"
    }
}


# -- Run gitmap ------------------------------------------------
function Invoke-Run {
    param($Config, $BinaryPath, [string[]]$CliArgs)

    Write-Host ""
    Write-Step "RUN" "Executing gitmap"

    # Always run from the local bin build, never from the deploy target
    $binDir = Split-Path $BinaryPath -Parent
    $dataDir = Join-Path $binDir "data"

    $resolvedArgs = Resolve-RunArgs -CliArgs $CliArgs
    $argString = $resolvedArgs -join ' '
    $currentDir = (Get-Location).Path
    Write-Info "Binary: $BinaryPath"
    Write-Info "Runner CWD: $currentDir"
    Write-Info "Command: gitmap $argString"
    if ($resolvedArgs.Count -ge 2 -and $resolvedArgs[0] -eq "scan") {
        Write-Info "Scan target: $($resolvedArgs[1])"
    }
    Write-Host ("  " + ("-" * 50)) -ForegroundColor DarkGray
    Write-Host ""

    $proc = Start-Process -FilePath $BinaryPath -ArgumentList $resolvedArgs -WorkingDirectory $binDir -NoNewWindow -Wait -PassThru

    Write-Host ""
    if ($proc.ExitCode -eq 0) {
        Write-Success "Run complete"
    } else {
        Write-Fail "gitmap exited with code $($proc.ExitCode)"
    }
}

# -- Resolve run arguments -------------------------------------
function Resolve-RunArgs {
    param([string[]]$CliArgs)

    if ($CliArgs.Count -eq 0) {
        $parentDir = Split-Path $RepoRoot -Parent
        Write-Info "No args provided, defaulting to: scan $parentDir"

        return @("scan", $parentDir)
    }

    # Resolve relative paths to absolute so Start-Process always receives correct targets
    $baseDir = (Get-Location).Path
    $resolved = @()
    foreach ($arg in $CliArgs) {
        if ($arg -match '^(\.\.[\\/]|\.[\\/]|\.\.?$)' -and -not $arg.StartsWith('-')) {
            $path = Resolve-Path -LiteralPath $arg -ErrorAction SilentlyContinue
            if ($path) {
                $resolved += $path.Path
            } else {
                $resolved += [System.IO.Path]::GetFullPath((Join-Path $baseDir $arg))
            }
        } else {
            $resolved += $arg
        }
    }

    return $resolved
}

# -- Run tests -------------------------------------------------
function Invoke-Tests {
    Write-Step "TEST" "Running unit tests"

    $reportDir = Join-Path (Join-Path $GitMapDir "data") "unit-test-reports"
    if (-not (Test-Path $reportDir)) {
        New-Item -ItemType Directory -Path $reportDir -Force | Out-Null
        Write-Info "Created report directory: $reportDir"
    }

    $overallLog = Join-Path $reportDir "overall.log.txt"
    $failingLog = Join-Path $reportDir "failingTest.log.txt"

    Push-Location $GitMapDir
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"

        Write-Info "Running: go test ./..."
        $testOutput = go test ./... -v -count=1 2>&1
        $testExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        # Write overall report
        $testOutput | Out-File -FilePath $overallLog -Encoding UTF8
        Write-Info "Overall report: $overallLog"

        # Extract failing tests
        $failLines = @()
        $currentTest = ""
        $inFail = $false
        foreach ($line in $testOutput) {
            $text = "$line"
            if ($text -match "^--- FAIL:") {
                $inFail = $true
                $currentTest = $text
                $failLines += ""
                $failLines += $text
            } elseif ($text -match "^--- PASS:" -or $text -match "^=== RUN") {
                $inFail = $false
            } elseif ($text -match "^FAIL\s") {
                $failLines += $text
            } elseif ($inFail) {
                $failLines += $text
            }
        }

        if ($failLines.Count -gt 0) {
            $failLines | Out-File -FilePath $failingLog -Encoding UTF8
            Write-Fail "Some tests failed. See: $failingLog"
        } else {
            "No failing tests." | Out-File -FilePath $failingLog -Encoding UTF8
            Write-Success "All tests passed"
        }

        # Print summary
        $passCount = ($testOutput | Where-Object { "$_" -match "^--- PASS:" }).Count
        $failCount = ($testOutput | Where-Object { "$_" -match "^--- FAIL:" }).Count
        $skipCount = ($testOutput | Where-Object { "$_" -match "^--- SKIP:" }).Count
        Write-Info "Results: $passCount passed, $failCount failed, $skipCount skipped"

        # Show test output in terminal
        foreach ($line in $testOutput) {
            $text = "$line".Trim()
            if ($text -match "^--- FAIL:") {
                Write-Host "  $text" -ForegroundColor Red
            } elseif ($text -match "^--- PASS:") {
                Write-Host "  $text" -ForegroundColor Green
            } elseif ($text -match "^FAIL") {
                Write-Host "  $text" -ForegroundColor Red
            } elseif ($text -match "^ok\s") {
                Write-Host "  $text" -ForegroundColor Green
            } elseif ($text.Length -gt 0) {
                Write-Host "  $text" -ForegroundColor Gray
            }
        }

        if ($testExit -ne 0) {
            Write-Fail "Tests failed (exit code $testExit)"
        }
    } finally {
        Pop-Location
    }
}

# -- Main ------------------------------------------------------
Show-Banner
Get-DeployManifest
$config = Load-Config
Test-GoToolchain

# -- Uninstall / Reinstall handlers ----------------------------
# -Uninstall  : delegate to ./uninstall-quick.ps1 -Yes and exit.
# -Reinstall  : run the uninstall, then re-invoke this very script with
#               NO arguments so the user gets a clean pull/build/deploy/setup.
# Both flags short-circuit before pull/build so they never touch git.
function Invoke-UninstallScript {
    $uninstallScript = Join-Path $RepoRoot "uninstall-quick.ps1"
    if (-not (Test-Path $uninstallScript)) {
        Write-Fail "uninstall-quick.ps1 not found at $uninstallScript"
        exit 1
    }
    Write-Step "uninstall" "Running uninstall-quick.ps1 -Yes"
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & $uninstallScript -Yes
    $uninstallExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref
    if ($uninstallExit -ne 0) {
        Write-Fail "uninstall-quick.ps1 exited with code $uninstallExit"
        exit $uninstallExit
    }
    Write-Success "Uninstall complete"
}

if ($Uninstall) {
    Invoke-UninstallScript
    Write-Host ""
    Write-Success "All done!"
    Write-Host ""
    exit 0
}

if ($Reinstall) {
    Invoke-UninstallScript
    Write-Host ""
    Write-Step "reinstall" "Re-invoking run.ps1 with no arguments"
    $selfPath = $MyInvocation.MyCommand.Path
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & $selfPath
    $reinstallExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref
    exit $reinstallExit
}

if ($Test) {
    Write-Info "Test mode enabled (-t)"
    Resolve-Dependencies
    Invoke-Tests
    Write-Host ""
    Write-Success "All done!"
    Write-Host ""
    exit 0
}

if ($Update) {
    Write-Info "Update mode enabled (-Update)"
}

if (-not $NoPull) {
    Invoke-GitPull
} else {
    Write-Info "Skipping git pull (-NoPull)"
}

Resolve-Dependencies
$binaryPath = Build-Binary -Config $config

# Show built version
$prevEAP = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
try {
    $versionOutput = & $binaryPath version 2>$null
} finally {
    $ErrorActionPreference = $prevEAP
}
Write-Info "Version: $versionOutput"

$deployedBinaryPath = $null
if ($Deploy) { $NoDeploy = $false }
if (-not $NoDeploy) {
    Deploy-Binary -Config $config -BinaryPath $binaryPath -OverridePath $DeployPath

    $effectiveDeployPath = Resolve-DeployTarget -Config $config -OverridePath $DeployPath
    $deployedAppDir = Join-Path $effectiveDeployPath $script:AppSubdir
    $deployedBinaryPath = Join-Path $deployedAppDir $config.binaryName

    # Persist the resolved target so future runs (and the config-binary
    # readout below) reflect reality (DFD-9).
    Sync-ConfigDeployPath -EffectiveDeployTarget $effectiveDeployPath

    $activeCmd = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($activeCmd -and (Test-Path $deployedBinaryPath)) {
        $activeBinaryPath = $activeCmd.Source
        if (Test-Path $activeBinaryPath) {
            $activeResolved = (Resolve-Path $activeBinaryPath).Path
            $deployedResolved = (Resolve-Path $deployedBinaryPath).Path
            if ($activeResolved -ine $deployedResolved) {
                Write-Warn "PATH points to a different gitmap binary."
                Write-Info "Active:   $activeResolved"
                Write-Info "Deployed: $deployedResolved"

                # New behavior (DFD-8): do NOT copy the new build into the
                # stale location — that perpetuates the wrong path. Delete
                # the stale binary, prune empty parents, and strip the dir
                # from user PATH. The deployed dir is already on PATH via
                # Register-OnPath above.
                Migrate-StaleActiveBinary `
                    -StaleBinaryPath $activeBinaryPath `
                    -DeployedAppDir $deployedAppDir `
                    -BinaryName $config.binaryName

                Write-Info "PATH: open a NEW shell to pick up '$deployedAppDir'"
            }
        }
    }
} else {
    Write-Info "Skipping deploy (-NoDeploy)"
}

# -- Auto-run setup after deploy (DFD-10) ----------------------
# After a successful deploy, invoke `gitmap setup` so completion,
# cd-function, PATH snippet, and gitignore steps are applied without
# requiring a second manual command. Skip with -NoSetup.
if (-not $NoDeploy -and -not $NoSetup -and $deployedBinaryPath -and (Test-Path $deployedBinaryPath)) {
    Write-Host ""
    Write-Info "Running 'gitmap setup' automatically..."
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & $deployedBinaryPath setup
    $setupExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref
    if ($setupExit -ne 0) {
        Write-Warn "gitmap setup exited with code $setupExit"
    } else {
        Write-Success "Setup completed"
    }
} elseif ($NoSetup) {
    Write-Info "Skipping setup (-NoSetup)"
}

$changelogBinaryPath = $binaryPath
$activeCmdForChangelog = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
if ($activeCmdForChangelog -and (Test-Path $activeCmdForChangelog.Source)) {
    $changelogBinaryPath = $activeCmdForChangelog.Source
} elseif ($deployedBinaryPath -and (Test-Path $deployedBinaryPath)) {
    $changelogBinaryPath = $deployedBinaryPath
}

if (Test-Path $changelogBinaryPath) {
    Write-Host ""
    Write-Info "Latest changelog:"
    & $changelogBinaryPath changelog --latest

    if ($Update) {
        Write-Host ""
        Write-Info "Running update cleanup"
        & $changelogBinaryPath update-cleanup
    }
}

if ($R) {
    Invoke-Run -Config $config -BinaryPath $binaryPath -CliArgs $RunArgs
}

Write-Host ""
Write-Success "All done!"
Write-Host ""

# -- Last release info -----------------------------------------
$lastReleaseScript = Join-Path (Join-Path (Join-Path $RepoRoot "gitmap") "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    $lrBinary = $changelogBinaryPath
    & $lastReleaseScript -BinaryPath $lrBinary -RepoRoot $RepoRoot
    Write-Host ""
}
