# gitmap cluster export

Export the cluster node registry to a JSON or CSV file.

## Usage

```bash
gitmap cluster export [--format json|csv] [--output <file>]
```

## Flags

| Flag | Description |
|------|-------------|
| `--format <type>` | Output format: `json` (default) or `csv` |
| `--output <file>` | Destination file path (defaults to stdout if omitted) |

## Examples

```bash
gitmap cluster export --format json --output nodes.json
gitmap cluster export --format csv > nodes.csv
```

See also: `gitmap cluster import`, `gitmap cluster nodes`
