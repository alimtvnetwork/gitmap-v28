# `gitmap fix-git`

Diagnose and repair Git repository commit failures, stale lockfiles, read-only permission flags, Windows NTFS ACLs, and corrupted index files.

## Summary

When Git operations fail with errors like:
- `fatal: Could not write new index file.`
- `fatal: Unable to create '.git/index.lock': File exists.`
- `fatal: index file corrupt`
- `fatal: detected dubious ownership in repository at ...`

`gitmap fix-git` provides an automated self-healing pipeline that repairs ownership, strips read-only flags, clears abandoned lock files, and rebuilds damaged indexes without risking uncommitted changes.

## Syntax

```bash
gitmap fix-git [path] [flags]
gitmap --fix-git
gitmap fg
```

## Options & Flags

| Flag | Shorthand | Description |
| :--- | :--- | :--- |
| `--dry-run` | `-n` | Scan and diagnose without applying file changes |
| `--permissions` | `--perms` | Only fix user permissions, ACLs, and read-only attributes |
| `--locks` | `--locks-only` | Only clean orphaned `.lock` files |
| `--index` | `--index-only` | Only inspect and rebuild damaged index files |
| `--safe-dir` | | Only register the repository under `git config --global safe.directory` |
| `--json` | | Output report as structured JSON |
| `--verbose` | `-v` | Enable detailed diagnostic output |

## Common Scenarios & Fixes

### 1. `fatal: Could not write new index file.`
**Cause:** The current user lacks Full Control NTFS permissions on `.git/` or `.git/index` has the read-only attribute set.
**Solution:** Run `gitmap fix-git`. It grants Full Control to `%USERNAME%` via `icacls` and removes read-only attributes with `attrib -r`.

### 2. `fatal: Unable to create '.git/index.lock': File exists.`
**Cause:** A background tool (e.g. GitHub Desktop, VS Code, indexer) terminated abruptly while holding a lockfile.
**Solution:** Run `gitmap fix-git`. It safely removes orphaned `.lock` files.

### 3. Corrupt or 0-byte `.git/index`
**Cause:** Crash during an index write.
**Solution:** Run `gitmap fix-git`. It backs up the corrupt index to `.git/index.corrupt.<ts>` and cleanly rebuilds it using `git reset HEAD`, preserving all modified files in your working directory.
