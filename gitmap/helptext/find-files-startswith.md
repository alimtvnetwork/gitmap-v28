# find-files-startswith

Searches and filters file names across the repository or directory matching file names that start with a prefix string with optional extension filtering.

## Usage

```bash
gitmap find-files-startswith <prefix> [-ext <extensions>] [--limit <n>] [--json]
```

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `-ext`, `--ext`   | string  | ""      | Comma-separated extension filter (e.g. `"md, go"`, `"ts, py"`) |
| `--limit`, `-l`   | integer | 0       | Maximum number of results to return                             |
| `--json`          | boolean | false   | Output machine-readable JSON array of matched file paths        |
| `--dir`, `-d`     | string  | "."     | Root search directory                                           |

## Examples

### Prefix File Search

Find all files starting with `01` in Markdown and Python files:

```bash
$ gitmap find-files-startswith "01" -ext "md, py"
.lovable/ai-fix-scripts/01-index.md
```

### Prefix Search in Code Files

```bash
$ gitmap find-files-startswith "folder" -ext "go"
gitmap/cmd/folder/folder.go
gitmap/cmd/folder/folder_test.go
```

## See Also

- [find-files](find-files.md) — Exact filename search
- [find-files-any](find-files-any.md) — Substring contains file search
- [find-files-endswith](find-files-endswith.md) — Suffix file search
