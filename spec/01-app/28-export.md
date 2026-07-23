# Export — Full Database Backup

## Overview

The `export` command dumps the entire gitmap database into a single
portable JSON file for backup, sharing, or migration.

---

## Command

### `gitmap export` (alias: `ex`)

Export all database tables as a single JSON file.

**Synopsis:**

```
gitmap export [file]
```

**Arguments:**

| Argument | Description                          | Default              |
|----------|--------------------------------------|----------------------|
| `file`   | Output file path                     | `gitmap-export.json` |

**Examples:**

```bash
# Export to default file
gitmap export
gitmap ex

# Export to custom path
gitmap export backup-2026-03.json
```

---

## Export Format

```json
{
  "version": "2.23.0",
  "exportedAt": "2026-03-09T12:00:00Z",
  "repos": [ ... ],
  "groups": [
    {
      "id": "...",
      "name": "frontend",
      "description": "...",
      "color": "blue",
      "createdAt": "...",
      "repoSlugs": ["my-app", "ui-lib"]
    }
  ],
  "releases": [ ... ],
  "history": [ ... ],
  "bookmarks": [ ... ]
}
```

Groups include a `repoSlugs` array listing the slugs of member repos,
making the export self-contained and human-readable.

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_export.go` | Command names, messages, defaults |
| `model/export.go` | DatabaseExport and GroupExport structs |
| `store/export.go` | ExportAll aggregation method |
| `cmd/export.go` | Export command handler |

---

## Constraints

- Read-only operation — does not modify the database.
- Includes all tables: Repos, Groups (with member slugs), Releases,
  CommandHistory, and Bookmarks.
- All files under 200 lines, all functions 8–15 lines.
