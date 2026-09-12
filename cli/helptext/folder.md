# folder

Scans, visualizes, filters, and exports directory hierarchies and nested file structures in ASCII Tree, Markdown list, JSON, YAML, or flat text formats.

## Aliases

tree

## Usage

```bash
gitmap folder [directory] [output-file] [flags]
gitmap tree [directory] [flags]
```

## Flags

| Flag                  | Type    | Default | Description                                                        |
|-----------------------|---------|---------|--------------------------------------------------------------------|
| `--tree`              | boolean | true    | Render visual ASCII/Unicode directory tree                         |
| `--md`, `--markdown`  | boolean | false   | Render nested Markdown bullet list with sub-indented folders/files |
| `--json`              | boolean | false   | Output machine-readable JSON payload containing full file metadata |
| `--yaml`, `--yml`     | boolean | false   | Output nested YAML representation                                  |
| `--flat`, `--list`    | boolean | false   | Output flat list of paths with optional inline metadata            |
| `--details`, `-l`     | boolean | false   | Append rich metadata (sequence number, lines of code, size)        |
| `--except`, `--exclude` | string | ""      | Comma-separated exclude patterns (e.g. `"vendor/**, *.png, *.min.js"`)|
| `--ext`               | string  | ""      | Filter by comma-separated file extensions (e.g. `"go,ts,md"`)      |
| `--max-depth`         | integer | 0       | Maximum folder traversal depth (0 for unlimited)                   |
| `--only-text`         | boolean | false   | Filter to include only text-based files                            |
| `--only-binary`       | boolean | false   | Filter to include only binary assets                               |
| `-o`, `--out`         | string  | ""      | Output file destination (e.g. `structure.md`, `tree.json`)         |

## Examples

### Clean ASCII Tree View

```bash
$ gitmap folder .lovable

.lovable/
├── ai-fix-scripts/
│   ├── 01-index.md
│   ├── 02-shared-engine.py
│   ├── 03-file-manipulator.py
│   └── 06-cicd-local-runner.py
```

### Detailed ASCII Tree (Filename First, Details After)

```bash
$ gitmap folder .lovable --tree --details

.lovable/
├── ai-fix-scripts/
│   ├── 01-index.md (seq: 01, 106 lines, 3.8 KB)
│   ├── 02-shared-engine.py (seq: 02, 107 lines, 3.1 KB)
│   ├── 03-file-manipulator.py (seq: 03, 195 lines, 6.4 KB)
│   └── 06-cicd-local-runner.py (seq: 06, 227 lines, 10.5 KB)
```

### Nested Markdown List with Metadata (No Backticks)

```bash
$ gitmap folder .lovable --md --details

- .lovable/
  - ai-fix-scripts/
    - 01-index.md (seq: 01, 106 lines, 3.8 KB)
    - 02-shared-engine.py (seq: 02, 107 lines, 3.1 KB)
    - 03-file-manipulator.py (seq: 03, 195 lines, 6.4 KB)
    - 06-cicd-local-runner.py (seq: 06, 227 lines, 10.5 KB)
```

### Multi-Glob Condition Filtering (`--except`)

Exclude dependencies, images, and minified bundles across all subdirectories:

```bash
gitmap folder . --except "vendor/**, node_modules/**, *.png, *.min.js" --tree --details
```

### Machine-Readable JSON for AI Agents & UI Visualization

```bash
$ gitmap folder .lovable/ai-fix-scripts --json
{
  "root": ".lovable/ai-fix-scripts",
  "totalFiles": 13,
  "totalLines": 1420,
  "totalSizeBytes": 48200,
  "totalSizeFormatted": "47.1 KB",
  "files": [
    {
      "path": ".lovable/ai-fix-scripts/01-index.md",
      "filename": "01-index.md",
      "directory": ".lovable/ai-fix-scripts",
      "extension": ".md",
      "sizeBytes": 3828,
      "sizeFormatted": "3.8 KB",
      "linesOfCode": 106,
      "isBinary": false,
      "sequence": 1
    }
  ]
}
```

### Direct File Export

Save detailed Markdown directory map to a file:

```bash
gitmap folder . project-structure.md --md --details --except "vendor/**, *.png"
```

## See Also

- [sequence](sequence.md) — Sequence numbering and repo-scoped database caching
- [find](find.md) — Fast indexed SQLite search across repository files
- [llm](llm.md) — LLM integration capabilities

