# cat

Display file contents, line ranges, or search filters cleanly in the terminal across all operating systems.

## Aliases

`view`, `type`

## Usage

    gitmap cat <filepath> [flags]

## Flags

- `-L, --lines <start:end>`: Slices lines from 1-based start line to end line (e.g. `-L 140:220`, `-L 100-150`).
- `-n, --head <count>`: Shows the first N lines (e.g. `-n 30`, `--head 30`).
- `--tail <count>`: Shows the last N lines (e.g. `--tail 30`).
- `-num, --number, --numbers`: Displays 1-based line numbers in the gutter (`  140 | ...`).
- `-q, --query, --filter, -s <string>`: Filters lines matching query string or multiple pipe-delimited tokens (e.g. `-q "foo|bar|baz"`).
- `-l, --limit <count>`: Maximum matching lines to display (default: 50 when filtering).
- `-i, --info, --stats`: Displays file metadata (line count, byte size, format) without dumping full contents.
- `--raw`: Emits content directly to stdout without Tokyo banners or borders.

## Description

The `cat` command provides an ultra-fast, cross-platform file reader and line inspector. It eliminates the need for ad-hoc Python one-liners (`open(...).readlines()`) or PowerShell `Get-Content` commands. It expands paths, environment variables, tildes, markdown URIs, and brackets, and automatically presents safe image summary cards for raster assets.

## Examples

```bash
# View contents of a configuration file
gitmap cat .env

# Slices lines 140 to 220 with formatted line numbers
gitmap cat src-tauri/src/modules/repo_db/dispatch.rs -L 140:220 --num

# Search for "ban" in skill.md and limit to 15 matches
gitmap cat .agents/skills/execute-parent-task-with-n-steps-v7/skill.md -q "ban" -l 15

# Search multiple keywords (OR matching) and limit to 15 matches
gitmap cat .agents/skills/execute-parent-task-with-n-steps-v7/skill.md -q "Get-ChildItem|gitmap find|ban|search" -l 15

# View file length, line count, and byte size metadata
gitmap cat .agents/skills/execute-parent-task-with-n-steps-v7/skill.md --info

# View first 30 lines of a file
gitmap cat cli/cmd/root.go --head 30

# View last 20 lines of a file
gitmap cat go.mod --tail 20

# View raw file content without decorative Tokyo banners
gitmap cat config.json --raw
```
