# gitmap revert

Two modes:

1. **Version-tag mode** (legacy): revert the repo to a specific release tag.
2. **Transaction-journal mode**: replay the stored reverse-operation for
   any filesystem-mutating gitmap command (clone, mv, merge-*, etc.) by
   transaction id, the most recent one, or the last N.

## Alias

None

## Usage

    gitmap revert <version>                          # version-tag mode
    gitmap revert --list-txn                          # list recent transactions
    gitmap revert --show-txn <id>                     # inspect one transaction
    gitmap revert --txn <id> [--force]                # revert one transaction
    gitmap revert --last-txn [--force]                # revert the most recent
    gitmap revert --last-n-txn <N> [--force]          # revert N most recent
    gitmap revert --prune-txn                         # force a prune cycle

## Flags

| Flag | Purpose |
|------|---------|
| `--list-txn` | List the last 50 transactions and exit. |
| `--show-txn <id>` | Print one transaction (header + every captured file) by id. |
| `--txn <id>` | Revert the named transaction id. |
| `--last-txn` | Revert the most recent committed transaction. |
| `--last-n-txn <N>` | Revert the N most recent committed transactions, newest first. Stops on the first failure (already-reverted rows are preserved). |
| `--prune-txn` | Force a prune cycle now (drops everything beyond the 50-row cap). |
| `--force` | Skip the confirm prompt and skip backup-sha verification. |

## Prerequisites

- Must be inside a Git repository with release tags
- Run `gitmap list-versions` to see available versions (see list-versions.md)

## Examples

### Example 1: Revert to a specific version

    gitmap revert v2.20.0

**Output:**

    Current version: v2.22.0
    Reverting to v2.20.0...
    Checking out tag v2.20.0... done
    Rebuilding gitmap.exe... done
    Deploying to E:\bin-run\gitmap.exe... done
    ✓ Reverted to v2.20.0
    → Run 'gitmap version' to confirm

### Example 2: Revert to an older version

    gitmap revert v2.15.0

**Output:**

    Current version: v2.22.0
    Reverting to v2.15.0 (7 versions back)...
    Checking out tag v2.15.0... done
    Rebuilding gitmap.exe... done
    Deploying to E:\bin-run\gitmap.exe... done
    ✓ Reverted to v2.15.0

### Example 3: Version tag not found

    gitmap revert v9.9.9

**Output:**

    ✗ Error: tag v9.9.9 not found
    Available versions:
      v2.22.0, v2.21.0, v2.20.0, v2.19.0, ...
    → Use 'gitmap list-versions' to see all available tags

### Example 4: Undo the last 3 filesystem-mutating commands

    gitmap revert --last-n-txn 3

**Output:**

    About to revert 3 transaction(s), newest first:
      #42   mv         2026-05-06T14:22:01Z  rename "old/" ← "new/"
      #41   merge      2026-05-06T14:18:44Z  restore left + right pre-merge bytes
      #40   clone      2026-05-06T14:15:09Z  remove cloned repo at /work/foo
    Type 'yes' to continue: yes
      ✓ reverted transaction #42 (mv)
      ✓ reverted transaction #41 (merge)
      ✓ reverted transaction #40 (clone)
      ✓ reverted 3 transaction(s)

## See Also

- [list-versions](list-versions.md) — List available versions to revert to
- [release](release.md) — Create a new release
- [changelog](changelog.md) — View release notes before reverting
- [update](update.md) — Update to the latest version instead

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter revert
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
