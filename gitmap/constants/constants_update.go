package constants

// Update handoff file patterns.
const (
	UpdateCopyFmtExe  = "gitmap-update-%d.exe"
	UpdateCopyFmtUnix = "gitmap-update-%d"
	UpdateCopyGlob    = "gitmap-update-*"
	UpdateScriptGlob  = "gitmap-update-*.ps1"
)

// Update flags.
const FlagVerbose = "--verbose"
const FlagRepoPath = "--repo-path"
const FlagReportErrors = "--report-errors"
const FlagReportErrorsFile = "--report-errors-file"
const FlagDebugRepoDetect = "--debug-repo-detect"
const FlagDebugWindows = "--debug-windows"
const FlagDebugWindowsJSON = "--debug-windows-json"

// Debug-windows env bridge — propagated to the detached cleanup child so
// the dump runs on both sides of the handoff without requiring users to
// pass the flag twice.
const EnvDebugWindows = "GITMAP_DEBUG_WINDOWS"

// Debug-windows JSON-sink env bridge. Set automatically when the sink
// opens so the Phase 3 cleanup child appends to the SAME file as the
// parent, giving one consolidated NDJSON trace per handoff.
const EnvDebugWindowsJSON = "GITMAP_DEBUG_WINDOWS_JSON"

// Default filename pattern for the JSON sink. %s = local timestamp
// formatted as 2006-01-02_15-04-05.
const DebugWindowsJSONFileFmt = "gitmap-debug-windows-%s.jsonl"

// Report-errors values and env-var bridge to run.ps1 / run.sh.
const (
	ReportErrorsJSON        = "json"
	EnvReportErrorsFormat   = "GITMAP_REPORT_ERRORS"
	EnvReportErrorsFile     = "GITMAP_REPORT_ERRORS_FILE"
	EnvUpdateCleanupDelayMS = "GITMAP_UPDATE_CLEANUP_DELAY_MS"
	UpdateHandoffLogNameFmt = "gitmap-update-handoff-%s.log"
	ReportErrorsFilePrefix  = "gitmap-update-report-"
	ReportErrorsFileSuffix  = ".jsonl"
)

// Report-errors UI messages.
const (
	MsgReportErrorsEnabled = "  → Error report (json): %s\n"
	MsgReportErrorsSummary = "  → Wrote %d non-fatal failure entr%s to %s\n"
	WarnReportErrorsCreate = "  ⚠ Could not create report-errors file at %s: %v\n"
	ErrReportErrorsValue   = "  ✗ Unsupported --report-errors value: %q (supported: json)\n"
)

// Debug-repo-detect bridge.
const (
	EnvDebugRepoDetect   = "GITMAP_DEBUG_REPO_DETECT"
	MsgDebugRepoDetectOn = "  → Repo-detect debug enabled (will print marker checks)\n"
)

// Updater fallback.
const (
	UpdaterBin                 = "gitmap-updater"
	MsgUpdaterFallback         = "  → No source repo found. Delegating to %s...\n\n"
	UpdateCleanupSourceConfig  = "config"
	UpdateCleanupSourceSibling = "sibling"
	UpdateCleanupSourcePath    = "PATH"
	UpdateCleanupSourceUnknown = "unknown"
)

