# gitmap ram (memory)

Query system memory capacity, active RAM usage, free memory, and swap utilization.

## Synopsis

```bash
gitmap ram [flags]
gitmap os ram [flags]
```

## Description

Retrieves physical RAM and swap statistics. On Windows, interrogates `Win32_OperatingSystem` via PowerShell CIM. On Linux, parses `/proc/meminfo`. Features a 5-second cross-process cache.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output memory metrics in structured JSON format |
| `--no-cache` | `false` | Bypass 5-second telemetry cache and force live hardware query |

## Examples

```bash
# Display system RAM usage and capacity
gitmap ram

# JSON output
gitmap ram --json
```
