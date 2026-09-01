# Watch — Live Repository Monitor

## Overview

The `watch` command provides a live-refreshing dashboard that monitors
all tracked repositories for local and remote changes. It periodically
checks working tree status and remote ahead/behind counts.

---

## How It Works

1. Load all tracked repos from the database (or filter by `--group`).
2. For each repo, run `git status --porcelain` and `git rev-list`
   counts against the upstream tracking branch.
3. Display a formatted table to the terminal.
4. Sleep for the configured interval (default 30 seconds).
5. Clear the screen and repeat until the user presses Ctrl+C.

---

## Commands

### `gitmap watch` (alias: `w`)

Start the live monitoring dashboard.

```bash
gitmap watch
gitmap w
gitmap watch --interval 60
gitmap watch --group work
```

---

## Dashboard Output

```
gitmap watch — refreshing every 30s (Ctrl+C to stop)
Last updated: 2026-03-09 14:32:05 UTC

REPO                STATUS     BRANCH          AHEAD  BEHIND  STASH
api-gateway         dirty      main            0      3       0
frontend-app        clean      feature/nav     2      0       1
shared-lib          dirty      develop         0      0       0
docs-site           clean      main            0      0       0

Repos: 4 | Dirty: 2 | Behind: 1 | Stash: 1
```

### Status Values

| Status | Meaning |
|--------|---------|
| `clean` | No uncommitted changes |
| `dirty` | Uncommitted changes in working tree or index |
| `error` | Git command failed (e.g., missing remote) |

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval <seconds>` | `30` | Refresh interval in seconds (minimum 5) |
| `--group <name>` | — | Monitor only repos in a specific group |
| `--no-fetch` | `false` | Skip `git fetch`; use local refs only |
| `--json` | `false` | Output single snapshot as JSON and exit (no loop) |

### JSON Mode

With `--json`, watch runs once, outputs structured JSON, and exits.
Useful for scripting and external dashboards.

```json
{
  "timestamp": "2026-03-09T14:32:05Z",
  "repos": [
    {
      "name": "api-gateway",
      "path": "C:\\repos\\api-gateway",
      "branch": "main",
      "status": "dirty",
      "ahead": 0,
      "behind": 3,
      "stash": 0
    }
  ],
  "summary": {
    "total": 4,
    "dirty": 2,
    "behind": 1,
    "stash": 1
  }
}
```

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_watch.go` | Command names, column headers, messages |
| `cmd/watch.go` | Command entry, flag parsing, refresh loop |
| `cmd/watchops.go` | Per-repo status collection logic |
| `cmd/watchformat.go` | Table formatting and summary line |
| `gitutil/watchstatus.go` | Git status and ahead/behind queries |

---

## Constraints

- Minimum interval is 5 seconds to avoid hammering disk/network.
- `git fetch` runs once per cycle, not per repo, using `--all`.
- Terminal clear uses ANSI escape `\033[2J\033[H` (cross-platform).
- All files under 200 lines, all functions 8–15 lines.
- Graceful shutdown on SIGINT (Ctrl+C) with a brief summary.