// Update UI messages.
const (
	MsgUpdateActive              = "  → Active: %s\n  → Handoff: %s\n"
	MsgUpdateCleanStart          = "\n  Cleaning up update artifacts..."
	MsgUpdateCleanBinary         = "  → Cleanup binary: %s\n"
	MsgUpdateCleanDelay          = "  → Waiting %dms for deploy handles to settle...\n"
	MsgUpdateCleanDone           = "  ✓ Removed %d file(s)\n\n"
	MsgUpdateCleanNone           = "  ✓ Nothing to clean up"
	MsgUpdateTempRemoved         = "  → Removed temp copy: %s\n"
	MsgUpdateOldRemoved          = "  → Removed backup: %s\n"
	MsgUpdatePhase3Handoff       = "\n  → Handing off cleanup to deployed binary: %s update-cleanup\n"
	MsgUpdatePhase3Resolve       = "  → Cleanup target resolved via: %s\n"
	MsgUpdatePhase3Target        = "  → Cleanup target path: %s\n"
	MsgUpdatePhase3Started       = "  → Cleanup process started (pid=%d)\n"
	MsgUpdatePhase3LogFile       = "  → Handoff log file: %s\n"
	UpdateRunnerLogStart         = "update-runner starting, repo=%s"
	UpdateScriptLogExec          = "executing update script: %s"
	UpdateScriptLogExit          = "update script exited: err=%v"
	UpdatePhase3LogResolve       = "update-cleanup target resolved via %s: %s"
	UpdatePhase3LogStarted       = "update-cleanup started pid=%d target=%s"
	UpdatePhase3LogStartFail     = "update-cleanup launch failed target=%s err=%v"
	UpdatePhase3LogTargetMissing = "update-cleanup handoff skipped: deployed target could not be resolved"
	UpdatePhase3LogInline        = "update-cleanup running inline from deployed binary: %s"
	UpdateCleanupLogStart        = "update-cleanup starting: self=%s"
	UpdateCleanupLogDone         = "update-cleanup finished: removed=%d"
	UpdateCleanupLogDelayInvalid = "update-cleanup ignored invalid delay value: %q"
)

// --debug-windows dump messages. All routed via fmt.Fprintf(os.Stderr, ...)
// so the dump survives even when stdout is being captured by a parent
// process. The dump is intentionally verbose and prefixed with
// `[debug-windows]` so it's grep-friendly in CI logs and bug reports.
const (
	MsgDebugWinHeader       = "\n[debug-windows] ===== update-cleanup handoff diagnostics =====\n"
	MsgDebugWinPhase        = "[debug-windows] phase            : %s\n"
	MsgDebugWinGOOS         = "[debug-windows] GOOS             : %s\n"
	MsgDebugWinSelf         = "[debug-windows] self executable  : %s\n"
	MsgDebugWinPID          = "[debug-windows] self pid         : %d\n"
	MsgDebugWinPPID         = "[debug-windows] parent pid       : %d\n"
	MsgDebugWinLogFile      = "[debug-windows] handoff log file : %s\n"
	MsgDebugWinSource       = "[debug-windows] resolution source: %s\n"
	MsgDebugWinTarget       = "[debug-windows] resolved target  : %s\n"
	MsgDebugWinTargetExists = "[debug-windows] target exists    : %t\n"
	MsgDebugWinChildArgv    = "[debug-windows] child argv       : %v\n"
	MsgDebugWinChildEnv     = "[debug-windows] %-16s : %s\n"
	MsgDebugWinChildPID     = "[debug-windows] spawned child pid: %d\n"
	MsgDebugWinNote         = "[debug-windows] %s\n"
	MsgDebugWinFooter       = "[debug-windows] ============================================\n\n"

	// Command-line rendering: exact invocation that Phase 3 will run,
	// formatted so the user can copy-paste it into a shell to reproduce.
	MsgDebugWinCmdLine = "[debug-windows] spawn command    : %s\n"
	MsgDebugWinCmdNote = "[debug-windows] (no `git` subprocess is launched by update-cleanup; only the line above plus os.Remove/os.RemoveAll on the paths below)\n"

	// Cleanup plan: the actual filesystem operations the deployed binary
	// will perform, enumerated BEFORE any deletion happens.
	MsgDebugWinCleanHdr      = "[debug-windows] ----- planned cleanup operations -----\n"
	MsgDebugWinCleanGlob     = "[debug-windows] glob             : %s\n"
	MsgDebugWinCleanMatch    = "[debug-windows]   → os.Remove    : %s\n"
	MsgDebugWinCleanEmpty    = "[debug-windows]   (no matches)\n"
	MsgDebugWinCleanSwap     = "[debug-windows]   → os.RemoveAll : %s\n"
	MsgDebugWinCleanShim     = "[debug-windows] drive-root shim  : %s (%s)\n"
	MsgDebugWinCleanShimSkip = "skipped"
	MsgDebugWinCleanShimDel  = "will os.Remove"
	MsgDebugWinCleanFooter   = "[debug-windows] --------------------------------------\n"

	// JSON sink — appended NDJSON events alongside the console dump.
	MsgDebugWinJSONFile     = "[debug-windows] json sink file   : %s\n"
	MsgDebugWinJSONOpenFail = "[debug-windows] json sink open failed: %s (%v)\n"
)

