# find-files

Searches and filters file names across the repository or directory using exact matching with optional extension filtering.

## Aliases

find-files-any, find-files-startswith, find-files-endswith, list-files

## Usage

```bash
gitmap find-files <exact-name> [-ext <extensions>] [--limit <n>] [--json]
```

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `-ext`, `--ext`   | string  | ""      | Comma-separated extension filter (e.g. `"md, go"`, `"ts, py"`) |
| `--limit`, `-l`   | integer | 0       | Maximum number of results to return                             |
| `--json`          | boolean | false   | Output machine-readable JSON array of matched file paths        |
| `--dir`, `-d`     | string  | "."     | Root search directory                                           |

## Examples

### Exact Filename Search

Find exact file names matching `01-index.md` restricted to Markdown files:

```bash
$ gitmap find-files "01-index.md" -ext "md, go"
.lovable/ai-fix-scripts/01-index.md
```

### Exact Go File Match

```bash
$ gitmap find-files "folder.go" -ext "go"
gitmap/cmd/folder/folder.go
```

## See Also

- [find-files-any](find-files-any.md) — Substring contains file search
- [find-files-startswith](find-files-startswith.md) — Prefix file search
- [find-files-endswith](find-files-endswith.md) — Suffix file search
- [find](find.md) — Universal repository file search
