<#
.SYNOPSIS Apply + verify visibility change via gh / glab.
#>

$ErrorActionPreference = 'Stop'

function Invoke-VisibilityApply {
    param([string]$Provider, [string]$Slug, [string]$Target)
    if ($Provider -eq 'github') {
        & gh repo edit $Slug --visibility $Target --accept-visibility-change-consequences
        return $LASTEXITCODE
    }
    & glab repo edit $Slug --visibility $Target
    return $LASTEXITCODE
}

function Test-VisibilityMatches {
    param([string]$Provider, [string]$Slug, [string]$Target)
    $actual = Get-CurrentVisibility -Provider $Provider -Slug $Slug
    return $actual -eq $Target
}

function Confirm-PublicChange {
    param([string]$Slug, [string]$Provider)
    if ([Console]::IsInputRedirected) { return $false }
    Write-Host ""
    Write-Host ("WARNING: About to make {0} PUBLIC on {1}." -f $Slug, $Provider)
    Write-Host "  Type 'yes' to continue, anything else aborts:"
    $answer = Read-Host
    return $answer -eq 'yes'
}
