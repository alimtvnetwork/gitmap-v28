# child-path

High-speed directory inventory and AST-aware file search designed as a native Go replacement for slow PowerShell `Get-ChildItem -Filter` and repository-bound `git grep`.

## Overview

During autonomous code refactoring and architecture audits, developer tools and AI coding agents frequently need to inventory directories, discover matching source files, and filter entries by extension, depth, or type. Historically, workflows relied on tools that introduced severe performance and reliability problems:

1. **PowerShell `Get-ChildItem -Filter` Bottlenecks:**
   - High runtime overhead from heavy .NET `FileInfo` object allocations, causing multi-second lags in medium-to-large workspaces.
   - Cross-platform incompatibilities and parameter disparities between Windows PowerShell 5.1 and modern PowerShell Core (`pwsh`).
   - UTF-16 / ANSI encoding artifacts, localized date formatting, and line-wrapping that breaks regex scrapers and AI agents.

2. **`git grep` Limitations:**
   - Inability to scan untracked files, local scratch directories, newly scaffolded templates, or `.ai-memory` directories excluded by `.gitignore`.
   - Cannot query directories outside the git worktree root.
   - Does not emit file metadata such as byte size, modification timestamp, or structured directory hierarchy.

`gitmap child-path` addresses these challenges with a compiled, zero-dependency Go implementation:
- **0ms Fast Path (`--depth 1`):** Direct kernel directory read (`os.ReadDir`) that bypasses recursive tree walking for immediate child inspection.
- **Deep Pruned Traversal (`--depth > 1`):** Bounded `filepath.WalkDir` with automatic pruning of heavy non-source directories (`.git`, `node_modules`, `.gitmap`, `.tmp`).
- **Normalized Cross-Platform Paths:** Forward-slash (`/`) relative paths standard across Windows, macOS, and Linux.
- **Machine-Readable JSON (`--json`):** Emits structured JSON array `[{path, name, isDir, size, modTime}]` for direct ingestion by AI agents and scripts.

## Aliases

`childpath`, `cpth`, `gitmap aum child-path`

## Usage

```bash
gitmap child-path [dir] [glob] [flags]
gitmap childpath [dir] [glob] [flags]
gitmap cpth [dir] [glob] [flags]
gitmap aum child-path [dir] [glob] [flags]
```

## Flags

| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--depth <n>` | `-d <n>` | integer | `1` | Traversal depth (`1` = immediate children, `0` = unlimited recursion, `>1` = bounded depth). |
| `--ext <exts>` | `-e <exts>` | string | `""` | Comma-delimited list of extensions to match (e.g. `go,ts`, `.json,.md`). |
| `--type <t>` | `-t <t>` | string | `all` | Filter by entry kind: `f` (files only), `d` (directories only), or `all` (both). |
| `--limit <n>` | `-l <n>`, `-n <n>` | integer | `0` | Maximum number of entries to emit (`0` = unlimited). |
| `--json` | `-j` | boolean | `false` | Output structured JSON array `[{path, name, isDir, size, modTime}]`. |
| `--abs` | | boolean | `false` | Emit absolute filesystem paths instead of workspace-relative paths. |
| `--rel` | | boolean | `true` | Emit forward-slash normalized relative paths (default). |

## Examples

### 1. Immediate Child Listing (ls / dir Replacement)

List all immediate child files and directories in the current working directory:

```bash
$ gitmap child-path
```

### 2. Replacing PowerShell `Get-ChildItem -Filter`

Instead of running slow PowerShell cmdlets:

```powershell
# Slow PowerShell equivalent:
Get-ChildItem -Path cli -Filter *.go
```

Use fast native Go child-path:

```bash
$ gitmap child-path cli "*.go"
```

For recursive file-only searches replacing `Get-ChildItem -Recurse -File -Depth 2`:

```powershell
# Slow PowerShell equivalent:
Get-ChildItem -Path . -Filter *.ts -Recurse -Depth 2 -File
```

Run:

```bash
$ gitmap child-path . "*.ts" --depth 2 --type f
```

### 3. Replacing `git grep` on Untracked or Scratch Files

When inspecting new source files or scratch directories before they are tracked by git:

```bash
# Searches all Go files up to 3 levels deep, including untracked and .ai-memory files:
$ gitmap child-path . "*.go" --depth 3 --type f
```

### 4. Extension Filtering

Filter by multiple file extensions across subdirectories:

```bash
$ gitmap child-path cli -e go,md --depth 2
```

### 5. Directory-Only Exploration

Find only immediate subdirectories in `cli`:

```bash
$ gitmap child-path cli --type d
```

### 6. Bounded Result Set

Limit output to the first 5 matches:

```bash
$ gitmap child-path cli "*.go" --limit 5
```

### 7. Structured JSON Output for AI Agents

Emit machine-readable JSON array for programmatic consumption:

```bash
$ gitmap child-path cli "*.go" --depth 1 --json --limit 3
```

Output:
```json
[
  {
    "path": "cli/locate_entry.go",
    "name": "locate_entry.go",
    "isDir": false,
    "size": 203,
    "modTime": "2026-10-08T16:02:13+08:00"
  },
  {
    "path": "cli/main.go",
    "name": "main.go",
    "isDir": false,
    "size": 1420,
    "modTime": "2026-10-08T15:40:00+08:00"
  },
  {
    "path": "cli/root.go",
    "name": "root.go",
    "isDir": false,
    "size": 25803,
    "modTime": "2026-10-08T16:08:13+08:00"
  }
]
```

### 8. Automation Subsystem Invocation

Run under the `gitmap aum` command tree:

```bash
$ gitmap aum child-path . "*.go" --depth 2 --type f
$ gitmap aum cpth cli -e go -l 10
```

## See Also

- [find-files](find-files.md) — Exact filename lookup across tracked trees
- [find](find.md) — Universal repository glob search
- [automation](automation.md) — High-performance native Go automation suite
