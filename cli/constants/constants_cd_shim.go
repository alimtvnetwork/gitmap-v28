package constants

// PowerShell command-shim files make `gitmap cd` work even before the
// profile-loaded function wrapper is active.
const (
	PowerShellShimFile        = "gitmap.ps1"
	PowerShellSingleQuote     = "'"
	PowerShellEscapedQuote    = "''"
	PowerShellShimTemplateFmt = `$real = Join-Path -Path '%[1]s' -ChildPath 'gitmap.exe'
if (-not (Test-Path -LiteralPath $real)) { Write-Error "gitmap executable not found: $real"; return }
$handoff = [IO.Path]::Combine([IO.Path]::GetTempPath(), "gitmap-handoff-$([Guid]::NewGuid().ToString('N')).txt")
try {
  $env:GITMAP_HANDOFF_FILE = $handoff; $env:GITMAP_WRAPPER = "1"; $env:GITMAP_COMMAND_WRAPPER = "1"
  & $real @args
  $exitCode = $LASTEXITCODE
  if ((Test-Path -LiteralPath $handoff) -and ((Get-Item -LiteralPath $handoff).Length -gt 0)) {
    $target = [string](Get-Content -LiteralPath $handoff -Raw); $target = $target.Trim()
    if ($target -and (Test-Path -LiteralPath ([string]$target))) { Set-Location -LiteralPath ([string]$target) }
  }

  $global:LASTEXITCODE = $exitCode

  return
}
finally {
  Remove-Item -LiteralPath $handoff -ErrorAction SilentlyContinue
  Remove-Item Env:\GITMAP_HANDOFF_FILE -ErrorAction SilentlyContinue
}

`
)