// Update error messages.
const (
	ErrUpdateExecFind            = "Error finding executable: %v\n"
	ErrUpdateCopyFail            = "Error creating update copy: %v\n"
	ErrUpdateNoRunSH             = "  ✗ run.sh not found at %s — cannot update on this platform without it.\n"
	ErrUpdateCleanupExecPath     = "Error: could not resolve executable path at active-binary: %v (operation: resolve executable, reason: os.Executable failed)\n"
	ErrUpdateCleanupConfigRead   = "Error: could not read cleanup config at %s: %v (operation: read config, reason: cleanup path resolution unavailable)\n"
	ErrUpdateCleanupGlob         = "Error: could not enumerate cleanup matches at %s: %v (operation: glob, reason: invalid cleanup pattern)\n"
	ErrUpdateCleanupRemove       = "Error: could not remove cleanup artifact at %s: %v (operation: remove, reason: file may be locked or missing)\n"
	ErrUpdateCleanDelayInvalid   = "Error: ignored invalid update-cleanup delay value %q (operation: parse cleanup delay, reason: env must be a positive integer)\n"
	ErrUpdatePhase3TargetMissing = "Error: could not resolve deployed binary for update-cleanup (operation: resolve cleanup handoff target, reason: no deployed binary path found)\n"
	ErrUpdatePhase3Handoff       = "Error: could not run update cleanup via deployed binary at %s: %v (operation: launch cleanup handoff, reason: target start or execution failed)\n"
)

// Unix update messages.
const (
	MsgUpdateInstallDir = "  → Installed directory: %s\n"
)

// Update path resolution messages.
const (
	MsgUpdatePathMissing = "\n  ⚠ The saved source repository path no longer exists on disk.\n"
	MsgUpdatePathPrompt  = "  Enter the new path to the gitmap source repo: "
	ErrUpdatePathInvalid = "  ✗ Directory not found at %s (operation: resolve, reason: file does not exist)\n"
)

// Clone-on-missing-path constants.
const (
	SourceRepoCloneURL   = "https://github.com/alimtvnetwork/gitmap-v27.git"
	MsgUpdateCloning     = "\n  ■ Path does not exist. Cloning gitmap source into %s...\n"
	MsgUpdateCloneOK     = "  ✓ Cloned successfully.\n"
	ErrUpdateCloneFailed = "  ✗ Clone failed: %v\n"
)

// Remote-installer update path (v5.51.0+).
//
// `gitmap update` downloads + runs the canonical install script straight
// from the gitmap repo. The installer itself handles the parallel
// `-v<N+i>` sibling-repo probe to find the latest published gitmap-vN
// release, so we don't reimplement that logic here.
//
// Pass `--source-rebuild` to fall back to the legacy in-tree rebuild
// flow (requires a local source checkout + Go toolchain).
const (
	UpdateRemoteInstallerPwsh = "https://raw.githubusercontent.com/" +
		"alimtvnetwork/gitmap-v27/main/install.ps1"
	UpdateRemoteInstallerBash = "https://raw.githubusercontent.com/" +
		"alimtvnetwork/gitmap-v27/main/install.sh"

	FlagSourceRebuild = "--source-rebuild"

	MsgUpdateRemoteFetch    = "\n  ■ Fetching remote installer: %s\n"
	MsgUpdateRemoteRun      = "  ■ Running installer: %s\n"
	MsgUpdateRemoteDone     = "\n  ✓ Update complete (via remote installer).\n\n"
	MsgUpdateRemoteFallback = "  ⚠ Remote installer failed — falling back to source rebuild...\n"

	ErrUpdateRemoteDownload = "  ✗ Failed to download installer: %v\n"
	ErrUpdateRemoteRun      = "  ✗ Installer failed: %v\n"
)

