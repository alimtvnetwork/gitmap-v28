# gitmap github-desktop

Register git repositories with GitHub Desktop.

> **`github-desktop` (`gd`) and `desktop-sync` (`ds`) are the same command.**
> Use whichever you remember.

## Alias

gd (also: ds, desktop-sync)

## Usage

    gitmap gd                       # register CWD or every DB-tracked repo under CWD
    gitmap ds                       # same thing
    gitmap gd D:\path\to\repo       # register an explicit folder
    gitmap gd --all                 # register every repo in the gitmap DB

## Flags

    --all      Register every repo currently tracked in the gitmap database,
               regardless of where you ran the command from.

## Prerequisites

- GitHub Desktop installed with the `github` CLI on PATH.
- A git repo at CWD, an explicit path, or under a registered scan root.
- **No prior `gitmap scan` is required.**

## Resolution order (no args)

1. Is GitHub Desktop installed? If not → exit with install hint.
2. Is CWD itself a git repo (`.git` dir OR `.git` file for worktrees)? Register it.
3. Is CWD inside a registered scan root? Bulk-register every tracked repo under it.
4. Otherwise → friendly hint + exit 3.

## Examples

### Example 1: Register the current repo

    cd D:\wp-work\riseup-asia\macro-ahk
    gitmap gd

**Output:**

    ✓ Registered: macro-ahk
    GitHub Desktop: 1 added · 0 skipped · 0 failed

### Example 2: Register an explicit path

    gitmap github-desktop D:\projects\billing-svc

**Output:**

    ✓ Registered: billing-svc
    GitHub Desktop: 1 added · 0 skipped · 0 failed

### Example 3: Bulk-register under a scan root

    cd D:\wp-work
    gitmap gd

**Output:**

    [1/14] my-api ............ ✓
    [2/14] web-app ........... ✓
    ...
    GitHub Desktop: 14 added · 0 skipped · 0 failed

## See Also

- [desktop-sync](desktop-sync.md) — same command, different name
- [scan](scan.md) — `--github-desktop` registers during scan
- [clone](clone.md) — `--github-desktop` registers during clone
- [scan-gd (spec 102)](../spec/01-app/102-scan-gd.md) — design doc for bulk mode

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter github-desktop
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
