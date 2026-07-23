# gitmap branch

Manage branch checkout shortcuts. Subcommand-based.

## Alias

b

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| default    | def   | Checkout the repo's default branch (origin/HEAD or `main`) |

## Usage

    gitmap branch default
    gitmap b def

## Prerequisites

- Must be inside a Git repository

## How `default` is resolved

1. `git symbolic-ref refs/remotes/origin/HEAD` — uses the upstream's
   declared default (commonly `main`, sometimes `master` / `trunk` /
   `develop`).
2. Falls back to the built-in `main` if step 1 fails (e.g. no origin
   configured yet).

The `origin/` prefix is stripped automatically before `git checkout`,
so Git's DWIM rules create a local tracking branch when needed.

## Examples

### Example 1: Jump to the default branch

    gitmap b def

**Output:**

      ▶ Switching to default branch main...
    Switched to branch 'main'
    Your branch is up to date with 'origin/main'.

### Example 2: Equivalent long form

    gitmap branch default

## See Also

- [latest-branch](latest-branch.md) — `gitmap lb -s` jumps to the
  freshest branch instead of the default one.
- [status](status.md) — View repo branch and status info

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter branch
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
