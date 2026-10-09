# gitmap os dns

Inspect, benchmark, switch, and revert system DNS servers across network adapters.

## Usage

```bash
gitmap os dns [subcommand]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| status (st) | Inspect active DNS servers on the default network adapter |
| set <provider> | Configure primary and secondary DNS (cloudflare, google, quad9, adguard) |
| dhcp (auto) | Revert DNS configuration to automatic DHCP assignment |
| benchmark (bench) | Measure lookup and ping latency across all supported secure resolvers |
| help | Display help and usage information |

## Supported Providers

| Provider | Primary IPv4 | Secondary IPv4 |
|----------|--------------|----------------|
| cloudflare | 1.1.1.1 | 1.0.0.1 |
| google | 8.8.8.8 | 8.8.4.4 |
| quad9 | 9.9.9.9 | 149.112.112.112 |
| adguard | 94.140.14.14 | 94.140.15.15 |

## Examples

### Inspect Current DNS

```bash
gitmap os dns status
```

### Switch to Cloudflare DNS

```bash
gitmap os dns cloudflare
```

### Switch to Google DNS

```bash
gitmap os dns set google
```

### Benchmark DNS Provider Latency

```bash
gitmap os dns benchmark
```

### Revert to Automatic DHCP

```bash
gitmap os dns dhcp
```
