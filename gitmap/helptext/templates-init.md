# gitmap templates init

Scaffold `.gitignore` and `.gitattributes` for one or more languages by
merging curated templates (overlay > embed) into the current directory.
Optionally include a Git LFS attributes block in the same pass.

This command is the **template-driven scaffolder**. It writes idempotent,
gitmap-managed marker blocks via the same primitive that powers
`gitmap add lfs-install`, so re-runs are byte-stable no-ops.

## Alias

ti

## Usage

    gitmap templates init <lang> [<lang>...] [--lfs] [--dry-run] [--force]

## Flags

| Flag | Purpose |
|------|---------|
| `--lfs` | Also merge `lfs/common.gitattributes` into `.gitattributes`. |
| `--dry-run` | Print every block that would be written; do not touch disk. |
| `--force` | Replace pre-existing `.gitignore` / `.gitattributes` outright. Discards hand edits OUTSIDE the gitmap marker block. |

Flag and positional order do not matter — `init --lfs go` and
`init go --lfs` are equivalent.

## What gets written

For each `<lang>`:

| Template | Required? | Target |
|----------|-----------|--------|
| `ignore/<lang>.gitignore` | **required** (hard error if missing) | `./.gitignore` |
| `attributes/<lang>.gitattributes` | optional (soft skip with notice) | `./.gitattributes` |

With `--lfs`:

| Template | Target |
|----------|--------|
| `lfs/common.gitattributes` | `./.gitattributes` (separate marker block) |

Each block is bracketed by stable markers:

    # >>> gitmap:<kind>/<lang> >>>
    ... template body ...
    # <<< gitmap:<kind>/<lang> <<<

Re-running `templates init` for the same languages updates the body in
place when the template has changed and is a complete no-op when it has
not. Hand edits OUTSIDE the markers survive untouched.

## Examples

### Example 1: Scaffold a Go project

    gitmap templates init go

Writes `.gitignore` (Go template) and `.gitattributes` (Go template)
into the current directory.

### Example 2: Multi-language project + LFS

    gitmap templates init go node --lfs

Writes:

- `.gitignore` with two blocks: `ignore/go` and `ignore/node`.
- `.gitattributes` with three blocks: `attributes/go`, `attributes/node`,
  and `lfs/common`.

### Example 3: Preview before writing

    gitmap templates init python --lfs --dry-run

Prints every block that would be written and exits without touching disk.
Useful for auditing the curated bytes against your project's conventions
before committing.

### Example 4: Reset a corrupted scaffold

    gitmap templates init go --force

Discards any pre-existing `.gitignore` / `.gitattributes` (including hand
edits outside the marker block) and re-creates them from the curated
template. Use with care — without `--force`, hand edits outside the
gitmap block always survive.

### Example 5: Use the short alias

    gitmap tpl ti rust --lfs

Round-trips identically to the long form.

## Relationship to `add lfs-install`

`templates init --lfs` and `gitmap add lfs-install` both write the
`lfs/common` block to `.gitattributes` using the **same marker tag**, so
they are interoperable: running one after the other is idempotent.

The difference:

- `add lfs-install` ALSO runs `git lfs install --local` (wires up the
  per-repo LFS hooks) and requires being inside a Git repo.
- `templates init --lfs` does NOT shell out to git; it is a pure
  template-merge operation and works before `git init`.

Use `templates init --lfs` for greenfield scaffolding, then run
`gitmap add lfs-install` once the repo is initialized to wire the hooks.

## Notes

- Operates on the **current working directory**. Does not require being
  inside a Git repository — scaffolding before `git init` is supported.
- Resolution is overlay-first: drop a file at
  `~/.gitmap/templates/ignore/<lang>.gitignore` to override the embedded
  template. See `gitmap templates list` for which entries are currently
  forked.
- Missing `attributes/<lang>` is a **soft skip** with a dim notice in the
  summary. Missing `ignore/<lang>` is a **hard error** — every language
  in the embed corpus has an ignore template, so a miss almost always
  means a typo.
- `--force` only touches the per-target files (`.gitignore`,
  `.gitattributes`). It does not delete other repo files.

## See Also

- [templates](templates.md) — `list` and `show` subcommands
- [add lfs-install](add-lfs-install.md) — Install LFS hooks + merge `lfs/common`
- [lfs-common](lfs-common.md) — Per-pattern `git lfs track` (no template)
- [setup](setup.md) — Configure Git global settings

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter templates-init
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
