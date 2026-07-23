# gitmap diff

Read-only preview of what `gitmap merge-both / merge-left /
merge-right` would change between two folders. Lists files
present on only one side and files whose content differs on both
sides. Writes nothing, commits nothing, pushes nothing.

Spec: companion to `spec/01-app/97-move-and-merge.md`

## Alias

df

## Usage

    gitmap diff LEFT RIGHT [flags]
    gitmap df   LEFT RIGHT [flags]

LEFT and RIGHT must both be local folder paths. URL endpoints are
intentionally rejected — clone them first with `gitmap clone` so
`diff` stays strictly side-effect-free.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Emit a JSON object `{summary, entries}` instead of text |
| --only-conflicts | false | Show only files that differ on both sides |
| --only-missing | false | Show only files present on one side |
| --include-identical | false | Include byte-equal files in the output |
| --include-vcs | false | Walk `.git/` (default: skipped) |
| --include-node-modules | false | Walk `node_modules/` (default: skipped) |

## Prerequisites

None. Both endpoints are folders that already exist on disk.

## Examples

### Example 1: Plain diff between two local folders

    gitmap diff ./gitmap-v28 ./gitmap-v28

**Output:**

      Conflicts (different content on both sides):
        README.md  (L: 4.2 KB @ 2026-04-17 14:02 | R: 4.1 KB @ 2026-04-18 09:11)
        src/app.ts  (L: 2.0 KB @ 2026-04-16 11:00 | R: 2.3 KB @ 2026-04-18 09:55)

      Missing on RIGHT (would be added by merge-right / merge-both):
        docs/changelog.md  (L: 1.1 KB @ 2026-04-15 08:30)

      Missing on LEFT (would be added by merge-left / merge-both):
        scripts/build.sh  (R: 512 B @ 2026-04-17 22:45)

    [diff] summary: 1 missing-on-left, 1 missing-on-right, 2 conflicts, 137 identical

### Example 2: Conflicts only (preview before merge-both)

    gitmap diff ./gitmap-v28 ./gitmap-v28 --only-conflicts

**Output:**

      Conflicts (different content on both sides):
        README.md  (L: 4.2 KB @ 2026-04-17 14:02 | R: 4.1 KB @ 2026-04-18 09:11)
        src/app.ts  (L: 2.0 KB @ 2026-04-16 11:00 | R: 2.3 KB @ 2026-04-18 09:55)

    [diff] summary: 0 missing-on-left, 0 missing-on-right, 2 conflicts, 137 identical

### Example 3: Machine-readable output

    gitmap df ./gitmap-v28 ./gitmap-v28 --json

**Output:**

    {
      "summary": {
        "missing_left": 1,
        "missing_right": 1,
        "conflicts": 2,
        "identical": 137
      },
      "entries": [
        { "path": "README.md", "kind": "conflict",
          "left_size": 4288, "right_size": 4203,
          "left_mtime": 1745935320, "right_mtime": 1745994660 }
      ]
    }

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Diff produced (regardless of whether differences were found) |
| 1 | One endpoint missing, not a directory, or walk failed |
| 2 | Wrong number of positional arguments (need exactly LEFT RIGHT) |

## Notes

- `gitmap diff` is the recommended dry-run preview before
  `gitmap merge-both` — every conflict listed here will trigger
  the `[L]eft / [R]ight / [S]kip / [A]ll-left / [B]all-right /
  [Q]uit` prompt during merge-both.
- The same default ignore list as `merge-*` applies: `.git/`,
  `node_modules/`, and `.gitmap/release-assets/` are skipped
  unless the corresponding `--include-*` flag is set.
- Identical files are tallied in the summary but not listed by
  default (use `--include-identical` to dump them).

## See Also

- [merge-both](merge-both.md) — Apply a two-way merge after previewing
- [merge-left](merge-left.md) — Apply RIGHT's changes into LEFT
- [merge-right](merge-right.md) — Apply LEFT's changes into RIGHT
- [mv](mv.md) — Move LEFT into RIGHT and delete LEFT

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter diff
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
