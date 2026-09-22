# gitmap lowercase-readme (Root README Renamer)

Safely rename root-level README files (e.g. `README.md`, `README`, `README.txt`, `Readme.md` -> `readme.md`, `readme`) to lowercase via a safe two-step Git move.

## Syntax

```bash
gitmap lowercase-readme [flags]
```

## Aliases

`readme-lower`, `readme-lowercase`, `lcr`, `lower-case-readme`

## Scope & Target Isolation

- Strictly isolates the root directory (`.`): only targets files at the project root matching `readme*`.
- Ignores all subdirectories (e.g. `docs/README.md`, `sub/README.md` are left untouched).
- Employs the two-step move (`README.md` -> `README.md.tmp-lcf` -> `readme.md`) to avoid case-insensitive collision on Windows NTFS and macOS APFS.

## Flags & Options

- `-d`, `--dry-run`: Preview root README rename without modifying disk or Git.
- `--no-commit`: Rename and stage in Git index without creating an automated Git commit.
- `-m`, `--message <text>`: Custom commit message for the automated commit.
- `-y`, `--yes`: Proceed non-interactively without prompt confirmation.

## Usage Examples

```bash
# Rename root README.md to readme.md and commit to Git
gitmap lowercase-readme

# Shorthand alias
gitmap readme-lower

# Preview without modifying disk
gitmap lowercase-readme --dry-run

# Rename and stage in Git index without committing
gitmap lowercase-readme --no-commit
```
