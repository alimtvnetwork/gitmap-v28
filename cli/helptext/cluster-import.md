# gitmap cluster import

Import or restore a cluster node registry from a JSON or CSV file.

## Usage

```bash
gitmap cluster import <file> [--merge|--replace]
```

## Flags

- `--merge`: Add new nodes and update existing without deleting (default).
- `--replace`: Overwrite the node registry with the imported file.

## Examples

```bash
gitmap cluster import nodes.json
gitmap cluster import nodes.csv --replace
```

See also: `gitmap cluster export`, `gitmap cluster nodes`
