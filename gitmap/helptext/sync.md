# gitmap sync

Union-merge curated defaults for `.gitignore`, `.gitattributes`,
`.prettierignore`, and `.prettierrc` into the current directory
**without** using gitmap marker blocks. Existing content stays
untouched; only missing entries are appended.

Use this when the target file is hand-maintained (or comes from a
framework generator) and you just want to top it up with any missing
essentials, rather than replace a marker block like `add ignore` does.

## Alias

sy

## Usage

    gitmap sync <target> [--dry-run] [--force]

## Targets

| Target | File / Action | Merge strategy |
|--------|---------------|----------------|
| `ignore` | `.gitignore` | Line union (trim-compare) |
| `attributes` | `.gitattributes` | Line union (trim-compare) |
| `lfs-install` | `git lfs install --local` + `lfs/common` block in `.gitattributes` | Marker-block merge (delegates to `add lfs-install`) |
| `prettier-ignore` | `.prettierignore` | Line union (trim-compare) |
| `prettier-rc` | `.prettierrc` | JSON key union (existing keys win) |
| `all` | all of the above | Runs each target in order |

## Flags

| Flag | Purpose |
|------|---------|
| `--dry-run`, `-n` | Print planned additions without touching disk |
| `--force`, `-f` | For `prettier-rc` only: overwrite conflicting keys instead of preserving them |

## What it does

1. Reads the target file (missing file counts as empty).
2. For line-based targets: splits both current + baseline on newlines,
   trims whitespace, and appends any baseline line that isn't already
   present. When the file already existed, comment/blank baseline
   lines are skipped so re-runs stay quiet.
3. For `.prettierrc`: unmarshals both, adds missing keys, and keeps
   existing values. Passing `--force` overwrites conflicting values.
4. Appends an `# added by gitmap sync` sentinel above the new lines so
   you can spot the delta in `git diff`.

## Idempotency

Re-running the same target on an unchanged repo prints `ok` and writes
nothing. Safe in CI.

## Examples

### Example 1: top up a hand-written .gitignore

    gitmap sync ignore

Adds any missing curated entries (node_modules, dist, .env, etc.)
without disturbing your existing rules.

### Example 2: preview all four in one shot

    gitmap sync all --dry-run

Prints planned additions for `.gitignore`, `.gitattributes`,
`.prettierignore`, and `.prettierrc` without writing.

### Example 3: reset prettier defaults

    gitmap sync prettier-rc --force

Sets `printWidth`, `semi`, `singleQuote`, `trailingComma` to the
curated values even if you had different ones.

## Relationship to `gitmap add ignore`

- `add ignore` writes a gitmap-managed **marker block** into
  `.gitignore` that is rewritten on every re-run.
- `sync ignore` does a **line-level union** into the raw file, with no
  markers, and never modifies existing lines.

Use `add ignore` when you want gitmap to own the ignore file. Use
`sync ignore` when the file is hand-maintained and you just want to
patch in essentials.

## See Also

- [add ignore](add-ignore.md) — Marker-block managed .gitignore
- [add attributes](add-attributes.md) — Marker-block managed .gitattributes
- [templates init](templates-init.md) — Scaffold from scratch

## Examples

```bash
gitmap sync ignore
gitmap sync attributes --dry-run
gitmap sync lfs-install
gitmap sync lfs-install --dry-run
gitmap sync prettier-rc --force
gitmap sync all
```
