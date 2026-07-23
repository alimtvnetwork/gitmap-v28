# gitmap commons

One-shot shortcut for `gitmap sync all`. Adds or dedupe-merges the
curated baselines for `.gitignore`, `.gitattributes`, `.prettierignore`,
`.prettierrc`, and runs `git lfs install --local` + the `lfs/common`
`.gitattributes` block, all in a single command.

Use this the moment you land in a fresh (or half-configured) repo and
want the standard hygiene files in place without thinking about which
targets to run.

## Alias

co

## Usage

    gitmap commons [--dry-run] [--force]

## Flags

| Flag | Purpose |
|------|---------|
| `--dry-run`, `-n` | Print planned additions without touching disk |
| `--force`, `-f` | For `.prettierrc` only: overwrite conflicting keys instead of preserving them |

## What it does

Runs the following targets in order (identical to `gitmap sync all`):

1. `.gitignore`      line-union merge (curated defaults)
2. `.gitattributes`  line-union merge (curated defaults)
3. `git lfs install --local` + `lfs/common` marker block
4. `.prettierignore` line-union merge (curated defaults)
5. `.prettierrc`     JSON key-union (existing keys win unless `--force`)

Existing lines/keys are preserved. Missing ones are appended below an
`# added by gitmap commons` sentinel so the delta is easy to spot in
`git diff`.

## Idempotency

Re-running on an unchanged repo prints `ok` for every target and
writes nothing. Safe in CI.

## Relationship to other commands

- `gitmap sync <target>` — same engine, but pick one target at a time.
- `gitmap add ignore` / `add attributes` — write gitmap-managed marker
  blocks instead of a raw union. Use when you want gitmap to own the
  file.
- `gitmap templates init` — scaffold from scratch when the file does
  not exist yet.

## Examples

```bash
gitmap commons
gitmap commons --dry-run
gitmap co -n
gitmap commons --force
```

## See Also

- [sync](sync.md) — Per-target union merge with the same engine
- [add ignore](add-ignore.md) — Marker-block managed .gitignore
- [add attributes](add-attributes.md) — Marker-block managed .gitattributes
- [templates init](templates-init.md) — Scaffold from scratch
