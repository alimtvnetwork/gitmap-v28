# CI/CD Fix: Release Script Unicode & Regex Corruption

## Error Summary
The `bump_versions.py` release script failed with a `UnicodeDecodeError: 'charmap' codec can't decode byte 0x90 in position 2851: character maps to <undefined>` on Windows runners. 
Additionally, a logic bug in the script caused it to greedily overwrite *every* historical version string in `readme.md`, causing massive documentation corruption.

## Root Cause Analysis (RCA)
- **Unicode Error**: Python's `open()` function defaults to the system's local encoding (e.g., `cp1252` on Windows). Because `readme.md` contains UTF-8 emoji and special characters, reading it with `cp1252` raised a decode exception.
- **Regex Corruption**: The script naively executed `re.sub(r'v\d+\.\d+\.\d+', f'v{new_version}', content)` across the entire `readme.md` file. This obliterated all historical version references, CLI examples, and changelog headers, replacing them all with the current version.

## Solution Applied
- Enforced `encoding="utf-8"` explicitly on all `open()` file handlers in `.lovable/release/bump_versions.py`.
- Restructured the regex replacement for `readme.md` to strictly target the exact header: `re.sub(r'Pinned version: v\d+\.\d+\.\d+', f'Pinned version: v{new_version}', content)`.
- Reverted the corrupted `readme.md` via `git checkout readme.md`.

## What NOT to Repeat
- NEVER open files in Python scripts (especially for cross-platform repositories) without explicitly specifying `encoding="utf-8"`.
- NEVER use global or greedy regex replacements on documentation files (like `readme.md` or `changelog.md`) when bumping versions; always target a highly specific anchor (e.g., `Pinned version: vX.Y.Z`).
