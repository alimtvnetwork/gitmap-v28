# gitmap cluster export

Export the cluster node registry to a file.

## Usage

```
gitmap cluster export [--format json|csv] [--output <file>]
```

## Description

Dumps the `ClusterNode` SQLite table. 

**Security Note:** Passwords are never exported. The `PasswordHash` field is strictly omitted from all export formats.

## Formats

- `json` (default): Exports an array of JSON objects.
- `csv`: Exports a CSV table.

## Examples

```
gitmap cluster export --format json --output nodes.json
gitmap cluster export --format csv > nodes.csv
```
