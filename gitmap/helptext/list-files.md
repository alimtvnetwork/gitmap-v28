# list-files

Searches, filters, and lists file names across the repository or directory using exact matching, wildcards, or extension filtering.

## Usage

```bash
gitmap list-files <pattern> [-ext <extensions>] [--limit <n>] [--json]
```

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `-ext`, `--ext`   | string  | ""      | Comma-separated extension filter (e.g. `"md, go"`, `"ts, py"`) |
| `--limit`, `-l`   | integer | 0       | Maximum number of results to return                             |
| `--json`          | boolean | false   | Output machine-readable JSON array of matched file paths        |
| `--dir`, `-d`     | string  | "."     | Root search directory                                           |

## Examples

### Wildcard List

```bash
$ gitmap list-files "*sequence*"
gitmap/cmd/sequence_cmd.go
gitmap/cmd/sequence_cmd_test.go
gitmap/helptext/sequence.md
```

### Extension Filtering

```bash
$ gitmap list-files "*" -ext "md"
.lovable/ai-fix-scripts/01-index.md
README.md
```

## See Also

- [find-files](find-files.md) — Exact filename search
- [find](find.md) — Universal repository file search
- [folder](folder.md) — Directory hierarchy visualization and metadata export
