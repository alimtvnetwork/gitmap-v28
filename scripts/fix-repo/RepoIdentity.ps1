<#
.SYNOPSIS  Repo-identity detection helpers for fix-repo.ps1.
#>

$ErrorActionPreference = 'Stop'

function Get-RepoRoot {
    $root = & git rev-parse --show-toplevel 2>$null
    if ($LASTEXITCODE -ne 0) { return $null }
    return $root.Trim()
}

function Get-RemoteUrl {
    $url = & git remote get-url origin 2>$null
    if ($LASTEXITCODE -eq 0 -and $url) { return $url.Trim() }
    $first = (& git remote -v 2>$null) | Where-Object { $_ -match '\(fetch\)$' } | Select-Object -First 1
    if (-not $first) { return $null }
    $parts = $first -split '\s+'
    if ($parts.Count -lt 2) { return $null }
    return $parts[1]
}

function _Try-MatchRemoteUrl {
    param([string]$Trimmed, [string]$Pattern)
    $m = [regex]::Match($Trimmed, $Pattern)
    if (-not $m.Success) { return $null }
    return [pscustomobject]@{
        Host  = $m.Groups['host'].Value
        Owner = $m.Groups['owner'].Value
        Repo  = $m.Groups['repo'].Value
    }
}

# lint-allow: function-length reason="flat URL-shape parser"
function ConvertFrom-RemoteUrl {
    param([string]$Url)
    if (-not $Url) { return $null }
    $trimmed = $Url.TrimEnd('/')
    if ($trimmed.EndsWith('.git')) { $trimmed = $trimmed.Substring(0, $trimmed.Length - 4) }
    $patterns = @(
        '^https?://(?<host>[^/:]+)(?::\d+)?/(?<owner>[^/]+)/(?<repo>[^/]+)$',
        '^git@(?<host>[^:]+):(?<owner>[^/]+)/(?<repo>[^/]+)$',
        '^ssh://git@(?<host>[^/:]+)(?::\d+)?/(?<owner>[^/]+)/(?<repo>[^/]+)$'
    )
    foreach ($pat in $patterns) {
        $r = _Try-MatchRemoteUrl -Trimmed $trimmed -Pattern $pat
        if ($r) { return $r }
    }
    return $null
}

function Split-RepoVersion {
    param([string]$RepoFull)
    $m = [regex]::Match($RepoFull, '^(?<base>.+)-v(?<n>\d+)$')
    if (-not $m.Success) { return $null }
    return [pscustomobject]@{
        Base    = $m.Groups['base'].Value
        Version = [int]$m.Groups['n'].Value
    }
}
