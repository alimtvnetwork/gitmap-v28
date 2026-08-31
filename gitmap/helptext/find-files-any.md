# find-files-any

Searches and filters file names across the repository or directory matching any file name that contains the search words/substring with optional extension filtering.

## Usage

```bash
gitmap find-files-any <substring> [-ext <extensions>] [--limit <n>] [--json]
```

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `-ext`, `--ext`   | string  | ""      | Comma-separated extension filter (e.g. `"md, go"`, `"ts, py"`) |
| `--limit`, `-l`   | integer | 0       | Maximum number of results to return                             |
| `--json`          | boolean | false   | Output machine-readable JSON array of matched file paths        |
| `--dir`, `-d`     | string  | "."     | Root search directory                                           |

## Examples

### Substring File Search

Find all files containing `runner` in Python and Go files:

```bash
$ gitmap find-files-any "runner" -ext "py, go"
.lovable/ai-fix-scripts/06-cicd-local-runner.py
```

### JSON Format for Automation

```bash
$ gitmap find-files-any "sequence" --json
[
  "gitmap/cmd/sequence_cmd.go",
  "gitmap/cmd/sequence_cmd_test.go"
]
```

## See Also

- [find-files](find-files.md) — Exact filename search
- [find-files-startswith](find-files-startswith.md) — Prefix file search
- [find-files-endswith](find-files-endswith.md) — Suffix file search
