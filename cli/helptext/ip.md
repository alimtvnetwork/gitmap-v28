# gitmap ip

Inspect, set, change, switch, or revert network IP configuration across Windows, Ubuntu, Debian, CentOS, and Fedora.

## Usage

```bash
gitmap ip [command] [arguments] [flags]
```

When invoked without arguments, \`gitmap ip\` prints the primary local IPv4 address.

## Commands

| Command | Description |
|---------|-------------|
| `show [iface]` | Display network interfaces, IPv4, gateway, DNS, and status |
| `set <ip> [gateway]` | Configure static IP with ICMP ping connectivity validation and auto-rollback |
| `change <ip> [gateway]` | Guided IP change workflow |
| `switch <dhcp|static>` | Switch network mode between DHCP and Static IP |
| `revert [iface]` | Restore previous network configuration from snapshot |
| `help` | Show command usage and options |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i, --interface <name>` | auto | Target network interface |
| `-g, --gateway <ip>` | "" | Default gateway IPv4 address |
| `-m, --mask <mask>` | `255.255.255.0` | Subnet mask |
| `--dns <ip,...>` | `8.8.8.8,1.1.1.1` | Comma-separated DNS servers |
| `--no-validate` | false | Bypass ICMP ping connectivity verification |
| `-n, --dry-run` | false | Inspect changes without applying to host network |

## Examples

```bash
# Print current local IP
gitmap ip

# Inspect all network interfaces
gitmap ip show

# Set static IP with ping verification and automatic rollback
gitmap ip set 192.168.1.150 192.168.1.1

# Switch interface to DHCP
gitmap ip switch dhcp

# Revert to last captured IP configuration snapshot
gitmap ip revert
```
