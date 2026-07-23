<#
.SYNOPSIS  Post-rewrite audit: catches paired-literal desync where a
           {base}-v{Current} token was rewritten but a sibling bare
           digit literal representing the previous version was left
           behind (e.g. `"gitmap-v28", "12"`). The Go-native rewriter
           catches these via -Strict + go test, but this audit gives
           a sub-second precise diagnostic that runs always.

           Audit is scoped to *_test.go files only — production code
           legitimately uses small integers next to module paths. See
           .lovable/memory/issues/2026-05-02-fixrepo-paired-literal-desync.md.
#>

$ErrorActionPreference = 'Stop'

# Look-ahead window: how many lines after a {base}-v{Current} hit we
# scan for a stale sibling digit. 2 covers the common Go map/slice
# layout (`"key": {field, "value"},` spans 1-2 lines).
$Script:AuditLookaheadLines = 2

function Test-IsTestFile {
    param([string]$Path)
    return $Path.ToLowerInvariant().EndsWith('_test.go')
}

# Returns $true when $Line contains a stale sibling — a quoted "<prev>"
# OR a bare integer <prev> not adjacent to other digits or a `v`.
# The `v` exclusion prevents matching `v12` itself.
function Test-LineHasStaleSibling {
    param([string]$Line, [int]$Prev)
    $rxQuoted = '"' + [regex]::Escape("$Prev") + '"'
    $rxBare   = '(^|[^v0-9])' + [regex]::Escape("$Prev") + '($|[^0-9])'
    return ($Line -match $rxQuoted) -or ($Line -match $rxBare)
}

function Find-PairedLiteralHits {
    param([string]$FullPath, [string]$Base, [int]$Current)
    $hits   = @()
    $lines  = [System.IO.File]::ReadAllLines($FullPath)
    $needle = "$Base-v$Current"
    $prev   = $Current - 1
    if ($prev -lt 1) { return $hits }
    for ($i = 0; $i -lt $lines.Length; $i++) {
        if (-not $lines[$i].Contains($needle)) { continue }
        $end = [Math]::Min($lines.Length - 1, $i + $Script:AuditLookaheadLines)
        for ($j = $i; $j -le $end; $j++) {
            if (Test-LineHasStaleSibling -Line $lines[$j] -Prev $prev) {
                $hits += [pscustomobject]@{ LineNo = $j + 1; Line = $lines[$j] }
                break
            }
        }
    }
    return $hits
}

function Invoke-PairedLiteralAudit {
    param([System.Collections.Generic.List[string]]$ChangedFiles, [string]$Base, [int]$Current, [bool]$DryRun)
    if ($DryRun) { Write-Host 'audit:   skipped (dry-run)'; return $true }
    $totalHits = 0; $filesWithHits = 0
    foreach ($f in $ChangedFiles) {
        if (-not (Test-IsTestFile -Path $f)) { continue }
        $hits = Find-PairedLiteralHits -FullPath $f -Base $Base -Current $Current
        if ($hits.Count -eq 0) { continue }
        $filesWithHits++
        $totalHits += $hits.Count
        foreach ($h in $hits) {
            [Console]::Error.WriteLine(("fix-repo: AUDIT paired-literal at {0}:{1}: '{2}-v{3}' on/near sibling literal '{4}'`n  -> line: {5}" -f $f, $h.LineNo, $Base, $Current, ($Current-1), $h.Line.Trim()))
        }
    }
    if ($totalHits -eq 0) { Write-Host 'audit:   no paired-literal desync detected'; return $true }
    [Console]::Error.WriteLine(("fix-repo: ERROR paired-literal audit failed: {0} hit(s) in {1} file(s) (E_PAIRED_LITERAL)" -f $totalHits, $filesWithHits))
    [Console]::Error.WriteLine('  see .lovable/memory/issues/2026-05-02-fixrepo-paired-literal-desync.md')
    [Console]::Error.WriteLine('  fix: derive sibling literals from the same int via fmt.Sprintf')
    return $false
}
