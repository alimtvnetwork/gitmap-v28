# gitmap push

Run `git push` in the current repository, optionally rewriting the
`origin` remote to SSH or HTTPS first.

## Alias

ph

## Usage

    gitmap push [--ssh|--https] [git push args...]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --ssh, -ssh, --sh | false | Rewrite `remote.origin.url` to SSH (`git@host:owner/repo.git`) and persist via `git remote set-url` before pushing |
| --https, -https, --ht | false | Rewrite `remote.origin.url` to HTTPS (`https://host/owner/repo.git`) and persist before pushing |

`--ssh` and `--https` are mutually exclusive; if both are set, `--ssh`
wins and a one-line warning is printed to stderr.

## Prerequisites

- Run inside a git repository.

## Examples

### Example 1: Plain push

    gitmap push

**Output:**

    → Running: git push (cwd: /repos/my-api)
    Everything up-to-date

### Example 2: Convert origin to SSH and push

    gitmap push --ssh

**Output:**

    → Rewrote remote.origin.url → git@github.com:owner/my-api.git
    → Running: git push (cwd: /repos/my-api)
    To github.com:owner/my-api.git
       a1b2c3d..e4f5g6h  main -> main

### Example 3: Forward extra args to git

    gitmap push --https origin main

**Output:**

    → Rewrote remote.origin.url → https://github.com/owner/my-api.git
    → Running: git push origin main (cwd: /repos/my-api)
    Everything up-to-date

## See Also

- [pull](pull.md) — Pull from origin with the same `--ssh`/`--https` flags
- [clone](clone.md) — Clone with `--ssh` / `--https` transport coercion

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter push
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
