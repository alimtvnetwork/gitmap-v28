# gitmap release-alias-pull

Pull-then-release shortcut for an aliased repository. Equivalent to running
`gitmap release-alias <alias> <version> --pull`.

## Aliases

`rap`

## Usage

    gitmap release-alias-pull <alias> <version> [--no-stash] [--dry-run]

## What it does

Identical to `gitmap release-alias` with `--pull` always enabled:

1. Resolves `<alias>` to an absolute repo path.
2. Runs `git pull --ff-only` inside that repo.
3. Auto-stashes dirty changes (unless `--no-stash`).
4. `chdir`s into the repo and runs the standard `gitmap release` workflow
   with `<version>`.
5. Pops the auto-stash on exit.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--no-stash` | false | Abort if dirty (skip auto-stash). |
| `--dry-run` | false | Forwarded to `gitmap release` — preview only. |

## Examples

    gitmap release-alias-pull my-api v1.4.0
    gitmap rap backend v0.9.0 --dry-run

## See also

- `gitmap release-alias` — same command without forced `--pull`.
- `gitmap as` — register an alias for the current repo.
- `gitmap release` — the underlying release workflow.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter release-alias-pull
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
