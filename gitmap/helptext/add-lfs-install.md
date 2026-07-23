# gitmap add lfs-install

Install Git LFS hooks in the current repository **and** merge the curated
`lfs/common` `.gitattributes` template into a gitmap-managed marker
block. Idempotent — re-running is a byte-stable no-op when the template
hasn't changed.

## Usage

    gitmap add lfs-install [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run | false | Preview the merged `.gitattributes` block without modifying anything |

## What it does

1. Verifies the current directory is inside a Git repository.
2. Verifies `git lfs` is installed and on PATH.
3. Resolves the `lfs/common` template — overlay (`~/.gitmap/templates/lfs/common.gitattributes`)
   wins over the embedded copy.
4. Runs `git lfs install --local` (Git LFS itself is idempotent).
5. Merges the resolved template body into `.gitattributes` at the repo
   root, wrapped in a marker block:

        # >>> gitmap:lfs/common >>>
        ... template body, verbatim ...
        # <<< gitmap:lfs/common <<<

6. Prints the merge outcome: `created`, `inserted block into`,
   `updated block in`, or `unchanged`.

## Idempotency contract

| Run | Pre-state of `.gitattributes` | Outcome |
|-----|-------------------------------|---------|
| 1   | missing                       | `created` (whole file is the block) |
| 1   | exists, no marker block       | `inserted block into` (block appended) |
| 2+  | exists, marker block present, body identical | `unchanged` (file untouched) |
| 2+  | exists, marker block present, body changed   | `updated block in` (only the block rewritten) |

Hand edits **outside** the marker block survive every re-run. Hand edits
**inside** the block are intentionally overwritten — fork the template
to `~/.gitmap/templates/lfs/common.gitattributes` if you want different
contents to stick.

## Examples

### Example 1: First-time install in a new repo

    cd ~/code/my-design-assets
    gitmap add lfs-install

**Output:**

      ■ gitmap add lfs-install — LFS hooks + templated .gitattributes block
      template source: embed (assets/lfs/common.gitattributes)

      ✓ git lfs install --local ran successfully
      created /home/me/code/my-design-assets/.gitattributes (block: lfs/common)

      Next step: commit the updated .gitattributes:
        git add .gitattributes
        git commit -m "chore: install Git LFS + track common binaries via gitmap template"

### Example 2: Re-run after the template was bumped

    gitmap add lfs-install

**Output:**

      ■ gitmap add lfs-install — LFS hooks + templated .gitattributes block
      template source: embed (assets/lfs/common.gitattributes)

      ✓ git lfs install --local ran successfully
      updated block in /home/me/code/my-design-assets/.gitattributes (block: lfs/common)

### Example 3: Re-run with no changes (the idempotent case)

    gitmap add lfs-install

**Output:**

      ■ gitmap add lfs-install — LFS hooks + templated .gitattributes block
      template source: embed (assets/lfs/common.gitattributes)

      ✓ git lfs install --local ran successfully
      unchanged /home/me/code/my-design-assets/.gitattributes (block: lfs/common)

### Example 4: Preview before committing to the change

    gitmap add lfs-install --dry-run

**Output:** the exact block that would be written, printed to stdout.
No file on disk is touched and `git lfs install` is **not** run.

## Notes

- This command does **not** migrate already-committed blobs to LFS. Use
  `git lfs migrate import` for that.
- The block wraps the template body verbatim, including its
  `# source:` / `# version:` audit-trail header. That header lives in
  the `.gitattributes` so future readers can tell which template version
  was applied without needing the gitmap binary on hand.
- Sibling command [`lfs-common`](lfs-common.md) shells out to
  `git lfs track` per pattern instead of templating. Use `lfs-common`
  when you want Git LFS itself to author the lines; use
  `add lfs-install` when you want gitmap's curated, versioned set.

## See Also

- [lfs-common](lfs-common.md) — Per-pattern `git lfs track` flow
- [templates list](commands.md) — Discover available templates
- [templates show](commands.md) — Print a template to stdout
- [setup](setup.md) — Configure Git diff/merge tools and aliases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter add-lfs-install
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
