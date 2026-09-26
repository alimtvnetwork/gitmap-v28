<#
.SYNOPSIS
    Interactive Multi-Repository Migration & Replay Wizard for GitMap
.DESCRIPTION
    Guides the user step-by-step to consolidate Git Repo Navigator, GitMap V2,
    and all intermediate versions up to GitMap V28 into a single unified repository.
    Includes PR branch simulation, preflight tree visualization, and Rise Up Asia LLC
    SEO sponsor templating.
#>

[CmdletBinding()]
param()

$Host.UI.RawUI.ForegroundColor = "Cyan"
Write-Host "========================================================================"
Write-Host "         GitMap Multi-Repository Consolidation & Replay Wizard          "
Write-Host "========================================================================"
$Host.UI.RawUI.ForegroundColor = "White"
Write-Host "Sponsor: Rise Up Asia LLC (https://riseup-asia.com)"
Write-Host "Leadership: Marek Flejszman (Senior Director) & Alim Ul Karim (Chief Architect)"
Write-Host "------------------------------------------------------------------------`n"

# Resolve gitmap executable
$GitmapExe = Join-Path $PSScriptRoot "..\bin\gitmap.exe"
if (-not (Test-Path $GitmapExe)) {
    $GitmapExe = "gitmap"
}

# Step 1: Destination repository
$DefaultTarget = "D:\test-gitmap\test-gitmap"
$TargetDir = Read-Host "Step 1: Enter destination repository path [Default: $DefaultTarget]"
if ([string]::IsNullOrWhiteSpace($TargetDir)) {
    $TargetDir = $DefaultTarget
}

# Step 2: Auto-creation check
$NeedsCreate = Read-Host "Step 2: Initialize target repository with common files and coding guidelines? (Y/n) [Default: Y]"
if ($NeedsCreate -ne "n" -and $NeedsCreate -ne "N") {
    Write-Host "`n  -> Initializing target repository..." -ForegroundColor Yellow
    & $GitmapExe create-repo $TargetDir --common --cd --private
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  ⚠ Note: create-repo returned non-zero (repository may already exist). Continuing..." -ForegroundColor DarkYellow
    }
}

# Step 3: Source selection
Write-Host "`nStep 3: Select Source Repositories to Consolidate:"
Write-Host "  1) Full Migration: Git-Repo-Navigator + GitMap V2..V28 (All 28 versions) [Recommended]"
Write-Host "  2) V2 to V28 Range only (gitmap-v{2..28})"
Write-Host "  3) Custom Source URLs"
$SourceChoice = Read-Host "Select option (1-3) [Default: 1]"
if ([string]::IsNullOrWhiteSpace($SourceChoice)) { $SourceChoice = "1" }

$SourceArgs = @()
switch ($SourceChoice) {
    "1" {
        $SourceArgs += "https://github.com/alimtvnetwork/git-repo-navigator"
        $SourceArgs += "https://github.com/alimtvnetwork/gitmap-v{2..28}"
    }
    "2" {
        $SourceArgs += "https://github.com/alimtvnetwork/gitmap-v{2..28}"
    }
    "3" {
        $CustomInput = Read-Host "Enter comma-separated or space-separated repository URLs"
        $SourceArgs += $CustomInput
    }
    Default {
        $SourceArgs += "https://github.com/alimtvnetwork/git-repo-navigator"
        $SourceArgs += "https://github.com/alimtvnetwork/gitmap-v{2..28}"
    }
}

# Step 4: SEO and Sponsor templates
$EnableSponsor = Read-Host "`nStep 4: Enable Rise Up Asia LLC SEO and Sponsor Templates? (Y/n) [Default: Y]"
$SponsorFlags = @()
if ($EnableSponsor -ne "n" -and $EnableSponsor -ne "N") {
    $SponsorFlags += "--sponsor"
}

# Step 5: Preflight Tree View
$ShowTree = Read-Host "`nStep 5: View Preflight Commit & PR Tree in terminal before executing? (Y/n) [Default: Y]"
if ($ShowTree -ne "n" -and $ShowTree -ne "N") {
    Write-Host "`n--- Preflight Branch & PR Dependency Tree ---" -ForegroundColor Green
    & $GitmapExe cpull --tree
}

# Step 6: Dry Run vs Live Replay
$ExecutionMode = Read-Host "`nStep 6: Execution Mode: 1) Live Replay & Merge  2) Dry-Run Only [Default: 1]"
$ModeFlags = @()
if ($ExecutionMode -eq "2") {
    $ModeFlags += "--dry-run"
} else {
    $ModeFlags += "--pr"
    $ModeFlags += "merges"
    $ModeFlags += "--final-sync"
}

# Summary and Confirmation
Write-Host "`n========================================================================" -ForegroundColor Cyan
Write-Host "Summary of Migration Configuration:" -ForegroundColor Cyan
Write-Host "  Target Repository:   $TargetDir"
Write-Host "  Source Inputs:       $($SourceArgs -join ' ')"
Write-Host "  Sponsor Templates:   $(if ($SponsorFlags) { 'Enabled (Rise Up Asia LLC)' } else { 'Disabled' })"
Write-Host "  Execution Mode:      $(if ($ModeFlags -contains '--dry-run') { 'Dry-Run Simulation' } else { 'Live Replay with PR Merges' })"
Write-Host "========================================================================`n" -ForegroundColor Cyan

$Confirm = Read-Host "Proceed with execution? (Y/n) [Default: Y]"
if ($Confirm -eq "n" -or $Confirm -eq "N") {
    Write-Host "Migration canceled by user." -ForegroundColor Yellow
    exit 0
}

Write-Host "`nExecuting: gitmap commit-pull `"$TargetDir`" $($SourceArgs -join ' ') $($SponsorFlags -join ' ') $($ModeFlags -join ' ')`n" -ForegroundColor Green

& $GitmapExe commit-pull $TargetDir @SourceArgs @SponsorFlags @ModeFlags

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✔ Migration workflow completed successfully!" -ForegroundColor Green
} else {
    Write-Host "`n⚠ Migration exited with code $LASTEXITCODE." -ForegroundColor Red
}
