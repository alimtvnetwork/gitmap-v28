# Spec 362: Component & CLI Specification

- **Module:** `cli/cmdmacro`, `cli/helpdoc`, `03-ai-scripts`
- **Spec Number:** 362
- **Document:** Component and CLI Contracts

---

## 1. CLI Commands & Flags

### 1.1 `gitmap cat` Flags and Syntax

```bash
gitmap cat <filepath> [flags]
```

#### Flags
- `-L, --lines <start:end>`: Slices lines from 1-based start line to end line (e.g. `-L 100:150`, `-L 140-220`).
- `-n, --head <count>`: Shows the first N lines (e.g. `-n 30`, `--head 30`).
- `--tail <count>`: Shows the last N lines (e.g. `--tail 30`).
- `-num, --number, --numbers`: Toggles line numbering in the output gutter (`  140 | ...`).
- `-q, --query, --filter, -s <string>`: Filters lines containing query string (case-insensitive).
- `-l, --limit <count>`: Caps maximum matching lines displayed (default: 50 when querying, unlimited otherwise).
- `--raw`: Emits content directly to `output.Raw()` without Tokyo banners or borders.

### 1.2 Examples

```bash
# View entire file cleanly
gitmap cat cli/cmdmacro/cat.go

# Slices lines 140 to 220 with line numbers
gitmap cat src-tauri/src/modules/repo_db/dispatch.rs -L 140:220 --num

# Search for "ban" in skill.md and limit to 15 matches
gitmap cat .agents/skills/execute-parent-task-with-n-steps-v7/skill.md -q "ban" -l 15

# View first 40 lines
gitmap cat cli/cmd/root.go --head 40

# View last 25 lines
gitmap cat go.mod --tail 25
```

---

## 2. Sync Script CLI Contracts

### 2.1 V6 Unrestricted Fleet Sync (`38-sync-prompts-skills-scripts.py`)

```bash
python 03-ai-scripts/38-sync-prompts-skills-scripts.py [--repo <name>] [--workers <N>] [--dry-run] [--force]
```
- Performs multi-branch backup, pre-change tag, asset mirror, commit, post-change tag, and push.
- Bypasses calendar constraints.

### 2.2 V7 Time-Restricted Fleet Sync (`38-sync-prompts-skills-scripts-v7.py`)

```bash
python 03-ai-scripts/38-sync-prompts-skills-scripts-v7.py [--tz <timezone>] [--schedule] [--auto-window] [--repo <name>] [--workers <N>]
```
- Validates active bot identity via `gitmap login --status` or `gitmap user info`.
- Verifies time window constraint: weekday morning < 08:00 AM in target timezone or automatically slots next morning window.
- Passes `--tz`, `--schedule`, and `--auto-window` to subagent git commands.

---

## 3. Go Internal Architecture

### 3.1 `cli/cmdmacro/cat_options.go`
- `CatOptions`: Struct holding parsed flags, paths, ranges, and filters.
- `ParseCatArgs(args []string) (CatOptions, error)`: Parses flags with zero nested ifs.
- `SliceLines(lines []string, opts CatOptions) []FilteredLine`: Extracts target window.
- `FilterLines(lines []string, query string, limit int) []FilteredLine`: Filters lines.
- `FormatCatOutput(lines []FilteredLine, opts CatOptions, totalLines int, fileSize int64) string`: Tokyo styled terminal renderer.
