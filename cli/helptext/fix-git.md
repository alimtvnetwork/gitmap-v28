# gitmap fix-git

Diagnose and repair Git repository commit failures, stale lockfiles,
read-only permission flags, Windows NTFS ACLs, and corrupted index files.

## Synopsis

```
gitmap fix-git [path] [--dry-run] [--permissions] [--locks] [--index] [--safe-dir] [--json]
gitmap --fix-git                                         # root-level flag alias
gitmap fg                                                # short alias
```

## Behavior

1. **Permissions & Windows ACLs:** Automatically grants Full Control to the active
   user on `.git` (`icacls` on Windows, `chmod -R u+rwX` on POSIX) and removes
   read-only file attributes (`attrib -r`).
2. **Stale Lockfile Removal:** Recursively locates and cleans orphaned `.lock` files
   (`.git/index.lock`, `HEAD.lock`, `config.lock`, etc.).
3. **Index Rebuild:** Detects 0-byte or corrupted `.git/index` files, creates a
   timestamped backup, and safely reconstructs the index from `HEAD` using `git reset`.
4. **Safe Directory Registration:** Resolves `fatal: detected dubious ownership`
   errors by registering the repository in global Git configuration.

## Examples

```
$ gitmap fix-git
  [FIX-GIT] Git Diagnostic & Self-Healing Engine
  Repository: D:\projects\my-app

  ✓ [Permissions  ] Read-only or restricted ACL permissions detected on .git
     ↳ Action: Granted Full Control to current user and stripped read-only file attributes
  ✓ [Lockfile     ] Stale lock file found: .git/index.lock
     ↳ Action: Removed stale lock file .git/index.lock

  ✓ Successfully remediated 2 Git issue(s)!
  Repository ready for staging and committing.
```

## Flags

- `-n`, `--dry-run`: Scan and report detected issues without modifying any files.
- `--permissions`, `--perms`: Only repair file attributes and OS user permissions.
- `--locks`, `--locks-only`: Only clean orphaned `.lock` files.
- `--index`, `--index-only`: Only inspect and repair corrupted index files.
- `--safe-dir`: Only register repository in global `safe.directory`.
- `--json`: Output diagnostic and remediation results in structured JSON format.
