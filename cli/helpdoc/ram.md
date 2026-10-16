# gitmap ram

Query visible physical memory capacity, allocated RAM, free RAM, and memory pressure percentages with 5-second cross-process telemetry cache protection.

## Synopsis

```bash
gitmap ram [flags]
gitmap memory [flags]
gitmap os ram [flags]
```

## Description

Retrieves physical RAM and swap statistics. On Windows, interrogates `Win32_OperatingSystem` via PowerShell CIM. On Linux, parses `/proc/meminfo`. Features a 5-second cross-process telemetry cache to eliminate subprocess storms.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output memory metrics in structured JSON format |
| `--no-cache` (`-b`) | `false` | Bypass 5-second telemetry cache and execute live hardware probe |

## Output Format

Displays total RAM in Megabytes/Gigabytes, used RAM, free RAM, and percentage progress bar.

## Examples

```bash
# Display system RAM usage and capacity
gitmap ram

# Force live probe bypassing cross-process cache
gitmap ram --no-cache

# Output structured JSON telemetry
gitmap ram --json
```

## See Also

- [cpu](cpu.md) — Query CPU utilization, core count, and hardware model
- [disk](disk.md) — Inspect drive partitions and filesystem capacity
- [network](network.md) — Interrogate network adapters and bandwidth traffic