// Native sibling-repo probe (v5.52.0+). Spec: spec/01-app/111-update-remote-probe.md.
//
// `gitmap update` resolves the winning gitmap-vN repo slug in Go using a
// 20-parallel HEAD probe, then downloads THAT repo's installer. The
// downloaded installer no longer needs to probe — Go is the single
// source of truth.
const (
	UpdateProbeMaxSiblings = 20
	UpdateProbeTimeoutSec  = 5
	UpdateProbeRepoBase    = "gitmap"
	UpdateRepoOwner        = "alimtvnetwork"
	UpdateCurrentRepoSlug  = "gitmap-v27"

	UpdateRepoHEADTmpl     = "https://github.com/%s/%s"
	UpdateRawInstallerTmpl = "https://raw.githubusercontent.com/%s/%s/main/%s"
	UpdateReleasesAPITmpl  = "https://api.github.com/repos/%s/%s/releases/latest"

	UpdateInstallerNamePwsh = "install.ps1"
	UpdateInstallerNameBash = "install.sh"

	FlagProbeOnly = "--probe-only"
	FlagNoProbe   = "--no-probe"

	UpdateProbeSourceSibling = "sibling-probe"
	UpdateProbeSourceRelease = "current-release"
	UpdateProbeSourceMain    = "current-main"

	MsgUpdateProbeStart   = "\n  ■ Probing sibling repos (%s-v%d..%s-v%d)...\n"
	MsgUpdateProbeHit     = "    hit: %s (HTTP %d)\n"
	MsgUpdateProbeResolve = "  ✓ Resolved: %s (source: %s)\n"
	MsgUpdateProbeSkipped = "  → Probe skipped (--no-probe); using current repo: %s\n"
	MsgUpdateProbeOnly    = "  → --probe-only: %s (source: %s)\n"

	ErrUpdateProbeParseSlug = "could not parse current repo slug %q: %w"
	ErrUpdateProbeNoResolve = "  ✗ Could not resolve any gitmap-vN repo (probe, release, main all failed)\n"
)

