# gitmap add attributes

Merge a curated `.gitattributes` template into a gitmap-managed marker
block inside `./.gitattributes`. Always merges the `common` template
(text eol normalization for known source extensions) plus any languages
you list on the command line. Idempotent — re-running with the same
language set is a byte-stable no-op.

## Usage

    gitmap add attributes [langs...] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run | false | Preview the merged `.gitattributes` block without writing anything |

## What it does

1. Verifies the current directory is inside a Git repository.
2. Resolves `common` plus each requested language (overlay > embed).
   Available languages: `go`, `node`, `python`, `rust`, `csharp`,
   `java`, `kotlin`, `php`, `ruby`, `swift` (plus anything you drop
   into `~/.gitmap/templates/attributes/`).
3. Concatenates bodies in stable order (`common` first, then your
   languages alphabetically), deduping repeated lines.
4. Writes a marker block at the repo root:

        # >>> gitmap:attributes/<sorted-langs> >>>
        ... merged body ...
        # <<< gitmap:attributes/<sorted-langs> <<<

5. Prints the merge outcome: `created`, `inserted block into`,
   `updated block in`, or `unchanged`.

## Idempotency contract

Same rules as `add ignore`: language list is sorted into the tag,
re-runs with the same set are no-ops, hand edits outside the block
survive, hand edits inside are overwritten.

## Examples

### Example 1: Node + TypeScript repo

    cd ~/code/my-ts-app
    gitmap add attributes node

**Output:**

      ■ gitmap add attributes — merge curated .gitattributes template block
        attributes/common  source=embed  (assets/attributes/common.gitattributes)
        attributes/node    source=embed  (assets/attributes/node.gitattributes)
      block tag: attributes/node

      created /home/me/code/my-ts-app/.gitattributes (block: attributes/node)

### Example 2: Pair with `add lfs-install`

    gitmap add attributes go
    gitmap add lfs-install

The two commands write **two separate** marker blocks
(`attributes/go` and `lfs/common`) into the same `.gitattributes` so
they can be refreshed independently.

### Example 3: Preview before committing

    gitmap add attributes rust --dry-run

Prints the exact block that would be written. Touches no files.

## Notes

- The `common` template asserts `text eol=lf` for ~30 source
  extensions and `binary` for common image / archive types — but
  treats `*.svg` as **text**, not LFS, since SVGs diff and merge well
  as text.
- `add attributes` does **not** install Git LFS. Use
  [`add lfs-install`](add-lfs-install.md) when you also want LFS
  hooks + the binary-pattern block.

## See Also

- [add ignore](add-ignore.md) — Same flow for `.gitignore`
- [add lfs-install](add-lfs-install.md) — LFS hooks + curated binary patterns
- [templates](templates.md) — Discover available templates / overlays

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter add-attributes
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
