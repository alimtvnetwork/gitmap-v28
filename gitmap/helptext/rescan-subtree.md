# gitmap rescan-subtree

Re-run `gitmap scan` against a single subtree -- typically the
`absolutePath` of an at-cap row from a previous scan -- with a deeper
default `--max-depth` so the previously hidden nested repos surface in
one step.

## Alias

rss

## Usage

    gitmap rescan-subtree <absolutePath> [scan flags...]
    gitmap rss <absolutePath> [scan flags...]

## When to use this

`gitmap scan` defaults to `--max-depth 4`. Rows in the resulting
CSV/JSON whose `depth` column equals the cap (`4` by default) are
boundary discoveries: the walker found the repo right at the edge and
did **not** descend into its children, so any nested repos below were
silently skipped (see `gitmap help scan`, section "Depth column and
`--max-depth`").

The recommended workflow is:

1. `gitmap scan ~/work` -- produces `.gitmap/output/gitmap.csv`.
2. Find at-cap rows:

       awk -F, 'NR>1 && $NF==4' .gitmap/output/gitmap.csv

3. Copy the `absolutePath` column from one of those rows.
4. `gitmap rescan-subtree <that path>` -- this command. Done.

## Behavior

- Validates `<absolutePath>` exists and is a directory. Bad input exits
  with code `2` and a `Did you copy the absolutePath from a row that
  has since moved?` hint when the directory is missing.
- Resolves the path with `filepath.Abs` so a relative path also works.
- Forwards every other flag verbatim to `gitmap scan`. Anything
  documented in `gitmap help scan` is accepted here.
- If you do **not** pass `--max-depth`, this command injects
  `--max-depth 8` before invoking the scan. That is intentionally
  deeper than the scan default (4) so the typical "go deeper than the
  cap" case finishes in one shot. Pass an explicit `--max-depth N` to
  override (use `-1` for unlimited).
- The scan target is the supplied subtree, **not** the original wide
  root -- so the resulting `.gitmap/output/` artifacts describe just
  that subtree.

## Examples

### Example 1: typical at-cap rescan

    gitmap scan ~/work                           # default --max-depth 4
    awk -F, 'NR>1 && $NF==4' .gitmap/output/gitmap.csv
    # ...says /home/me/work/monorepo is at the cap...
    gitmap rescan-subtree /home/me/work/monorepo # implicit --max-depth 8

**Output:**

      ▶ gitmap rescan-subtree — /home/me/work/monorepo (max-depth=8)
      ▶ gitmap scan vX.Y.Z — /home/me/work/monorepo
    ...
    ✓ Walked 412 directories · found 17 repositories
    ✓ .gitmap/output/gitmap.csv written

### Example 2: explicit unlimited depth

    gitmap rescan-subtree /home/me/work/monorepo --max-depth -1

The user-supplied `--max-depth` wins; no synthetic cap is added.

### Example 3: forward arbitrary scan flags

    gitmap rss /srv/repos --output json --quiet --no-probe

All flags after the path go straight to `gitmap scan`.

## Exit codes

- `0` -- the underlying scan ran (scan itself decides 0 vs. non-zero on
  hard walk errors).
- `2` -- usage error: missing path, more than one positional, or path
  is not an existing directory. Distinct from scan failures so shell
  wrappers can tell "you invoked me wrong" apart from "the walk
  itself failed".

## See Also

- [scan](scan.md) -- the underlying command, including the "Depth
  column and `--max-depth`" section that defines at-cap rows.
- [rescan](rescan.md) -- replay the **last** scan with its cached
  flags (no path argument).

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter rescan-subtree
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
