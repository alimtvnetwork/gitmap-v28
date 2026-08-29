# Smoke test wrapper: forwards to Python 3 cross-platform runner.
[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [ValidateSet('source', 'release')]
    [string]$Mode = 'source'
)

$ErrorActionPreference = 'Stop'
$script = Join-Path $PSScriptRoot 'smoke-installer.py'
& python $script $Mode
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}
