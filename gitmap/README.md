# GitMap

[![CI](https://github.com/alimtvnetwork/gitmap-v27/actions/workflows/ci.yml/badge.svg)](https://github.com/alimtvnetwork/gitmap-v27/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alimtvnetwork/gitmap-v27?style=flat-square)](https://goreportcard.com/report/github.com/alimtvnetwork/gitmap-v27)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/license-MIT-green?style=flat-square)](./LICENSE)

> Scan directories for Git repositories, generate clone instructions, and re-clone them anywhere.

## Quick Start

### Build

```powershell
# From the repo root:
.\run.ps1

# Skip git pull:
.\run.ps1 -NoPull

# Build only, no deploy:
.\run.ps1 -NoPull -NoDeploy

# Deploy to custom path:
.\run.ps1 -DeployPath "D:\tools"

# Build and run immediately:
.\run.ps1 -Run
.\run.ps1 -Run -RunPath "D:\projects"
.\run.ps1 -Run -RunArgs "--mode ssh"
```

The binary and `data/` config folder are output to `./bin/`. By default, the binary is also copied to the deploy path in `powershell.json` (default: `E:\bin-run`).

### Manual Build

```bash
cd gitmap
go build -o ../bin/gitmap.exe .
```

---

## Usage

### Scan a directory

Every scan **always produces all outputs** — terminal, CSV, JSON, and a folder structure Markdown file. They are written to a `gitmap-output/` folder at the root of the scanned directory.

```bash
# Scan current directory (outputs everything to ./gitmap-output/)
gitmap scan
gitmap s                    # shorthand

# Scan a specific folder with SSH URLs
gitmap scan ./projects --mode ssh

# Scan and add repos to GitHub Desktop
gitmap scan ./projects --github-desktop

# Scan and auto-open output folder
gitmap scan ./projects --open

# Custom output directory
gitmap scan ./projects --output-path ./my-exports
```

### Output files

When you run `gitmap scan ./projects`, the following is created:

```
projects/
└── gitmap-output/
    ├── gitmap.csv              # All repos in CSV format
    ├── gitmap.json             # All repos in JSON format
    └── folder-structure.md     # Tree view of repo hierarchy
```

The **folder-structure.md** shows a visual tree of all discovered repos:

```
# Folder Structure

Git repositories discovered by gitmap.

├── 📦 **my-app** (`main`) — https://github.com/user/my-app.git
├── libs/
│   ├── 📦 **core-lib** (`develop`) — https://github.com/user/core-lib.git
│   └── 📦 **utils** (`main`) — https://github.com/user/utils.git
└── 📦 **docs** (`main`) — https://github.com/user/docs.git
```

### Output path behavior

| Flag | Behavior |
|------|----------|
| No flags | Creates `gitmap-output/` inside the scanned directory |
| `--output-path ./exports` | Writes to `./exports/` |
| `--out-file report.csv` | Overrides CSV file path only |

### Clone from a previous scan

```bash
# Clone using shorthand (auto-resolves to ./gitmap-output/gitmap.json)
gitmap clone json
gitmap c json               # shorthand alias

# Clone using CSV shorthand
gitmap clone csv

# Clone from JSON with explicit path (preserves original folder structure)
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored

# Safe-pull existing clones (retries unlink/read-only failures)
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored --safe-pull

# Clone and add all repos to GitHub Desktop
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored --github-desktop
```

The clone command recreates the exact folder hierarchy from the `relativePath` field in each record. **Safe-pull is automatically enabled** when existing repos are detected in the target directory — it retries failed pulls, clears read-only attributes, and diagnoses Windows unlink issues. The `--safe-pull` flag can also be set explicitly. With `--github-desktop`, successfully cloned repos are automatically registered in GitHub Desktop.

---

## Configuration

### `data/config.json`

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": "./gitmap-output",
  "excludeDirs": [".cache", "node_modules", "vendor", ".venv"],
  "notes": ""
}
```

CLI flags override config values.

### `powershell.json`

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "gitmap.exe",
  "copyData": true
}
```

---

## CLI Reference

### `gitmap scan [dir]`

| Flag | Description | Default |
|------|-------------|---------|
| `--config <path>` | Config file path | `./data/config.json` |
| `--mode ssh\|https` | Clone URL style | `https` |
| `--output-path <dir>` | Output directory | `gitmap-output/` in scan dir |
| `--out-file <path>` | Exact CSV output file path | — |
| `--github-desktop` | Add discovered repos to GitHub Desktop | `false` |
| `--open` | Open output folder after scan | `false` |

### `gitmap clone <source|json|csv>`

**Shorthands:** `gitmap clone json` and `gitmap clone csv` auto-resolve to `./gitmap-output/gitmap.json` and `./gitmap-output/gitmap.csv`.

| Flag | Description | Default |
|------|-------------|---------|
| `--target-dir <path>` | Base clone directory | `.` |
| `--safe-pull` | Pull existing repos with retries, read-only clear, and diagnosis (auto-enabled) | `false` |
| `--github-desktop` | Add cloned repos to GitHub Desktop | `false` |
| `--verbose` | Write detailed debug log to a timestamped file | `false` |

**Note:** `--safe-pull` is automatically enabled when existing repos are detected in the target directory.

### `gitmap update [--verbose]`

| Flag | Description | Default |
|------|-------------|---------|
| `--verbose` | Write detailed debug log to a timestamped file | `false` |

---

## CSV Output Columns

`repoName, httpsUrl, sshUrl, branch, relativePath, absolutePath, cloneInstruction, notes`

---

## Project Structure

```
gitmap/
├── main.go              # Entry point
├── cmd/                  # CLI commands
│   ├── root.go           # Routing & flags
│   ├── scan.go           # Scan command
│   └── clone.go          # Clone command
├── config/               # Config loading
├── constants/            # All shared string literals
├── scanner/              # Directory walking
├── gitutil/              # Git command wrappers
├── mapper/               # Record building
├── formatter/            # Output (terminal, CSV, JSON, folder structure)
│   ├── terminal.go
│   ├── csv.go
│   ├── json.go
│   └── structure.go      # Folder tree Markdown
├── desktop/              # GitHub Desktop integration
├── cloner/               # Re-clone logic
├── model/                # Data structures
├── data/                 # Default config
│   └── config.json
├── powershell.json       # Build/deploy config
└── go.mod
```

## Command History

Every CLI command is automatically logged to the SQLite database. View the audit trail:

```bash
# Show recent history
gitmap history
gitmap hi

# Basic view (just command + time + status)
gitmap history --detail basic

# Detailed view filtered to scan commands
gitmap history --detail detailed --command scan

# Last 5 entries as JSON
gitmap history --json --limit 5

# Clear all history
gitmap history-reset --confirm
```

## Bookmarks

Save and replay frequently-used command+flag combinations:

```bash
# Save a bookmark
gitmap bookmark save ssh-scan scan --mode ssh
gitmap bk save quick-status status

# List all bookmarks
gitmap bookmark list
gitmap bk list --json

# Replay a saved bookmark
gitmap bookmark run ssh-scan
gitmap bk run quick-status

# Delete a bookmark
gitmap bookmark delete ssh-scan
```

## Usage Statistics

View aggregated command usage patterns:

```bash
# Show all command stats
gitmap stats
gitmap ss

# Stats for a specific command
gitmap stats --command scan

# JSON output
gitmap stats --json
```

## Database Export

Export the full database for backup or sharing:

```bash
# Export to default file (gitmap-export.json)
gitmap export
gitmap ex

# Export to custom path
gitmap export backup-2026-03.json
```

## Database Import

Restore a database from a backup file:

```bash
# Import from default file
gitmap import --confirm
gitmap im --confirm

# Import from custom path
gitmap import backup-2026-03.json --confirm
```

## Database Profiles

Manage multiple separate database environments:

```bash
# Create a new profile
gitmap profile create work
gitmap pf create personal

# List all profiles
gitmap profile list

# Switch active profile
gitmap profile switch work

# Show current profile
gitmap profile show

# Delete a profile
gitmap profile delete personal
```

## Specs

See [spec/01-app/](../spec/01-app/) for detailed specifications.

## License

Released under the [MIT License](./LICENSE). © 2026 Md. Alim Ul Karim.
