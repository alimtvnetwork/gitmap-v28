# gitmap add ignore

Merge a curated `.gitignore` template into a gitmap-managed marker block
inside `./.gitignore`. Always merges the `common` template (OS junk, IDE
noise, gitmap artifacts) plus any languages you list on the command
line. Idempotent — re-running with the same language set is a
byte-stable no-op.

## Usage

    gitmap add ignore [langs...] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run | false | Preview the merged `.gitignore` block without writing anything |

## What it does

1. Verifies the current directory is inside a Git repository.
2. Resolves `common` plus each requested language (overlay > embed).
   Available languages: `go`, `node`, `python`, `rust`, `csharp`,
   `java`, `kotlin`, `php`, `ruby`, `swift` (plus anything you drop
   into `~/.gitmap/templates/ignore/`).
3. Concatenates bodies in stable order (`common` first, then your
   languages alphabetically), deduping repeated lines while preserving
   blank-line separators.
4. Writes a marker block at the repo root:

        # >>> gitmap:ignore/<sorted-langs> >>>
        ... merged body ...
        # <<< gitmap:ignore/<sorted-langs> <<<

5. Prints the merge outcome: `created`, `inserted block into`,
   `updated block in`, or `unchanged`.

## Idempotency contract

Running `gitmap add ignore go node` twice writes the **same** bytes the
second time. The marker tag is `ignore/go+node` regardless of whether
you typed `go node` or `node go` — the language list is sorted into
the tag so order doesn't fork the block.

Switching the language set (`add ignore go` → `add ignore go node`)
creates a *new* marker block under a *new* tag rather than rewriting
the old one. Run `gitmap add ignore go` again later and the old block
refreshes; the new one is left alone. If that's not what you want,
delete the stale block by hand before re-running.

Hand edits **outside** any marker block survive every re-run. Hand
edits **inside** a block are intentionally overwritten — fork the
template to `~/.gitmap/templates/ignore/<lang>.gitignore` if you want
custom content to stick.

## Examples

### Example 1: New Go repo

    cd ~/code/my-go-svc
    gitmap add ignore go

**Output:**

      ■ gitmap add ignore — merge curated .gitignore template block
        ignore/common  source=embed  (assets/ignore/common.gitignore)
        ignore/go      source=embed  (assets/ignore/go.gitignore)
      block tag: ignore/go

      created /home/me/code/my-go-svc/.gitignore (block: ignore/go)

      Next step: commit the updated .gitignore:
        git add .gitignore
        git commit -m "chore: refresh .gitignore via gitmap template"

### Example 2: Multi-language monorepo

    gitmap add ignore go node python

Writes a single marker block tagged `ignore/go+node+python` with the
union of all four templates (common + the three langs), deduped.

### Example 3: Re-run is a no-op

    gitmap add ignore go node
    gitmap add ignore go node     # → "unchanged"

### Example 4: Preview before committing

    gitmap add ignore rust --dry-run

Prints the exact block that would be written. Touches no files.

## Notes

- `common` is always first. Passing `common` explicitly is allowed but
  redundant — it's silently de-duped.
- Unknown languages abort with a clear error (`template ignore/foo: not
  found`) so typos surface immediately instead of silently dropping the
  language from the block.
- The block wraps the template body verbatim, including the
  `# source: ... # version: 1` audit-trail header. Future readers can
  tell which template version was applied without needing the gitmap
  binary on hand.

## See Also

- [add attributes](add-attributes.md) — Same flow for `.gitattributes`
- [add lfs-install](add-lfs-install.md) — LFS hooks + curated binary patterns
- [templates](templates.md) — Discover available templates / overlays

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter add-ignore
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
