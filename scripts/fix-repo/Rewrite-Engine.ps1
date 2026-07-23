<#
.SYNOPSIS  Token-rewrite engine for fix-repo.ps1.
#>

$ErrorActionPreference = 'Stop'

function Get-TargetVersions {
    param([int]$Current, [int]$Span)
    $start = [Math]::Max(1, $Current - $Span)
    $end   = $Current - 1
    if ($end -lt $start) { return @() }
    return $start..$end
}

function Get-RewritePattern {
    param([string]$Base, [int]$N)
    $escaped = [regex]::Escape("$Base-v$N")
    return "$escaped(?!\d)"
}

function _Apply-OneTarget {
    param([string]$Text, [string]$Base, [int]$N, [int]$Current)
    $pattern  = Get-RewritePattern -Base $Base -N $N
    $replaced = [regex]::Replace($Text, $pattern, "$Base-v$Current")
    if ($replaced -eq $Text) { return [pscustomobject]@{ Text=$Text; Added=0 } }
    $added = ([regex]::Matches($Text, $pattern)).Count
    return [pscustomobject]@{ Text=$replaced; Added=$added }
}

function Invoke-FileRewrite {
    param([string]$FullPath, [string]$Base, [int[]]$Targets, [int]$Current, [bool]$DryRun)
    $original = [System.IO.File]::ReadAllText($FullPath)
    $updated  = $original
    $count    = 0
    foreach ($n in $Targets) {
        $r       = _Apply-OneTarget -Text $updated -Base $Base -N $n -Current $Current
        $updated = $r.Text
        $count  += $r.Added
    }
    if ($count -eq 0) { return 0 }
    if (-not $DryRun) { [System.IO.File]::WriteAllText($FullPath, $updated) }
    return $count
}
