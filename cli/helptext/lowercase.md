# `gitmap lowercase` (aliases: `lower-case-fix`, `lowercase-fix`, `lcf`, `lower`, `lowercase-readme`, `readme-lower`, `readme-lowercase`, `lcr`)

Scan the repository or directory for uppercase or mixed-case files matching specified patterns, safely rename them to lowercase using a two-step git mv to avoid case-collision on case-insensitive filesystems (Windows/macOS), and optionally commit the changes.

## Usage

```bash
# Rename all markdown files matching *.md across the repo (default)
gitmap lowercase

# Rename files matching specific globs (e.g. *.md, *md, *, SKILL*, README*)
gitmap lowercase *.md
gitmap lowercase *md
gitmap lowercase *
gitmap lowercase SKILL*
gitmap lowercase "docs/*.md"

# Target root README specifically (README.md, README, Readme.md -> readme.md)
gitmap lowercase-readme
gitmap readme-lower
gitmap lowercase --readme

# Preview matching files without modifying disk or Git
gitmap lowercase *.md --dry-run

# Rename without creating an automated Git commit
gitmap lowercase *.md --no-commit

# Custom commit message
gitmap lowercase *.md -m "chore: lowercase markdown files"
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-d`, `--dry-run` | Preview matching uppercase files without modifying disk or Git | `false` |
| `--no-commit` | Rename files on disk and in Git index without creating a commit | `false` (auto-commits) |
| `-r`, `--readme` | Target root README files only (`README.md`, `README`, etc.) | `false` |
| `-m`, `--message <text>` | Custom Git commit message | Auto-generated |
| `-y`, `--yes` | Proceed without interactive confirmation | `false` |

## Features & Git Manipulation
- **Two-Step Git Move**: On Windows and macOS, Git cannot directly rename `README.md` to `readme.md` because the filesystem is case-insensitive. `gitmap lowercase` performs a safe 2-step move:
  1. `git mv FILE FILE.tmp-lcf`
  2. `git mv FILE.tmp-lcf file`
- **Non-Git & Untracked Support**: If the directory is not a Git repository, or the file is untracked, a safe 2-step filesystem rename is performed. If in a Git repo, untracked files are staged with `git add`.
- **Status Reporting**: Displays a step-by-step log of each file's rename operation and outputs an informative summary of scanned, matched, and renamed files.
