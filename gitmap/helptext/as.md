# gitmap as

Register a short name (alias) for the current Git repository so you can
target it from any directory on disk.

## Aliases

`s-alias`

## Usage

    gitmap as [alias-name] [--force]

When `alias-name` is omitted, the basename of the repo's top-level folder
is used. The command must be run from inside the working tree of a Git
repository.

## What it does

1. Resolves the repo root via `git rev-parse --show-toplevel`.
2. Builds a `ScanRecord` for the repo (slug, HTTPS/SSH URLs, branch, paths)
   and upserts it into the gitmap SQLite database — equivalent to running
   `gitmap scan` for that one repo.
3. Creates (or updates) an alias row mapping `alias-name -> repo`.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--force` / `-f` | false | Overwrite an existing alias that points to a different repo. |

## Examples

    cd ~/code/my-app
    gitmap as              # registers alias "my-app"
    gitmap as backend      # registers alias "backend"
    gitmap as backend -f   # overwrite an existing "backend" alias

## Errors

| Condition | Exit | Message |
|-----------|------|---------|
| Not inside a Git repo | 1 | `error: not inside a Git repository ...` |
| Alias already in use, no `--force` | 1 | `error: alias '...' is already mapped to a different repo ...` |

## See also

- `gitmap alias list` — show every registered alias.
- `gitmap release-alias` — release a repo by its alias from anywhere.
- `gitmap release-alias-pull` — pull-then-release shortcut.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter as
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
