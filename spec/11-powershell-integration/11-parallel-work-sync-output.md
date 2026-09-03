# Parallel Work Sync Output

## Problem

When running parallel background jobs in PowerShell (e.g., `-uas` multi-site uploads), output arrives in non-deterministic order. This makes it difficult to:

1. Map results back to their original task order
2. Display a coherent summary table
3. Track progress across concurrent operations
4. Diagnose which specific job failed and why

## Solution: Indexed Result Array Pattern

Pre-allocate a synchronized result array with one slot per job. Each parallel job writes to its assigned index, guaranteeing ordered output regardless of completion order.

### Data Structure

```powershell

# Pre-allocate result slots in sequential order

$totalJobs = $plugins.Count * $sites.Count
$results = [System.Collections.ArrayList]::new()

$jobIndex = 0
foreach ($site in $sites) {
    foreach ($plugin in $plugins) {
        [void]$results.Add(@{
            Index    = $jobIndex
            Site     = $site.name
            Plugin   = $plugin.Name
            Version  = "unknown"
            Status   = "PENDING"
            ExitCode = $null
            Output   = ""
            Duration = $null
            Error    = $null
        })
        $jobIndex++
    }
}
```

### Job Execution

Each job receives its index and writes back a complete result object:

```powershell
$jobs = @()
$jobIndex = 0

foreach ($site in $sites) {
    foreach ($plugin in $plugins) {
        $currentIndex = $jobIndex

        $jobs += Start-Job -Name "job-$currentIndex" -ScriptBlock {
            param($Index, $Script, $PluginPath, $SiteUrl, $Username, $Password, $PluginName, $SiteName, $Version)

            $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
            $output = ""
            $exitCode = 1
            $error = $null

            try {
                $output = (& $Script -jc $JsonConfig -a 2>&1 | Out-String)
                $exitCode = if ($null -ne $LASTEXITCODE) { $LASTEXITCODE } else { 0 }
            } catch {
                $error = $_.Exception.Message
                $output = ($_ | Out-String).Trim()
            }

            $stopwatch.Stop()

            return @{
                Index    = $Index
                Site     = $SiteName
                Plugin   = $PluginName
                Version  = $Version
                Status   = if ($exitCode -eq 0) { "OK" } else { "FAILED" }
                ExitCode = $exitCode
                Output   = $output
                Duration = $stopwatch.Elapsed.TotalSeconds
                Error    = $error
            }
        } -ArgumentList $currentIndex, $quploadScript, $pluginPath, $siteUrl, $username, $password, $pluginName, $siteName, $version

        $jobIndex++
    }
}
```

### Result Collection (Ordered)

Collect results as jobs complete, then sort by index for ordered display:

```powershell
$completedResults = @()

foreach ($job in $jobs) {
    $result = Receive-Job -Job $job -Wait | Select-Object -First 1

    if ($null -eq $result) {
        # Extract index from job name
        $idx = [int]($job.Name -replace '^job-', '')
        $result = $results[$idx]  # Use pre-allocated slot as fallback
        $result.Status = "FAILED (no result)"
        $result.ExitCode = 1
    }

    $completedResults += $result
    Remove-Job -Job $job -Force
}

# Sort by original index for deterministic display

$orderedResults = $completedResults | Sort-Object { $_.Index }
```

### Structured Summary Display

```powershell
Write-Host ""
Write-Host "========================================" -ForegroundColor Magenta
Write-Host "  Multi-Site Upload Summary" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta

# Table header

$col1 = 22  # Site
$col2 = 28  # Plugin
$col3 = 10  # Version
$col4 = 8   # Status
$col5 = 8   # Time

Write-Host ("  {0,-$col1} {1,-$col2} {2,-$col3} {3,-$col4} {4,$col5}" -f "Site", "Plugin", "Version", "Status", "Time")
Write-Host ("  " + ("-" * ($col1 + $col2 + $col3 + $col4 + $col5 + 4)))

foreach ($r in $orderedResults) {
    $color = if ($r.Status -eq "OK") { "Green" } else { "Red" }
    $duration = if ($r.Duration) { "{0:N1}s" -f $r.Duration } else { "-" }

    Write-Host ("  {0,-$col1} {1,-$col2} {2,-$col3} {3,-$col4} {4,$col5}" -f `
        $r.Site, $r.Plugin, ("v" + $r.Version), $r.Status, $duration) -ForegroundColor $color
}

Write-Host ""
$successCount = ($orderedResults | Where-Object { $_.Status -eq "OK" }).Count
$failCount = $orderedResults.Count - $successCount
Write-Host "  Total: $($orderedResults.Count) | Success: $successCount | Failed: $failCount" -ForegroundColor $(if ($failCount -eq 0) { "Green" } else { "Yellow" })
Write-Host "========================================" -ForegroundColor Magenta
```

### Expected Output

```
========================================
  Multi-Site Upload Summary
========================================
  Site                   Plugin                       Version    Status     Time
  ---------------------------------------------------------------------------------
  Atto Property Demo     riseup-asia-uploader         v2.11.0    OK        3.2s
  Atto Property Demo     plugins-onboard              v1.0.0     OK        2.1s
  Test V1                riseup-asia-uploader         v2.11.0    OK        4.5s
  Test V1                plugins-onboard              v1.0.0     FAILED    1.8s
  Test V2                riseup-asia-uploader         v2.11.0    OK        3.8s
  Test V2                plugins-onboard              v1.0.0     OK        2.4s

  Total: 6 | Success: 5 | Failed: 1
========================================
```

## Design Principles

1. **Pre-allocate**: Create all result slots before starting any jobs. This guarantees every job has a known index and a fallback result if the job crashes.

2. **Immutable index**: The `Index` field is assigned once at creation and never changes. Jobs receive it as a parameter; it is not derived from completion order.

3. **Sort-on-display**: Results are collected in completion order (fastest first) but sorted by `Index` before display. This gives the user a predictable, sequential view.

4. **Rich result objects**: Each result contains enough context (site, plugin, version, duration, output, error) for both summary display and failure diagnostics without needing to re-query.

5. **Fallback slots**: If a job returns `$null` (crash, timeout), the pre-allocated slot provides enough metadata to display a meaningful failure row rather than silently dropping it.

## Applicability

This pattern applies to any parallel PowerShell workload that needs ordered output:

- Multi-site plugin uploads (`-uas`)
- Parallel ZIP creation (`-za`)
- Bulk site health checks
- Parallel test execution

## Trade-offs

| Aspect | Benefit | Cost |
|--------|---------|------|
| Pre-allocation | Guaranteed fallback for crashed jobs | Slightly more memory (negligible for <100 jobs) |
| Index sorting | Deterministic output order | O(n log n) sort (trivial for typical job counts) |
| Rich result objects | Self-contained diagnostics | More data passed between job boundaries |
| ArrayList vs array | Thread-safe append | Requires `[void]` cast on `.Add()` |

## Future Extensions

- **Progress tracking**: A shared `[hashtable]::Synchronized(@{})` can be polled from the main thread to show real-time "3/6 complete" progress.
- **Retry logic**: Failed slots can be re-queued with the same index, preserving display order.
- **JSON export**: The ordered results array can be serialized to `logs/uas-upload/results-{timestamp}.json` for CI/CD integration.
