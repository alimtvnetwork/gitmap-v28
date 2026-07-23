# List Command Database Diagnostic

## Overview

The `gitmap-v27 ls` command should provide transparent feedback about
which database file it is querying, especially when no repos are found.

---

## Diagnostic Output

When `--verbose` is passed or when zero repos are found, the list
command should print the resolved database path:

```
  → Database: C:\Users\Alim\bin\data\gitmap.db
  No repos tracked. Run 'gitmap-v27 scan' first.
```

This helps users verify that scan and list are targeting the same
database file.

---

## Path Resolution Contract

All database access via `store.OpenDefault()` must resolve to:

```
<binary-directory>/data/<profile-db-file>
```

Where:
- `<binary-directory>` = `filepath.Dir(filepath.EvalSymlinks(os.Executable()))`
- `<profile-db-file>` = `gitmap.db` (default) or `gitmap-<profile>.db`

The profile config file (`profiles.json`) lives at:

```
<binary-directory>/data/profiles.json
```

**Not** `<binary-directory>/data/data/profiles.json` (the previous bug).

---

## Acceptance Criteria

- **Given** `gitmap-v27 ls --verbose` with repos, **then** DB path is printed
  above the repo table.
- **Given** `gitmap-v27 ls` with zero repos, **then** DB path is printed
  alongside the empty message.
- **Given** scan then list on the same binary, **then** both resolve
  to the identical database path.
- **Given** `ActiveProfileDBFile`, **then** profile config is read from
  `<binary-dir>/data/profiles.json` (no double nesting).
