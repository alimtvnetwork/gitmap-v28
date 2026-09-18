# Output & Formatting Patterns

## Overview

This document describes reusable patterns for CLI tools that produce
multiple output formats simultaneously from a single data source.

## Multi-Format Output Strategy

### Principle: All Formats, Every Time

When a command runs, produce **all** output formats in one pass.
Don't make the user choose — generate everything and let them
pick what they need.

| Format | Destination | Purpose |
|--------|-------------|---------|
| Terminal (colored) | stdout | Immediate human feedback |
| CSV | file | Spreadsheet / data import |
| JSON | file | Machine-readable, re-import |
| Markdown | file | Documentation / review |
| Scripts | file | Automation / re-execution |

### Output Directory

All file outputs go into a dedicated subfolder (e.g., `toolname-output/`)
inside the scanned/processed directory. This keeps outputs organized
and avoids polluting the working directory.

```
target-dir/
├── project-a/
├── project-b/
└── toolname-output/       ← all outputs here
    ├── data.csv
    ├── data.json
    ├── structure.md
    └── scripts/
```

## Terminal Output

### Colored, Structured Reports

Use ANSI codes for visual hierarchy:

| Element | Color | Purpose |
|---------|-------|---------|
| Banner/headers | Cyan | Visual identity |
| Success markers | Green | Confirmed items |
| Warnings | Yellow | Non-fatal issues |
| Data values | White | Primary content |
| Metadata | Dim/Gray | Secondary info |

### All color codes live in `constants`:

```go
const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorCyan   = "\033[36m"
    ColorDim    = "\033[90m"
)
```

### Terminal Report Sections

Structure terminal output as distinct sections:

1. **Banner** — tool name + version + item count
2. **Item list** — each item with icon, path, and key data
3. **Tree visualization** — hierarchical folder structure
4. **Output file list** — what files were generated and where (with full paths)
5. **Action instructions** — step-by-step how to use the outputs
6. **Script shortcuts** — direct commands for automation scripts
7. **Related commands** — other commands the user can run next

### Action Instructions Section

Always end terminal output with actionable next-step instructions.
Show **multiple options** so the user can pick the best approach:

```
■ How to Use the Output
──────────────────────────────────────────
  1. Copy the output files to another machine:
     toolname-output/data.json  (or data.csv)

  2. Restore via JSON:
     toolname restore ./toolname-output/data.json --target-dir ./projects

  3. Restore via CSV:
     toolname restore ./toolname-output/data.csv --target-dir ./projects

  4. Or run a generated script directly:
     .\restore.ps1              # With progress & error handling
     .\restore-quick.ps1        # Raw commands only
```

### Banner Pattern

```
╔══════════════════════════════════════╗
║            toolname v1.0.0          ║
╚══════════════════════════════════════╝
  ✓ Found 12 items
```

## Template-Based Script Generation

### Approach: `go:embed` Templates

For complex script outputs (PowerShell, Bash), use Go's embedded
templates rather than string concatenation:

```go
//go:embed templates/restore.ps1.tmpl
var restoreTemplate string

func WriteRestoreScript(w io.Writer, data RestoreData) error {
    tmpl := template.Must(template.New("restore").Parse(restoreTemplate))
    return tmpl.Execute(w, data)
}
```

### Template Data Structures

Define clear data structures for template rendering:

```go
type RestoreData struct {
    Items      []ItemEntry
    BaseDir    string
    TotalCount int
}

type ItemEntry struct {
    URL    string
    Branch string
    Path   string
    Name   string
}
```

### Script Categories

| Type | Content | Use Case |
|------|---------|----------|
| Logic scripts | Progress bars, error handling, summaries | Interactive restoration |
| Direct scripts | Raw commands, no logic | Quick copy-paste execution |
| Registration scripts | Tool-specific integrations | Third-party tool registration |

## CSV Output

### Conventions

- Always include a header row.
- Use consistent column ordering: name, identifiers, metadata, paths.
- Quote fields that may contain commas or special characters.
- Use standard Go `encoding/csv` writer.

```
name,primaryUrl,altUrl,branch,relativePath,absolutePath,instruction,notes
```

## JSON Output

### Conventions

- Use 2-space indentation for readability.
- Output an array of record objects.
- Field names match the Go struct's `json` tags.
- The JSON output should be directly re-importable by the tool's
  restore/import command.

```go
encoder := json.NewEncoder(w)
encoder.SetIndent("", constants.JSONIndent)
encoder.Encode(records)
```

## Markdown Output

### Folder Structure Visualization

Render a tree using Unicode box-drawing characters:

```
├── project-a/
│   ├── 📦 **service** (`main`) — https://example.com/service.git
│   └── 📦 **api** (`develop`) — https://example.com/api.git
└── project-b/
    └── 📦 **frontend** (`main`) — https://example.com/frontend.git
```

| Character | Constant | Usage |
|-----------|----------|-------|
| `├──` | `TreeBranch` | Non-last child |
| `└──` | `TreeCorner` | Last child |
| `│   ` | `TreePipe` | Vertical continuation |
| `    ` | `TreeSpace` | No continuation |

## Formatter Package Structure

```
formatter/
├── terminal.go       Terminal (colored stdout)
├── csv.go            CSV file output
├── json.go           JSON file output
├── structure.go      Markdown folder tree
├── logicscript.go    Logic-based restore script
├── directscript.go   Raw command scripts
├── template.go       Shared template loading
└── templates/        Embedded .tmpl files
    ├── restore.ps1.tmpl
    └── restore-quick.ps1.tmpl
```

### Rules

- Each format has its own file.
- All formatters accept `io.Writer` as first argument (testable).
- Templates are embedded via `go:embed`, not loaded from disk.
- No format string literals in formatter files — use `constants`.
