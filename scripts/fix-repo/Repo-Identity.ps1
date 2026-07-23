# Repo-Identity.ps1 - resolve git repo root, remote URL, base/version split.

function Get-RepoRoot {
    try {
        $root = & git rev-parse --show-toplevel 2>$null
        if ($LASTEXITCODE -ne 0) { return $null }
        return ($root | Select-Object -First 1).Trim()
    } catch { return $null }
}

function Get-RemoteUrl {
    try {
        $url = & git config --get remote.origin.url 2>$null
        if ($LASTEXITCODE -ne 0 -or -not $url) { return $null }
        return ($url | Select-Object -First 1).Trim()
    } catch { return $null }
}

function ConvertFrom-RemoteUrl {
    param([string]$Url)
    if (-not $Url) { return $null }
    $u = $Url.Trim()
    if ($u.EndsWith('.git')) { $u = $u.Substring(0, $u.Length - 4) }

    $sshMatch = [regex]::Match($u, '^[^@]+@([^:]+):([^/]+)/(.+)$')
    if ($sshMatch.Success) {
        return [pscustomobject]@{
            RepoHost = $sshMatch.Groups[1].Value
            Owner    = $sshMatch.Groups[2].Value
            Repo     = $sshMatch.Groups[3].Value
        }
    }

    $httpMatch = [regex]::Match($u, '^https?://([^/]+)/([^/]+)/(.+)$')
    if ($httpMatch.Success) {
        return [pscustomobject]@{
            RepoHost = $httpMatch.Groups[1].Value
            Owner    = $httpMatch.Groups[2].Value
            Repo     = $httpMatch.Groups[3].Value
        }
    }

    return $null
}

function Split-RepoVersion {
    param([string]$RepoFull)
    if (-not $RepoFull) { return $null }
    $m = [regex]::Match($RepoFull, '^(.+)-v(\d+)$')
    if (-not $m.Success) { return $null }
    return [pscustomobject]@{
        Base    = $m.Groups[1].Value
        Version = [int]$m.Groups[2].Value
    }
}

function Get-TargetVersions {
    param([int]$Current, [int]$Span)
    if ($Span -le 0 -or $Current -le 1) { return @() }
    $start = [Math]::Max(1, $Current - $Span)
    $end   = $Current - 1
    if ($start -gt $end) { return @() }
    return @($start..$end)
}
