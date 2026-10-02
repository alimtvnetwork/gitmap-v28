# gitmap devtools / gitmap devtool

Developer Tools Cache Explorer & Purger -- scans, renders hierarchical disk trees, and reclaims space from build caches across Go, Node, Python, Rust, .NET, and Java ecosystems.

## Usage

```bash
gitmap devtools [flags]
gitmap devtool [flags]
gitmap clean-dev [flags]
gitmap clear devtools [flags]
```

## Aliases

- `devtool`, `devtools`, `clean-dev`, `cleandev`, `dev-clean`, `clear devtools`

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--tree` | `-t` | Render hierarchical tree view of discovered caches, file counts, and reclaimable sizes |
| `--dry-run` | `-n`, `-d` | Preview discoverable cache paths without deleting anything |
| `--force` | `-f` | Invalidate SQLite Split-DB cache (`.gitmap/data/devtools/cache/sql.db`) and force a fresh disk discovery scan |
| `--only <cats>` | | Limit cache discovery/cleanup to specific ecosystems (e.g. `--only go,pnpm,npm`) |
| `--yes` | `-y` | Execute deletion non-interactively without prompt |
| `--json` | `-j` | Output discovered cache paths and size metrics in JSON format |
| `--help` | `-h` | Display usage instructions |

## Examples

### 1. Render Hierarchical Cache Tree (Dry-Run)

```bash
gitmap devtools --tree --dry-run
```

Output:

```text
  Developer Tool Cache Hierarchy
  ==============================
  ├── Go
  │   ├── Build Cache (%LOCALAPPDATA%\go-build) [3.20 GB, 12,410 files]
  │   └── Module Cache (%USERPROFILE%\go\pkg\mod) [1.80 GB, 8,920 files]
  ├── Node
  │   ├── pnpm Store (C:\dev-tool\pnpm\store) [2.10 GB, 6,100 files]
  │   └── npm Cache (%LOCALAPPDATA%\npm-cache) [450.00 MB, 1,200 files]
  └── Python
      └── pip Cache (%LOCALAPPDATA%\pip\cache) [310.00 MB, 412 files]

  ✔ Total Reclaimable: 7.86 GB across 29,042 files
```

### 2. Purge Only Go Build Cache Non-Interactively

```bash
gitmap clear devtools --only go -y
```

### 3. Force Fresh Scan Bypassing Split-DB Cache

```bash
gitmap devtools --force --dry-run
```
