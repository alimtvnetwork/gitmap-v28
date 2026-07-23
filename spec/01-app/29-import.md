# Import — Restore Database from Backup

## Overview

The `import` command reads a `gitmap-export.json` file (produced by
`gitmap export`) and restores all data into the database using
merge/upsert semantics. Existing records are preserved; duplicates
are skipped or updated.

---

## Command

### `gitmap import` (alias: `im`)

Import database from a portable JSON export file.

**Synopsis:**

```
gitmap import [file] --confirm
```

**Arguments:**

| Argument | Description                          | Default              |
|----------|--------------------------------------|----------------------|
| `file`   | Input file path                      | `gitmap-export.json` |

**Flags:**

| Flag        | Description                                | Required |
|-------------|--------------------------------------------|----------|
| `--confirm` | Confirm the import (prevents accidents)    | Yes      |

**Examples:**

```bash
# Import from default file
gitmap import --confirm
gitmap im --confirm

# Import from custom path
gitmap import backup-2026-03.json --confirm
```

---

## Import Behavior

| Table          | Strategy                                    |
|----------------|---------------------------------------------|
| Repos          | Upsert by ID — updates if exists            |
| Groups         | INSERT OR IGNORE — skips if name exists      |
| GroupRepos     | Resolved by slug — links repos to groups    |
| Releases       | Upsert by Tag — updates if exists           |
| CommandHistory | INSERT OR IGNORE by ID — skips duplicates   |
| Bookmarks      | INSERT OR IGNORE by ID — skips duplicates   |

Group members are linked by resolving `repoSlugs` against the Repos
table. If a slug doesn't exist (e.g., repo not imported), the link
is silently skipped.

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_import.go` | Command names, messages |
| `store/import.go` | ImportAll with per-table restore methods |
| `cmd/importcmd.go` | Import command handler |

---

## Constraints

- `--confirm` required to prevent accidental data overwrites.
- Merge semantics: never deletes existing data, only adds/updates.
- All files under 200 lines, all functions 8–15 lines.
