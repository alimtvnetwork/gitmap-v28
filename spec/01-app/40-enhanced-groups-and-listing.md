# Enhanced Groups, Multi-Group, and List Filtering

## Overview

This spec adds: (1) group activation stored in SQLite, (2) multi-group
selection, (3) project-type filtering on `list`, (4) helper hint
footers on output, and (5) group-scoped pull/exec/status operations.

---

## 1. Active Group (`gitmap-v27 g`)

### Behavior

| Command | Effect |
|---------|--------|
| `gitmap-v27 g <name>` | Sets `<name>` as the active group (persisted in DB) |
| `gitmap-v27 g` | Prints the currently active group name |
| `gitmap-v27 g clear` | Clears the active group |

The active group is stored in a `Settings` table (key-value):

```sql
CREATE TABLE IF NOT EXISTS Settings (
    Key   TEXT PRIMARY KEY,
    Value TEXT NOT NULL
);
```

Key: `active_group`. Value: group name.

### Multi-Group: `gitmap-v27 multi-group` (alias: `mg`)

| Command | Effect |
|---------|--------|
| `gitmap-v27 mg g1,g2,g3` | Sets multiple active groups (comma-separated) |
| `gitmap-v27 mg` | Prints the currently active multi-group list |
| `gitmap-v27 mg clear` | Clears multi-group selection |
| `gitmap-v27 mg pull` | Pulls repos from all active multi-groups |
| `gitmap-v27 mg status` | Shows status for all active multi-group repos |
| `gitmap-v27 mg exec <args>` | Runs git command across all active multi-group repos |

Key: `active_multi_group`. Value: comma-separated group names.

### Scoped Operations via Active Group

When an active group or multi-group is set, commands that accept
`--group` can operate on it implicitly:

| Command | Behavior |
|---------|----------|
| `gitmap-v27 g pull` | Pull all repos in the active group |
| `gitmap-v27 g status` | Show status for active group repos |
| `gitmap-v27 g exec <args>` | Run git across active group repos |

---

## 2. List with Type Filter

### `gitmap-v27 ls <type>`

Accepts an optional positional type keyword to filter by project type:

| Keyword | Maps to |
|---------|---------|
| `go` | `ProjectKeyGo` |
| `node`, `nodejs` | `ProjectKeyNode` |
| `react` | `ProjectKeyReact` |
| `cpp` | `ProjectKeyCpp` |
| `csharp` | `ProjectKeyCSharp` |
| `groups` | Lists all groups (same as `gitmap-v27 group list`) |

When a type is provided, the output shows repos detected as that
project type with helper hints at the bottom.

### `gitmap-v27 ls groups`

Shows all defined groups (equivalent to `gitmap-v27 group list`).

---

## 3. Helper Hint Footers

After command output, print 2-3 contextual hints to stderr.
These help users discover related commands.

### When hints appear

| After command | Hints shown |
|--------------|-------------|
| `gitmap-v27 go-repos` / `gr` | CD navigation, grouping, listing |
| `gitmap-v27 ls` | Grouping, filtering by type, CD |
| `gitmap-v27 ls go` | Group add, CD, pull by group |
| `gitmap-v27 ls groups` | Group create, group show |
| `gitmap-v27 g` (active group) | Pull, status, exec, clear |
| `gitmap-v27 cd <name>` | Set default, list repos |
| `gitmap-v27 group list` | Group create, group show, delete |

### Format

```
Hints:
  → gitmap-v27 cd <repo-name>       Navigate to a repo
  → gitmap-v27 g create <name>      Create a group
  → gitmap-v27 ls go                List only Go projects
```

Hints are suppressed with `--quiet` or `-q`.

---

## 4. Data Model Changes

### Settings Table

```sql
CREATE TABLE IF NOT EXISTS Settings (
    Key   TEXT PRIMARY KEY,
    Value TEXT NOT NULL
);
```

### Store Methods

| Method | Description |
|--------|-------------|
| `GetSetting(key)` | Returns value or empty string |
| `SetSetting(key, value)` | Upserts a setting |
| `DeleteSetting(key)` | Removes a setting |

---

## 5. Package Structure

### New Files

| File | Contents |
|------|----------|
| `store/settings.go` | Settings CRUD |
| `cmd/multigroup.go` | `multi-group` command routing |
| `cmd/multigroupops.go` | Multi-group operations (pull/status/exec) |
| `cmd/hints.go` | Helper hint printing |
| `constants/constants_hints.go` | Hint message constants |
| `constants/constants_multigroup.go` | Multi-group constants |
| `constants/constants_settings.go` | Settings table SQL |

### Updated Files

| File | Change |
|------|--------|
| `store/store.go` | Add `SQLCreateSettings` to migration |
| `cmd/root.go` | Register `multi-group`/`mg` dispatch |
| `cmd/group.go` | Add activation and scoped operations |
| `cmd/list.go` | Add type filter support |
| `cmd/projectrepos.go` | Add hint footer |
| `constants/constants_cli.go` | New command names |
| `constants/constants_store.go` | Settings SQL |

---

## 6. Error Handling

| Scenario | Behavior |
|----------|----------|
| Active group not set | `"No active group. Use 'gitmap-v27 g <name>' to set one."` |
| Invalid type keyword | `"Unknown type: %s. Supported: go, node, react, cpp, csharp"` |
| Multi-group name not found | `"Group not found: %s"` |

---

## 7. Constraints

- All string literals in `constants` package.
- All files under 200 lines.
- All functions 8–15 lines.
- Positive conditions only.
- Blank line before `return`.
