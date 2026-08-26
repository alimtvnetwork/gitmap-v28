<#
.SYNOPSIS  Provider/auth helpers for visibility-change.ps1.
#>

$ErrorActionPreference = 'Stop'

function Get-OriginUrl {
    $url = & git remote get-url origin 2>$null
    if ($LASTEXITCODE -ne 0 -or -not $url) { return $null }
    return $url.Trim()
}

function Resolve-Provider {
    param([string]$Url)
    if (-not $Url) { return $null }
    $allow = ($env:VISIBILITY_GITLAB_HOSTS -split ',') | ForEach-Object { $_.Trim() } | Where-Object { $_ }
    $host = ([regex]::Match($Url, '^(?:https?://|ssh://[^@]+@|[^@]+@)([^/:]+)')).Groups[1].Value
    if (-not $host) { return $null }
    if ($host -in 'github.com','ssh.github.com')   { return 'github' }
    if ($host -eq 'gitlab.com' -or $allow -contains $host) { return 'gitlab' }
    return $null
}

# lint-allow: function-length reason="flat URL-shape parser"
function Resolve-OwnerRepo {
    param([string]$Url)
    $trimmed = ($Url.TrimEnd('/'))
    if ($trimmed.EndsWith('.git')) { $trimmed = $trimmed.Substring(0, $trimmed.Length - 4) }
    $patterns = @(
        '^https?://[^/]+/(?<owner>[^/]+)/(?<repo>[^/]+)$',
        '^[^@]+@[^:]+:(?<owner>[^/]+)/(?<repo>[^/]+)$',
        '^ssh://[^@]+@[^/]+/(?<owner>[^/]+)/(?<repo>[^/]+)$'
    )
    foreach ($p in $patterns) {
        $m = [regex]::Match($trimmed, $p)
        if ($m.Success) { return "$($m.Groups['owner'].Value)/$($m.Groups['repo'].Value)" }
    }
    return $null
}

function Test-CliAvailable {
    param([string]$Name)
    return (Get-Command $Name -ErrorAction SilentlyContinue) -ne $null
}

# lint-allow: function-length reason="flat REST GET + parse"
function Get-CurrentVisibility {
    param([string]$Provider, [string]$Slug)
    if ($Provider -eq 'github') {
        $v = & gh repo view $Slug --json visibility -q .visibility 2>$null
        if ($LASTEXITCODE -ne 0) { return $null }
        return $v.Trim().ToLowerInvariant()
    }
    $json = & glab repo view $Slug -F json 2>$null
    if ($LASTEXITCODE -ne 0) { return $null }
    $obj = $json | ConvertFrom-Json
    return ($obj.visibility).ToString().ToLowerInvariant()
}
