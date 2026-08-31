# find-files-endswith

Searches and filters file names across the repository or directory matching file names that end with a suffix string with optional extension filtering.

## Usage

```bash
gitmap find-files-endswith <suffix> [-ext <extensions>] [--limit <n>] [--json]
```

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `-ext`, `--ext`   | string  | ""      | Comma-separated extension filter (e.g. `"md, go"`, `"ts, py"`) |
| `--limit`, `-l`   | integer | 0       | Maximum number of results to return                             |
| `--json`          | boolean | false   | Output machine-readable JSON array of matched file paths        |
| `--dir`, `-d`     | string  | "."     | Root search directory                                           |

## Examples

### Suffix File Search

Find all Go test files ending with `_test.go`:

```bash
$ gitmap find-files-endswith "_test.go" -ext "go"
gitmap/cmd/sequence_cmd_test.go
gitmap/cmd/folder/folder_test.go
```

### JSON Format for Automation

```bash
$ gitmap find-files-endswith "_cmd.go" --json
[
  "gitmap/cmd/sequence_cmd.go",
  "gitmap/cmd/vscode_cmd.go"
]
```

## See Also

- [find-files](find-files.md) — Exact filename search
- [find-files-any](find-files-any.md) — Substring contains file search
- [find-files-startswith](find-files-startswith.md) — Prefix file search
