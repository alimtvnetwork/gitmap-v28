# gitmap profile

Manage database profiles (separate repo databases for different contexts).

## Alias

pf

## Usage

    gitmap profile <create|list|switch|delete|show> [name]

## Flags

None.

## Prerequisites

- None

## Examples

### Example 1: Create a new profile and switch to it

    gitmap profile create work
    gitmap profile switch work

**Output:**

    ✓ Profile 'work' created (empty database)

    ✓ Switched to profile 'work'
    Active profile: work (0 repos)
    → Run 'gitmap scan' to populate this profile

### Example 2: List all profiles

    gitmap pf list

**Output:**

    PROFILE     REPOS   GROUPS  STATUS
    default     42      3       
    work        18      2       ✓ active
    personal    7       1       
    3 profiles

### Example 3: Show current profile details

    gitmap profile show

**Output:**

    Active profile: work
    Repos:    18
    Groups:   2 (backend, frontend)
    Aliases:  5
    Created:  2025-03-01
    Database: ~/.gitmap/profiles/work.db

### Example 4: Delete a profile

    gitmap profile delete old-project

**Output:**

    Delete profile 'old-project' and all its data? [y/N]: y
    ✓ Profile 'old-project' deleted (12 repos, 1 group removed)

## See Also

- [diff-profiles](diff-profiles.md) — Compare repos across profiles
- [export](export.md) — Export current profile data
- [import](import.md) — Import data into a profile
- [db-reset](db-reset.md) — Reset the current profile database

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter profile
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
