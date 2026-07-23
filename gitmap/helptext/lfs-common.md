# gitmap lfs-common

Track a curated set of common binary file types with Git LFS in the
current repository.

## Alias

lfsc

## Usage

    gitmap lfs-common [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run | false | Preview what would be tracked without modifying anything |

## What it does

1. Verifies the current directory is inside a Git repository.
2. Verifies `git lfs` is installed and on PATH.
3. Runs `git lfs install --local` (idempotent).
4. Calls `git lfs track "<pattern>"` for each entry below, which appends
   the canonical line `<pattern> filter=lfs diff=lfs merge=lfs -text`
   to `.gitattributes` (skipping patterns already tracked).
5. Prints a per-pattern summary: added, already tracked, or failed.

## Default tracked patterns

    *.pptx    *.ppt     *.eps     *.psd
    *.ttf     *.wott    *.svg     *.ai
    *.jpg     *.bmp     *.png
    *.zip     *.gz      *.tar     *.rar    *.7z
    *.mp4     *.aep

## After running

Commit the updated `.gitattributes` so collaborators pick up the LFS rules:

    git add .gitattributes
    git commit -m "chore: track common binary types with Git LFS"

## Notes

- Existing files already committed as plain Git blobs are NOT migrated to
  LFS by this command. Use `git lfs migrate import` for that.
- The command is idempotent: re-running adds only new patterns and
  reports the rest as `already tracked`.

## See Also

- [setup](setup.md) — Configure Git diff/merge tool, aliases & core settings
- [doctor](doctor.md) — Diagnose PATH, deploy, and version issues

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter lfs-common
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap lfs-common
```
