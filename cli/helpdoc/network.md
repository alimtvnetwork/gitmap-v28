# gitmap network

Interrogate physical and virtual network adapters, MAC addresses, IPv4/IPv6 configurations, link speeds, and cumulative byte transfer metrics.

## Synopsis

```bash
gitmap network [flags]
gitmap os network [flags]
gitmap net [flags]
```

## Description

Inspects physical and virtual network adapters, hardware MAC addresses, assigned IPv4 and IPv6 addresses, operational link status (UP / DOWN), interface speed (e.g. 1.0 Gbps), and transmitted/received byte counters.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output network telemetry as a structured JSON document |
| `--all` (`-a`) | `false` | Include inactive (DOWN) or loopback interfaces |

## Examples

```bash
# Display formatted table of active network adapters
gitmap network

# Include all adapters (down or loopback)
gitmap network --all

# Output structured JSON telemetry
gitmap network --json
```

## See Also

- [os](os.md) — Operating system configuration and environment diagnostics
- [ip](ip.md) — Inspect and manage IP configuration
- [nodes-tasks](nodes-tasks.md) — Inspect active and pending AI agent tasks across fleet
