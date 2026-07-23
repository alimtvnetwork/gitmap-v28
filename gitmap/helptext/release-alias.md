# gitmap release-alias

Release a repository by its registered alias from any working directory.

## Aliases

`ra`

## Usage

    gitmap release-alias <alias> <version> [--pull] [--no-stash] [--dry-run]

## What it does

1. Resolves `<alias>` to an absolute repo path via the gitmap database.
2. (Optional, `--pull`) Runs `git pull --ff-only` inside that repo.
3. Auto-stashes any dirty changes (untracked included), unless `--no-stash`.
4. `chdir`s into the repo and invokes the standard `gitmap release` pipeline
   with `<version>`.
5. Pops the auto-stash on the way out (warning printed if pop fails).

This means you can release any tracked repo from anywhere — there is no
need to `cd` first.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--pull` | false | Run `git pull --ff-only` before releasing. |
| `--no-stash` | false | Abort if the working tree is dirty (skip auto-stash). |
| `--dry-run` | false | Forwarded to `gitmap release` — preview only. |

## Examples

    # Register first (one-time, from inside the repo)
    cd ~/code/my-api
    gitmap as my-api

    # Now release from anywhere
    gitmap release-alias my-api v1.4.0
    gitmap ra my-api v1.4.0 --pull
    gitmap ra my-api v1.4.0 --dry-run

## Sibling command

`gitmap release-alias-pull` (`rap`) is a thin verb that always implies
`--pull`. The two commands share the same code path:

    gitmap release-alias-pull my-api v1.4.0
    gitmap rap my-api v1.4.0

## Errors

| Condition | Exit | Message |
|-----------|------|---------|
| Alias not registered | 1 | `error: alias '...' is not registered. Run 'gitmap as ...' from the repo first.` |
| `chdir` to repo failed | 1 | `error: could not change directory to '...': ...` |
| `git pull` failed | 1 | `error: git pull failed in ...: ...` |
| `git stash` failed | 1 | `error: auto-stash failed in ...: ...` |

## See also

- `gitmap as` — register an alias for the current repo.
- `gitmap release` — the underlying release workflow.
- `gitmap alias list` — show every registered alias.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter release-alias
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