// Update PowerShell script template sections.
const (
	UpdatePSHeader = `# gitmap self-update script (auto-generated)
Set-Location "%s"

# Refresh run.ps1 from origin BEFORE invoking it, so a stale/buggy local
# copy can't break the update flow (e.g. positional-binding errors from
# old code paths). Best-effort: silently skip if git is unavailable or
# the repo has uncommitted run.ps1 changes.
try {
    $gitCmd = Get-Command git -ErrorAction SilentlyContinue
    if ($gitCmd) {
        $statusOut = & git status --porcelain -- run.ps1 2>$null
        if ([string]::IsNullOrWhiteSpace($statusOut)) {
            & git fetch --quiet origin 2>$null | Out-Null
            $headBranch = (& git symbolic-ref --quiet --short HEAD 2>$null)
            if ($headBranch) {
                & git checkout --quiet "origin/$headBranch" -- run.ps1 2>$null | Out-Null
                Write-Host "  [INFO] Refreshed run.ps1 from origin/$headBranch" -ForegroundColor DarkGray
            }
        } else {
            Write-Host "  [INFO] Local run.ps1 has uncommitted changes; skipping refresh" -ForegroundColor DarkGray
        }
    }
} catch {
    Write-Host "  [WARN] Could not refresh run.ps1: $_" -ForegroundColor Yellow
}
`
	// UpdatePSDeployDetect format args (in order):
	//   %[1]s — repo path           (e.g. C:\dev\gitmap-v27)
	//   %[2]s — gitmap subdir       (sourceRepoSubdir from manifest)
	//   %[3]s — app subdir          (appSubdir from manifest, e.g. gitmap-cli)
	//   %[4]s — binary name         (gitmap.exe)
	//   %[5]s — known-subdir array  (PowerShell @("a","b") literal of app + legacy)
	UpdatePSDeployDetect = `
$configPath = Join-Path "%[1]s" "%[2]s\powershell.json"
$deployedBinary = $null
$configDeployedBinary = $null
if (Test-Path $configPath) {
    $cfg = Get-Content $configPath | ConvertFrom-Json
    if ($cfg.deployPath) {
	    $configDeployedBinary = Join-Path $cfg.deployPath "%[3]s\%[4]s"
    }
}

$activeCmdForDeploy = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
if ($activeCmdForDeploy -and (Test-Path $activeCmdForDeploy.Source)) {
    $resolvedActiveBinary = (Resolve-Path $activeCmdForDeploy.Source).Path
    $resolvedActiveDir = Split-Path $resolvedActiveBinary -Parent
    if ((Split-Path $resolvedActiveDir -Leaf) -in %[5]s) {
        $effectiveDeployTarget = Split-Path $resolvedActiveDir -Parent
    } else {
        $effectiveDeployTarget = Split-Path $resolvedActiveDir -Parent
    }
    if ($effectiveDeployTarget) {
        $deployedBinary = Join-Path $effectiveDeployTarget "%[3]s\%[4]s"
    }
}

if ((-not $deployedBinary) -and $configDeployedBinary) {
    $deployedBinary = $configDeployedBinary
}
`
	UpdatePSVersionBefore = `
$activeBinary = $null
$activeBefore = "unknown"
$cmdBefore = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
if ($cmdBefore -and (Test-Path $cmdBefore.Source)) {
    $activeBinary = $cmdBefore.Source
    $activeBefore = & $activeBinary version 2>&1
}
`
	UpdatePSRunUpdate = `
Write-Host ""
Write-Host "  Starting update via run.ps1 -Update" -ForegroundColor Cyan
& "%s" -Update
$runExit = $LASTEXITCODE
if (($runExit -ne 0) -and ($runExit -ne $null)) {
    exit $runExit
}
`
	UpdatePSSync = `
# Auto-sync deployed binary to active PATH binary if they differ.
if ($activeBinary -and $deployedBinary -and (Test-Path $deployedBinary)) {
    $resolvedActive = (Resolve-Path $activeBinary -ErrorAction SilentlyContinue).Path
    $resolvedDeployed = (Resolve-Path $deployedBinary -ErrorAction SilentlyContinue).Path
    if ($resolvedActive -and $resolvedDeployed -and ($resolvedActive -ne $resolvedDeployed)) {
        $deployedVer = & $deployedBinary version 2>&1
        $activeVer = & $activeBinary version 2>&1
        if ($deployedVer -ne $activeVer) {
            Write-Host ""
            Write-Host "  Syncing deployed binary to active PATH location..." -ForegroundColor Cyan
            Write-Host "    From: $resolvedDeployed" -ForegroundColor DarkGray
            Write-Host "    To:   $resolvedActive" -ForegroundColor DarkGray
            $syncOK = $false

            # Step 1: Try direct Copy-Item.
            try {
                Copy-Item -Path $resolvedDeployed -Destination $resolvedActive -Force
                $syncOK = $true
                Write-Host "  [OK] Synced successfully." -ForegroundColor Green
            } catch {
                Write-Host "  [WARN] Copy-Item failed: $_" -ForegroundColor Yellow
            }

            # Step 2: Rename-then-copy fallback.
            if (-not $syncOK) {
                Write-Host "  Trying rename-then-copy fallback..." -ForegroundColor Cyan
                $backupPath = "$resolvedActive.old"
                try {
                    if (Test-Path $backupPath) { Remove-Item $backupPath -Force -ErrorAction SilentlyContinue }
                    Move-Item -Path $resolvedActive -Destination $backupPath -Force
                    Copy-Item -Path $resolvedDeployed -Destination $resolvedActive -Force
                    $syncOK = $true
                    Write-Host "  [OK] Synced via rename fallback." -ForegroundColor Green
                } catch {
                    Write-Host "  [WARN] Rename fallback failed: $_" -ForegroundColor Yellow
                    # Restore backup if rename succeeded but copy failed.
                    if ((Test-Path $backupPath) -and (-not (Test-Path $resolvedActive))) {
                        Move-Item -Path $backupPath -Destination $resolvedActive -Force -ErrorAction SilentlyContinue
                    }
                }
            }

            # Step 3: Kill stale processes and retry.
            if (-not $syncOK) {
                Write-Host "  Killing stale gitmap processes..." -ForegroundColor Cyan
                $currentPID = $PID
                $stale = Get-CimInstance Win32_Process -Filter "Name='gitmap.exe'" -ErrorAction SilentlyContinue |
                    Where-Object { $_.ProcessId -ne $currentPID }
                foreach ($proc in $stale) {
                    try {
                        Stop-Process -Id $proc.ProcessId -Force -ErrorAction SilentlyContinue
                        Write-Host "    Stopped PID $($proc.ProcessId)" -ForegroundColor DarkGray
                    } catch {}
                }
                if ($stale) { Start-Sleep -Milliseconds 500 }
                try {
                    Copy-Item -Path $resolvedDeployed -Destination $resolvedActive -Force
                    $syncOK = $true
                    Write-Host "  [OK] Synced after killing stale processes." -ForegroundColor Green
                } catch {
                    Write-Host "  [WARN] Still could not sync: $_" -ForegroundColor Yellow
                }
            }

            if (-not $syncOK) {
                Write-Host "  [HINT] Run 'gitmap doctor --fix-path' manually." -ForegroundColor Yellow
            }
        }
    }
}
`
	UpdatePSVersionAfter = `
$activeAfter = "unknown"
$deployedAfter = "unknown"
$cmdAfter = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
if ($cmdAfter -and (Test-Path $cmdAfter.Source)) {
    $activeBinary = $cmdAfter.Source
    $activeAfter = & $activeBinary version 2>&1
} else {
    Write-Host "  [TRACE] Get-Command gitmap: not found in PATH" -ForegroundColor DarkGray
}
if ($deployedBinary -and (Test-Path $deployedBinary)) {
    $deployedAfter = & $deployedBinary version 2>&1
} else {
    if (-not $deployedBinary) {
        Write-Host "  [TRACE] deployedBinary: not resolved (check powershell.json deployPath)" -ForegroundColor DarkGray
    } else {
        Write-Host "  [TRACE] deployedBinary: path not found: $deployedBinary" -ForegroundColor DarkGray
    }
}
`
	UpdatePSVerify = `
Write-Host ""
Write-Host "  Version before:   $activeBefore" -ForegroundColor DarkGray
Write-Host "  Version active:   $activeAfter" -ForegroundColor DarkGray
Write-Host "  Version deployed: $deployedAfter" -ForegroundColor DarkGray
Write-Host "  Active binary:    $activeBinary" -ForegroundColor DarkGray
Write-Host "  Deployed binary:  $(if ($deployedBinary) { $deployedBinary } else { '(not resolved)' })" -ForegroundColor DarkGray
if ($configDeployedBinary -and $deployedBinary -and ($configDeployedBinary -ne $deployedBinary)) {
    Write-Host "  Config binary:    $configDeployedBinary" -ForegroundColor DarkGray
}

$lastReleaseScript = Join-Path (Join-Path (Join-Path "%s" "gitmap") "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    & $lastReleaseScript -BinaryPath $activeBinary -RepoRoot "%s"
}

if ($activeAfter -ne "unknown" -and $deployedAfter -eq "unknown") {
    Write-Host ""
    Write-Host "  [WARN] Deployed binary could not be verified (not resolved or missing)." -ForegroundColor Yellow
    Write-Host "  [TRACE] activeAfter=$activeAfter  deployedAfter=$deployedAfter" -ForegroundColor DarkGray
    if ($configDeployedBinary -and $deployedBinary -and ($configDeployedBinary -ne $deployedBinary)) {
        Write-Host "  [HINT] powershell.json points to an older deploy location; using PATH-derived target for verification." -ForegroundColor Yellow
    } else {
        Write-Host "  [HINT] Check that powershell.json 'deployPath' points to the correct directory" -ForegroundColor Yellow
        Write-Host "         and that the binary exists at: $deployedBinary" -ForegroundColor Yellow
    }
    Write-Host "  [OK] Active PATH binary updated successfully: $activeAfter" -ForegroundColor Green
} elseif (($activeAfter -eq "unknown") -or ($activeAfter -ne $deployedAfter)) {
    Write-Host ""
    Write-Host "  [FAIL] Active PATH version does not match deployed version." -ForegroundColor Red
    Write-Host "  [TRACE] activeAfter=$activeAfter  deployedAfter=$deployedAfter" -ForegroundColor DarkGray
    if ($activeAfter -eq "unknown") {
        Write-Host "  [HINT] Active binary not found in PATH." -ForegroundColor Yellow
    } elseif ($configDeployedBinary -and $deployedBinary -and ($configDeployedBinary -ne $deployedBinary)) {
        Write-Host "  [HINT] powershell.json still references a different deploy location than the active PATH binary." -ForegroundColor Yellow
    }
    exit 1
} else {
    Write-Host "  [OK] Active PATH binary matches deployed version." -ForegroundColor Green
}
`
	UpdatePSPostActions = `
if ($activeBinary -and (Test-Path $activeBinary)) {
    Write-Host ""
    Write-Host "  Latest changelog:" -ForegroundColor Cyan
    & $activeBinary changelog --latest
}

# NOTE: cleanup is intentionally NOT invoked here. The handoff copy
# (gitmap-update-<pid>.exe) is still alive at this point because the
# Go process that spawned this PowerShell script has not yet exited,
# so any "& $activeBinary update-cleanup" call would race against the
# still-locked handoff exe + freshly-renamed .old backup and emit two
# scary "Access is denied" lines for files Windows simply cannot
# release until our process tree winds down.
#
# Cleanup is handled exactly once by Phase 3 (scheduleDeployedCleanupHandoff
# in gitmap/cmd/updatehandoff_phase3.go), which resolves the deployed binary
# from powershell.json first (sibling/path fallback only if needed), then
# launches "<deployed> update-cleanup" directly in a hidden process with a
# short delay. That avoids brittle cmd.exe quoting AND avoids stale PATH
# binaries hijacking cleanup on machines with duplicate installs. See
# spec/08-generic-update/06-cleanup.md.

Write-Host ""
exit 0
`
)

