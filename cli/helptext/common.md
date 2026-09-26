# gitmap common

Union-merge curated defaults for `.gitignore`, `.gitattributes`,
`.prettierignore`, and `.prettierrc` into the current directory
**without** using gitmap marker blocks. Existing content stays
untouched; only missing entries are appended. Bare invocation defaults
to running all targets in sequence.

## Aliases

sync, sy, cm

## Usage

    gitmap common [target] [--dry-run] [--force]

## Targets

| Target | File / Action | Merge strategy |
|--------|---------------|----------------|
| `all` (default) | all targets below | Runs each target in order |
| `ignore` | `.gitignore` | Line union (trim-compare) |
| `attributes` | `.gitattributes` | Line union (trim-compare) |
| `lfs-install` | `git lfs install --local` + `.gitattributes` | Git LFS tracking configuration |
| `prettier-ignore` | `.prettierignore` | Line union (trim-compare) |
| `prettier-rc` | `.prettierrc` | JSON key union (existing keys win) |

## Flags

| Flag | Purpose |
|------|---------|
| `--dry-run`, `-n` | Print planned additions without touching disk |
| `--force`, `-f` | For `prettier-rc` only: overwrite conflicting keys |
| `--help`, `-h` | Display this help reference |

## Examples

```bash
# Run full baseline sync across all files
gitmap common

# Sync only gitignore defaults
gitmap common ignore

# Dry run preview
gitmap common all --dry-run
```

