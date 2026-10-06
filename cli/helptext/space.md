# gitmap space

Space operations namespace. Groups repo-setup operations that apply
GitMap's curated shared baselines, starting with `common`.

## Subcommands

| Subcommand | Purpose |
|------------|---------|
| `common` | Apply the curated common baselines — same logic as `gitmap commons` |

## Usage

    gitmap space common [--dry-run] [--force]

## `space common`

Adds or dedupe-merges the curated baselines for `.gitignore`,
`.gitattributes`, `.prettierignore`, `.prettierrc`, and runs
`git lfs install --local` + the `lfs/common` `.gitattributes`
block, all in one pass. Identical engine to `gitmap commons`.

Behavior:
  - Line-based targets append MISSING lines only; existing entries
    are preserved verbatim. Safe to re-run.
  - `.prettierrc` is JSON key-union: missing keys added, existing kept
    unless `--force` is passed.
  - Idempotent — a second run on an unchanged repo writes nothing.

## Flags

| Flag | Purpose |
|------|---------|
| `--dry-run`, `-n` | Print planned additions without touching disk |
| `--force`, `-f` | For `.prettierrc` only: overwrite conflicting keys instead of preserving them |

## Examples

```bash
gitmap space common
gitmap space common --dry-run
gitmap space common --force
```

## See Also

- [commons](commons.md) — Same baselines as a top-level shortcut (`co`)
- [sync](sync.md) — Per-target union merge with the same engine