// Revert PowerShell script template sections.
const (
	RevertPSHeader = `# gitmap revert script (auto-generated)
Set-Location "%s"
`
	RevertPSBuild = `
Write-Host ""
Write-Host "  Building from checked-out version..." -ForegroundColor Cyan
& "%s"
$runExit = $LASTEXITCODE
if (($runExit -ne 0) -and ($runExit -ne $null)) {
    exit $runExit
}
`
	RevertPSPostActions = `
$cmdAfter = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
if ($cmdAfter -and (Test-Path $cmdAfter.Source)) {
    $activeAfter = & $cmdAfter.Source version 2>&1
    Write-Host "  Active version: $activeAfter" -ForegroundColor DarkGray
}

# NOTE: cleanup is intentionally NOT invoked here, mirroring the same
# decision in UpdatePSPostActions. Phase 3 (scheduleDeployedCleanupHandoff
# in runRevertRunner) handles cleanup after this PowerShell script and
# the surrounding handoff process exit, so the deployed binary can
# safely remove the still-locked handoff exe + .old backup. See
# spec/08-generic-update/06-cleanup.md.

Write-Host ""
exit 0
`
)

// Set-source-repo messages.
const (
	ErrSetSourceRepoNoPath  = "  ✗ set-source-repo requires a path argument\n"
	ErrSetSourceRepoInvalid = "  ✗ Invalid source repo path: %s\n"
	MsgSetSourceRepoDone    = "  ✓ Source repo path saved: %s\n"
)

// Backup file extension glob.
const OldBackupGlob = "*.old"

// PowerShell execution arguments.
const (
	PSBin            = "powershell"
	PSExecPolicy     = "-ExecutionPolicy"
	PSBypass         = "Bypass"
	PSNoProfile      = "-NoProfile"
	PSNoLogo         = "-NoLogo"
	PSFile           = "-File"
	PSNonInteractive = "-NonInteractive"
	PSCommand        = "-Command"
)
