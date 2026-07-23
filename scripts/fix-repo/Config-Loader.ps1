<#
.SYNOPSIS  Config loader + ignore-matching for fix-repo.ps1.

.DESCRIPTION
    Loads a JSON config (default: <RepoRoot>/fix-repo.config.json) with:
      ignoreDirs:     array of repo-relative directory prefixes to skip
      ignorePatterns: array of glob patterns (** = any depth, * = within segment)

    Exposes:
      Import-FixRepoConfig -ConfigPath <opt> -RepoRoot <root>
      Test-IsIgnoredPath   -RelPath <rel>
#>

$ErrorActionPreference = 'Stop'

$Script:FixRepoIgnoreDirs     = @()
$Script:FixRepoIgnorePatterns = @()
$Script:FixRepoConfigPath     = $null

function Resolve-FixRepoConfigPath {
    param([string]$ExplicitPath, [string]$RepoRoot)
    if ($ExplicitPath) {
        if (-not (Test-Path -LiteralPath $ExplicitPath -PathType Leaf)) {
            throw "fix-repo: ERROR config file not found: $ExplicitPath"
        }
        return (Resolve-Path -LiteralPath $ExplicitPath).Path
    }
    $default = Join-Path $RepoRoot 'fix-repo.config.json'
    if (Test-Path -LiteralPath $default -PathType Leaf) { return $default }
    return $null
}

function Import-FixRepoConfig {
    param([string]$ConfigPath, [string]$RepoRoot)
    $resolved = Resolve-FixRepoConfigPath -ExplicitPath $ConfigPath -RepoRoot $RepoRoot
    $Script:FixRepoConfigPath     = $resolved
    $Script:FixRepoIgnoreDirs     = @()
    $Script:FixRepoIgnorePatterns = @()
    if (-not $resolved) { return }
    $raw  = Get-Content -LiteralPath $resolved -Raw -Encoding UTF8
    $data = $raw | ConvertFrom-Json
    if ($data.ignoreDirs)     { $Script:FixRepoIgnoreDirs     = @($data.ignoreDirs     | Where-Object { $_ }) }
    if ($data.ignorePatterns) { $Script:FixRepoIgnorePatterns = @($data.ignorePatterns | Where-Object { $_ }) }
}

function Test-PathStartsWithDir {
    param([string]$Rel, [string]$Dir)
    $d = $Dir.TrimEnd('/','\')
    if (-not $d) { return $false }
    $norm = $Rel -replace '\\','/'
    if ($norm -eq $d) { return $true }
    return $norm.StartsWith($d + '/')
}

function Test-IsIgnoredDir {
    param([string]$RelPath)
    foreach ($d in $Script:FixRepoIgnoreDirs) {
        if (Test-PathStartsWithDir -Rel $RelPath -Dir $d) { return $true }
    }
    return $false
}

function ConvertTo-FixRepoRegex {
    param([string]$Pattern)
    $sb = New-Object System.Text.StringBuilder
    [void]$sb.Append('^')
    $i = 0
    while ($i -lt $Pattern.Length) {
        $ch = $Pattern[$i]
        if ($ch -eq '*') {
            if ($i + 1 -lt $Pattern.Length -and $Pattern[$i+1] -eq '*') {
                [void]$sb.Append('.*'); $i += 2; continue
            }
            [void]$sb.Append('[^/]*'); $i++; continue
        }
        if ($ch -eq '?') { [void]$sb.Append('[^/]'); $i++; continue }
        if ('.+()[]{}^$|\'.Contains([string]$ch)) { [void]$sb.Append('\').Append($ch); $i++; continue }
        [void]$sb.Append($ch); $i++
    }
    [void]$sb.Append('$')
    return $sb.ToString()
}

function Test-IsIgnoredPattern {
    param([string]$RelPath)
    $norm = $RelPath -replace '\\','/'
    foreach ($p in $Script:FixRepoIgnorePatterns) {
        $re = ConvertTo-FixRepoRegex -Pattern $p
        if ($norm -match $re) { return $true }
    }
    return $false
}

function Test-IsIgnoredPath {
    param([string]$RelPath)
    if (Test-IsIgnoredDir     -RelPath $RelPath) { return $true }
    if (Test-IsIgnoredPattern -RelPath $RelPath) { return $true }
    return $false
}
