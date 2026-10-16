# gitmap cpu

Query real-time CPU utilization percentage, logical core count, model identity, and load averages with 5-second cross-process telemetry cache protection.

## Synopsis

```bash
gitmap cpu [flags]
gitmap processor [flags]
gitmap os cpu [flags]
```

## Description

Executes non-intrusive CPU telemetry interrogation across Windows, Linux, and macOS platforms. Implements a 5-second cross-process telemetry cache to prevent CPU spikes under repeated queries.

- **Windows:** Queries `Win32_Processor` via PowerShell CIM cmdlets for `LoadPercentage` and clock speed.
- **Linux:** Reads and computes deltas from `/proc/stat` and parses `/proc/cpuinfo`.
- **macOS:** Queries `sysctl` for model identity and active core counts.
- **Cache Invariant:** Results are cached in OS temp directory for 5 seconds (`HardwareMetricsTTL = 5 * time.Second`) to eliminate subprocess storms.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output processor metrics as structured JSON |
| `--no-cache` (`-b`) | `false` | Bypass 5-second telemetry cache and execute live hardware probe |

## Examples

```bash
# Query CPU utilization and processor specs
gitmap cpu

# Force live hardware query bypassing cache
gitmap cpu --no-cache

# Output raw JSON metrics for autonomous AI agents
gitmap cpu --json
```

## See Also

- [ram](ram.md) — Query system memory capacity and active utilization
- [disk](disk.md) — Inspect drive partitions and filesystem capacity
- [network](network.md) — Interrogate network adapters and bandwidth traffic
