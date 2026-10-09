# gitmap lowercase (Safe 2-Step Renamer)

Scan repository or directory for uppercase or mixed-case files, safely rename them to lowercase using an atomic two-step move to avoid filesystem collisions on case-insensitive filesystems (Windows NTFS, macOS APFS), and optionally stage and commit the renames in Git.

## Syntax

```bash
gitmap lowercase [patterns...] [flags]
```

## Aliases

`lcf`, `lower-case-fix`, `lowercase-fix`, `lower`, `lc-fix`

## Wildcard & Filter Engine

- `*` or no pattern: Matches all uppercase files across the repository or directory.
- `prefix*` (e.g. `SKILL*`, `TEST*`, `READ*`): Matches any uppercase file starting with `prefix` (ending wildcard).
- `*suffix` or `*.ext` (e.g. `*.md`, `*doc`, `*.txt`): Matches uppercase files with the specified extension or ending.
- `*pattern*` (e.g. `*readme*`): Matches uppercase files containing the substring.
- `docs/*.md`: Scans within a specific relative subdirectory.

## The Two-Step Rename Architecture

On case-insensitive filesystems like Windows (NTFS) and macOS (APFS), running a direct rename such as `git mv README.md readme.md` fails, silently no-ops, or produces Git index corruption because the OS filesystem considers `README.md` and `readme.md` to be the exact same file path.

GitMap solves this with an atomic two-step rename pipeline:

1. **Step 1 (Safe Intermediate Move)**:
   - Git repository: `git mv <FILE> <FILE>.tmp-lcf`
   - Non-Git directory: `os.Rename(<FILE>, <FILE>.tmp-lcf)`
2. **Step 2 (Target Rename)**:
   - Git repository: `git mv <FILE>.tmp-lcf <file>`
   - Non-Git directory: `os.Rename(<FILE>.tmp-lcf, <file>)`
3. **Step 3 (Pre-Flight Verification & User Confirmation)**:
   - When running on a Git repository, GitMap displays a pre-flight event summarizing matched files, the 2-step move plan, and after-effects (Git index staging and automated commit).
   - Prompts the user to enter `confirm`, `yes`, or `y` to proceed. Use `-y` or `--yes` to proceed without interactive confirmation.
4. **Step 4 (Git Synchronization)**:
   - Staged with `git add -A` to guarantee index tree consistency.
5. **Step 5 (Atomic Commit)**:
   - Automatically commits changes with a standardized message (can be disabled with `--no-commit`).

## Flags & Options

- `-d`, `--dry-run`: Preview matching files and planned rename operations without modifying files on disk or in Git.
- `--no-commit`: Perform disk renames and stage in Git index, but skip creating an automated Git commit.
- `-r`, `--readme`: Restrict scan to root-level README files only (`README.md`, `README`, etc.).
- `-m`, `--message <text>`: Provide a custom Git commit message for the automated commit.
- `-y`, `--yes`: Proceed non-interactively without prompt confirmation.

## Usage Examples

```bash
# Rename all uppercase files across the repository to lowercase
gitmap lowercase

# Rename uppercase markdown files only
gitmap lowercase *.md

# Rename files starting with SKILL (e.g. SKILL.md -> skill.md)
gitmap lowercase SKILL*

# Target root README only
gitmap lowercase --readme

# Preview renames without modifying files
gitmap lowercase --dry-run

# Rename and stage in Git index without committing
gitmap lowercase --no-commit

# Custom commit message
gitmap lowercase *.md -m "chore: normalize markdown filenames to lowercase"
```

## Manipulation Summary Reporting

At the conclusion of every run, GitMap outputs a structured summary box:
- Total files scanned in directory tree
- Files matched by filter
- Files successfully renamed
- Git commit SHA / Staged status / Filesystem status
- Explicit breakdown of the 2-step manipulation steps executed
