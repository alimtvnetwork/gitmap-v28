# Bookmarks — Save & Replay Commands

## Overview

The `bookmark` command lets users save frequently-used command+flag
combinations and replay them by name. Bookmarks are stored in the
`Bookmarks` SQLite table.

---

## Table Schema

### Bookmarks Table

| Column    | Type | Constraints               | Notes                       |
|-----------|------|---------------------------|-----------------------------|
| Id        | TEXT | PRIMARY KEY               | Timestamp-based unique ID   |
| Name      | TEXT | NOT NULL UNIQUE           | User-chosen bookmark name   |
| Command   | TEXT | NOT NULL                  | Command name (e.g. `scan`)  |
| Args      | TEXT | DEFAULT ''                | Positional arguments        |
| Flags     | TEXT | DEFAULT ''                | Flags (e.g. `--mode ssh`)   |
| CreatedAt | TEXT | DEFAULT CURRENT_TIMESTAMP |                             |

---

## Commands

### `gitmap bookmark` (alias: `bk`)

Manage saved command bookmarks.

**Subcommands:**

#### `gitmap bookmark save <name> <command> [args...] [--flags...]`

Save a command+flags combination under a name.

```bash
gitmap bookmark save ssh-scan scan --mode ssh
gitmap bk save quick-status status
gitmap bk save scan-projects scan ./projects --mode ssh --open
```

#### `gitmap bookmark list [--json]`

Show all saved bookmarks.

```bash
gitmap bookmark list
gitmap bk list --json
```

#### `gitmap bookmark run <name>`

Replay a saved bookmark.

```bash
gitmap bookmark run ssh-scan
gitmap bk run quick-status
```

#### `gitmap bookmark delete <name>`

Remove a saved bookmark.

```bash
gitmap bookmark delete ssh-scan
gitmap bk delete quick-status
```

---

## Replay Behavior

When `bookmark run <name>` is executed:

1. The bookmark record is loaded from the database.
2. `os.Args` is reconstructed from the saved command, args, and flags.
3. The standard `dispatch()` function is called, which also triggers
   the audit hook — so replayed commands appear in `gitmap history`.

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_bookmark.go` | SQL, command names, messages |
| `model/bookmark.go` | BookmarkRecord struct |
| `store/bookmark.go` | Bookmark CRUD operations |
| `cmd/bookmark.go` | Bookmark command routing |
| `cmd/bookmarksave.go` | Save subcommand |
| `cmd/bookmarklist.go` | List and delete subcommands |
| `cmd/bookmarkrun.go` | Run (replay) subcommand |

---

## Constraints

- Bookmark names must be unique (enforced by UNIQUE constraint).
- Save refuses if name exists — user must delete first.
- `db-reset --confirm` also clears the Bookmarks table.
- PascalCase table and column names.
- All files under 200 lines, all functions 8–15 lines.
