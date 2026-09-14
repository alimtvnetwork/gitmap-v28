# gitmap profile

Manage database profiles (separate repo databases for different contexts) and tool installation profiles.

## Alias

pf

## Usage

    gitmap profile <create|list|switch|delete|show|install|import|export|inspect> [args]
    gitmap pf in <profile-name> [flags]

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| create | | Create a new database profile |
| list | ls, status | List all database profiles |
| switch | | Switch the active database profile |
| delete | | Delete a database profile |
| show | | Show details of the active profile |
| install | in | Execute an installation profile with idempotency tracking |
| import | cpi | Import Chrome browser profiles |
| export | cpe | Export Chrome browser profiles |
| inspect | check-import | Inspect Chrome profiles |

## Flags (for `profile install`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| --tree | -t | false | Preview full tool hierarchy of a profile without executing |
| --force | -f | false | Force reinstallation even if already recorded in SQLite |
| -y, --yes | -y | false | Auto-confirm package installation without prompting |

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

### Example 5: Import and inspect Chrome profiles

    gitmap profile import
    gitmap profile import erfan.office.n@gmail.com
    gitmap profile inspect ./chrome-ext --json

### Example 6: Install an installation profile

    gitmap profile install dev

**Output:**

    === Installing Profile: dev (Developer Core) ===
    Description: Essential developer toolchain
    Total Tools: 3

    ├─ git
    ├─ curl
    └─ ripgrep

      ✓ [1/3] git is already installed (2.44.0)
      → [2/3] Installing curl...
      → [3/3] Installing ripgrep...

    ✓ Profile 'dev' setup complete! (3/3 tools processed)

### Example 7: Idempotent profile execution (already installed)

When a profile is already installed and `--force` is omitted, GitMap detects the existing state from `installation.db`, displays the install timestamp, and prints the tool tree preview without reinstalling:

    gitmap profile install dev

**Output:**

    [INFO] Profile 'dev' is already installed (installed at: 2026-09-14T09:15:00Z)
    Profile: dev (Developer Core)
    Description: Essential developer toolchain
    Total Tools: 3

    ├─ git [installed]
    ├─ curl [installed]
    └─ ripgrep [installed]

## See Also

- [diff-profiles](diff-profiles.md) — Compare repos across profiles
- [export](export.md) — Export current profile data
- [import](import.md) — Import data into a profile
- [db-reset](db-reset.md) — Reset the current profile database
- [install](install.md) — Install individual developer tools

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter profile
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
