# Profiles — Multiple Database Environments

## Overview

The `profile` command lets users maintain separate database environments
(e.g., work, personal, client projects) and switch between them. Each
profile has its own SQLite database file while sharing the same
`.gitmap/output/` directory.

---

## How It Works

- Profiles are tracked in `.gitmap/output/data/profiles.json`
- Each profile maps to a separate DB file:
  - `default` → `gitmap.db` (backward compatible)
  - `work` → `gitmap-work.db`
  - `personal` → `gitmap-personal.db`
- Switching profiles changes which DB file all commands use
- The `default` profile always exists and cannot be deleted

---

## Commands

### `gitmap profile` (alias: `pf`)

Manage database profiles.

**Subcommands:**

#### `gitmap profile create <name>`

Create a new profile with its own empty database.

```bash
gitmap profile create work
gitmap pf create personal
```

#### `gitmap profile list`

Show all profiles with active marker.

```bash
gitmap profile list
gitmap pf list
```

Output:
```
PROFILE              STATUS
default              (active)
work
personal
```

#### `gitmap profile switch <name>`

Switch to a different profile.

```bash
gitmap profile switch work
gitmap pf switch personal
```

#### `gitmap profile show`

Display the currently active profile.

```bash
gitmap profile show
gitmap pf show
```

#### `gitmap profile delete <name>`

Delete a profile and its database file.

```bash
gitmap profile delete personal
```

Cannot delete the `default` profile or the currently active profile.

---

## Configuration File

`.gitmap/output/data/profiles.json`:

```json
{
  "active": "default",
  "profiles": ["default", "work", "personal"]
}
```

If the file doesn't exist, gitmap assumes `default` profile
(backward compatible with existing installations).

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_profile.go` | Command names, messages |
| `model/profile.go` | ProfileConfig struct |
| `store/profile.go` | Profile config read/write, DB file resolution |
| `cmd/profile.go` | Profile command routing |
| `cmd/profileops.go` | Create, list, switch, delete, show handlers |
| `cmd/profileutil.go` | Shared profile helper functions |

---

## Constraints

- Default profile always exists and cannot be deleted.
- Cannot delete the currently active profile.
- Profile names must be unique.
- Backward compatible: no `profiles.json` = `default` profile.
- All files under 200 lines, all functions 8–15 lines.
