<#
.SYNOPSIS
    Imports SSH fleet nodes and credentials from repo-secrets into GitMap and optionally deploys across the cluster.

.DESCRIPTION
    Reads the SSH node configuration JSON (defaulting to D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json),
    imports all nodes with their RSA-encrypted passwords and private keys into GitMap's Split-DB,
    and optionally broadcasts the node configuration across all online remote machines.

.PARAMETER Path
    Path to the SSH nodes JSON file. Defaults to D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json.

.PARAMETER Deploy
    If set, automatically deploys the configuration to all reachable remote nodes via 'gitmap ssh deploy node-config'.

.PARAMETER Except
    Comma-separated list of node IDs, IPs, or aliases to exclude during deployment.

.EXAMPLE
    .\scripts\import-ssh-nodes.ps1
    Imports all 5 nodes from D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json.

.EXAMPLE
    .\scripts\import-ssh-nodes.ps1 -Deploy
    Imports nodes and deploys configuration across the entire cluster fleet.
#>

[CmdletBinding()]
param(
    [string]$Path = "D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json",
    [switch]$Deploy,
    [string]$Except
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path -LiteralPath $Path)) {
    $Fallback = "D:\work\repo-secrets\gitmap-ssh-nodes.json"
    if (Test-Path -LiteralPath $Fallback) {
        $Path = $Fallback
    } else {
        Write-Error "SSH nodes configuration JSON not found at '$Path' or '$Fallback'."
        exit 1
    }
}

Write-Host "[*] Importing SSH node configurations from: $Path" -ForegroundColor Cyan

$GitMapCmd = Get-Command "gitmap" -ErrorAction SilentlyContinue
if ($null -ne $GitMapCmd) {
    & gitmap ssh nodes import-json $Path
    if ($Deploy) {
        Write-Host "[*] Deploying node configuration across online fleet..." -ForegroundColor Cyan
        if ($Except) {
            & gitmap ssh deploy node-config --except $Except
        } else {
            & gitmap ssh deploy node-config
        }
    }
} else {
    $CliDir = Join-Path $PSScriptRoot "..\cli"
    if (Test-Path (Join-Path $CliDir "main.go")) {
        Push-Location $CliDir
        try {
            go run . ssh nodes import-json $Path
            if ($Deploy) {
                Write-Host "[*] Deploying node configuration across online fleet..." -ForegroundColor Cyan
                if ($Except) {
                    go run . ssh deploy node-config --except $Except
                } else {
                    go run . ssh deploy node-config
                }
            }
        } finally {
            Pop-Location
        }
    } else {
        Write-Error "gitmap CLI not found in PATH and Go source directory not found."
        exit 1
    }
}

Write-Host "[OK] SSH nodes successfully imported into GitMap." -ForegroundColor Green
