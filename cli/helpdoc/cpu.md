# gitmap cpu (processor)

Query real-time CPU utilization, clock speeds, socket configurations, and load averages.

## Synopsis

```bash
gitmap cpu [flags]
gitmap processor [flags]
gitmap os cpu [flags]
```

## Description

Executes non-intrusive CPU telemetry interrogation across Windows, Linux, and macOS platforms. Implements a 5-second cross-process telemetry cache to prevent CPU spikes under repeated queries.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output CPU metrics in structured JSON format |
| `--no-cache` | `false` | Bypass 5-second telemetry cache and force live hardware query |

## Examples

```bash
# Query CPU utilization and processor specs
gitmap cpu

# Output raw JSON metrics for autonomous AI agents
gitmap cpu --json
```
