# find

Finds files within the repository using exact matches, wildcards (`*ends`, `starts*`, `*contains*`), or extension filtering.

## Usage

```bash
gitmap find <pattern> [-ext <extensions>] [--limit <n>] [--json]
```

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `-ext`, `--ext`   | string  | ""      | Comma-separated extension filter (e.g. `"md, go"`, `"ts, py"`) |
| `--limit`, `-l`   | integer | 0       | Maximum number of results to return                             |
| `--json`          | boolean | false   | Output machine-readable JSON array of matched file paths        |
| `--dir`, `-d`     | string  | "."     | Root search directory                                           |

## Wildcard Support

- **Ends-with:** `gitmap find "*_test.go"`
- **Starts-with:** `gitmap find "01*"`
- **Contains:** `gitmap find "*sequence*"`
- **Exact:** `gitmap find "folder.go"`

## Examples

### Suffix Search with Extension Filter

```bash
gitmap find "*_test.go" -ext "go"
```

### Prefix Search with Multiple Extensions

```bash
gitmap find "01*" -ext "md, py"
```

### Substring / Contains Search

```bash
gitmap find "*runner*" -ext "py"
```

### JSON Format for AI & Scripts

```bash
gitmap find "*.go" --limit 5 --json
```

## See Also

- [find-files](find_files.md) — Dedicated exact, starts-with, and contains file search
- [folder](folder.md) — Directory hierarchy visualization and metadata export
- [search](search.md) — Content search across files

