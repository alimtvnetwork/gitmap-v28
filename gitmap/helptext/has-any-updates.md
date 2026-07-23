# gitmap has-any-updates

Check if the current repository has new commits on the remote that you haven't pulled yet.

## Alias

hau, hac (has-any-changes)

## Usage

    gitmap has-any-updates
    gitmap hau
    gitmap hac

## Prerequisites

- Must be run inside a Git repository
- Remote must be configured with an upstream tracking branch

## Examples

### Example 1: Remote has new commits

    gitmap hau

**Output:**

    Checking for updates...

    ✓ Yes, you have 3 new update(s) from remote.
      Run 'git pull' to sync.

### Example 2: Already up to date

    gitmap hau

**Output:**

    Checking for updates...

    ✓ You are up to date. No new changes.

### Example 3: Local is ahead

    gitmap hac

**Output:**

    Checking for updates...

    ✓ You are 2 commit(s) ahead of remote. No incoming changes.

### Example 4: Branch has diverged

    gitmap hau

**Output:**

    Checking for updates...

    ⚠ Branch has diverged: 1 ahead, 4 behind remote.
      Run 'git pull --rebase' or 'git pull' to reconcile.

## See Also

- [status](status.md) — Show dirty/clean status for tracked repos
- [watch](watch.md) — Live-refresh dashboard
- [pull](pull.md) — Pull a specific repo

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter has-any-updates
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
